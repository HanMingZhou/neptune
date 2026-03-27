#!/usr/bin/env python3
"""
Horovod MPI 分布式训练脚本示例
适用于 Neptune 平台的 MPI 模式

使用方法:
1. 在创建训练任务时选择 MPI 框架
2. 启动命令填写: 
   mpirun --allow-run-as-root -np $WORLD_SIZE \
          -H $MPI_HOSTS \
          python /path/to/horovod_train.py --epochs 10

   或更简洁 (MPI_HOSTS 已由平台注入):
   mpirun -np 4 -H worker-0:2,worker-1:2 python /path/to/horovod_train.py

特点:
- 使用 Horovod 实现 MPI 分布式训练
- 支持多种后端 (NCCL, Gloo, MPI)
- 兼容 PyTorch, TensorFlow, Keras
- 梯度压缩和通信优化
"""

import argparse
import os
import time
import torch
import torch.nn as nn
import torch.optim as optim
from torch.utils.data import DataLoader, TensorDataset

# 尝试导入 Horovod
try:
    import horovod.torch as hvd
    HAS_HOROVOD = True
except ImportError:
    HAS_HOROVOD = False
    print('[WARN] Horovod 未安装，将以单进程模式运行')
    print('       安装方法: pip install horovod')


class SimpleNet(nn.Module):
    """简单的神经网络模型"""
    def __init__(self, input_size=784, hidden_size=256, num_classes=10):
        super(SimpleNet, self).__init__()
        self.fc1 = nn.Linear(input_size, hidden_size)
        self.relu = nn.ReLU()
        self.fc2 = nn.Linear(hidden_size, hidden_size)
        self.fc3 = nn.Linear(hidden_size, num_classes)
        self.dropout = nn.Dropout(0.2)
        
    def forward(self, x):
        x = x.view(x.size(0), -1)
        x = self.fc1(x)
        x = self.relu(x)
        x = self.dropout(x)
        x = self.fc2(x)
        x = self.relu(x)
        x = self.fc3(x)
        return x


def create_dummy_data(num_samples=10000, input_size=784, num_classes=10):
    """创建模拟数据集"""
    X = torch.randn(num_samples, input_size)
    y = torch.randint(0, num_classes, (num_samples,))
    return TensorDataset(X, y)


def get_rank_info():
    """获取进程 Rank 信息"""
    if HAS_HOROVOD:
        return hvd.rank(), hvd.size(), hvd.local_rank()
    else:
        return 0, 1, 0


def setup_horovod():
    """初始化 Horovod"""
    if not HAS_HOROVOD:
        return 0, 1, 0, torch.device('cuda' if torch.cuda.is_available() else 'cpu')
    
    # 初始化 Horovod
    hvd.init()
    
    rank = hvd.rank()
    world_size = hvd.size()
    local_rank = hvd.local_rank()
    
    print(f'[Rank {rank}] Horovod 初始化完成')
    print(f'  - Rank: {rank}')
    print(f'  - World Size: {world_size}')
    print(f'  - Local Rank: {local_rank}')
    
    # 绑定 GPU
    if torch.cuda.is_available():
        torch.cuda.set_device(local_rank)
        device = torch.device(f'cuda:{local_rank}')
        print(f'  - 使用 GPU: {torch.cuda.get_device_name(local_rank)}')
    else:
        device = torch.device('cpu')
        print(f'  - 使用 CPU')
    
    return rank, world_size, local_rank, device


def train_epoch(model, train_loader, criterion, optimizer, device, epoch, rank, world_size):
    """训练一个 epoch"""
    model.train()
    total_loss = 0.0
    correct = 0
    total = 0
    
    for batch_idx, (data, target) in enumerate(train_loader):
        data, target = data.to(device), target.to(device)
        
        optimizer.zero_grad()
        output = model(data)
        loss = criterion(output, target)
        loss.backward()
        optimizer.step()
        
        total_loss += loss.item()
        _, predicted = output.max(1)
        total += target.size(0)
        correct += predicted.eq(target).sum().item()
        
        if batch_idx % 50 == 0 and rank == 0:
            print(f'Epoch: {epoch} | Batch: {batch_idx}/{len(train_loader)} | '
                  f'Loss: {loss.item():.4f} | Acc: {100.*correct/total:.2f}%')
    
    # 聚合所有进程的结果
    if HAS_HOROVOD:
        total_loss = hvd.allreduce(torch.tensor([total_loss]), average=True).item()
        correct = hvd.allreduce(torch.tensor([correct]), average=False).item()
        total = hvd.allreduce(torch.tensor([total]), average=False).item()
    
    avg_loss = total_loss / len(train_loader)
    accuracy = 100. * correct / total
    
    return avg_loss, accuracy


def validate(model, val_loader, criterion, device):
    """验证模型"""
    model.eval()
    total_loss = 0.0
    correct = 0
    total = 0
    
    with torch.no_grad():
        for data, target in val_loader:
            data, target = data.to(device), target.to(device)
            output = model(data)
            loss = criterion(output, target)
            
            total_loss += loss.item()
            _, predicted = output.max(1)
            total += target.size(0)
            correct += predicted.eq(target).sum().item()
    
    # 聚合所有进程的结果
    if HAS_HOROVOD:
        total_loss = hvd.allreduce(torch.tensor([total_loss]), average=True).item()
        correct = hvd.allreduce(torch.tensor([correct]), average=False).item()
        total = hvd.allreduce(torch.tensor([total]), average=False).item()
    
    avg_loss = total_loss / len(val_loader)
    accuracy = 100. * correct / total
    
    return avg_loss, accuracy


def main():
    parser = argparse.ArgumentParser(description='Horovod MPI 分布式训练示例')
    parser.add_argument('--epochs', type=int, default=10, help='训练轮数')
    parser.add_argument('--batch-size', type=int, default=64, help='每个进程的批次大小')
    parser.add_argument('--lr', type=float, default=0.001, help='基础学习率')
    parser.add_argument('--num-samples', type=int, default=10000, help='样本数量')
    parser.add_argument('--save-dir', type=str, default='/workspace/output', help='模型保存目录')
    parser.add_argument('--log-dir', type=str, default='/workspace/logs', help='TensorBoard日志目录')
    parser.add_argument('--fp16', action='store_true', help='使用混合精度训练')
    parser.add_argument('--compression', type=str, default='none', 
                        choices=['none', 'fp16'], help='梯度压缩方式')
    args = parser.parse_args()
    
    # 初始化 Horovod
    rank, world_size, local_rank, device = setup_horovod()
    
    # 调整学习率 (Horovod 推荐: lr * world_size)
    scaled_lr = args.lr * world_size
    
    # 只在 Rank 0 打印配置信息
    if rank == 0:
        print(f'\n[CONFIG] 训练配置:')
        print(f'  - Epochs: {args.epochs}')
        print(f'  - Batch Size (per GPU): {args.batch_size}')
        print(f'  - Global Batch Size: {args.batch_size * world_size}')
        print(f'  - Base Learning Rate: {args.lr}')
        print(f'  - Scaled Learning Rate: {scaled_lr}')
        print(f'  - World Size: {world_size}')
        print(f'  - FP16: {args.fp16}')
        print(f'  - Gradient Compression: {args.compression}')
    
    # 创建数据集
    if rank == 0:
        print(f'\n[INFO] 创建模拟数据集，样本数: {args.num_samples}')
    
    dataset = create_dummy_data(num_samples=args.num_samples)
    train_size = int(0.8 * len(dataset))
    val_size = len(dataset) - train_size
    train_dataset, val_dataset = torch.utils.data.random_split(dataset, [train_size, val_size])
    
    # 使用 Horovod DistributedSampler
    if HAS_HOROVOD:
        train_sampler = torch.utils.data.distributed.DistributedSampler(
            train_dataset, num_replicas=world_size, rank=rank, shuffle=True
        )
        val_sampler = torch.utils.data.distributed.DistributedSampler(
            val_dataset, num_replicas=world_size, rank=rank, shuffle=False
        )
    else:
        train_sampler = None
        val_sampler = None
    
    train_loader = DataLoader(
        train_dataset, 
        batch_size=args.batch_size, 
        sampler=train_sampler,
        shuffle=(train_sampler is None),
        num_workers=2,
        pin_memory=True
    )
    val_loader = DataLoader(
        val_dataset, 
        batch_size=args.batch_size, 
        sampler=val_sampler,
        num_workers=2,
        pin_memory=True
    )
    
    # 创建模型
    model = SimpleNet().to(device)
    
    if rank == 0:
        print(f'[INFO] 模型已创建')
    
    # 损失函数和优化器
    criterion = nn.CrossEntropyLoss()
    optimizer = optim.Adam(model.parameters(), lr=scaled_lr)
    
    # Horovod: 包装优化器
    if HAS_HOROVOD:
        # 选择梯度压缩算法
        if args.compression == 'fp16':
            compression = hvd.Compression.fp16
        else:
            compression = hvd.Compression.none
        
        # 包装优化器
        optimizer = hvd.DistributedOptimizer(
            optimizer,
            named_parameters=model.named_parameters(),
            compression=compression,
            op=hvd.Average
        )
        
        # 广播初始参数
        hvd.broadcast_parameters(model.state_dict(), root_rank=0)
        hvd.broadcast_optimizer_state(optimizer, root_rank=0)
        
        if rank == 0:
            print(f'[INFO] Horovod 优化器已配置')
    
    scheduler = optim.lr_scheduler.StepLR(optimizer, step_size=5, gamma=0.5)
    
    # TensorBoard (仅 Rank 0)
    writer = None
    if rank == 0:
        try:
            from torch.utils.tensorboard import SummaryWriter
            os.makedirs(args.log_dir, exist_ok=True)
            writer = SummaryWriter(args.log_dir)
            print(f'[INFO] TensorBoard 日志目录: {args.log_dir}')
        except ImportError:
            print('[WARN] 未安装 tensorboard，跳过日志记录')
    
    # 混合精度训练
    scaler = None
    if args.fp16 and torch.cuda.is_available():
        scaler = torch.cuda.amp.GradScaler()
        if rank == 0:
            print('[INFO] 已启用混合精度训练 (FP16)')
    
    # 训练循环
    if rank == 0:
        print(f'\n[INFO] 开始训练，共 {args.epochs} 轮')
        print('=' * 60)
    
    best_acc = 0.0
    start_time = time.time()
    
    for epoch in range(1, args.epochs + 1):
        epoch_start = time.time()
        
        # 设置 epoch
        if HAS_HOROVOD and train_sampler:
            train_sampler.set_epoch(epoch)
        
        # 训练
        train_loss, train_acc = train_epoch(
            model, train_loader, criterion, optimizer, device, epoch, rank, world_size
        )
        
        # 验证
        val_loss, val_acc = validate(model, val_loader, criterion, device)
        
        # 更新学习率
        scheduler.step()
        
        epoch_time = time.time() - epoch_start
        
        # 只在 Rank 0 打印日志和保存模型
        if rank == 0:
            print(f'\nEpoch {epoch}/{args.epochs} 完成 (用时: {epoch_time:.1f}s)')
            print(f'  Train - Loss: {train_loss:.4f}, Acc: {train_acc:.2f}%')
            print(f'  Val   - Loss: {val_loss:.4f}, Acc: {val_acc:.2f}%')
            
            # TensorBoard 记录
            if writer:
                writer.add_scalar('Loss/train', train_loss, epoch)
                writer.add_scalar('Loss/val', val_loss, epoch)
                writer.add_scalar('Accuracy/train', train_acc, epoch)
                writer.add_scalar('Accuracy/val', val_acc, epoch)
                writer.add_scalar('LearningRate', scheduler.get_last_lr()[0], epoch)
            
            # 保存最佳模型
            if val_acc > best_acc:
                best_acc = val_acc
                os.makedirs(args.save_dir, exist_ok=True)
                model_path = os.path.join(args.save_dir, 'best_model.pth')
                torch.save({
                    'epoch': epoch,
                    'model_state_dict': model.state_dict(),
                    'optimizer_state_dict': optimizer.state_dict(),
                    'best_acc': best_acc,
                }, model_path)
                print(f'  [SAVE] 最佳模型已保存: {model_path}')
    
    # 训练完成
    if rank == 0:
        total_time = time.time() - start_time
        print('=' * 60)
        print(f'\n[INFO] 训练完成!')
        print(f'  总用时: {total_time:.1f}s')
        print(f'  最佳验证准确率: {best_acc:.2f}%')
        
        if writer:
            writer.close()


if __name__ == '__main__':
    main()

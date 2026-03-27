#!/usr/bin/env python3
"""
PyTorch DDP 分布式训练脚本示例
适用于 Neptune 平台的 PYTORCH_DDP 模式

使用方法:
1. 在创建训练任务时选择 PYTORCH_DDP 框架
2. 启动命令填写: 
   torchrun --nnodes=$WORLD_SIZE --nproc_per_node=1 \
            --rdzv_id=$JOB_ID --rdzv_backend=c10d \
            --rdzv_endpoint=$MASTER_ADDR:$MASTER_PORT \
            /path/to/ddp_train.py --epochs 10

   或简化版 (Volcano pytorch 插件会自动注入环境变量):
   python /path/to/ddp_train.py --epochs 10

特点:
- 多机多卡分布式数据并行
- 自动从环境变量获取分布式配置
- 支持梯度同步和模型同步
- 支持 TensorBoard 日志（仅 Rank 0 写入）
"""

import argparse
import os
import time
import torch
import torch.nn as nn
import torch.optim as optim
import torch.distributed as dist
from torch.nn.parallel import DistributedDataParallel as DDP
from torch.utils.data import DataLoader, TensorDataset
from torch.utils.data.distributed import DistributedSampler


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


def setup_distributed():
    """初始化分布式环境"""
    # 从环境变量获取分布式配置
    # Volcano pytorch 插件会自动注入这些变量
    rank = int(os.environ.get('RANK', 0))
    world_size = int(os.environ.get('WORLD_SIZE', 1))
    local_rank = int(os.environ.get('LOCAL_RANK', 0))
    master_addr = os.environ.get('MASTER_ADDR', 'localhost')
    master_port = os.environ.get('MASTER_PORT', '29500')
    
    print(f'[Rank {rank}] 初始化分布式环境...')
    print(f'  - RANK: {rank}')
    print(f'  - WORLD_SIZE: {world_size}')
    print(f'  - LOCAL_RANK: {local_rank}')
    print(f'  - MASTER_ADDR: {master_addr}')
    print(f'  - MASTER_PORT: {master_port}')
    
    # 设置当前进程使用的 GPU
    if torch.cuda.is_available():
        torch.cuda.set_device(local_rank)
        device = torch.device(f'cuda:{local_rank}')
        backend = 'nccl'  # GPU 使用 NCCL 后端
    else:
        device = torch.device('cpu')
        backend = 'gloo'  # CPU 使用 Gloo 后端
    
    # 初始化进程组
    if not dist.is_initialized():
        dist.init_process_group(
            backend=backend,
            init_method=f'tcp://{master_addr}:{master_port}',
            world_size=world_size,
            rank=rank
        )
    
    print(f'[Rank {rank}] 分布式环境初始化完成，使用设备: {device}')
    
    return rank, world_size, local_rank, device


def cleanup_distributed():
    """清理分布式环境"""
    if dist.is_initialized():
        dist.destroy_process_group()


def train_epoch(model, train_loader, criterion, optimizer, device, epoch, rank):
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
    
    # 聚合所有进程的损失和准确率
    total_loss_tensor = torch.tensor([total_loss]).to(device)
    correct_tensor = torch.tensor([correct]).to(device)
    total_tensor = torch.tensor([total]).to(device)
    
    dist.all_reduce(total_loss_tensor, op=dist.ReduceOp.SUM)
    dist.all_reduce(correct_tensor, op=dist.ReduceOp.SUM)
    dist.all_reduce(total_tensor, op=dist.ReduceOp.SUM)
    
    avg_loss = total_loss_tensor.item() / (len(train_loader) * dist.get_world_size())
    accuracy = 100. * correct_tensor.item() / total_tensor.item()
    
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
    total_loss_tensor = torch.tensor([total_loss]).to(device)
    correct_tensor = torch.tensor([correct]).to(device)
    total_tensor = torch.tensor([total]).to(device)
    
    dist.all_reduce(total_loss_tensor, op=dist.ReduceOp.SUM)
    dist.all_reduce(correct_tensor, op=dist.ReduceOp.SUM)
    dist.all_reduce(total_tensor, op=dist.ReduceOp.SUM)
    
    avg_loss = total_loss_tensor.item() / (len(val_loader) * dist.get_world_size())
    accuracy = 100. * correct_tensor.item() / total_tensor.item()
    
    return avg_loss, accuracy


def main():
    parser = argparse.ArgumentParser(description='PyTorch DDP 分布式训练示例')
    parser.add_argument('--epochs', type=int, default=10, help='训练轮数')
    parser.add_argument('--batch-size', type=int, default=64, help='每个进程的批次大小')
    parser.add_argument('--lr', type=float, default=0.001, help='学习率')
    parser.add_argument('--num-samples', type=int, default=10000, help='样本数量')
    parser.add_argument('--save-dir', type=str, default='/workspace/output', help='模型保存目录')
    parser.add_argument('--log-dir', type=str, default='/workspace/logs', help='TensorBoard日志目录')
    args = parser.parse_args()
    
    # 初始化分布式环境
    rank, world_size, local_rank, device = setup_distributed()
    
    try:
        # 只在 Rank 0 打印配置信息
        if rank == 0:
            print(f'\n[CONFIG] 训练配置:')
            print(f'  - Epochs: {args.epochs}')
            print(f'  - Batch Size (per GPU): {args.batch_size}')
            print(f'  - Global Batch Size: {args.batch_size * world_size}')
            print(f'  - Learning Rate: {args.lr}')
            print(f'  - World Size: {world_size}')
        
        # 创建数据集
        if rank == 0:
            print(f'\n[INFO] 创建模拟数据集，样本数: {args.num_samples}')
        
        dataset = create_dummy_data(num_samples=args.num_samples)
        train_size = int(0.8 * len(dataset))
        val_size = len(dataset) - train_size
        train_dataset, val_dataset = torch.utils.data.random_split(dataset, [train_size, val_size])
        
        # 使用 DistributedSampler 确保每个进程处理不同的数据
        train_sampler = DistributedSampler(train_dataset, num_replicas=world_size, rank=rank, shuffle=True)
        val_sampler = DistributedSampler(val_dataset, num_replicas=world_size, rank=rank, shuffle=False)
        
        train_loader = DataLoader(train_dataset, batch_size=args.batch_size, sampler=train_sampler, num_workers=2)
        val_loader = DataLoader(val_dataset, batch_size=args.batch_size, sampler=val_sampler, num_workers=2)
        
        # 创建模型
        model = SimpleNet().to(device)
        
        # 包装为 DDP 模型
        model = DDP(model, device_ids=[local_rank] if torch.cuda.is_available() else None)
        
        if rank == 0:
            print(f'[INFO] 模型已包装为 DistributedDataParallel')
        
        # 损失函数和优化器
        criterion = nn.CrossEntropyLoss()
        optimizer = optim.Adam(model.parameters(), lr=args.lr)
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
        
        # 训练循环
        if rank == 0:
            print(f'\n[INFO] 开始训练，共 {args.epochs} 轮')
            print('=' * 60)
        
        best_acc = 0.0
        start_time = time.time()
        
        for epoch in range(1, args.epochs + 1):
            epoch_start = time.time()
            
            # 设置 epoch，确保每个 epoch 的数据打乱不同
            train_sampler.set_epoch(epoch)
            
            # 训练
            train_loss, train_acc = train_epoch(model, train_loader, criterion, optimizer, device, epoch, rank)
            
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
                        'model_state_dict': model.module.state_dict(),
                        'optimizer_state_dict': optimizer.state_dict(),
                        'best_acc': best_acc,
                    }, model_path)
                    print(f'  [SAVE] 最佳模型已保存: {model_path}')
            
            # 同步所有进程
            dist.barrier()
        
        # 训练完成
        if rank == 0:
            total_time = time.time() - start_time
            print('=' * 60)
            print(f'\n[INFO] 训练完成!')
            print(f'  总用时: {total_time:.1f}s')
            print(f'  最佳验证准确率: {best_acc:.2f}%')
            
            if writer:
                writer.close()
    
    finally:
        cleanup_distributed()


if __name__ == '__main__':
    main()

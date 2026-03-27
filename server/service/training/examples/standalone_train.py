#!/usr/bin/env python3
"""
单机训练脚本示例
适用于 Neptune 平台的 STANDALONE 模式

使用方法:
1. 在创建训练任务时选择 STANDALONE 框架
2. 启动命令填写: python /path/to/standalone_train.py --epochs 10

特点:
- 单卡或单机多卡训练
- 使用 DataParallel 实现单机多卡
- 支持 GPU 自动检测
"""

import argparse
import os
import time
import torch
import torch.nn as nn
import torch.optim as optim
from torch.utils.data import DataLoader, TensorDataset


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


def train_epoch(model, train_loader, criterion, optimizer, device, epoch):
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
        
        if batch_idx % 100 == 0:
            print(f'Epoch: {epoch} | Batch: {batch_idx}/{len(train_loader)} | '
                  f'Loss: {loss.item():.4f} | Acc: {100.*correct/total:.2f}%')
    
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
    
    avg_loss = total_loss / len(val_loader)
    accuracy = 100. * correct / total
    return avg_loss, accuracy


def main():
    parser = argparse.ArgumentParser(description='单机训练示例')
    parser.add_argument('--epochs', type=int, default=10, help='训练轮数')
    parser.add_argument('--batch-size', type=int, default=64, help='批次大小')
    parser.add_argument('--lr', type=float, default=0.001, help='学习率')
    parser.add_argument('--num-samples', type=int, default=10000, help='样本数量')
    parser.add_argument('--save-dir', type=str, default='/workspace/output', help='模型保存目录')
    parser.add_argument('--log-dir', type=str, default='/workspace/logs', help='TensorBoard日志目录')
    args = parser.parse_args()
    
    # 设备检测
    if torch.cuda.is_available():
        device = torch.device('cuda')
        num_gpus = torch.cuda.device_count()
        print(f'[INFO] 使用 GPU 训练，检测到 {num_gpus} 个 GPU')
        for i in range(num_gpus):
            print(f'  - GPU {i}: {torch.cuda.get_device_name(i)}')
    else:
        device = torch.device('cpu')
        print('[INFO] 使用 CPU 训练')
    
    # 创建数据集
    print(f'[INFO] 创建模拟数据集，样本数: {args.num_samples}')
    dataset = create_dummy_data(num_samples=args.num_samples)
    train_size = int(0.8 * len(dataset))
    val_size = len(dataset) - train_size
    train_dataset, val_dataset = torch.utils.data.random_split(dataset, [train_size, val_size])
    
    train_loader = DataLoader(train_dataset, batch_size=args.batch_size, shuffle=True, num_workers=4)
    val_loader = DataLoader(val_dataset, batch_size=args.batch_size, shuffle=False, num_workers=4)
    
    # 创建模型
    model = SimpleNet()
    
    # 单机多卡: 使用 DataParallel
    if torch.cuda.device_count() > 1:
        print(f'[INFO] 使用 DataParallel，GPU 数量: {torch.cuda.device_count()}')
        model = nn.DataParallel(model)
    
    model = model.to(device)
    
    # 损失函数和优化器
    criterion = nn.CrossEntropyLoss()
    optimizer = optim.Adam(model.parameters(), lr=args.lr)
    scheduler = optim.lr_scheduler.StepLR(optimizer, step_size=5, gamma=0.5)
    
    # TensorBoard (可选)
    try:
        from torch.utils.tensorboard import SummaryWriter
        os.makedirs(args.log_dir, exist_ok=True)
        writer = SummaryWriter(args.log_dir)
        print(f'[INFO] TensorBoard 日志目录: {args.log_dir}')
    except ImportError:
        writer = None
        print('[WARN] 未安装 tensorboard，跳过日志记录')
    
    # 训练循环
    print(f'\n[INFO] 开始训练，共 {args.epochs} 轮')
    print('=' * 60)
    
    best_acc = 0.0
    start_time = time.time()
    
    for epoch in range(1, args.epochs + 1):
        epoch_start = time.time()
        
        # 训练
        train_loss, train_acc = train_epoch(model, train_loader, criterion, optimizer, device, epoch)
        
        # 验证
        val_loss, val_acc = validate(model, val_loader, criterion, device)
        
        # 更新学习率
        scheduler.step()
        
        epoch_time = time.time() - epoch_start
        
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
            # 如果是 DataParallel，保存原始模型
            state_dict = model.module.state_dict() if hasattr(model, 'module') else model.state_dict()
            torch.save({
                'epoch': epoch,
                'model_state_dict': state_dict,
                'optimizer_state_dict': optimizer.state_dict(),
                'best_acc': best_acc,
            }, model_path)
            print(f'  [SAVE] 最佳模型已保存: {model_path}')
    
    # 训练完成
    total_time = time.time() - start_time
    print('=' * 60)
    print(f'\n[INFO] 训练完成!')
    print(f'  总用时: {total_time:.1f}s')
    print(f'  最佳验证准确率: {best_acc:.2f}%')
    
    if writer:
        writer.close()


if __name__ == '__main__':
    main()

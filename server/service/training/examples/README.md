# 训练脚本示例

本目录包含适用于 Neptune AI 平台的训练脚本示例。

## 脚本列表

| 脚本 | 框架类型 | 说明 |
|------|---------|------|
| `standalone_train.py` | STANDALONE | 单机训练，支持单卡和单机多卡（DataParallel） |
| `ddp_train.py` | PYTORCH_DDP | PyTorch DistributedDataParallel 分布式训练 |
| `horovod_train.py` | MPI | Horovod MPI 分布式训练，支持梯度压缩 |

---

## 1. 单机训练 (STANDALONE)

### 使用方法

在创建训练任务时：
- **框架类型**: `STANDALONE`
- **启动命令**:
```bash
python /workspace/code/standalone_train.py --epochs 10 --batch-size 64
```

### 参数说明

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--epochs` | 10 | 训练轮数 |
| `--batch-size` | 64 | 批次大小 |
| `--lr` | 0.001 | 学习率 |
| `--num-samples` | 10000 | 模拟数据样本数 |
| `--save-dir` | /workspace/output | 模型保存目录 |
| `--log-dir` | /workspace/logs | TensorBoard 日志目录 |

### 特点
- 自动检测 GPU 并使用
- 单机多卡时自动使用 `DataParallel`
- 支持 TensorBoard 日志记录
- 自动保存最佳模型

---

## 2. PyTorch DDP 分布式训练 (PYTORCH_DDP)

### 使用方法

在创建训练任务时：
- **框架类型**: `PYTORCH_DDP`
- **Worker 数量**: 4（根据需要调整）
- **启动命令**:
```bash
# 方式1: 使用 torchrun (推荐)
torchrun --nnodes=$WORLD_SIZE --nproc_per_node=1 \
         --rdzv_id=$JOB_ID --rdzv_backend=c10d \
         --rdzv_endpoint=$MASTER_ADDR:$MASTER_PORT \
         /workspace/code/ddp_train.py --epochs 10

# 方式2: 直接运行 (Volcano pytorch 插件会自动设置环境变量)
python /workspace/code/ddp_train.py --epochs 10
```

### 环境变量

Volcano `pytorch` 插件会自动注入：

| 环境变量 | 说明 |
|---------|------|
| `MASTER_ADDR` | Master 节点地址 |
| `MASTER_PORT` | Master 端口（默认 29500） |
| `WORLD_SIZE` | 总进程数 |
| `RANK` | 当前进程全局排名 |
| `LOCAL_RANK` | 当前进程在本节点的排名 |

### 参数说明

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--epochs` | 10 | 训练轮数 |
| `--batch-size` | 64 | 每个进程的批次大小 |
| `--lr` | 0.001 | 学习率 |

### 特点
- 多机多卡分布式数据并行
- 使用 `DistributedSampler` 自动分割数据
- 使用 `DistributedDataParallel` 包装模型
- 梯度自动同步
- 仅 Rank 0 保存模型和写入日志

---

## 3. Horovod MPI 分布式训练 (MPI)

### 使用方法

在创建训练任务时：
- **框架类型**: `MPI`
- **Worker 数量**: 4（根据需要调整）
- **启动命令**:
```bash
# 平台会自动注入 MPI_HOSTS 环境变量
mpirun --allow-run-as-root -np 4 \
       -H $MPI_HOSTS \
       -bind-to none -map-by slot \
       -x NCCL_DEBUG=INFO \
       python /workspace/code/horovod_train.py --epochs 10

# 简化版（平台使用）
mpirun -np 4 -H worker-0:2,worker-1:2 python /workspace/code/horovod_train.py
```

### 参数说明

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--epochs` | 10 | 训练轮数 |
| `--batch-size` | 64 | 每个进程的批次大小 |
| `--lr` | 0.001 | 基础学习率（会自动缩放） |
| `--fp16` | False | 启用混合精度训练 |
| `--compression` | none | 梯度压缩方式（none/fp16） |

### 特点
- 使用 Horovod 实现 MPI 分布式训练
- 学习率自动按 world_size 缩放
- 支持 FP16 梯度压缩
- 支持混合精度训练
- 兼容 NCCL/Gloo/MPI 后端

### 安装 Horovod

```bash
# CUDA 环境
HOROVOD_GPU_OPERATIONS=NCCL pip install horovod

# CPU 环境
pip install horovod
```

---

## 4. 目录结构建议

在训练任务中，建议使用以下目录结构：

```
/workspace/
├── code/           # 代码目录（挂载代码 PVC）
│   ├── train.py
│   └── model.py
├── data/           # 数据目录（挂载数据集 PVC）
│   ├── train/
│   └── val/
├── output/         # 输出目录（挂载输出 PVC）
│   ├── checkpoints/
│   └── best_model.pth
└── logs/           # 日志目录（TensorBoard）
    └── events.out.tfevents.*
```

---

## 5. TensorBoard 集成

所有脚本都支持 TensorBoard 日志记录：

1. 创建训练任务时勾选 **启用 TensorBoard**
2. 填写 **日志路径**: `/workspace/logs`（与 `--log-dir` 参数一致）
3. 训练开始后，可通过平台访问 TensorBoard

### 记录的指标

| 指标 | 说明 |
|------|------|
| `Loss/train` | 训练损失 |
| `Loss/val` | 验证损失 |
| `Accuracy/train` | 训练准确率 |
| `Accuracy/val` | 验证准确率 |
| `LearningRate` | 当前学习率 |

---

## 6. 常见问题

### Q: DDP 训练时某个节点卡住
**A**: 检查 NCCL 网络配置，确保所有节点间网络互通。可以设置 `NCCL_DEBUG=INFO` 查看详细日志。

### Q: MPI 训练时 SSH 连接失败
**A**: 平台使用 Volcano ssh 插件自动配置免密 SSH。如果仍然失败，检查 Worker Pod 的 sshd 服务是否正常启动。

### Q: 共享内存不足导致 DataLoader 失败
**A**: 确保训练任务启用了共享内存（UseSHM=true），平台会自动配置 `/dev/shm`。

### Q: 如何使用真实数据集？
**A**: 
1. 将数据上传到 PVC
2. 创建训练任务时添加挂载配置
3. 修改脚本中的数据加载路径

```python
# 替换模拟数据
from torchvision import datasets, transforms

transform = transforms.Compose([
    transforms.ToTensor(),
    transforms.Normalize((0.5,), (0.5,))
])

train_dataset = datasets.MNIST('/workspace/data', train=True, download=False, transform=transform)
```

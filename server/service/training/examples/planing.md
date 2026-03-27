# 开发规范
## 项目结构
```bash
项目结构分为api router model service，所有的router url以api/v1为起始，然后加业务名称，除了必须用get外，其余全部用post。
所有的service 入参的req和response结构体都要在model层的requset和response创建，不要直接返回db的model字段。
db model和其他model用id关联，model创建完后需要在/server/initialize/gorm.go中通过gorm的迁移方法自动迁移到数据库中
工具放到对应的utils包中
如无特殊要求，不要总结并创建文档
所有的列表都需要有刷新按钮，并且5s/10s刷新一次请求。这个按钮全平台列表统一。
```


# 菜单视图
```bash
平台是按照用户角色来展示菜单的，创建完菜单后需要在平台上的超级管理员菜单下的菜单管理模块添加菜单，然后在角色管理绑定对应的菜单权限。
服务启动时，会给超级管理员添加一些默认的菜单和api权限，代码在neptune/server/source/system/menu.go，neptune/server/source/system/api.go。
```

# 开发功能
```bash
机器学习平台 - 训练任务模块概要设计 
1. 设计原则 (Design Principles)
在进入具体设计前，确立四条不可撼动的核心原则：
 - 强类型约束 (Strict Typing): 
 - 严禁使用字符串拼接 (fmt.Sprintf) 或简单的 yaml.Marshal 生成 K8s 资源。
 - 必须使用 Go Structs (volcano.sh/apis + k8s.io/api) 配合 Typed Client。
 - 数据库禁用json字段类型存储。
设计方案可以根据开发实际过程调整，但是四条核心原则不动！
收益: 编译期查错、防止注入攻击、完美的 IDE 补全、平滑的版本升级。
K8s 为事实来源 (Single Source of Truth): 数据库仅作为用户意图的持久化存储和历史归档。任务的实时状态（Running/Failed/Completed）以 K8s 集群内的实际状态为准，后端充当 Controller 角色，而非简单的 CRUD。
事件驱动架构 (Event-Driven): 状态同步不依赖定时轮询（Polling），而是依赖 K8s 的 List-Watch (Informer) 机制。

2. 核心业务架构图Code snippetsequenceDiagram
    participant User as 用户 (Frontend)
    participant API as 后端 API
    participant DB as 数据库 (MySQL)
    participant Builder as Job Builder (Go)
    participant K8s as K8s API (Volcano)
    participant Informer as K8s Watch/Informer

    %% 提交阶段
    User->>API: 1. 提交训练任务 (DDP, 4 GPU)
    API->>DB: 2. 插入元数据 (status=SUBMITTED)
    API->>Builder: 3. 构建 Volcano Job Struct (内存对象)
    Builder-->>API: 返回 *vcbatch.Job 对象
    API->>K8s: 4. client.BatchV1alpha1().Jobs(ns).Create(ctx, job)
    K8s-->>API: 5. 返回 Created Job (UID)
    API->>DB: 6. 更新 k8s_job_id, status=CREATING

    %% 状态同步阶段 (异步)
    loop Watch Loop
        K8s->>Informer: 7. Job Update Event (Pending -> Running)
        Informer->>API: 8. 回调 UpdateFunc
        API->>DB: 9. 更新 DB status=RUNNING, updated_at=Now
    end

    %% 日志查看
    User->>API: 10. 请求日志
    API->>K8s: 11. client.CoreV1().Pods().GetLogs()
    K8s-->>User: 返回 Log Stream

3. 数据库设计 (Data Model, 以下数据库均设计为暂定的数据库DDL，可以根据实现过程中发现的问题进行调整，但是字段不要存 json！！！！！！！)
3.1 主表：训练任务 (training_jobs)
SQLCREATE TABLE training_jobs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    -- 业务ID与归属
    job_name VARCHAR(63) NOT NULL, -- 需符合DNS-1123规范 (小写字母数字短横线)
    user_id BIGINT NOT NULL,
    namespace VARCHAR(63) NOT NULL DEFAULT 'default', -- 租户隔离
    
    -- 核心配置
    framework_type VARCHAR(50) NOT NULL, -- 'PYTORCH_DDP', 'MPI', 'STANDALONE'
    image_id int NOT NULL,
    startup_command TEXT, 
    
    -- 资源配置 (展平关键字段以便统计，详情存JSON)
    total_gpu_count INT DEFAULT 0,
    gpu_type VARCHAR(50) DEFAULT 'nvidia.com/gpu', -- 兼容多厂商
    resource_id int,
    resource_type VARCHAR(50) DEFAULT 'gpu',    
    -- 高级网络与环境
    envs JSON, 
    use_shm BOOLEAN DEFAULT TRUE,
    
    -- 状态管理
    k8s_job_id VARCHAR(128),      -- K8s UID，防止同名任务混淆
    status VARCHAR(50) NOT NULL,  -- SUBMITTED, PENDING, RUNNING, SUCCEEDED, FAILED, KILLING, KILLED, UNKNOWN
    error_msg TEXT,               -- 记录 ImagePullBackOff 等具体错误
    
    -- 时间审计
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    started_at TIMESTAMP NULL,    -- 真正变为 Running 的时间
    finished_at TIMESTAMP NULL,    
    UNIQUE KEY uk_ns_name (namespace, job_name) -- 命名空间下唯一
);

3.2 关联表：挂载配置 (training_job_mounts)
将“数据集”抽象化，而不是把路径写死在 JSON 里。
SQLCREATE TABLE training_job_mounts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    job_id BIGINT NOT NULL,
    
    -- 挂载类型
    mount_type VARCHAR(20),  -- 'DATASET', 'MODEL', 'CODE', 'OUTPUT'
    
    -- 资源引用 (指向平台管理的 Dataset/Model 表)
    source_id VARCHAR(128),  
    
    -- 底层存储映射 (由后端根据 source_id 解析)
    pvc_id VARCHAR(255),   
    sub_path VARCHAR(255) DEFAULT '', -- PVC内的子路径
    
    -- 容器内目标
    mount_path VARCHAR(512) NOT NULL,
    read_only BOOLEAN DEFAULT TRUE,
);

4. 后端核心开发 (Backend， serivce层统一采用interface方式实现，可以参考service/notebook.go的实现)
采用 Builder Pattern 和 Typed Client。
4.1 核心接口定义 (JobBuilder)
后端不直接写逻辑，而是定义接口，针对不同框架实现不同的 Builder。Gopackage builder
import (
	"context"
	vcbatch "volcano.sh/apis/pkg/apis/batch/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)
// TrainingJobSpec 是业务层的 DTO
type TrainingJobSpec struct {
	Name        string
	Namespace   string
	Framework   string
	Image       string
	Command     []string
	WorkerCount int32
	Volumes     []corev1.Volume // 预处理后的 Volumes
    // ...
}

// 核心构建器接口
type JobBuilder interface {
	Build(spec TrainingJobSpec) (*vcbatch.Job, error)
}

// 框架抽象接口 (Strategy Pattern)
type FrameworkStrategy interface {
	BuildTasks(spec TrainingJobSpec) ([]vcbatch.TaskSpec, error)
	GetPlugins() map[string][]string
	GetMinAvailable(spec TrainingJobSpec) int32
}
4.2 实现示例 (PyTorch DDP Builder)
核心逻辑：
Plugins: 注入 `pytorch`, `svc`, `env`。
Tasks: 拆分为 `master` (1 replica) 和 `worker` (N-1 replicas)。
Command: **后端不封装启动命令**。直接透传用户输入的命令。
环境变量: 自动注入 `MASTER_ADDR`, `MASTER_PORT`, `WORLD_SIZE`, `RANK`。

```go
func (s *PyTorchDDPStrategy) BuildTasks(spec *trainingReq.TrainingJobSpec) ([]vcbatch.TaskSpec, error) {
    // ... 构建 Master 和 Worker Task
    // 容器 Command 直接使用 spec.Command
    // 容器 Args 直接使用 spec.Args
}
```

GetPlugins 返回：
```go
func (s *PyTorchDDPStrategy) GetPlugins() map[string][]string {
    return map[string][]string{
        "pytorch": {"--master=master", "--worker=worker", "--port=29500"},
        "svc":     {}, // 启用 Service Discovery
        "env":     {}, // 自动注入环境变量
    }
}
```
4.3 状态同步 (Informer/Watch)
后端启动时初始化 SharedInformerFactory。Gofunc (c *JobController) StartInformer(stopCh <-chan struct{}) {
    // 监听带有特定 Label 的 Volcano Jobs
    c.vcInformer.Batch().V1alpha1().Jobs().Informer().AddEventHandler(
        cache.ResourceEventHandlerFuncs{
            UpdateFunc: func(oldObj, newObj interface{}) {
                newJob := newObj.(*vcbatch.Job)
                // 仅当状态发生实质变化时更新 DB
                c.syncJobStatusToDB(newJob)
            },
            DeleteFunc: func(obj interface{}) {
                // 处理外部删除的情况，更新 DB 为 KILLED 或 UNKNOWN
            },
        },
    )
}

func (c *JobController) syncJobStatusToDB(job *vcbatch.Job) {
    // 映射 K8s 状态到 DB 状态
    var dbStatus string
    switch job.Status.State.Phase {
    case vcbatch.Pending:
        dbStatus = "PENDING"
    case vcbatch.Running:
        dbStatus = "RUNNING"
    case vcbatch.Completed:
        dbStatus = "SUCCEEDED"
    case vcbatch.Failed, vcbatch.Terminated: // Terminated 通常视为 Failed
        dbStatus = "FAILED"
    default:
        dbStatus = "UNKNOWN"
    }
    
    // 记录 Event，更新 DB
    // ...
}

5. 存储与容错 (Storage & Resilience)
5.1 SHM 
优化为了避免 OOM，代码层构建 Volume 时必须添加 SizeLimit。Goif spec.UseSHM {
    volumes = append(volumes, corev1.Volume{
        Name: "shm",
        VolumeSource: corev1.VolumeSource{
            EmptyDir: &corev1.EmptyDirVolumeSource{
                Medium: corev1.StorageMediumMemory,
                SizeLimit: resource.NewQuantity(2*1024*1024*1024, resource.BinarySI), // 2Gi
            },
        },
    })
}
5.2 容错策略利用 
Volcano 的 Policies 处理部分故障。MPI/DDP: 任意 Pod 失败 -> RestartJob (因为需要全员参与)。
Standalone: Pod 失败 -> RestartPod (或根据 maxRetry 决定)。Gojob.Spec.Policies = []vcbatch.LifecyclePolicy{
    {Event: vcbatch.PodFailed, Action: vcbatch.RestartJob},
    {Event: vcbatch.PodEvicted, Action: vcbatch.RestartJob},
}
job.Spec.MaxRetry = 3


### 5.3 MPI 训练模式专项设计 (MPI Specialization)
5.3.1 功能定位与对比
MPI 模式专用于传统 HPC 程序、OpenMPI、以及基于 Horovod 的分布式训练。
维度               PyTorch DDP                                         MPI / Horovod
通信后端            NCCL / Gloo                                         MPI (OpenMPI, MPICH)
拓扑结构            Master + Worker (角色模糊，通常全员计算)                mpimaster (Launcher) + worker (Launcher 不计算)
启动方式            torch.distributed.launch                            mpirun (由 Launcher 触发 SSH 远程执行)
资源分配            每 Pod 申请 GPU                                      Launcher: 0 GPU; Worker: N GPU
Volcano插件         pytorch, svc                                        ssh (免密), svc (DNS)

5.3.2 用户侧交互设计 (Frontend)
当用户在“框架选择”中选中 [MPI] 时，表单发生联动。
增加提示：Master 节点将自动额外创建，不消耗 GPU。
**启动命令提示**：系统不自动封装 mpirun，用户需手动输入。提供示例按钮，利用环境变量 `${MPI_HOSTS}` 简化主机配置。

5.3.3 更新核心模型 (Backend Data Model)
后端在构建 MPI 任务时，会动态生成所有节点的 Host 列表，并以 `MPI_HOSTS` 环境变量注入到 Master 容器中。

5.3.4 MPI 构建策略 (MPIStrategy)
核心逻辑：
Plugins: 必须注入 `ssh` 和 `svc`。
Tasks: 拆分为 `mpimaster` (1 replica) 和 `worker` (N replicas)。
Command: **后端不封装 mpirun**。直接透传用户输入的命令。
环境变量: 自动注入 `MPI_HOSTS` (格式为 `host1:slots,host2:slots...`)。

```go
func (b *MPIStrategy) BuildTasks(spec *trainingReq.TrainingJobSpec) ([]vcbatch.TaskSpec, error) {
    // 1. 生成 Host 列表 (包含 Master 和所有 Workers)
    var hosts []string
    hosts = append(hosts, fmt.Sprintf("%s-mpimaster-0.%s", spec.Name, spec.Name))
    for i := 0; i < int(spec.WorkerCount); i++ {
        slots := 1
        if spec.WorkerGPU > 0 { slots = int(spec.WorkerGPU) }
        hosts = append(hosts, fmt.Sprintf("%s-worker-%d.%s:%d", spec.Name, i, spec.Name, slots))
    }
    hostStr := strings.Join(hosts, ",")

    // 2. 注入环境变量
    masterEnvs := append(spec.Envs, corev1.EnvVar{Name: "MPI_HOSTS", Value: hostStr})

    // 3. 构建 Master (Launcher)
    master := vcbatch.TaskSpec{
        Name: "mpimaster",
        Template: corev1.PodTemplateSpec{
            Spec: corev1.PodSpec{
                Containers: []corev1.Container{{
                    Command: spec.Command, // 用户输入的 mpirun -H ${MPI_HOSTS} ...
                    Env:     masterEnvs,
                    // ... 资源限制 (1C2G)
                }},
            },
        },
    }
    // 4. 构建 Worker (sshd 守护进程)
    // ...
}
```

5.3.5 日志与可观测性 (Observability)
MPI 日志策略：
Master (mpimaster): 聚合了所有 rank 的标准输出。用户查看日志时，后端默认拉取 mpimaster Pod 的日志。
Worker: 仅作为计算节点，通常只运行 sshd。

5.3.6 TensorBoard 集成
与 DDP 模式一致，通过共享 PVC 挂载日志目录。

5.3.7 调试 (Web Terminal)
用户连接到 `mpimaster` Pod 的 Shell 后，可以直接执行 `mpirun` 进行手动调试。
示例：`mpirun -np 2 -H ${MPI_HOSTS} hostname`

5.3.8 总结：开发 Check List
API: 确认 POST /training-jobs 接口接收 framework_type。
Builder: 实现 MPIBuilder，重点测试 MPI_HOSTS 注入和 ssh/svc 插件。
K8s: 确保底层集群安装了 Volcano 且启用了 ssh/svc plugin。
Image: 确保提供的 Base Image 已安装 openssh-server 和 openmpi/mpich。
Network: 确保 APISix 已部署，且后端有权限调用 APISix Admin API 创建路由。


6. TensorBoard 与监控 (Observability)
TensorBoard 
前端访问:TensorBoard 支持子路径访问：https://tb.ai-platform.com/sys/tensorboard/ (需配置 TB 的 path_prefix,参考serivec/notebook.go的实现)。
（TODO 监控指标:集成 Prometheus + Grafana。后端不直接存监控数据，而是代理查询 Prometheus (rate(container_gpu_utilization[5m])) 并在前端绘图。）

7. ApiSix || SSHPiper (功能)
Apisix 要求和notebook创建时一样，具有路由访问功能（包括tensorboard），jupyter应该也是需要加上的吧？如果训练完了怎么进入容器内查看训练的数据呢？ 
SSHPiper 和jupyter一样，从我个人角度来理解，应该也是需要加入的吧？要不然用户怎么连接容器呢？
 
8. 前端交互 (Frontend)
 资源预检: 选择 GPU 数量时，调用后端 API,如果 request > available，显示黄色警告“任务将排队”。
 动态表单:选择 PyTorch DDP -> 显示 Worker Count。选择 Standalone -> 隐藏 Worker Count。
 日志查看:支持 Auto Refresh (轮询 API) 和 Follow (WebSocket)。对于已结束的任务，显示“下载完整日志”按钮。

9. apisix和Jupyter以及ssh的考量
Volcano Job (训练任务)：瞬态 (Ephemeral) 的 Batch Job。用于全量数据高并发训练
Jupyter Notebook Pod
   ↕ 共享 PVC
Training Job (Volcano)
用户：在 Jupyter 里写代码，提交训练任务
训练：使用同一个代码 / 数据 PVC
连接桥梁：共享存储 (RWX PVC)

APISix
    TensorBoard 访问 (核心需求)
推荐 - Web Terminal，前端提供一个 Web Shell，前端提供一个 Web Shell（基于 kubectl exec + WebSocket，例如 xterm.js）。这是最云原生的方式，不需要开放 SSH 端口，权限完全由 K8s RBAC 控制
```





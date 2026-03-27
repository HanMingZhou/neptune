package consts

// Training
const (
	TrainingInstance = "training"

	// 框架类型常量
	FrameworkPyTorchDDP = "PYTORCH_DDP" // PyTorch 分布式数据并行
	FrameworkMPI        = "MPI"         // MPI 分布式
	FrameworkStandalone = "STANDALONE"  // 单机训练

	// 任务角色名称常量
	TaskSpecMaster    = "master"    // 主节点任务角色
	TaskSpecWorker    = "worker"    // 工作节点任务角色
	TaskSpecMPIMaster = "mpimaster" // MPI的主节点任务角色
	TaskSpecMPIWorker = "mpiworker" // MPI的工作节点任务角色

	// 挂载类型常量
	MountTypeDataset = "DATASET" // 数据集
	MountTypeModel   = "MODEL"   // 模型
	MountTypeCode    = "CODE"    // 代码
	MountTypeOutput  = "OUTPUT"  // 输出目录

	TrainingStatusSubmitted = "SUBMITTED"
	TrainingStatusPending   = "PENDING"
	TrainingStatusRunning   = "RUNNING"
	TrainingStatusSucceeded = "SUCCEEDED"
	TrainingStatusFailed    = "FAILED"
	TrainingStatusKilling   = "KILLING"
	TrainingStatusKilled    = "KILLED"
	TrainingStatusUnknown   = "UNKNOWN"
)

// Notebook
const (
	NotebookInstance = "notebook"

	NotebookStatusCreating     = "CREATING"
	NotebookStatusRunning      = "RUNNING"
	NotebookStatusPending      = "PENDING"
	NotebookStatusStopped      = "STOPPED"
	NotebookStatusFailed       = "FAILED"
	NotebookStatusSucceeded    = "SUCCEEDED"
	NotebookStatusDeleting     = "DELETING"
	NotebookStatusDeleteFailed = "DELETE_FAILED"
	NotebookStatusUpdating     = "UPDATING"
	NotebookStatusUpdateFailed = "UPDATE_FAILED"
	NotebookStatusUpdated      = "UPDATED"
)

// Inference
const (
	InferenceInstance = "inference"

	DefaultSHM = 2

	// 部署类型常量
	DeployTypeDistributed = "DISTRIBUTED" // 分布式推理 (VCJob)
	DeployTypeStandalone  = "STANDALONE"  // 单体推理 (Deployment)

	// 推理框架常量
	FrameworkSGLang = "SGLANG" // SGLang 框架
	FrameworkVLLM   = "VLLM"   // vLLM 框架

	InferenceStatusCreating   = "CREATING"
	InferenceStatusPending    = "PENDING"
	InferenceStatusRunning    = "RUNNING"
	InferenceStatusFailed     = "FAILED"
	InferenceStatusStopped    = "STOPPED"
	InferenceStatusDeleting   = "DELETING"
	InferenceStatusRestarting = "RESTARTING"
)

// Sheculer
const (
	VolcanoScheduler = "volcano"
)

// PodGroup
const (
	// PodGroup Phase 常量（Volcano 原生状态）
	PhasePending = "Pending"
	PhaseRunning = "Running"
	PhaseUnknown = "Unknown"
	PhaseInqueue = "Inqueue"

	PodGroupStatusCreating  = "CREATING"
	PodGroupStatusPending   = "PENDING"
	PodGroupStatusRunning   = "RUNNING"
	PodGroupStatusCompleted = "COMPLETED"
	PodGroupStatusFailed    = "FAILED"
	PodGroupStatusDeleting  = "DELETING"
	PodGroupStatusDeleted   = "DELETED"
)

// Volcano 状态常量
const (
	VolcanoPhasePending   = "Pending"
	VolcanoPhaseInqueue   = "Inqueue"
	VolcanoPhaseRunning   = "Running"
	VolcanoPhaseSucceeded = "Succeeded"
	VolcanoPhaseFailed    = "Failed"
	VolcanoPhaseUnknown   = "Unknown"

	VolcanoVGPUMemory = "volcano.sh/vgpu-memory"
	VolcanoVGPUCores  = "volcano.sh/vgpu-cores"
	VolcanoVGPUNumber = "volcano.sh/vgpu-number"
)

// Queue
const (
	DefalutQueue = "default"
)

// PVC
const (
	// 模型来源常量 (当前仅支持 PVC)
	ModelSourcePVC = "PVC" // PVC 挂载

	DefaultStorageClass = "standard"
	PVCStatusReady      = "READY"
	PVCStatusBound      = "BOUND"
	PVCStatusError      = "ERROR"
)

// Order Order Types
const (
	OrderTypeNotebook  = 1
	OrderTypeTraining  = 2
	OrderTypeInference = 3
	OrderTypeStorage   = 4
)

// Invoice Status
const (
	InvoiceStatusProcessing = "PROCESSING"
	InvoiceStatusSent       = "SENT"
)

// Invoice Type
const (
	InvoiceTypePersonal   = "PERSONAL"
	InvoiceTypeEnterprise = "ENTERPRISE"
)

// Pay Method
const (
	PayMethodWechat  = "WECHAT"
	PayMethodAlipay  = "ALIPAY"
	PayMethodBalance = "BALANCE"
	PayMethodSystem  = "SYSTEM"
)

var PayMethodToInt64 = map[string]int64{
	PayMethodWechat:  1,
	PayMethodAlipay:  2,
	PayMethodBalance: 3,
	PayMethodSystem:  4,
}

// Default Namespaces
const (
	DefaultNamespace = "default"
)

// Nvidia GPU Type
const (
	LabelCPUModel   = "beta.kubernetes.io/arch"
	NvidiaGPUType   = "nvidia.com/gpu"
	AmdGPUType      = "amd.com/gpu"
	VGPUType        = "volcano.sh/vgpu"
	LabelGPUModel   = "gpu-model"
	LabelGPUMemAlt  = "gpu-memory"
	LabelGPUMem     = "nvidia.com/gpu.memory"
	LabelGpuProduct = "nvidia.com/gpu.product"
)

// TensorBoard
const (
	TensorBoardInstance        = "tensorboard"
	StatusCreating             = "CREATING"
	StatusRunning              = "RUNNING"
	StatusFailed               = "FAILED"
	DefaultTensorBoardLogsPath = "logs"
)

// 标签
const (
	LabelCreatedBy = "created-by"
	LabelInstance  = "instance"
	LabelType      = "type"

	LabelApp             = "app"
	LabelJobID           = "job-id"
	LabelInstanceType    = "instance-type"
	LabelOwnerType       = "owner-type"
	LabelFramework       = "framework"
	LabelVolcanoJob      = "volcano.sh/job-type"
	LabelVolcanoJobName  = "volcano.sh/job-name"
	LabelVolcanoTaskSpec = "volcano.sh/task-spec"

	// 标签值
	LabelValuePlatform = "kubeflow-platform"
	LabelValueSSH      = "ssh"

	// 节点标签
	LabelNodeInstanceType = "node.kubernetes.io/instance-type"
	LabelAccelerator      = "accelerator"

	// 节点角色标签
	LabelNodeRoleMaster       = "node-role.kubernetes.io/master"
	LabelNodeRoleControlPlane = "node-role.kubernetes.io/control-plane"
)

// SSHPiper
const (
	SSHPiperNamespace = "kubeflow"
)

package inference

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model/consts"
	"time"
)

// 认证类型常量
const (
	AuthTypeJWT    = 1 // Authorization (JWT Token)
	AuthTypeApiKey = 2 // API-Key 认证
)

// VolcanoJobStatusMap Volcano Job 状态映射表
var VolcanoJobStatusMap = map[string]string{
	"Pending":    consts.InferenceStatusPending,
	"Running":    consts.InferenceStatusRunning,
	"Succeeded":  consts.InferenceStatusStopped,
	"Failed":     consts.InferenceStatusFailed,
	"Aborted":    consts.InferenceStatusFailed,
	"Terminated": consts.InferenceStatusStopped,
}

// DeploymentStatusMap Deployment 状态映射表
var DeploymentStatusMap = map[string]string{
	"Available":   consts.InferenceStatusRunning,
	"Progressing": consts.InferenceStatusPending,
	"Failed":      consts.InferenceStatusFailed,
}

// Inference 推理服务数据库模型
type Inference struct {
	global.GVA_MODEL
	DisplayName  string `json:"displayName" gorm:"column:display_name;type:varchar(63);not null;comment:展示名称"`
	InstanceName string `json:"instanceName" gorm:"column:instance_name;type:varchar(100);not null;uniqueIndex;comment:实例名称"`
	UserId       uint   `json:"userId" gorm:"column:user_id;not null;index;comment:用户ID"`
	Namespace    string `json:"namespace" gorm:"column:namespace;type:varchar(63);not null;default:default;comment:命名空间"`
	ClusterID    uint   `json:"clusterId" gorm:"column:cluster_id;not null;index;comment:集群ID"`

	// 部署配置
	DeployType string `json:"deployType" gorm:"column:deploy_type;type:varchar(20);not null;comment:部署类型(DISTRIBUTED/STANDALONE)"`
	Framework  string `json:"framework" gorm:"column:framework;type:varchar(20);not null;comment:推理框架(SGLANG/VLLM)"`

	// 模型存储卷配置
	ModelMountPath string `json:"modelMountPath" gorm:"column:model_mount_path;type:varchar(256);default:/model;comment:PVC挂载路径"`
	ModelPvcId     uint   `json:"modelPvcId" gorm:"column:model_pvc_id;not null;comment:模型PVC ID"`

	// 镜像配置
	ImageId uint `json:"imageId" gorm:"column:image_id;not null;comment:镜像ID"`

	// 并行配置 (分布式模式)
	TensorParallel   int `json:"tensorParallel" gorm:"column:tensor_parallel;default:1;comment:张量并行度"`
	PipelineParallel int `json:"pipelineParallel" gorm:"column:pipeline_parallel;default:1;comment:流水线并行度"`
	WorkerCount      int `json:"workerCount" gorm:"column:worker_count;default:0;comment:Worker数量"`

	// 资源配置
	ProductId  uint   `json:"productId" gorm:"column:product_id;comment:产品ID"`
	CPU        int64  `json:"cpu" gorm:"column:cpu;comment:CPU配置"`
	Memory     int64  `json:"memory" gorm:"column:memory;comment:内存配置(GiB)"`
	GPU        int64  `json:"gpu" gorm:"column:gpu;comment:单节点GPU数量"`
	GPUModel   string `json:"gpuModel" gorm:"column:gpu_model;comment:GPU型号"`
	VGPUNumber int64  `json:"vGpuNumber" gorm:"column:v_gpu_number;comment:vGPU数量;default:0"`
	VGPUMemory int64  `json:"vGpuMemory" gorm:"column:v_gpu_memory;comment:vGPU显存 (GB);default:0"`
	VGPUCores  int64  `json:"vGpuCores" gorm:"column:v_gpu_cores;comment:vGPU核心数;default:0"`

	// 服务配置
	ServicePort int    `json:"servicePort" gorm:"column:service_port;default:30000;comment:服务端口"`
	Command     string `json:"command" gorm:"column:command;type:text;comment:启动命令(原始字符串)"`
	Args        string `json:"args" gorm:"column:args;type:text;comment:启动参数(JSON格式)"`

	// 自愈配置 (分布式模式)
	AutoRestart  bool `json:"autoRestart" gorm:"column:auto_restart;default:false;comment:是否允许自动重启"`
	RestartCount int  `json:"restartCount" gorm:"column:restart_count;default:0;comment:已重启次数"`
	MaxRestarts  int  `json:"maxRestarts" gorm:"column:max_restarts;default:3;comment:最大重启次数"`

	// 认证配置
	AuthType int `json:"authType" gorm:"column:auth_type;default:1;comment:认证类型 1-Authorization 2-API-Key"`

	// 状态信息
	Status         string `json:"status" gorm:"column:status;type:varchar(50);not null;default:CREATING;index;comment:状态"`
	ErrorMsg       string `json:"errorMsg" gorm:"column:error_msg;type:text;comment:错误信息"`
	K8sResourceUid string `json:"k8sResourceUid" gorm:"column:k8s_resource_uid;type:varchar(128);index;comment:K8s资源UID"`

	// 时间信息
	StartedAt *time.Time `json:"startedAt" gorm:"column:started_at;comment:启动时间"`

	// 计费相关
	PayType int64   `json:"payType" gorm:"column:pay_type;default:1;comment:付费类型(1-按量 2-包日 3-包周 4-包月)"`
	Price   float64 `json:"price" gorm:"column:price;type:decimal(10,4);comment:单价"`
	OrderId uint    `json:"orderId" gorm:"column:order_id;comment:关联订单ID"`
}

func (Inference) TableName() string {
	return "inferences"
}

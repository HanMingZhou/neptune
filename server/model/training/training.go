package training

import (
	"gin-vue-admin/global"
	"time"
)

// TrainingJob 训练任务数据库模型
// 资源规格（CPU/Memory/GPU等）不冗余存储，通过 ProductId 查询产品表获取
type TrainingJob struct {
	global.GVA_MODEL
	DisplayName      string     `json:"displayName" gorm:"column:display_name;type:varchar(63);not null;comment:展示名称"`
	UserId           uint       `json:"userId" gorm:"column:user_id;not null;index;comment:用户ID"`
	Namespace        string     `json:"namespace" gorm:"column:namespace;type:varchar(63);not null;default:default;comment:命名空间"`
	ClusterID        uint       `json:"clusterId" gorm:"column:cluster_id;not null;index;comment:集群ID"`
	FrameworkType    string     `json:"frameworkType" gorm:"column:framework_type;type:varchar(50);not null;comment:框架类型(PYTORCH_DDP/MPI/STANDALONE)"`
	ImageId          uint       `json:"imageId" gorm:"column:image_id;not null;comment:镜像ID"`
	StartupCommand   string     `json:"startupCommand" gorm:"column:startup_command;type:text;comment:启动命令"`
	TotalGPUCount    int64      `json:"totalGpuCount" gorm:"column:total_gpu_count;default:0;comment:GPU总数(快照,用于统计)"`
	ProductId        uint       `json:"productId" gorm:"column:product_id;comment:产品ID;index"`
	InstanceCount    int64      `json:"instanceCount" gorm:"column:instance_count;default:1;comment:实例数量"`
	ScheduleStrategy string     `json:"scheduleStrategy" gorm:"column:schedule_strategy;type:varchar(20);default:BALANCED;comment:调度策略"`
	WorkerCount      int64      `json:"workerCount" gorm:"column:worker_count;comment:Worker数量"`
	K8sJobUid        string     `json:"k8sJobUid" gorm:"column:k8s_job_uid;type:varchar(128);index;comment:K8s Job UID"`
	JobName          string     `json:"jobName" gorm:"column:job_name;type:varchar(128);comment:K8s Job名称"`
	Status           string     `json:"status" gorm:"column:status;type:varchar(50);not null;default:SUBMITTED;index;comment:状态"`
	ErrorMsg         string     `json:"errorMsg" gorm:"column:error_msg;type:text;comment:错误信息"`
	StartedAt        *time.Time `json:"startedAt" gorm:"column:started_at;comment:开始时间"`
	FinishedAt       *time.Time `json:"finishedAt" gorm:"column:finished_at;comment:结束时间"`
	// TensorBoard 相关
	EnableTensorboard  bool   `json:"enableTensorboard" gorm:"column:enable_tensorboard;default:false;comment:是否启用TensorBoard"`
	TensorboardLogPath string `json:"tensorboardLogPath" gorm:"column:tensorboard_log_path;type:varchar(512);comment:TensorBoard日志路径"`
	TensorboardId      uint   `json:"tensorboardId" gorm:"column:tensorboard_id;comment:关联的TensorBoard ID"`
	// 计费相关
	PayType int64   `json:"payType" gorm:"column:pay_type;default:1;comment:付费类型(1-按量 2-包日 3-包周 4-包月)"`
	Price   float64 `json:"price" gorm:"column:price;type:decimal(10,4);comment:单价"`
	OrderId uint    `json:"orderId" gorm:"column:order_id;comment:关联订单ID"`
}

func (TrainingJob) TableName() string {
	return "training_jobs"
}

// TrainingJobMount 训练任务挂载配置
type TrainingJobMount struct {
	global.GVA_MODEL
	JobId     uint   `json:"jobId" gorm:"column:job_id;not null;index;comment:关联的训练任务ID"`
	MountType string `json:"mountType" gorm:"column:mount_type;type:varchar(20);comment:挂载类型(DATASET/MODEL/CODE/OUTPUT)"`
	SourceId  uint   `json:"sourceId" gorm:"column:source_id;comment:资源引用ID"`
	PvcId     uint   `json:"pvcId" gorm:"column:pvc_id;comment:PVC ID"`
	PvcName   string `json:"pvcName" gorm:"column:pvc_name;type:varchar(255);comment:PVC名称"`
	SubPath   string `json:"subPath" gorm:"column:sub_path;type:varchar(255);default:'';comment:PVC内子路径"`
	MountPath string `json:"mountPath" gorm:"column:mount_path;type:varchar(512);not null;comment:容器内挂载路径"`
	ReadOnly  bool   `json:"readOnly" gorm:"column:read_only;default:true;comment:是否只读"`
}

func (TrainingJobMount) TableName() string {
	return "training_job_mounts"
}

// TrainingJobEnv 训练任务环境变量
type TrainingJobEnv struct {
	global.GVA_MODEL
	JobId uint   `json:"jobId" gorm:"column:job_id;not null;index;comment:关联的训练任务ID"`
	Name  string `json:"name" gorm:"column:name;type:varchar(255);not null;comment:环境变量名"`
	Value string `json:"value" gorm:"column:value;type:text;comment:环境变量值"`
}

func (TrainingJobEnv) TableName() string {
	return "training_job_envs"
}

package request

import (
	helper "gin-vue-admin/utils/k8s"

	corev1 "k8s.io/api/core/v1"
)

// CreateTrainingJobReq 创建训练任务请求
type CreateTrainingJobReq struct {
	Name              string                      `json:"name"`                              // 任务名称
	Namespace         string                      `json:"namespace"`                         // 命名空间（可选，默认使用用户namespace）
	FrameworkType     string                      `json:"frameworkType" binding:"required"`  // 框架类型
	ImageId           uint                        `json:"imageId" binding:"required"`        // 镜像ID
	StartupCommand    string                      `json:"startupCommand" binding:"required"` // 启动命令
	WorkerCount       int64                       `json:"workerCount"`                       // 兼容旧字段：实例数量
	InstanceCount     int64                       `json:"instanceCount"`                     // 新字段：实例数量
	ResourceId        uint                        `json:"resourceId"`                        // 兼容旧字段：产品ID
	TemplateProductId uint                        `json:"templateProductId"`                 // 新字段：模板商品ID
	ScheduleStrategy  string                      `json:"scheduleStrategy"`                  // 调度策略
	Mounts            []CreateTrainingJobMountReq `json:"mounts"`                            // 挂载配置
	Envs              []CreateTrainingJobEnvReq   `json:"envs"`                              // 环境变量
	// TensorBoard 配置
	EnableTensorboard  bool   `json:"enableTensorboard"`  // 是否启用TensorBoard
	TensorboardLogPath string `json:"tensorboardLogPath"` // TensorBoard日志路径（可选，默认为/workspace/logs）
	UserId             uint   `json:"userId"`             // 用户ID
	PayType            int64  `json:"payType"`            // 付费类型
}

// CreateTrainingJobMountReq 挂载配置请求
type CreateTrainingJobMountReq struct {
	MountType string `json:"mountType"` // 挂载类型
	SourceId  uint   `json:"sourceId"`  // 资源引用ID
	PvcId     uint   `json:"pvcId"`     // PVC ID
	PvcName   string `json:"pvcName"`   // PVC名称 (后端填充)
	SubPath   string `json:"subPath"`   // 子路径
	MountPath string `json:"mountPath"` // 容器内挂载路径
	ReadOnly  bool   `json:"readOnly"`  // 是否只读
}

// CreateTrainingJobEnvReq 环境变量请求
type CreateTrainingJobEnvReq struct {
	Name  string `json:"name"`  // 环境变量名
	Value string `json:"value"` // 环境变量值
}

// DeleteTrainingJobReq 删除训练任务请求
type DeleteTrainingJobReq struct {
	ID uint `json:"id" binding:"required"` // 任务ID
}

// GetTrainingJobListReq 获取训练任务列表请求
type GetTrainingJobListReq struct {
	Page      int    `json:"page" form:"page"`           // 页码
	PageSize  int    `json:"pageSize" form:"pageSize"`   // 每页数量
	Name      string `json:"name" form:"name"`           // 任务名称（模糊搜索）
	Status    string `json:"status" form:"status"`       // 状态筛选
	ClusterId uint   `json:"clusterId" form:"clusterId"` // 集群ID筛选
	UserId    uint   `json:"userId" form:"userId"`       // 用户ID
}

// GetTrainingJobDetailReq 获取训练任务详情请求
type GetTrainingJobDetailReq struct {
	ID uint `json:"id" form:"id" binding:"required"` // 任务ID
}

// GetTrainingJobPodsReq 获取训练任务Pod列表请求
type GetTrainingJobPodsReq struct {
	ID uint `json:"id" form:"id" binding:"required"` // 任务ID
}

// GetTrainingJobLogsReq 获取训练任务日志请求
type GetTrainingJobLogsReq struct {
	ID         uint   `json:"id" form:"id" binding:"required"` // 任务ID
	TaskName   string `json:"taskName" form:"taskName"`        // Task名称（master/worker等）
	PodIndex   *int   `json:"podIndex" form:"podIndex"`        // Pod索引
	Container  string `json:"container" form:"container"`      // 容器名称
	TailLines  int64  `json:"tailLines" form:"tailLines"`      // 尾部行数
	Follow     bool   `json:"follow" form:"follow"`            // 是否实时跟踪
	Timestamps bool   `json:"timestamps" form:"timestamps"`    // 是否显示时间戳
}

// StopTrainingJobReq 停止训练任务请求
type StopTrainingJobReq struct {
	ID uint `json:"id" binding:"required"` // 任务ID
}

// ProductSpec 产品规格别名，引用共享定义
type ProductSpec = helper.ProductSpec

// TrainingJobSpec 训练任务规格
type TrainingJobSpec struct {
	Name         string               // 任务名称
	Namespace    string               // 命名空间
	Framework    string               // 框架类型：PYTORCH_DDP, MPI, STANDALONE
	Image        string               // 镜像
	Command      []string             // 启动命令
	Args         []string             // 命令参数
	WorkerCount  int64                // Worker 数量
	AllowedNodes []string             // 允许调度的节点集合
	StrictSpread bool                 // 是否要求严格一节点一实例
	Product      ProductSpec          // 产品规格（builder 根据此构建资源请求）
	Volumes      []corev1.Volume      // 卷配置
	VolumeMounts []corev1.VolumeMount // 卷挂载配置
	Envs         []corev1.EnvVar      // 环境变量
	UseSHM       bool                 // 是否使用共享内存
	SHMSize      int64                // 共享内存大小
	Labels       map[string]string    // 标签
	Annotations  map[string]string    // 注解
	MaxRetry     int64                // 最大重试次数
}

type TerminalReq struct {
	ID        int64  `form:"id" binding:"required"`
	TaskName  string `form:"taskName"`
	Container string `form:"container"`
}

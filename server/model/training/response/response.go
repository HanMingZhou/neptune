package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response 通用响应结构
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{code, data, msg})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "成功", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

// TrainingJobItem 训练任务列表项
type TrainingJobItem struct {
	ID            uint       `json:"id"`            // 任务ID
	DisplayName   string     `json:"displayName"`   // 展示名称
	JobName       string     `json:"jobName"`       // K8s Job名称
	Namespace     string     `json:"namespace"`     // 命名空间
	FrameworkType string     `json:"frameworkType"` // 框架类型
	Image         string     `json:"image"`         // 镜像地址
	ImageName     string     `json:"imageName"`     // 镜像名称
	Status        string     `json:"status"`        // 状态
	TotalGPUCount int64      `json:"totalGpuCount"` // GPU总数
	WorkerCount   int64      `json:"workerCount"`   // Worker数量
	ClusterName   string     `json:"clusterName"`   // 集群名称
	CreatedAt     time.Time  `json:"createdAt"`     // 创建时间
	StartedAt     *time.Time `json:"startedAt"`     // 开始时间
	FinishedAt    *time.Time `json:"finishedAt"`    // 结束时间
	Duration      string     `json:"duration"`      // 运行时长
	ErrorMsg      string     `json:"errorMsg"`      // 错误信息
	// TensorBoard 和 Web Terminal
	TensorboardUrl string `json:"tensorboardUrl"` // TensorBoard访问地址
}

// GetTrainingJobListResp 获取训练任务列表响应
type GetTrainingJobListResp struct {
	List  []TrainingJobItem `json:"list"`  // 任务列表
	Total int64             `json:"total"` // 总数
}

// TrainingJobDetail 训练任务详情
type TrainingJobDetail struct {
	ID             uint                     `json:"id"`             // 任务ID
	DisplayName    string                   `json:"displayName"`    // 展示名称
	JobName        string                   `json:"jobName"`        // K8s Job名称
	Namespace      string                   `json:"namespace"`      // 命名空间
	ClusterID      uint                     `json:"clusterId"`      // 集群ID
	ClusterName    string                   `json:"clusterName"`    // 集群名称
	FrameworkType  string                   `json:"frameworkType"`  // 框架类型
	ImageId        uint                     `json:"imageId"`        // 镜像ID
	Image          string                   `json:"image"`          // 镜像地址
	ImageName      string                   `json:"imageName"`      // 镜像名称
	StartupCommand string                   `json:"startupCommand"` // 启动命令
	Status         string                   `json:"status"`         // 状态
	TotalGPUCount  int64                    `json:"totalGpuCount"`  // GPU总数
	GPUType        string                   `json:"gpuType"`        // GPU类型
	GPUModel       string                   `json:"gpuModel"`       // GPU型号
	ProductId      uint                     `json:"productId"`      // 产品ID
	WorkerCount    int64                    `json:"workerCount"`    // Worker数量
	CPU            int64                    `json:"cpu"`            // 每节点CPU(从产品获取)
	Memory         int64                    `json:"memory"`         // 每节点内存(从产品获取)
	WorkerGPU      int64                    `json:"workerGpu"`      // 每节点GPU数(从产品获取)
	K8sJobUid      string                   `json:"k8sJobUid"`      // K8s Job UID
	ErrorMsg       string                   `json:"errorMsg"`       // 错误信息
	CreatedAt      time.Time                `json:"createdAt"`      // 创建时间
	StartedAt      *time.Time               `json:"startedAt"`      // 开始时间
	FinishedAt     *time.Time               `json:"finishedAt"`     // 结束时间
	Duration       string                   `json:"duration"`       // 运行时长
	Mounts         []TrainingJobMountDetail `json:"mounts"`         // 挂载配置
	Envs           []TrainingJobEnvDetail   `json:"envs"`           // 环境变量
	Tasks          []TrainingJobTaskDetail  `json:"tasks"`          // 任务详情（从K8s获取）
	// TensorBoard
	EnableTensorboard  bool   `json:"enableTensorboard"`  // 是否启用TensorBoard
	TensorboardUrl     string `json:"tensorboardUrl"`     // TensorBoard访问地址
	TensorboardLogPath string `json:"tensorboardLogPath"` // TensorBoard日志路径
	// 计费与集群
	Area    string  `json:"area"`    // 地域
	PayType int64   `json:"payType"` // 付费类型
	Price   float64 `json:"price"`   // 单价
}

// TrainingJobMountDetail 挂载配置详情
type TrainingJobMountDetail struct {
	MountType string `json:"mountType"` // 挂载类型
	SourceId  uint   `json:"sourceId"`  // 资源引用ID
	PvcName   string `json:"pvcName"`   // PVC名称
	SubPath   string `json:"subPath"`   // 子路径
	MountPath string `json:"mountPath"` // 容器内挂载路径
	ReadOnly  bool   `json:"readOnly"`  // 是否只读
}

// TrainingJobEnvDetail 环境变量详情
type TrainingJobEnvDetail struct {
	Name  string `json:"name"`  // 环境变量名
	Value string `json:"value"` // 环境变量值
}

// TrainingJobTaskDetail 任务Task详情（从K8s获取）
type TrainingJobTaskDetail struct {
	TaskName string    `json:"taskName"` // Task名称
	Replicas int32     `json:"replicas"` // 副本数
	Status   string    `json:"status"`   // 状态
	Pods     []PodInfo `json:"pods"`     // Pod列表
}

// PodInfo Pod信息
type PodInfo struct {
	Name      string     `json:"name"`      // Pod名称
	Status    string     `json:"status"`    // Pod状态
	HostIP    string     `json:"hostIP"`    // 主机IP
	PodIP     string     `json:"podIP"`     // Pod IP
	StartTime *time.Time `json:"startTime"` // 启动时间
}

// CreateTrainingJobResp 创建训练任务响应
type CreateTrainingJobResp struct {
	ID          uint   `json:"id"`          // 任务ID
	DisplayName string `json:"displayName"` // 展示名称
	JobName     string `json:"jobName"`     // K8s Job名称
}

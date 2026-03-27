package podgroup

import (
	"time"

	"gorm.io/gorm"
)

// PodGroup Volcano PodGroup 资源的数据库模型
// 用于记录和监控 Volcano 调度的 PodGroup 状态，便于计费和资源管理
type PodGroup struct {
	gorm.Model

	// 基本信息
	Name         string `json:"name" gorm:"type:varchar(255);not null;index:idx_name_namespace,unique;comment:PodGroup名称"`
	Namespace    string `json:"namespace" gorm:"type:varchar(255);not null;index:idx_name_namespace,unique;comment:命名空间"`
	InstanceName string `json:"instanceName" gorm:"type:varchar(255);not null;index;comment:关联的实例名称(Notebook/Training等)"`
	InstanceType string `json:"instanceType" gorm:"type:varchar(50);not null;comment:实例类型(notebook/training/inference)"`

	// 关联信息
	OwnerID   uint   `json:"ownerId" gorm:"not null;index;comment:所属资源ID(NotebookID等)"`
	OwnerType string `json:"ownerType" gorm:"type:varchar(50);not null;comment:所属资源类型"`
	UserID    uint   `json:"userId" gorm:"not null;index;comment:用户ID"`
	ClusterID uint   `json:"clusterId" gorm:"not null;index;comment:集群ID"`

	// 调度信息
	MinMember         int32  `json:"minMember" gorm:"not null;default:1;comment:最小成员数"`
	Queue             string `json:"queue" gorm:"type:varchar(100);comment:队列名称"`
	PriorityClassName string `json:"priorityClassName" gorm:"type:varchar(100);comment:优先级类名"`

	// 状态信息
	Phase     string `json:"phase" gorm:"type:varchar(50);comment:调度阶段(Pending/Running/Unknown)"`
	Status    string `json:"status" gorm:"type:varchar(50);not null;default:'Creating';comment:状态(Creating/Running/Completed/Failed)"`
	Scheduled int32  `json:"scheduled" gorm:"default:0;comment:已调度的Pod数量"`
	Running   int32  `json:"running" gorm:"default:0;comment:运行中的Pod数量"`
	Succeeded int32  `json:"succeeded" gorm:"default:0;comment:成功的Pod数量"`
	Failed    int32  `json:"failed" gorm:"default:0;comment:失败的Pod数量"`

	// 资源信息（用于计费）
	CPU      string `json:"cpu" gorm:"type:varchar(50);comment:CPU配置"`
	Memory   string `json:"memory" gorm:"type:varchar(50);comment:内存配置"`
	GPU      string `json:"gpu" gorm:"type:varchar(50);comment:GPU数量"`
	GPUModel string `json:"gpuModel" gorm:"type:varchar(100);comment:GPU型号"`

	// 计费相关
	ProductID uint       `json:"productId" gorm:"index;comment:产品ID"`
	OrderID   uint       `json:"orderId" gorm:"index;comment:订单ID"`
	StartTime *time.Time `json:"startTime" gorm:"comment:开始运行时间(用于计费)"`
	EndTime   *time.Time `json:"endTime" gorm:"comment:结束时间(用于计费)"`

	// K8s 元数据
	UID             string `json:"uid" gorm:"type:varchar(255);index;comment:K8s UID"`
	ResourceVersion string `json:"resourceVersion" gorm:"type:varchar(50);comment:资源版本"`
	Generation      int64  `json:"generation" gorm:"comment:Generation"`

	// 备注
	Remark string `json:"remark" gorm:"type:text;comment:备注"`
}

// TableName 指定表名
func (PodGroup) TableName() string {
	return "podgroups"
}

// 查询资源的 product_id 和释放所需信息
type ResourceResult struct {
	ProductId     uint
	WorkerCount   int
	FrameworkType string
	DeployType    string
}

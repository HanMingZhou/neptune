package product

import (
	"time"

	"gorm.io/gorm"
)

// ResourceAllocation 资源分配记录
// 记录训练/推理实例在实际节点产品上的库存扣减明细。
type ResourceAllocation struct {
	ID uint `gorm:"primarykey"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	InstanceType string `json:"instanceType" gorm:"column:instance_type;size:32;not null;index:idx_allocation_owner;comment:实例类型(training/inference)"`
	InstanceID   uint   `json:"instanceId" gorm:"column:instance_id;not null;index:idx_allocation_owner;comment:实例ID"`

	ClusterID         uint   `json:"clusterId" gorm:"column:cluster_id;not null;index;comment:集群ID"`
	TemplateProductID uint   `json:"templateProductId" gorm:"column:template_product_id;not null;comment:模板商品ID"`
	ProductID         uint   `json:"productId" gorm:"column:product_id;not null;index;comment:实际扣减的节点产品ID"`
	NodeName          string `json:"nodeName" gorm:"column:node_name;size:100;not null;index;comment:节点名称"`

	ScheduleStrategy string `json:"scheduleStrategy" gorm:"column:schedule_strategy;size:20;not null;default:BALANCED;comment:调度策略"`
	ReplicaIndex     int    `json:"replicaIndex" gorm:"column:replica_index;not null;comment:副本索引"`
	TaskRole         string `json:"taskRole" gorm:"column:task_role;size:32;comment:任务角色(master/worker/head/standalone)"`
	ReservedCount    int64  `json:"reservedCount" gorm:"column:reserved_count;not null;default:1;comment:占用数量"`
}

func (ResourceAllocation) TableName() string {
	return "resource_allocations"
}

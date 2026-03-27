package pvc

import (
	"gin-vue-admin/model/cluster"
	"gin-vue-admin/model/product"
	"time"

	"gorm.io/gorm"
)

// Volume 数据盘，用于独立管理的存储资源
type Volume struct {
	ID         uint                `gorm:"primarykey"`
	CreatedAt  time.Time           `json:"createdAt"`
	UpdatedAt  time.Time           `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt      `json:"-" gorm:"uniqueIndex:idx_volume_name_ns_deleted"`
	Name       string              `json:"name" gorm:"column:name;comment:数据盘名称;uniqueIndex:idx_volume_name_ns_deleted"`
	Namespace  string              `json:"namespace" gorm:"column:namespace;comment:命名空间;uniqueIndex:idx_volume_name_ns_deleted"`
	Size       int64               `json:"size" gorm:"column:size;comment:大小"`
	Status     string              `json:"status" gorm:"column:status;comment:状态(Ready/Bound/Error)"`
	PVCName    string              `json:"pvcName" gorm:"column:pvc_name;comment:K8s PVC名称"`
	Type       int                 `json:"type" gorm:"column:type;comment:类型(1:dataset, 2:model);default:1"`
	UserId     uint                `json:"userId" gorm:"column:user_id;comment:用户ID;index"`
	K8sCluster *cluster.K8sCluster `gorm:"foreignKey:ClusterId;references:ID"`
	ClusterId  uint                `json:"clusterId" gorm:"column:cluster_id;comment:集群ID;index"`
	Product    *product.Product    `gorm:"foreignKey:ProductId;references:ID"`
	ProductId  uint                `json:"productId" gorm:"column:product_id;comment:产品ID;index"`
}

func (Volume) TableName() string {
	return "volumes"
}

var VolumeTypeToString = map[int]string{
	1: "dataset",
	2: "model",
}

var VolumeTypeToInt = map[string]int{
	"dataset": 1,
	"model":   2,
}

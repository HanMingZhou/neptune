package cluster

import "gorm.io/gorm"

// 集群状态常量
const (
	ClusterStatusEnabled  = 1 // 启用
	ClusterStatusDisabled = 0 // 停用
)

// K8sCluster 集群配置表 - 存储多集群的kubeconfig
type K8sCluster struct {
	gorm.Model
	Name         string `json:"name" gorm:"column:name;size:100;uniqueIndex;comment:集群名称"`
	Area         string `json:"area" gorm:"column:area;size:50;comment:地域区域"`
	Description  string `json:"description" gorm:"column:description;size:500;comment:描述"`
	KubeConfig   string `json:"kubeconfig" gorm:"column:kubeconfig;type:text;comment:kubeconfig内容"`
	ApiServer    string `json:"apiServer" gorm:"column:api_server;size:200;comment:API Server地址"`
	Status       int    `json:"status" gorm:"column:status;default:1;comment:状态(1-启用 0-停用)"`
	HarborAddr   string `json:"harborAddr" gorm:"column:harbor_addr;comment:Harbor地址"`
	StorageClass string `json:"storageClass" gorm:"column:storage_class;size:100;comment:K8s StorageClass名称"`
}

func (K8sCluster) TableName() string {
	return "k8s_clusters"
}

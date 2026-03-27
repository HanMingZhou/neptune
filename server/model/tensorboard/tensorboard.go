package tensorboard

import (
	"time"

	"gorm.io/gorm"
)

// Tensorboard 数据库模型，用于存储Tensorboard的元数据
type Tensorboard struct {
	gorm.Model
	InstanceName string `json:"instanceName" gorm:"column:instance_name;comment:实例名称"`
	OwnerType    string `json:"ownerType" gorm:"column:owner_type;comment:所有者类型"`
	OwnerID      uint   `json:"ownerId" gorm:"column:owner_id;comment:所有者ID"`
	Namespace    string `json:"namespace" gorm:"column:namespace;comment:命名空间"`
	LogsPath     string `json:"logsPath" gorm:"column:logs_path;comment:日志路径"`
	Status       string `json:"status" gorm:"column:status;comment:状态"`
	UserId       uint   `json:"userId" gorm:"column:user_id;comment:用户ID"`
	ClusterID    uint   `json:"clusterId" gorm:"column:cluster_id;comment:集群ID"`
}

// TensorBoardCRD 完整的TensorBoard CRD结构
type TensorBoardCRD struct {
	APIVersion string              `json:"apiVersion"`
	Kind       string              `json:"kind"`
	Metadata   TensorBoardMetadata `json:"metadata"`
	Spec       TensorBoardSpec     `json:"spec"`
	Status     TensorBoardStatus   `json:"status,omitempty"`
}

// TensorBoardMetadata TensorBoard元数据
type TensorBoardMetadata struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	Labels            map[string]string `json:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
	UID               string            `json:"uid,omitempty"`
	ResourceVersion   string            `json:"resourceVersion,omitempty"`
	CreationTimestamp time.Time         `json:"creationTimestamp,omitempty"`
}

// TensorBoardSpec TensorBoard规格
type TensorBoardSpec struct {
	LogsPath string `json:"logspath"`
}

// TensorBoardStatus TensorBoard状态
type TensorBoardStatus struct {
	Conditions    []TensorBoardCondition `json:"conditions,omitempty"`
	ReadyReplicas int32                  `json:"readyReplicas,omitempty"`
}

// TensorBoardCondition TensorBoard条件
type TensorBoardCondition struct {
	Type               string    `json:"type"`
	Status             string    `json:"status"`
	LastTransitionTime time.Time `json:"lastTransitionTime,omitempty"`
	Reason             string    `json:"reason,omitempty"`
	Message            string    `json:"message,omitempty"`
}

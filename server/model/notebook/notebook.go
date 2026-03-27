package notebook

import (
	"time"

	"gorm.io/gorm"
)

// Default 默认值
const (
	Workspace               = "workspace"
	DefaultNamespace        = "default"
	DefaultImage            = "kubeflownotebookswg/jupyter-scipy:v1.9.2"
	DefaultCPU              = "1"
	DefaultMemory           = "2Gi"
	DefaultWorkspaceSize    = 10
	DefaultWorkspacePath    = "/home/system"
	DefaultStorageClass     = "standard"
	DefaultTensorboardImage = "tensorflow/tensorflow:2.1.0"
	DefaultDataMountPath    = "/home/notebook/neptune"
)

// 卷类型常量
const (
	VolumeTypeDataset = "dataset"
	VolumeTypeModel   = "model"
)

// Notebook 数据库模型，用于存储Notebook的元数据
type Notebook struct {
	gorm.Model
	DisplayName        string           `json:"displayName" gorm:"column:display_name;comment:Notebook名称"`
	InstanceName       string           `json:"instanceName" gorm:"column:instance_name;comment:实例名称"`
	Namespace          string           `json:"namespace" gorm:"column:namespace;comment:命名空间"`
	Image              string           `json:"image" gorm:"column:image;comment:镜像地址"`
	ImageId            uint             `json:"imageId" gorm:"column:image_id;comment:镜像ID;index"`
	CPU                int64            `json:"cpu" gorm:"column:cpu;comment:CPU资源限制"`
	Memory             int64            `json:"memory" gorm:"column:memory;comment:内存资源限制"`
	GPU                int64            `json:"gpu" gorm:"column:gpu;comment:GPU资源限制"`
	GPUModel           string           `json:"gpuModel" gorm:"column:gpu_model;comment:GPU型号"`
	VGPUNumber         int64            `json:"vGpuNumber" gorm:"column:v_gpu_number;comment:vGPU数量;default:0"`
	VGPUMemory         int64            `json:"vGpuMemory" gorm:"column:v_gpu_memory;comment:vGPU显存 (GB);default:0"`
	VGPUCores          int64            `json:"vGpuCores" gorm:"column:v_gpu_cores;comment:vGPU核心数;default:0"`
	StorageSize        int64            `json:"storageSize" gorm:"column:storage_size;comment:存储大小"`
	Status             string           `json:"status" gorm:"column:status;comment:状态"`
	UserId             uint             `json:"userId" gorm:"column:user_id;comment:用户ID;index"`
	ClusterID          uint             `json:"clusterId" gorm:"column:cluster_id;comment:集群ID"`
	ProductId          uint             `json:"productId" gorm:"column:product_id;comment:产品ID;index"`
	PayType            int              `json:"payType" gorm:"column:pay_type;comment:付费类型(1-按量 2-包日 3-包周 4-包月)"`
	OrderId            uint             `json:"orderId" gorm:"column:order_id;comment:订单ID;index"`
	SSHKeyId           uint             `json:"sshKeyId" gorm:"column:ssh_key_id;comment:SSH密钥ID"`
	EnableSSHPassword  bool             `json:"enableSshPassword" gorm:"column:enable_ssh_password;comment:是否启用SSH密码登录"`
	SSHPassword        string           `json:"sshPassword" gorm:"column:ssh_password;comment:SSH登录密码"`
	EnableTensorboard  bool             `json:"enableTensorboard" gorm:"column:enable_tensorboard;comment:是否启用Tensorboard"`
	TensorboardLogPath string           `json:"tensorboardLogPath" gorm:"column:tensorboard_log_path;comment:Tensorboard日志路径"`
	TensorboardID      uint             `json:"tensorboardId" gorm:"column:tensorboard_id;comment:Tensorboard ID;index"`
	VolumeMounts       []NotebookVolume `json:"volumeMounts" gorm:"foreignKey:NotebookID"`
}

// NotebookVolume Notebook挂载卷表
type NotebookVolume struct {
	gorm.Model
	NotebookID uint   `json:"notebookId" gorm:"column:notebook_id;comment:Notebook ID;index"`
	Name       string `json:"name" gorm:"column:name;comment:卷名称"`
	MountsPath string `json:"mountsPath" gorm:"column:mounts_path;comment:挂载路径"`
	Size       int64  `json:"size" gorm:"column:size;comment:大小"`
	Type       string `json:"type" gorm:"column:type;comment:类型(dataset/model/workspace)"`
	PVCId      uint   `json:"pvcId" gorm:"column:pvc_id;comment:PVC ID"`
	PVCName    string `json:"pvcName" gorm:"column:pvc_name;comment:PVC名称"`
}

func (Notebook) TableName() string {
	return "notebooks"
}

// NotebookSpec Kubeflow Notebook CRD的Spec定义
type NotebookSpec struct {
	Template NotebookTemplateSpec `json:"template"`
}

// NotebookTemplateSpec Notebook模板规格
type NotebookTemplateSpec struct {
	Spec PodSpec `json:"spec"`
}

// PodSpec Pod规格
type PodSpec struct {
	Containers         []Container `json:"containers"`
	Volumes            []Volume    `json:"volumes,omitempty"`
	ServiceAccount     string      `json:"serviceAccount,omitempty"`
	ServiceAccountName string      `json:"serviceAccountName,omitempty"`
}

// Container 容器定义
type Container struct {
	Name            string               `json:"name"`
	Image           string               `json:"image"`
	ImagePullPolicy string               `json:"imagePullPolicy,omitempty"`
	Resources       ResourceRequirements `json:"resources,omitempty"`
	VolumeMounts    []VolumeMount        `json:"volumeMounts,omitempty"`
	Env             []EnvVar             `json:"env,omitempty"`
	Ports           []ContainerPort      `json:"ports,omitempty"`
}

// ResourceRequirements 资源需求
type ResourceRequirements struct {
	Limits   ResourceList `json:"limits,omitempty"`
	Requests ResourceList `json:"requests,omitempty"`
}

// ResourceList 资源列表
type ResourceList map[string]string

// VolumeMount 卷挂载
type VolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
}

// EnvVar 环境变量
type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ContainerPort 容器端口
type ContainerPort struct {
	ContainerPort int32  `json:"containerPort"`
	Name          string `json:"name,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
}

// Volume 卷定义
type Volume struct {
	Name                  string                             `json:"name"`
	PersistentVolumeClaim *PersistentVolumeClaimVolumeSource `json:"persistentVolumeClaim,omitempty"`
	EmptyDir              *EmptyDirVolumeSource              `json:"emptyDir,omitempty"`
}

// PersistentVolumeClaimVolumeSource PVC卷源
type PersistentVolumeClaimVolumeSource struct {
	ClaimName string `json:"claimName"`
}

// EmptyDirVolumeSource EmptyDir卷源
type EmptyDirVolumeSource struct {
	Medium    string `json:"medium,omitempty"`
	SizeLimit string `json:"sizeLimit,omitempty"`
}

// NotebookStatus Notebook状态
type NotebookStatus struct {
	Conditions     []NotebookCondition `json:"conditions,omitempty"`
	ReadyReplicas  int32               `json:"readyReplicas,omitempty"`
	ContainerState ContainerState      `json:"containerState,omitempty"`
}

// NotebookCondition Notebook条件
type NotebookCondition struct {
	Type               string    `json:"type"`
	Status             string    `json:"status"`
	LastTransitionTime time.Time `json:"lastTransitionTime,omitempty"`
	Reason             string    `json:"reason,omitempty"`
	Message            string    `json:"message,omitempty"`
}

// ContainerState 容器状态
type ContainerState struct {
	Waiting    *ContainerStateWaiting    `json:"waiting,omitempty"`
	Running    *ContainerStateRunning    `json:"running,omitempty"`
	Terminated *ContainerStateTerminated `json:"terminated,omitempty"`
}

// ContainerStateWaiting 容器等待状态
type ContainerStateWaiting struct {
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// ContainerStateRunning 容器运行状态
type ContainerStateRunning struct {
	StartedAt time.Time `json:"startedAt,omitempty"`
}

// ContainerStateTerminated 容器终止状态
type ContainerStateTerminated struct {
	ExitCode   int32     `json:"exitCode"`
	Signal     int32     `json:"signal,omitempty"`
	Reason     string    `json:"reason,omitempty"`
	Message    string    `json:"message,omitempty"`
	StartedAt  time.Time `json:"startedAt,omitempty"`
	FinishedAt time.Time `json:"finishedAt,omitempty"`
}

// NotebookCRD 完整的Notebook CRD结构
type NotebookCRD struct {
	APIVersion string           `json:"apiVersion"`
	Kind       string           `json:"kind"`
	Metadata   NotebookMetadata `json:"metadata"`
	Spec       NotebookSpec     `json:"spec"`
	Status     NotebookStatus   `json:"status,omitempty"`
}

// NotebookMetadata Notebook元数据
type NotebookMetadata struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	Labels            map[string]string `json:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
	UID               string            `json:"uid,omitempty"`
	ResourceVersion   string            `json:"resourceVersion,omitempty"`
	CreationTimestamp time.Time         `json:"creationTimestamp,omitempty"`
}

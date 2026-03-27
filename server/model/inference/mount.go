package inference

import "gin-vue-admin/global"

// InferenceMount 推理服务挂载配置
type InferenceMount struct {
	global.GVA_MODEL
	ServiceId uint   `json:"serviceId" gorm:"column:service_id;not null;index;comment:关联的推理服务ID"`
	MountType string `json:"mountType" gorm:"column:mount_type;type:varchar(20);comment:挂载类型(MODEL/DATA)"`
	PvcId     uint   `json:"pvcId" gorm:"column:pvc_id;comment:PVC ID"`
	PvcName   string `json:"pvcName" gorm:"column:pvc_name;type:varchar(255);comment:PVC名称"`
	SubPath   string `json:"subPath" gorm:"column:sub_path;type:varchar(255);default:'';comment:PVC内子路径"`
	MountPath string `json:"mountPath" gorm:"column:mount_path;type:varchar(512);not null;comment:容器内挂载路径"`
	ReadOnly  bool   `json:"readOnly" gorm:"column:read_only;default:true;comment:是否只读"`
}

func (InferenceMount) TableName() string {
	return "inference_mounts"
}

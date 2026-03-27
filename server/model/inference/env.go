package inference

import "gin-vue-admin/global"

// InferenceEnv 推理服务环境变量
type InferenceEnv struct {
	global.GVA_MODEL
	ServiceId uint   `json:"serviceId" gorm:"column:service_id;not null;index;comment:关联的推理服务ID"`
	Name      string `json:"name" gorm:"column:name;type:varchar(255);not null;comment:环境变量名"`
	Value     string `json:"value" gorm:"column:value;type:text;comment:环境变量值"`
}

func (InferenceEnv) TableName() string {
	return "inference_envs"
}

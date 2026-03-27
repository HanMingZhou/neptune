package inference

import "gin-vue-admin/global"

// InferenceApiKeyPolicy API Key 权限策略
// 支持一个 Key 访问多个服务，每个服务可配置不同权限
type InferenceApiKeyPolicy struct {
	global.GVA_MODEL
	ApiKeyId  uint   `json:"apiKeyId" gorm:"column:api_key_id;not null;index;comment:关联的API Key ID"`
	ServiceId uint   `json:"serviceId" gorm:"column:service_id;not null;index;comment:关联的服务ID(0表示所有服务)"`
	Actions   string `json:"actions" gorm:"column:actions;type:varchar(100);default:'inference';comment:允许的操作(inference,read,write)"`
}

func (InferenceApiKeyPolicy) TableName() string {
	return "inference_api_key_policies"
}

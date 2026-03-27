package inference

import (
	"gin-vue-admin/global"
	"time"
)

// InferenceApiKey API Key 实体
type InferenceApiKey struct {
	global.GVA_MODEL
	ServiceId   uint       `json:"serviceId" gorm:"column:service_id;not null;index;comment:关联的推理服务ID"`
	ApiKey      string     `json:"apiKey" gorm:"column:api_key;type:varchar(64);not null;uniqueIndex;comment:API Key"`
	Name        string     `json:"name" gorm:"column:name;type:varchar(100);comment:Key名称"`
	Description string     `json:"description" gorm:"column:description;type:varchar(500);comment:描述"`
	Status      string     `json:"status" gorm:"column:status;type:varchar(20);default:'active';comment:状态(active/disabled)"`
	Scopes      string     `json:"scopes" gorm:"column:scopes;type:varchar(100);default:'read,write';comment:权限范围(read,write)"`
	RateLimit   int        `json:"rateLimit" gorm:"column:rate_limit;default:0;comment:每分钟请求限制(0表示不限制)"`
	LastUsedAt  *time.Time `json:"lastUsedAt" gorm:"column:last_used_at;comment:最后使用时间"`
	ExpiredAt   *time.Time `json:"expiredAt" gorm:"column:expired_at;comment:过期时间"`
	UserId      uint       `json:"userId" gorm:"column:user_id;not null;index;comment:创建者用户ID"`
}

func (InferenceApiKey) TableName() string {
	return "inference_api_keys"
}

// API Key 状态常量
const (
	ApiKeyStatusActive   = "active"
	ApiKeyStatusDisabled = "disabled"
)

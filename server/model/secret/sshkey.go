package secret

import (
	"gorm.io/gorm"
)

// SSHKey SSH密钥表 - 用于存储用户的SSH公钥
type SSHKey struct {
	gorm.Model
	Name        string `json:"name" gorm:"column:name;size:100;comment:密钥名称"`
	UserId      uint   `json:"userId" gorm:"column:user_id;index;comment:用户ID"`
	PublicKey   string `json:"publicKey" gorm:"column:public_key;type:text;comment:公钥内容"`
	Fingerprint string `json:"fingerprint" gorm:"column:fingerprint;size:100;comment:密钥指纹"`
	IsDefault   bool   `json:"isDefault" gorm:"column:is_default;default:false;comment:是否默认密钥"`
}

func (SSHKey) TableName() string {
	return "ssh_keys"
}

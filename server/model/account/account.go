package account

import (
	"gin-vue-admin/global"
	"time"
)

const (
	// Phone/Email Status
	StatusBound = "已绑定"
	StatusUnset = "未设置"

	// Multi-factor Auth Status
	StatusMfaEnabled  = "已开启"
	StatusMfaDisabled = "未启用"

	// Third-party Binding Status
	StatusLinked   = "已关联"
	StatusUnlinked = "未关联"

	// Access Key Status
	StatusKeyGenerated    = "已生成"
	StatusKeyNotGenerated = "未生成"

	// Log Status
	LogStatusSuccess = "Success"
	LogStatusFailed  = "Failed"

	// Log Type
	LogTypeLogin     = "Login"
	LogTypeOperation = "Operation"

	// Log Method
	LogMethodPassword = "Password"
	LogMethodGithub   = "Github"

	// Bind Type
	BindTypePhone = 1
	BindTypeEmail = 2
)

type AccountAccessLog struct {
	global.GVA_MODEL
	UserId    uint      `json:"userId" gorm:"index;comment:用户ID"`
	Ip        string    `json:"ip" gorm:"comment:登录IP"`
	Location  string    `json:"location" gorm:"comment:IP归属地"`
	Device    string    `json:"device" gorm:"comment:登录设备"`
	Browser   string    `json:"browser" gorm:"comment:浏览器"`
	Os        string    `json:"os" gorm:"comment:操作系统"`
	Method    string    `json:"method" gorm:"comment:登录方式(密码/Github/Wechat等)"`
	Status    string    `json:"status" gorm:"comment:状态"` // constant.LogStatus*
	Reason    string    `json:"reason" gorm:"comment:失败原因"`
	LogType   string    `json:"logType" gorm:"comment:日志类型"` // constant.LogType*
	LoginTime time.Time `json:"loginTime" gorm:"comment:登录时间"`
}

func (AccountAccessLog) TableName() string {
	return "sys_account_access_logs"
}

type AccountSecurity struct {
	global.GVA_MODEL
	UserId          uint   `json:"userId" gorm:"uniqueIndex;comment:用户ID"`
	MfaEnabled      bool   `json:"mfaEnabled" gorm:"default:false;comment:是否开启MFA"`
	MfaSecret       string `json:"-" gorm:"comment:MFA密钥"`
	GithubId        string `json:"githubId" gorm:"index;comment:Github ID"`
	GithubUsername  string `json:"githubUsername" gorm:"comment:Github用户名"`
	AccessKeyId     string `json:"accessKeyId" gorm:"uniqueIndex;comment:AccessKey ID"`
	AccessKeySecret string `json:"-" gorm:"comment:AccessKey Secret"`
}

func (AccountSecurity) TableName() string {
	return "sys_account_securities"
}

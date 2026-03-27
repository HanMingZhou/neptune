package request

import (
	"gin-vue-admin/model/common/request"
)

type GetAccessLogListReq struct {
	request.PageInfo
	Status string `json:"status" form:"status"` // 成功/失败
	Ip     string `json:"ip" form:"ip"`         // IP搜索
	Device string `json:"device" form:"device"` // 设备搜索
}

type GetSecurityStatusReq struct {
	// 不需要字段，从上下文获取
}

type BindAccountReq struct {
	Type  int64  `json:"type" binding:"required"`  // 1:phone / 2:email
	Value string `json:"value" binding:"required"` // 手机号或邮箱
	Code  string `json:"code" binding:"required"`  // 验证码
}

type UpdatePasswordReq struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

type ToggleMfaReq struct {
	Enabled bool   `json:"enabled"`
	Code    string `json:"code" binding:"required"` // 关闭MFA时需验证TOTP码
}

// MFA初始化请求（生成密钥和二维码）
type MfaSetupReq struct{}

// MFA激活请求（扫码后输入验证码确认绑定）
type MfaActivateReq struct {
	Code string `json:"code" binding:"required"` // 6位TOTP验证码
}

// MFA二次登录验证请求
type MfaLoginReq struct {
	MfaToken string `json:"mfaToken" binding:"required"` // 临时令牌
	Code     string `json:"code" binding:"required"`     // 6位TOTP验证码
}

type BindGithubReq struct {
	Code string `json:"code" binding:"required"` // OAuth 授权码
}

type GenerateAccessKeyReq struct {
	// 按需添加
}

type KillSessionReq struct {
	LogId uint `json:"logId" binding:"required"`
}

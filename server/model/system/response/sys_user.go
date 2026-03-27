package response

import (
	"gin-vue-admin/model/system"
)

type UserResponse struct {
	ID          uint                  `json:"id"`
	UUID        string                `json:"uuid"`
	Username    string                `json:"userName"`
	NickName    string                `json:"nickName"`
	HeaderImg   string                `json:"headerImg"`
	AuthorityId uint                  `json:"authorityId"`
	Authority   system.SysAuthority   `json:"authority"`
	Authorities []system.SysAuthority `json:"authorities"`
	Phone       string                `json:"phone"`
	Email       string                `json:"email"`
	Enable      int                   `json:"enable"`
	Namespace   string                `json:"namespace"`
}

type SysUserResponse struct {
	User UserResponse `json:"user"`
}

type LoginResponse struct {
	User      UserResponse `json:"user"`
	Token     string       `json:"token"`
	ExpiresAt int64        `json:"expiresAt"`
	NeedMfa   bool         `json:"needMfa"`  // 是否需要MFA二次验证
	MfaToken  string       `json:"mfaToken"` // MFA临时令牌（needMfa为true时返回）
}

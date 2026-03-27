package request

// AddSSHKeyReq 创建SSH密钥请求
type AddSSHKeyReq struct {
	Name      string `json:"name" form:"name" binding:"required"`           // 密钥名称
	PublicKey string `json:"publicKey" form:"publicKey" binding:"required"` // 公钥内容
}

// GetSSHKeyListReq 获取SSH密钥列表请求
type GetSSHKeyListReq struct {
	Page     int    `json:"page" form:"page"`         // 页码
	PageSize int    `json:"pageSize" form:"pageSize"` // 每页数量
	Name     string `json:"name" form:"name"`         // 搜索名称
}

// DeleteSSHKeyReq 删除SSH密钥请求
type DeleteSSHKeyReq struct {
	ID int64 `json:"id" form:"id" binding:"required"` // 密钥ID
}

// SetDefaultSSHKeyReq 设置默认密钥请求
type SetDefaultSSHKeyReq struct {
	ID int64 `json:"id" form:"id" binding:"required"` // 密钥ID
}

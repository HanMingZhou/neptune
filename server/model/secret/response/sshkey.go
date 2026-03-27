package response

// SSHKeyItem SSH密钥列表项
type SSHKeyItem struct {
	ID          uint   `json:"id"`          // 密钥ID
	Name        string `json:"name"`        // 密钥名称
	Fingerprint string `json:"fingerprint"` // 密钥指纹
	IsDefault   bool   `json:"isDefault"`   // 是否默认密钥
	CreatedAt   string `json:"createdAt"`   // 创建时间
}

// GetSSHKeyListResp 获取SSH密钥列表响应
type GetSSHKeyListResp struct {
	List  []SSHKeyItem `json:"list"`  // 密钥列表
	Total int64        `json:"total"` // 总数
}

// AddSSHKeyResp 创建SSH密钥响应
type AddSSHKeyResp struct {
	ID          uint   `json:"id"`          // 密钥ID
	Fingerprint string `json:"fingerprint"` // 密钥指纹
}

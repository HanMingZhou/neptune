package request

import "gorm.io/gorm"

// PageInfo Paging common input parameter structure
type PageInfo struct {
	Page     int  `json:"page" form:"page"`         // 页码,默认1
	PageSize int  `json:"pageSize" form:"pageSize"` // 每页大小,默认10
	Id       uint `json:"id" form:"id"`             // 关键字
}

func (r *PageInfo) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.Page <= 0 {
			r.Page = 1
		}
		switch {
		case r.PageSize > 100:
			r.PageSize = 100
		case r.PageSize <= 0:
			r.PageSize = 10
		}
		offset := (r.Page - 1) * r.PageSize
		return db.Offset(offset).Limit(r.PageSize)
	}
}

// 获取分类列表
type GetPipeListReq struct {
	PageInfo
	Name      string `json:"name" form:"name"`
	Type      string `json:"type" form:"type"`
	Status    string `json:"status" form:"status"`
	CreatedBy string `json:"createdBy" form:"createdBy"`
}

type AddPipeReq struct {
	InstanceName         string            `json:"instanceName" form:"instanceName" binding:"required"` // Pipe 代理POD名称
	Namespace            string            `json:"namespace" form:"namespace"`                          // Pipe 所在的 Namespace (通常是 sshpiper-system)
	TargetNamespace      string            `json:"targetNamespace" form:"targetNamespace"`              // 目标 Notebook 所在的 Namespace
	UserSSHKey           string            `json:"userSSHKey" form:"userSSHKey"`                        // 用户的 SSH 公钥（公钥登录时使用）
	PrivateKeySecretName string            `json:"privateKeySecretName" form:"privateKeySecretName"`    // SSHPiper 连接后端使用的私钥（公钥登录时使用）
	EnablePasswordAuth   bool              `json:"enablePasswordAuth" form:"enablePasswordAuth"`        // 是否启用密码登录（密码透传给后端验证）
	TargetHost           string            `json:"targetHost" form:"targetHost"`                        // 目标后端地址 (可选，默认使用 InstanceName)
	TargetUsername       string            `json:"targetUsername" form:"targetUsername"`                // 后端容器内的用户名 (默认 root，可选 jovyan 等)
	Labels               map[string]string `json:"labels" form:"labels"`                                // Pipe 的标签
}

type UpdatePipeReq struct {
	Id           uint              `json:"id" form:"id"`
	InstanceName string            `json:"instanceName" form:"instanceName" binding:"required"`
	Namespace    string            `json:"namespace" form:"namespace"`
	Labels       map[string]string `json:"labels" form:"labels"`
}

type DeletePipeReq struct {
	Id           uint   `json:"id" form:"id"`
	InstanceName string `json:"instanceName" form:"instanceName" binding:"required"`
	Namespace    string `json:"namespace" form:"namespace"`
}

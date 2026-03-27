package request

// GetClusterListReq 获取集群列表请求
type GetClusterListReq struct {
	Keyword string `json:"keyword" form:"keyword"` // 搜索关键词
	Status  *int   `json:"status" form:"status"`   // 状态筛选
}

// CreateClusterReq 创建集群请求
type CreateClusterReq struct {
	Name         string `json:"name" form:"name" binding:"required"` // 集群名称
	Area         string `json:"area" form:"area"`                    // 地域区域
	Description  string `json:"description" form:"description"`      // 描述
	KubeConfig   string `json:"kubeconfig" form:"kubeconfig"`        // kubeconfig内容
	ApiServer    string `json:"apiServer" form:"apiServer"`          // API Server地址
	Status       int    `json:"status" form:"status"`                // 状态
	HarborAddr   string `json:"harborAddr" form:"harborAddr"`        // Harbor地址
	StorageClass string `json:"storageClass" form:"storageClass"`    // StorageClass名称
}

// UpdateClusterReq 更新集群请求
type UpdateClusterReq struct {
	ID           uint   `json:"id" form:"id" binding:"required"`  // 集群ID
	Name         string `json:"name" form:"name"`                 // 集群名称
	Area         string `json:"area" form:"area"`                 // 地域区域
	Description  string `json:"description" form:"description"`   // 描述
	KubeConfig   string `json:"kubeconfig" form:"kubeconfig"`     // kubeconfig内容
	ApiServer    string `json:"apiServer" form:"apiServer"`       // API Server地址
	Status       int    `json:"status" form:"status"`             // 状态
	HarborAddr   string `json:"harborAddr" form:"harborAddr"`     // Harbor地址
	StorageClass string `json:"storageClass" form:"storageClass"` // StorageClass名称
}

// DeleteClusterReq 删除集群请求
type DeleteClusterReq struct {
	ID uint `json:"id" form:"id" binding:"required"` // 集群ID
}

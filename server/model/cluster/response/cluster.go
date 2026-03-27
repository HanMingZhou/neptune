package response

// ClusterItem 集群列表项
type ClusterItem struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Area         string `json:"area"`
	Description  string `json:"description"`
	KubeConfig   string `json:"kubeConfig"`
	ApiServer    string `json:"apiServer"`
	Status       int    `json:"status"`
	HarborAddr   string `json:"harborAddr"`
	StorageClass string `json:"storageClass"`
	NodeCount    int    `json:"nodeCount"`  // 节点总数
	InternalIP   string `json:"internalIp"` // 内网IP（Master节点）
	CreatedAt    string `json:"createdAt"`
}

// ClusterListResponse 集群列表响应
type ClusterListResponse struct {
	List  []ClusterItem `json:"list"`
	Total int64         `json:"total"`
}

package request

// UncordonNodeReq 恢复节点调度请求
type UncordonNodeReq struct {
	ClusterId uint   `json:"clusterId" binding:"required"`
	NodeName  string `json:"nodeName" binding:"required"`
}

// DrainNodeReq 驱逐节点请求
type DrainNodeReq struct {
	ClusterId uint   `json:"clusterId" binding:"required"`
	NodeName  string `json:"nodeName" binding:"required"`
}

// GetClusterNodesReq 获取集群节点请求
type GetClusterNodesReq struct {
	ClusterId uint   `json:"clusterId" binding:"required"`
	Keyword   string `json:"keyword"`
	CPU       int64  `json:"cpu"`      // 产品CPU规格
	Memory    int64  `json:"memory"`   // 产品内存规格
	GPUCount  int64  `json:"gpuCount"` // 产品GPU规格
}

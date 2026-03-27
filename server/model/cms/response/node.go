package response

type NodeInfoResponse struct {
	NodeName          string `json:"nodeName"`
	InternalIP        string `json:"internalIp"`
	ClusterName       string `json:"clusterName"` // 所属集群名称
	NodeRole          string `json:"nodeRole"`    // 节点角色: master / worker
	Schedulable       bool   `json:"schedulable"`
	CPUAllocatable    int64  `json:"cpuAllocatable"`
	CPUAvailable      int64  `json:"cpuAvailable"`
	CPUModel          string `json:"cpuModel"` // CPU型号
	MemoryAllocatable int64  `json:"memoryAllocatable"`
	MemoryAvailable   int64  `json:"memoryAvailable"`
	GPUModel          string `json:"gpuModel"`
	GPUMemory         int64  `json:"gpuMemory"`
	GPUCount          int64  `json:"gpuCount"`
	GPUAvailable      int64  `json:"gpuAvailable"`
	VGPUNumber        int64  `json:"vGpuNumber"` // volcano.sh/vgpu-number
	VGPUMemory        int64  `json:"vGpuMemory"` // volcano.sh/vgpu-memory
	VGPUCores         int64  `json:"vGpuCores"`  // volcano.sh/vgpu-cores
	MaxInstances      int64  `json:"maxInstances"`
	Area              string `json:"area"`
}

type NodeListResponse struct {
	Nodes []NodeInfoResponse `json:"nodes"`
	Total int                `json:"total"`
}

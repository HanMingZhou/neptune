package request

type ProductPriceItem struct {
	PriceType int     `json:"priceType"`
	Price     float64 `json:"price"`
}

// CreateProductReq 创建产品请求（对应一个K8s Node节点）
type CreateProductReq struct {
	ProductType   int                `json:"productType"` // 1-计算资源 2-文件存储
	Name          string             `json:"name" binding:"required"`
	Description   string             `json:"description"`
	ClusterId     uint               `json:"clusterId" binding:"required"`
	Area          string             `json:"area" binding:"required"`
	NodeName      string             `json:"nodeName"` // 计算资源必填，文件存储可选
	NodeType      string             `json:"nodeType"`
	CPUModel      string             `json:"cpuModel"`
	CPU           int64              `json:"cpu"`
	Memory        int64              `json:"memory"`
	GPUModel      string             `json:"gpuModel"`
	GPUCount      int64              `json:"gpuCount"` // GPU总卡数
	GPUMemory     int64              `json:"gpuMemory"`
	VGPUNumber    int64              `json:"vGpuNumber"`
	VGPUMemory    int64              `json:"vGpuMemory"`
	VGPUCores     int64              `json:"vGpuCores"`
	Prices        []ProductPriceItem `json:"prices"`
	DriverVersion string             `json:"driverVersion"`
	CUDAVersion   string             `json:"cudaVersion"`
	SystemDisk    int64              `json:"systemDisk"`
	DataDisk      int64              `json:"dataDisk"`
	Status        int                `json:"status"`
	SortOrder     int                `json:"sortOrder"`
	MaxInstances  int64              `json:"maxInstances"`
	// 文件存储字段
	StorageClass   string  `json:"storageClass"`
	StoragePriceGB float64 `json:"storagePriceGb"`
}

// BatchCreateComputeProductReq 批量创建计算产品请求
type BatchCreateComputeProductReq struct {
	ProductType   int                `json:"productType"` // 默认计算产品
	NodeNames     []string           `json:"nodeNames" binding:"required,min=1"`
	Name          string             `json:"name" binding:"required"`
	Description   string             `json:"description"`
	ClusterId     uint               `json:"clusterId" binding:"required"`
	Area          string             `json:"area" binding:"required"`
	NodeType      string             `json:"nodeType"`
	CPUModel      string             `json:"cpuModel"`
	CPU           int64              `json:"cpu"`
	Memory        int64              `json:"memory"`
	GPUModel      string             `json:"gpuModel"`
	GPUCount      int64              `json:"gpuCount"`
	GPUMemory     int64              `json:"gpuMemory"`
	VGPUNumber    int64              `json:"vGpuNumber"`
	VGPUMemory    int64              `json:"vGpuMemory"`
	VGPUCores     int64              `json:"vGpuCores"`
	Prices        []ProductPriceItem `json:"prices"`
	DriverVersion string             `json:"driverVersion"`
	CUDAVersion   string             `json:"cudaVersion"`
	SystemDisk    int64              `json:"systemDisk"`
	DataDisk      int64              `json:"dataDisk"`
	Status        int                `json:"status"`
	SortOrder     int                `json:"sortOrder"`
	MaxInstances  int64              `json:"maxInstances"`
}

// UpdateProductReq 更新产品请求
type UpdateProductReq struct {
	ID            uint               `json:"id" binding:"required"`
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	Area          string             `json:"area"`
	NodeType      string             `json:"nodeType"`
	CPUModel      string             `json:"cpuModel"`
	CPU           int64              `json:"cpu"`
	Memory        int64              `json:"memory"`
	GPUModel      string             `json:"gpuModel"`
	GPUCount      int64              `json:"gpuCount"`
	GPUMemory     int64              `json:"gpuMemory"`
	VGPUNumber    int64              `json:"vGpuNumber"`
	VGPUMemory    int64              `json:"vGpuMemory"`
	VGPUCores     int64              `json:"vGpuCores"`
	Prices        []ProductPriceItem `json:"prices"`
	DriverVersion string             `json:"driverVersion"`
	CUDAVersion   string             `json:"cudaVersion"`
	SystemDisk    int64              `json:"systemDisk"`
	DataDisk      int64              `json:"dataDisk"`
	Status        int                `json:"status"`
	SortOrder     int                `json:"sortOrder"`
	MaxInstances  int64              `json:"maxInstances"`
	// 文件存储字段
	StorageClass   string  `json:"storageClass"`
	StoragePriceGB float64 `json:"storagePriceGb"`
}

// UpdatePriceReq 更新价格请求
type UpdatePriceReq struct {
	ID     uint               `json:"id" binding:"required"`
	Prices []ProductPriceItem `json:"prices"`
}

// DeleteProductReq 删除产品请求
type DeleteProductReq struct {
	ID uint `json:"id" binding:"required"`
}

// GetProductListReq 获取产品列表请求
type GetProductListReq struct {
	Page            int      `json:"page" form:"page"`
	PageSize        int      `json:"pageSize" form:"pageSize"`
	ProductType     int      `json:"productType" form:"productType"` // 1-计算资源 2-文件存储
	ClusterId       uint     `json:"clusterId" form:"clusterId"`
	Area            string   `json:"area" form:"area"`
	ResourceType    string   `json:"resourceType" form:"resourceType"`
	GPUModel        string   `json:"gpuModel" form:"gpuModel"`
	Status          *int     `json:"status" form:"status"`
	AvailableMin    *int64   `json:"availableMin" form:"availableMin"`
	AvailableMax    *int64   `json:"availableMax" form:"availableMax"`
	MaxInstancesMin *int64   `json:"maxInstancesMin" form:"maxInstancesMin"`
	MaxInstancesMax *int64   `json:"maxInstancesMax" form:"maxInstancesMax"`
	UsedCapacityMin *int64   `json:"usedCapacityMin" form:"usedCapacityMin"`
	UsedCapacityMax *int64   `json:"usedCapacityMax" form:"usedCapacityMax"`
	PriceType       int      `json:"priceType" form:"priceType"`
	PriceMin        *float64 `json:"priceMin" form:"priceMin"`
	PriceMax        *float64 `json:"priceMax" form:"priceMax"`
	Keyword         string   `json:"keyword" form:"keyword"`
}

// GetProductNodeCandidatesReq 获取产品节点候选列表请求
type GetProductNodeCandidatesReq struct {
	ClusterId        uint   `json:"clusterId" form:"clusterId" binding:"required"`
	ResourceType     string `json:"resourceType" form:"resourceType"`
	CPU              int64  `json:"cpu" form:"cpu"`
	Memory           int64  `json:"memory" form:"memory"`
	GPUCount         int64  `json:"gpuCount" form:"gpuCount"`
	GPUMemory        int64  `json:"gpuMemory" form:"gpuMemory"`
	VGPUNumber       int64  `json:"vGpuNumber" form:"vGpuNumber"`
	VGPUMemory       int64  `json:"vGpuMemory" form:"vGpuMemory"`
	VGPUCores        int64  `json:"vGpuCores" form:"vGpuCores"`
	ExcludeProductId uint   `json:"excludeProductId" form:"excludeProductId"`
}

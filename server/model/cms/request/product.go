package request

// CreateProductReq 创建产品请求（对应一个K8s Node节点）
type CreateProductReq struct {
	ProductType   int     `json:"productType"` // 1-计算资源 2-文件存储
	Name          string  `json:"name" binding:"required"`
	Description   string  `json:"description"`
	ClusterId     uint    `json:"clusterId" binding:"required"`
	Area          string  `json:"area" binding:"required"`
	NodeName      string  `json:"nodeName"` // 计算资源必填，文件存储可选
	NodeType      string  `json:"nodeType"`
	CPUModel      string  `json:"cpuModel"`
	CPU           int64   `json:"cpu"`
	Memory        int64   `json:"memory"`
	GPUModel      string  `json:"gpuModel"`
	GPUCount      int64   `json:"gpuCount"` // GPU总卡数
	GPUMemory     int64   `json:"gpuMemory"`
	VGPUNumber    int64   `json:"vGpuNumber"`
	VGPUMemory    int64   `json:"vGpuMemory"`
	VGPUCores     int64   `json:"vGpuCores"`
	PriceHourly   float64 `json:"priceHourly"`
	PriceDaily    float64 `json:"priceDaily"`
	PriceWeekly   float64 `json:"priceWeekly"`
	PriceMonthly  float64 `json:"priceMonthly"`
	DriverVersion string  `json:"driverVersion"`
	CUDAVersion   string  `json:"cudaVersion"`
	SystemDisk    int64   `json:"systemDisk"`
	DataDisk      int64   `json:"dataDisk"`
	Status        int     `json:"status"`
	MaxInstances  int64   `json:"maxInstances"`
	// 文件存储字段
	StorageClass   string  `json:"storageClass"`
	StoragePriceGB float64 `json:"storagePriceGb"`
}

// UpdateProductReq 更新产品请求
type UpdateProductReq struct {
	ID            uint    `json:"id" binding:"required"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Area          string  `json:"area"`
	NodeType      string  `json:"nodeType"`
	CPUModel      string  `json:"cpuModel"`
	CPU           int64   `json:"cpu"`
	Memory        int64   `json:"memory"`
	GPUModel      string  `json:"gpuModel"`
	GPUCount      int64   `json:"gpuCount"`
	GPUMemory     int64   `json:"gpuMemory"`
	VGPUNumber    int64   `json:"vGpuNumber"`
	VGPUMemory    int64   `json:"vGpuMemory"`
	VGPUCores     int64   `json:"vGpuCores"`
	PriceHourly   float64 `json:"priceHourly"`
	PriceDaily    float64 `json:"priceDaily"`
	PriceWeekly   float64 `json:"priceWeekly"`
	PriceMonthly  float64 `json:"priceMonthly"`
	DriverVersion string  `json:"driverVersion"`
	CUDAVersion   string  `json:"cudaVersion"`
	SystemDisk    int64   `json:"systemDisk"`
	DataDisk      int64   `json:"dataDisk"`
	Status        int     `json:"status"`
	MaxInstances  int64   `json:"maxInstances"`
	// 文件存储字段
	StorageClass   string  `json:"storageClass"`
	StoragePriceGB float64 `json:"storagePriceGb"`
}

// UpdatePriceReq 更新价格请求
type UpdatePriceReq struct {
	ID           uint    `json:"id" binding:"required"`
	PriceHourly  float64 `json:"priceHourly"`
	PriceDaily   float64 `json:"priceDaily"`
	PriceWeekly  float64 `json:"priceWeekly"`
	PriceMonthly float64 `json:"priceMonthly"`
}

// DeleteProductReq 删除产品请求
type DeleteProductReq struct {
	ID uint `json:"id" binding:"required"`
}

// GetProductListReq 获取产品列表请求
type GetProductListReq struct {
	Page        int    `json:"page" form:"page"`
	PageSize    int    `json:"pageSize" form:"pageSize"`
	ProductType int    `json:"productType" form:"productType"` // 1-计算资源 2-文件存储
	ClusterId   uint   `json:"clusterId" form:"clusterId"`
	Area        string `json:"area" form:"area"`
	GPUModel    string `json:"gpuModel" form:"gpuModel"`
	Status      *int   `json:"status" form:"status"`
	Keyword     string `json:"keyword" form:"keyword"`
}

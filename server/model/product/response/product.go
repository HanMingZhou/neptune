package response

type ProductFiltersResponse struct {
	Areas      []string        `json:"areas"`
	GPUModels  []GPUModelInfo  `json:"gpuModels"`
	VGPUModels []VGPUModelInfo `json:"vgpuModels"`
	CPUModels  []CPUModelInfo  `json:"cpuModels"`
	CPUOnly    CPUInfo         `json:"cpuOnly"`
}

type CPUModelInfo struct {
	Model     string `json:"model"`
	Available int    `json:"available"`
	Total     int    `json:"total"`
}

type GPUModelInfo struct {
	Model     string `json:"model"`
	Available int    `json:"available"`
	Total     int    `json:"total"`
}

type VGPUModelInfo struct {
	Memory    string `json:"memory"`
	Available int    `json:"available"`
	Total     int    `json:"total"`
}

type CPUInfo struct {
	Available int `json:"available"`
	Total     int `json:"total"`
}

type ProductListResponse struct {
	List  []ProductDetailResponse `json:"list"`
	Total int64                   `json:"total"`
}

type ProductDetailResponse struct {
	ID             uint    `json:"id"`
	ProductType    int     `json:"productType"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Area           string  `json:"area"`
	NodeName       string  `json:"nodeName"`
	NodeType       string  `json:"nodeType"`
	CPUModel       string  `json:"cpuModel"`
	CPU            int64   `json:"cpu"`
	Memory         int64   `json:"memory"`
	GPUModel       string  `json:"gpuModel"`
	GPUCount       int64   `json:"gpuCount"`
	GPUMemory      int64   `json:"gpuMemory"`
	VGPUNumber     int64   `json:"vGpuNumber"`
	VGPUMemory     int64   `json:"vGpuMemory"`
	VGPUCores      int64   `json:"vGpuCores"`
	PriceHourly    float64 `json:"priceHourly"`
	PriceDaily     float64 `json:"priceDaily"`
	PriceWeekly    float64 `json:"priceWeekly"`
	PriceMonthly   float64 `json:"priceMonthly"`
	DriverVersion  string  `json:"driverVersion"`
	CUDAVersion    string  `json:"cudaVersion"`
	SystemDisk     int64   `json:"systemDisk"`
	DataDisk       int64   `json:"dataDisk"`
	Status         int     `json:"status"`
	StorageClass   string  `json:"storageClass"`
	StoragePriceGB float64 `json:"storagePriceGb"`
	MaxInstances   int64   `json:"maxInstances"`
	Available      int64   `json:"available"` // 可用库存
	ClusterId      uint    `json:"clusterId"`
}

type ClusterResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Area        string `json:"area"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

type NodeInfoResponse struct {
	Area            string `json:"area"`
	NodeName        string `json:"nodeName"`
	CPUModel        string `json:"cpuModel"`
	CPU             int64  `json:"cpu"`
	Memory          int64  `json:"memory"`
	GPUModel        string `json:"gpuModel"`
	GPUCount        int64  `json:"gpuCount"`
	GPUMemory       int64  `json:"gpuMemory"`
	Schedulable     bool   `json:"schedulable"`
	CPUAvailable    int64  `json:"cpuAvailable"`
	MemoryAvailable int64  `json:"memoryAvailable"`
	GPUAvailable    int64  `json:"gpuAvailable"`
	MaxInstances    int64  `json:"maxInstances"`
}

type ReserveResult struct {
	ResourceCount int64                 `json:"resourceCount"` // 实际锁定的资源数量
	Product       ProductDetailResponse `json:"product"`       // 产品信息（用于获取 CPU/Memory/GPU 等配置）
}

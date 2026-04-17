package response

type ProductListResponse struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type ExistingComputeProductSummary struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	ResourceType string  `json:"resourceType"`
	CPU          int64   `json:"cpu"`
	Memory       int64   `json:"memory"`
	GPUModel     string  `json:"gpuModel"`
	GPUCount     int64   `json:"gpuCount"`
	GPUMemory    int64   `json:"gpuMemory"`
	VGPUNumber   int64   `json:"vGpuNumber"`
	VGPUMemory   int64   `json:"vGpuMemory"`
	VGPUCores    int64   `json:"vGpuCores"`
	PriceHourly  float64 `json:"priceHourly"`
	PriceDaily   float64 `json:"priceDaily"`
	PriceWeekly  float64 `json:"priceWeekly"`
	PriceMonthly float64 `json:"priceMonthly"`
	Status       int     `json:"status"`
	MaxInstances int64   `json:"maxInstances"`
	UsedCapacity int64   `json:"usedCapacity"`
	Available    int64   `json:"available"`
}

type ProductNodeCandidateResponse struct {
	NodeInfoResponse
	ExistingComputeProducts []ExistingComputeProductSummary `json:"existingComputeProducts"`
	CanCreateComputeProduct bool                            `json:"canCreateComputeProduct"`
	Compatible              bool                            `json:"compatible"`
	DisableReason           string                          `json:"disableReason,omitempty"`
}

type BatchCreateComputeProductResponse struct {
	CreatedIds   []uint   `json:"createdIds"`
	CreatedCount int      `json:"createdCount"`
	CreatedNodes []string `json:"createdNodes"`
	SkippedNodes []string `json:"skippedNodes,omitempty"`
}

package response

// StatsData 统计概览数据
type StatsData struct {
	RunningNotebooks   int64 `json:"runningNotebooks"`
	RunningTraining    int64 `json:"runningTraining"`
	RunningInference   int64 `json:"runningInference"`
	TotalNotebooks     int64 `json:"totalNotebooks"`
	TotalTraining      int64 `json:"totalTraining"`
	TotalInference     int64 `json:"totalInference"`
	TotalGpuInUse      int64 `json:"totalGpuInUse"`
	StorageUsed        int64 `json:"storageUsed"` // GB
	StorageVolumeCount int64 `json:"storageVolumeCount"`
}

// RecentInstance 最近实例
type RecentInstance struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Type      string `json:"type"` // notebook / training / inference
	CreatedAt string `json:"createdAt"`
	GPU       int64  `json:"gpu"`
}

// UsageTrend 使用趋势
type UsageTrend struct {
	Date         string  `json:"date"`
	GpuUsage     float64 `json:"gpuUsage"`
	RunningTasks int     `json:"runningTasks"`
	StorageUsage float64 `json:"storageUsage"`
}

// DashboardData 仪表盘聚合数据
type DashboardData struct {
	Stats           StatsData        `json:"stats"`
	RecentInstances []RecentInstance `json:"recentInstances"`
	UsageTrends     []UsageTrend     `json:"usageTrends"`
}

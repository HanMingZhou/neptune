package dashboard

import "gin-vue-admin/global"

// DailyMetric 每日指标记录 (用于历史趋势图)
type DailyMetric struct {
	global.GVA_MODEL
	Date         string  `json:"date" gorm:"index;comment:日期(YYYY-MM-DD)"`
	GpuUsage     float64 `json:"gpuUsage" gorm:"comment:GPU使用率"`
	RunningTasks int     `json:"runningTasks" gorm:"comment:运行中任务数"`
	StorageUsage float64 `json:"storageUsage" gorm:"comment:存储使用量(GB)"`
	TotalCost    float64 `json:"totalCost" gorm:"comment:总花费"`
}

func (DailyMetric) TableName() string {
	return "sys_daily_metrics"
}

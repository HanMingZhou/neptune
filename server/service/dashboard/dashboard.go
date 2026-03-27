package dashboard

import (
	"context"
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/model/consts"
	dashboardModel "gin-vue-admin/model/dashboard"
	dashboardReq "gin-vue-admin/model/dashboard/request"
	dashboardRes "gin-vue-admin/model/dashboard/response"
	inferenceModel "gin-vue-admin/model/inference"
	notebookModel "gin-vue-admin/model/notebook"
	pvcModel "gin-vue-admin/model/pvc"
	trainingModel "gin-vue-admin/model/training"
	"time"
)

type DashboardService struct{}

// GetDashboardData 获取仪表盘聚合数据
func (s *DashboardService) GetDashboardData(ctx context.Context, req *dashboardReq.GetDashboardDataReq) (data dashboardRes.DashboardData, err error) {
	db := global.GVA_DB
	userId := req.UserId

	// ========== 1. 统计概览 (Stats) ==========

	// Notebook 统计
	var runningNotebooks, totalNotebooks int64
	db.Model(&notebookModel.Notebook{}).Where("user_id = ? AND status = ?", userId, consts.NotebookStatusRunning).Count(&runningNotebooks)
	db.Model(&notebookModel.Notebook{}).Where("user_id = ?", userId).Count(&totalNotebooks)

	// Training 统计
	var runningTraining, totalTraining int64
	db.Model(&trainingModel.TrainingJob{}).Where("user_id = ? AND status IN ?", userId, []string{consts.TrainingStatusRunning, consts.TrainingStatusPending}).Count(&runningTraining)
	db.Model(&trainingModel.TrainingJob{}).Where("user_id = ?", userId).Count(&totalTraining)

	// Inference 统计
	var runningInference, totalInference int64
	db.Model(&inferenceModel.Inference{}).Where("user_id = ? AND status = ?", userId, consts.InferenceStatusRunning).Count(&runningInference)
	db.Model(&inferenceModel.Inference{}).Where("user_id = ?", userId).Count(&totalInference)

	// GPU 使用量：统计当前用户所有 Running 状态资源的 GPU 之和
	var notebookGpuSum, trainingGpuSum, inferenceGpuSum int64

	// Notebook GPU
	row := db.Model(&notebookModel.Notebook{}).Where("user_id = ? AND status = ?", userId, consts.NotebookStatusRunning).Select("COALESCE(SUM(gpu), 0)").Row()
	_ = row.Scan(&notebookGpuSum)

	// Training GPU
	row = db.Model(&trainingModel.TrainingJob{}).Where("user_id = ? AND status IN ?", userId, []string{consts.TrainingStatusRunning, consts.TrainingStatusPending}).Select("COALESCE(SUM(total_gpu_count), 0)").Row()
	_ = row.Scan(&trainingGpuSum)

	// Inference GPU
	row = db.Model(&inferenceModel.Inference{}).Where("user_id = ? AND status = ?", userId, consts.InferenceStatusRunning).Select("COALESCE(SUM(gpu), 0)").Row()
	_ = row.Scan(&inferenceGpuSum)

	totalGpuInUse := notebookGpuSum + trainingGpuSum + inferenceGpuSum

	// 存储统计
	var storageUsed int64
	var storageVolumeCount int64
	row = db.Model(&pvcModel.Volume{}).Where("user_id = ?", userId).Select("COALESCE(SUM(size), 0)").Row()
	_ = row.Scan(&storageUsed)
	db.Model(&pvcModel.Volume{}).Where("user_id = ?", userId).Count(&storageVolumeCount)

	data.Stats = dashboardRes.StatsData{
		RunningNotebooks:   runningNotebooks,
		RunningTraining:    runningTraining,
		RunningInference:   runningInference,
		TotalNotebooks:     totalNotebooks,
		TotalTraining:      totalTraining,
		TotalInference:     totalInference,
		TotalGpuInUse:      totalGpuInUse,
		StorageUsed:        storageUsed,
		StorageVolumeCount: storageVolumeCount,
	}

	// ========== 2. 最近实例（最近创建的 Notebook / Training / Inference，合并排序） ==========
	limit := 10
	if req.Days > 0 && req.Days < 10 {
		limit = req.Days
	}

	// 最近 Notebook
	var notebooks []notebookModel.Notebook
	db.Where("user_id = ?", userId).Order("created_at desc").Limit(limit).Find(&notebooks)
	for _, nb := range notebooks {
		data.RecentInstances = append(data.RecentInstances, dashboardRes.RecentInstance{
			ID:        nb.ID,
			Name:      nb.DisplayName,
			Status:    nb.Status,
			Type:      "notebook",
			CreatedAt: nb.CreatedAt.Format(time.DateTime),
			GPU:       nb.GPU,
		})
	}

	// 最近 Training
	var trainings []trainingModel.TrainingJob
	db.Where("user_id = ?", userId).Order("created_at desc").Limit(limit).Find(&trainings)
	for _, tj := range trainings {
		data.RecentInstances = append(data.RecentInstances, dashboardRes.RecentInstance{
			ID:        tj.ID,
			Name:      tj.DisplayName,
			Status:    tj.Status,
			Type:      "training",
			CreatedAt: tj.CreatedAt.Format(time.DateTime),
			GPU:       tj.TotalGPUCount,
		})
	}

	// 最近 Inference
	var inferences []inferenceModel.Inference
	db.Where("user_id = ?", userId).Order("created_at desc").Limit(limit).Find(&inferences)
	for _, inf := range inferences {
		data.RecentInstances = append(data.RecentInstances, dashboardRes.RecentInstance{
			ID:        inf.ID,
			Name:      inf.DisplayName,
			Status:    inf.Status,
			Type:      "inference",
			CreatedAt: inf.CreatedAt.Format(time.DateTime),
			GPU:       inf.GPU,
		})
	}

	// 按时间倒序排列，只取前 limit 条
	sortRecentInstances(data.RecentInstances)
	if len(data.RecentInstances) > limit {
		data.RecentInstances = data.RecentInstances[:limit]
	}

	// ========== 3. 使用趋势（从 DailyMetric 表查询） ==========
	days := 7
	if req.Days > 0 {
		days = req.Days
	}
	var metrics []dashboardModel.DailyMetric
	db.Where("date >= ?", time.Now().AddDate(0, 0, -days).Format("2006-01-02")).Order("date asc").Find(&metrics)

	for _, m := range metrics {
		data.UsageTrends = append(data.UsageTrends, dashboardRes.UsageTrend{
			Date:         m.Date,
			GpuUsage:     m.GpuUsage,
			RunningTasks: m.RunningTasks,
			StorageUsage: m.StorageUsage,
		})
	}

	// 如果没有历史数据，返回当天的实时快照作为唯一数据点
	if len(data.UsageTrends) == 0 {
		data.UsageTrends = append(data.UsageTrends, dashboardRes.UsageTrend{
			Date:         time.Now().Format("2006-01-02"),
			GpuUsage:     float64(totalGpuInUse),
			RunningTasks: int(runningNotebooks + runningTraining + runningInference),
			StorageUsage: float64(storageUsed),
		})
	}

	return data, nil
}

// sortRecentInstances 按 CreatedAt 降序排序
func sortRecentInstances(items []dashboardRes.RecentInstance) {
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].CreatedAt < items[j].CreatedAt {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

// RecordDailyMetric 记录每日指标快照（可由定时任务调用）
func (s *DashboardService) RecordDailyMetric(ctx context.Context) error {
	db := global.GVA_DB
	today := time.Now().Format("2006-01-02")

	// 检查是否已存在今天的记录
	var count int64
	db.Model(&dashboardModel.DailyMetric{}).Where("date = ?", today).Count(&count)
	if count > 0 {
		return fmt.Errorf("today's metric already recorded")
	}

	// 统计当前数据
	var runningNotebooks, runningTraining, runningInference int64
	db.Model(&notebookModel.Notebook{}).Where("status = ?", consts.NotebookStatusRunning).Count(&runningNotebooks)
	db.Model(&trainingModel.TrainingJob{}).Where("status IN ?", []string{consts.TrainingStatusRunning, consts.TrainingStatusPending}).Count(&runningTraining)
	db.Model(&inferenceModel.Inference{}).Where("status = ?", consts.InferenceStatusRunning).Count(&runningInference)

	var gpuSum int64
	row := db.Model(&notebookModel.Notebook{}).Where("status = ?", consts.NotebookStatusRunning).Select("COALESCE(SUM(gpu), 0)").Row()
	_ = row.Scan(&gpuSum)

	var storageSum int64
	row = db.Model(&pvcModel.Volume{}).Select("COALESCE(SUM(size), 0)").Row()
	_ = row.Scan(&storageSum)

	metric := dashboardModel.DailyMetric{
		Date:         today,
		GpuUsage:     float64(gpuSum),
		RunningTasks: int(runningNotebooks + runningTraining + runningInference),
		StorageUsage: float64(storageSum),
		TotalCost:    0, // 需要计费系统接入后补充
	}

	return db.Create(&metric).Error
}

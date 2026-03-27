package training

import (
	v1 "gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type TrainingRouter struct{}

// InitTrainingRouter 初始化训练任务路由
func (r *TrainingRouter) InitTrainingRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	trainingRouter := Router.Group("training").Use(middleware.OperationRecord())
	trainingRouterWithoutRecord := Router.Group("training")
	trainingPublicRouter := PublicRouter.Group("training")

	trainingApi := v1.ApiGroupApp.TrainingApiGroup.TrainingJobApi

	{
		trainingRouter.POST("add", trainingApi.CreateTrainingJob)    // 创建训练任务
		trainingRouter.POST("delete", trainingApi.DeleteTrainingJob) // 删除训练任务
		trainingRouter.POST("stop", trainingApi.StopTrainingJob)     // 停止训练任务
	}

	{
		trainingRouterWithoutRecord.GET("list", trainingApi.GetTrainingJobList)         // 获取训练任务列表
		trainingRouterWithoutRecord.GET("get", trainingApi.GetTrainingJobDetail)        // 获取训练任务详情
		trainingRouterWithoutRecord.GET("pod/list", trainingApi.GetTrainingJobPods)     // 获取训练任务 Pod 列表
		trainingRouterWithoutRecord.GET("log/list", trainingApi.GetTrainingJobLogs)     // 获取训练任务日志
		trainingRouterWithoutRecord.GET("log/download", trainingApi.GetTrainingJobLogs) // 下载训练任务日志
	}

	{
		trainingPublicRouter.GET("terminal", trainingApi.HandleTerminal)     // Web Terminal (WebSocket)
		trainingPublicRouter.GET("log/stream", trainingApi.HandleStreamLogs) // 日志流 (WebSocket)
	}
}

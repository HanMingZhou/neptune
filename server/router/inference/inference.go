package inference

import (
	v1 "gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type InferenceRouter struct{}

// InitInferenceRouter 初始化推理服务路由
func (r *InferenceRouter) InitInferenceRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	inferenceRouter := Router.Group("inference").Use(middleware.OperationRecord())
	inferenceWithNotRecorderRouter := Router.Group("inference")
	inferencePublicRouter := PublicRouter.Group("inference")
	inferenceApi := v1.ApiGroupApp.InferenceApiGroup

	{
		// 推理服务管理 (需要操作记录)
		inferenceRouter.POST("add", inferenceApi.CreateInferenceService)
		inferenceRouter.POST("delete", inferenceApi.DeleteInferenceService)
		inferenceRouter.POST("stop", inferenceApi.StopInferenceService)
		inferenceRouter.POST("start", inferenceApi.StartInferenceService)

		// API Key 管理
		inferenceRouter.POST("api/key/add", inferenceApi.CreateApiKey)
		inferenceRouter.POST("api/key/delete", inferenceApi.DeleteApiKey)
	}

	{
		// WebSocket 路由（无法携带 Header，通过 query token 鉴权）
		inferencePublicRouter.GET("terminal", inferenceApi.HandleTerminal)
		inferencePublicRouter.GET("log/stream", inferenceApi.HandleStreamLogs)
	}

	{
		inferenceWithNotRecorderRouter.POST("list", inferenceApi.GetInferenceServiceList)
		inferenceWithNotRecorderRouter.GET("get", inferenceApi.GetInferenceServiceDetail)
		inferenceWithNotRecorderRouter.GET("pod/list", inferenceApi.GetInferenceServicePods)
		inferenceWithNotRecorderRouter.GET("log/list", inferenceApi.GetInferenceServiceLogs)
		inferenceWithNotRecorderRouter.POST("api/key/list", inferenceApi.ListApiKeys)
	}
}

package notebook

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type NotebookRouter struct{}

func (s *NotebookRouter) InitNotebookApiRouter(PrivateGroup *gin.RouterGroup, PublicGroup *gin.RouterGroup) {
	notebookRouter := PrivateGroup.Group("notebook").Use(middleware.OperationRecord())
	notebookWithNotRecorderRouter := PrivateGroup.Group("notebook")
	notebookPublicRouter := PublicGroup.Group("notebook")
	{
		notebookRouter.POST("add", notebookApi.AddNoteBook)       // 创建容器实例
		notebookRouter.POST("delete", notebookApi.DeleteNoteBook) // 删除容器实例
		notebookRouter.POST("update", notebookApi.UpdateNoteBook) // 更新容器实例
		notebookRouter.POST("stop", notebookApi.StopNotebook)     // 停止容器实例
		notebookRouter.POST("start", notebookApi.StartNotebook)   // 启动容器实例
	}
	{
		// WebSocket 路由（无法携带 Header，通过 query token 鉴权）
		notebookPublicRouter.GET("terminal", notebookApi.HandleTerminal)     // Web Terminal (WebSocket)
		notebookPublicRouter.GET("log/stream", notebookApi.HandleStreamLogs) // 日志流 (WebSocket)
	}
	{
		notebookWithNotRecorderRouter.GET("list", notebookApi.GetNotebookList)     // 获取容器实例
		notebookWithNotRecorderRouter.GET("get", notebookApi.GetNotebookDetail)    // 获取容器实例详情
		notebookWithNotRecorderRouter.GET("pod/list", notebookApi.GetNotebookPods) // 获取容器实例 Pod 列表
		notebookWithNotRecorderRouter.GET("log/list", notebookApi.GetNotebookLogs) // 获取容器实例日志
	}
}

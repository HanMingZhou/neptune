package piper

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type PipeRouter struct{}

func (s *PipeRouter) InitPipeRouter(PrivateGroup *gin.RouterGroup) {
	apiRouter := PrivateGroup.Group("piper").Use(middleware.OperationRecord())
	apiWithNotRecorderRouter := PrivateGroup.Group("piper")
	{
		apiRouter.POST("add", pipeApi.AddPipe)       // 创建Pipe
		apiRouter.POST("delete", pipeApi.DeletePipe) // 删除Pipe
		apiRouter.POST("update", pipeApi.UpdatePipe) // 更新Pipe
	}
	{
		apiWithNotRecorderRouter.GET("list", pipeApi.GetPipeList) // 获取Pipe列表
	}
}

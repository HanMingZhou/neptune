package image

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type RouterGroup struct {
	ImageRouter
}

type ImageRouter struct{}

func (r *ImageRouter) InitImageRouter(Router *gin.RouterGroup) {
	imageRouter := Router.Group("image").Use(middleware.OperationRecord())
	imageWithNotRecorderRouter := Router.Group("image")
	imageApi := imageApiGroup.ImageApi
	{
		imageRouter.POST("add", imageApi.CreateImage)    // 创建镜像
		imageRouter.POST("update", imageApi.UpdateImage) // 更新镜像
		imageRouter.POST("delete", imageApi.DeleteImage) // 删除镜像
	}
	{
		imageWithNotRecorderRouter.GET("list", imageApi.GetImageList) // 获取镜像列表
	}
}

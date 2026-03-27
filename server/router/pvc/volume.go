package pvc

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type VolumeRouter struct{}

// InitVolumeRouter 初始化 Volume 路由
func (r *VolumeRouter) InitVolumeRouter(Router *gin.RouterGroup) {
	volumeRouter := Router.Group("volume").Use(middleware.OperationRecord())
	{
		volumeRouter.POST("add", volumeApi.CreateVolume)     // 创建文件存储
		volumeRouter.GET("list", volumeApi.GetVolumeList)    // 获取文件存储列表
		volumeRouter.POST("expand", volumeApi.ExpandVolume)  // 扩容文件存储
		volumeRouter.POST("delete", volumeApi.DeleteVolume)  // 删除文件存储
		volumeRouter.GET("area/list", volumeApi.GetAreaList) // 获取可用区域列表
	}
}

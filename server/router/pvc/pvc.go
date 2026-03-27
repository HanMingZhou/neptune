package pvc

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type PVCRouter struct{}

func (s *PVCRouter) InitPVCRouter(PrivateGroup *gin.RouterGroup) {
	apiRouter := PrivateGroup.Group("pvc").Use(middleware.OperationRecord())
	apiWithNotRecorderRouter := PrivateGroup.Group("pvc")
	{
		apiRouter.POST("add", pvcApi.AddPVC)       // 创建PVC
		apiRouter.POST("delete", pvcApi.DeletePVC) // 删除PVC
		apiRouter.POST("update", pvcApi.UpdatePVC) // 更新PVC
	}
	{
		apiWithNotRecorderRouter.GET("list", pvcApi.GetPVCList) // 获取PVC列表
	}
}

package system

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type SysParamsRouter struct{}

// InitSysParamsRouter 初始化 参数 路由信息
func (s *SysParamsRouter) InitSysParamsRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	sysParamsRouter := Router.Group("params").Use(middleware.OperationRecord())
	{
		sysParamsRouter.POST("add", sysParamsApi.CreateSysParams) // 新建参数

		sysParamsRouter.POST("delete", sysParamsApi.DeleteSysParams)            // 删除参数
		sysParamsRouter.POST("delete/multi", sysParamsApi.DeleteSysParamsByIds) // 批量删除参数
		sysParamsRouter.POST("update", sysParamsApi.UpdateSysParams)            // 更新参数
		sysParamsRouter.GET("get", sysParamsApi.FindSysParams)                  // 根据ID获取参数
		sysParamsRouter.GET("list", sysParamsApi.GetSysParamsList)              // 获取参数列表
		sysParamsRouter.GET("get/key", sysParamsApi.GetSysParam)                // 根据Key获取参数
	}
}

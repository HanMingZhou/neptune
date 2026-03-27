package system

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type SysErrorRouter struct{}

// InitSysErrorRouter 初始化 错误日志 路由信息
func (s *SysErrorRouter) InitSysErrorRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	sysErrorRouter := Router.Group("error").Use(middleware.OperationRecord())
	sysErrorRouterWithoutAuth := PublicRouter.Group("error")
	{

		sysErrorRouter.POST("delete", sysErrorApi.DeleteSysError)            // 删除错误日志
		sysErrorRouter.POST("delete/multi", sysErrorApi.DeleteSysErrorByIds) // 批量删除错误日志
		sysErrorRouter.POST("update", sysErrorApi.UpdateSysError)            // 更新错误日志
		sysErrorRouter.GET("solution/get", sysErrorApi.GetSysErrorSolution)  // 触发错误日志处理
		sysErrorRouter.GET("get", sysErrorApi.FindSysError)                  // 根据ID获取错误日志
		sysErrorRouter.GET("list", sysErrorApi.GetSysErrorList)              // 获取错误日志列表
	}
	{
		sysErrorRouterWithoutAuth.POST("add", sysErrorApi.CreateSysError) // 新建错误日志
	}
}

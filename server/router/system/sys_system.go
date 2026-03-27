package system

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type SysRouter struct{}

func (s *SysRouter) InitSystemRouter(Router *gin.RouterGroup) {
	sysRouter := Router.Group("system").Use(middleware.OperationRecord())

	{
		sysRouter.POST("config/update", systemApi.SetSystemConfig) // 设置配置文件内容
		sysRouter.POST("reload", systemApi.ReloadSystem)           // 重启服务
		sysRouter.POST("config/get", systemApi.GetSystemConfig)    // 获取配置文件内容
		sysRouter.POST("info", systemApi.GetServerInfo)            // 获取服务器信息
	}
}

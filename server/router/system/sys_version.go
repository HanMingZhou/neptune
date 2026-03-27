package system

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type SysVersionRouter struct{}

// InitSysVersionRouter 初始化 版本管理 路由信息
func (s *SysVersionRouter) InitSysVersionRouter(Router *gin.RouterGroup) {
	sysVersionRouter := Router.Group("version").Use(middleware.OperationRecord())
	{
		sysVersionRouter.GET("get", sysVersionApi.FindSysVersion)                      // 获取单一版本
		sysVersionRouter.GET("list", sysVersionApi.GetSysVersionList)                  // 获取版本列表
		sysVersionRouter.GET("download/json", sysVersionApi.DownloadVersionJson)       // 下载版本json
		sysVersionRouter.POST("export", sysVersionApi.ExportVersion)                   // 创建版本
		sysVersionRouter.POST("import", sysVersionApi.ImportVersion)                   // 同步版本
		sysVersionRouter.POST("delete", sysVersionApi.DeleteSysVersion)                // 删除版本
		sysVersionRouter.POST("delete/multi", sysVersionApi.DeleteSysVersionByIds)     // 批量删除版本
		sysVersionRouter.GET("downloadVersionJson", sysVersionApi.DownloadVersionJson) // 下载版本JSON数据
	}
}

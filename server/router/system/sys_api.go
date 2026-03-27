package system

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type ApiRouter struct{}

func (s *ApiRouter) InitApiRouter(Router *gin.RouterGroup, RouterPub *gin.RouterGroup) {
	apiRouter := Router.Group("api").Use(middleware.OperationRecord())

	apiPublicRouterWithoutRecord := RouterPub.Group("api")
	{
		apiRouter.GET("group/list", apiRouterApi.GetApiGroups)       // 获取路由组
		apiRouter.GET("sync", apiRouterApi.SyncApi)                  // 同步Api
		apiRouter.POST("ignore", apiRouterApi.IgnoreApi)             // 忽略Api
		apiRouter.POST("sync/enter", apiRouterApi.EnterSyncApi)      // 确认同步Api
		apiRouter.POST("add", apiRouterApi.CreateApi)                // 创建Api
		apiRouter.POST("delete", apiRouterApi.DeleteApi)             // 删除Api
		apiRouter.POST("get", apiRouterApi.GetApiById)               // 获取单条Api消息
		apiRouter.POST("update", apiRouterApi.UpdateApi)             // 更新api
		apiRouter.POST("delete/multi", apiRouterApi.DeleteApisByIds) // 删除选中api
		apiRouter.POST("all", apiRouterApi.GetAllApis)               // 获取所有api
		apiRouter.POST("list", apiRouterApi.GetApiList)              // 获取Api列表
	}
	{
		apiPublicRouterWithoutRecord.GET("casbin/fresh", apiRouterApi.FreshCasbin) // 刷新casbin权限
	}
}

package system

import (
	"github.com/gin-gonic/gin"
)

type AutoCodeHistoryRouter struct{}

func (s *AutoCodeRouter) InitAutoCodeHistoryRouter(Router *gin.RouterGroup) {
	autoCodeHistoryRouter := Router.Group("autocode")
	{
		autoCodeHistoryRouter.POST("rollback", autocodeHistoryApi.RollBack) // 回滚
		autoCodeHistoryRouter.POST("meta/get", autocodeHistoryApi.First)    // 获取meta
		autoCodeHistoryRouter.POST("list", autocodeHistoryApi.GetList)      // 查询历史记录
		autoCodeHistoryRouter.POST("delete", autocodeHistoryApi.Delete)     // 删除历史记录
	}
}

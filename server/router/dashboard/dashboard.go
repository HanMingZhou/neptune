package dashboard

import (
	v1 "gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type DashboardRouter struct{}

func (s *DashboardRouter) InitDashboardRouter(Router *gin.RouterGroup) {
	dashboardRouter := Router.Group("dashboard").Use(middleware.OperationRecord())
	dashboardApi := v1.ApiGroupApp.DashboardApiGroup.DashboardApi
	{
		dashboardRouter.POST("get", dashboardApi.GetDashboardData) // 获取仪表盘数据
	}
}

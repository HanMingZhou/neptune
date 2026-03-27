package dashboard

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model/common/response"
	dashboardReq "gin-vue-admin/model/dashboard/request"
	"gin-vue-admin/service"
	"gin-vue-admin/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DashboardApi struct{}

// GetDashboardData
// @Tags      Dashboard
// @Summary   获取仪表盘数据
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      dashboardReq.GetDashboardDataReq  true  "查询参数"
// @Success   200   {object}  response.Response{data=dashboardRes.DashboardData,msg=string}  "获取仪表盘数据"
// @Router    /dashboard/getData [post]
func (api *DashboardApi) GetDashboardData(c *gin.Context) {
	var req dashboardReq.GetDashboardDataReq
	_ = c.ShouldBindJSON(&req)
	req.UserId = utils.GetUserID(c)

	dashboardService := service.ServiceGroupApp.DashboardServiceGroup.DashboardService
	if data, err := dashboardService.GetDashboardData(c.Request.Context(), &req); err != nil {
		global.GVA_LOG.Error("获取仪表盘数据失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(data, "获取成功", c)
	}
}

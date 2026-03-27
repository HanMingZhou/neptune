package order

import (
	v1 "gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type OrderRouter struct{}

func (s *OrderRouter) InitOrderRouter(Router *gin.RouterGroup) {
	orderRouter := Router.Group("order").Use(middleware.OperationRecord())
	orderWithNotRecorderRouter := Router.Group("order")
	orderApi := v1.ApiGroupApp.OrderApiGroup.OrderApi
	{
		orderRouter.POST("invoice/apply", orderApi.ApplyInvoice) // 申请发票
	}
	{
		orderWithNotRecorderRouter.GET("overview", orderApi.GetOrderOverview)            // 获取财务总览
		orderWithNotRecorderRouter.POST("usage/list", orderApi.GetUsageList)             // 使用详情列表
		orderWithNotRecorderRouter.POST("transaction/list", orderApi.GetTransactionList) // 收支明细列表
		orderWithNotRecorderRouter.POST("order/list", orderApi.GetOrderList)             // 订单列表
		orderWithNotRecorderRouter.POST("invoice/list", orderApi.GetInvoiceList)         // 发票记录列表
	}
}

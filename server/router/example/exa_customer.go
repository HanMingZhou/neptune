package example

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type CustomerRouter struct{}

func (e *CustomerRouter) InitCustomerRouter(Router *gin.RouterGroup) {
	customerRouter := Router.Group("customer").Use(middleware.OperationRecord())
	customerRouterWithoutRecord := Router.Group("customer")
	{
		customerRouter.POST("add", exaCustomerApi.CreateExaCustomer) // 创建客户

		customerRouter.POST("update", exaCustomerApi.UpdateExaCustomer) // 更新客户
		customerRouter.POST("delete", exaCustomerApi.DeleteExaCustomer) // 删除客户
	}
	{
		customerRouterWithoutRecord.GET("get", exaCustomerApi.GetExaCustomer)      // 获取单一客户信息
		customerRouterWithoutRecord.GET("list", exaCustomerApi.GetExaCustomerList) // 获取客户列表
	}
}

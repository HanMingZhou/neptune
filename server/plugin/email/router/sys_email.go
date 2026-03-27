package router

import (
	"gin-vue-admin/middleware"
	"gin-vue-admin/plugin/email/api"

	"github.com/gin-gonic/gin"
)

type EmailRouter struct{}

func (s *EmailRouter) InitEmailRouter(Router *gin.RouterGroup) {
	emailRouter := Router.Use(middleware.OperationRecord())
	EmailApi := api.ApiGroupApp.EmailApi.EmailTest
	SendEmail := api.ApiGroupApp.EmailApi.SendEmail
	{
		emailRouter.POST("test", EmailApi)  // 发送测试邮件
		emailRouter.POST("send", SendEmail) // 发送邮件
	}
}

package apisix

import (
	"github.com/gin-gonic/gin"
)

type ApisixRouter struct{}

// InitApisixRouter 初始化 Apisix 路由信息
func (s *ApisixRouter) InitApisixRouter(PublicGroup *gin.RouterGroup) {
	apisixRouter := PublicGroup.Group("apisix")
	{
		apisixRouter.GET("auth", apisixApi.AuthApisix)
		apisixRouter.POST("auth", apisixApi.AuthApisix)
	}
}

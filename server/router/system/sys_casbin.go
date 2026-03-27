package system

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type CasbinRouter struct{}

func (s *CasbinRouter) InitCasbinRouter(Router *gin.RouterGroup) {
	casbinRouter := Router.Group("casbin").Use(middleware.OperationRecord())
	{
		casbinRouter.POST("update", casbinApi.UpdateCasbin)
		casbinRouter.POST("policy/list", casbinApi.GetPolicyPathByAuthorityId)
	}
}

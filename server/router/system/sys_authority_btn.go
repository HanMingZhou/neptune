package system

import (
	"github.com/gin-gonic/gin"
)

type AuthorityBtnRouter struct{}

var AuthorityBtnRouterApp = new(AuthorityBtnRouter)

func (s *AuthorityBtnRouter) InitAuthorityBtnRouterRouter(Router *gin.RouterGroup) {
	authorityRouter := Router.Group("authority/btn")
	{
		authorityRouter.POST("update", authorityBtnApi.SetAuthorityBtn)       // 设置按钮权限
		authorityRouter.POST("get", authorityBtnApi.GetAuthorityBtn)          // 获取已有按钮权限
		authorityRouter.POST("delete", authorityBtnApi.CanRemoveAuthorityBtn) // 删除按钮
	}
}

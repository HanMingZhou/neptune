package system

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type AuthorityRouter struct{}

func (s *AuthorityRouter) InitAuthorityRouter(Router *gin.RouterGroup) {
	authorityRouter := Router.Group("authority").Use(middleware.OperationRecord())
	{
		authorityRouter.POST("add", authorityApi.CreateAuthority)    // 创建角色
		authorityRouter.POST("delete", authorityApi.DeleteAuthority) // 删除角色

		authorityRouter.POST("update", authorityApi.UpdateAuthority)                 // 更新角色
		authorityRouter.POST("copy", authorityApi.CopyAuthority)                     // 拷贝角色
		authorityRouter.POST("data/authority/update", authorityApi.SetDataAuthority) // 设置角色资源权限
		authorityRouter.POST("list", authorityApi.GetAuthorityList)                  // 获取角色列表
	}
}

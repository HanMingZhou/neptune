package system

import (
	v1 "gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (s *UserRouter) InitUserRouter(Router *gin.RouterGroup, v1PublicGroup *gin.RouterGroup) {
	userRouter := Router.Group("user").Use(middleware.OperationRecord())
	publicRouter := v1PublicGroup.Group("user")
	userApi := v1.ApiGroupApp.SystemApiGroup.BaseApi

	{

		userRouter.POST("password/update", userApi.ChangePassword)    // 用户修改密码
		userRouter.POST("authority/update", userApi.SetUserAuthority) // 设置用户权限

		userRouter.POST("delete", userApi.DeleteUser)                     // 删除用户
		userRouter.POST("info/update", userApi.SetUserInfo)               // 设置用户信息
		userRouter.POST("self/info/update", userApi.SetSelfInfo)          // 设置自身信息
		userRouter.POST("authorities/update", userApi.SetUserAuthorities) // 设置用户权限组
		userRouter.POST("password/reset", userApi.ResetPassword)          // 重置用户密码
		userRouter.POST("self/setting/update", userApi.SetSelfSetting)    // 用户界面配置
		userRouter.POST("list", userApi.GetUserList)                      // 分页获取用户列表
		userRouter.GET("info", userApi.GetUserInfo)                       // 获取自身信息
	}
	{
		publicRouter.POST("register", userApi.Register) // 注册账号
	}
}

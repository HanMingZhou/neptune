package account

import (
	v1 "gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type AccountRouter struct{}

func (r *AccountRouter) InitAccountRouter(Router *gin.RouterGroup) {
	accountRouter := Router.Group("account").Use(middleware.OperationRecord())
	accountWithNotRecorderRouter := Router.Group("account")
	accountApi := v1.ApiGroupApp.AccountApiGroup.AccountApi

	{
		accountRouter.GET("security/status", accountApi.GetSecurityStatus) // 获取安全状态
		accountRouter.POST("password/update", accountApi.UpdatePassword)   // 修改密码
		accountRouter.POST("bind", accountApi.BindAccount)                 // 绑定手机/邮箱
		accountRouter.POST("mfa/setup", accountApi.MfaSetup)               // MFA初始化（生成密钥和二维码）
		accountRouter.POST("mfa/activate", accountApi.MfaActivate)         // MFA激活（验证码确认绑定）
		accountRouter.POST("mfa/toggle", accountApi.ToggleMfa)             // 关闭MFA
		accountRouter.POST("ak/generate", accountApi.GenerateAccessKey)    // 生成AccessKey
		accountRouter.POST("active/session/kill", accountApi.KillSession)  // 强制下线
	}
	{
		accountWithNotRecorderRouter.POST("access/log/list", accountApi.GetAccessLogList)         // 获取访问日志列表
		accountWithNotRecorderRouter.POST("active/session/list", accountApi.GetActiveSessionList) // 获取活跃会话列表
	}
}

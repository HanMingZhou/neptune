package secret

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type SecretRouter struct{}

func (s *SecretRouter) InitSecretRouter(PrivateGroup *gin.RouterGroup) {
	apiRouter := PrivateGroup.Group("secret").Use(middleware.OperationRecord())
	apiWithNotRecorderRouter := PrivateGroup.Group("secret")
	{
		apiRouter.POST("add", secretApi.AddSecret)       // 创建Secret
		apiRouter.POST("delete", secretApi.DeleteSecret) // 删除Secret
		apiRouter.POST("update", secretApi.UpdateSecret) // 更新Secret
	}
	{
		apiWithNotRecorderRouter.GET("list", secretApi.GetSecretList) // 获取Secret列表
	}
}

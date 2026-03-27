package sshkey

import (
	v1 "gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type SSHKeyRouter struct{}

func (s *SSHKeyRouter) InitSSHKeyRouter(Router *gin.RouterGroup) {
	sshKeyRouter := Router.Group("sshkey").Use(middleware.OperationRecord())
	sshKeyWithNotRecorderRouter := Router.Group("sshkey")
	sshKeyApi := v1.ApiGroupApp.SSHKeyApiGroup.SSHKeyApi

	{
		sshKeyRouter.POST("add", sshKeyApi.CreateSSHKey)                // 创建SSH密钥
		sshKeyRouter.POST("delete", sshKeyApi.DeleteSSHKey)             // 删除SSH密钥
		sshKeyRouter.POST("default/update", sshKeyApi.SetDefaultSSHKey) // 设置默认密钥
	}
	{
		sshKeyWithNotRecorderRouter.POST("list", sshKeyApi.GetSSHKeyList) // 获取SSH密钥列表
	}
}

package sshkey

import (
	"gin-vue-admin/model/common/response"
	secretReq "gin-vue-admin/model/secret/request"
	sshkeyService "gin-vue-admin/service/sshkey"
	"gin-vue-admin/utils"

	"github.com/gin-gonic/gin"
)

type SSHKeyApi struct{}

var sshKeyServiceInstance = new(sshkeyService.SSHKeyService)

// CreateSSHKey 创建SSH密钥
func (s *SSHKeyApi) CreateSSHKey(c *gin.Context) {
	var req secretReq.AddSSHKeyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)

	resp, err := sshKeyServiceInstance.CreateSSHKey(c.Request.Context(), &req, userId)
	if err != nil {
		utils.HandleError(c, err, "创建SSH密钥失败")
		return
	}

	response.OkWithData(resp, c)
}

// GetSSHKeyList 获取SSH密钥列表
func (s *SSHKeyApi) GetSSHKeyList(c *gin.Context) {
	var req secretReq.GetSSHKeyListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	resp, err := sshKeyServiceInstance.GetSSHKeyList(c.Request.Context(), &req, userId)
	if err != nil {
		utils.HandleError(c, err, "获取SSH密钥列表失败")
		return
	}

	response.OkWithData(resp, c)
}

// DeleteSSHKey 删除SSH密钥
func (s *SSHKeyApi) DeleteSSHKey(c *gin.Context) {
	var req secretReq.DeleteSSHKeyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)

	if err := sshKeyServiceInstance.DeleteSSHKey(c.Request.Context(), &req, userId); err != nil {
		utils.HandleError(c, err, "删除SSH密钥失败")
		return
	}

	response.OkWithMessage("删除成功", c)
}

// SetDefaultSSHKey 设置默认密钥
func (s *SSHKeyApi) SetDefaultSSHKey(c *gin.Context) {
	var req secretReq.SetDefaultSSHKeyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)

	if err := sshKeyServiceInstance.SetDefaultSSHKey(c.Request.Context(), &req, userId); err != nil {
		utils.HandleError(c, err, "设置默认密钥失败")
		return
	}

	response.OkWithMessage("设置成功", c)
}

package secret

import (
	"gin-vue-admin/model/common/response"
	"gin-vue-admin/model/secret/request"
	"gin-vue-admin/utils"

	"github.com/gin-gonic/gin"
)

type SecretApi struct{}

// AddSecret 创建Secret
func (a *SecretApi) AddSecret(c *gin.Context) {
	var req request.AddSecretReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
		return
	}

	if err := secretService.CreateSSHSecret(c.Request.Context(), &req); err != nil {
		utils.HandleError(c, err, "创建Secret失败")
		return
	}
	response.OkWithMessage("创建Secret成功", c)
}

// DeleteSecret 删除Secret
func (a *SecretApi) DeleteSecret(c *gin.Context) {
	var req request.DeleteSecretReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if req.InstanceName == "" {
		response.FailWithMessage("Secret名称不能为空", c)
		return
	}
	if err := secretService.DeleteSSHSecret(c.Request.Context(), &req); err != nil {
		utils.HandleError(c, err, "删除Secret失败")
		return
	}
	response.OkWithMessage("删除Secret成功", c)
}

// GetSecretList 获取Secret列表
func (a *SecretApi) GetSecretList(c *gin.Context) {
	// TODO: 实现获取Secret列表逻辑
	response.OkWithMessage("获取Secret列表成功", c)
}

// UpdateSecret 更新Secret
func (a *SecretApi) UpdateSecret(c *gin.Context) {
	// TODO: 实现更新Secret逻辑
	response.OkWithMessage("更新Secret成功", c)
}

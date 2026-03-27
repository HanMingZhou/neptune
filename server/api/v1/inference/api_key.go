package inference

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model/common/response"
	inferenceReq "gin-vue-admin/model/inference/request"
	"gin-vue-admin/service"
	"gin-vue-admin/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ApiKeyApi struct{}

var apiKeyService = service.ServiceGroupApp.InferenceServiceGroup.InferenceApiKeyService

// CreateApiKey 创建 API Key
// @Summary 创建 API Key
// @Router /inference/keys [post]
func (a *ApiKeyApi) CreateApiKey(c *gin.Context) {
	var req inferenceReq.CreateApiKeyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	req.UserId = utils.GetUserID(c)

	resp, err := apiKeyService.CreateApiKey(c.Request.Context(), &req)
	if err != nil {
		global.GVA_LOG.Error("创建API Key失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(resp, c)
}

// ListApiKeys 获取 API Key 列表
// @Summary 获取 API Key 列表
// @Router /inference/keys/list [post]
func (a *ApiKeyApi) ListApiKeys(c *gin.Context) {
	var req inferenceReq.ListApiKeysReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	resp, err := apiKeyService.ListApiKeys(c.Request.Context(), &req)
	if err != nil {
		global.GVA_LOG.Error("获取API Key列表失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(resp, c)
}

// DeleteApiKey 删除 API Key
// @Summary 删除 API Key
// @Router /inference/keys [delete]
func (a *ApiKeyApi) DeleteApiKey(c *gin.Context) {
	var req inferenceReq.DeleteApiKeyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := apiKeyService.DeleteApiKey(c.Request.Context(), &req); err != nil {
		global.GVA_LOG.Error("删除API Key失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

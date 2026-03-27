package account

import (
	"gin-vue-admin/global"
	accountReq "gin-vue-admin/model/account/request"
	"gin-vue-admin/model/common/response"
	"gin-vue-admin/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AccountApi struct{}

// GetAccessLogList
// @Tags      Account
// @Summary   获取访问日志列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Producer  application/json
// @Param     data  body      accountReq.GetAccessLogListReq  true  "查询条件"
// @Success   200   {object}  response.Response{data=accountResp.GetAccessLogListResp}
// @Router    /account/access/log/list [post]
func (a *AccountApi) GetAccessLogList(c *gin.Context) {
	var req accountReq.GetAccessLogListReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	list, err := accountService.GetAccessLogList(c, req, userId)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(list, c)
}

// GetSecurityStatus
// @Tags      Account
// @Summary   获取账号安全状态
// @Security  ApiKeyAuth
// @accept    application/json
// @Producer  application/json
// @Router    /account/security/status [get]
func (a *AccountApi) GetSecurityStatus(c *gin.Context) {
	userId := utils.GetUserID(c)
	status, err := accountService.GetSecurityStatus(c, userId)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(status, c)
}

// GetActiveSessionList
// @Tags      Account
// @Summary   获取活跃会话列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Producer  application/json
// @Param     data  body      accountReq.GetAccessLogListReq  true  "查询条件"
// @Success   200   {object}  response.Response{data=accountResp.GetAccessLogListResp}
// @Router    /account/active/session/list [post]
func (a *AccountApi) GetActiveSessionList(c *gin.Context) {
	var req accountReq.GetAccessLogListReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	list, err := accountService.GetActiveSessionList(c, req, userId)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(list, c)
}

// UpdatePassword
// @Tags      Account
// @Summary   修改密码
// @Security  ApiKeyAuth
// @accept    application/json
// @Producer  application/json
// @Param     data  body      accountReq.UpdatePasswordReq  true  "旧密码, 新密码"
// @Success   200   {object}  response.Response{}
// @Router    /account/password/update [post]
func (a *AccountApi) UpdatePassword(c *gin.Context) {
	var req accountReq.UpdatePasswordReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	err = accountService.UpdatePassword(userId, req)
	if err != nil {
		global.GVA_LOG.Error("修改失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("修改成功", c)
}

// BindAccount
// @Tags      Account
// @Summary   绑定/修改手机号或邮箱
// @Security  ApiKeyAuth
// @accept    application/json
// @Producer  application/json
// @Param     data  body      accountReq.BindAccountReq  true  "类型, 值, 验证码"
// @Success   200   {object}  response.Response{}
// @Router    /account/bind [post]
func (a *AccountApi) BindAccount(c *gin.Context) {
	var req accountReq.BindAccountReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	err = accountService.BindAccount(userId, req)
	if err != nil {
		global.GVA_LOG.Error("绑定失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("绑定成功", c)
}

// MfaSetup
// @Tags      Account
// @Summary   初始化MFA（生成密钥和二维码）
// @Security  ApiKeyAuth
// @accept    application/json
// @Producer  application/json
// @Success   200   {object}  response.Response{data=accountResp.MfaSecretResp}
// @Router    /account/mfa/setup [post]
func (a *AccountApi) MfaSetup(c *gin.Context) {
	userId := utils.GetUserID(c)
	resp, err := accountService.MfaSetup(userId)
	if err != nil {
		global.GVA_LOG.Error("MFA初始化失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(resp, c)
}

// MfaActivate
// @Tags      Account
// @Summary   激活MFA（扫码后输入验证码确认绑定）
// @Security  ApiKeyAuth
// @accept    application/json
// @Producer  application/json
// @Param     data  body      accountReq.MfaActivateReq  true  "验证码"
// @Success   200   {object}  response.Response{}
// @Router    /account/mfa/activate [post]
func (a *AccountApi) MfaActivate(c *gin.Context) {
	var req accountReq.MfaActivateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	if err := accountService.MfaActivate(userId, req.Code); err != nil {
		global.GVA_LOG.Error("MFA激活失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("MFA已成功开启", c)
}

// ToggleMfa
// @Tags      Account
// @Summary   关闭MFA
// @Security  ApiKeyAuth
// @accept    application/json
// @Producer  application/json
// @Param     data  body      accountReq.ToggleMfaReq  true  "状态, 验证码"
// @Success   200   {object}  response.Response{}
// @Router    /account/mfa/toggle [post]
func (a *AccountApi) ToggleMfa(c *gin.Context) {
	var req accountReq.ToggleMfaReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	err = accountService.ToggleMfa(userId, req)
	if err != nil {
		global.GVA_LOG.Error("操作失败!", zap.Error(err))
		response.FailWithMessage("操作失败", c)
		return
	}
	response.OkWithMessage("操作成功", c)
}

// GenerateAccessKey
// @Tags      Account
// @Summary   生成AccessKey
// @Security  ApiKeyAuth
// @accept    application/json
// @Producer  application/json
// @Success   200   {object}  response.Response{data=accountResp.GenerateAccessKeyResp}
// @Router    /account/ak/generate [post]
func (a *AccountApi) GenerateAccessKey(c *gin.Context) {
	userId := utils.GetUserID(c)
	resp, err := accountService.GenerateAccessKey(userId)
	if err != nil {
		global.GVA_LOG.Error("生成失败!", zap.Error(err))
		response.FailWithMessage("生成失败", c)
		return
	}
	response.OkWithData(resp, c)
}

// KillSession
// @Tags      Account
// @Summary   强制退出会话
// @Security  ApiKeyAuth
// @accept    application/json
// @Producer  application/json
// @Param     data  body      accountReq.KillSessionReq  true  "logId"
// @Success   200   {object}  response.Response{}
// @Router    /account/active/session/kill [post]
func (a *AccountApi) KillSession(c *gin.Context) {
	var req accountReq.KillSessionReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	err = accountService.KillSession(c, req, userId)
	if err != nil {
		global.GVA_LOG.Error("退出失败!", zap.Error(err))
		response.FailWithMessage("退出失败", c)
		return
	}
	response.OkWithMessage("成功下线该设备", c)
}

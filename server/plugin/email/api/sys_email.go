package api

import (
	"gin-vue-admin/model/common/response"
	email_response "gin-vue-admin/plugin/email/model/response"
	"gin-vue-admin/plugin/email/service"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

type EmailApi struct{}

// EmailTest
// @Tags      System
// @Summary   发送测试邮件
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200  {string}  string  "{"success":true,"data":{},"msg":"发送成功"}"
// @Router    /email/emailTest [post]
func (s *EmailApi) EmailTest(c *gin.Context) {
	err := service.ServiceGroupApp.EmailTest()
	if err != nil {
		logx.Error("发送失败!")
		response.FailWithMessage("发送失败", c)
		return
	}
	response.OkWithMessage("发送成功", c)
}

// SendEmail
// @Tags      System
// @Summary   发送邮件
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      email_response.Email  true  "发送邮件必须的参数"
// @Success   200   {string}  string                "{"success":true,"data":{},"msg":"发送成功"}"
// @Router    /email/sendEmail [post]
func (s *EmailApi) SendEmail(c *gin.Context) {
	var email email_response.Email
	err := c.ShouldBindJSON(&email)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = service.ServiceGroupApp.SendEmail(email.To, email.Subject, email.Body)
	if err != nil {
		logx.Error("发送失败!")
		response.FailWithMessage("发送失败", c)
		return
	}
	response.OkWithMessage("发送成功", c)
}

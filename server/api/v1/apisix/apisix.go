package apisix

import (
	"fmt"
	"gin-vue-admin/model/apisix/request"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

type ApisixApi struct{}

// AuthApisix Apisix forward-auth 认证接口（Notebook / TensorBoard / Inference 统一入口）
func (a *ApisixApi) AuthApisix(c *gin.Context) {
	// 1. 获取原始请求路径
	originalUri := c.GetHeader("X-Forwarded-Uri")
	if originalUri == "" {
		originalUri = c.GetHeader("X-Original-URI")
	}

	// 2. 提取 API-Key（推理服务可用）
	apiKey := c.GetHeader("API-Key")

	// 3. 提取 Token（优先 Cookie，其次 Header，最后 URL 参数）
	token, _ := c.Cookie("x-token")
	if token == "" {
		token = c.GetHeader("x-token")
	}
	if token == "" {
		token = c.GetHeader("Authorization")
	}
	if token == "" {
		if idx := strings.Index(originalUri, "token="); idx != -1 {
			tokenStart := idx + 6
			tokenEnd := strings.IndexAny(originalUri[tokenStart:], "& ")
			if tokenEnd == -1 {
				token = originalUri[tokenStart:]
			} else {
				token = originalUri[tokenStart : tokenStart+tokenEnd]
			}
		}
	}

	// 必须有 Token 或 API-Key 之一
	if token == "" && apiKey == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// 4. 调用 service 层完成认证
	authReq := &request.AuthApisixReq{
		OriginalUri: originalUri,
		Token:       token,
		ApiKey:      apiKey,
	}
	authResp, err := apisixService.AuthApiSix(c.Request.Context(), authReq)
	if err != nil {
		logx.Error("Apisix 资源认证失败", err)
		errMsg := err.Error()
		if errMsg == "Token 无效" || errMsg == "缺少认证凭据" {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
		return
	}

	// 5. 认证通过，设置响应 Header
	c.Header("X-User-Id", fmt.Sprintf("%d", authResp.UserID))
	c.Header("X-User-Namespace", authResp.Namespace)
	if authResp.Token != "" {
		c.Header("X-Set-Token", authResp.Token)
	}
	if authResp.RateLimit > 0 {
		c.Header("X-Rate-Limit", fmt.Sprintf("%d", authResp.RateLimit))
	}
	c.Status(http.StatusOK)
}

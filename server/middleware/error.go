package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logx.Error(c.Request.URL.Path,
						logx.Field("error", err),
						logx.Field("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logx.Error("[Recovery from panic]",
						logx.Field("error", err),
						logx.Field("request", string(httpRequest)),
						logx.Field("stack", string(debug.Stack())),
					)
				} else {
					logx.Error("[Recovery from panic]",
						logx.Field("error", err),
						logx.Field("request", string(httpRequest)),
					)
				}
				// 返回500错误响应
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code": 500,
					"data": nil,
					"msg":  "服务器内部错误",
				})
			}
		}()
		c.Next()
	}
}

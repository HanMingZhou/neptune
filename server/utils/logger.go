package utils

import (
	"fmt"
	"gin-vue-admin/model/common/response"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

// HandleError 统一处理错误：记录堆栈日志 + 返回错误响应
func HandleError(c *gin.Context, err error, msg string) {
	logx.Error(err)
	response.FailWithMessage(fmt.Sprintf("%s: %v", msg, err), c)
}

// ErrorWithStack 记录错误日志，并将 error 格式化为 %+v 以显示堆栈信息
// 适用于 pkg/errors 创建的带堆栈的错误
func ErrorWithStack(msg string, err error) {
	logx.Error(err)
}

package middleware

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model/common/response"
	"gin-vue-admin/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CasbinHandler 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		waitUse, _ := utils.GetClaims(c)
		//获取请求的PATH
		path := c.Request.URL.Path
		prefix := global.GVA_CONFIG.System.RouterPrefix
		if prefix != "" && !strings.HasPrefix(prefix, "/") {
			prefix = "/" + prefix
		}
		obj := strings.TrimPrefix(path, prefix)
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := strconv.Itoa(int(waitUse.AuthorityId))
		e := utils.GetCasbin() // 判断策略中是否存在
		success, _ := e.Enforce(sub, obj, act)
		if !success {
			response.FailWithDetailed(gin.H{}, "权限不足", c)
			c.Abort()
			return
		}
		c.Next()
	}
}

package router

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

var Info = new(info)

type info struct{}

// Init 初始化 公告 路由信息
func (r *info) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	{
		group := private.Group("info").Use(middleware.OperationRecord())
		group.POST("add", apiInfo.CreateInfo)                 // 新建公告
		group.DELETE("delete", apiInfo.DeleteInfo)            // 删除公告
		group.DELETE("delete/multi", apiInfo.DeleteInfoByIds) // 批量删除公告
		group.PUT("update", apiInfo.UpdateInfo)               // 更新公告
	}
	{
		group := private.Group("info")
		group.GET("get", apiInfo.FindInfo)     // 根据ID获取公告
		group.GET("list", apiInfo.GetInfoList) // 获取公告列表
	}
	{
		group := public.Group("info")
		group.GET("get/datasource", apiInfo.GetInfoDataSource) // 获取公告数据源
		group.GET("get/public", apiInfo.GetInfoPublic)         // 获取公告列表
	}
}

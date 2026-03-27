package system

import (
	"github.com/gin-gonic/gin"
)

type AutoCodeRouter struct{}

func (s *AutoCodeRouter) InitAutoCodeRouter(Router *gin.RouterGroup, RouterPublic *gin.RouterGroup) {
	autoCodeRouter := Router.Group("autocode")
	publicAutoCodeRouter := RouterPublic.Group("autocode")
	{
		autoCodeRouter.GET("db/list", autoCodeApi.GetDB)         // 获取数据库
		autoCodeRouter.GET("table/list", autoCodeApi.GetTables)  // 获取对应数据库的表
		autoCodeRouter.GET("column/list", autoCodeApi.GetColumn) // 获取指定表所有字段信息
	}
	{
		autoCodeRouter.POST("preview", autoCodeTemplateApi.Preview)  // 获取自动创建代码预览
		autoCodeRouter.POST("add", autoCodeTemplateApi.Create)       // 创建自动化代码
		autoCodeRouter.POST("func/add", autoCodeTemplateApi.AddFunc) // 为代码插入方法
	}
	{
		autoCodeRouter.POST("mcp/add", autoCodeTemplateApi.MCP)      // 自动创建Mcp Tool模板
		autoCodeRouter.POST("mcp/list", autoCodeTemplateApi.MCPList) // 获取MCP ToolList
		autoCodeRouter.POST("mcp/test", autoCodeTemplateApi.MCPTest) // MCP 工具测试
	}
	{
		autoCodeRouter.POST("package/list", autoCodePackageApi.All)      // 获取package包
		autoCodeRouter.POST("package/delete", autoCodePackageApi.Delete) // 删除package包
		autoCodeRouter.POST("package/add", autoCodePackageApi.Create)    // 创建package包
	}
	{
		autoCodeRouter.GET("template/list", autoCodePackageApi.Templates) // 创建package包
	}
	{
		autoCodeRouter.POST("plugin/pack", autoCodePluginApi.Packaged)   // 打包插件
		autoCodeRouter.POST("plugin/install", autoCodePluginApi.Install) // 自动安装插件

	}
	{
		publicAutoCodeRouter.POST("llm/add", autoCodeApi.LLMAuto)
		publicAutoCodeRouter.POST("menu/init", autoCodePluginApi.InitMenu) // 同步插件菜单
		publicAutoCodeRouter.POST("api/init", autoCodePluginApi.InitAPI)   // 同步插件API
	}
}

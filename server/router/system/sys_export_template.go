package system

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type SysExportTemplateRouter struct {
}

// InitSysExportTemplateRouter 初始化 导出模板 路由信息
func (s *SysExportTemplateRouter) InitSysExportTemplateRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	sysExportTemplateRouter := Router.Group("export/template").Use(middleware.OperationRecord())
	{
		sysExportTemplateRouter.POST("add", exportTemplateApi.CreateSysExportTemplate) // 新增导出模板

		sysExportTemplateRouter.POST("delete", exportTemplateApi.DeleteSysExportTemplate)            // 删除导出模板
		sysExportTemplateRouter.POST("delete/multi", exportTemplateApi.DeleteSysExportTemplateByIds) // 批量删除导出模板
		sysExportTemplateRouter.POST("update", exportTemplateApi.UpdateSysExportTemplate)            // 更新导出模板
		sysExportTemplateRouter.POST("import", exportTemplateApi.ImportExcel)                        // 导入Excel
		sysExportTemplateRouter.GET("get", exportTemplateApi.FindSysExportTemplate)                  // 根据ID获取导出模板
		sysExportTemplateRouter.GET("list", exportTemplateApi.GetSysExportTemplateList)              // 获取导出模板列表
		sysExportTemplateRouter.GET("export", exportTemplateApi.ExportExcel)                         // 导出Excel
		sysExportTemplateRouter.GET("download", exportTemplateApi.ExportTemplate)                    // 下载模板
		sysExportTemplateRouter.GET("sql/preview", exportTemplateApi.PreviewSQL)                     // 预览SQL
	}
}

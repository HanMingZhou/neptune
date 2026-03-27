package system

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type DictionaryRouter struct{}

func (s *DictionaryRouter) InitSysDictionaryRouter(Router *gin.RouterGroup) {
	sysDictionaryRouter := Router.Group("dictionary").Use(middleware.OperationRecord())
	{
		sysDictionaryRouter.POST("add", dictionaryApi.CreateSysDictionary)    // 新建SysDictionary
		sysDictionaryRouter.POST("delete", dictionaryApi.DeleteSysDictionary) // 删除SysDictionary
		sysDictionaryRouter.POST("update", dictionaryApi.UpdateSysDictionary) // 更新SysDictionary
		sysDictionaryRouter.POST("import", dictionaryApi.ImportSysDictionary) // 导入SysDictionary
		sysDictionaryRouter.GET("export", dictionaryApi.ExportSysDictionary)  // 导出SysDictionary
		sysDictionaryRouter.GET("get", dictionaryApi.FindSysDictionary)       // 根据ID获取SysDictionary
		sysDictionaryRouter.GET("list", dictionaryApi.GetSysDictionaryList)   // 获取SysDictionary列表
	}
}

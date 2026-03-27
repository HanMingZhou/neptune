package system

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type DictionaryDetailRouter struct{}

func (s *DictionaryDetailRouter) InitSysDictionaryDetailRouter(Router *gin.RouterGroup) {
	dictionaryDetailRouter := Router.Group("dictionary/detail").Use(middleware.OperationRecord())
	{
		dictionaryDetailRouter.POST("add", dictionaryDetailApi.CreateSysDictionaryDetail) // 新建SysDictionaryDetail

		dictionaryDetailRouter.POST("delete", dictionaryDetailApi.DeleteSysDictionaryDetail) // 删除SysDictionaryDetail
		dictionaryDetailRouter.POST("update", dictionaryDetailApi.UpdateSysDictionaryDetail) // 更新SysDictionaryDetail

		dictionaryDetailRouter.GET("get", dictionaryDetailApi.FindSysDictionaryDetail)                // 根据ID获取SysDictionaryDetail
		dictionaryDetailRouter.GET("list", dictionaryDetailApi.GetSysDictionaryDetailList)            // 获取SysDictionaryDetail列表
		dictionaryDetailRouter.GET("tree/list", dictionaryDetailApi.GetDictionaryTreeList)            // 获取字典详情树形结构
		dictionaryDetailRouter.GET("tree/type/list", dictionaryDetailApi.GetDictionaryTreeListByType) // 根据字典类型获取字典详情树形结构
		dictionaryDetailRouter.GET("parent/get", dictionaryDetailApi.GetDictionaryDetailsByParent)    // 根据父级ID获取字典详情
		dictionaryDetailRouter.GET("path/get", dictionaryDetailApi.GetDictionaryPath)                 // 获取字典详情的完整路径
	}
}

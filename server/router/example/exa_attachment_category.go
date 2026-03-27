package example

import (
	"github.com/gin-gonic/gin"
)

type AttachmentCategoryRouter struct{}

func (r *AttachmentCategoryRouter) InitAttachmentCategoryRouterRouter(Router *gin.RouterGroup) {
	router := Router.Group("attachment/category")
	{
		router.GET("list", attachmentCategoryApi.GetCategoryList)   // 分类列表
		router.POST("add", attachmentCategoryApi.AddCategory)       // 添加/编辑分类
		router.POST("delete", attachmentCategoryApi.DeleteCategory) // 删除分类
	}
}

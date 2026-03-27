package cms

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type ProductRouter struct{}

func (r *ProductRouter) InitCMSProductRouter(Router *gin.RouterGroup) {
	cmsRouter := Router.Group("cms/product").Use(middleware.OperationRecord())
	cmsWithNotRecorderRouter := Router.Group("cms/product")
	{
		// 产品管理
		cmsRouter.POST("add", cmsApi.ProductApi.CreateProduct)
		cmsRouter.POST("update", cmsApi.ProductApi.UpdateProduct)
		cmsRouter.POST("price/update", cmsApi.ProductApi.UpdatePrice)
		cmsRouter.POST("delete", cmsApi.ProductApi.DeleteProduct)
	}
	{
		cmsWithNotRecorderRouter.GET("list", cmsApi.ProductApi.GetProductList)
		cmsWithNotRecorderRouter.GET("detail", cmsApi.ProductApi.GetProductDetail)
		cmsWithNotRecorderRouter.GET("cluster/list", cmsApi.ProductApi.GetClusterList)
		cmsWithNotRecorderRouter.GET("area/list", cmsApi.ProductApi.GetAreaList)
		cmsWithNotRecorderRouter.GET("node/list", cmsApi.ProductApi.GetClusterNodes) // 获取集群K8s节点列表
	}
}

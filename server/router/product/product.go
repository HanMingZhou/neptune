package product

import (
	"github.com/gin-gonic/gin"
)

type ProductRouter struct{}

func (r *ProductRouter) InitProductRouter(Router *gin.RouterGroup) {
	productWithNotRecorderRouter := Router.Group("product")

	{
		productWithNotRecorderRouter.GET("list", productApi.GetProductList)
		productWithNotRecorderRouter.GET("get", productApi.GetProductDetail)
		productWithNotRecorderRouter.GET("filter/list", productApi.GetProductFilters)
	}
}

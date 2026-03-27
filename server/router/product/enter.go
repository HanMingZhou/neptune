package product

import v1 "gin-vue-admin/api/v1"

type RouterGroup struct {
	ProductRouter
}

var productApi = v1.ApiGroupApp.ProductApiGroup.ProductApi

package product

import "gin-vue-admin/service"

type ApiGroup struct {
	ProductApi
}

var productService = service.ServiceGroupApp.ProductServiceGroup.ProductService

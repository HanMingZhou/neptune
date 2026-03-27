package image

import "gin-vue-admin/service"

type ApiGroup struct {
	ImageApi
}

var imageService = service.ServiceGroupApp.ImageServiceGroup.ImageService

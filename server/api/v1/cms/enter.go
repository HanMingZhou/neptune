package cms

import "gin-vue-admin/service"

type ApiGroup struct {
	ProductApi
	NodeApi
}

var cmsService = service.ServiceGroupApp.CMSServiceGroup

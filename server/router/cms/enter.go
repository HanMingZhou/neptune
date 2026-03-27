package cms

import v1 "gin-vue-admin/api/v1"

type RouterGroup struct {
	ProductRouter
	NodeRouter
}

var cmsApi = v1.ApiGroupApp.CMSApiGroup

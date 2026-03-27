package apisix

import (
	api "gin-vue-admin/api/v1"
)

type RouterGroup struct {
	ApisixRouter
}

var (
	apisixApi = api.ApiGroupApp.ApisixApiGroup.ApisixApi
)

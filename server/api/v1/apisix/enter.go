package apisix

import (
	"gin-vue-admin/service"
	apisixSvc "gin-vue-admin/service/apisix"
)

type ApiGroup struct {
	ApisixApi
}

// 使用 interface 类型，解耦 API 层和 Service 层，方便单元测试 mock
var apisixService apisixSvc.ApisixManager = &service.ServiceGroupApp.ApisixServiceGroup.ApisixService

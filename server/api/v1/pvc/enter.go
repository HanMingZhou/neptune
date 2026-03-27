package pvc

import (
	"gin-vue-admin/service"
	"gin-vue-admin/service/pvc"
)

type ApiGroup struct {
	PVCApi
	VolumeApi
}

// 使用 interface 类型，解耦 API 层和 Service 层，方便单元测试 mock
var pvcService pvc.PVCManager = &service.ServiceGroupApp.PVCServiceGroup.K8sPVCService

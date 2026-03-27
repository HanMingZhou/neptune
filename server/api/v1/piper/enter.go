package piper

import (
	"gin-vue-admin/service"
	"gin-vue-admin/service/piper"
)

type ApiGroup struct {
	PiperApi
}

// 使用 interface 类型，解耦 API 层和 Service 层，方便单元测试 mock
var piperService piper.SSHPiperManager = &service.ServiceGroupApp.PiperServiceGroup.K8sSSHPiperService

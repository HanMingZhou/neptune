package secret

import (
	"gin-vue-admin/service"
	"gin-vue-admin/service/secret"
)

type ApiGroup struct {
	SecretApi
}

// 使用 interface 类型，解耦 API 层和 Service 层，方便单元测试 mock
var secretService secret.SecretManager = &service.ServiceGroupApp.SecretServiceGroup.K8sSecretService

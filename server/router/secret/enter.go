package secret

import api "gin-vue-admin/api/v1"

type RouterGroup struct {
	SecretRouter
}

var (
	secretApi = api.ApiGroupApp.SecretApiGroup.SecretApi
)

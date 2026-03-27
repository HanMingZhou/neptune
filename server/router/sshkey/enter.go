package sshkey

import api "gin-vue-admin/api/v1"

type RouterGroup struct {
	SSHKeyRouter
}

var (
	SSHKeyApi = api.ApiGroupApp.SSHKeyApiGroup.SSHKeyApi
)

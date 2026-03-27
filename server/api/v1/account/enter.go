package account

import "gin-vue-admin/service"

type ApiGroup struct {
	AccountApi
}

var accountService = service.ServiceGroupApp.AccountServiceGroup.AccountService

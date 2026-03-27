package piper

import api "gin-vue-admin/api/v1"

type RouterGroup struct {
	PipeRouter
}

var (
	pipeApi = api.ApiGroupApp.PiperApiGroup.PiperApi
)

package notebook

import (
	"gin-vue-admin/service"
	notebookSvc "gin-vue-admin/service/notebook"
)

type ApiGroup struct {
	NoteBookApi
}

// 使用 interface 类型，解耦 API 层和 Service 层，方便单元测试 mock
var noteBookService notebookSvc.NotebookManager = &service.ServiceGroupApp.NotebookServiceGroup.NotebookService

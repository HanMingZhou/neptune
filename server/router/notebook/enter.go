package notebook

import api "gin-vue-admin/api/v1"

type RouterGroup struct {
	NotebookRouter
}

var (
	notebookApi = api.ApiGroupApp.NotebookApiGroup.NoteBookApi
)

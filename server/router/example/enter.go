package example

import (
	api "gin-vue-admin/api/v1"
)

type RouterGroup struct {
	CustomerRouter
	AttachmentCategoryRouter
	FileUploadAndDownloadRouter
}

var (
	exaCustomerApi              = api.ApiGroupApp.ExampleApiGroup.CustomerApi
	attachmentCategoryApi       = api.ApiGroupApp.ExampleApiGroup.AttachmentCategoryApi
	exaFileUploadAndDownloadApi = api.ApiGroupApp.ExampleApiGroup.FileUploadAndDownloadApi
)

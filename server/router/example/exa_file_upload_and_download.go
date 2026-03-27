package example

import (
	"github.com/gin-gonic/gin"
)

type FileUploadAndDownloadRouter struct{}

func (e *FileUploadAndDownloadRouter) InitFileUploadAndDownloadRouter(Router *gin.RouterGroup) {
	fileUploadAndDownloadRouter := Router.Group("file")
	{
		fileUploadAndDownloadRouter.POST("upload", exaFileUploadAndDownloadApi.UploadFile)                          // 上传文件
		fileUploadAndDownloadRouter.POST("list", exaFileUploadAndDownloadApi.GetFileList)                           // 获取上传文件列表
		fileUploadAndDownloadRouter.POST("delete", exaFileUploadAndDownloadApi.DeleteFile)                          // 删除指定文件
		fileUploadAndDownloadRouter.POST("update", exaFileUploadAndDownloadApi.EditFileName)                        // 编辑文件名或者备注
		fileUploadAndDownloadRouter.POST("breakpoint/continue", exaFileUploadAndDownloadApi.BreakpointContinue)     // 断点续传
		fileUploadAndDownloadRouter.GET("find", exaFileUploadAndDownloadApi.FindFile)                               // 查询当前文件成功的切片
		fileUploadAndDownloadRouter.POST("breakpoint/finish", exaFileUploadAndDownloadApi.BreakpointContinueFinish) // 切片传输完成
		fileUploadAndDownloadRouter.POST("chunk/delete", exaFileUploadAndDownloadApi.RemoveChunk)                   // 删除切片
		fileUploadAndDownloadRouter.POST("import", exaFileUploadAndDownloadApi.ImportURL)                           // 导入URL
	}
}

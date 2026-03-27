package tensorboard

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type TensorBoardRouter struct{}

func (s *TensorBoardRouter) InitTensorBoardApiRouter(PrivateGroup *gin.RouterGroup) {
	apiRouter := PrivateGroup.Group("tensorboard").Use(middleware.OperationRecord())
	apiWithNotRecorderRouter := PrivateGroup.Group("tensorboard")
	{
		apiRouter.POST("add", tensorBoardApi.AddTensorBoard)       // 创建TensorBoard
		apiRouter.POST("delete", tensorBoardApi.DeleteTensorBoard) // 删除TensorBoard
		apiRouter.POST("update", tensorBoardApi.UpdateTensorBoard) // 更新TensorBoard
	}
	{
		apiWithNotRecorderRouter.GET("list", tensorBoardApi.GetTensorBoardList) // 获取TensorBoard列表
	}
}

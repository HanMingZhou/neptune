package tensorboard

import api "gin-vue-admin/api/v1"

type RouterGroup struct {
	TensorBoardRouter
}

var (
	tensorBoardApi = api.ApiGroupApp.TensorBoardApiGroup.TensorBoardApi
)

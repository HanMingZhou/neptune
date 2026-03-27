package tensorboard

import (
	"gin-vue-admin/service"
	"gin-vue-admin/service/tensorboard"
)

type ApiGroup struct {
	TensorBoardApi
}

// 使用 interface 类型，解耦 API 层和 Service 层，方便单元测试 mock
var tensorBoardService tensorboard.TensorboardManager = &service.ServiceGroupApp.TensorBoardServiceGroup.K8sTensorboardService

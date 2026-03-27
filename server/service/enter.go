package service

import (
	"gin-vue-admin/global"
	"gin-vue-admin/service/account"
	"gin-vue-admin/service/apisix"
	"gin-vue-admin/service/cluster"
	"gin-vue-admin/service/cms"
	"gin-vue-admin/service/dashboard"
	"gin-vue-admin/service/example"
	"gin-vue-admin/service/image"
	"gin-vue-admin/service/inference"
	"gin-vue-admin/service/notebook"
	"gin-vue-admin/service/order"
	"gin-vue-admin/service/piper"
	"gin-vue-admin/service/product"
	"gin-vue-admin/service/pvc"
	"gin-vue-admin/service/secret"
	"gin-vue-admin/service/system"
	"gin-vue-admin/service/tensorboard"
	"gin-vue-admin/service/training"
)

// ServiceGroup 聚合所有服务组
type ServiceGroup struct {
	SystemServiceGroup      system.ServiceGroup
	ExampleServiceGroup     example.ServiceGroup
	ClusterServiceGroup     cluster.ServiceGroup
	CMSServiceGroup         cms.ServiceGroup
	ProductServiceGroup     product.ServiceGroup
	ImageServiceGroup       image.ServiceGroup
	NotebookServiceGroup    notebook.ServiceGroup
	PVCServiceGroup         pvc.ServiceGroup
	SecretServiceGroup      secret.ServiceGroup
	PiperServiceGroup       piper.ServiceGroup
	TensorBoardServiceGroup tensorboard.ServiceGroup
	ApisixServiceGroup      apisix.ServiceGroup
	OrderServiceGroup       order.ServiceGroup
	InferenceServiceGroup   inference.ServiceGroup
	DashboardServiceGroup   dashboard.ServiceGroup
	AccountServiceGroup     account.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)

// InitServiceDependencies 初始化服务间依赖
func InitServiceDependencies() {
	ServiceGroupApp.NotebookServiceGroup.NotebookService.SetApisixService(&ServiceGroupApp.ApisixServiceGroup.ApisixService)
	notebook.NotebookServiceApp.SetApisixService(&ServiceGroupApp.ApisixServiceGroup.ApisixService)

	ServiceGroupApp.InferenceServiceGroup.InferenceServiceService.SetApisixService(&ServiceGroupApp.ApisixServiceGroup.ApisixService)
	inference.InferenceServiceServiceApp.SetApisixService(&ServiceGroupApp.ApisixServiceGroup.ApisixService)

	training.TrainingJobServiceApp.SetApisixService(&ServiceGroupApp.ApisixServiceGroup.ApisixService)

	if global.GVA_K8S_CLUSTER_MANAGER != nil {
		global.GVA_K8S_CLUSTER_MANAGER.SetOrderManager(&ServiceGroupApp.OrderServiceGroup.OrderService)
		global.GVA_K8S_CLUSTER_MANAGER.SetProductManager(&ServiceGroupApp.ProductServiceGroup.ProductService)
	}
}

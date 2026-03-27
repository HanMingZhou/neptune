package v1

import (
	"gin-vue-admin/api/v1/account"
	"gin-vue-admin/api/v1/apisix"
	"gin-vue-admin/api/v1/cluster"
	"gin-vue-admin/api/v1/cms"
	"gin-vue-admin/api/v1/dashboard"
	"gin-vue-admin/api/v1/example"
	"gin-vue-admin/api/v1/image"
	"gin-vue-admin/api/v1/inference"
	"gin-vue-admin/api/v1/notebook"
	"gin-vue-admin/api/v1/order"
	"gin-vue-admin/api/v1/piper"
	"gin-vue-admin/api/v1/product"
	"gin-vue-admin/api/v1/pvc"
	"gin-vue-admin/api/v1/secret"
	"gin-vue-admin/api/v1/sshkey"
	"gin-vue-admin/api/v1/system"
	"gin-vue-admin/api/v1/tensorboard"
	"gin-vue-admin/api/v1/training"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	ApisixApiGroup      apisix.ApiGroup
	ClusterApiGroup     cluster.ApiGroup
	CMSApiGroup         cms.ApiGroup
	ExampleApiGroup     example.ApiGroup
	ImageApiGroup       image.ApiGroup
	InferenceApiGroup   inference.ApiGroup
	NotebookApiGroup    notebook.ApiGroup
	PiperApiGroup       piper.ApiGroup
	ProductApiGroup     product.ApiGroup
	PVCApiGroup         pvc.ApiGroup
	SecretApiGroup      secret.ApiGroup
	SSHKeyApiGroup      sshkey.ApiGroup
	SystemApiGroup      system.ApiGroup
	TensorBoardApiGroup tensorboard.ApiGroup
	TrainingApiGroup    training.ApiGroup
	DashboardApiGroup   dashboard.ApiGroup
	AccountApiGroup     account.ApiGroup
	OrderApiGroup       order.ApiGroup
}

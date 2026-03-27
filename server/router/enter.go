package router

import (
	"gin-vue-admin/router/account"
	"gin-vue-admin/router/apisix"
	"gin-vue-admin/router/cluster"
	"gin-vue-admin/router/cms"
	"gin-vue-admin/router/dashboard"
	"gin-vue-admin/router/example"
	"gin-vue-admin/router/image"
	"gin-vue-admin/router/inference"
	"gin-vue-admin/router/notebook"
	"gin-vue-admin/router/order"
	"gin-vue-admin/router/piper"
	"gin-vue-admin/router/product"
	"gin-vue-admin/router/pvc"
	"gin-vue-admin/router/secret"
	"gin-vue-admin/router/sshkey"
	"gin-vue-admin/router/system"
	"gin-vue-admin/router/tensorboard"
	"gin-vue-admin/router/training"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	Apisix      apisix.RouterGroup
	Cluster     cluster.RouterGroup
	CMS         cms.RouterGroup
	Example     example.RouterGroup
	Image       image.RouterGroup
	Inference   inference.RouterGroup
	NoteBook    notebook.RouterGroup
	Piper       piper.RouterGroup
	Product     product.RouterGroup
	PVC         pvc.RouterGroup
	Secret      secret.RouterGroup
	System      system.RouterGroup
	TensorBoard tensorboard.RouterGroup
	SSHKey      sshkey.RouterGroup
	Training    training.RouterGroup
	Dashboard   dashboard.RouterGroup
	Account     account.RouterGroup
	Order       order.RouterGroup
}

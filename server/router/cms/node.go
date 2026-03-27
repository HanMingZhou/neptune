package cms

import (
	v1 "gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type NodeRouter struct{}

func (r *NodeRouter) InitNodeRouter(Router *gin.RouterGroup) {
	nodeRouter := Router.Group("cms").Use(middleware.OperationRecord())
	nodeWithNotRecorderRouter := Router.Group("cms")
	nodeApi := v1.ApiGroupApp.CMSApiGroup.NodeApi
	{
		nodeRouter.POST("node/uncordon", nodeApi.UncordonNode) // 恢复节点调度
		nodeRouter.POST("node/drain", nodeApi.DrainNode)       // 驱逐节点
	}
	{
		nodeWithNotRecorderRouter.POST("node/list", nodeApi.GetClusterNodes) // 获取集群节点列表
	}
}

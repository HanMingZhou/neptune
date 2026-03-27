package cluster

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type ClusterRouter struct{}

// InitClusterRouter 初始化集群管理路由
func (r *ClusterRouter) InitClusterRouter(Router *gin.RouterGroup) {
	clusterRouter := Router.Group("cms/cluster").Use(middleware.OperationRecord())
	clusterWithNotRecorderRouter := Router.Group("cms/cluster")
	clusterApiGroup := clusterApi.ClusterApi
	{
		clusterRouter.POST("add", clusterApiGroup.CreateCluster)    // 创建集群
		clusterRouter.POST("update", clusterApiGroup.UpdateCluster) // 更新集群
		clusterRouter.POST("delete", clusterApiGroup.DeleteCluster) // 删除集群
	}
	{
		clusterWithNotRecorderRouter.POST("list", clusterApiGroup.GetClusterList) // 获取集群列表
	}
}

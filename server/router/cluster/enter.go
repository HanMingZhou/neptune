package cluster

import v1 "gin-vue-admin/api/v1"

type RouterGroup struct {
	ClusterRouter
}

var clusterApi = v1.ApiGroupApp.ClusterApiGroup

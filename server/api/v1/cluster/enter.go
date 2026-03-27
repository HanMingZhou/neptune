package cluster

import "gin-vue-admin/service"

type ApiGroup struct {
	ClusterApi
}

var clusterService = service.ServiceGroupApp.ClusterServiceGroup.ClusterService

package pvc

import api "gin-vue-admin/api/v1"

type RouterGroup struct {
	PVCRouter
	VolumeRouter
}

var (
	pvcApi    = api.ApiGroupApp.PVCApiGroup.PVCApi
	volumeApi = api.ApiGroupApp.PVCApiGroup.VolumeApi
)

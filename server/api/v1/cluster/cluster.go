package cluster

import (
	"gin-vue-admin/model/cluster/request"
	"gin-vue-admin/model/common/response"
	"gin-vue-admin/utils"

	"github.com/gin-gonic/gin"
)

type ClusterApi struct{}

// GetClusterList 获取集群列表
// @Tags CMS-Cluster
// @Summary 获取集群列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param keyword query string false "搜索关键词"
// @Param status query int false "状态筛选"
// @Success 200 {object} response.Response{data=object,msg=string}
// @Router /api/v1/cms/cluster/list [post]
func (a *ClusterApi) GetClusterList(c *gin.Context) {
	var req request.GetClusterListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	list, err := clusterService.GetClusterList(&req)
	if err != nil {
		utils.HandleError(c, err, "获取集群列表失败")
		return
	}

	response.OkWithDetailed(list, "获取成功", c)
}

// CreateCluster 创建集群
// @Tags CMS-Cluster
// @Summary 创建集群
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.CreateClusterReq true "创建集群请求"
// @Success 200 {object} response.Response{msg=string}
// @Router /api/v1/cms/cluster/add [post]
func (a *ClusterApi) CreateCluster(c *gin.Context) {
	var req request.CreateClusterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	if err := clusterService.CreateCluster(&req); err != nil {
		utils.HandleError(c, err, "创建集群失败")
		return
	}

	response.OkWithMessage("创建成功", c)
}

// UpdateCluster 更新集群
// @Tags CMS-Cluster
// @Summary 更新集群
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.UpdateClusterReq true "更新集群请求"
// @Success 200 {object} response.Response{msg=string}
// @Router /api/v1/cms/cluster/update [post]
func (a *ClusterApi) UpdateCluster(c *gin.Context) {
	var req request.UpdateClusterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	if err := clusterService.UpdateCluster(&req); err != nil {
		utils.HandleError(c, err, "更新集群失败")
		return
	}

	response.OkWithMessage("更新成功", c)
}

// DeleteCluster 删除集群
// @Tags CMS-Cluster
// @Summary 删除集群
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.DeleteClusterReq true "删除集群请求"
// @Success 200 {object} response.Response{msg=string}
// @Router /api/v1/cms/cluster/delete [post]
func (a *ClusterApi) DeleteCluster(c *gin.Context) {
	var req request.DeleteClusterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	if err := clusterService.DeleteCluster(&req); err != nil {
		utils.HandleError(c, err, "删除集群失败")
		return
	}

	response.OkWithMessage("删除成功", c)
}

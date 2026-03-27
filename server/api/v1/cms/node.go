package cms

import (
	"gin-vue-admin/model/cms/request"
	cmsRes "gin-vue-admin/model/cms/response"
	"gin-vue-admin/model/common/response"
	"gin-vue-admin/service"
	"gin-vue-admin/utils"

	"github.com/gin-gonic/gin"
)

type NodeApi struct{}

// GetClusterNodes 获取集群节点列表
// @Tags CMS-Node
// @Summary 获取集群节点列表
// @Accept json
// @Produce json
// @Param clusterId query int true "集群ID"
// @Success 200 {object} response.Response{}
// @Router /api/v1/cms/node/list [post]
func (api *NodeApi) GetClusterNodes(c *gin.Context) {
	var req request.GetClusterNodesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
		return
	}

	nodes, err := service.ServiceGroupApp.CMSServiceGroup.NodeService.GetClusterNodesWithResources(c.Request.Context(), req.ClusterId, req.Keyword)
	if err != nil {
		utils.HandleError(c, err, "获取集群节点失败")
		return
	}

	response.OkWithDetailed(cmsRes.NodeListResponse{
		Nodes: nodes,
		Total: len(nodes),
	}, "获取成功", c)
}

// UncordonNode 恢复节点调度
// @Tags CMS-Node
// @Summary 恢复节点调度
// @Accept json
// @Produce json
// @Param clusterId body int true "集群ID"
// @Param nodeName body string true "节点名称"
// @Success 200 {object} response.Response{}
// @Router /api/v1/cms/node/uncordon [post]
func (api *NodeApi) UncordonNode(c *gin.Context) {
	var req request.UncordonNodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
		return
	}

	err := service.ServiceGroupApp.CMSServiceGroup.NodeService.UncordonNode(c.Request.Context(), req.ClusterId, req.NodeName)
	if err != nil {
		utils.HandleError(c, err, "恢复节点调度失败")
		return
	}

	response.OkWithMessage("恢复成功", c)
}

// DrainNode 驱逐节点
// @Tags CMS-Node
// @Summary 驱逐节点
// @Accept json
// @Produce json
// @Param clusterId body int true "集群ID"
// @Param nodeName body string true "节点名称"
// @Success 200 {object} response.Response{}
// @Router /api/v1/cms/node/drain [post]
func (api *NodeApi) DrainNode(c *gin.Context) {
	var req request.DrainNodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
		return
	}

	err := service.ServiceGroupApp.CMSServiceGroup.NodeService.DrainNode(c.Request.Context(), req.ClusterId, req.NodeName)
	if err != nil {
		utils.HandleError(c, err, "驱逐节点失败")
		return
	}

	response.OkWithMessage("驱逐成功", c)
}

package cms

import (
	cmsReq "gin-vue-admin/model/cms/request"
	"gin-vue-admin/model/common/response"
	"gin-vue-admin/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type ProductApi struct{}

// CreateProduct 创建产品
// @Tags CMS-Product
// @Summary 创建产品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cmsReq.CreateProductReq true "创建产品请求"
// @Success 200 {object} response.Response{msg=string}
// @Router /api/v1/cms/product/add [post]
func (p *ProductApi) CreateProduct(c *gin.Context) {
	var req cmsReq.CreateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	if err := cmsService.ProductService.CreateProduct(&req); err != nil {
		utils.HandleError(c, err, "创建产品失败")
		return
	}

	response.OkWithMessage("创建成功", c)
}

// UpdateProduct 更新产品
// @Tags CMS-Product
// @Summary 更新产品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cmsReq.UpdateProductReq true "更新产品请求"
// @Success 200 {object} response.Response{msg=string}
// @Router /api/v1/cms/product/update [post]
func (p *ProductApi) UpdateProduct(c *gin.Context) {
	var req cmsReq.UpdateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	if err := cmsService.ProductService.UpdateProduct(&req); err != nil {
		utils.HandleError(c, err, "更新产品失败")
		return
	}

	response.OkWithMessage("更新成功", c)
}

// UpdatePrice 更新价格
// @Tags CMS-Product
// @Summary 更新产品价格
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cmsReq.UpdatePriceReq true "更新价格请求"
// @Success 200 {object} response.Response{msg=string}
// @Router /api/v1/cms/product/price/update [post]
func (p *ProductApi) UpdatePrice(c *gin.Context) {
	var req cmsReq.UpdatePriceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	if err := cmsService.ProductService.UpdatePrice(&req); err != nil {
		utils.HandleError(c, err, "更新价格失败")
		return
	}

	response.OkWithMessage("更新成功", c)
}

// DeleteProduct 删除产品
// @Tags CMS-Product
// @Summary 删除产品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body cmsReq.DeleteProductReq true "删除产品请求"
// @Success 200 {object} response.Response{msg=string}
// @Router /api/v1/cms/product/delete [post]
func (p *ProductApi) DeleteProduct(c *gin.Context) {
	var req cmsReq.DeleteProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	if err := cmsService.ProductService.DeleteProduct(&req); err != nil {
		utils.HandleError(c, err, "删除产品失败")
		return
	}

	response.OkWithMessage("删除成功", c)
}

// GetProductList 获取产品列表
// @Tags CMS-Product
// @Summary 获取产品列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param page query int false "页码"
// @Param pageSize query int false "每页大小"
// @Param clusterId query int false "集群ID"
// @Param area query string false "地区"
// @Param gpuModel query string false "GPU型号"
// @Param status query int false "状态"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} response.Response{data=object,msg=string}
// @Router /api/v1/cms/product/list [get]
func (p *ProductApi) GetProductList(c *gin.Context) {
	var req cmsReq.GetProductListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	list, err := cmsService.ProductService.GetProductList(&req)
	if err != nil {
		utils.HandleError(c, err, "获取产品列表失败")
		return
	}

	response.OkWithDetailed(list, "获取成功", c)
}

// GetProductDetail 获取产品详情
// @Tags CMS-Product
// @Summary 获取产品详情
// @Security ApiKeyAuth
// @Param id query int true "产品ID"
// @Success 200 {object} response.Response{data=object,msg=string}
// @Router /api/v1/cms/product/detail [get]
func (p *ProductApi) GetProductDetail(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.FailWithMessage("产品ID不能为空", c)
		return
	}

	product, err := cmsService.ProductService.GetProductDetail(idStr)
	if err != nil {
		utils.HandleError(c, err, "获取产品详情失败")
		return
	}

	response.OkWithDetailed(product, "获取成功", c)
}

// GetClusterList 获取集群列表（用于下拉选择）
// @Tags CMS-Product
// @Summary 获取集群列表
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=object,msg=string}
// @Router /api/v1/cms/product/cluster/list [get]
func (p *ProductApi) GetClusterList(c *gin.Context) {
	list, err := cmsService.ProductService.GetClusterList()
	if err != nil {
		utils.HandleError(c, err, "获取集群列表失败")
		return
	}

	response.OkWithDetailed(list, "获取成功", c)
}

// GetAreaList 获取地区列表（用于下拉选择）
// @Tags CMS-Product
// @Summary 获取地区列表
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=object,msg=string}
// @Router /api/v1/cms/product/area/list [get]
func (p *ProductApi) GetAreaList(c *gin.Context) {
	list, err := cmsService.ProductService.GetAreaList()
	if err != nil {
		utils.HandleError(c, err, "获取地区列表失败")
		return
	}

	response.OkWithDetailed(list, "获取成功", c)
}

// GetClusterNodes 获取集群下的K8s节点列表
// @Tags CMS-Product
// @Summary 获取集群K8s节点列表
// @Security ApiKeyAuth
// @Param clusterId query int true "集群ID"
// @Success 200 {object} response.Response{data=object,msg=string}
// @Router /api/v1/cms/product/node/list [get]
func (p *ProductApi) GetClusterNodes(c *gin.Context) {
	clusterIdStr := c.Query("clusterId")
	if clusterIdStr == "" {
		response.FailWithMessage("集群ID不能为空", c)
		return
	}

	clusterId := cast.ToUint(clusterIdStr)
	if clusterId == 0 {
		response.FailWithMessage("集群ID格式错误", c)
		return
	}

	nodes, err := cmsService.NodeService.GetClusterNodes(clusterId, "")
	if err != nil {
		utils.HandleError(c, err, "获取集群节点失败")
		return
	}

	response.OkWithDetailed(nodes, "获取成功", c)
}

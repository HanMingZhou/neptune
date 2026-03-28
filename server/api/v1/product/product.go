package product

import (
	"context"
	"fmt"
	"gin-vue-admin/model/common/response"
	productSvc "gin-vue-admin/service/product"
	"gin-vue-admin/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type ProductApi struct{}

// GetProductList 获取产品列表
// @Tags Product
// @Summary 获取产品列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param page query int false "页码"
// @Param pageSize query int false "每页大小"
// @Success 200 {object} response.Response{data=object,msg=string}
// @Router /product/list [get]
func (p *ProductApi) GetProductList(c *gin.Context) {
	var page, pageSize, productType int
	page = 1
	pageSize = 100

	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}
	if ps := c.Query("pageSize"); ps != "" {
		pageSize = parseInt(ps, 100)
	}
	if pt := c.Query("productType"); pt != "" {
		productType = parseInt(pt, 0)
	}

	filters := aggregateProductFiltersFromContext(c)

	products, err := productService.GetProductList(context.Background(), page, pageSize, productType, filters)
	if err != nil {
		utils.HandleError(c, err, "查询产品列表失败")
		return
	}

	response.OkWithDetailed(products, "获取成功", c)
}

// GetAggregateProductList 获取集群级聚合商品列表
func (p *ProductApi) GetAggregateProductList(c *gin.Context) {
	var page, pageSize, productType int
	page = 1
	pageSize = 100

	if pageQuery := c.Query("page"); pageQuery != "" {
		page = parseInt(pageQuery, 1)
	}
	if pageSizeQuery := c.Query("pageSize"); pageSizeQuery != "" {
		pageSize = parseInt(pageSizeQuery, 100)
	}
	if productTypeQuery := c.Query("productType"); productTypeQuery != "" {
		productType = parseInt(productTypeQuery, 0)
	}

	products, err := productService.GetAggregateProductList(
		context.Background(),
		page,
		pageSize,
		productType,
		aggregateProductFiltersFromContext(c),
	)
	if err != nil {
		utils.HandleError(c, err, "查询聚合商品列表失败")
		return
	}

	response.OkWithDetailed(products, "获取成功", c)
}

// GetProductDetail 获取产品详情
// @Tags Product
// @Summary 获取产品详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id query int true "产品ID"
// @Success 200 {object} response.Response{data=object,msg=string}
// @Router /product/detail [get]
func (p *ProductApi) GetProductDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.FailWithMessage("产品ID不能为空", c)
		return
	}

	toUint := cast.ToUint(id)
	product, err := productService.GetProductById(context.Background(), toUint)
	if err != nil {
		utils.HandleError(c, err, "查询产品详情失败")
		return
	}

	response.OkWithDetailed(product, "获取成功", c)
}

// GetProductFilters 获取产品筛选条件
// @Tags Product
// @Summary 获取产品筛选条件（地区、GPU型号等）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string}
// @Router /product/filters [get]
func (p *ProductApi) GetProductFilters(c *gin.Context) {
	productType := 0
	if pt := c.Query("productType"); pt != "" {
		productType = parseInt(pt, 0)
	}
	filters, err := productService.GetProductFilters(context.Background(), productType)
	if err != nil {
		utils.HandleError(c, err, "获取筛选条件失败")
		return
	}

	response.OkWithDetailed(filters, "获取成功", c)
}

func parseInt(s string, defaultVal int) int {
	var v int
	if _, err := fmt.Sscanf(s, "%d", &v); err != nil {
		return defaultVal
	}
	return v
}

func aggregateProductFiltersFromContext(c *gin.Context) productSvc.AggregateProductFilters {
	return productSvc.AggregateProductFilters{
		Area:     c.Query("area"),
		GPUModel: c.Query("gpuModel"),
		CPUModel: c.Query("cpuModel"),
	}
}

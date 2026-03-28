package order

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model/common/response"
	"gin-vue-admin/model/order/request"
	"gin-vue-admin/service"
	"gin-vue-admin/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrderApi struct{}

var orderService = service.ServiceGroupApp.OrderServiceGroup.OrderService

// GetOrderOverview 获取财务总览
// @Tags Order
// @Summary 获取财务总览
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=orderResp.OrderOverviewResp} "成功"
// @Router /api/v1/order/overview [get]
func (a *OrderApi) GetOrderOverview(c *gin.Context) {
	userId := utils.GetUserID(c)
	resp, err := orderService.GetOrderOverview(c, userId)
	if err != nil {
		global.GVA_LOG.Error("获取财务总览失败!", zap.Error(err))
		response.FailWithMessage("获取财务总览失败", c)
		return
	}
	response.OkWithDetailed(resp, "获取成功", c)
}

// GetUsageList 获取使用详情列表
// @Tags Order
// @Summary 获取使用详情列表
// @accept application/json
// @Produce application/json
// @Param data body request.GetOrderUsageListReq true "分页和过滤条件"
// @Success 200 {object} response.Response{data=orderResp.UsageListResp} "成功"
// @Router /api/v1/order/usage/list [post]
func (a *OrderApi) GetUsageList(c *gin.Context) {
	var req request.GetOrderUsageListReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	list, total, err := orderService.GetUsageList(c, userId, req)
	if err != nil {
		global.GVA_LOG.Error("获取使用详情失败!", zap.Error(err))
		response.FailWithMessage("获取使用详情失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// GetTransactionList 获取收支明细列表
// @Tags Order
// @Summary 获取收支明细列表
// @accept application/json
// @Produce application/json
// @Param data body request.GetTransactionListReq true "分页和过滤条件"
// @Success 200 {object} response.Response{data=orderResp.TransactionListResp} "成功"
// @Router /api/v1/order/transaction/list [post]
func (a *OrderApi) GetTransactionList(c *gin.Context) {
	var req request.GetTransactionListReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	list, total, err := orderService.GetTransactionList(c, userId, req)
	if err != nil {
		global.GVA_LOG.Error("获取收支明细失败!", zap.Error(err))
		response.FailWithMessage("获取收支明细失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// RechargeBalance 余额充值
// @Tags Order
// @Summary 余额充值
// @accept application/json
// @Produce application/json
// @Param data body request.RechargeBalanceReq true "充值信息"
// @Success 200 {object} response.Response{data=orderResp.RechargeBalanceResp} "成功"
// @Router /api/v1/order/recharge [post]
func (a *OrderApi) RechargeBalance(c *gin.Context) {
	var req request.RechargeBalanceReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)
	authorityId := utils.GetUserAuthorityId(c)
	resp, err := orderService.RechargeBalance(c, userId, authorityId, req)
	if err != nil {
		global.GVA_LOG.Error("充值失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithDetailed(resp, "充值成功", c)
}

// GetInvoiceList 获取发票列表
// @Tags Order
// @Summary 获取发票列表
// @accept application/json
// @Produce application/json
// @Param data body request.GetInvoiceListReq true "分页"
// @Success 200 {object} response.Response{data=orderResp.InvoiceListResp} "成功"
// @Router /api/v1/order/invoice/list [post]
func (a *OrderApi) GetInvoiceList(c *gin.Context) {
	var req request.GetInvoiceListReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	list, total, err := orderService.GetInvoiceList(c, userId, req)
	if err != nil {
		global.GVA_LOG.Error("获取发票列表失败!", zap.Error(err))
		response.FailWithMessage("获取发票列表失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// GetOrderList 获取用户订单列表
// @Tags Order
// @Summary 获取用户订单列表
// @accept application/json
// @Produce application/json
// @Param data body request.GetOrderUsageListReq true "分页和过滤条件"
// @Success 200 {object} response.Response{data=response.PageResult} "成功"
// @Router /api/v1/order/order/list [post]
func (a *OrderApi) GetOrderList(c *gin.Context) {
	var req request.GetOrderListReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	list, total, err := orderService.GetOrderList(c, int64(userId), req)
	if err != nil {
		global.GVA_LOG.Error("获取订单列表失败!", zap.Error(err))
		response.FailWithMessage("获取订单列表失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// ApplyInvoice 申请发票
// @Tags Order
// @Summary 申请发票
// @accept application/json
// @Produce application/json
// @Param data body request.ApplyInvoiceReq true "发票详情"
// @Success 200 {object} response.Response{} "成功"
// @Router /api/v1/order/invoice/apply [post]
func (a *OrderApi) ApplyInvoice(c *gin.Context) {
	var req request.ApplyInvoiceReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	err = orderService.ApplyInvoice(c, userId, req)
	if err != nil {
		global.GVA_LOG.Error("申请发票失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("申请提交成功", c)
}

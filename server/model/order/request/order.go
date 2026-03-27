package request

import "gin-vue-admin/model/common/request"

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	UserId       uint
	ProduceType  int64 // 订单类型：1-计算资源，2-文件存储
	ResourceType int64 // 资源类型：1-notebook，2-training，3-inference, 4-storage
	ResourceId   uint
	ProductId    uint
	ImageId      uint
	PayType      int64 // 支付类型：1-wx 2-alipay 3-余额 4-系统
	ChargeType   int64 // 计费类型：1-按量付费，2-包日，3-包周，4-包月
	Quantity     int   // 时长：按量付费为1，包日/周/月为对应数量
	Area         string
	ClusterId    uint
	Remark       string
}

// GetOrderListReq 获取订单列表请求
type GetOrderListReq struct {
	request.PageInfo
	Status *int `json:"status" form:"status"` // 0-待支付 1-已支付 2-已取消 3-已退款 4-已过期
}

// GetOrderUsageListReq 计费流水查询请求
type GetOrderUsageListReq struct {
	request.PageInfo
	StartTime string `json:"startTime" form:"startTime"`
	EndTime   string `json:"endTime" form:"endTime"`
}

// GetTransactionListReq 交易流水查询请求
type GetTransactionListReq struct {
	request.PageInfo
	Type      int    `json:"type" form:"type"` // 1-充值 2-消费 3-退款 4-系统调账
	StartTime string `json:"startTime" form:"startTime"`
	EndTime   string `json:"endTime" form:"endTime"`
}

// GetInvoiceListReq 发票列表请求
type GetInvoiceListReq struct {
	request.PageInfo
}

// ApplyInvoiceReq 申请发票请求
type ApplyInvoiceReq struct {
	Amount    float64 `json:"amount" binding:"required"`
	TitleId   uint    `json:"titleId" binding:"required"`
	AddressId uint    `json:"addressId" binding:"required"`
}

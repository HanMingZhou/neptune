package response

import "gin-vue-admin/model/order"

// OrderOverviewResp 账单概览响应
type OrderOverviewResp struct {
	Balance             float64           `json:"balance"`             // 账户余额
	MtdEstimate         float64           `json:"mtdEstimate"`         // 本月预估消费
	LastMonthSettlement float64           `json:"lastMonthSettlement"` // 上月结算金额
	CompareLastMonth    float64           `json:"compareLastMonth"`    // 与上月对比
	Status              string            `json:"status"`              // 状态
	ConsumptionTrend    []ConsumptionData `json:"consumptionTrend"`    // 消费趋势（7天）
	MonthlyRanking      []RankingData     `json:"monthlyRanking"`      // 月度排名（Top 4）
}

// ConsumptionData 消费趋势数据（别名 TrendData）
type ConsumptionData struct {
	Date   string `json:"date"`   // 日期 (MM-DD 格式)
	Amount string `json:"amount"` // 金额（格式化为2位小数的字符串）
}

// TrendData 趋势数据（与 ConsumptionData 相同）
type TrendData = ConsumptionData

// RankingData 排名数据
type RankingData struct {
	Name    string  `json:"name"`    // 资源名称
	Amount  string  `json:"amount"`  // 金额（格式化为2位小数的字符串）
	Percent float64 `json:"percent"` // 百分比
}

// GetUsageListResponse 使用列表响应
type GetUsageListResponse struct {
	Records []order.Order `json:"records"` // 消费记录列表
	Total   int64         `json:"total"`   // 总记录数
}

// UsageListResp 使用列表响应（旧版兼容）
type UsageListResp struct {
	List  []UnifiedUsageItem `json:"list"`  // 使用详情列表
	Total int64              `json:"total"` // 总记录数
}

// UnifiedUsageItem 统一使用详情条目（合并按量计费和预付费数据源）
type UnifiedUsageItem struct {
	ID             uint    `json:"id"`
	Source         string  `json:"source"`         // 数据来源: "按量计费" / "预付费"
	ResourceType   int64   `json:"resourceType"`   // 资源类型(1-notebook 2-training 3-inference 4-storage)
	ResourceName   string  `json:"resourceName"`   // 实例名称（容器实例/训练任务/推理服务名称）
	ChargeType     int64   `json:"chargeType"`     // 计费类型(1-按量 2-包日 3-包周 4-包月)
	ChargeTypeName string  `json:"chargeTypeName"` // 计费类型名称
	UnitPrice      float64 `json:"unitPrice"`      // 单价
	Duration       int64   `json:"duration"`       // 使用时长(秒) - 按量计费用
	Quantity       int     `json:"quantity"`       // 数量/时长 - 预付费用
	Amount         float64 `json:"amount"`         // 费用金额
	Status         int     `json:"status"`         // 状态
	StatusText     string  `json:"statusText"`     // 状态文本
	StartTime      string  `json:"startTime"`      // 开始时间
	EndTime        string  `json:"endTime"`        // 结束时间（预付费到期时间）
	CreatedAt      string  `json:"createdAt"`      // 创建时间
	UpdatedAt      string  `json:"updatedAt"`      // 更新时间
	OrderNo        string  `json:"orderNo"`        // 订单号（用于关联交易流水）
}

type TransactionDetail struct {
	order.Transaction
	OrderStatus  int    `json:"orderStatus"`
	OrderNo      string `json:"orderNo"`      // 关联订单号
	ResourceType int64  `json:"resourceType"` // 资源类型(1-notebook 2-training 3-inference 4-storage)
	ResourceName string `json:"resourceName"` // 资源名称（批量查询填充）
	ChargeType   int64  `json:"chargeType"`   // 计费类型
}

type TransactionListResp struct {
	List  []TransactionDetail `json:"list"`
	Total int64               `json:"total"`
}

type InvoiceListResp struct {
	List  []order.Invoice `json:"list"`  // 发票列表
	Total int64           `json:"total"` // 总记录数
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`    // 错误代码
	Message string `json:"message"` // 错误消息
}

// ProductTypeConsumption 产品类型消费数据
type ProductTypeConsumption struct {
	ProductType int64   `json:"productType"` // 产品类型：1-计算资源 2-存储资源
	TypeName    string  `json:"typeName"`    // 类型名称
	Amount      float64 `json:"amount"`      // 总金额
}

// GetConsumptionByProductTypeResp 按产品类型查询消费响应
type GetConsumptionByProductTypeResp struct {
	Data map[int64]float64 `json:"data"` // 产品类型 -> 总金额的映射
}

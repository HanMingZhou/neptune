package order

import (
	"time"

	"gorm.io/gorm"
)

// 订单状态常量
const (
	OrderStatusCreating  = -1 // 资源创建中
	OrderStatusPending   = 0  // 待支付
	OrderStatusPaid      = 1  // 已支付/已完成
	OrderStatusCancelled = 2  // 已取消
	OrderStatusRefunded  = 3  // 已退款
	OrderStatusExpired   = 4  // 已过期
	OrderStatusFailed    = 5  // 创建失败
)

// 订单类型常量
const (
	OrderTypeNotebook  = 1 // 容器实例订单
	OrderTypeTraining  = 2 // 训练任务订单
	OrderTypeInference = 3 // 推理服务订单
	OrderTypeStorage   = 4 // 存储订单
)

var OrderTypeToString = map[int]string{
	OrderTypeNotebook:  "容器实例",
	OrderTypeTraining:  "训练任务",
	OrderTypeInference: "推理服务",
	OrderTypeStorage:   "存储",
}

// 产品类型常量
const (
	ProductTypeCompute = 1 // 计算资源
	ProductTypeStorage = 2 // 存储资源
)

var ProductTypeToString = map[int]string{
	ProductTypeCompute: "计算资源",
	ProductTypeStorage: "存储资源",
}

// 计费状态常量
const (
	// Common 状态
	// 资源创建后：使用中
	OrderStatusInUse = 1 // 使用中

	// 按量状态
	// 按量资源停止/删除后：已停止 ——> 扣费完成后:已结算
	OrderStatusStopped  = 2 // 已停止
	OrderStatusFinished = 3 // 已结算

	// 预付费状态
	// 预付费资源到期后：已到期
	// OrderStatusExpired = 4 // 已到期
)

var OrderStatusToString = map[int]string{
	OrderStatusInUse:    "使用中",
	OrderStatusStopped:  "已停止",
	OrderStatusFinished: "已结算",
	OrderStatusExpired:  "已到期",
}

// 交易类型常量
const (
	TransactionTypeRecharge    = 1 // 充值
	TransactionTypeConsumption = 2 // 消费
	TransactionTypeRefund      = 3 // 退款
	TransactionTypeAdjustment  = 4 // 系统调账
)

// 发票状态常量
const (
	InvoiceStatusProcessing = "PROCESSING" // 开票中
	InvoiceStatusSent       = "SENT"       // 已寄出
)

// 发票类型常量
const (
	InvoiceTypePersonal   = "PERSONAL"   // 个人发票
	InvoiceTypeEnterprise = "ENTERPRISE" // 企业发票
)

// 支付方式常量
const (
	PayMethodAlipay  = 1 // 支付宝
	PayMethodWechat  = 2 // 微信
	PayMethodBalance = 3 // 余额
	PayMethodSystem  = 4 // 系统
)

// Order 表
type Order struct {
	gorm.Model

	// === 核心标识字段 ===
	RecordNo string `json:"recordNo" gorm:"column:record_no;size:64;uniqueIndex;comment:记录编号(唯一)"`
	UserId   uint   `json:"userId" gorm:"column:user_id;index:idx_user_time;comment:用户ID"`

	// === 来源标识 ===
	OrderNo string `json:"orderNo" gorm:"column:order_no;size:64;index:idx_order_no;comment:订单号"`

	// === 资源信息 ===
	ResourceType int64 `json:"resourceType" gorm:"column:resource_type;index:idx_resource;comment:资源类型(1-notebook 2-training 3-inference 4-storage)"`
	ResourceId   uint  `json:"resourceId" gorm:"column:resource_id;index:idx_resource_id;comment:资源ID"`
	ProductId    uint  `json:"productId" gorm:"column:product_id;comment:产品ID"`
	ProductType  int64 `json:"productType" gorm:"column:product_type;comment:产品类型(1-计算资源 2-存储资源)"`
	ImageId      uint  `json:"imageId" gorm:"column:image_id;comment:镜像ID"`

	// === 计费信息 ===
	ChargeType     int64   `json:"chargeType" gorm:"column:charge_type;comment:计费类型(1-按量 2-包日 3-包周 4-包月)"`
	UnitPrice      float64 `json:"unitPrice" gorm:"column:unit_price;type:decimal(20,6);comment:单价"`
	Quantity       int     `json:"quantity" gorm:"column:quantity;default:0;comment:数量(预付费用)"`
	Duration       int64   `json:"duration" gorm:"column:duration;default:0;comment:累计时长秒数(按量计费用)"`
	Amount         float64 `json:"amount" gorm:"column:amount;type:decimal(20,6);comment:金额"`
	DiscountAmount float64 `json:"discountAmount" gorm:"column:discount_amount;type:decimal(20,6);default:0;comment:优惠金额"`

	// === 时间信息 ===
	StartTime  time.Time  `json:"startTime" gorm:"column:start_time;index:idx_user_time;index:idx_time_range;comment:服务生效时间"`
	EndTime    *time.Time `json:"endTime" gorm:"column:end_time;comment:服务到期时间（预付费）"`
	PayTime    *time.Time `json:"payTime" gorm:"column:pay_time;comment:支付时间(预付费)"`
	SettleTime *time.Time `json:"settleTime" gorm:"column:settle_time;comment:结算时间(按量计费)"`

	// === 状态信息 ===
	Status int `json:"status" gorm:"column:status;index:idx_status;comment:状态(1-使用中 2-已停止 3-已结算 4-已到期)"`

	// === 其他信息 ===
	Area      string `json:"area" gorm:"column:area;size:50;comment:区域"`
	ClusterId uint   `json:"clusterId" gorm:"column:cluster_id;comment:集群ID"`
	Remark    string `json:"remark" gorm:"column:remark;size:500;comment:备注"`
}

func (Order) TableName() string {
	return "orders"
}

// OrderSummary 账单汇总表 - 按天/月汇总用户的费用
type OrderSummary struct {
	gorm.Model
	UserId       uint    `json:"userId" gorm:"column:user_id;index;comment:用户ID"`
	OrderDate    string  `json:"orderDate" gorm:"column:order_date;size:10;index;comment:账单日期(YYYY-MM-DD)"`
	OrderMonth   string  `json:"orderMonth" gorm:"column:order_month;size:7;index;comment:账单月份(YYYY-MM)"`
	TotalAmount  float64 `json:"totalAmount" gorm:"column:total_amount;type:decimal(20,6);default:0;comment:总金额"`
	PaidAmount   float64 `json:"paidAmount" gorm:"column:paid_amount;type:decimal(20,6);default:0;comment:已支付金额"`
	UnpaidAmount float64 `json:"unpaidAmount" gorm:"column:unpaid_amount;type:decimal(20,6);default:0;comment:未支付金额"`
	RefundAmount float64 `json:"refundAmount" gorm:"column:refund_amount;type:decimal(20,6);default:0;comment:退款金额"`
}

func (OrderSummary) TableName() string {
	return "order_summaries"
}

// Transaction 交易流水表 - 记录所有充值、消费、退款等资金变动
type Transaction struct {
	gorm.Model
	TransactionNo string    `json:"transactionNo" gorm:"column:transaction_no;size:64;uniqueIndex;comment:交易流水号"`
	UserId        uint      `json:"userId" gorm:"column:user_id;index;comment:用户ID"`
	Type          int       `json:"type" gorm:"column:type;comment:交易类型(1-充值 2-消费 3-退款 4-系统调账)"`
	Amount        float64   `json:"amount" gorm:"column:amount;type:decimal(20,6);comment:交易金额(正数为收入,负数为支出)"`
	BalanceBefore float64   `json:"balanceBefore" gorm:"column:balance_before;type:decimal(20,6);comment:交易前余额"`
	BalanceAfter  float64   `json:"balanceAfter" gorm:"column:balance_after;type:decimal(20,6);comment:交易后余额"`
	OrderId       uint      `json:"orderId" gorm:"column:order_id;comment:关联订单ID"`
	Method        int64     `json:"method" gorm:"column:method;comment:支付方式(1-支付宝 2-微信 3-余额 4-系统)"`
	Remark        string    `json:"remark" gorm:"column:remark;size:500;comment:备注"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at;comment:创建时间"`
}

func (Transaction) TableName() string {
	return "transactions"
}

// Wallet 用户钱包 - 存储余额
type Wallet struct {
	gorm.Model
	UserId  uint    `json:"userId" gorm:"column:user_id;uniqueIndex;comment:用户ID"`
	Balance float64 `json:"balance" gorm:"column:balance;type:decimal(20,6);default:0;comment:可用余额"`
	Frozen  float64 `json:"frozen" gorm:"column:frozen;type:decimal(20,6);default:0;comment:锁定余额"`
	Version int64   `json:"version" gorm:"column:version;default:0;comment:乐观锁版本号"`
}

func (Wallet) TableName() string {
	return "wallets"
}

// Invoice 发票记录
type Invoice struct {
	gorm.Model
	RequestId string  `json:"requestId" gorm:"column:request_id;size:64;uniqueIndex;comment:申请单号"`
	UserId    uint    `json:"userId" gorm:"column:user_id;index;comment:用户ID"`
	Amount    float64 `json:"amount" gorm:"column:amount;type:decimal(20,6);comment:发票金额"`
	Title     string  `json:"title" gorm:"column:title;size:200;comment:发票抬头"`
	Type      string  `json:"type" gorm:"column:type;size:50;comment:发票类型"`
	Status    string  `json:"status" gorm:"column:status;size:20;comment:状态(Processing/Sent)"`
	Code      string  `json:"code" gorm:"column:code;size:50;comment:发票代码(电子发票)"`
	Number    string  `json:"number" gorm:"column:number;size:50;comment:发票号码"`
	FileUrl   string  `json:"fileUrl" gorm:"column:file_url;size:500;comment:电子发票下载地址"`
}

func (Invoice) TableName() string {
	return "order_invoices"
}

// InvoiceTitle 发票抬头
type InvoiceTitle struct {
	gorm.Model
	UserId uint   `json:"userId" gorm:"column:user_id;index;comment:用户ID"`
	Title  string `json:"title" gorm:"column:title;size:200;comment:抬头名称"`
	TaxNo  string `json:"taxNo" gorm:"column:tax_no;size:50;comment:纳税人识别号"`
	Type   string `json:"type" gorm:"column:type;size:20;comment:类型(Personal/Enterprise)"`
}

func (InvoiceTitle) TableName() string {
	return "order_invoice_titles"
}

// InvoiceAddress 收票地址
type InvoiceAddress struct {
	gorm.Model
	UserId    uint   `json:"userId" gorm:"column:user_id;index;comment:用户ID"`
	Consignee string `json:"consignee" gorm:"column:consignee;size:50;comment:收货人"`
	Phone     string `json:"phone" gorm:"column:phone;size:20;comment:联系电话"`
	Province  string `json:"province" gorm:"column:province;size:50;comment:省"`
	City      string `json:"city" gorm:"column:city;size:50;comment:市"`
	District  string `json:"district" gorm:"column:district;size:50;comment:区"`
	Detail    string `json:"detail" gorm:"column:detail;size:200;comment:详细地址"`
}

func (InvoiceAddress) TableName() string {
	return "order_invoice_addresses"
}

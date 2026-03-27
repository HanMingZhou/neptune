package order

import (
	"context"
	"gin-vue-admin/global"
	orderModel "gin-vue-admin/model/order"
	orderReq "gin-vue-admin/model/order/request"
	orderResp "gin-vue-admin/model/order/response"
	productModel "gin-vue-admin/model/product"
	"gin-vue-admin/utils/order"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// 表名常量
const (
	tableOrders       = "orders"
	tableTransactions = "transactions"

	aliasOrder = "o"
	aliasTx    = "tx"
)

type OrderManager interface {
	CreateOrder(ctx context.Context, req *orderReq.CreateOrderRequest) (*orderModel.Order, error)
	CheckBalanceSufficient(ctx context.Context, userId uint, productId uint, chargeType int64, quantity int64) error
	DeductBalance(ctx context.Context, userId uint, amount float64, orderId uint, remark string) error

	PayOrder(ctx context.Context, orderId int64, userId int64) error
	StartOrder(ctx context.Context, resourceId int64, resourceType int) error
	StopOrder(ctx context.Context, resourceId int64, resourceType int) error
	GetOrderList(ctx context.Context, userId int64, req orderReq.GetOrderListReq) ([]orderModel.Order, int64, error)
	GetOrderOverview(ctx context.Context, userId uint) (*orderResp.OrderOverviewResp, error)
	GetUsageList(ctx context.Context, userId uint, req orderReq.GetOrderUsageListReq) ([]orderResp.UnifiedUsageItem, int64, error)
	GetTransactionList(ctx context.Context, userId uint, req orderReq.GetTransactionListReq) ([]orderResp.TransactionDetail, int64, error)
	GetInvoiceList(ctx context.Context, userId uint, req orderReq.GetInvoiceListReq) ([]orderModel.Invoice, int64, error)
	ApplyInvoice(ctx context.Context, userId uint, req orderReq.ApplyInvoiceReq) error
}

var _ OrderManager = &OrderService{}

type OrderService struct{}

// CreateOrder 创建订单
// 预付费：创建订单并立即扣费
// 按量付费：创建订单，金额为0，后续按实际使用扣费
func (s *OrderService) CreateOrder(ctx context.Context, req *orderReq.CreateOrderRequest) (*orderModel.Order, error) {
	// 1. 获取产品信息
	var product productModel.Product
	if err := global.GVA_DB.Where("id = ?", req.ProductId).First(&product).Error; err != nil {
		logx.Error("获取产品信息失败", err)
		return nil, errors.Errorf("产品不存在")
	}

	// 2. 计算价格
	unitPrice := product.GetPrice(int64(req.ChargeType))
	var amount float64
	var endTime *time.Time
	var payTime *time.Time
	now := time.Now()

	// 预付费需要计算总金额和到期时间
	if req.ChargeType != productModel.ChargeTypeHourly {
		amount, _ = decimal.NewFromFloat(unitPrice).Mul(decimal.NewFromInt(int64(req.Quantity))).Float64()
		switch req.ChargeType {
		case productModel.ChargeTypeDaily:
			t := now.AddDate(0, 0, req.Quantity)
			endTime = &t
		case productModel.ChargeTypeWeekly:
			t := now.AddDate(0, 0, req.Quantity*7)
			endTime = &t
		case productModel.ChargeTypeMonthly:
			t := now.AddDate(0, req.Quantity, 0)
			endTime = &t
		}
		payTime = &now
	} else {
		// 按量付费初始金额为0
		amount = 0
	}

	// 3. 生成订单号和记录编号
	orderNo := order.GenerateOrderNo()
	recordNo := order.GenerateOrderNo()

	// 4. 创建订单
	newOrder := &orderModel.Order{
		RecordNo:       recordNo,
		OrderNo:        orderNo,
		UserId:         req.UserId,
		ProductType:    req.ProduceType,
		ResourceType:   req.ResourceType,
		ResourceId:     req.ResourceId,
		ProductId:      req.ProductId,
		ImageId:        req.ImageId,
		ChargeType:     req.ChargeType,
		UnitPrice:      unitPrice,
		Quantity:       req.Quantity,
		Duration:       0,
		Amount:         amount,
		DiscountAmount: 0,
		StartTime:      now,
		EndTime:        endTime,
		PayTime:        payTime,
		Status:         orderModel.OrderStatusInUse,
		Area:           req.Area,
		ClusterId:      req.ClusterId,
		Remark:         req.Remark,
	}

	// 5. 事务处理
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 创建订单
		if err := tx.Create(newOrder).Error; err != nil {
			logx.Error("创建订单失败", err)
			return err
		}

		// 预付费立即扣费
		if req.ChargeType != productModel.ChargeTypeHourly && amount > 0 {
			if err := s.DeductBalance(ctx, req.UserId, amount, newOrder.ID, "预付费订单支付"); err != nil {
				logx.Error("预付费扣费失败", logx.Field("err", err))
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return newOrder, nil
}

// CheckBalanceSufficient 通用余额检查
// 预付费：检查余额 >= unitPrice * quantity
// 按量付费：检查余额 >= 单价（至少够1小时）
func (s *OrderService) CheckBalanceSufficient(ctx context.Context, userId uint, productId uint, chargeType int64, quantity int64) error {
	// 1. 获取产品信息
	var product productModel.Product
	if err := global.GVA_DB.Where("id = ?", productId).First(&product).Error; err != nil {
		return errors.Errorf("产品不存在")
	}

	// 2. 计算所需金额
	unitPrice := product.GetPrice(chargeType)
	var requiredAmount float64

	if chargeType != productModel.ChargeTypeHourly {
		// 预付费：需要全额
		if quantity <= 0 {
			quantity = 1
		}
		requiredAmount, _ = decimal.NewFromFloat(unitPrice).Mul(decimal.NewFromInt(quantity)).Float64()
	} else {
		// 按量付费：至少需要1小时的费用
		requiredAmount = unitPrice
	}

	if requiredAmount <= 0 {
		return nil
	}

	// 3. 查询钱包余额
	var wallet orderModel.Wallet
	if err := global.GVA_DB.Where("user_id = ?", userId).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Errorf("余额不足，需要 %.6f 元，当前余额 0.00 元", requiredAmount)
		}
		return errors.Wrap(err, "查询钱包失败")
	}

	currentBalance := decimal.NewFromFloat(wallet.Balance)
	required := decimal.NewFromFloat(requiredAmount)

	if currentBalance.LessThan(required) {
		return errors.Errorf("余额不足，需要 %.6f 元，当前余额 %.6f 元", requiredAmount, wallet.Balance)
	}

	return nil
}

// DeductBalance 扣费
func (s *OrderService) DeductBalance(ctx context.Context, userId uint, amount float64, orderId uint, remark string) error {
	if amount <= 0 {
		return nil
	}

	deductAmount := decimal.NewFromFloat(amount)

	// 最多重试3次（乐观锁冲突时）
	for range 3 {
		// 1. 读取钱包（不加锁）
		var wallet orderModel.Wallet
		if err := global.GVA_DB.Where("user_id = ?", userId).First(&wallet).Error; err != nil {
			return err
		}

		// 2. 检查余额
		currentBalance := decimal.NewFromFloat(wallet.Balance)
		if currentBalance.LessThan(deductAmount) {
			return errors.Errorf("余额不足，需支付 %.4f，当前余额 %.4f", amount, wallet.Balance)
		}

		// 3. 乐观锁更新：WHERE version = oldVersion
		newBalance := currentBalance.Sub(deductAmount)
		newBalanceFloat, _ := newBalance.Round(6).Float64()
		result := global.GVA_DB.Model(&orderModel.Wallet{}).
			Where("id = ? AND version = ?", wallet.ID, wallet.Version).
			Updates(map[string]any{
				"balance": newBalanceFloat,
				"version": gorm.Expr("version + 1"),
			})

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			// 版本冲突，重试
			logx.Info("DeductBalance 乐观锁冲突，重试", logx.Field("userId", userId))
			continue
		}

		// 4. 记录交易流水
		balanceBefore, _ := currentBalance.Round(6).Float64()
		balanceAfter, _ := newBalance.Round(6).Float64()
		amountRounded, _ := deductAmount.Round(6).Float64()
		now := time.Now()
		transaction := orderModel.Transaction{
			TransactionNo: order.GenerateOrderNo(),
			UserId:        userId,
			Type:          orderModel.TransactionTypeConsumption,
			Amount:        -amountRounded,
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
			OrderId:       orderId,
			Method:        orderModel.PayMethodBalance,
			Remark:        remark,
			CreatedAt:     now,
		}
		if err := global.GVA_DB.Create(&transaction).Error; err != nil {
			logx.Error("记录交易流水失败", logx.Field("err", err))
		}
		return nil
	}
	return errors.Errorf("扣费失败：乐观锁冲突次数过多")
}

// PayOrder 支付订单（已废弃，现在预付费在CreateOrder时直接扣费）
func (s *OrderService) PayOrder(ctx context.Context, orderId int64, userId int64) error {
	return errors.New("PayOrder已废弃，预付费在CreateOrder时直接扣费")
}

// StartOrder 开始计费（PodGroup进入Running状态时由Informer调用）
// 预付费：无需操作（已在CreateOrder时扣费）
// 按量付费：确保订单状态正确，准备开始累计计费
func (s *OrderService) StartOrder(ctx context.Context, resourceId int64, resourceType int) error {
	// 查找该资源的订单
	var orderInfo orderModel.Order
	err := global.GVA_DB.Where("resource_id = ? AND resource_type = ? AND status = ?",
		resourceId, resourceType, orderModel.OrderStatusInUse).
		First(&orderInfo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没有找到使用中的订单，可能是预付费已完成或订单不存在
			logx.Info("StartOrder: 未找到使用中的订单",
				logx.Field("resourceId", resourceId),
				logx.Field("resourceType", resourceType))
			return nil
		}
		return err
	}

	// 预付费订单：已在CreateOrder时扣费，无需额外操作
	if orderInfo.ChargeType != productModel.ChargeTypeHourly {
		logx.Info("StartOrder: 预付费订单无需操作",
			logx.Field("orderId", orderInfo.ID),
			logx.Field("chargeType", orderInfo.ChargeType))
		return nil
	}

	// 按量付费订单：确保StartTime已设置（通常在CreateOrder时已设置）
	// 这里只是防御性检查，正常情况下StartTime应该已经存在
	if orderInfo.StartTime.IsZero() {
		now := time.Now()
		if err := global.GVA_DB.Model(&orderModel.Order{}).
			Where("id = ?", orderInfo.ID).
			Update("start_time", now).Error; err != nil {
			logx.Error("StartOrder: 更新开始时间失败",
				logx.Field("orderId", orderInfo.ID),
				logx.Field("err", err))
			return err
		}
		logx.Info("StartOrder: 已更新按量计费订单开始时间",
			logx.Field("orderId", orderInfo.ID))
	}

	logx.Info("StartOrder: 按量计费订单已就绪",
		logx.Field("orderId", orderInfo.ID),
		logx.Field("resourceId", resourceId))
	return nil
}

// StopOrder 停止计费（资源释放时调用）
// 按量付费：计算实际使用时长和费用，扣费并更新订单
// 预付费：只更新状态
func (s *OrderService) StopOrder(ctx context.Context, resourceId int64, resourceType int) error {
	now := time.Now()

	// 查找使用中的订单
	var orderInfo orderModel.Order
	err := global.GVA_DB.Where("resource_id = ? AND resource_type = ? AND status = ?",
		resourceId, resourceType, orderModel.OrderStatusInUse).
		First(&orderInfo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logx.Info("StopOrder: 未找到使用中的订单",
				logx.Field("resourceId", resourceId),
				logx.Field("resourceType", resourceType))
			return nil // 没有使用中的订单
		}
		return err
	}

	logx.Info("StopOrder: 找到订单，准备结算",
		logx.Field("orderId", orderInfo.ID),
		logx.Field("resourceId", resourceId),
		logx.Field("chargeType", orderInfo.ChargeType))

	// 按量付费：计算费用并扣费
	if orderInfo.ChargeType == productModel.ChargeTypeHourly {
		// 计算使用时长（秒）
		duration := int64(now.Sub(orderInfo.StartTime).Seconds())

		// 计算费用（按小时计费），全程 decimal 避免精度丢失
		hours := decimal.NewFromInt(duration).Div(decimal.NewFromInt(3600))
		amountDec := hours.Mul(decimal.NewFromFloat(orderInfo.UnitPrice)).Round(6)
		amount, _ := amountDec.Float64()

		// 扣费
		if amount > 0 {
			if err := s.DeductBalance(ctx, orderInfo.UserId, amount, orderInfo.ID, "按量计费结算"); err != nil {
				logx.Error("StopOrder 扣费失败", logx.Field("err", err), logx.Field("resourceId", resourceId))
				// 扣费失败也要更新状态，避免重复扣费
			}
		}

		// 更新订单：累加时长和金额，状态改为已结算
		result := global.GVA_DB.Model(&orderModel.Order{}).
			Where("id = ?", orderInfo.ID).
			Updates(map[string]any{
				"duration":    duration,
				"amount":      amount,
				"status":      orderModel.OrderStatusFinished,
				"settle_time": now,
			})
		if result.Error != nil {
			logx.Error("StopOrder: 更新订单状态失败",
				logx.Field("orderId", orderInfo.ID),
				logx.Field("err", result.Error))
		} else {
			logx.Info("StopOrder: 订单已结算",
				logx.Field("orderId", orderInfo.ID),
				logx.Field("rowsAffected", result.RowsAffected),
				logx.Field("amount", amount),
				logx.Field("duration", duration))
		}
		return result.Error
	}

	// 预付费：计算已使用时长，退还剩余金额，状态改为已结算
	duration := int64(now.Sub(orderInfo.StartTime).Seconds())

	// 计算退款金额（仅当有到期时间且未过期时退款）
	var refundAmount float64
	if orderInfo.EndTime != nil && now.Before(*orderInfo.EndTime) && orderInfo.Amount > 0 {
		totalDuration := orderInfo.EndTime.Sub(orderInfo.StartTime).Seconds()
		if totalDuration > 0 {
			usedRatio := float64(duration) / totalDuration
			if usedRatio > 1 {
				usedRatio = 1
			}
			refundDec := decimal.NewFromFloat(orderInfo.Amount).
				Mul(decimal.NewFromFloat(1 - usedRatio)).
				Round(6)
			refundAmount, _ = refundDec.Float64()
		}
	}

	// 退款到钱包
	if refundAmount > 0 {
		refundDec := decimal.NewFromFloat(refundAmount)
		// 读取钱包
		var wallet orderModel.Wallet
		if err := global.GVA_DB.Where("user_id = ?", orderInfo.UserId).First(&wallet).Error; err == nil {
			newBalance := decimal.NewFromFloat(wallet.Balance).Add(refundDec)
			newBalanceFloat, _ := newBalance.Round(6).Float64()
			global.GVA_DB.Model(&orderModel.Wallet{}).
				Where("id = ?", wallet.ID).
				Updates(map[string]any{
					"balance": newBalanceFloat,
					"version": gorm.Expr("version + 1"),
				})

			// 记录退款流水
			refundRounded, _ := refundDec.Round(6).Float64()
			transaction := orderModel.Transaction{
				TransactionNo: order.GenerateOrderNo(),
				UserId:        orderInfo.UserId,
				Type:          orderModel.TransactionTypeRefund,
				Amount:        refundRounded,
				BalanceBefore: wallet.Balance,
				BalanceAfter:  newBalanceFloat,
				OrderId:       orderInfo.ID,
				Method:        orderModel.PayMethodBalance,
				Remark:        "预付费订单提前停止退款",
				CreatedAt:     now,
			}
			if err := global.GVA_DB.Create(&transaction).Error; err != nil {
				logx.Error("记录退款流水失败", logx.Field("err", err))
			}
		}
	}

	// 更新订单状态为已结算
	actualAmount, _ := decimal.NewFromFloat(orderInfo.Amount).Sub(decimal.NewFromFloat(refundAmount)).Round(6).Float64()
	return global.GVA_DB.Model(&orderModel.Order{}).
		Where("id = ?", orderInfo.ID).
		Updates(map[string]any{
			"duration":    duration,
			"amount":      actualAmount,
			"status":      orderModel.OrderStatusFinished,
			"settle_time": now,
		}).Error
}

// GetOrderList 获取用户订单列表
func (s *OrderService) GetOrderList(ctx context.Context, userId int64, req orderReq.GetOrderListReq) ([]orderModel.Order, int64, error) {
	var orders []orderModel.Order
	var total int64

	db := global.GVA_DB.Model(&orderModel.Order{}).Where("user_id = ?", userId)

	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := db.Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetOrderOverview 获取用户账单概览
func (s *OrderService) GetOrderOverview(ctx context.Context, userId uint) (*orderResp.OrderOverviewResp, error) {
	var resp orderResp.OrderOverviewResp

	// 1. 获取余额
	var wallet orderModel.Wallet
	err := global.GVA_DB.Where("user_id = ?", userId).FirstOrCreate(&wallet, orderModel.Wallet{UserId: userId}).Error
	if err != nil {
		return nil, err
	}
	resp.Balance = wallet.Balance
	resp.Status = "Healthy"

	// 2. 本月消费（按订单开始时间统计）
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var mtdAmount float64
	err = global.GVA_DB.Table(tableOrders).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND start_time >= ?", userId, firstDayOfMonth).
		Scan(&mtdAmount).Error
	if err != nil {
		return nil, err
	}
	resp.MtdEstimate = mtdAmount

	// 3. 上月结算
	lastMonth := now.AddDate(0, -1, 0)
	lastMonthStr := lastMonth.Format("2006-01")
	var lastMonthSummary orderModel.OrderSummary
	global.GVA_DB.Where("user_id = ? AND order_month = ?", userId, lastMonthStr).Limit(1).Find(&lastMonthSummary)
	resp.LastMonthSettlement = lastMonthSummary.TotalAmount

	// 4. 消费趋势（最近7天）
	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i).Format("01-02")
		dateStart := time.Date(now.Year(), now.Month(), now.Day()-i, 0, 0, 0, 0, now.Location())
		dateEnd := dateStart.Add(24 * time.Hour)

		var dayAmount float64
		global.GVA_DB.Table(tableOrders).
			Select("COALESCE(SUM(amount), 0)").
			Where("user_id = ? AND start_time >= ? AND start_time < ?", userId, dateStart, dateEnd).
			Scan(&dayAmount)

		resp.ConsumptionTrend = append(resp.ConsumptionTrend, orderResp.ConsumptionData{
			Date:   date,
			Amount: decimal.NewFromFloat(dayAmount).StringFixed(2),
		})
	}

	// 5. 费用排行（本月，按资源类型）
	type RankingRaw struct {
		ResourceType int64
		Amount       float64
	}
	var rankingsRaw []RankingRaw

	err = global.GVA_DB.Table(tableOrders).
		Select("resource_type, SUM(amount) as amount").
		Where("user_id = ? AND start_time >= ?", userId, firstDayOfMonth).
		Group("resource_type").
		Order("amount DESC").
		Limit(4).
		Scan(&rankingsRaw).Error
	if err != nil {
		return nil, err
	}

	resourceTypeNames := map[int64]string{
		int64(orderModel.OrderTypeNotebook):  "容器实例",
		int64(orderModel.OrderTypeTraining):  "训练任务",
		int64(orderModel.OrderTypeInference): "推理服务",
		int64(orderModel.OrderTypeStorage):   "数据存储",
	}

	for _, r := range rankingsRaw {
		name := resourceTypeNames[r.ResourceType]
		if name == "" {
			name = "其他"
		}
		percent := 0.0
		if resp.MtdEstimate > 0 {
			percent = (r.Amount / resp.MtdEstimate) * 100
		}
		resp.MonthlyRanking = append(resp.MonthlyRanking, orderResp.RankingData{
			Name:    name,
			Amount:  decimal.NewFromFloat(r.Amount).StringFixed(2),
			Percent: percent,
		})
	}

	return &resp, nil
}

// GetUsageList 获取使用列表（直接从Order表查询）
func (s *OrderService) GetUsageList(ctx context.Context, userId uint, req orderReq.GetOrderUsageListReq) ([]orderResp.UnifiedUsageItem, int64, error) {
	var result []orderResp.UnifiedUsageItem

	// 查询订单
	var orders []orderModel.Order
	orderDb := global.GVA_DB.Model(&orderModel.Order{}).Where("user_id = ?", userId)

	if req.StartTime != "" {
		orderDb = orderDb.Where("start_time >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		orderDb = orderDb.Where("start_time <= ?", req.EndTime)
	}

	if err := orderDb.Order("id DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	// 按资源类型分组收集 ResourceId，批量查询名称
	resourceNames := s.batchQueryResourceNames(orders)

	// 转换为统一格式
	now := time.Now()
	for _, o := range orders {
		// 判断状态（预付费可能已过期）
		status := o.Status
		if o.ChargeType != productModel.ChargeTypeHourly && o.EndTime != nil && now.After(*o.EndTime) && status == orderModel.OrderStatusInUse {
			status = orderModel.OrderStatusExpired
		}

		// 格式化时间
		startStr := o.StartTime.Format("2006-01-02 15:04:05")
		endStr := ""
		if o.EndTime != nil {
			endStr = o.EndTime.Format("2006-01-02 15:04:05")
		}

		// 确定数据来源
		source := "预付费"
		if o.ChargeType == productModel.ChargeTypeHourly {
			source = "按量计费"
		}

		// 查找资源名称
		nameKey := resourceNameKey{ResourceType: o.ResourceType, ResourceId: o.ResourceId}
		resourceName := resourceNames[nameKey]

		result = append(result, orderResp.UnifiedUsageItem{
			ID:             o.ID,
			Source:         source,
			ResourceType:   o.ResourceType,
			ResourceName:   resourceName,
			ChargeType:     o.ChargeType,
			ChargeTypeName: productModel.ChargeTypeToString[int(o.ChargeType)],
			UnitPrice:      o.UnitPrice,
			Duration:       o.Duration,
			Quantity:       o.Quantity,
			Amount:         o.Amount,
			Status:         status,
			StatusText:     orderModel.OrderStatusToString[status],
			StartTime:      startStr,
			EndTime:        endStr,
			CreatedAt:      o.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:      o.UpdatedAt.Format("2006-01-02 15:04:05"),
			OrderNo:        o.OrderNo,
		})
	}

	// 手动分页
	total := int64(len(result))
	offset := (req.Page - 1) * req.PageSize
	end := offset + req.PageSize
	if offset >= int(total) {
		return []orderResp.UnifiedUsageItem{}, total, nil
	}
	if end > int(total) {
		end = int(total)
	}

	return result[offset:end], total, nil
}

// resourceNameKey 资源名称查询键
type resourceNameKey struct {
	ResourceType int64
	ResourceId   uint
}

// batchQueryResourceNames 批量查询资源名称，避免 N+1 查询
func (s *OrderService) batchQueryResourceNames(orders []orderModel.Order) map[resourceNameKey]string {
	result := make(map[resourceNameKey]string)

	// 按资源类型分组收集 ID
	notebookIds := make([]uint, 0)
	trainingIds := make([]uint, 0)
	inferenceIds := make([]uint, 0)
	storageIds := make([]uint, 0)

	for _, o := range orders {
		if o.ResourceId == 0 {
			continue
		}
		switch o.ResourceType {
		case orderModel.OrderTypeNotebook:
			notebookIds = append(notebookIds, o.ResourceId)
		case orderModel.OrderTypeTraining:
			trainingIds = append(trainingIds, o.ResourceId)
		case orderModel.OrderTypeInference:
			inferenceIds = append(inferenceIds, o.ResourceId)
		case orderModel.OrderTypeStorage:
			storageIds = append(storageIds, o.ResourceId)
		}
	}

	// 批量查询 Notebook 名称
	if len(notebookIds) > 0 {
		var rows []struct {
			ID   uint
			Name string
		}
		global.GVA_DB.Table("notebooks").Select("id, display_name as name").Where("id IN ?", notebookIds).Scan(&rows)
		for _, r := range rows {
			result[resourceNameKey{ResourceType: orderModel.OrderTypeNotebook, ResourceId: r.ID}] = r.Name
		}
	}

	// 批量查询 TrainingJob 名称
	if len(trainingIds) > 0 {
		var rows []struct {
			ID          uint
			DisplayName string
		}
		global.GVA_DB.Table("training_jobs").Select("id, display_name").Where("id IN ?", trainingIds).Scan(&rows)
		for _, r := range rows {
			result[resourceNameKey{ResourceType: orderModel.OrderTypeTraining, ResourceId: r.ID}] = r.DisplayName
		}
	}

	// 批量查询 Inference 名称
	if len(inferenceIds) > 0 {
		var rows []struct {
			ID          uint
			DisplayName string
		}
		global.GVA_DB.Table("inferences").Select("id, display_name").Where("id IN ?", inferenceIds).Scan(&rows)
		for _, r := range rows {
			result[resourceNameKey{ResourceType: orderModel.OrderTypeInference, ResourceId: r.ID}] = r.DisplayName
		}
	}

	// 批量查询 Volume 名称
	if len(storageIds) > 0 {
		var rows []struct {
			ID   uint
			Name string
		}
		global.GVA_DB.Table("volumes").Select("id, name").Where("id IN ?", storageIds).Scan(&rows)
		for _, r := range rows {
			result[resourceNameKey{ResourceType: orderModel.OrderTypeStorage, ResourceId: r.ID}] = r.Name
		}
	}

	return result
}

// GetTransactionList 获取交易流水列表
func (s *OrderService) GetTransactionList(ctx context.Context, userId uint, req orderReq.GetTransactionListReq) ([]orderResp.TransactionDetail, int64, error) {
	var list []orderResp.TransactionDetail
	var total int64

	db := global.GVA_DB.Table(tableTransactions+" "+aliasTx).
		Select(aliasTx+".*, "+
			aliasOrder+".status as order_status, "+
			aliasOrder+".order_no as order_no, "+
			aliasOrder+".resource_type as resource_type, "+
			aliasOrder+".resource_id as resource_id, "+
			aliasOrder+".charge_type as charge_type").
		Joins("LEFT JOIN "+tableOrders+" "+aliasOrder+" ON "+aliasTx+".order_id = "+aliasOrder+".id").
		Where(aliasTx+".user_id = ?", userId)

	if req.Type > 0 {
		db = db.Where(aliasTx+".type = ?", req.Type)
	}

	if req.StartTime != "" {
		db = db.Where(aliasTx+".created_at >= ?", req.StartTime)
	}

	if req.EndTime != "" {
		db = db.Where(aliasTx+".created_at <= ?", req.EndTime)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 临时结构体用于接收 resource_id（不暴露给前端）
	type txWithResourceId struct {
		orderResp.TransactionDetail
		ResourceId uint `json:"resourceId"`
	}
	var rawList []txWithResourceId

	err = db.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).Order(aliasTx + ".id DESC").Scan(&rawList).Error
	if err != nil {
		return nil, 0, err
	}

	// 构造虚拟 Order 列表用于批量查询资源名称
	var fakeOrders []orderModel.Order
	for _, item := range rawList {
		if item.ResourceType > 0 && item.ResourceId > 0 {
			fakeOrders = append(fakeOrders, orderModel.Order{
				ResourceType: item.ResourceType,
				ResourceId:   item.ResourceId,
			})
		}
	}
	resourceNames := s.batchQueryResourceNames(fakeOrders)

	// 填充资源名称
	list = make([]orderResp.TransactionDetail, len(rawList))
	for i, item := range rawList {
		list[i] = item.TransactionDetail
		if item.ResourceType > 0 && item.ResourceId > 0 {
			nameKey := resourceNameKey{ResourceType: item.ResourceType, ResourceId: item.ResourceId}
			list[i].ResourceName = resourceNames[nameKey]
		}
	}

	return list, total, nil
}

// GetInvoiceList 获取发票列表
func (s *OrderService) GetInvoiceList(ctx context.Context, userId uint, req orderReq.GetInvoiceListReq) ([]orderModel.Invoice, int64, error) {
	var list []orderModel.Invoice
	var total int64
	db := global.GVA_DB.Model(&orderModel.Invoice{}).Where("user_id = ?", userId)
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).Order("id desc").Find(&list).Error
	return list, total, err
}

// ApplyInvoice 申请发票
func (s *OrderService) ApplyInvoice(ctx context.Context, userId uint, req orderReq.ApplyInvoiceReq) error {
	var title orderModel.InvoiceTitle
	if err := global.GVA_DB.Where("id = ? AND user_id = ?", req.TitleId, userId).First(&title).Error; err != nil {
		return errors.New("发票抬头不存在")
	}

	invoice := orderModel.Invoice{
		RequestId: order.GenerateOrderNo(),
		UserId:    userId,
		Amount:    req.Amount,
		Title:     title.Title,
		Type:      title.Type,
		Status:    orderModel.InvoiceStatusProcessing,
	}
	return global.GVA_DB.Create(&invoice).Error
}

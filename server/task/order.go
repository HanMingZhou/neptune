package task

import (
	"context"
	"gin-vue-admin/global"
	"gin-vue-admin/model/consts"
	orderModel "gin-vue-admin/model/order"
	podgroupModel "gin-vue-admin/model/podgroup"
	productModel "gin-vue-admin/model/product"
	pvcModel "gin-vue-admin/model/pvc"
	pvcReq "gin-vue-admin/model/pvc/request"
	orderService "gin-vue-admin/service/order"
	"gin-vue-admin/utils/order"
	"time"

	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
)

var orderSvc = &orderService.OrderService{}

// HourlyOrderTask 按量计费定时扣费任务（每小时执行）
// 扫描 podgroups 表中 Running 的资源，关联 orders 表进行扣费
func HourlyOrderTask() error {
	logx.Info("开始执行按量计费定时扣费任务")
	ctx := context.Background()

	// 查询所有 Running 状态且有开始时间的 PodGroup
	var podgroups []podgroupModel.PodGroup
	if err := global.GVA_DB.Where("status = ? AND start_time IS NOT NULL", consts.PodGroupStatusRunning).Find(&podgroups).Error; err != nil {
		logx.Error("查询运行中的 PodGroup 失败", logx.Field("err", err))
		return err
	}

	if len(podgroups) == 0 {
		logx.Info("没有运行中的按量计费资源")
		return nil
	}

	// 一次性批量查出所有使用中的按量计费订单（避免 N+1 查询）
	type activeOrder struct {
		orderModel.Order
	}
	var allOrders []activeOrder
	if err := global.GVA_DB.Model(&orderModel.Order{}).
		Where("status = ? AND charge_type = ?", orderModel.OrderStatusInUse, productModel.ChargeTypeHourly).
		Find(&allOrders).Error; err != nil {
		logx.Error("批量查询活跃按量计费订单失败", logx.Field("err", err))
		return err
	}

	// 构建 map：key = "resourceId:resourceType" → order
	type orderKey struct {
		resourceId   uint
		resourceType int64
	}
	orderMap := make(map[orderKey]*activeOrder, len(allOrders))
	for i := range allOrders {
		key := orderKey{resourceId: allOrders[i].ResourceId, resourceType: allOrders[i].ResourceType}
		orderMap[key] = &allOrders[i]
	}

	now := time.Now()
	successCount := 0
	failCount := 0

	for _, pg := range podgroups {
		if pg.OwnerID == 0 {
			continue
		}

		// O(1) map 查找，替代原来每次循环的 JOIN 查询
		resourceType := getResourceType(pg.InstanceType)
		orderInfo, ok := orderMap[orderKey{resourceId: pg.OwnerID, resourceType: resourceType}]
		if !ok {
			continue // 没有关联的订单，跳过
		}

		// 计算本次应扣金额：从上次结算到现在的时长
		lastSettleTime := orderInfo.StartTime
		if orderInfo.SettleTime != nil {
			lastSettleTime = *orderInfo.SettleTime
		}

		elapsed := now.Sub(lastSettleTime)
		if elapsed.Seconds() < 60 {
			continue // 不足1分钟不扣费
		}

		// 精确到秒计算本次金额
		seconds := decimal.NewFromFloat(elapsed.Seconds())
		hourlyPrice := decimal.NewFromFloat(orderInfo.UnitPrice)
		incrementAmount, _ := seconds.Div(decimal.NewFromInt(3600)).Mul(hourlyPrice).Float64()

		if incrementAmount <= 0 {
			continue
		}

		// 调用通用扣费
		if err := orderSvc.DeductBalance(ctx, orderInfo.UserId, incrementAmount, orderInfo.ID, "按量计费（定时扣费）"); err != nil {
			logx.Error("按量计费扣费失败",
				logx.Field("userId", orderInfo.UserId),
				logx.Field("resourceId", orderInfo.ResourceId),
				logx.Field("amount", incrementAmount),
				logx.Field("err", err))
			failCount++
			continue
		}

		// 更新 Order 累计时长、累计金额和结算时间
		newDuration := int64(now.Sub(orderInfo.StartTime).Seconds())
		newAmount := orderInfo.Amount + incrementAmount

		if err := global.GVA_DB.Model(&orderModel.Order{}).
			Where("id = ?", orderInfo.ID).
			Updates(map[string]interface{}{
				"duration":    newDuration,
				"amount":      newAmount,
				"settle_time": now,
			}).Error; err != nil {
			logx.Error("更新订单失败", logx.Field("orderId", orderInfo.ID), logx.Field("err", err))
			failCount++
			continue
		}

		successCount++
	}

	logx.Info("按量计费定时扣费任务完成",
		logx.Field("total", len(podgroups)),
		logx.Field("success", successCount),
		logx.Field("fail", failCount))
	return nil
}

// DailyStorageOrderTask 存储按天计费任务（每天执行）
// 扫描 volumes 表中活跃的数据盘，按大小计费
func DailyStorageOrderTask() error {
	logx.Info("开始执行存储按天计费任务")
	ctx := context.Background()

	// 获取存储单价配置（从系统配置或使用默认值）
	// TODO: 后续从配置中心获取动态单价
	var storagePricePerGBPerDay float64 = 0.01 // 默认 0.01 元/GB/天

	// 查询所有未删除的数据盘
	var volumes []pvcReq.TaskVolumeInfo
	if err := global.GVA_DB.Model(&pvcModel.Volume{}).
		Select("id, name, size, user_id, cluster_id").
		Where("deleted_at IS NULL").
		Find(&volumes).Error; err != nil {
		logx.Error("查询数据盘失败", logx.Field("err", err))
		return err
	}

	if len(volumes) == 0 {
		logx.Info("没有活跃的存储资源")
		return nil
	}

	successCount := 0
	failCount := 0
	now := time.Now()
	today := now.Format("2006-01-02")

	for _, vol := range volumes {
		if vol.Size <= 0 || vol.UserId == 0 {
			continue
		}

		// 按照 GB 计费
		sizeGB := float64(vol.Size)
		amount := sizeGB * storagePricePerGBPerDay

		if amount <= 0 {
			continue
		}

		// 检查今天是否已有该资源的订单（防止重复扣费）
		var existCount int64
		global.GVA_DB.Model(&orderModel.Order{}).
			Where("resource_id = ? AND resource_type = ? AND DATE(start_time) = ?",
				vol.ID, orderModel.OrderTypeStorage, today).
			Count(&existCount)

		if existCount > 0 {
			continue // 今天已扣费
		}

		// 生成订单号和记录编号
		orderNo := order.GenerateOrderNo()
		recordNo := order.GenerateOrderNo()

		// 创建存储订单（存储按天计费，视为按量计费模式）
		storageOrder := &orderModel.Order{
			RecordNo:       recordNo,
			OrderNo:        orderNo,
			UserId:         vol.UserId,
			ProductType:    orderModel.ProductTypeStorage,
			ResourceType:   orderModel.OrderTypeStorage,
			ResourceId:     vol.ID,
			ChargeType:     productModel.ChargeTypeDaily, // 存储按天计费
			UnitPrice:      storagePricePerGBPerDay,
			Quantity:       int(vol.Size), // 数量为存储大小（GB）
			Duration:       86400,         // 1天的秒数
			Amount:         amount,
			DiscountAmount: 0,
			StartTime:      now,
			EndTime:        nil, // 存储没有固定结束时间
			PayTime:        nil, // 按量计费没有PayTime
			SettleTime:     nil, // 扣费成功后会设置
			Status:         orderModel.OrderStatusInUse,
			ClusterId:      vol.ClusterId,
			Remark:         "存储按天计费",
		}

		if err := global.GVA_DB.Create(storageOrder).Error; err != nil {
			logx.Error("创建存储订单失败", logx.Field("err", err))
			failCount++
			continue
		}

		// 扣费
		if err := orderSvc.DeductBalance(ctx, vol.UserId, amount, storageOrder.ID, "存储按天计费"); err != nil {
			logx.Error("存储扣费失败",
				logx.Field("userId", vol.UserId),
				logx.Field("volumeId", vol.ID),
				logx.Field("amount", amount),
				logx.Field("err", err))
			// 扣费失败，更新订单状态为已停止（保留失败订单）
			global.GVA_DB.Model(storageOrder).Update("status", orderModel.OrderStatusStopped)
			failCount++
			continue
		}

		// 扣费成功，更新结算时间
		global.GVA_DB.Model(storageOrder).Update("settle_time", now)
		successCount++
	}

	logx.Info("存储按天计费任务完成",
		logx.Field("total", len(volumes)),
		logx.Field("success", successCount),
		logx.Field("fail", failCount))
	return nil
}

// getResourceType 根据 PodGroup 实例类型获取计费资源类型常量
func getResourceType(instanceType string) int64 {
	switch instanceType {
	case consts.NotebookInstance:
		return orderModel.OrderTypeNotebook
	case consts.TrainingInstance:
		return orderModel.OrderTypeTraining
	case consts.InferenceInstance:
		return orderModel.OrderTypeInference
	default:
		return 0
	}
}

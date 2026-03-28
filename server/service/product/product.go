package product

import (
	"context"
	"errors"
	"gin-vue-admin/global"
	productModel "gin-vue-admin/model/product"
	productRes "gin-vue-admin/model/product/response"
	redisUtils "gin-vue-admin/utils/redis"

	"gorm.io/gorm"
)

type ProductService struct{}

// toProductDetailResponse 将 Product Model 统一转换为 ProductDetailResponse DTO
// 集中维护映射关系，避免多处重复代码导致字段遗漏
func toProductDetailResponse(p *productModel.Product) productRes.ProductDetailResponse {
	return productRes.ProductDetailResponse{
		ID:                p.ID,
		TemplateProductId: p.ID,
		ProductType:       p.ProductType,
		Name:              p.Name,
		Description:       p.Description,
		Area:              p.Area,
		NodeName:          p.NodeName,
		NodeType:          p.NodeType,
		CPUModel:          p.CPUModel,
		CPU:               p.CPU,
		Memory:            p.Memory,
		GPUModel:          p.GPUModel,
		GPUCount:          p.GPUCount,
		GPUMemory:         p.GPUMemory,
		VGPUNumber:        p.VGPUNumber,
		VGPUMemory:        p.VGPUMemory,
		VGPUCores:         p.VGPUCores,
		PriceHourly:       p.PriceHourly,
		PriceDaily:        p.PriceDaily,
		PriceWeekly:       p.PriceWeekly,
		PriceMonthly:      p.PriceMonthly,
		DriverVersion:     p.DriverVersion,
		CUDAVersion:       p.CUDAVersion,
		SystemDisk:        p.SystemDisk,
		DataDisk:          p.DataDisk,
		Status:            p.Status,
		StorageClass:      p.StorageClass,
		StoragePriceGB:    p.StoragePriceGB,
		MaxInstances:      p.MaxInstances,
		Available:         p.AvailableCapacity(),
		ClusterId:         p.ClusterId,
	}
}

// GetProductById 根据产品ID获取产品信息
func (s *ProductService) GetProductById(ctx context.Context, productId uint) (*productRes.ProductDetailResponse, error) {
	var product productModel.Product
	if err := global.GVA_DB.WithContext(ctx).Where("id = ? AND status = ?", productId, productModel.ProductStatusEnabled).First(&product).Error; err != nil {
		return nil, err
	}

	resp := toProductDetailResponse(&product)
	return &resp, nil
}

// GetProductList 获取产品列表
func (s *ProductService) GetProductList(
	ctx context.Context,
	page, pageSize, productType int,
	filters AggregateProductFilters,
) (*productRes.ProductListResponse, error) {
	// 分页参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var products []productModel.Product
	var total int64
	db := s.buildProductListDB(ctx, productType, filters)

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("sort_order ASC, id ASC").Find(&products).Error; err != nil {
		return nil, err
	}

	list := make([]productRes.ProductDetailResponse, 0, len(products))
	for _, p := range products {
		list = append(list, toProductDetailResponse(&p))
	}

	resp := &productRes.ProductListResponse{
		List:  list,
		Total: total,
	}

	return resp, nil
}

// GetProductFilters 获取产品筛选条件
func (s *ProductService) GetProductFilters(ctx context.Context, productType int) (productRes.ProductFiltersResponse, error) {
	var resp productRes.ProductFiltersResponse

	baseDb := global.GVA_DB.WithContext(ctx).Model(&productModel.Product{}).Where("status = ?", productModel.ProductStatusEnabled)
	if productType > 0 {
		baseDb = baseDb.Where("product_type = ?", productType)
	}

	// 获取地区列表
	resp.Areas = make([]string, 0)
	if err := baseDb.Session(&gorm.Session{}).
		Distinct("area").
		Pluck("area", &resp.Areas).Error; err != nil {
		return resp, err
	}

	// 获取GPU型号列表（带可用数量统计）
	resp.GPUModels = make([]productRes.GPUModelInfo, 0)
	if err := baseDb.Session(&gorm.Session{}).
		Select("gpu_model as model, SUM(max_instances - used_capacity) as available, SUM(max_instances) as total").
		Where("gpu_model != '' AND gpu_model IS NOT NULL").
		Group("gpu_model").
		Scan(&resp.GPUModels).Error; err != nil {
		return resp, err
	}

	// 获取vGPU型号列表
	resp.VGPUModels = make([]productRes.VGPUModelInfo, 0)
	if err := baseDb.Session(&gorm.Session{}).
		Select("v_gpu_memory as memory, SUM(max_instances - used_capacity) as available, SUM(max_instances) as total").
		Where("v_gpu_number > 0 OR v_gpu_memory > 0 OR v_gpu_cores > 0").
		Group("v_gpu_memory").
		Scan(&resp.VGPUModels).Error; err != nil {
		return resp, err
	}

	// 获取CPU型号列表
	resp.CPUModels = make([]productRes.CPUModelInfo, 0)
	if err := baseDb.Session(&gorm.Session{}).
		Select("cpu_model as model, SUM(max_instances - used_capacity) as available, SUM(max_instances) as total").
		Where("cpu_model != '' AND cpu_model IS NOT NULL AND (gpu_model = '' OR gpu_model IS NULL)").
		Group("cpu_model").
		Scan(&resp.CPUModels).Error; err != nil {
		return resp, err
	}

	// CPU产品总计（无GPU的产品）
	if err := baseDb.Session(&gorm.Session{}).
		Select("SUM(max_instances - used_capacity) as available, SUM(max_instances) as total").
		Where("(gpu_model = '' OR gpu_model IS NULL)").
		Scan(&resp.CPUOnly).Error; err != nil {
		return resp, err
	}

	return resp, nil
}

// ReserveCapacity 检查并锁定产品资源（锁定 1 个实例配额）
// 所有产品类型统一使用 MaxInstances 管理库存
func (s *ProductService) ReserveCapacity(ctx context.Context, productId uint) (*productRes.ReserveResult, error) {
	return s.ReserveCapacityWithCount(ctx, productId, 1)
}

// ReserveCapacityWithCount 检查并锁定指定数量的产品资源
// 使用 Redis 分布式锁 + DB 乐观锁双重保护，防止并发超卖
func (s *ProductService) ReserveCapacityWithCount(ctx context.Context, productId uint, count int64) (*productRes.ReserveResult, error) {
	if count <= 0 {
		return nil, errors.New("锁定数量必须大于0")
	}

	// 1. 查询产品
	var p productModel.Product
	if err := global.GVA_DB.WithContext(ctx).Where("id = ? AND status = ?", productId, productModel.ProductStatusEnabled).First(&p).Error; err != nil {
		return nil, errors.New("产品不存在或已下架")
	}

	// 2. 快速检查（无锁，避免不必要的加锁开销）
	if p.AvailableCapacity() < count {
		return nil, ErrCapacityInsufficient
	}

	// 3. 获取分布式锁
	lock, err := redisUtils.AcquireProductLock(ctx, productId)
	if err != nil {
		return nil, err
	}
	defer lock.Unlock(ctx)

	// 4. 重新查询 + 乐观锁更新（原子操作）
	// 二次查询同样需要带 status 过滤，防止加锁期间产品被下架仍能锁定资源
	if err = global.GVA_DB.WithContext(ctx).Where("id = ? AND status = ?", productId, productModel.ProductStatusEnabled).First(&p).Error; err != nil {
		return nil, errors.New("产品不存在或已下架")
	}
	if p.AvailableCapacity() < count {
		return nil, ErrCapacityInsufficient
	}

	updates := map[string]interface{}{
		"used_capacity": gorm.Expr("used_capacity + ?", count),
		"version":       gorm.Expr("version + 1"),
	}

	// 所有产品类型统一使用 max_instances 管理库存
	result := global.GVA_DB.WithContext(ctx).Model(&productModel.Product{}).
		Where("id = ? AND version = ? AND (max_instances - used_capacity) >= ?", productId, p.Version, count).
		Updates(updates)

	if result.Error != nil {
		return nil, errors.New("系统繁忙，请重试")
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("资源已被抢占，请重试")
	}

	detailResp := toProductDetailResponse(&p)
	// 修正可用容量为扣减后的真实值
	detailResp.Available = p.AvailableCapacity() - count

	return &productRes.ReserveResult{Product: detailResp, ResourceCount: count}, nil
}

// ReleaseCapacity 释放产品资源（需指定数量）
func (s *ProductService) ReleaseCapacity(ctx context.Context, productId uint, count int64) error {
	if count <= 0 {
		return errors.New("释放数量必须大于0")
	}

	result := global.GVA_DB.WithContext(ctx).Model(&productModel.Product{}).
		Where("id = ? AND used_capacity >= ?", productId, count).
		Update("used_capacity", gorm.Expr("used_capacity - ?", count))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("释放失败：产品不存在或可释放数量不足")
	}
	return nil
}

// ReleaseCapacityAuto 自动释放 1 个实例配额
func (s *ProductService) ReleaseCapacityAuto(ctx context.Context, productId uint) error {
	return s.ReleaseCapacity(ctx, productId, 1)
}

var ErrCapacityInsufficient = errors.New("资源配额不足，请选择其他产品")

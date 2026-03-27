package cms

import (
	"context"
	"gin-vue-admin/global"
	clusterModel "gin-vue-admin/model/cluster"
	cmsReq "gin-vue-admin/model/cms/request"
	nbModel "gin-vue-admin/model/notebook"
	productModel "gin-vue-admin/model/product"
	productRes "gin-vue-admin/model/product/response"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ProductService struct{}

// CreateProduct 创建产品（对应一个Node节点或StorageClass）
func (s *ProductService) CreateProduct(req *cmsReq.CreateProductReq) error {
	// 默认产品类型为计算资源
	if req.ProductType == 0 {
		return errors.New("产品类型不能为空")
	}

	// 检查是否已存在（仅检查未删除的记录）
	var existing productModel.Product
	var err error
	switch req.ProductType {
	case productModel.ProductTypeCompute:
		// 计算资源：检查 NodeName
		if req.NodeName == "" {
			return errors.New("计算资源必须指定节点名称")
		}
		err = global.GVA_DB.Where("cluster_id = ? AND node_name = ? AND product_type = ? AND name = ?",
			req.ClusterId, req.NodeName, req.ProductType, req.Name).First(&existing).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Errorf("ProductService.CreateProduct global.GVA_DB.First(Calculation) failed: %v", err)
		}
	case productModel.ProductTypeStorage:
		// 文件存储：检查 StorageClass 和 Name
		if req.StorageClass == "" {
			return errors.New("文件存储必须指定StorageClass")
		}
		err = global.GVA_DB.Where("cluster_id = ? AND storage_class = ? AND product_type = ? AND name = ?",
			req.ClusterId, req.StorageClass, req.ProductType, req.Name).First(&existing).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Errorf("ProductService.CreateProduct global.GVA_DB.First(Storage) failed: %v", err)
		}
	}

	if err == nil {
		// 记录存在且未删除
		switch req.ProductType {
		case productModel.ProductTypeCompute:
			return errors.Errorf("该集群下已存在节点 %s 的产品名称为 %s 的配置", req.NodeName, req.Name)
		case productModel.ProductTypeStorage:
			return errors.Errorf("该集群下已存在 StorageClass %s 的产品名称为 %s 的配置", req.StorageClass, req.Name)
		}
	}

	// 记录不存在，创建新记录
	// 自动计算 MaxInstances：根据节点实际资源 / 产品规格
	maxInstances := req.MaxInstances
	if req.ProductType == productModel.ProductTypeCompute && req.NodeName != "" && maxInstances == 0 {
		calc := &NodeService{}
		cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(req.ClusterId)
		if cluster == nil {
			return errors.New("集群不存在或未连接")
		}

		isVGPU := req.VGPUNumber > 0 || req.VGPUMemory > 0 || req.VGPUCores > 0
		isGPU := req.GPUCount > 0 && !isVGPU

		nodeInfo, err := calc.CalculateNodeResources(context.Background(), cluster.ClientSet, req.NodeName)
		if err != nil {
			return errors.Wrap(err, "查询节点资源失败")
		}

		maxInstances = calc.CalculateMaxInstances(nodeInfo, req.GPUCount, req.VGPUNumber, req.VGPUMemory, req.VGPUCores, req.CPU, req.Memory, isGPU, isVGPU)

		if maxInstances <= 0 {
			return errors.New("该节点资源不足以创建此规格的产品")
		}
	}

	product := productModel.Product{
		ProductType:    req.ProductType,
		Name:           req.Name,
		Description:    req.Description,
		ClusterId:      req.ClusterId,
		Area:           req.Area,
		NodeName:       req.NodeName,
		NodeType:       req.NodeType,
		CPUModel:       req.CPUModel,
		CPU:            req.CPU,
		Memory:         req.Memory,
		GPUModel:       req.GPUModel,
		GPUCount:       req.GPUCount,
		GPUMemory:      req.GPUMemory,
		UsedCapacity:   0,
		VGPUNumber:     req.VGPUNumber,
		VGPUMemory:     req.VGPUMemory,
		PriceHourly:    req.PriceHourly,
		PriceDaily:     req.PriceDaily,
		PriceWeekly:    req.PriceWeekly,
		PriceMonthly:   req.PriceMonthly,
		DriverVersion:  req.DriverVersion,
		CUDAVersion:    req.CUDAVersion,
		SystemDisk:     req.SystemDisk,
		DataDisk:       req.DataDisk,
		Status:         req.Status,
		MaxInstances:   maxInstances,
		StorageClass:   req.StorageClass,
		StoragePriceGB: req.StoragePriceGB,
		VGPUCores:      req.VGPUCores,
	}

	if err := global.GVA_DB.Create(&product).Error; err != nil {
		return errors.Errorf("ProductService.CreateProduct global.GVA_DB.Create failed: %v", err)
	}

	logx.Info("创建产品成功", logx.Field("name", req.Name), logx.Field("maxInstances", maxInstances))
	return nil
}

// UpdateProduct 更新产品
// GPU / vGPU / CPU-only 三种类型互斥，更新某一类型的规格字段时自动清零其他类型的字段
func (s *ProductService) UpdateProduct(req *cmsReq.UpdateProductReq) error {
	// 1. 查询现有产品
	var existing productModel.Product
	if err := global.GVA_DB.Where("id = ?", req.ID).First(&existing).Error; err != nil {
		return errors.Wrap(err, "产品不存在")
	}

	updates := make(map[string]interface{})

	// 基础字段
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Area != "" {
		updates["area"] = req.Area
	}
	if req.Status >= 0 {
		updates["status"] = req.Status
	}
	if req.NodeType != "" {
		updates["node_type"] = req.NodeType
	}
	if req.CPUModel != "" {
		updates["cpu_model"] = req.CPUModel
	}
	if req.CPU > 0 {
		updates["cpu"] = req.CPU
	}
	if req.Memory > 0 {
		updates["memory"] = req.Memory
	}

	// 2. 处理 GPU / vGPU / CPU-only 互斥
	// 判断本次请求要设置的产品类型
	hasGPU := req.GPUCount > 0 || req.GPUModel != ""
	hasVGPU := req.VGPUNumber > 0 || req.VGPUMemory > 0 || req.VGPUCores > 0

	if hasGPU && hasVGPU {
		return errors.New("GPU 和 vGPU 不能同时设置，产品类型互斥")
	}

	if hasGPU {
		// 设置 GPU 字段，清零 vGPU 字段
		if req.GPUModel != "" {
			updates["gpu_model"] = req.GPUModel
		}
		if req.GPUCount > 0 {
			updates["gpu_count"] = req.GPUCount
		}
		if req.GPUMemory > 0 {
			updates["gpu_memory"] = req.GPUMemory
		}
		updates["v_gpu_number"] = 0
		updates["v_gpu_memory"] = 0
		updates["v_gpu_cores"] = 0
	} else if hasVGPU {
		// 设置 vGPU 字段，清零 GPU 字段
		if req.VGPUNumber > 0 {
			updates["v_gpu_number"] = req.VGPUNumber
		}
		if req.VGPUMemory > 0 {
			updates["v_gpu_memory"] = req.VGPUMemory
		}
		if req.VGPUCores > 0 {
			updates["v_gpu_cores"] = req.VGPUCores
		}
		updates["gpu_model"] = ""
		updates["gpu_count"] = 0
		updates["gpu_memory"] = 0
	}
	// 如果都没传，保持现有类型不变

	// 价格字段
	if req.PriceHourly > 0 {
		updates["price_hourly"] = req.PriceHourly
	}
	if req.PriceDaily > 0 {
		updates["price_daily"] = req.PriceDaily
	}
	if req.PriceWeekly > 0 {
		updates["price_weekly"] = req.PriceWeekly
	}
	if req.PriceMonthly > 0 {
		updates["price_monthly"] = req.PriceMonthly
	}
	if req.DriverVersion != "" {
		updates["driver_version"] = req.DriverVersion
	}
	if req.CUDAVersion != "" {
		updates["cuda_version"] = req.CUDAVersion
	}
	if req.SystemDisk > 0 {
		updates["system_disk"] = req.SystemDisk
	}
	if req.DataDisk > 0 {
		updates["data_disk"] = req.DataDisk
	}

	// 文件存储字段
	if req.StorageClass != "" {
		updates["storage_class"] = req.StorageClass
	}
	if req.StoragePriceGB > 0 {
		updates["storage_price_gb"] = req.StoragePriceGB
	}

	if len(updates) == 0 {
		return errors.New("没有需要更新的字段")
	}

	// 3. 重新计算 MaxInstances（如果规格字段有变更）
	specChanged := req.GPUCount > 0 || req.VGPUNumber > 0 || req.VGPUMemory > 0 || req.VGPUCores > 0 || req.CPU > 0 || req.Memory > 0
	if req.MaxInstances > 0 {
		// 管理员手动指定
		updates["max_instances"] = req.MaxInstances
	} else if specChanged && existing.NodeName != "" {
		// 规格变了，自动重新计算
		cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(existing.ClusterId)
		if cluster != nil && cluster.ClientSet != nil {
			calc := &NodeService{}
			nodeInfo, err := calc.CalculateNodeResources(context.Background(), cluster.ClientSet, existing.NodeName)
			if err == nil {
				// 合并：用请求值覆盖现有值
				gpuCount := existing.GPUCount
				if req.GPUCount > 0 {
					gpuCount = req.GPUCount
				}
				vgpuNumber := existing.VGPUNumber
				if req.VGPUNumber > 0 {
					vgpuNumber = req.VGPUNumber
				}
				vgpuMemory := existing.VGPUMemory
				if req.VGPUMemory > 0 {
					vgpuMemory = req.VGPUMemory
				}
				vgpuCores := existing.VGPUCores
				if req.VGPUCores > 0 {
					vgpuCores = req.VGPUCores
				}
				cpu := existing.CPU
				if req.CPU > 0 {
					cpu = req.CPU
				}
				memory := existing.Memory
				if req.Memory > 0 {
					memory = req.Memory
				}

				// 根据互斥后的最终类型计算
				finalIsVGPU := hasVGPU || (!hasGPU && existing.IsVGPU())
				finalIsGPU := hasGPU || (!hasVGPU && existing.IsGPUOnly())

				newMax := calc.CalculateMaxInstances(nodeInfo, gpuCount, vgpuNumber, vgpuMemory, vgpuCores, cpu, memory, finalIsGPU, finalIsVGPU)
				if newMax > 0 {
					updates["max_instances"] = newMax
				}
			}
		}
	}

	if err := global.GVA_DB.Model(&productModel.Product{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		return errors.Errorf("ProductService.UpdateProduct failed: %v", err)
	}

	logx.Info("更新产品成功", logx.Field("id", req.ID))
	return nil
}

// UpdatePrice 更新价格
func (s *ProductService) UpdatePrice(req *cmsReq.UpdatePriceReq) error {
	updates := make(map[string]interface{})

	if req.PriceHourly > 0 {
		updates["price_hourly"] = req.PriceHourly
	}
	if req.PriceDaily > 0 {
		updates["price_daily"] = req.PriceDaily
	}
	if req.PriceWeekly > 0 {
		updates["price_weekly"] = req.PriceWeekly
	}
	if req.PriceMonthly > 0 {
		updates["price_monthly"] = req.PriceMonthly
	}

	if len(updates) == 0 {
		return errors.New("没有需要更新的价格")
	}

	if err := global.GVA_DB.Model(&productModel.Product{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		return errors.Errorf("ProductService.UpdatePrice global.GVA_DB.Updates failed: %v", err)
	}

	logx.Info("更新价格成功")
	return nil
}

// DeleteProduct 删除产品
func (s *ProductService) DeleteProduct(req *cmsReq.DeleteProductReq) error {
	// 检查是否有实例在使用该产品
	var count int64
	if err := global.GVA_DB.Model(&nbModel.Notebook{}).Where("product_id = ?", req.ID).Count(&count).Error; err != nil {
		return errors.Errorf("ProductService.DeleteProduct global.GVA_DB.Count(Notebooks) failed: %v", err)
	}

	if count > 0 {
		return errors.Errorf("该产品有 %d 个实例正在使用，无法删除", count)
	}

	// 删除产品
	if err := global.GVA_DB.Delete(&productModel.Product{}, req.ID).Error; err != nil {
		return errors.Errorf("ProductService.DeleteProduct global.GVA_DB.Delete failed: %v", err)
	}

	logx.Info("删除产品成功")
	return nil
}

// GetProductList 获取产品列表
func (s *ProductService) GetProductList(req *cmsReq.GetProductListReq) (*productRes.ProductListResponse, error) {
	var products []productModel.Product
	var total int64

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	db := global.GVA_DB.Model(&productModel.Product{})

	// 按产品类型筛选
	if req.ProductType > 0 {
		db = db.Where("product_type = ?", req.ProductType)
	}

	// 按集群筛选
	if req.ClusterId > 0 {
		db = db.Where("cluster_id = ?", req.ClusterId)
	}

	// 按地区筛选
	if req.Area != "" {
		db = db.Where("area = ?", req.Area)
	}

	// 按GPU型号筛选
	if req.GPUModel != "" {
		db = db.Where("gpu_model = ?", req.GPUModel)
	}

	// 按状态筛选
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	// 关键词搜索（包括节点名称）
	if req.Keyword != "" {
		db = db.Where("name LIKE ? OR description LIKE ? OR node_name LIKE ? OR storage_class LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		logx.Error("查询产品总数失败", err)
		return nil, errors.Errorf("ProductService.GetProductList db.Count failed: %v", err)
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := db.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&products).Error; err != nil {
		logx.Error("查询产品列表失败", err)
		return nil, errors.Errorf("ProductService.GetProductList db.Find failed: %v", err)
	}

	var list []productRes.ProductDetailResponse
	for _, p := range products {
		list = append(list, productRes.ProductDetailResponse{
			ID:             p.ID,
			ClusterId:      p.ClusterId,
			ProductType:    p.ProductType,
			Name:           p.Name,
			Description:    p.Description,
			Area:           p.Area,
			NodeName:       p.NodeName,
			NodeType:       p.NodeType,
			CPUModel:       p.CPUModel,
			CPU:            p.CPU,
			Memory:         p.Memory,
			GPUModel:       p.GPUModel,
			GPUCount:       p.GPUCount,
			GPUMemory:      p.GPUMemory,
			VGPUNumber:     p.VGPUNumber,
			VGPUMemory:     p.VGPUMemory,
			VGPUCores:      p.VGPUCores,
			PriceHourly:    p.PriceHourly,
			PriceDaily:     p.PriceDaily,
			PriceWeekly:    p.PriceWeekly,
			PriceMonthly:   p.PriceMonthly,
			DriverVersion:  p.DriverVersion,
			CUDAVersion:    p.CUDAVersion,
			SystemDisk:     p.SystemDisk,
			DataDisk:       p.DataDisk,
			Status:         p.Status,
			StorageClass:   p.StorageClass,
			StoragePriceGB: p.StoragePriceGB,
			MaxInstances:   p.MaxInstances,
		})
	}
	resp := &productRes.ProductListResponse{
		List:  list,
		Total: total,
	}

	return resp, nil
}

// GetProductDetail 获取产品详情
func (s *ProductService) GetProductDetail(idStr string) (*productRes.ProductDetailResponse, error) {
	var product productModel.Product
	if err := global.GVA_DB.Where("id = ?", idStr).First(&product).Error; err != nil {
		return nil, errors.Errorf("ProductService.GetProductDetail global.GVA_DB.First failed: %v", err)
	}

	resp := &productRes.ProductDetailResponse{
		ID:             product.ID,
		ProductType:    product.ProductType,
		Name:           product.Name,
		Description:    product.Description,
		Area:           product.Area,
		NodeName:       product.NodeName,
		NodeType:       product.NodeType,
		CPUModel:       product.CPUModel,
		CPU:            product.CPU,
		Memory:         product.Memory,
		GPUModel:       product.GPUModel,
		GPUCount:       product.GPUCount,
		GPUMemory:      product.GPUMemory,
		VGPUNumber:     product.VGPUNumber,
		VGPUMemory:     product.VGPUMemory,
		PriceHourly:    product.PriceHourly,
		PriceDaily:     product.PriceDaily,
		PriceWeekly:    product.PriceWeekly,
		PriceMonthly:   product.PriceMonthly,
		DriverVersion:  product.DriverVersion,
		CUDAVersion:    product.CUDAVersion,
		SystemDisk:     product.SystemDisk,
		DataDisk:       product.DataDisk,
		Status:         product.Status,
		StorageClass:   product.StorageClass,
		StoragePriceGB: product.StoragePriceGB,
		MaxInstances:   product.MaxInstances,
	}
	return resp, nil
}

// GetClusterList 获取集群列表
func (s *ProductService) GetClusterList() ([]productRes.ClusterResponse, error) {
	var clusters []clusterModel.K8sCluster
	if err := global.GVA_DB.Where("status = ?", clusterModel.ClusterStatusEnabled).Find(&clusters).Error; err != nil {
		return nil, errors.Errorf("ProductService.GetClusterList global.GVA_DB.Find failed: %v", err)
	}

	var list []productRes.ClusterResponse
	for _, c := range clusters {
		list = append(list, productRes.ClusterResponse{
			ID:          c.ID,
			Name:        c.Name,
			Area:        c.Area,
			Description: c.Description,
			Status:      c.Status,
		})
	}
	return list, nil
}

// GetAreaList 获取地区列表
func (s *ProductService) GetAreaList() ([]string, error) {
	var areas []string
	if err := global.GVA_DB.Model(&productModel.Product{}).Distinct("area").Pluck("area", &areas).Error; err != nil {
		return nil, errors.Errorf("ProductService.GetAreaList global.GVA_DB.Pluck failed: %v", err)
	}
	return areas, nil
}

// SyncNodeGPUUsage 同步节点GPU使用情况（用于更新已用卡数）
func (s *ProductService) SyncNodeGPUUsage(productId int64, usedGPU int) error {
	if err := global.GVA_DB.Model(&productModel.Product{}).Where("id = ?", productId).
		Update("used_capacity", usedGPU).Error; err != nil {
		return errors.Errorf("ProductService.SyncNodeGPUUsage global.GVA_DB.Update failed: %v", err)
	}
	return nil
}

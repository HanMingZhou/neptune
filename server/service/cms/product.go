package cms

import (
	"context"
	"gin-vue-admin/global"
	clusterModel "gin-vue-admin/model/cluster"
	cmsReq "gin-vue-admin/model/cms/request"
	cmsRes "gin-vue-admin/model/cms/response"
	nbModel "gin-vue-admin/model/notebook"
	productModel "gin-vue-admin/model/product"
	productRes "gin-vue-admin/model/product/response"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ProductService struct{}

type productNodeKey struct {
	clusterID uint
	nodeName  string
}

type productNodeMetadata struct {
	nodeIP        string
	driverVersion string
	cudaVersion   string
}

func (s *ProductService) resolveCatalogPriceType(req *cmsReq.GetProductListReq) int {
	return productModel.NormalizeCatalogPriceType(req.ProductType, req.PriceType)
}

func buildComputePriceItems(productID uint, items []cmsReq.ProductPriceItem) ([]productModel.ProductPrice, error) {
	if len(items) == 0 {
		return nil, nil
	}

	priceMap := make(map[int]float64, len(items))
	for _, item := range items {
		if !productModel.IsComputePriceType(item.PriceType) {
			return nil, errors.Errorf("不支持的价格类型: %d", item.PriceType)
		}
		if _, exists := priceMap[item.PriceType]; exists {
			return nil, errors.Errorf("价格类型 %d 重复", item.PriceType)
		}
		priceMap[item.PriceType] = item.Price
	}

	result := make([]productModel.ProductPrice, 0, len(priceMap))
	for _, priceType := range productModel.ComputePriceTypes() {
		price, ok := priceMap[priceType]
		if !ok {
			continue
		}
		result = append(result, productModel.ProductPrice{
			ProductID: productID,
			PriceType: priceType,
			Price:     price,
		})
	}

	return result, nil
}

func normalizeNodeNames(nodeNames []string) []string {
	seen := make(map[string]struct{}, len(nodeNames))
	normalized := make([]string, 0, len(nodeNames))
	for _, nodeName := range nodeNames {
		trimmed := strings.TrimSpace(nodeName)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; exists {
			continue
		}
		seen[trimmed] = struct{}{}
		normalized = append(normalized, trimmed)
	}

	sort.Strings(normalized)
	return normalized
}

func determineComputeResourceMode(gpuCount, gpuMemory, vgpuNumber, vgpuMemory, vgpuCores int64) (bool, bool, error) {
	isVGPU := vgpuNumber > 0 || vgpuMemory > 0 || vgpuCores > 0
	isGPU := gpuCount > 0 || gpuMemory > 0
	if isGPU && isVGPU {
		return false, false, errors.New("GPU 和 vGPU 不能同时设置")
	}
	return isGPU, isVGPU, nil
}

func buildUpdatedComputeSpec(existing productModel.Product, req *cmsReq.UpdateProductReq, hasGPU, hasVGPU bool) *cmsReq.CreateProductReq {
	spec := &cmsReq.CreateProductReq{
		ProductType:   existing.ProductType,
		ClusterId:     existing.ClusterId,
		Area:          existing.Area,
		NodeName:      existing.NodeName,
		CPUModel:      firstNonEmpty(req.CPUModel, existing.CPUModel),
		CPU:           existing.CPU,
		Memory:        existing.Memory,
		GPUModel:      firstNonEmpty(req.GPUModel, existing.GPUModel),
		GPUCount:      existing.GPUCount,
		GPUMemory:     existing.GPUMemory,
		VGPUNumber:    existing.VGPUNumber,
		VGPUMemory:    existing.VGPUMemory,
		VGPUCores:     existing.VGPUCores,
		DriverVersion: firstNonEmpty(req.DriverVersion, existing.DriverVersion),
		CUDAVersion:   firstNonEmpty(req.CUDAVersion, existing.CUDAVersion),
		SystemDisk:    existing.SystemDisk,
		DataDisk:      existing.DataDisk,
		Status:        existing.Status,
		MaxInstances:  existing.MaxInstances,
	}

	if req.CPU > 0 {
		spec.CPU = req.CPU
	}
	if req.Memory > 0 {
		spec.Memory = req.Memory
	}

	switch {
	case hasGPU:
		if req.GPUCount > 0 {
			spec.GPUCount = req.GPUCount
		}
		if req.GPUMemory > 0 {
			spec.GPUMemory = req.GPUMemory
		}
		spec.VGPUNumber = 0
		spec.VGPUMemory = 0
		spec.VGPUCores = 0
	case hasVGPU:
		if req.VGPUNumber > 0 {
			spec.VGPUNumber = req.VGPUNumber
		}
		if req.VGPUMemory > 0 {
			spec.VGPUMemory = req.VGPUMemory
		}
		if req.VGPUCores > 0 {
			spec.VGPUCores = req.VGPUCores
		}
		spec.GPUCount = 0
		spec.GPUMemory = 0
	}

	return spec
}

func buildExistingComputeProductSummary(product productModel.Product) cmsRes.ExistingComputeProductSummary {
	resourceType := "cpu"
	if product.IsVGPU() {
		resourceType = "vgpu"
	} else if product.IsGPUOnly() {
		resourceType = "gpu"
	}

	return cmsRes.ExistingComputeProductSummary{
		ID:           product.ID,
		Name:         product.Name,
		Description:  product.Description,
		ResourceType: resourceType,
		CPU:          product.CPU,
		Memory:       product.Memory,
		GPUModel:     product.GPUModel,
		GPUCount:     product.GPUCount,
		GPUMemory:    product.GPUMemory,
		VGPUNumber:   product.VGPUNumber,
		VGPUMemory:   product.VGPUMemory,
		VGPUCores:    product.VGPUCores,
		PriceHourly:  product.PriceHourly,
		PriceDaily:   product.PriceDaily,
		PriceWeekly:  product.PriceWeekly,
		PriceMonthly: product.PriceMonthly,
		Status:       product.Status,
		MaxInstances: product.MaxInstances,
		UsedCapacity: product.UsedCapacity,
		Available:    product.AvailableCapacity(),
	}
}

func (s *ProductService) loadExistingComputeProducts(tx *gorm.DB, clusterID uint, nodeNames []string, excludeProductID uint) (map[string][]productModel.Product, error) {
	result := make(map[string][]productModel.Product)
	normalizedNodeNames := normalizeNodeNames(nodeNames)
	if clusterID == 0 || len(normalizedNodeNames) == 0 {
		return result, nil
	}

	db := tx
	if db == nil {
		db = global.GVA_DB
	}

	query := db.Where("cluster_id = ? AND product_type = ? AND node_name IN ?", clusterID, productModel.ProductTypeCompute, normalizedNodeNames)
	if excludeProductID > 0 {
		query = query.Where("id <> ?", excludeProductID)
	}

	var products []productModel.Product
	if err := query.Order("id ASC").Find(&products).Error; err != nil {
		return nil, err
	}
	if err := productModel.LoadPriceItemsForProducts(context.Background(), db, products); err != nil {
		return nil, err
	}

	for _, product := range products {
		result[product.NodeName] = append(result[product.NodeName], product)
	}

	return result, nil
}

func (s *ProductService) hydrateComputeProductRequest(req *cmsReq.CreateProductReq, cluster *global.ClusterClientInfo, nodeInfo *cmsRes.NodeInfoResponse) {
	req.NodeName = strings.TrimSpace(req.NodeName)
	if req.Area == "" && cluster != nil {
		req.Area = cluster.Area
	}
	if req.CPUModel == "" {
		req.CPUModel = nodeInfo.CPUModel
	}
	if req.GPUModel == "" {
		req.GPUModel = nodeInfo.GPUModel
	}
	if req.DriverVersion == "" {
		req.DriverVersion = nodeInfo.DriverVersion
	}
	if req.CUDAVersion == "" {
		req.CUDAVersion = nodeInfo.CUDAVersion
	}
}

func (s *ProductService) validateComputeProductSpec(nodeInfo *cmsRes.NodeInfoResponse, req *cmsReq.CreateProductReq) (int64, error) {
	if nodeInfo == nil {
		return 0, errors.New("节点资源信息不存在")
	}
	if req.CPU <= 0 {
		return 0, errors.New("CPU 规格必须大于 0")
	}
	if req.Memory <= 0 {
		return 0, errors.New("内存规格必须大于 0")
	}

	isGPU, isVGPU, err := determineComputeResourceMode(req.GPUCount, req.GPUMemory, req.VGPUNumber, req.VGPUMemory, req.VGPUCores)
	if err != nil {
		return 0, err
	}

	if req.CPU > nodeInfo.CPUAllocatable {
		return 0, errors.Errorf("节点 %s 的 CPU 资源不足", nodeInfo.NodeName)
	}
	if req.Memory > nodeInfo.MemoryAllocatable {
		return 0, errors.Errorf("节点 %s 的内存资源不足", nodeInfo.NodeName)
	}

	switch {
	case isGPU:
		if req.GPUCount <= 0 {
			return 0, errors.New("GPU 产品必须配置卡数")
		}
		if nodeInfo.GPUCount <= 0 {
			return 0, errors.Errorf("节点 %s 不支持 GPU 产品", nodeInfo.NodeName)
		}
		if req.GPUCount > nodeInfo.GPUCount {
			return 0, errors.Errorf("节点 %s 的 GPU 卡数不足", nodeInfo.NodeName)
		}
		if req.GPUMemory > 0 && nodeInfo.GPUMemory > 0 && req.GPUMemory > nodeInfo.GPUMemory {
			return 0, errors.Errorf("节点 %s 的单卡显存不足", nodeInfo.NodeName)
		}
	case isVGPU:
		if nodeInfo.VGPUNumber <= 0 && nodeInfo.VGPUMemory <= 0 && nodeInfo.VGPUCores <= 0 {
			return 0, errors.Errorf("节点 %s 不支持 vGPU 产品", nodeInfo.NodeName)
		}
		if req.VGPUNumber > 0 && req.VGPUNumber > nodeInfo.VGPUNumber {
			return 0, errors.Errorf("节点 %s 的 vGPU 数量不足", nodeInfo.NodeName)
		}
		if req.VGPUMemory > 0 && req.VGPUMemory > nodeInfo.VGPUMemory {
			return 0, errors.Errorf("节点 %s 的 vGPU 显存不足", nodeInfo.NodeName)
		}
		if req.VGPUCores > 0 && req.VGPUCores > nodeInfo.VGPUCores {
			return 0, errors.Errorf("节点 %s 的 vGPU 算力不足", nodeInfo.NodeName)
		}
	}

	calc := &NodeService{}
	maxInstances := calc.CalculateMaxInstances(
		nodeInfo,
		req.GPUCount,
		req.VGPUNumber,
		req.VGPUMemory,
		req.VGPUCores,
		req.CPU,
		req.Memory,
		isGPU,
		isVGPU,
	)
	if maxInstances <= 0 {
		return 0, errors.New("该节点资源不足以创建此规格的产品")
	}

	return maxInstances, nil
}

func (s *ProductService) prepareComputeProduct(tx *gorm.DB, req *cmsReq.CreateProductReq) (*productModel.Product, error) {
	if req.ProductType == 0 {
		req.ProductType = productModel.ProductTypeCompute
	}
	if req.ProductType != productModel.ProductTypeCompute {
		return nil, errors.New("仅支持计算产品")
	}

	req.NodeName = strings.TrimSpace(req.NodeName)
	if req.NodeName == "" {
		return nil, errors.New("计算资源必须指定节点名称")
	}

	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(req.ClusterId)
	if cluster == nil || cluster.ClientSet == nil {
		return nil, errors.New("集群不存在或未连接")
	}

	nodeService := &NodeService{}
	nodeInfo, err := nodeService.CalculateNodeResources(context.Background(), cluster.ClientSet, req.NodeName)
	if err != nil {
		return nil, errors.Wrap(err, "查询节点资源失败")
	}
	if !nodeInfo.Schedulable {
		return nil, errors.Errorf("节点 %s 当前不可调度，不能新增计算产品", req.NodeName)
	}

	s.hydrateComputeProductRequest(req, cluster, nodeInfo)

	existingProducts, err := s.loadExistingComputeProducts(tx, req.ClusterId, []string{req.NodeName}, 0)
	if err != nil {
		return nil, errors.Wrap(err, "查询节点已有计算产品失败")
	}
	if len(existingProducts[req.NodeName]) > 0 {
		return nil, errors.Errorf("节点 %s 已存在计算产品，不能重复新增", req.NodeName)
	}

	maxInstances, err := s.validateComputeProductSpec(nodeInfo, req)
	if err != nil {
		return nil, err
	}
	if req.MaxInstances > 0 {
		maxInstances = req.MaxInstances
	}

	product := &productModel.Product{
		ProductType:   productModel.ProductTypeCompute,
		Name:          req.Name,
		Description:   req.Description,
		ClusterId:     req.ClusterId,
		Area:          req.Area,
		NodeName:      req.NodeName,
		NodeType:      req.NodeType,
		CPUModel:      req.CPUModel,
		CPU:           req.CPU,
		Memory:        req.Memory,
		GPUModel:      req.GPUModel,
		GPUCount:      req.GPUCount,
		GPUMemory:     req.GPUMemory,
		UsedCapacity:  0,
		VGPUNumber:    req.VGPUNumber,
		VGPUMemory:    req.VGPUMemory,
		DriverVersion: req.DriverVersion,
		CUDAVersion:   req.CUDAVersion,
		SystemDisk:    req.SystemDisk,
		DataDisk:      req.DataDisk,
		Status:        req.Status,
		MaxInstances:  maxInstances,
		VGPUCores:     req.VGPUCores,
	}

	return product, nil
}

func (s *ProductService) resolveNodeCandidateCompatibility(node *cmsRes.NodeInfoResponse, req *cmsReq.GetProductNodeCandidatesReq) (bool, string) {
	if node == nil {
		return false, "节点资源信息不存在"
	}
	if !node.Schedulable {
		return false, "节点当前不可调度"
	}

	resourceType := strings.ToLower(strings.TrimSpace(req.ResourceType))
	switch {
	case resourceType == "vgpu" || req.VGPUNumber > 0 || req.VGPUMemory > 0 || req.VGPUCores > 0:
		if node.VGPUNumber <= 0 && node.VGPUMemory <= 0 && node.VGPUCores <= 0 {
			return false, "该节点不支持 vGPU 规格"
		}
		if req.VGPUNumber > 0 && req.VGPUNumber > node.VGPUNumber {
			return false, "节点 vGPU 数量不足"
		}
		if req.VGPUMemory > 0 && req.VGPUMemory > node.VGPUMemory {
			return false, "节点 vGPU 显存不足"
		}
		if req.VGPUCores > 0 && req.VGPUCores > node.VGPUCores {
			return false, "节点 vGPU 算力不足"
		}
	case resourceType == "gpu" || req.GPUCount > 0 || req.GPUMemory > 0:
		if node.GPUCount <= 0 {
			return false, "该节点不支持 GPU 规格"
		}
		if req.GPUCount > 0 && req.GPUCount > node.GPUCount {
			return false, "节点 GPU 卡数不足"
		}
		if req.GPUMemory > 0 && node.GPUMemory > 0 && req.GPUMemory > node.GPUMemory {
			return false, "节点单卡显存不足"
		}
	}

	if req.CPU > 0 && req.CPU > node.CPUAllocatable {
		return false, "节点 CPU 不满足当前规格"
	}
	if req.Memory > 0 && req.Memory > node.MemoryAllocatable {
		return false, "节点内存不满足当前规格"
	}

	return true, ""
}

// CreateProduct 创建产品（对应一个Node节点或StorageClass）
func (s *ProductService) CreateProduct(req *cmsReq.CreateProductReq) error {
	// 默认产品类型为计算资源
	if req.ProductType == 0 {
		return errors.New("产品类型不能为空")
	}

	switch req.ProductType {
	case productModel.ProductTypeCompute:
		if len(req.Prices) == 0 {
			return errors.New("计算产品价格不能为空")
		}
		product, err := s.prepareComputeProduct(global.GVA_DB, req)
		if err != nil {
			return err
		}
		if err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(product).Error; err != nil {
				return errors.Errorf("ProductService.CreateProduct global.GVA_DB.Create failed: %v", err)
			}

			priceItems, err := buildComputePriceItems(product.ID, req.Prices)
			if err != nil {
				return err
			}
			if err := productModel.UpsertPriceItems(context.Background(), tx, priceItems); err != nil {
				return errors.Errorf("ProductService.CreateProduct save prices failed: %v", err)
			}
			return nil
		}); err != nil {
			return err
		}

		logx.Info("创建产品成功", logx.Field("name", req.Name), logx.Field("maxInstances", product.MaxInstances))
		return nil
	case productModel.ProductTypeStorage:
		// 检查是否已存在（仅检查未删除的记录）
		var existing productModel.Product
		// 文件存储：检查 StorageClass 和 Name
		if req.StorageClass == "" {
			return errors.New("文件存储必须指定StorageClass")
		}
		err := global.GVA_DB.Where("cluster_id = ? AND storage_class = ? AND product_type = ? AND name = ?",
			req.ClusterId, req.StorageClass, req.ProductType, req.Name).First(&existing).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Errorf("ProductService.CreateProduct global.GVA_DB.First(Storage) failed: %v", err)
		}
		if err == nil {
			return errors.Errorf("该集群下已存在 StorageClass %s 的产品名称为 %s 的配置", req.StorageClass, req.Name)
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
			DriverVersion:  req.DriverVersion,
			CUDAVersion:    req.CUDAVersion,
			SystemDisk:     req.SystemDisk,
			DataDisk:       req.DataDisk,
			Status:         req.Status,
			MaxInstances:   req.MaxInstances,
			StorageClass:   req.StorageClass,
			StoragePriceGB: req.StoragePriceGB,
			VGPUCores:      req.VGPUCores,
		}

		if err := global.GVA_DB.Create(&product).Error; err != nil {
			return errors.Errorf("ProductService.CreateProduct global.GVA_DB.Create failed: %v", err)
		}

		logx.Info("创建产品成功", logx.Field("name", req.Name))
		return nil
	default:
		return errors.New("不支持的产品类型")
	}
}

func (s *ProductService) GetProductNodeCandidates(req *cmsReq.GetProductNodeCandidatesReq) ([]cmsRes.ProductNodeCandidateResponse, error) {
	if req.ClusterId == 0 {
		return nil, errors.New("集群ID不能为空")
	}

	nodeService := &NodeService{}
	nodes, err := nodeService.GetClusterNodesWithResources(context.Background(), req.ClusterId, "")
	if err != nil {
		return nil, err
	}

	nodeNames := make([]string, 0, len(nodes))
	for _, node := range nodes {
		nodeNames = append(nodeNames, node.NodeName)
	}

	existingProducts, err := s.loadExistingComputeProducts(global.GVA_DB, req.ClusterId, nodeNames, req.ExcludeProductId)
	if err != nil {
		return nil, errors.Wrap(err, "查询节点已有计算产品失败")
	}

	candidates := make([]cmsRes.ProductNodeCandidateResponse, 0, len(nodes))
	for _, node := range nodes {
		compatible, compatibilityReason := s.resolveNodeCandidateCompatibility(&node, req)
		summaries := make([]cmsRes.ExistingComputeProductSummary, 0, len(existingProducts[node.NodeName]))
		for _, product := range existingProducts[node.NodeName] {
			summaries = append(summaries, buildExistingComputeProductSummary(product))
		}

		disableReason := ""
		canCreate := compatible && len(summaries) == 0
		if len(summaries) > 0 {
			disableReason = "该节点已存在计算产品"
		} else if !compatible {
			disableReason = compatibilityReason
		}

		candidates = append(candidates, cmsRes.ProductNodeCandidateResponse{
			NodeInfoResponse:        node,
			ExistingComputeProducts: summaries,
			CanCreateComputeProduct: canCreate,
			Compatible:              compatible,
			DisableReason:           disableReason,
		})
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].CanCreateComputeProduct != candidates[j].CanCreateComputeProduct {
			return candidates[i].CanCreateComputeProduct
		}
		if candidates[i].Compatible != candidates[j].Compatible {
			return candidates[i].Compatible
		}
		return strings.ToLower(candidates[i].NodeName) < strings.ToLower(candidates[j].NodeName)
	})

	return candidates, nil
}

func (s *ProductService) BatchCreateComputeProducts(req *cmsReq.BatchCreateComputeProductReq) (*cmsRes.BatchCreateComputeProductResponse, error) {
	if req.ProductType == 0 {
		req.ProductType = productModel.ProductTypeCompute
	}
	if req.ProductType != productModel.ProductTypeCompute {
		return nil, errors.New("仅支持批量创建计算产品")
	}
	if len(req.Prices) == 0 {
		return nil, errors.New("计算产品价格不能为空")
	}

	nodeNames := normalizeNodeNames(req.NodeNames)
	if len(nodeNames) == 0 {
		return nil, errors.New("请至少选择一个节点")
	}

	createdIDs := make([]uint, 0, len(nodeNames))
	if err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for _, nodeName := range nodeNames {
			createReq := cmsReq.CreateProductReq{
				ProductType:   productModel.ProductTypeCompute,
				Name:          req.Name,
				Description:   req.Description,
				ClusterId:     req.ClusterId,
				Area:          req.Area,
				NodeName:      nodeName,
				NodeType:      req.NodeType,
				CPUModel:      req.CPUModel,
				CPU:           req.CPU,
				Memory:        req.Memory,
				GPUModel:      req.GPUModel,
				GPUCount:      req.GPUCount,
				GPUMemory:     req.GPUMemory,
				VGPUNumber:    req.VGPUNumber,
				VGPUMemory:    req.VGPUMemory,
				VGPUCores:     req.VGPUCores,
				Prices:        append([]cmsReq.ProductPriceItem(nil), req.Prices...),
				DriverVersion: req.DriverVersion,
				CUDAVersion:   req.CUDAVersion,
				SystemDisk:    req.SystemDisk,
				DataDisk:      req.DataDisk,
				Status:        req.Status,
				MaxInstances:  req.MaxInstances,
			}

			product, err := s.prepareComputeProduct(tx, &createReq)
			if err != nil {
				return err
			}
			if err := tx.Create(product).Error; err != nil {
				return errors.Errorf("批量创建节点 %s 的产品失败: %v", nodeName, err)
			}
			priceItems, err := buildComputePriceItems(product.ID, createReq.Prices)
			if err != nil {
				return err
			}
			if err := productModel.UpsertPriceItems(context.Background(), tx, priceItems); err != nil {
				return errors.Errorf("批量创建节点 %s 的价格失败: %v", nodeName, err)
			}

			createdIDs = append(createdIDs, product.ID)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &cmsRes.BatchCreateComputeProductResponse{
		CreatedIds:   createdIDs,
		CreatedCount: len(createdIDs),
		CreatedNodes: nodeNames,
	}, nil
}

// UpdateProduct 更新产品
// GPU / vGPU / CPU ONLY 三种类型互斥，更新某一类型的规格字段时自动清零其他类型的字段
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

	// 2. 处理 GPU / vGPU / CPU ONLY 互斥
	// 判断本次请求要设置的产品类型
	hasGPU := req.GPUCount > 0 || req.GPUMemory > 0
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
		// 设置 vGPU 字段，保留底层 GPU 型号元数据，清零 GPU 规格字段
		if req.VGPUNumber > 0 {
			updates["v_gpu_number"] = req.VGPUNumber
		}
		if req.VGPUMemory > 0 {
			updates["v_gpu_memory"] = req.VGPUMemory
		}
		if req.VGPUCores > 0 {
			updates["v_gpu_cores"] = req.VGPUCores
		}
		if req.GPUModel != "" {
			updates["gpu_model"] = req.GPUModel
		}
		updates["gpu_count"] = 0
		updates["gpu_memory"] = 0
	}
	// 如果都没传，保持现有类型不变

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

	if len(updates) == 0 && len(req.Prices) == 0 {
		return errors.New("没有需要更新的字段")
	}

	// 3. 重新计算 MaxInstances（如果规格字段有变更）
	specChanged := req.GPUCount > 0 || req.GPUMemory > 0 || req.VGPUNumber > 0 || req.VGPUMemory > 0 || req.VGPUCores > 0 || req.CPU > 0 || req.Memory > 0
	if specChanged && existing.NodeName != "" {
		// 规格变了，自动重新计算
		cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(existing.ClusterId)
		if cluster != nil && cluster.ClientSet != nil {
			calc := &NodeService{}
			nodeInfo, err := calc.CalculateNodeResources(context.Background(), cluster.ClientSet, existing.NodeName)
			if err == nil {
				validationReq := buildUpdatedComputeSpec(existing, req, hasGPU, hasVGPU)

				newMax, validateErr := s.validateComputeProductSpec(nodeInfo, validationReq)
				if validateErr != nil {
					return validateErr
				}
				if req.MaxInstances > 0 {
					newMax = req.MaxInstances
				}
				updates["max_instances"] = newMax
			}
		}
	} else if req.MaxInstances > 0 {
		// 管理员手动指定
		updates["max_instances"] = req.MaxInstances
	}

	if err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if len(updates) > 0 {
			if err := tx.Model(&productModel.Product{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
				return errors.Errorf("ProductService.UpdateProduct failed: %v", err)
			}
		}

		if existing.ProductType == productModel.ProductTypeCompute && len(req.Prices) > 0 {
			priceItems, err := buildComputePriceItems(req.ID, req.Prices)
			if err != nil {
				return err
			}
			if err := productModel.UpsertPriceItems(context.Background(), tx, priceItems); err != nil {
				return errors.Errorf("ProductService.UpdateProduct save prices failed: %v", err)
			}
		}
		return nil
	}); err != nil {
		return err
	}

	logx.Info("更新产品成功", logx.Field("id", req.ID))
	return nil
}

// UpdatePrice 更新价格
func (s *ProductService) UpdatePrice(req *cmsReq.UpdatePriceReq) error {
	if len(req.Prices) == 0 {
		return errors.New("没有需要更新的价格")
	}

	priceItems, err := buildComputePriceItems(req.ID, req.Prices)
	if err != nil {
		return err
	}
	if len(priceItems) == 0 {
		return errors.New("没有需要更新的价格")
	}

	if err := productModel.UpsertPriceItems(context.Background(), global.GVA_DB, priceItems); err != nil {
		return errors.Errorf("ProductService.UpdatePrice save prices failed: %v", err)
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
	if err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := productModel.DeletePriceItems(context.Background(), tx, req.ID); err != nil {
			return errors.Errorf("ProductService.DeleteProduct delete prices failed: %v", err)
		}
		if err := tx.Delete(&productModel.Product{}, req.ID).Error; err != nil {
			return errors.Errorf("ProductService.DeleteProduct global.GVA_DB.Delete failed: %v", err)
		}
		return nil
	}); err != nil {
		return err
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

	// 按资源类型筛选
	switch req.ResourceType {
	case "cpu":
		db = db.Where("COALESCE(gpu_count, 0) = 0").
			Where("COALESCE(v_gpu_number, 0) = 0").
			Where("COALESCE(v_gpu_memory, 0) = 0").
			Where("COALESCE(v_gpu_cores, 0) = 0")
	case "gpu":
		db = db.Where("COALESCE(gpu_count, 0) > 0").
			Where("COALESCE(v_gpu_number, 0) = 0").
			Where("COALESCE(v_gpu_memory, 0) = 0").
			Where("COALESCE(v_gpu_cores, 0) = 0")
	case "vgpu":
		db = db.Where("(COALESCE(v_gpu_number, 0) > 0 OR COALESCE(v_gpu_memory, 0) > 0 OR COALESCE(v_gpu_cores, 0) > 0)")
	}

	// 按GPU型号筛选
	if req.GPUModel != "" {
		db = db.Where("gpu_model LIKE ?", "%"+req.GPUModel+"%")
	}

	// 按状态筛选
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	if req.AvailableMin != nil {
		db = db.Where("(max_instances - used_capacity) >= ?", *req.AvailableMin)
	}

	if req.AvailableMax != nil {
		db = db.Where("(max_instances - used_capacity) <= ?", *req.AvailableMax)
	}

	if req.MaxInstancesMin != nil {
		db = db.Where("max_instances >= ?", *req.MaxInstancesMin)
	}

	if req.MaxInstancesMax != nil {
		db = db.Where("max_instances <= ?", *req.MaxInstancesMax)
	}

	if req.UsedCapacityMin != nil {
		db = db.Where("used_capacity >= ?", *req.UsedCapacityMin)
	}

	if req.UsedCapacityMax != nil {
		db = db.Where("used_capacity <= ?", *req.UsedCapacityMax)
	}

	priceType := s.resolveCatalogPriceType(req)
	if req.ProductType == productModel.ProductTypeStorage || priceType == productModel.PriceTypeStorageGBDaily {
		if req.PriceMin != nil {
			db = db.Where("storage_price_gb >= ?", *req.PriceMin)
		}
		if req.PriceMax != nil {
			db = db.Where("storage_price_gb <= ?", *req.PriceMax)
		}
	} else if req.PriceMin != nil || req.PriceMax != nil {
		priceQuery := global.GVA_DB.Model(&productModel.ProductPrice{}).
			Select("product_id").
			Where("price_type = ?", priceType)
		if req.PriceMin != nil {
			priceQuery = priceQuery.Where("price >= ?", *req.PriceMin)
		}
		if req.PriceMax != nil {
			priceQuery = priceQuery.Where("price <= ?", *req.PriceMax)
		}
		db = db.Where("id IN (?)", priceQuery)
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
	if err := productModel.LoadPriceItemsForProducts(context.Background(), global.GVA_DB, products); err != nil {
		return nil, errors.Errorf("ProductService.GetProductList load prices failed: %v", err)
	}

	clusterNames, err := s.loadClusterNames(products)
	if err != nil {
		return nil, errors.Errorf("ProductService.GetProductList loadClusterNames failed: %v", err)
	}
	nodeMetadata := s.loadProductNodeMetadata(context.Background(), products)

	var list []productRes.ProductDetailResponse
	for _, p := range products {
		metadata := nodeMetadata[productNodeKey{
			clusterID: p.ClusterId,
			nodeName:  p.NodeName,
		}]
		list = append(list, productRes.ProductDetailResponse{
			ID:             p.ID,
			ClusterId:      p.ClusterId,
			ClusterName:    clusterNames[p.ClusterId],
			ProductType:    p.ProductType,
			Name:           p.Name,
			Description:    p.Description,
			Area:           p.Area,
			NodeName:       p.NodeName,
			NodeIP:         metadata.nodeIP,
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
			DriverVersion:  firstNonEmpty(p.DriverVersion, metadata.driverVersion),
			CUDAVersion:    firstNonEmpty(p.CUDAVersion, metadata.cudaVersion),
			SystemDisk:     p.SystemDisk,
			DataDisk:       p.DataDisk,
			Status:         p.Status,
			StorageClass:   p.StorageClass,
			StoragePriceGB: p.StoragePriceGB,
			MaxInstances:   p.MaxInstances,
			UsedCapacity:   p.UsedCapacity,
			Available:      p.AvailableCapacity(),
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
	if err := productModel.LoadPriceItems(context.Background(), global.GVA_DB, &product); err != nil {
		return nil, errors.Errorf("ProductService.GetProductDetail load prices failed: %v", err)
	}

	clusterNames, err := s.loadClusterNames([]productModel.Product{product})
	if err != nil {
		return nil, errors.Errorf("ProductService.GetProductDetail loadClusterNames failed: %v", err)
	}
	nodeMetadata := s.loadProductNodeMetadata(context.Background(), []productModel.Product{product})
	metadata := nodeMetadata[productNodeKey{
		clusterID: product.ClusterId,
		nodeName:  product.NodeName,
	}]

	resp := &productRes.ProductDetailResponse{
		ID:             product.ID,
		ClusterId:      product.ClusterId,
		ClusterName:    clusterNames[product.ClusterId],
		ProductType:    product.ProductType,
		Name:           product.Name,
		Description:    product.Description,
		Area:           product.Area,
		NodeName:       product.NodeName,
		NodeIP:         metadata.nodeIP,
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
		DriverVersion:  firstNonEmpty(product.DriverVersion, metadata.driverVersion),
		CUDAVersion:    firstNonEmpty(product.CUDAVersion, metadata.cudaVersion),
		SystemDisk:     product.SystemDisk,
		DataDisk:       product.DataDisk,
		Status:         product.Status,
		StorageClass:   product.StorageClass,
		StoragePriceGB: product.StoragePriceGB,
		MaxInstances:   product.MaxInstances,
		UsedCapacity:   product.UsedCapacity,
		Available:      product.AvailableCapacity(),
	}
	return resp, nil
}

func (s *ProductService) loadClusterNames(products []productModel.Product) (map[uint]string, error) {
	clusterIDs := make([]uint, 0)
	seen := make(map[uint]struct{})
	for _, p := range products {
		if _, ok := seen[p.ClusterId]; ok {
			continue
		}
		seen[p.ClusterId] = struct{}{}
		clusterIDs = append(clusterIDs, p.ClusterId)
	}

	if len(clusterIDs) == 0 {
		return map[uint]string{}, nil
	}

	var clusters []clusterModel.K8sCluster
	if err := global.GVA_DB.Where("id IN ?", clusterIDs).Find(&clusters).Error; err != nil {
		return nil, err
	}

	clusterNames := make(map[uint]string, len(clusters))
	for _, cluster := range clusters {
		clusterNames[cluster.ID] = cluster.Name
	}

	return clusterNames, nil
}

func (s *ProductService) loadProductNodeMetadata(ctx context.Context, products []productModel.Product) map[productNodeKey]productNodeMetadata {
	clusterIDs := make(map[uint]struct{})
	for _, product := range products {
		if product.ClusterId == 0 || product.NodeName == "" {
			continue
		}
		clusterIDs[product.ClusterId] = struct{}{}
	}

	if len(clusterIDs) == 0 {
		return map[productNodeKey]productNodeMetadata{}
	}

	nodeService := &NodeService{}
	result := make(map[productNodeKey]productNodeMetadata)
	for clusterID := range clusterIDs {
		nodes, err := nodeService.GetClusterNodesWithResources(ctx, clusterID, "")
		if err != nil {
			logx.Errorf("加载集群节点 IP 失败 clusterId=%d: %v", clusterID, err)
			continue
		}

		for _, node := range nodes {
			result[productNodeKey{
				clusterID: clusterID,
				nodeName:  node.NodeName,
			}] = productNodeMetadata{
				nodeIP:        node.InternalIP,
				driverVersion: node.DriverVersion,
				cudaVersion:   node.CUDAVersion,
			}
		}
	}

	return result
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
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
			HarborAddr:  c.HarborAddr,
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

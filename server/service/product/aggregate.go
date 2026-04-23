package product

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"gin-vue-admin/global"
	clusterModel "gin-vue-admin/model/cluster"
	"gin-vue-admin/model/consts"
	productModel "gin-vue-admin/model/product"
	productRes "gin-vue-admin/model/product/response"
	redisUtils "gin-vue-admin/utils/redis"

	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AggregateProductFilters struct {
	ClusterId       uint
	Area            string
	GPUModel        string
	CPUModel        string
	GPUResourceType string
	VGPUNumber      int64
	VGPUMemory      int64
	VGPUCores       int64
}

type aggregateProductKey struct {
	ProductType   int
	ClusterID     uint
	Area          string
	CPUModel      string
	CPU           int64
	Memory        int64
	GPUModel      string
	GPUCount      int64
	GPUMemory     int64
	VGPUNumber    int64
	VGPUMemory    int64
	VGPUCores     int64
	DriverVersion string
	CUDAVersion   string
	PriceHourly   float64
	PriceDaily    float64
	PriceWeekly   float64
	PriceMonthly  float64
}

type clusterNodeState struct {
	Schedulable   bool
	DriverVersion string
	CUDAVersion   string
}

type candidateProduct struct {
	ProductID uint
	NodeName  string
	Available int64
	ClusterID uint
	Product   productModel.Product
}

type PlannedReservation struct {
	ProductID uint
	NodeName  string
	Count     int64
}

type AllocationPlan struct {
	TemplateProduct  *productModel.Product
	ScheduleStrategy string
	InstanceCount    int64
	AllowedNodes     []string
	Reservations     []PlannedReservation
	StrictSpread     bool
}

func (s *ProductService) GetAggregateProductList(
	ctx context.Context,
	page, pageSize, productType int,
	filters AggregateProductFilters,
) (*productRes.AggregateProductListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var products []productModel.Product
	db := s.buildProductListDB(ctx, productType, filters)
	if err := db.Order("sort_order ASC, id ASC").Find(&products).Error; err != nil {
		return nil, err
	}
	if err := productModel.LoadPriceItemsForProducts(ctx, global.GVA_DB, products); err != nil {
		return nil, err
	}

	clusterNames, err := s.loadClusterNames(ctx, products)
	if err != nil {
		return nil, err
	}
	nodeStates := s.loadClusterNodeStates(products)

	grouped := make(map[aggregateProductKey]*productRes.AggregateProductResponse)
	for _, p := range products {
		key := buildAggregateProductKey(&p)
		nodeState := nodeStates[p.ClusterId][p.NodeName]
		resolvedDriverVersion := firstNonEmpty(p.DriverVersion, nodeState.DriverVersion)
		resolvedCUDAVersion := firstNonEmpty(p.CUDAVersion, nodeState.CUDAVersion)
		item, exists := grouped[key]
		if !exists {
			item = &productRes.AggregateProductResponse{
				ID:                p.ID,
				TemplateProductId: p.ID,
				ProductType:       p.ProductType,
				Name:              p.Name,
				Description:       p.Description,
				Area:              p.Area,
				ClusterId:         p.ClusterId,
				ClusterName:       clusterNames[p.ClusterId],
				CPUModel:          p.CPUModel,
				CPU:               p.CPU,
				Memory:            p.Memory,
				GPUModel:          p.GPUModel,
				GPUCount:          p.GPUCount,
				GPUMemory:         p.GPUMemory,
				VGPUNumber:        p.VGPUNumber,
				VGPUMemory:        p.VGPUMemory,
				VGPUCores:         p.VGPUCores,
				DriverVersion:     resolvedDriverVersion,
				CUDAVersion:       resolvedCUDAVersion,
				PriceHourly:       p.PriceHourly,
				PriceDaily:        p.PriceDaily,
				PriceWeekly:       p.PriceWeekly,
				PriceMonthly:      p.PriceMonthly,
			}
			grouped[key] = item
		}

		_, ok := nodeStates[p.ClusterId][p.NodeName]
		if !ok || !nodeState.Schedulable {
			continue
		}

		available := p.AvailableCapacity()
		if available <= 0 {
			continue
		}

		item.Available += available
		item.TotalSlots += available
		item.BalancedMax += available
		item.StrictMax++
		item.PhysicalNodeCount++
	}

	list := make([]productRes.AggregateProductResponse, 0, len(grouped))
	for _, item := range grouped {
		list = append(list, *item)
	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].Area != list[j].Area {
			return list[i].Area < list[j].Area
		}
		if list[i].ClusterName != list[j].ClusterName {
			return list[i].ClusterName < list[j].ClusterName
		}
		if list[i].GPUModel != list[j].GPUModel {
			return list[i].GPUModel < list[j].GPUModel
		}
		if list[i].CPU != list[j].CPU {
			return list[i].CPU < list[j].CPU
		}
		return list[i].TemplateProductId < list[j].TemplateProductId
	})

	total := int64(len(list))
	offset := (page - 1) * pageSize
	if offset > len(list) {
		offset = len(list)
	}
	end := offset + pageSize
	if end > len(list) {
		end = len(list)
	}

	return &productRes.AggregateProductListResponse{
		List:  list[offset:end],
		Total: total,
	}, nil
}

func (s *ProductService) PlanAllocations(
	ctx context.Context,
	templateProductID uint,
	instanceCount int64,
	scheduleStrategy string,
) (*AllocationPlan, error) {
	templateProduct, err := s.getEnabledProductModel(ctx, templateProductID)
	if err != nil {
		return nil, err
	}

	if instanceCount <= 0 {
		instanceCount = 1
	}
	if scheduleStrategy != consts.ScheduleStrategyStrict {
		scheduleStrategy = consts.ScheduleStrategyBalanced
	}

	candidates, err := s.loadCandidateProducts(ctx, templateProduct)
	if err != nil {
		return nil, err
	}

	reservations, allowedNodes, strictSpread, err := planReservations(candidates, instanceCount, scheduleStrategy)
	if err != nil {
		return nil, err
	}

	return &AllocationPlan{
		TemplateProduct:  templateProduct,
		ScheduleStrategy: scheduleStrategy,
		InstanceCount:    instanceCount,
		AllowedNodes:     allowedNodes,
		Reservations:     reservations,
		StrictSpread:     strictSpread,
	}, nil
}

func (s *ProductService) SaveResourceAllocations(ctx context.Context, allocations []productModel.ResourceAllocation) error {
	if len(allocations) == 0 {
		return nil
	}
	return global.GVA_DB.WithContext(ctx).Create(&allocations).Error
}

func (s *ProductService) DeleteResourceAllocations(ctx context.Context, instanceType string, instanceID uint) error {
	return global.GVA_DB.WithContext(ctx).
		Where("instance_type = ? AND instance_id = ?", instanceType, instanceID).
		Delete(&productModel.ResourceAllocation{}).Error
}

func (s *ProductService) ReleaseResourceAllocations(ctx context.Context, instanceType string, instanceID uint) (bool, error) {
	lock, err := redisUtils.AcquireResourceLock(ctx, fmt.Sprintf("resource-allocation-release:%s", instanceType), instanceID)
	if err != nil {
		return false, err
	}
	defer lock.Unlock(ctx)

	released := false
	err = global.GVA_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var allocations []productModel.ResourceAllocation
		if err := tx.
			Where("instance_type = ? AND instance_id = ?", instanceType, instanceID).
			Order("id ASC").
			Find(&allocations).Error; err != nil {
			return err
		}
		if len(allocations) == 0 {
			return nil
		}

		counts := make(map[uint]int64)
		ids := make([]uint, 0, len(allocations))
		for _, item := range allocations {
			counts[item.ProductID] += item.ReservedCount
			ids = append(ids, item.ID)
		}

		productIDs := make([]uint, 0, len(counts))
		for productID := range counts {
			productIDs = append(productIDs, productID)
		}
		sort.Slice(productIDs, func(i, j int) bool { return productIDs[i] < productIDs[j] })
		for _, productID := range productIDs {
			result := tx.Model(&productModel.Product{}).
				Where("id = ? AND used_capacity >= ?", productID, counts[productID]).
				Update("used_capacity", gorm.Expr("used_capacity - ?", counts[productID]))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New("释放失败：产品不存在或可释放数量不足")
			}
		}

		if err := tx.Where("id IN ?", ids).Delete(&productModel.ResourceAllocation{}).Error; err != nil {
			return err
		}

		released = true
		return nil
	})
	if err != nil {
		return false, err
	}
	return released, nil
}

func (s *ProductService) buildProductListDB(ctx context.Context, productType int, filters AggregateProductFilters) *gorm.DB {
	db := global.GVA_DB.WithContext(ctx).
		Model(&productModel.Product{}).
		Where("status = ?", productModel.ProductStatusEnabled)

	if productType > 0 {
		db = db.Where("product_type = ?", productType)
	}
	if filters.ClusterId > 0 {
		db = db.Where("cluster_id = ?", filters.ClusterId)
	}
	if filters.Area != "" {
		db = db.Where("area = ?", filters.Area)
	}
	if filters.GPUModel != "" {
		db = db.Where("gpu_model = ?", filters.GPUModel)
	}
	if filters.CPUModel != "" {
		db = db.Where("cpu_model = ?", filters.CPUModel)
	}
	if filters.GPUResourceType == "gpu" {
		db = db.Where("COALESCE(v_gpu_number, 0) = 0 AND COALESCE(v_gpu_memory, 0) = 0 AND COALESCE(v_gpu_cores, 0) = 0")
	}
	if filters.GPUResourceType == "vgpu" {
		db = db.Where("(COALESCE(v_gpu_number, 0) > 0 OR COALESCE(v_gpu_memory, 0) > 0 OR COALESCE(v_gpu_cores, 0) > 0)")
	}
	if filters.VGPUNumber > 0 {
		db = db.Where("v_gpu_number = ?", filters.VGPUNumber)
	}
	if filters.VGPUMemory > 0 {
		db = db.Where("v_gpu_memory = ?", filters.VGPUMemory)
	}
	if filters.VGPUCores > 0 {
		db = db.Where("v_gpu_cores = ?", filters.VGPUCores)
	}
	return db
}

func (s *ProductService) getEnabledProductModel(ctx context.Context, productID uint) (*productModel.Product, error) {
	var product productModel.Product
	if err := global.GVA_DB.WithContext(ctx).
		Where("id = ? AND status = ?", productID, productModel.ProductStatusEnabled).
		First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("产品不存在或已下架")
		}
		return nil, err
	}
	if err := productModel.LoadPriceItems(ctx, global.GVA_DB, &product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *ProductService) loadClusterNames(ctx context.Context, products []productModel.Product) (map[uint]string, error) {
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
	if err := global.GVA_DB.WithContext(ctx).Where("id IN ?", clusterIDs).Find(&clusters).Error; err != nil {
		return nil, err
	}

	result := make(map[uint]string, len(clusters))
	for _, cluster := range clusters {
		result[cluster.ID] = cluster.Name
	}
	return result, nil
}

func (s *ProductService) loadCandidateProducts(ctx context.Context, templateProduct *productModel.Product) ([]candidateProduct, error) {
	var products []productModel.Product
	if err := s.matchTemplateProductsQuery(ctx, templateProduct).Order("id ASC").Find(&products).Error; err != nil {
		return nil, err
	}
	if err := productModel.LoadPriceItemsForProducts(ctx, global.GVA_DB, products); err != nil {
		return nil, err
	}

	nodeStates := s.loadClusterNodeStates(products)
	candidates := make([]candidateProduct, 0, len(products))
	for _, product := range products {
		if !productModel.HasSameComputePrices(templateProduct, &product) {
			continue
		}
		if !nodeStates[product.ClusterId][product.NodeName].Schedulable {
			continue
		}
		available := product.AvailableCapacity()
		if available <= 0 {
			continue
		}
		candidates = append(candidates, candidateProduct{
			ProductID: product.ID,
			NodeName:  product.NodeName,
			Available: available,
			ClusterID: product.ClusterId,
			Product:   product,
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Available != candidates[j].Available {
			return candidates[i].Available > candidates[j].Available
		}
		if candidates[i].NodeName != candidates[j].NodeName {
			return candidates[i].NodeName < candidates[j].NodeName
		}
		return candidates[i].ProductID < candidates[j].ProductID
	})
	return candidates, nil
}

func (s *ProductService) matchTemplateProductsQuery(ctx context.Context, template *productModel.Product) *gorm.DB {
	return global.GVA_DB.WithContext(ctx).
		Model(&productModel.Product{}).
		Where("status = ?", productModel.ProductStatusEnabled).
		Where("product_type = ?", template.ProductType).
		Where("cluster_id = ?", template.ClusterId).
		Where("cpu_model = ? AND cpu = ? AND memory = ?", template.CPUModel, template.CPU, template.Memory).
		Where("gpu_model = ? AND gpu_count = ? AND gpu_memory = ?", template.GPUModel, template.GPUCount, template.GPUMemory).
		Where("v_gpu_number = ? AND v_gpu_memory = ? AND v_gpu_cores = ?", template.VGPUNumber, template.VGPUMemory, template.VGPUCores).
		Where("driver_version = ? AND cuda_version = ?", template.DriverVersion, template.CUDAVersion)
}

func (s *ProductService) loadClusterNodeStates(products []productModel.Product) map[uint]map[string]clusterNodeState {
	result := make(map[uint]map[string]clusterNodeState)
	clusterIDs := make(map[uint]struct{})
	for _, p := range products {
		clusterIDs[p.ClusterId] = struct{}{}
	}

	for clusterID := range clusterIDs {
		result[clusterID] = make(map[string]clusterNodeState)
		cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(clusterID)
		if cluster == nil || cluster.ClientSet == nil {
			continue
		}
		nodeList, err := cluster.ClientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			continue
		}
		for _, node := range nodeList.Items {
			driverVersion, cudaVersion := parseNodeVersionInfo(&node)
			result[clusterID][node.Name] = clusterNodeState{
				Schedulable:   !node.Spec.Unschedulable && isNodeReady(&node),
				DriverVersion: driverVersion,
				CUDAVersion:   cudaVersion,
			}
		}
	}

	return result
}

func isNodeReady(node *corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}

func parseNodeVersionInfo(node *corev1.Node) (string, string) {
	driverVersion := resolveNodeMetadataValue(node,
		[]string{
			"nvidia.com/cuda.driver-version.full",
			"nvidia.com/cuda.driver.version",
			"nvidia.com/cuda.driver-version",
			"nvidia.com/gpu.driver.version",
			"nvidia.com/gpu.driver-version",
			"nvidia.com/driver.version",
			"gpu.nvidia.com/driver-version",
		},
		[][]string{
			{"nvidia.com/cuda.driver-version.major", "nvidia.com/cuda.driver-version.minor", "nvidia.com/cuda.driver-version.revision"},
			{"nvidia.com/cuda.driver.major", "nvidia.com/cuda.driver.minor", "nvidia.com/cuda.driver.rev"},
			{"nvidia.com/cuda.driver.major", "nvidia.com/cuda.driver.minor", "nvidia.com/cuda.driver.patch"},
		},
	)

	cudaVersion := resolveNodeMetadataValue(node,
		[]string{
			"nvidia.com/cuda.runtime-version.full",
			"nvidia.com/cuda.runtime.version",
			"nvidia.com/cuda.runtime-version",
			"nvidia.com/cuda.version",
			"nvidia.com/cuda-version",
		},
		[][]string{
			{"nvidia.com/cuda.runtime-version.major", "nvidia.com/cuda.runtime-version.minor"},
			{"nvidia.com/cuda.runtime.major", "nvidia.com/cuda.runtime.minor"},
			{"nvidia.com/cuda.major", "nvidia.com/cuda.minor"},
		},
	)

	return driverVersion, cudaVersion
}

func resolveNodeMetadataValue(node *corev1.Node, fullKeys []string, partKeys [][]string) string {
	for _, key := range fullKeys {
		if value := lookupNodeMetadataValue(node, key); value != "" {
			return value
		}
	}

	for _, keys := range partKeys {
		parts := make([]string, 0, len(keys))
		missing := false
		for _, key := range keys {
			value := lookupNodeMetadataValue(node, key)
			if value == "" {
				missing = true
				break
			}
			parts = append(parts, value)
		}
		if !missing && len(parts) > 0 {
			return strings.Join(parts, ".")
		}
	}

	return ""
}

func lookupNodeMetadataValue(node *corev1.Node, key string) string {
	if value := strings.TrimSpace(node.Labels[key]); value != "" {
		return value
	}
	return strings.TrimSpace(node.Annotations[key])
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func buildAggregateProductKey(p *productModel.Product) aggregateProductKey {
	return aggregateProductKey{
		ProductType:   p.ProductType,
		ClusterID:     p.ClusterId,
		Area:          p.Area,
		CPUModel:      p.CPUModel,
		CPU:           p.CPU,
		Memory:        p.Memory,
		GPUModel:      p.GPUModel,
		GPUCount:      p.GPUCount,
		GPUMemory:     p.GPUMemory,
		VGPUNumber:    p.VGPUNumber,
		VGPUMemory:    p.VGPUMemory,
		VGPUCores:     p.VGPUCores,
		DriverVersion: p.DriverVersion,
		CUDAVersion:   p.CUDAVersion,
		PriceHourly:   p.PriceHourly,
		PriceDaily:    p.PriceDaily,
		PriceWeekly:   p.PriceWeekly,
		PriceMonthly:  p.PriceMonthly,
	}
}

func planReservations(candidates []candidateProduct, instanceCount int64, scheduleStrategy string) ([]PlannedReservation, []string, bool, error) {
	if instanceCount <= 0 {
		instanceCount = 1
	}

	if scheduleStrategy == consts.ScheduleStrategyStrict {
		nodeSeen := make(map[string]struct{})
		reservations := make([]PlannedReservation, 0, instanceCount)
		allowedNodes := make([]string, 0, instanceCount)
		for _, candidate := range candidates {
			if _, exists := nodeSeen[candidate.NodeName]; exists {
				continue
			}
			nodeSeen[candidate.NodeName] = struct{}{}
			reservations = append(reservations, PlannedReservation{
				ProductID: candidate.ProductID,
				NodeName:  candidate.NodeName,
				Count:     1,
			})
			allowedNodes = append(allowedNodes, candidate.NodeName)
			if int64(len(reservations)) == instanceCount {
				return reservations, allowedNodes, true, nil
			}
		}
		return nil, nil, false, ErrCapacityInsufficient
	}

	type mutableCandidate struct {
		ProductID uint
		NodeName  string
		Available int64
	}
	mutables := make([]mutableCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		mutables = append(mutables, mutableCandidate{
			ProductID: candidate.ProductID,
			NodeName:  candidate.NodeName,
			Available: candidate.Available,
		})
	}

	selectedCounts := make(map[uint]int64)
	allowedNodes := make([]string, 0, len(mutables))
	nodeSeen := make(map[string]struct{})
	remaining := instanceCount

	for idx := range mutables {
		if remaining == 0 {
			break
		}
		if mutables[idx].Available <= 0 {
			continue
		}
		if _, exists := nodeSeen[mutables[idx].NodeName]; exists {
			continue
		}
		nodeSeen[mutables[idx].NodeName] = struct{}{}
		allowedNodes = append(allowedNodes, mutables[idx].NodeName)
		selectedCounts[mutables[idx].ProductID]++
		mutables[idx].Available--
		remaining--
	}

	for remaining > 0 {
		best := -1
		for idx := range mutables {
			if mutables[idx].Available <= 0 {
				continue
			}
			if best == -1 || mutables[idx].Available > mutables[best].Available {
				best = idx
			}
		}
		if best == -1 {
			return nil, nil, false, ErrCapacityInsufficient
		}
		selectedCounts[mutables[best].ProductID]++
		mutables[best].Available--
		remaining--
	}

	reservations := make([]PlannedReservation, 0, len(selectedCounts))
	for _, candidate := range candidates {
		count := selectedCounts[candidate.ProductID]
		if count <= 0 {
			continue
		}
		reservations = append(reservations, PlannedReservation{
			ProductID: candidate.ProductID,
			NodeName:  candidate.NodeName,
			Count:     count,
		})
	}

	return reservations, allowedNodes, false, nil
}

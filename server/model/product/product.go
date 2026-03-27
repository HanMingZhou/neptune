package product

import (
	"time"

	"gorm.io/gorm"
)

// 产品状态常量
const (
	ProductStatusEnabled  = 1 // 上架
	ProductStatusDisabled = 0 // 下架
)

// 付费类型常量
const (
	ChargeTypeHourly  = 1 // 按量付费（每小时）
	ChargeTypeDaily   = 2 // 包日
	ChargeTypeWeekly  = 3 // 包周
	ChargeTypeMonthly = 4 // 包月
)

var ChargeTypeToString = map[int]string{
	ChargeTypeHourly:  "按量付费",
	ChargeTypeDaily:   "包日",
	ChargeTypeWeekly:  "包周",
	ChargeTypeMonthly: "包月",
}

// 产品类型常量
const (
	ProductTypeCompute = 1 // 计算资源
	ProductTypeStorage = 2 // 文件存储
)

// Product 产品表 - 对应集群中每个Node节点的配置和价格信息
// 一个产品 = 一个K8s Node节点（计算资源）或一个 StorageClass（文件存储）
// 产品类型互斥：GPU / vGPU / CPU-only 三选一
// 库存统一使用 MaxInstances（自动计算）和 UsedCapacity 管理
type Product struct {
	ID            uint           `gorm:"primarykey"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index;uniqueIndex:idx_product_unique"`
	ProductType   int            `json:"productType" gorm:"column:product_type;default:1;comment:产品类型(1-计算资源 2-文件存储);uniqueIndex:idx_product_unique"`
	Name          string         `json:"name" gorm:"column:name;size:100;comment:产品名称;uniqueIndex:idx_product_unique"`
	Description   string         `json:"description" gorm:"column:description;size:500;comment:产品描述"`
	ClusterId     uint           `json:"clusterId" gorm:"column:cluster_id;comment:集群ID;index;uniqueIndex:idx_product_unique"`
	Area          string         `json:"area" gorm:"column:area;size:50;comment:地域区域;index"`
	NodeName      string         `json:"nodeName" gorm:"column:node_name;size:100;comment:K8s节点名称;uniqueIndex:idx_product_unique"`
	NodeType      string         `json:"nodeType" gorm:"column:node_type;size:50;comment:节点类型"`
	CPUModel      string         `json:"cpuModel" gorm:"column:cpu_model;size:100;comment:CPU型号"`
	CPU           int64          `json:"cpu" gorm:"column:cpu;size:20;comment:每个实例的CPU核数"`
	Memory        int64          `json:"memory" gorm:"column:memory;size:20;comment:每个实例的内存大小(GB)"`
	GPUModel      string         `json:"gpuModel" gorm:"column:gpu_model;size:50;comment:GPU型号(GPU产品必填,与vGPU/CPU-only互斥)"`
	GPUCount      int64          `json:"gpuCount" gorm:"column:gpu_count;comment:每个实例的GPU卡数(GPU产品必填);default:0"`
	GPUMemory     int64          `json:"gpuMemory" gorm:"column:gpu_memory;size:20;comment:单卡GPU显存(GB)"`
	VGPUNumber    int64          `json:"vGpuNumber" gorm:"column:v_gpu_number;comment:每个实例的vGPU数量(vGPU产品,与GPU/CPU-only互斥);default:0"`
	VGPUMemory    int64          `json:"vGpuMemory" gorm:"column:v_gpu_memory;comment:每个实例的vGPU显存(vGPU产品);default:0"`
	VGPUCores     int64          `json:"vGpuCores" gorm:"column:v_gpu_cores;comment:每个实例的vGPU核心数(vGPU产品);default:0"`
	PriceHourly   float64        `json:"priceHourly" gorm:"column:price_hourly;type:decimal(20,6);comment:每小时价格(单个实例)"`
	PriceDaily    float64        `json:"priceDaily" gorm:"column:price_daily;type:decimal(20,6);comment:包日价格(单个实例)"`
	PriceWeekly   float64        `json:"priceWeekly" gorm:"column:price_weekly;type:decimal(20,6);comment:包周价格(单个实例)"`
	PriceMonthly  float64        `json:"priceMonthly" gorm:"column:price_monthly;type:decimal(20,6);comment:包月价格(单个实例)"`
	DriverVersion string         `json:"driverVersion" gorm:"column:driver_version;size:50;comment:显卡驱动版本"`
	CUDAVersion   string         `json:"cudaVersion" gorm:"column:cuda_version;size:50;comment:CUDA版本"`
	SystemDisk    int64          `json:"systemDisk" gorm:"column:system_disk;comment:系统盘大小(GB)"`
	DataDisk      int64          `json:"dataDisk" gorm:"column:data_disk;comment:数据盘大小(GB)"`
	Status        int            `json:"status" gorm:"column:status;default:1;comment:状态(1-上架 0-下架)"`
	SortOrder     int            `json:"sortOrder" gorm:"column:sort_order;default:0;comment:排序"`
	MaxInstances  int64          `json:"maxInstances" gorm:"column:max_instances;comment:最大实例数(自动计算:节点资源/产品规格);default:0"`
	UsedCapacity  int64          `json:"usedCapacity" gorm:"column:used_capacity;comment:已使用实例数;default:0"`
	// 文件存储专用字段
	StorageClass   string  `json:"storageClass" gorm:"column:storage_class;size:100;comment:K8s StorageClass名称;uniqueIndex:idx_product_unique"`
	StoragePriceGB float64 `json:"storagePriceGb" gorm:"column:storage_price_gb;type:decimal(20,6);comment:每GB每日价格"`
	Version        int64   `json:"version" gorm:"column:version;default:0;comment:乐观锁版本号"`
}

func (Product) TableName() string {
	return "products"
}

// IsCPUOnly 判断是否为CPU-only产品（无GPU、无vGPU）
func (p *Product) IsCPUOnly() bool {
	return p.GPUCount == 0 && p.VGPUNumber == 0 && p.VGPUMemory == 0 && p.VGPUCores == 0
}

// IsVGPU 判断是否为vGPU产品
func (p *Product) IsVGPU() bool {
	return p.VGPUNumber > 0 || p.VGPUMemory > 0 || p.VGPUCores > 0
}

// IsGPUOnly 判断是否为纯GPU产品（有GPU卡，无vGPU）
func (p *Product) IsGPUOnly() bool {
	return p.GPUCount > 0 && !p.IsVGPU()
}

// AvailableCapacity 计算可用容量(CPU-ONLY GPU-ONLY vGPU)
func (p *Product) AvailableCapacity() int64 {
	available := p.MaxInstances - p.UsedCapacity
	if available < 0 {
		return 0
	}
	return available
}

// GetPrice 根据付费类型获取价格
func (p *Product) GetPrice(chargeType int64) float64 {
	switch chargeType {
	case ChargeTypeHourly:
		return p.PriceHourly
	case ChargeTypeDaily:
		return p.PriceDaily
	case ChargeTypeWeekly:
		return p.PriceWeekly
	case ChargeTypeMonthly:
		return p.PriceMonthly
	default:
		return p.PriceHourly
	}
}

// GetPayTypeName 获取付费类型名称
func GetPayTypeName(chargeType int64) string {
	switch chargeType {
	case ChargeTypeHourly:
		return "按量付费"
	case ChargeTypeDaily:
		return "包日"
	case ChargeTypeWeekly:
		return "包周"
	case ChargeTypeMonthly:
		return "包月"
	default:
		return "按量付费"
	}
}

// NodeInfo K8s节点信息（用于从K8s获取后返回给前端）
type NodeInfo struct {
	Area        string `json:"area"`
	NodeName    string `json:"nodeName"`
	CPUModel    string `json:"cpuModel"`
	CPU         string `json:"cpu"`
	Memory      string `json:"memory"`
	GPUModel    string `json:"gpuModel"`
	GPUCount    int    `json:"gpuCount"`
	GPUMemory   string `json:"gpuMemory"`
	Schedulable bool   `json:"schedulable"`
}

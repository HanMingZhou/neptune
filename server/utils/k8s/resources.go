package helper

import (
	"fmt"
	"gin-vue-admin/model/consts"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// ProductSpec 产品规格（传递给各 builder，避免 builder 依赖 DB）
// 统一用于 Training / Notebook / Inference 的资源构建
type ProductSpec struct {
	CPU        int64  // CPU 核数
	Memory     int64  // 内存大小(GB)
	GPUModel   string // GPU 型号
	GPUCount   int64  // GPU 卡数（GPU 产品）
	VGPUNumber int64  // vGPU 数量（vGPU 产品）
	VGPUMemory int64  // vGPU 显存（vGPU 产品，单位由调度器决定）
	VGPUCores  int64  // vGPU 核心数（vGPU 产品）
}

// IsGPU 是否为 GPU 产品
func (p *ProductSpec) IsGPU() bool {
	return p.GPUCount > 0 && p.VGPUNumber == 0
}

// IsVGPU 是否为 vGPU 产品
func (p *ProductSpec) IsVGPU() bool {
	return p.VGPUNumber > 0 || p.VGPUMemory > 0 || p.VGPUCores > 0
}

// HasGPU 是否有任何 GPU 资源（用于 Tolerations 判断）
func (p *ProductSpec) HasGPU() bool {
	return p.IsGPU() || p.IsVGPU()
}

// BuildResources 根据产品规格构建 K8s 资源请求（统一入口）
// GPU 产品 → nvidia.com/gpu
// vGPU 产品 → volcano.sh/vgpu-number + volcano.sh/vgpu-memory + volcano.sh/vgpu-cores
// CPU-only → 不设置 GPU 资源
func BuildResources(product *ProductSpec) corev1.ResourceRequirements {
	resources := corev1.ResourceRequirements{
		Limits:   corev1.ResourceList{},
		Requests: corev1.ResourceList{},
	}

	if product.CPU > 0 {
		cpuQty := resource.MustParse(fmt.Sprintf("%d", product.CPU))
		resources.Limits[corev1.ResourceCPU] = cpuQty
		resources.Requests[corev1.ResourceCPU] = cpuQty
	}

	if product.Memory > 0 {
		memQty := resource.MustParse(fmt.Sprintf("%dGi", product.Memory))
		resources.Limits[corev1.ResourceMemory] = memQty
		resources.Requests[corev1.ResourceMemory] = memQty
	}

	switch {
	case product.IsGPU():
		gpuQty := resource.MustParse(fmt.Sprintf("%d", product.GPUCount))
		resources.Limits[consts.NvidiaGPUType] = gpuQty
		resources.Requests[consts.NvidiaGPUType] = gpuQty
	case product.IsVGPU():
		if product.VGPUNumber > 0 {
			qty := resource.MustParse(fmt.Sprintf("%d", product.VGPUNumber))
			resources.Limits[consts.VolcanoVGPUNumber] = qty
			resources.Requests[consts.VolcanoVGPUNumber] = qty
		}
		if product.VGPUMemory > 0 {
			qty := resource.MustParse(fmt.Sprintf("%d", product.VGPUMemory))
			resources.Limits[consts.VolcanoVGPUMemory] = qty
			resources.Requests[consts.VolcanoVGPUMemory] = qty
		}
		if product.VGPUCores > 0 {
			qty := resource.MustParse(fmt.Sprintf("%d", product.VGPUCores))
			resources.Limits[consts.VolcanoVGPUCores] = qty
			resources.Requests[consts.VolcanoVGPUCores] = qty
		}
	}

	return resources
}

// BuildGPUTolerations 根据产品规格构建 GPU Tolerations
func BuildGPUTolerations(product *ProductSpec) []corev1.Toleration {
	if !product.HasGPU() {
		return nil
	}
	return []corev1.Toleration{
		{
			Key:      consts.NvidiaGPUType,
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoSchedule,
		},
	}
}

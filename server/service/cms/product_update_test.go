package cms

import (
	"testing"

	cmsReq "gin-vue-admin/model/cms/request"
	productModel "gin-vue-admin/model/product"
)

func TestBuildUpdatedComputeSpecSwitchesGPUToVGPU(t *testing.T) {
	existing := productModel.Product{
		ProductType: productModel.ProductTypeCompute,
		ClusterId:   1,
		Area:        "cn",
		NodeName:    "node-a",
		CPU:         16,
		Memory:      64,
		GPUModel:    "A100",
		GPUCount:    1,
		GPUMemory:   40,
	}

	spec := buildUpdatedComputeSpec(existing, &cmsReq.UpdateProductReq{
		ID:         1,
		VGPUNumber: 2,
		VGPUMemory: 10,
		VGPUCores:  25,
	}, false, true)

	if spec.GPUCount != 0 {
		t.Fatalf("expected gpu_count to be cleared, got %d", spec.GPUCount)
	}
	if spec.GPUMemory != 0 {
		t.Fatalf("expected gpu_memory to be cleared, got %d", spec.GPUMemory)
	}
	if spec.VGPUNumber != 2 || spec.VGPUMemory != 10 || spec.VGPUCores != 25 {
		t.Fatalf("expected vgpu spec to be applied, got %+v", spec)
	}
}

func TestBuildUpdatedComputeSpecSwitchesVGPUToGPU(t *testing.T) {
	existing := productModel.Product{
		ProductType: productModel.ProductTypeCompute,
		ClusterId:   1,
		Area:        "cn",
		NodeName:    "node-b",
		CPU:         16,
		Memory:      64,
		GPUModel:    "L20",
		VGPUNumber:  2,
		VGPUMemory:  12,
		VGPUCores:   30,
	}

	spec := buildUpdatedComputeSpec(existing, &cmsReq.UpdateProductReq{
		ID:        2,
		GPUCount:  1,
		GPUMemory: 24,
	}, true, false)

	if spec.VGPUNumber != 0 || spec.VGPUMemory != 0 || spec.VGPUCores != 0 {
		t.Fatalf("expected vgpu spec to be cleared, got %+v", spec)
	}
	if spec.GPUCount != 1 || spec.GPUMemory != 24 {
		t.Fatalf("expected gpu spec to be applied, got %+v", spec)
	}
}

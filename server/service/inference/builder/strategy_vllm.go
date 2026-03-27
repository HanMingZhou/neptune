package builder

import corev1 "k8s.io/api/core/v1"

// VLLMStrategy vLLM 框架环境变量策略
type VLLMStrategy struct{}

var _ FrameworkStrategy = (*VLLMStrategy)(nil)

// GetEnvVars 获取 vLLM 所需环境变量（NCCL 通信）
func (v *VLLMStrategy) GetEnvVars(spec *InferenceSpec) []corev1.EnvVar {
	if !spec.Product.HasGPU() {
		return []corev1.EnvVar{}
	}
	return []corev1.EnvVar{
		{Name: "NCCL_SOCKET_IFNAME", Value: "eth0"},
	}
}

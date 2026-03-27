package builder

import corev1 "k8s.io/api/core/v1"

// SGLangStrategy SGLang 框架环境变量策略
type SGLangStrategy struct{}

var _ FrameworkStrategy = (*SGLangStrategy)(nil)

// GetEnvVars 获取 SGLang 所需环境变量（NCCL 通信）
func (s *SGLangStrategy) GetEnvVars(spec *InferenceSpec) []corev1.EnvVar {
	if !spec.Product.HasGPU() {
		return []corev1.EnvVar{}
	}
	return []corev1.EnvVar{
		{Name: "NCCL_SOCKET_IFNAME", Value: "eth0"},
	}
}

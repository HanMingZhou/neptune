package builder

import corev1 "k8s.io/api/core/v1"

// NoopStrategy 空策略（未指定框架时使用）
type NoopStrategy struct{}

var _ FrameworkStrategy = (*NoopStrategy)(nil)

// GetEnvVars 返回空环境变量
func (n *NoopStrategy) GetEnvVars(spec *InferenceSpec) []corev1.EnvVar {
	return []corev1.EnvVar{}
}

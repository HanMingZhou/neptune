package builder

import (
	"fmt"
	"gin-vue-admin/model/consts"
	helper "gin-vue-admin/utils/k8s"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	vcbatch "volcano.sh/apis/pkg/apis/batch/v1alpha1"
)

// InferenceSpec 推理服务规格
type InferenceSpec struct {
	Name             string
	Namespace        string
	InstanceName     string
	ServiceID        uint   // 数据库记录 ID（用于 Pod 标签，Informer 识别）
	Framework        string // SGLANG/VLLM（仅用于分布式 NCCL 环境变量注入）
	Image            string
	ModelPvcName     string
	ModelMountPath   string // PVC 挂载路径（默认 /model）
	TensorParallel   int
	PipelineParallel int
	WorkerCount      int
	Product          helper.ProductSpec
	ServicePort      int
	Command          string               // 完整启动命令字符串
	Args             []string             // 启动参数（追加到命令末尾）
	UserVolumes      []corev1.Volume      // 用户自定义挂载 Volume
	UserVolumeMounts []corev1.VolumeMount // 用户自定义挂载 VolumeMount
	UserEnvVars      []corev1.EnvVar      // 用户自定义环境变量
}

// GetModelMountPath 获取 PVC 挂载路径（带默认值）
func (s *InferenceSpec) GetModelMountPath() string {
	if s.ModelMountPath != "" {
		return s.ModelMountPath
	}
	return "/model"
}

// InferenceBuilder 推理服务构建器接口
type InferenceBuilder interface {
	BuildVCJob(spec *InferenceSpec) (*vcbatch.Job, error)
	BuildDeployment(spec *InferenceSpec) (*appsv1.Deployment, error)
}

// FrameworkStrategy 框架环境变量策略（NCCL 等）
type FrameworkStrategy interface {
	GetEnvVars(spec *InferenceSpec) []corev1.EnvVar
}

// BaseInferenceBuilder 推理服务构建器基础实现
type BaseInferenceBuilder struct {
	Strategy FrameworkStrategy
}

// NewInferenceBuilder 根据框架创建构建器
func NewInferenceBuilder(framework string) InferenceBuilder {
	var strategy FrameworkStrategy
	switch framework {
	case consts.FrameworkSGLang:
		strategy = &SGLangStrategy{}
	case consts.FrameworkVLLM:
		strategy = &VLLMStrategy{}
	default:
		strategy = &NoopStrategy{}
	}
	return &BaseInferenceBuilder{Strategy: strategy}
}

// buildModelVolume 构建模型 PVC Volume
func buildModelVolume(pvcName, mountPath string) (corev1.Volume, corev1.VolumeMount) {
	volume := corev1.Volume{
		Name: "model-volume",
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: pvcName,
				ReadOnly:  true,
			},
		},
	}
	mount := corev1.VolumeMount{
		Name:      "model-volume",
		MountPath: mountPath,
		ReadOnly:  true,
	}
	return volume, mount
}

// buildSHMVolume 构建共享内存 Volume（NCCL 通信需要）
func buildSHMVolume(sizeGi int) (corev1.Volume, corev1.VolumeMount) {
	volume := corev1.Volume{
		Name: "shm",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{
				Medium:    corev1.StorageMediumMemory,
				SizeLimit: func() *resource.Quantity { q := resource.MustParse(fmt.Sprintf("%dGi", sizeGi)); return &q }(),
			},
		},
	}
	mount := corev1.VolumeMount{
		Name:      "shm",
		MountPath: "/dev/shm",
	}
	return volume, mount
}

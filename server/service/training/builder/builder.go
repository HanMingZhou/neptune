package builder

import (
	"fmt"
	"gin-vue-admin/model/consts"
	trainingReq "gin-vue-admin/model/training/request"
	helper "gin-vue-admin/utils/k8s"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	vcbatch "volcano.sh/apis/pkg/apis/batch/v1alpha1"
)

// JobBuilder 训练任务构建器接口
type JobBuilder interface {
	// Build 根据规格构建 Volcano Job
	Build(spec *trainingReq.TrainingJobSpec) (*vcbatch.Job, error)
}

// FrameworkStrategy 框架策略接口
// 不同的训练框架（PyTorch DDP, MPI, Standalone）实现不同的策略
type FrameworkStrategy interface {
	// BuildTasks 构建任务的 Tasks 列表
	BuildTasks(spec *trainingReq.TrainingJobSpec) ([]vcbatch.TaskSpec, error)
	// GetPlugins 获取 Volcano 插件配置
	GetPlugins() map[string][]string
	// GetMinAvailable 获取最小可用 Pod 数量
	GetMinAvailable(spec *trainingReq.TrainingJobSpec) int64
	// GetPolicies 获取生命周期策略
	GetPolicies() []vcbatch.LifecyclePolicy
}

// BaseBuilder 基础构建器，包含通用逻辑
type BaseBuilder struct {
	Strategy FrameworkStrategy
}

type taskTemplateOptions struct {
	containerName string
	image         string
	command       []string
	args          []string
	resources     corev1.ResourceRequirements
	volumeMounts  []corev1.VolumeMount
	envs          []corev1.EnvVar
	volumes       []corev1.Volume
	tolerations   []corev1.Toleration
	readiness     *corev1.Probe
}

// NewJobBuilder 根据框架类型创建对应的 Builder
func NewJobBuilder(framework string) JobBuilder {
	var strategy FrameworkStrategy
	switch framework {
	case consts.FrameworkPyTorchDDP:
		strategy = &PyTorchDDPStrategy{}
	case consts.FrameworkMPI:
		strategy = &MPIStrategy{}
	case consts.FrameworkStandalone:
		strategy = &StandaloneStrategy{}
	default:
		strategy = &StandaloneStrategy{}
	}
	return &BaseBuilder{Strategy: strategy}
}

// Build 构建 Volcano Job
func (b *BaseBuilder) Build(spec *trainingReq.TrainingJobSpec) (*vcbatch.Job, error) {
	// 1. 构建 Tasks
	tasks, err := b.Strategy.BuildTasks(spec)
	if err != nil {
		return nil, err
	}

	// 2. 添加共享内存卷（如果启用）
	if spec.UseSHM {
		shmVolume, shmMount := buildSHMVolume(spec.SHMSize)
		spec.Volumes = append(spec.Volumes, shmVolume)
		// 为每个 Task 的容器添加 SHM 挂载
		for i := range tasks {
			for j := range tasks[i].Template.Spec.Containers {
				tasks[i].Template.Spec.Containers[j].VolumeMounts = append(
					tasks[i].Template.Spec.Containers[j].VolumeMounts,
					shmMount,
				)
			}
			tasks[i].Template.Spec.Volumes = append(tasks[i].Template.Spec.Volumes, shmVolume)
		}
	}

	// 3. 构建 Volcano Job
	job := &vcbatch.Job{
		Spec: vcbatch.JobSpec{
			MinAvailable:  int32(b.Strategy.GetMinAvailable(spec)),
			SchedulerName: consts.VolcanoScheduler,
			Queue:         consts.DefalutQueue,
			Tasks:         tasks,
			Plugins:       b.Strategy.GetPlugins(),
			Policies:      b.Strategy.GetPolicies(),
			MaxRetry:      int32(spec.MaxRetry),
		},
	}

	// 4. 设置元数据
	job.Name = spec.Name
	job.Namespace = spec.Namespace
	if spec.Labels != nil {
		job.Labels = spec.Labels
	}
	if job.Labels == nil {
		job.Labels = make(map[string]string)
	}
	job.Labels[consts.LabelVolcanoJob] = consts.TrainingInstance
	job.Labels[consts.LabelFramework] = spec.Framework

	if spec.Annotations != nil {
		job.Annotations = spec.Annotations
	}

	return job, nil
}

// buildSHMVolume 创建共享内存卷
func buildSHMVolume(shmSize int64) (corev1.Volume, corev1.VolumeMount) {
	if shmSize <= 0 {
		shmSize = 2 // 默认 2Gi
	}
	sizeLimit := resource.MustParse(fmt.Sprintf("%dGi", shmSize))

	volume := corev1.Volume{
		Name: "shm",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{
				Medium:    corev1.StorageMediumMemory,
				SizeLimit: &sizeLimit,
			},
		},
	}

	mount := corev1.VolumeMount{
		Name:      "shm",
		MountPath: "/dev/shm",
	}

	return volume, mount
}

// buildResources 根据产品规格构建 K8s 资源请求（委托给共享实现）
func buildResources(product *trainingReq.ProductSpec) corev1.ResourceRequirements {
	return helper.BuildResources(product)
}

// buildGPUTolerations 根据产品规格构建 GPU Tolerations（委托给共享实现）
func buildGPUTolerations(product *trainingReq.ProductSpec) []corev1.Toleration {
	return helper.BuildGPUTolerations(product)
}

func buildTaskTemplate(opts taskTemplateOptions) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{
					Name:           opts.containerName,
					Image:          opts.image,
					Command:        append([]string(nil), opts.command...),
					Args:           append([]string(nil), opts.args...),
					Resources:      opts.resources,
					VolumeMounts:   append([]corev1.VolumeMount(nil), opts.volumeMounts...),
					Env:            cloneEnvVars(opts.envs),
					ReadinessProbe: opts.readiness,
				},
			},
			Volumes:     append([]corev1.Volume(nil), opts.volumes...),
			Tolerations: append([]corev1.Toleration(nil), opts.tolerations...),
		},
	}
}

func cloneEnvVars(envs []corev1.EnvVar) []corev1.EnvVar {
	if len(envs) == 0 {
		return nil
	}
	cloned := make([]corev1.EnvVar, len(envs))
	copy(cloned, envs)
	return cloned
}

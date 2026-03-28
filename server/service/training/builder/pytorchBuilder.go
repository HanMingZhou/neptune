package builder

import (
	"gin-vue-admin/model/consts"
	trainingReq "gin-vue-admin/model/training/request"
	helper "gin-vue-admin/utils/k8s"

	corev1 "k8s.io/api/core/v1"
	vcbatch "volcano.sh/apis/pkg/apis/batch/v1alpha1"
	vcbus "volcano.sh/apis/pkg/apis/bus/v1alpha1"
)

// PyTorchDDPStrategy PyTorch DDP 分布式训练策略
type PyTorchDDPStrategy struct{}

var _ FrameworkStrategy = (*PyTorchDDPStrategy)(nil)

// BuildTasks 构建 PyTorch DDP 任务
// PyTorch DDP 需要一个 Master 和多个 Worker
func (s *PyTorchDDPStrategy) BuildTasks(spec *trainingReq.TrainingJobSpec) ([]vcbatch.TaskSpec, error) {
	if spec.WorkerCount < 1 {
		spec.WorkerCount = 1
	}

	// 计算实际 PyTorch Worker 数量（不含 Master 本身占用的 1 个节点配额）
	workerReplicas := spec.WorkerCount - 1
	if workerReplicas < 0 {
		workerReplicas = 0
	}

	// GPU Tolerations
	tolerations := buildGPUTolerations(&spec.Product)

	// 构建容器资源（所有节点使用相同的产品规格）
	nodeResources := buildResources(&spec.Product)
	baseEnvs := s.buildEnvVars(spec)
	affinity := helper.BuildSchedulingAffinity(spec.Name, spec.AllowedNodes, spec.StrictSpread)

	tasks := []vcbatch.TaskSpec{}

	// 1. Master Task
	masterTask := vcbatch.TaskSpec{
		Name:     consts.TaskSpecMaster,
		Replicas: 1,
		Template: buildTaskTemplate(taskTemplateOptions{
			containerName: spec.Name + "-master",
			image:         spec.Image,
			command:       spec.Command,
			args:          spec.Args,
			resources:     nodeResources,
			volumeMounts:  spec.VolumeMounts,
			envs:          baseEnvs,
			volumes:       spec.Volumes,
			tolerations:   tolerations,
			labels:        spec.Labels,
			affinity:      affinity,
		}),
	}
	tasks = append(tasks, masterTask)

	// 2. Worker Tasks（如果有）
	if workerReplicas > 0 {
		workerTask := vcbatch.TaskSpec{
			Name:     consts.TaskSpecWorker,
			Replicas: int32(workerReplicas),
			Template: buildTaskTemplate(taskTemplateOptions{
				containerName: spec.Name + "-worker",
				image:         spec.Image,
				command:       spec.Command,
				args:          spec.Args,
				resources:     nodeResources,
				volumeMounts:  spec.VolumeMounts,
				envs:          baseEnvs,
				volumes:       spec.Volumes,
				tolerations:   tolerations,
				labels:        spec.Labels,
				affinity:      affinity,
			}),
		}
		tasks = append(tasks, workerTask)
	}

	return tasks, nil
}

// GetPlugins 获取 PyTorch DDP 插件配置
func (s *PyTorchDDPStrategy) GetPlugins() map[string][]string {
	return map[string][]string{
		"pytorch": {"--master=master", "--worker=worker", "--port=29500"},
		"svc":     {}, // 启用 Service Discovery
		"env":     {}, // 自动注入环境变量
	}
}

// GetMinAvailable 获取最小可用数量
// DDP 模式要求所有节点（1 个 Master + WorkerCount-1 个 Worker）都可用，所以要求可用总数为 spec.WorkerCount
func (s *PyTorchDDPStrategy) GetMinAvailable(spec *trainingReq.TrainingJobSpec) int64 {
	return int64(spec.WorkerCount)
}

// GetPolicies 获取生命周期策略
// DDP 模式任意 Pod 失败则重启整个 Job
func (s *PyTorchDDPStrategy) GetPolicies() []vcbatch.LifecyclePolicy {
	return []vcbatch.LifecyclePolicy{
		{
			Event:  vcbus.PodFailedEvent,
			Action: vcbus.RestartJobAction,
		},
		{
			Event:  vcbus.PodEvictedEvent,
			Action: vcbus.RestartJobAction,
		},
	}
}

// buildEnvVars 构建环境变量
func (s *PyTorchDDPStrategy) buildEnvVars(spec *trainingReq.TrainingJobSpec) []corev1.EnvVar {
	envs := cloneEnvVars(spec.Envs)

	// 添加 PyTorch DDP 相关环境变量
	// Volcano PyTorch 插件会自动注入以下变量：
	// - MASTER_ADDR: Master 节点地址
	// - MASTER_PORT: Master 端口
	// - WORLD_SIZE: 总节点数
	// - RANK: 当前节点排名

	// 添加自定义环境变量
	envs = append(envs,
		corev1.EnvVar{
			Name:  "NCCL_DEBUG",
			Value: "INFO",
		},
	)

	return envs
}

package builder

import (
	"gin-vue-admin/model/consts"
	trainingReq "gin-vue-admin/model/training/request"

	vcbatch "volcano.sh/apis/pkg/apis/batch/v1alpha1"
	vcbus "volcano.sh/apis/pkg/apis/bus/v1alpha1"
)

// StandaloneStrategy 单机训练策略
type StandaloneStrategy struct{}

var _ FrameworkStrategy = (*StandaloneStrategy)(nil)

// BuildTasks 构建单机训练任务
// 单机模式只有一个 Task，一个副本
func (s *StandaloneStrategy) BuildTasks(spec *trainingReq.TrainingJobSpec) ([]vcbatch.TaskSpec, error) {
	// 构建资源配置
	resources := buildResources(&spec.Product)

	// GPU Tolerations
	tolerations := buildGPUTolerations(&spec.Product)

	task := vcbatch.TaskSpec{
		Name:     consts.TaskSpecWorker,
		Replicas: 1,
		Template: buildTaskTemplate(taskTemplateOptions{
			containerName: spec.Name + "-worker",
			image:         spec.Image,
			command:       spec.Command,
			args:          spec.Args,
			resources:     resources,
			volumeMounts:  spec.VolumeMounts,
			envs:          spec.Envs,
			volumes:       spec.Volumes,
			tolerations:   tolerations,
		}),
	}

	return []vcbatch.TaskSpec{task}, nil
}

// GetPlugins 获取插件配置
// 单机模式不需要特殊插件
func (s *StandaloneStrategy) GetPlugins() map[string][]string {
	return map[string][]string{
		"env": {}, // 自动注入环境变量
	}
}

// GetMinAvailable 获取最小可用数量
// 单机模式只需要 1 个 Pod
func (s *StandaloneStrategy) GetMinAvailable(spec *trainingReq.TrainingJobSpec) int64 {
	return 1
}

// GetPolicies 获取生命周期策略
// 单机模式 Pod 失败可以重启 Pod（而不是整个 Job）
func (s *StandaloneStrategy) GetPolicies() []vcbatch.LifecyclePolicy {
	return []vcbatch.LifecyclePolicy{
		{
			Event:  vcbus.PodFailedEvent,
			Action: vcbus.RestartTaskAction,
		},
	}
}

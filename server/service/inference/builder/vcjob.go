package builder

import (
	"fmt"
	"gin-vue-admin/model/consts"
	helper "gin-vue-admin/utils/k8s"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	vcbatch "volcano.sh/apis/pkg/apis/batch/v1alpha1"
	vcbus "volcano.sh/apis/pkg/apis/bus/v1alpha1"
)

const (
	ncclPort   = 29500 // NCCL rendezvous 端口
	masterPort = 29500 // 与 ncclPort 一致，用于 MASTER_PORT 环境变量
)

// BuildVCJob 构建分布式推理 VCJob
//
// 分布式推理架构（基于 NCCL 原生多节点，非 Ray）：
//   - 所有节点（head + worker）运行相同的推理框架进程
//   - 通过环境变量 MASTER_ADDR / MASTER_PORT / WORLD_SIZE / NODE_RANK 做 NCCL 组网
//   - head 节点 NODE_RANK=0，暴露 API 端口；worker 节点 NODE_RANK 由 Volcano 注入
//
// Volcano svc 插件会为 Job 创建 headless service，Pod DNS 格式：
//
//	{jobName}-{taskName}-{index}.{jobName}
func (b *BaseInferenceBuilder) BuildVCJob(spec *InferenceSpec) (*vcbatch.Job, error) {
	envVars := b.Strategy.GetEnvVars(spec)
	envVars = append(envVars, spec.UserEnvVars...)

	resources := helper.BuildResources(&spec.Product)
	tolerations := helper.BuildGPUTolerations(&spec.Product)

	modelVolume, modelMount := buildModelVolume(spec.ModelPvcName, spec.GetModelMountPath())
	shmVolume, shmMount := buildSHMVolume(consts.DefaultSHM)
	volumes := []corev1.Volume{modelVolume, shmVolume}
	volumeMounts := []corev1.VolumeMount{modelMount, shmMount}
	volumes = append(volumes, spec.UserVolumes...)
	volumeMounts = append(volumeMounts, spec.UserVolumeMounts...)

	headAddr := fmt.Sprintf("%s-head-0.%s", spec.Name, spec.Name)
	totalNodes := spec.WorkerCount
	affinity := helper.BuildSchedulingAffinity(spec.InstanceName, spec.AllowedNodes, spec.StrictSpread)

	headTask := b.buildHeadTask(spec, envVars, headAddr, totalNodes, resources, volumes, volumeMounts, tolerations, affinity)
	tasks := []vcbatch.TaskSpec{headTask}

	workerReplicas := spec.WorkerCount - 1
	if workerReplicas > 0 {
		workerTask := b.buildWorkerTask(spec, envVars, headAddr, totalNodes, workerReplicas, resources, volumes, volumeMounts, tolerations, affinity)
		tasks = append(tasks, workerTask)
	}

	job := &vcbatch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.Name,
			Namespace: spec.Namespace,
			Labels: map[string]string{
				consts.LabelApp:          spec.InstanceName,
				consts.LabelInstanceType: consts.InferenceInstance,
				consts.LabelJobID:        fmt.Sprintf("%d", spec.ServiceID),
				consts.LabelFramework:    spec.Framework,
				"neptune.io/instance":    spec.InstanceName,
			},
		},
		Spec: vcbatch.JobSpec{
			MinAvailable:  int32(totalNodes),
			SchedulerName: consts.VolcanoScheduler,
			Queue:         consts.DefalutQueue,
			Tasks:         tasks,
			Plugins: map[string][]string{
				"svc": {},
			},
			Policies: []vcbatch.LifecyclePolicy{
				{Event: vcbus.PodFailedEvent, Action: vcbus.RestartJobAction},
			},
		},
	}

	return job, nil
}

// buildHeadTask 构建 head 节点 TaskSpec
//
// head 节点 NODE_RANK=0，运行推理框架进程并暴露 API 端口。
// 分布式模式下使用 shell 包装，以便注入 SGLang 的 --node-rank / --dist-init-addr 参数。
func (b *BaseInferenceBuilder) buildHeadTask(
	spec *InferenceSpec,
	envVars []corev1.EnvVar,
	headAddr string,
	totalNodes int,
	resources corev1.ResourceRequirements,
	volumes []corev1.Volume,
	volumeMounts []corev1.VolumeMount,
	tolerations []corev1.Toleration,
	affinity *corev1.Affinity,
) vcbatch.TaskSpec {
	// head 节点 NCCL 环境变量
	ncclEnvs := append(envVars,
		corev1.EnvVar{Name: "MASTER_ADDR", Value: headAddr},
		corev1.EnvVar{Name: "MASTER_PORT", Value: fmt.Sprintf("%d", masterPort)},
		corev1.EnvVar{Name: "WORLD_SIZE", Value: fmt.Sprintf("%d", totalNodes)},
		corev1.EnvVar{Name: "NODE_RANK", Value: "0"},
	)

	// 分布式模式下用 shell 包装，注入框架特定的分布式参数
	finalCommand := b.wrapDistributedCommand(spec, 0)

	return vcbatch.TaskSpec{
		Name:     "head",
		Replicas: 1,
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					consts.LabelApp:          spec.InstanceName,
					consts.LabelInstanceType: consts.InferenceInstance,
					consts.LabelJobID:        fmt.Sprintf("%d", spec.ServiceID),
					consts.LabelFramework:    spec.Framework,
					"neptune.io/instance":    spec.InstanceName,
					"role":                   "head",
				},
				Annotations: map[string]string{
					"prometheus.io/scrape": "true",
					"prometheus.io/port":   fmt.Sprintf("%d", spec.ServicePort),
					"prometheus.io/path":   "/metrics",
				},
			},
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyOnFailure,
				Containers: []corev1.Container{
					{
						Name:         spec.Name + "-head",
						Image:        spec.Image,
						Command:      []string{"/bin/sh", "-c", finalCommand},
						Env:          ncclEnvs,
						Resources:    resources,
						VolumeMounts: volumeMounts,
						Ports: []corev1.ContainerPort{
							{Name: "http", ContainerPort: int32(spec.ServicePort)},
							{Name: "nccl", ContainerPort: int32(ncclPort)},
						},
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/health",
									Port: intstr.FromInt(spec.ServicePort),
								},
							},
							InitialDelaySeconds: 120,
							PeriodSeconds:       10,
						},
						LivenessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/health",
									Port: intstr.FromInt(spec.ServicePort),
								},
							},
							InitialDelaySeconds: 180,
							PeriodSeconds:       30,
						},
					},
				},
				Volumes:     volumes,
				Tolerations: tolerations,
				Affinity:    affinity,
			},
		},
	}
}

// buildWorkerTask 构建 worker 节点 TaskSpec
//
// worker 节点运行与 head 相同的推理框架进程，通过 NCCL 环境变量加入集群。
// NODE_RANK 通过 Volcano 的 VK_TASK_INDEX 环境变量注入（从 0 开始），
// worker 的 NODE_RANK = VK_TASK_INDEX + 1（head 占 0）。
// worker 不暴露 API 端口，不需要 readiness/liveness probe。
func (b *BaseInferenceBuilder) buildWorkerTask(
	spec *InferenceSpec,
	envVars []corev1.EnvVar,
	headAddr string,
	totalNodes int,
	workerReplicas int,
	resources corev1.ResourceRequirements,
	volumes []corev1.Volume,
	volumeMounts []corev1.VolumeMount,
	tolerations []corev1.Toleration,
	affinity *corev1.Affinity,
) vcbatch.TaskSpec {
	// worker 节点 NCCL 环境变量（NODE_RANK 在 shell 中动态计算）
	ncclEnvs := append(envVars,
		corev1.EnvVar{Name: "MASTER_ADDR", Value: headAddr},
		corev1.EnvVar{Name: "MASTER_PORT", Value: fmt.Sprintf("%d", masterPort)},
		corev1.EnvVar{Name: "WORLD_SIZE", Value: fmt.Sprintf("%d", totalNodes)},
	)

	// 用 shell 包装：计算 NODE_RANK 并注入分布式参数
	// VK_TASK_INDEX 由 Volcano 自动注入（从 0 开始），worker NODE_RANK = VK_TASK_INDEX + 1
	shellScript := b.buildWorkerShellScript(spec)

	return vcbatch.TaskSpec{
		Name:     "worker",
		Replicas: int32(workerReplicas),
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					consts.LabelApp:          spec.InstanceName,
					consts.LabelInstanceType: consts.InferenceInstance,
					consts.LabelJobID:        fmt.Sprintf("%d", spec.ServiceID),
					consts.LabelFramework:    spec.Framework,
					"neptune.io/instance":    spec.InstanceName,
					"role":                   "worker",
				},
			},
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyOnFailure,
				Containers: []corev1.Container{
					{
						Name:         spec.Name + "-worker",
						Image:        spec.Image,
						Command:      []string{"/bin/sh", "-c"},
						Args:         []string{shellScript},
						Env:          ncclEnvs,
						Resources:    resources,
						VolumeMounts: volumeMounts,
						Ports: []corev1.ContainerPort{
							{Name: "nccl", ContainerPort: int32(ncclPort)},
						},
					},
				},
				Volumes:     volumes,
				Tolerations: tolerations,
				Affinity:    affinity,
			},
		},
	}
}

// wrapDistributedCommand 分布式模式下包装启动命令
//
// 将用户的完整命令字符串与 args 拼接，并注入框架特定的分布式参数。
// 返回完整的 shell 脚本字符串（由调用方用 /bin/sh -c 执行）。
func (b *BaseInferenceBuilder) wrapDistributedCommand(spec *InferenceSpec, nodeRank int) string {
	baseCmd := spec.Command
	if len(spec.Args) > 0 {
		for _, arg := range spec.Args {
			baseCmd += " " + arg
		}
	}

	// SGLang 需要额外的 --node-rank 和 --dist-init-addr 参数
	if spec.Framework == consts.FrameworkSGLang {
		return fmt.Sprintf(
			`exec %s --node-rank %d --dist-init-addr ${MASTER_ADDR}:${MASTER_PORT}`,
			baseCmd, nodeRank,
		)
	}

	// vLLM 通过 MASTER_ADDR/WORLD_SIZE/NODE_RANK 环境变量自动组网，无需额外参数
	return fmt.Sprintf(`exec %s`, baseCmd)
}

// buildWorkerShellScript 构建 worker 节点的 shell 启动脚本
//
// 计算 NODE_RANK = VK_TASK_INDEX + 1，然后注入框架特定的分布式参数。
func (b *BaseInferenceBuilder) buildWorkerShellScript(spec *InferenceSpec) string {
	baseCmd := spec.Command
	if len(spec.Args) > 0 {
		for _, arg := range spec.Args {
			baseCmd += " " + arg
		}
	}

	// SGLang worker 需要 --node-rank 和 --dist-init-addr
	if spec.Framework == "SGLANG" {
		return fmt.Sprintf(
			`export NODE_RANK=$((VK_TASK_INDEX + 1)) && exec %s --node-rank $NODE_RANK --dist-init-addr ${MASTER_ADDR}:${MASTER_PORT}`,
			baseCmd,
		)
	}

	// vLLM worker 通过环境变量自动组网
	return fmt.Sprintf(
		`export NODE_RANK=$((VK_TASK_INDEX + 1)) && exec %s`,
		baseCmd,
	)
}

package builder

import (
	"fmt"
	"gin-vue-admin/model/consts"
	"gin-vue-admin/service/podgroup"
	helper "gin-vue-admin/utils/k8s"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// BuildDeployment 构建单体推理 Deployment（使用 Volcano 调度器）
//
// Volcano 调度器会根据 Pod 的 schedulerName 和 annotations 自动创建 PodGroup，
// 无需手动创建。PodGroup 删除时 Informer 会自动处理资源释放。
func (b *BaseInferenceBuilder) BuildDeployment(spec *InferenceSpec) (*appsv1.Deployment, error) {
	envVars := b.Strategy.GetEnvVars(spec)
	envVars = append(envVars, spec.UserEnvVars...)

	resources := helper.BuildResources(&spec.Product)
	tolerations := helper.BuildGPUTolerations(&spec.Product)

	modelVolume, modelMount := buildModelVolume(spec.ModelPvcName, spec.GetModelMountPath())
	shmVolume, shmMount := buildSHMVolume(16)

	volumes := []corev1.Volume{modelVolume, shmVolume}
	volumeMounts := []corev1.VolumeMount{modelMount, shmMount}
	// 合并用户自定义挂载
	volumes = append(volumes, spec.UserVolumes...)
	volumeMounts = append(volumeMounts, spec.UserVolumeMounts...)

	replicas := int32(1)

	// Pod 标签：包含 Informer 识别所需的 instance-type 和 job-id，
	// 同时保持与 K8s Service selector 的兼容
	podLabels := podgroup.BuildVolcanoLabels(
		spec.InstanceName,
		consts.InferenceInstance,
		consts.InferenceInstance,
		spec.ServiceID,
	)
	podLabels["neptune.io/instance"] = spec.InstanceName

	// Volcano 调度器 annotations（Volcano 根据此 annotation 自动创建 PodGroup）
	podAnnotations := podgroup.BuildVolcanoAnnotations(spec.InstanceName, consts.DefalutQueue)
	podAnnotations["prometheus.io/scrape"] = "true"
	podAnnotations["prometheus.io/port"] = fmt.Sprintf("%d", spec.ServicePort)
	podAnnotations["prometheus.io/path"] = "/metrics"

	// 构建启动命令：用 /bin/sh -c 包装用户的完整命令字符串
	shellScript := spec.Command
	if len(spec.Args) > 0 {
		for _, arg := range spec.Args {
			shellScript += " " + arg
		}
	}
	containerCommand := []string{"/bin/sh", "-c", shellScript}

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.Name,
			Namespace: spec.Namespace,
			Labels:    podLabels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					consts.LabelApp:       spec.InstanceName,
					"neptune.io/instance": spec.InstanceName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      podLabels,
					Annotations: podAnnotations,
				},
				Spec: corev1.PodSpec{
					SchedulerName: consts.VolcanoScheduler,
					Containers: []corev1.Container{
						{
							Name:         spec.Name,
							Image:        spec.Image,
							Command:      containerCommand,
							Env:          envVars,
							Resources:    resources,
							VolumeMounts: volumeMounts,
							Ports: []corev1.ContainerPort{
								{Name: "http", ContainerPort: int32(spec.ServicePort)},
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/health",
										Port: intstr.FromInt(spec.ServicePort),
									},
								},
								InitialDelaySeconds: 60,
								PeriodSeconds:       10,
							},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/health",
										Port: intstr.FromInt(spec.ServicePort),
									},
								},
								InitialDelaySeconds: 120,
								PeriodSeconds:       30,
							},
						},
					},
					Volumes:     volumes,
					Tolerations: tolerations,
				},
			},
		},
	}

	return deploy, nil
}

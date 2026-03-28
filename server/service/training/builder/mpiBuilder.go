package builder

import (
	"fmt"
	"gin-vue-admin/model/consts"
	trainingReq "gin-vue-admin/model/training/request"
	helper "gin-vue-admin/utils/k8s"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	vcbatch "volcano.sh/apis/pkg/apis/batch/v1alpha1"
	vcbus "volcano.sh/apis/pkg/apis/bus/v1alpha1"
)

// MPIStrategy MPI 分布式训练策略
type MPIStrategy struct{}

var _ FrameworkStrategy = (*MPIStrategy)(nil)

const mpiWorkerEntrypoint = `
echo "=== MPI Worker Environment Ready ==="

# 净化权限池，避免 StrictModes=yes 直接拒绝连接
mkdir -p /var/run/sshd /etc/ssh /var/empty/sshd /root/.ssh
chmod 755 /var/run/sshd
chmod 711 /var/empty/sshd
chmod 700 /root/.ssh 2>/dev/null || true

# sshd 至少要有一对鉴权公钥存在才能拉起服务，防止镜像太干净起不来
[ -f /etc/ssh/ssh_host_rsa_key ] || ssh-keygen -A 2>/dev/null
chmod 600 /etc/ssh/ssh_host_*_key 2>/dev/null || true

echo "Starting sshd in main container on port 22..."
# - UsePAM=no 防止缺少 pam_env.so 导致登录秒退
# - StrictModes=no 忽视 ~/.ssh/authorized_keys （来自 volcano configmap）的强制 600 权限要求引发的鉴权拒绝
exec /usr/sbin/sshd -D -e -o UsePAM=no -o StrictModes=no -o Port=22
`

// BuildTasks 构建 MPI 任务
// MPI 需要一个 Master (Launcher) 和多个 Worker
func (s *MPIStrategy) BuildTasks(spec *trainingReq.TrainingJobSpec) ([]vcbatch.TaskSpec, error) {
	if spec.WorkerCount < 1 {
		spec.WorkerCount = 1
	}

	// GPU Tolerations
	tolerations := buildGPUTolerations(&spec.Product)
	affinity := helper.BuildSchedulingAffinity(spec.Name, spec.AllowedNodes, spec.StrictSpread)

	// MPI 环境变量（Master 专用）
	masterEnvs := cloneEnvVars(spec.Envs)
	masterEnvs = append(masterEnvs,
		// 允许以 root 身份运行 mpirun
		corev1.EnvVar{Name: "OMPI_ALLOW_RUN_AS_ROOT", Value: "1"},
		corev1.EnvVar{Name: "OMPI_ALLOW_RUN_AS_ROOT_CONFIRM", Value: "1"},
		// OpenMPI 网络配置
		corev1.EnvVar{Name: "OMPI_MCA_plm_rsh_no_tree_spawn", Value: "1"},
		corev1.EnvVar{Name: "OMPI_MCA_orte_keep_fqdn_hostnames", Value: "1"},
		corev1.EnvVar{Name: "OMPI_MCA_btl_tcp_if_include", Value: "eth0"},
		// 屏蔽 GPU 设备：Launcher 不参与计算，当与 Worker 调度到同一节点时
		// 防止 NVIDIA Container Runtime 暴露节点上的 GPU 给 Launcher 容器
		corev1.EnvVar{Name: "NVIDIA_VISIBLE_DEVICES", Value: "none"},
		corev1.EnvVar{Name: "CUDA_VISIBLE_DEVICES", Value: ""},
	)

	// --- Master Task (Launcher) ---
	masterCmd := s.buildMasterCommand(spec)
	masterTask := vcbatch.TaskSpec{
		Name:     consts.TaskSpecMPIMaster,
		Replicas: 1,
		// Master 等待所有 Worker Ready 后再启动
		DependsOn: &vcbatch.DependsOn{
			Name:      []string{consts.TaskSpecMPIWorker},
			Iteration: vcbatch.IterationAll,
		},
		// Master 完成即整个 Job 完成
		Policies: []vcbatch.LifecyclePolicy{
			{
				Event:  vcbus.TaskCompletedEvent,
				Action: vcbus.CompleteJobAction,
			},
		},
		Template: buildTaskTemplate(taskTemplateOptions{
			containerName: spec.Name + "-mpimaster",
			image:         spec.Image,
			command:       masterCmd,
			resources:     buildResources(&trainingReq.ProductSpec{CPU: 1, Memory: 2}),
			volumeMounts:  spec.VolumeMounts,
			envs:          masterEnvs,
			volumes:       spec.Volumes,
			tolerations:   tolerations,
			labels:        spec.Labels,
			affinity:      affinity,
		}),
	}

	// --- Worker Task ---
	workerTask := s.buildWorkerTask(spec, tolerations, affinity)

	return []vcbatch.TaskSpec{masterTask, workerTask}, nil
}

// buildMasterCommand 构建 Master 启动命令
func (s *MPIStrategy) buildMasterCommand(spec *trainingReq.TrainingJobSpec) []string {
	userCmd := flattenCommand(spec.Command, spec.Args)

	// 自动为 mpirun/mpiexec 添加 --allow-run-as-root（如果尚未包含）
	if !strings.Contains(userCmd, "--allow-run-as-root") {
		if strings.Contains(userCmd, "mpirun ") {
			userCmd = strings.Replace(userCmd, "mpirun ", "mpirun --allow-run-as-root ", 1)
		} else if strings.Contains(userCmd, "mpiexec ") {
			userCmd = strings.Replace(userCmd, "mpiexec ", "mpiexec --allow-run-as-root ", 1)
		}
	}

	// 注入原生的 MPI_HOSTS（切断 FQDN 域名后缀，解决 ORTE daemon 迷失路由的错误）
	script := fmt.Sprintf(`export MPI_HOSTS=$(awk -F'.' '{print $1}' /etc/volcano/%s.host | tr "\n" ",")
export MPI_HOSTS=${MPI_HOSTS%%,}
echo "MPI_HOSTS=$MPI_HOSTS"

%s`, consts.TaskSpecMPIWorker, userCmd)

	return []string{"/bin/bash", "-c", script}
}

// buildWorkerTask 构建 Worker TaskSpec
func (s *MPIStrategy) buildWorkerTask(
	spec *trainingReq.TrainingJobSpec,
	tolerations []corev1.Toleration,
	affinity *corev1.Affinity,
) vcbatch.TaskSpec {
	workerResources := buildResources(&spec.Product)
	workerCmd := []string{"/bin/sh", "-c", mpiWorkerEntrypoint}
	readinessProbe := &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt(22),
			},
		},
		InitialDelaySeconds: 2,
		PeriodSeconds:       2,
		FailureThreshold:    30,
	}

	return vcbatch.TaskSpec{
		Name:     consts.TaskSpecMPIWorker,
		Replicas: int32(spec.WorkerCount),
		Template: buildTaskTemplate(taskTemplateOptions{
			containerName: spec.Name + "-worker",
			image:         spec.Image,
			command:       workerCmd,
			resources:     workerResources,
			volumeMounts:  spec.VolumeMounts,
			envs:          spec.Envs,
			volumes:       spec.Volumes,
			tolerations:   tolerations,
			readiness:     readinessProbe,
			labels:        spec.Labels,
			affinity:      affinity,
		}),
	}
}

// GetPlugins 获取 MPI 插件配置
func (s *MPIStrategy) GetPlugins() map[string][]string {
	return map[string][]string{
		"ssh": {}, // SSH 免密插件
		"svc": {}, // DNS 发现插件
	}
}

// GetMinAvailable 获取最小可用数量
// 使用 DependsOn 时 minAvailable 只含 Worker，避免死锁
func (s *MPIStrategy) GetMinAvailable(spec *trainingReq.TrainingJobSpec) int64 {
	return int64(spec.WorkerCount)
}

// GetPolicies 获取 Job 级别生命周期策略
func (s *MPIStrategy) GetPolicies() []vcbatch.LifecyclePolicy {
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

// 处理 /bin/sh -c "..." 的情况，避免双重包装
func flattenCommand(command []string, args []string) string {
	if len(command) >= 2 && (command[0] == "/bin/sh" || command[0] == "/bin/bash") && command[1] == "-c" {
		if len(args) > 0 {
			return args[0]
		}
		if len(command) >= 3 {
			return command[2]
		}
	}

	parts := append(command, args...)
	return strings.Join(parts, " ")
}

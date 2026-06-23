package builder

import (
	"fmt"
	"gin-vue-admin/model/consts"
	trainingReq "gin-vue-admin/model/training/request"
	helper "gin-vue-admin/utils/k8s"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	labels        map[string]string
	affinity      *corev1.Affinity
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
	var hasPreflightType bool
	for _, env := range opts.envs {
		if env.Name == "PREFLIGHT_TYPE" {
			hasPreflightType = true
			break
		}
	}
	if hasPreflightType {
		wrapCommandWithPreflight(&opts)
	}

	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: cloneStringMap(opts.labels),
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{
					Name:                     opts.containerName,
					Image:                    opts.image,
					Command:                  append([]string(nil), opts.command...),
					Args:                     append([]string(nil), opts.args...),
					Resources:                opts.resources,
					VolumeMounts:             append([]corev1.VolumeMount(nil), opts.volumeMounts...),
					Env:                      cloneEnvVars(opts.envs),
					ReadinessProbe:           opts.readiness,
					TerminationMessagePath:   "/dev/termination-log",
					TerminationMessagePolicy: corev1.TerminationMessageFallbackToLogsOnError,
					/*
						FallbackToLogsOnError 的作用是：如果 termination log 文件为空，并且容器以错误退出，K8s 会把最后一部分容器日志放到 termination message 里；
						API 文档说明它最多使用 2048 字节或 80 行。
					*/
				},
			},
			Volumes:     append([]corev1.Volume(nil), opts.volumes...),
			Tolerations: append([]corev1.Toleration(nil), opts.tolerations...),
			Affinity:    opts.affinity,
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

func cloneStringMap(values map[string]string) map[string]string {
	if len(values) == 0 {
		return nil
	}
	cloned := make(map[string]string, len(values))
	for k, v := range values {
		cloned[k] = v
	}
	return cloned
}

func wrapCommandWithPreflight(opts *taskTemplateOptions) {
	bootstrapScript := `cat > /tmp/preflight-and-run.sh <<'EOF'
#!/bin/bash
set -euo pipefail

if [ "${1:-}" = "--" ]; then
  shift
fi

fail() {
  local code="$1"
  local msg="$2"
  echo "[platform] ${msg}" | tee /dev/termination-log
  exit "${code}"
}

run_basic_check() {
  echo "[platform] running basic checks..."
  python -V || python3 -V || fail 10 "python not found"
}

run_gpu_check() {
  echo "[platform] running GPU checks..."
  nvidia-smi || fail 11 "gpu not visible"
  python - <<'PY' || fail 12 "torch cuda check failed"
import torch
print("torch:", torch.__version__)
print("cuda available:", torch.cuda.is_available())
print("cuda device count:", torch.cuda.device_count())
if not torch.cuda.is_available():
    raise SystemExit("cuda is not available")
if torch.cuda.device_count() == 0:
    raise SystemExit("no cuda device found")
PY
}

write_ddp_preflight_py() {
  cat > /tmp/ddp_preflight.py <<'PY'
import os
import torch
import torch.distributed as dist

required = ["MASTER_ADDR", "MASTER_PORT", "WORLD_SIZE", "RANK"]
missing = [k for k in required if not os.environ.get(k)]
if missing:
    raise RuntimeError(f"missing env: {missing}")

local_rank = int(os.environ.get("LOCAL_RANK", "0"))
torch.cuda.set_device(local_rank)

dist.init_process_group(backend="nccl", init_method="env://")

rank = dist.get_rank()
world_size = dist.get_world_size()

x = torch.ones(1, device="cuda")
dist.all_reduce(x)

actual = float(x.item())
expected = float(world_size)

print(f"[ddp-preflight] rank={rank}, actual={actual}, expected={expected}", flush=True)

if actual != expected:
    raise RuntimeError(f"all_reduce invalid: {actual} != {expected}")

dist.barrier()
dist.destroy_process_group()
PY
}

run_mpi_master_check() {
  echo "[platform] running mpi master checks..."
  command -v mpirun &>/dev/null || fail 40 "mpirun not found"
  
  local host_file="/etc/volcano/mpiworker.host"
  if [ ! -f "${host_file}" ]; then
    fail 43 "mpi worker host file ${host_file} not found"
  fi
  
  export MPI_HOSTS=$(awk -F'.' '{print $1}' "${host_file}" | tr "\n" ",")
  export MPI_HOSTS=${MPI_HOSTS%,}
  echo "[platform] resolved MPI_HOSTS=${MPI_HOSTS}"
  
  echo "[platform] running basic mpirun check..."
  mpirun --allow-run-as-root \
    --host "${MPI_HOSTS}" \
    -np "${WORLD_SIZE:-1}" \
    hostname || fail 41 "mpi hostname check failed"
    
  local nccl_bin=""
  if [ -x /opt/nccl-tests/build/all_reduce_perf ]; then
    nccl_bin="/opt/nccl-tests/build/all_reduce_perf"
  elif command -v all_reduce_perf &>/dev/null; then
    nccl_bin="all_reduce_perf"
  fi

  if [ -n "${nccl_bin}" ]; then
    echo "[platform] running nccl-tests all_reduce_perf check using ${nccl_bin}..."
    mpirun --allow-run-as-root \
      --host "${MPI_HOSTS}" \
      -np "${WORLD_SIZE:-1}" \
      -N "${GPUS_PER_NODE:-1}" \
      "${nccl_bin}" -b 8M -e 128M -f 2 -g 1 || fail 42 "mpi NCCL test failed"
  else
    echo "[platform] nccl-tests not found or not executable, skipping NCCL check"
  fi
}

case "${PREFLIGHT_TYPE:-none}" in
  none)
    run_basic_check
    ;;

  single_gpu)
    run_basic_check
    run_gpu_check
    ;;

  single_node_ddp)
    run_basic_check
    run_gpu_check
    write_ddp_preflight_py
    echo "[platform] running single node DDP preflight check..."
    torchrun \
      --standalone \
      --nproc_per_node="${NPROC_PER_NODE:-1}" \
      /tmp/ddp_preflight.py || fail 30 "single node DDP preflight check failed"
    ;;

  multi_node_ddp)
    run_basic_check
    run_gpu_check
    write_ddp_preflight_py
    
    node_rank="${NODE_RANK:-${RANK:-}}"
    if [ -z "${node_rank}" ]; then
      fail 31 "NODE_RANK or RANK env is not set for multi_node_ddp"
    fi
    
    echo "[platform] running multi node DDP preflight check (rank ${node_rank})...."
    torchrun \
      --nnodes="${NNODES:-1}" \
      --nproc_per_node="${NPROC_PER_NODE:-1}" \
      --node_rank="${node_rank}" \
      --master_addr="${MASTER_ADDR:-}" \
      --master_port="${MASTER_PORT:-29500}" \
      /tmp/ddp_preflight.py || fail 32 "multi node DDP preflight check failed"
    ;;

  mpi_master)
    run_basic_check
    run_mpi_master_check
    ;;

  mpi_worker)
    run_basic_check
    run_gpu_check
    ;;

  *)
    fail 3 "unknown PREFLIGHT_TYPE=${PREFLIGHT_TYPE}"
    ;;
esac

echo "[platform] preflight passed"
echo "[platform] start user command: $*"

exec "$@"
EOF
chmod +x /tmp/preflight-and-run.sh
exec /tmp/preflight-and-run.sh -- "$@"
`

	userCmd := append([]string(nil), opts.command...)
	userCmd = append(userCmd, opts.args...)

	opts.command = []string{"/bin/bash", "-lc", bootstrapScript, "--"}
	opts.args = userCmd
}

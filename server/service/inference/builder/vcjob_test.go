package builder

import (
	"strings"
	"testing"

	"gin-vue-admin/model/consts"
)

func TestBuildVCJob_SvcPluginDisablesNetworkPolicy(t *testing.T) {
	t.Parallel()

	builder := NewInferenceBuilder(consts.FrameworkSGLang)
	b, ok := builder.(*BaseInferenceBuilder)
	if !ok {
		t.Fatal("NewInferenceBuilder() did not return *BaseInferenceBuilder")
	}
	spec := &InferenceSpec{
		Name:         "inference-e64025",
		InstanceName: "inference-e64025",
		Namespace:    "zzz",
		Framework:    consts.FrameworkSGLang,
		WorkerCount:  2,
	}

	job, err := b.BuildVCJob(spec)
	if err != nil {
		t.Fatalf("BuildVCJob() error = %v", err)
	}

	args, ok := job.Spec.Plugins["svc"]
	if !ok {
		t.Fatal("BuildVCJob() svc plugin not configured")
	}
	if len(args) != 2 {
		t.Fatalf("BuildVCJob() svc plugin args len = %d, want 2", len(args))
	}
	if args[0] != "--publish-not-ready-addresses=true" {
		t.Fatalf("BuildVCJob() svc plugin args[0] = %q, want %q", args[0], "--publish-not-ready-addresses=true")
	}
	if args[1] != "--disable-network-policy=true" {
		t.Fatalf("BuildVCJob() svc plugin args[1] = %q, want %q", args[1], "--disable-network-policy=true")
	}
}

func TestWrapDistributedCommand_SGLangInjectsMultiNodeFlags(t *testing.T) {
	t.Parallel()

	b := &BaseInferenceBuilder{}
	spec := &InferenceSpec{
		Framework:   consts.FrameworkSGLang,
		Command:     "python -m sglang.launch_server --model-path /model/foo --tp 4 --port 30000",
		WorkerCount: 2,
	}

	cmd := b.wrapDistributedCommand(spec, 0)

	for _, want := range []string{
		"--tp 4",
		"--nnodes 2",
		"--node-rank 0",
		"--dist-init-addr ${MASTER_ADDR}:${MASTER_PORT}",
	} {
		if !strings.Contains(cmd, want) {
			t.Fatalf("wrapDistributedCommand() = %q, want substring %q", cmd, want)
		}
	}
}

func TestWrapDistributedCommand_VLLMDoesNotInjectHeadlessForHead(t *testing.T) {
	t.Parallel()

	b := &BaseInferenceBuilder{}
	spec := &InferenceSpec{
		Framework: consts.FrameworkVLLM,
		Command:   "vllm serve /model/foo --host 0.0.0.0 --port 30000 --tensor-parallel-size 4 --pipeline-parallel-size 1 --distributed-executor-backend mp --nnodes 2 --node-rank ${NODE_RANK} --master-addr ${MASTER_ADDR} --master-port ${MASTER_PORT}",
	}

	cmd := b.wrapDistributedCommand(spec, 0)

	if strings.Contains(cmd, "--headless") {
		t.Fatalf("wrapDistributedCommand() = %q, should not contain --headless for head", cmd)
	}
	for _, want := range []string{
		"--distributed-executor-backend mp",
		"--node-rank ${NODE_RANK}",
		"--master-addr ${MASTER_ADDR}",
		"--master-port ${MASTER_PORT}",
	} {
		if !strings.Contains(cmd, want) {
			t.Fatalf("wrapDistributedCommand() = %q, want substring %q", cmd, want)
		}
	}
}

func TestBuildWorkerShellScript_VLLMAppendsHeadless(t *testing.T) {
	t.Parallel()

	b := &BaseInferenceBuilder{}
	spec := &InferenceSpec{
		Framework: consts.FrameworkVLLM,
		Command:   "vllm serve /model/foo --host 0.0.0.0 --port 30000 --tensor-parallel-size 4 --pipeline-parallel-size 1 --distributed-executor-backend mp --nnodes 2 --node-rank ${NODE_RANK} --master-addr ${MASTER_ADDR} --master-port ${MASTER_PORT}",
	}

	script := b.buildWorkerShellScript(spec)

	for _, want := range []string{
		"export NODE_RANK=$((VK_TASK_INDEX + 1))",
		"--distributed-executor-backend mp",
		"--node-rank ${NODE_RANK}",
		"--master-addr ${MASTER_ADDR}",
		"--master-port ${MASTER_PORT}",
		"--headless",
	} {
		if !strings.Contains(script, want) {
			t.Fatalf("buildWorkerShellScript() = %q, want substring %q", script, want)
		}
	}
}

func TestBuildWorkerShellScript_SGLangInjectsMultiNodeFlags(t *testing.T) {
	t.Parallel()

	b := &BaseInferenceBuilder{}
	spec := &InferenceSpec{
		Framework:   consts.FrameworkSGLang,
		Command:     "python -m sglang.launch_server --model-path /model/foo --tp 4 --port 30000",
		WorkerCount: 2,
	}

	script := b.buildWorkerShellScript(spec)

	for _, want := range []string{
		"export NODE_RANK=$((VK_TASK_INDEX + 1))",
		"--tp 4",
		"--nnodes 2",
		"--node-rank $NODE_RANK",
		"--dist-init-addr ${MASTER_ADDR}:${MASTER_PORT}",
	} {
		if !strings.Contains(script, want) {
			t.Fatalf("buildWorkerShellScript() = %q, want substring %q", script, want)
		}
	}
}

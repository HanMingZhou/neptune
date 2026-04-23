package builder

import (
	"strings"
	"testing"

	"gin-vue-admin/model/consts"
)

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

package inference

import (
	"testing"

	"gin-vue-admin/model/consts"
	inferenceReq "gin-vue-admin/model/inference/request"
)

func TestValidateCreateRequest_StandaloneFrameworkIsNormalized(t *testing.T) {
	t.Parallel()

	svc := &InferenceServiceService{}
	req := &inferenceReq.CreateInferenceServiceReq{
		DeployType: consts.DeployTypeStandalone,
		Framework:  " vllm ",
		Command:    "python3 -m vllm.entrypoints.openai.api_server",
	}

	if err := svc.validateCreateRequest(req); err != nil {
		t.Fatalf("validateCreateRequest() error = %v", err)
	}
	if req.Framework != consts.FrameworkVLLM {
		t.Fatalf("validateCreateRequest() framework = %q, want %q", req.Framework, consts.FrameworkVLLM)
	}
}

func TestValidateCreateRequest_DistributedRequiresFramework(t *testing.T) {
	t.Parallel()

	svc := &InferenceServiceService{}
	req := &inferenceReq.CreateInferenceServiceReq{
		DeployType:    consts.DeployTypeDistributed,
		Framework:     "",
		Command:       "python3 -m vllm.entrypoints.openai.api_server",
		InstanceCount: 2,
	}

	if err := svc.validateCreateRequest(req); err == nil {
		t.Fatal("validateCreateRequest() expected error, got nil")
	}
}

func TestValidateCreateRequest_StandaloneRejectsInvalidFramework(t *testing.T) {
	t.Parallel()

	svc := &InferenceServiceService{}
	req := &inferenceReq.CreateInferenceServiceReq{
		DeployType: consts.DeployTypeStandalone,
		Framework:  "INVALID",
		Command:    "python3 main.py",
	}

	if err := svc.validateCreateRequest(req); err == nil {
		t.Fatal("validateCreateRequest() expected error, got nil")
	}
}

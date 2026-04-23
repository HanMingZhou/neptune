package inference

import (
	"context"
	"testing"

	"gin-vue-admin/global"
	apisixReq "gin-vue-admin/model/apisix/request"
	"gin-vue-admin/model/consts"
	inferenceModel "gin-vue-admin/model/inference"
	inferenceReq "gin-vue-admin/model/inference/request"
)

type mockInferenceApisixService struct {
	createCalls int
	updateCalls int
	deleteCalls int
	updateReq   *apisixReq.CreateRouteReq
	updateErr   error
}

func (m *mockInferenceApisixService) CreateRoute(_ context.Context, _ *apisixReq.CreateRouteReq) error {
	m.createCalls++
	return nil
}

func (m *mockInferenceApisixService) UpdateRoute(_ context.Context, req *apisixReq.CreateRouteReq) error {
	m.updateCalls++
	copied := *req
	if req.Labels != nil {
		copied.Labels = make(map[string]string, len(req.Labels))
		for k, v := range req.Labels {
			copied.Labels[k] = v
		}
	}
	m.updateReq = &copied
	return m.updateErr
}

func (m *mockInferenceApisixService) DeleteRoute(_ context.Context, _ *apisixReq.DeleteRouteReq) error {
	m.deleteCalls++
	return nil
}

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

func TestInferenceServiceNamesForCleanup_DistributedIncludesLegacyAndAPIService(t *testing.T) {
	t.Parallel()

	service := &inferenceModel.Inference{
		InstanceName: "inference-e64025",
		DeployType:   consts.DeployTypeDistributed,
	}

	names := inferenceServiceNamesForCleanup(service)
	if len(names) != 2 {
		t.Fatalf("inferenceServiceNamesForCleanup() len = %d, want 2", len(names))
	}
	if names[0] != "inference-e64025-api" {
		t.Fatalf("inferenceServiceNamesForCleanup()[0] = %q, want %q", names[0], "inference-e64025-api")
	}
	if names[1] != "inference-e64025" {
		t.Fatalf("inferenceServiceNamesForCleanup()[1] = %q, want %q", names[1], "inference-e64025")
	}
}

func TestEnsureInferenceRoute_UsesUpdateRouteAndHonorsAuthSwitch(t *testing.T) {
	original := global.GVA_CONFIG.Apisix
	t.Cleanup(func() {
		global.GVA_CONFIG.Apisix = original
	})

	global.GVA_CONFIG.Apisix.Enabled = true
	global.GVA_CONFIG.Apisix.BaseDomain = "gateway.example.com"
	global.GVA_CONFIG.Apisix.AuthEnabled = false
	global.GVA_CONFIG.Apisix.AuthUri = "http://neptune-server.neptune.svc.cluster.local:8888/aiInfra/api/v1/apisix/auth"

	mockApisix := &mockInferenceApisixService{}
	svc := &InferenceServiceService{apisixSvc: mockApisix}
	service := &inferenceModel.Inference{
		InstanceName: "inference-e64025",
		Namespace:    "zzz",
		ClusterID:    7,
	}

	if err := svc.ensureInferenceRoute(context.Background(), service); err != nil {
		t.Fatalf("ensureInferenceRoute() error = %v", err)
	}
	if mockApisix.createCalls != 0 {
		t.Fatalf("ensureInferenceRoute() createCalls = %d, want 0", mockApisix.createCalls)
	}
	if mockApisix.updateCalls != 1 {
		t.Fatalf("ensureInferenceRoute() updateCalls = %d, want 1", mockApisix.updateCalls)
	}
	if mockApisix.updateReq == nil {
		t.Fatal("ensureInferenceRoute() updateReq is nil")
	}
	if mockApisix.updateReq.EnableAuth {
		t.Fatal("ensureInferenceRoute() unexpectedly enabled auth")
	}
	if mockApisix.updateReq.Path != "/inference/zzz/inference-e64025/*" {
		t.Fatalf("ensureInferenceRoute() path = %q", mockApisix.updateReq.Path)
	}
}

func TestEnsureInferenceRoute_DistributedUsesDedicatedAPIServiceName(t *testing.T) {
	original := global.GVA_CONFIG.Apisix
	t.Cleanup(func() {
		global.GVA_CONFIG.Apisix = original
	})

	global.GVA_CONFIG.Apisix.Enabled = true
	global.GVA_CONFIG.Apisix.BaseDomain = "gateway.example.com"
	global.GVA_CONFIG.Apisix.AuthEnabled = false
	global.GVA_CONFIG.Apisix.AuthUri = "http://neptune-server.neptune.svc.cluster.local:8888/aiInfra/api/v1/apisix/auth"

	mockApisix := &mockInferenceApisixService{}
	svc := &InferenceServiceService{apisixSvc: mockApisix}
	service := &inferenceModel.Inference{
		InstanceName: "inference-e64025",
		Namespace:    "zzz",
		ClusterID:    7,
		DeployType:   consts.DeployTypeDistributed,
	}

	if err := svc.ensureInferenceRoute(context.Background(), service); err != nil {
		t.Fatalf("ensureInferenceRoute() error = %v", err)
	}
	if mockApisix.updateReq == nil {
		t.Fatal("ensureInferenceRoute() updateReq is nil")
	}
	if mockApisix.updateReq.ServiceName != "inference-e64025-api" {
		t.Fatalf("ensureInferenceRoute() serviceName = %q, want %q", mockApisix.updateReq.ServiceName, "inference-e64025-api")
	}
}

func TestEnsureInferenceRoute_RejectsMissingAuthURIWhenAuthEnabled(t *testing.T) {
	original := global.GVA_CONFIG.Apisix
	t.Cleanup(func() {
		global.GVA_CONFIG.Apisix = original
	})

	global.GVA_CONFIG.Apisix.Enabled = true
	global.GVA_CONFIG.Apisix.AuthEnabled = true
	global.GVA_CONFIG.Apisix.AuthUri = ""

	mockApisix := &mockInferenceApisixService{}
	svc := &InferenceServiceService{apisixSvc: mockApisix}
	service := &inferenceModel.Inference{
		InstanceName: "inference-e64025",
		Namespace:    "zzz",
	}

	if err := svc.ensureInferenceRoute(context.Background(), service); err == nil {
		t.Fatal("ensureInferenceRoute() expected error, got nil")
	}
	if mockApisix.updateCalls != 0 {
		t.Fatalf("ensureInferenceRoute() updateCalls = %d, want 0", mockApisix.updateCalls)
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

package response

import "time"

// CreateInferenceServiceResp 创建推理服务响应
type CreateInferenceServiceResp struct {
	ID           uint   `json:"id"`
	InstanceName string `json:"instanceName"`
	Status       string `json:"status"`
	AccessUrl    string `json:"accessUrl"` // 访问地址
}

// GetInferenceServiceListResp 获取推理服务列表响应
type GetInferenceServiceListResp struct {
	Total int64                      `json:"total"`
	List  []InferenceServiceListItem `json:"list"`
}

// InferenceServiceListItem 列表项
type InferenceServiceListItem struct {
	ID           uint       `json:"id"`
	DisplayName  string     `json:"displayName"`
	InstanceName string     `json:"instanceName"`
	DeployType   string     `json:"deployType"`
	Framework    string     `json:"framework"`
	Status       string     `json:"status"`
	GPU          int64      `json:"gpu"`
	GPUModel     string     `json:"gpuModel"`
	CreatedAt    time.Time  `json:"createdAt"`
	StartedAt    *time.Time `json:"startedAt"`
}

// InferenceServiceDetail 推理服务详情
type InferenceServiceDetail struct {
	ID               uint                 `json:"id"`
	DisplayName      string               `json:"displayName"`
	InstanceName     string               `json:"instanceName"`
	Namespace        string               `json:"namespace"`
	DeployType       string               `json:"deployType"`
	Framework        string               `json:"framework"`
	ModelMountPath   string               `json:"modelMountPath"`
	ImageName        string               `json:"imageName"`
	TensorParallel   int                  `json:"tensorParallel"`
	PipelineParallel int                  `json:"pipelineParallel"`
	WorkerCount      int                  `json:"workerCount"`
	CPU              int64                `json:"cpu"`
	Memory           int64                `json:"memory"`
	GPU              int64                `json:"gpu"`
	GPUModel         string               `json:"gpuModel"`
	ServicePort      int                  `json:"servicePort"`
	Command          string               `json:"command"`
	Args             []string             `json:"args"`
	AutoRestart      bool                 `json:"autoRestart"`
	RestartCount     int                  `json:"restartCount"`
	MaxRestarts      int                  `json:"maxRestarts"`
	AuthType         int                  `json:"authType"`
	Status           string               `json:"status"`
	ErrorMsg         string               `json:"errorMsg"`
	AccessUrl        string               `json:"accessUrl"`
	CreatedAt        time.Time            `json:"createdAt"`
	StartedAt        *time.Time           `json:"startedAt"`
	Mounts           []InferenceMountItem `json:"mounts"`
	Envs             []InferenceEnvItem   `json:"envs"`
}

// InferenceMountItem 挂载项
type InferenceMountItem struct {
	MountType string `json:"mountType"`
	PvcName   string `json:"pvcName"`
	SubPath   string `json:"subPath"`
	MountPath string `json:"mountPath"`
	ReadOnly  bool   `json:"readOnly"`
}

// InferenceEnvItem 环境变量项
type InferenceEnvItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// AuthResult 认证结果
type AuthResult struct {
	Valid        bool   `json:"valid"`
	UserId       uint   `json:"userId"`
	ServiceId    uint   `json:"serviceId"`
	Message      string `json:"message,omitempty"`
	UpstreamHost string `json:"upstreamHost,omitempty"` // 动态 upstream 地址
	RateLimit    int    `json:"rateLimit,omitempty"`    // 限流配置
}

// CreateApiKeyResp 创建 API Key 响应
type CreateApiKeyResp struct {
	ID     uint   `json:"id"`
	ApiKey string `json:"apiKey"` // 仅创建时返回完整 Key
	Name   string `json:"name"`
}

// ListApiKeysResp 获取 API Key 列表响应
type ListApiKeysResp struct {
	Total int64        `json:"total"`
	List  []ApiKeyItem `json:"list"`
}

// ApiKeyItem API Key 列表项
type ApiKeyItem struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	ApiKey      string     `json:"apiKey"` // 脱敏显示，如 sk-***abc
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Scopes      string     `json:"scopes"`
	RateLimit   int        `json:"rateLimit"`
	LastUsedAt  *time.Time `json:"lastUsedAt"`
	ExpiredAt   *time.Time `json:"expiredAt"`
	CreatedAt   time.Time  `json:"createdAt"`
}

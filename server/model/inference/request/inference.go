package request

// CreateInferenceServiceReq 创建推理服务请求
type CreateInferenceServiceReq struct {
	Name       string `json:"name" binding:"required"`       // 服务名称
	DeployType string `json:"deployType" binding:"required"` // 部署类型: DISTRIBUTED/STANDALONE
	Framework  string `json:"framework"`                     // 框架: SGLANG/VLLM（分布式必填，单机可选，仅用于 NCCL 环境变量注入）

	// 模型存储卷配置
	ModelMountPath string `json:"modelMountPath"`                // PVC 挂载到容器的路径（默认 /model）
	ModelPvcId     uint   `json:"modelPvcId" binding:"required"` // 模型PVC ID

	// 镜像和产品
	ImageId   uint `json:"imageId" binding:"required"`   // 镜像ID
	ProductId uint `json:"productId" binding:"required"` // 产品ID

	// 并行配置 (分布式模式必填)
	TensorParallel   int `json:"tensorParallel"`   // 张量并行度
	PipelineParallel int `json:"pipelineParallel"` // 流水线并行度
	WorkerCount      int `json:"workerCount"`      // Worker数量（总节点数，包含 head）

	// 启动命令（用户完全控制，必填）
	Command string   `json:"command" binding:"required"` // 完整启动命令字符串，如 "python3 -m vllm.entrypoints.openai.api_server --model /model/Qwen2.5 --port 8000"
	Args    []string `json:"args"`                       // 启动参数（可选，每行一个，追加到命令末尾）

	// 服务配置
	ServicePort int `json:"servicePort"` // 服务端口，默认30000

	// 自愈配置 (分布式模式)
	AutoRestart bool `json:"autoRestart"` // 是否允许自动重启
	MaxRestarts int  `json:"maxRestarts"` // 最大重启次数，默认3

	// 认证配置
	AuthType int `json:"authType"` // 认证类型 1-JWT 2-API-Key

	// 挂载配置
	Mounts []CreateInferenceMountReq `json:"mounts"` // 额外挂载配置

	// 环境变量
	Envs []CreateInferenceEnvReq `json:"envs"` // 环境变量

	// 计费
	PayType int64 `json:"payType"` // 付费类型
	UserId  uint  `json:"userId"`  // 用户ID (由后端填充)
}

// CreateInferenceMountReq 挂载配置请求
type CreateInferenceMountReq struct {
	MountType string `json:"mountType"` // 挂载类型: MODEL/DATA
	PvcId     uint   `json:"pvcId"`     // PVC ID
	SubPath   string `json:"subPath"`   // 子路径
	MountPath string `json:"mountPath"` // 容器内挂载路径
	ReadOnly  bool   `json:"readOnly"`  // 是否只读
}

// CreateInferenceEnvReq 环境变量请求
type CreateInferenceEnvReq struct {
	Name  string `json:"name"`  // 环境变量名
	Value string `json:"value"` // 环境变量值
}

// GetInferenceServiceListReq 获取推理服务列表请求
type GetInferenceServiceListReq struct {
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"pageSize" form:"pageSize"`
	Name       string `json:"name" form:"name"`             // 名称模糊搜索
	Status     string `json:"status" form:"status"`         // 状态筛选
	DeployType string `json:"deployType" form:"deployType"` // 部署类型筛选
	Framework  string `json:"framework" form:"framework"`   // 框架筛选
	UserId     uint   `json:"userId" form:"userId"`         // 用户ID (由后端填充)
}

// GetInferenceServiceDetailReq 获取推理服务详情请求
type GetInferenceServiceDetailReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

// DeleteInferenceServiceReq 删除推理服务请求
type DeleteInferenceServiceReq struct {
	ID uint `json:"id" binding:"required"`
}

// StopInferenceServiceReq 停止推理服务请求
type StopInferenceServiceReq struct {
	ID uint `json:"id" binding:"required"`
}

// StartInferenceServiceReq 启动推理服务请求
type StartInferenceServiceReq struct {
	ID uint `json:"id" binding:"required"`
}

// GetInferenceServiceLogsReq 获取推理服务日志请求
type GetInferenceServiceLogsReq struct {
	ID        uint   `json:"id" form:"id" binding:"required"`
	TaskName  string `json:"taskName" form:"taskName"`   // Task名称 (master/worker)
	PodIndex  *int   `json:"podIndex" form:"podIndex"`   // Pod索引
	PodName   string `json:"podName" form:"podName"`     // 直接指定 Pod 名称（优先级高于 TaskName+PodIndex）
	Container string `json:"container" form:"container"` // 容器名称
	TailLines int64  `json:"tailLines" form:"tailLines"` // 尾部行数
	Follow    bool   `json:"follow" form:"follow"`       // 是否实时跟踪
}

// GetInferenceServicePodsReq 获取推理服务Pod列表请求
type GetInferenceServicePodsReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

// HandleTerminalReq 推理服务终端连接请求
type HandleTerminalReq struct {
	ID        uint   `json:"id" form:"id" binding:"required"`
	Token     string `json:"token" form:"token" binding:"required"` // JWT Token（WebSocket 无法携带 Header）
	PodName   string `json:"podName" form:"podName"`                // 指定 Pod 名称（可选）
	Container string `json:"container" form:"container"`            // 指定容器名称（可选）
}

// CreateApiKeyReq 创建 API Key 请求
type CreateApiKeyReq struct {
	ServiceId   uint   `json:"serviceId" binding:"required"` // 推理服务ID
	Name        string `json:"name" binding:"required"`      // Key名称
	Description string `json:"description"`                  // 描述
	ExpireDays  int    `json:"expireDays"`                   // 过期天数，0表示永不过期
	Scopes      string `json:"scopes"`                       // 权限范围
	RateLimit   int    `json:"rateLimit"`                    // 每分钟请求限制
	UserId      uint   `json:"userId"`                       // 用户ID (由后端填充)
}

// ListApiKeysReq 获取 API Key 列表请求
type ListApiKeysReq struct {
	ServiceId uint `json:"serviceId" form:"serviceId" binding:"required"`
	Page      int  `json:"page" form:"page"`
	PageSize  int  `json:"pageSize" form:"pageSize"`
}

// DeleteApiKeyReq 删除 API Key 请求
type DeleteApiKeyReq struct {
	ID uint `json:"id" binding:"required"`
}

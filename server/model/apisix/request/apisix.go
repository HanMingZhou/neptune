package request

// CreateRouteReq 创建路由请求
type CreateRouteReq struct {
	Name          string            `json:"name"`          // 路由名称
	Namespace     string            `json:"namespace"`     // 命名空间
	ClusterId     uint              `json:"clusterId"`     // 集群ID
	Host          string            `json:"host"`          // 域名
	Path          string            `json:"path"`          // 路径匹配（可选）
	RewriteRegex  string            `json:"rewriteRegex"`  // 路径重写正则（用于 proxy-rewrite 插件的 regex_uri）
	RewriteTarget string            `json:"rewriteTarget"` // 路径重写目标（可选，如 "/$1"）
	ServiceName   string            `json:"serviceName"`   // 后端服务名称
	ServicePort   int               `json:"servicePort"`   // 后端服务端口
	Labels        map[string]string `json:"labels"`        // 标签
	Websocket     bool              `json:"websocket"`     // 是否开启 WebSocket
	// forward-auth 认证配置
	EnableAuth bool   `json:"enableAuth"` // 是否启用 forward-auth 认证
	AuthUri    string `json:"authUri"`    // 认证服务地址，如 http://backend:8001/aiInfra/notebook/auth
}

// DeleteRouteReq 删除路由请求
type DeleteRouteReq struct {
	Name      string `json:"name"`      // 路由名称
	Namespace string `json:"namespace"` // 命名空间
	ClusterId uint   `json:"clusterId"` // 集群ID
}

// CreateStreamRouteReq 创建 TCP Stream 路由请求（用于 SSH 等 TCP 服务）
type CreateStreamRouteReq struct {
	Name             string            `json:"name"`             // 路由名称
	Namespace        string            `json:"namespace"`        // 命名空间（路由创建在 Apisix 所在的 namespace）
	ClusterId        uint              `json:"clusterId"`        // 集群ID
	IngressPort      int               `json:"ingressPort"`      // Apisix 监听的入口端口（如 9100，对外通过 Gateway 映射为 22）
	ServiceName      string            `json:"serviceName"`      // 后端服务名称（如 sshpiper）
	ServiceNamespace string            `json:"serviceNamespace"` // 后端服务所在的 namespace
	ServicePort      int               `json:"servicePort"`      // 后端服务端口（如 22）
	Labels           map[string]string `json:"labels"`           // 标签
}

// DeleteStreamRouteReq 删除 TCP Stream 路由请求
type DeleteStreamRouteReq struct {
	Name      string `json:"name"`      // 路由名称
	Namespace string `json:"namespace"` // 命名空间
	ClusterId uint   `json:"clusterId"` // 集群ID
}

type AuthApisixReq struct {
	OriginalUri string `json:"originalUri" form:"originalUri"` // 原始请求路径
	Token       string `json:"token" form:"token"`             // 用户 Token (JWT)
	ApiKey      string `json:"apiKey" form:"apiKey"`           // API-Key（推理服务可选）
}

package config

// Apisix Apisix 网关配置
type Apisix struct {
	// Enabled 是否启用 Apisix 路由自动创建
	Enabled bool `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	// BaseDomain 基础域名，用于生成访问 URL，如 ai.dev.com
	BaseDomain string `mapstructure:"base-domain" json:"base-domain" yaml:"base-domain"`
	// GatewayNamespace Apisix 网关所在的命名空间
	GatewayNamespace string `mapstructure:"gateway-namespace" json:"gateway-namespace" yaml:"gateway-namespace"`
	// AuthEnabled 是否启用访问认证 (forward-auth)
	AuthEnabled bool `mapstructure:"auth-enabled" json:"auth-enabled" yaml:"auth-enabled"`
	// AuthUri forward-auth 认证服务地址（Notebook / TensorBoard / Inference 统一入口）
	// 注意：此地址必须是 Apisix 能访问到的后端地址
	// 开发环境（Docker Desktop）: http://host.docker.internal:8001/aiInfra/api/v1/apisix/auth
	// 生产环境（K8s）: http://<backend-service>.<namespace>.svc.cluster.local:8001/aiInfra/api/v1/apisix/auth
	AuthUri string `mapstructure:"auth-uri" json:"auth-uri" yaml:"auth-uri"`
	// HttpPort Apisix 网关对外暴露的 HTTP 入口端口（Jupyter、TensorBoard、Inference 共用）
	HttpPort int `mapstructure:"http-port" json:"http-port" yaml:"http-port"`
}

package apisix

// Apisix 配置常量
const (
	// DefaultBaseDomain 默认基础域名
	DefaultBaseDomain = "ai.local"

	// RoutePrefix 路由名称前缀
	RoutePrefix = "route"

	// StreamRoutePrefix TCP Stream 路由名称前缀
	StreamRoutePrefix = "stream"

	// DefaultSSHIngressPort Apisix 监听的 SSH 入口端口
	DefaultSSHIngressPort = 22
)

// ApisixConfig Apisix 配置
type ApisixConfig struct {
	// Enabled 是否启用 Apisix 路由
	Enabled bool `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	// BaseDomain 基础域名，如 ai.dev.com
	BaseDomain string `mapstructure:"base-domain" json:"base_domain" yaml:"base-domain"`
	// GatewayNamespace Apisix 网关所在的命名空间
	GatewayNamespace string `mapstructure:"gateway-namespace" json:"gateway_namespace" yaml:"gateway-namespace"`
}

// NotebookRouteConfig Notebook 路由配置
type NotebookRouteConfig struct {
	// NotebookName Notebook 名称
	NotebookName string
	// Namespace 命名空间
	Namespace string
	// ServiceName 后端服务名称
	ServiceName string
	// ServicePort 后端服务端口
	ServicePort int
	// BaseDomain 基础域名
	BaseDomain string
	// EnableWebsocket 是否启用 WebSocket
	EnableWebsocket bool
}

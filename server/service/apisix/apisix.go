package apisix

import (
	"context"
	"fmt"
	"gin-vue-admin/global"
	apisixReq "gin-vue-admin/model/apisix/request"
	"gin-vue-admin/model/apisix/response"
	"gin-vue-admin/model/consts"
	inferenceModel "gin-vue-admin/model/inference"
	nbModel "gin-vue-admin/model/notebook"
	tensorboardModel "gin-vue-admin/model/tensorboard"
	"gin-vue-admin/utils"
	"strings"
	"time"

	apisixv2 "github.com/apache/apisix-ingress-controller/pkg/kube/apisix/apis/config/v2"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// ApisixManager Apisix 路由管理接口
type ApisixManager interface {
	// HTTP 路由
	CreateRoute(ctx context.Context, req *apisixReq.CreateRouteReq) error
	UpdateRoute(ctx context.Context, req *apisixReq.CreateRouteReq) error // 更新路由（删除后重建）
	DeleteRoute(ctx context.Context, req *apisixReq.DeleteRouteReq) error
	// TCP Stream 路由（用于 SSH 等 TCP 服务）
	CreateStreamRoute(ctx context.Context, req *apisixReq.CreateStreamRouteReq) error
	DeleteStreamRoute(ctx context.Context, req *apisixReq.DeleteStreamRouteReq) error
	// AUTH
	AuthApiSix(ctx context.Context, req *apisixReq.AuthApisixReq) (*response.AuthApisixResp, error)
}

var _ ApisixManager = &ApisixService{}

// ApisixService 基于 Kubernetes CRD 的 Apisix 路由管理服务
// 多集群模式：每次操作时根据ClusterId动态获取对应集群的client
type ApisixService struct {
	notebookSvc interface{} // Notebook 服务接口（预留）
}

// NewApisixManager 创建 Apisix 路由管理器
func NewApisixManager() *ApisixService {
	return &ApisixService{}
}

// SetNotebookService 设置 Notebook 服务依赖
func (m *ApisixService) SetNotebookService(svc interface{}) {
	m.notebookSvc = svc
}

// CreateRoute 创建 Apisix HTTP 路由
func (m *ApisixService) CreateRoute(ctx context.Context, req *apisixReq.CreateRouteReq) error {
	// 根据ClusterId获取集群客户端
	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(req.ClusterId)
	if cluster == nil {
		return errors.Wrapf(errors.New("没有可用的K8s集群"), "")
	}

	if cluster.ApisixClient == nil {
		return errors.Errorf("Apisix client 未初始化")
	}

	// 1. 构造路径：使用传入的 Path，如果没有则使用默认的 notebook 路径格式
	var paths []string
	if req.Path != "" {
		paths = []string{req.Path}
	} else {
		logx.Error("path is empty")
		return errors.Wrapf(errors.New("path is empty"), "")
	}

	// 2. 构造 Hosts (如果不传 Host，Apisix 会匹配所有域名，建议加上)
	var hosts []string
	if req.Host != "" {
		hosts = []string{req.Host}
	}

	serviceName := req.ServiceName

	// 3. 构建 HTTP 路由配置
	httpRoute := apisixv2.ApisixRouteHTTP{
		Name: "access",
		Match: apisixv2.ApisixRouteHTTPMatch{
			Paths: paths,
			Hosts: hosts,
		},
		Backends: []apisixv2.ApisixRouteHTTPBackend{
			{
				ServiceName: serviceName,
				ServicePort: intstr.FromInt(req.ServicePort),
			},
		},
		// 注意：不使用 Apisix 内置的 jwt-auth（需要额外创建 ApisixConsumer）
		// 而是使用 forward-auth 插件，复用后端现有的 JWT 验证逻辑
		Websocket: req.Websocket,
	}

	// 添加插件
	var plugins []apisixv2.ApisixRoutePlugin

	// 调试日志
	logx.Infof("CreateRoute: Name=%s, EnableAuth=%v, AuthUri=%s", req.Name, req.EnableAuth, req.AuthUri)

	// 如果启用 forward-auth 认证
	if req.EnableAuth && req.AuthUri != "" {
		// forward-auth 插件：验证用户身份
		plugins = append(plugins, apisixv2.ApisixRoutePlugin{
			Name:   "forward-auth",
			Enable: true,
			Config: apisixv2.ApisixRoutePluginConfig{
				"uri":              req.AuthUri,                                               // 认证服务地址
				"request_headers":  []string{"Cookie", "x-token", "Authorization", "API-Key"}, // 转发这些请求头到认证服务
				"upstream_headers": []string{"X-User-Id", "X-User-Namespace", "X-Set-Token"},  // 认证成功后传递给 upstream（和后续插件）
				"client_headers":   []string{"Location", "X-Auth-Error"},                      // 认证失败时返回给客户端
			},
		})

		// serverless-post-function 插件：在 header_filter 阶段设置 Cookie
		// 读取 X-Set-Token header，生成 Set-Cookie 响应头
		// 注意：使用 add_header 追加而不是 set_header，避免覆盖 JupyterLab 的 _xsrf Cookie
		luaCode := `
return function(conf, ctx)
    local core = require("apisix.core")
    local token = core.request.header(ctx, "X-Set-Token")
    if token and token ~= "" then
        local cookie = "x-token=" .. token .. "; Path=/; Max-Age=604800"
        core.response.add_header("Set-Cookie", cookie)
    end
end
`
		plugins = append(plugins, apisixv2.ApisixRoutePlugin{
			Name:   "serverless-post-function",
			Enable: true,
			Config: apisixv2.ApisixRoutePluginConfig{
				"phase":     "header_filter",
				"functions": []string{luaCode},
			},
		})
	}

	// proxy-rewrite 插件：剥离路径前缀
	// 将 /notebook/{namespace}/{name}/xxx 重写为 /xxx
	// 这样 Jupyter 不需要知道自己的 base_url
	if req.RewriteRegex != "" {
		rewriteTarget := req.RewriteTarget
		if rewriteTarget == "" {
			rewriteTarget = "$1"
		}
		plugins = append(plugins, apisixv2.ApisixRoutePlugin{
			Name:   "proxy-rewrite",
			Enable: true,
			Config: apisixv2.ApisixRoutePluginConfig{
				"regex_uri": []string{req.RewriteRegex, rewriteTarget},
			},
		})
	}

	if len(plugins) > 0 {
		httpRoute.Plugins = plugins
	}

	// 4. 构造 ApisixRoute CRD
	route := &apisixv2.ApisixRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels:    req.Labels,
		},
		Spec: apisixv2.ApisixRouteSpec{
			IngressClassName: "apisix", // Note: Don't Modify
			HTTP:             []apisixv2.ApisixRouteHTTP{httpRoute},
		},
	}

	// 创建 ApisixRoute CRD，Controller 会自动同步到 Apisix Gateway
	_, err := cluster.ApisixClient.ApisixV2().ApisixRoutes(req.Namespace).Create(ctx, route, metav1.CreateOptions{})
	if err != nil {
		return errors.Wrap(err, "创建 ApisixRoute 失败")
	}

	logx.Info("创建 Apisix HTTP 路由成功")
	return nil
}

// DeleteRoute 删除 Apisix HTTP 路由
func (m *ApisixService) DeleteRoute(ctx context.Context, req *apisixReq.DeleteRouteReq) error {
	// 根据ClusterId获取集群客户端
	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(req.ClusterId)
	if cluster == nil {
		return errors.Errorf("没有可用的K8s集群")
	}

	if cluster.ApisixClient == nil {
		return errors.Errorf("Apisix client 未初始化")
	}

	err := cluster.ApisixClient.ApisixV2().ApisixRoutes(req.Namespace).Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		return errors.Errorf("删除 ApisixRoute 失败: %v", err)
	}

	logx.Info("删除 Apisix HTTP 路由成功")
	return nil
}

// UpdateRoute 更新 Apisix HTTP 路由（删除后重建）
func (m *ApisixService) UpdateRoute(ctx context.Context, req *apisixReq.CreateRouteReq) error {
	// 先删除旧路由（忽略不存在的错误）
	deleteReq := &apisixReq.DeleteRouteReq{
		Name:      req.Name,
		Namespace: req.Namespace,
		ClusterId: req.ClusterId,
	}
	if err := m.DeleteRoute(ctx, deleteReq); err != nil {
		logx.Error("删除旧路由失败（可能不存在）", err)
	}

	// 创建新路由
	if err := m.CreateRoute(ctx, req); err != nil {
		return errors.Errorf("更新路由失败: %v", err)
	}

	logx.Info("更新 Apisix HTTP 路由成功")
	return nil
}

// CreateStreamRoute 创建 Apisix TCP Stream 路由（用于 SSH 等 TCP 服务）
// 注意：需要 Apisix 配置中启用 stream_proxy，并且 Apisix-ingress-controller 支持 ApisixRoute 的 stream 字段
func (m *ApisixService) CreateStreamRoute(ctx context.Context, req *apisixReq.CreateStreamRouteReq) error {
	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(req.ClusterId)
	if cluster == nil {
		return errors.Errorf("没有可用的K8s集群")
	}

	if cluster.ApisixClient == nil {
		return errors.Errorf("Apisix client 未初始化")
	}

	// 构造后端服务地址（跨 namespace 需要使用 FQDN）
	backendServiceName := req.ServiceName
	if req.ServiceNamespace != "" && req.ServiceNamespace != req.Namespace {
		backendServiceName = fmt.Sprintf("%s.%s.svc.cluster.local", req.ServiceName, req.ServiceNamespace)
	}

	// 构建 TCP Stream 路由配置
	// 使用 ApisixRoute 的 stream 字段来配置 TCP 路由
	streamRoute := apisixv2.ApisixRouteStream{
		Name:     "tcp-access",
		Protocol: "TCP",
		Match: apisixv2.ApisixRouteStreamMatch{
			IngressPort: int32(req.IngressPort), // Apisix 内部监听的端口 (9100)
		},
		Backend: apisixv2.ApisixRouteStreamBackend{
			ServiceName: backendServiceName,              // sshpiper
			ServicePort: intstr.FromInt(req.ServicePort), // kubeflow
		},
	}

	// 构造 ApisixRoute CRD（包含 stream 配置）
	route := &apisixv2.ApisixRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels:    req.Labels,
		},
		Spec: apisixv2.ApisixRouteSpec{
			IngressClassName: "Apisix",
			Stream:           []apisixv2.ApisixRouteStream{streamRoute},
		},
	}

	// 创建 ApisixRoute CRD
	_, err := cluster.ApisixClient.ApisixV2().ApisixRoutes(req.Namespace).Create(ctx, route, metav1.CreateOptions{})
	if err != nil {
		return errors.Errorf("创建 ApisixRoute (Stream) 失败: %v", err)
	}

	logx.Info("创建 Apisix TCP Stream 路由成功")
	return nil
}

// DeleteStreamRoute 删除 Apisix TCP Stream 路由
func (m *ApisixService) DeleteStreamRoute(ctx context.Context, req *apisixReq.DeleteStreamRouteReq) error {
	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(req.ClusterId)
	if cluster == nil {
		return errors.Errorf("没有可用的K8s集群")
	}

	if cluster.ApisixClient == nil {
		return errors.Errorf("Apisix client 未初始化")
	}

	err := cluster.ApisixClient.ApisixV2().ApisixRoutes(req.Namespace).Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		return errors.Errorf("删除 ApisixRoute (Stream) 失败: %v", err)
	}

	logx.Info("删除 Apisix TCP Stream 路由成功")
	return nil
}

// AuthApiSix Apisix forward-auth 认证逻辑（Notebook / TensorBoard / Inference）
func (s *ApisixService) AuthApiSix(ctx context.Context, req *apisixReq.AuthApisixReq) (*response.AuthApisixResp, error) {
	// 1. 解析路径，提取 资源类型、namespace 和 实例名称
	resourceType, namespace, instanceName, err := parseApisixPath(req.OriginalUri)
	if err != nil {
		return nil, errors.Wrap(err, "解析路径失败")
	}

	// 2. 推理服务支持 API-Key 认证，优先检查
	if resourceType == consts.InferenceInstance && req.ApiKey != "" {
		return s.authInferenceByApiKey(ctx, req.ApiKey, namespace, instanceName)
	}

	// 3. JWT Token 认证
	if req.Token == "" {
		return nil, errors.New("缺少认证凭据")
	}

	j := utils.NewJWT()
	claims, err := j.ParseToken(req.Token)
	if err != nil {
		return nil, errors.New("Token 无效")
	}

	// 4. 验证资源所有者
	switch resourceType {
	case consts.NotebookInstance:
		if claims.Namespace != namespace {
			return nil, errors.New("用户无权访问该资源")
		}
		var notebook nbModel.Notebook
		if err := global.GVA_DB.Where("instance_name = ? AND namespace = ?", instanceName, namespace).First(&notebook).Error; err != nil {
			return nil, errors.New("Notebook 不存在")
		}
		if notebook.UserId != claims.BaseClaims.ID {
			return nil, errors.New("用户不是 Notebook 所有者")
		}
	case consts.TensorBoardInstance:
		if claims.Namespace != namespace {
			return nil, errors.New("用户无权访问该资源")
		}
		var tb tensorboardModel.Tensorboard
		tbInstanceName := instanceName + "-tb"
		if err := global.GVA_DB.Where("instance_name = ? AND namespace = ?", tbInstanceName, namespace).First(&tb).Error; err != nil {
			return nil, errors.New("Tensorboard 不存在")
		}
		if tb.UserId != claims.BaseClaims.ID {
			return nil, errors.New("用户不是 Tensorboard 所有者")
		}
	case consts.InferenceInstance:
		var service inferenceModel.Inference
		if err := global.GVA_DB.Where("instance_name = ? AND namespace = ?", instanceName, namespace).First(&service).Error; err != nil {
			return nil, errors.New("推理服务不存在")
		}
		if service.UserId != claims.BaseClaims.ID {
			return nil, errors.New("用户不是推理服务所有者")
		}
	default:
		return nil, errors.New("不支持的资源类型鉴权")
	}

	return &response.AuthApisixResp{
		UserID:    claims.BaseClaims.ID,
		Namespace: claims.Namespace,
		Token:     req.Token,
	}, nil
}

// authInferenceByApiKey 推理服务 API-Key 认证
func (s *ApisixService) authInferenceByApiKey(_ context.Context, apiKey, namespace, serviceName string) (*response.AuthApisixResp, error) {
	var key inferenceModel.InferenceApiKey
	if err := global.GVA_DB.Where("api_key = ?", apiKey).First(&key).Error; err != nil {
		return nil, errors.New("无效的 API Key")
	}

	if key.Status != inferenceModel.ApiKeyStatusActive {
		return nil, errors.New("API Key 已禁用")
	}

	if key.ExpiredAt != nil && key.ExpiredAt.Before(time.Now()) {
		return nil, errors.New("API Key 已过期")
	}

	var service inferenceModel.Inference
	if err := global.GVA_DB.Where("id = ?", key.ServiceId).First(&service).Error; err != nil {
		return nil, errors.New("推理服务不存在")
	}

	if service.InstanceName != serviceName || service.Namespace != namespace {
		return nil, errors.New("API Key 与请求的服务不匹配")
	}

	// 更新最后使用时间
	now := time.Now()
	global.GVA_DB.Model(&key).Update("last_used_at", &now)

	return &response.AuthApisixResp{
		UserID:    key.UserId,
		Namespace: namespace,
		RateLimit: key.RateLimit,
	}, nil
}

// parseApisixPath 解析 Apisix 路径
// 路径格式: /notebook/{namespace}/{notebook-name}/... 或 /tensorboard/{namespace}/{job-name}/...
func parseApisixPath(uri string) (resourceType, namespace, instanceName string, err error) {
	if idx := strings.Index(uri, "?"); idx != -1 {
		uri = uri[:idx]
	}
	parts := strings.Split(strings.Trim(uri, "/"), "/")
	if len(parts) < 3 {
		return "", "", "", errors.New("invalid Apisix path format")
	}

	resourceType = parts[0]
	if resourceType != consts.NotebookInstance && resourceType != consts.TensorBoardInstance && resourceType != consts.InferenceInstance {
		return "", "", "", errors.New("unsupported resource type in path: " + resourceType)
	}

	return resourceType, parts[1], parts[2], nil
}

# 从 APISIX 迁移到 Kong 网关方案文档

## 1. 背景与目标

当前项目使用 Apache APISIX 作为 Kubernetes 的 Ingress Controller。但在实践中发现：
1. **架构复杂性**：APISIX 依赖单独的 etcd 集群，维护成本较高（如遇到 `ImagePullBackOff` 等环境问题）。
2. **兼容性与坑**：APISIX v2.0+ 引入了 `GatewayProxy` CRD，默认配置容易导致 `no GatewayProxy configs provided` 等同步问题。
3. **运维工具链**：Helm 安装和 `proxy_mode` 存在转义等小坑。

**迁移目标**：使用 **Kong (DB-less mode)** 替代 APISIX，实现更轻量（无需独立数据库）、更稳定、社区生态更完善的云原生网关方案。同时确保完全支持我们目前的三个核心诉求：
1. **HTTP/WebSocket 代理**：承载 Notebook 和 TensorBoard 流量。
2. **TCP 流代理**：承载 SSHPiper 的 SSH `22` 端口流量。
3. **动态 Header 改写与外部鉴权**：通过 Lua 脚本和自定义鉴权接口（forward-auth）拦截非法请求并将 Token 注入 Cookie。

---

## 2. 架构与资源映射对比

在代码层面（Go 后端操作 K8s 资源），我们需要将生成 APISIX 的 CRD 替换为生成 Kong 的周边资源：

| 需求场景 | APISIX 目前方案 | Kong 替代方案 | 备注 |
| :--- | :--- | :--- | :--- |
| **基础 HTTP 路由** | `ApisixRoute` (HTTP) | 原生 K8s `Ingress` | Kong KIC 完美支持标准 Ingress，更易维护。 |
| **WebSocket** | `ApisixRoute` (websocket: true) | 原生 K8s `Ingress` | Kong 默认支持升级至 WebSocket。 |
| **TCP 路由 (SSH)** | `ApisixRoute` (Stream) | `TCPIngress` (Kong CRD) | 需要在 Helm 安装 Kong 时暴露 Stream (22) 端口。 |
| **URL 重写** | `proxy-rewrite` 插件 | 注解 `konghq.com/rewrite` | TensorBoard 的剥离前缀可以通过标准 Ingress annotation 完成。 |
| **外部鉴权** | `forward-auth` 插件 | `KongPlugin` (定制 OAuth2/OIDC 或 Serverless Plugin) | 免费版最简单的做法是用 Lua (pre-function) 调内部 HTTP 接口。 |
| **Cookie 注入** | `serverless-post-function` (Lua) | `KongPlugin` (post-function Lua) | Lua 语法类似，从 `core.request.header` 改为 `kong.request.get_header` 等。 |

---

## 3. 环境迁移步骤

### 3.1 清理旧的 APISIX
为避免路由冲突，需先完全卸载 APISIX 以及相关路由：
```bash
# 1. 卸载 APISIX Helm Release
helm uninstall apisix -n apisix
kubectl delete ns apisix

# 2. 删除已经动态生成的 ApisixRoute / CRDs (如不需要保留)
kubectl delete apisixroutes -n hmz --all
kubectl delete apisixroutes -n kubeflow stream-sshpiper
kubectl delete crd apisixroutes.apisix.apache.org # 可选
```

### 3.2 部署 Kong (DB-less 模式)

使用 Helm 部署 Kong，并开启 Ingress 控制器和 TCP 暴露：

1. **添加 Helm 仓库**
```bash
helm repo add kong https://charts.konghq.com
helm repo update
```

2. **准备 `kong-values.yaml`**
创建一个精简的 values 文件用于部署 db-less 模式并暴露 tcp `22` 端口：
```yaml
# kong-values.yaml
env:
  database: "off" # 无数据库模式
proxy:
  type: LoadBalancer
  stream:
    - containerPort: 22
      servicePort: 22
      protocol: TCP
ingressController:
  enabled: true
  installCRDs: false
```

3. **创建 Namespace 并安装 Kong**
```bash
kubectl create namespace kong
helm install kong kong/kong -n kong -f kong-values.yaml
```

---

## 4. 后端 Go 代码改造示例 (CRD 生成)

将后端自动管理网关生命周期的代码（即 `api/v1/notebook` 下动态拉起/销毁路由的部分）进行对应改造：

### 4.1 Notebook 原生 Ingress 示例
我们不再拼接 `ApisixRoute` 的结构体，而是直接创建标准的 `networking.k8s.io/v1 Ingress`：

```yaml
# 原来使用 APISIX CRD，现改为原生 Ingress
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: route-notebook-86339a
  namespace: hmz
  annotations:
    kubernetes.io/ingress.class: "kong"
    # 挂载我们稍后创建的鉴权和 Header 处理插件
    konghq.com/plugins: "notebook-auth, notebook-set-cookie" 
spec:
  rules:
  - host: localhost
    http:
      paths:
      - path: /notebook/hmz/notebook-86339a/
        pathType: Prefix
        backend:
          service:
            name: notebook-86339a
            port:
              number: 80
```

### 4.2 TensorBoard 路由分离重写
Tensorboard 需要重写路径前缀，Kong 可以直接用注解解决：
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: route-tb-notebook-86339a
  namespace: hmz
  annotations:
    kubernetes.io/ingress.class: "kong"
    konghq.com/rewrite: "/" # 网关会将匹配的 Prefix 替换为 /
spec:
  rules:
  - host: localhost
    http:
      paths:
      - path: /tensorboard/hmz/notebook-86339a/
        pathType: Prefix
        backend:
          service:
            name: notebook-86339a-tb
            port:
              number: 80
```

### 4.3 SSH TCP 流代理
使用 Kong 的 `TCPIngress` 替代 `ApisixRoute (stream)`：
```yaml
apiVersion: configuration.konghq.com/v1beta1
kind: TCPIngress
metadata:
  name: stream-sshpiper
  namespace: kubeflow
  annotations:
    kubernetes.io/ingress.class: kong
spec:
  rules:
  - port: 22
    backend:
      serviceName: sshpiper
      servicePort: 22
```

### 4.4 Lua 插件平替 (鉴权与注入 Cookie)
原来分散在 APISIX 配置中的 Lua 代码，提取成复用的 `KongPlugin` CRD（你在创建 Instance 时可以注入，或在 Namespace 级别创建好供各路由引用）：

#### 动态鉴权与 Cookie 注入合并 (通过 Serverless Plugin / Pre-function)
```yaml
apiVersion: configuration.konghq.com/v1
kind: KongPlugin
metadata:
  name: notebook-auth-and-cookie
  namespace: hmz
config:
  functions:
    - |
      -- 请求拦截：发往后端的 Auth 接口 (同原 forward-auth)
      local http = require "resty.http"
      local httpc = http.new()
      local uri = "http://host.docker.internal:8001/aiInfra/api/v1/notebook/auth"
      
      -- 透传 token 可以从 query, header 或 cookie 获取
      local token = kong.request.get_header("x-token")
      
      -- 这里执行鉴权逻辑... (简化版)
      local res, err = httpc:request_uri(uri, {
          method = "GET",
          headers = {
            ["x-token"] = token,
            ["X-Original-URI"] = kong.request.get_path()
          }
      })
      if not res or res.status ~= 200 then
         return kong.response.exit(401, 'Unauthorized')
      end

      -- 如果后端返回了新的 Token (通过特定 header)，我们截获并种 Cookie
      local set_token_header = res.headers["x-set-token"]
      if set_token_header and set_token_header ~= "" then
         local cookie_val = "x-token=" .. set_token_header .. "; Path=/; Max-Age=604800"
         kong.response.add_header("Set-Cookie", cookie_val)
      end
plugin: pre-function
```

## 5. 常见问题 (FAQ)
*   **迁移过程中会导致原有的 Notebook 无法访问吗？**
    是的。卸载 APISIX 到 Kong 接管这段时间内，原本由网关暴露的端口（80/22）会短暂不可用，建议在维护窗口操作。
*   **需要重新拉取集群里的所有 Notebook 吗？**
    不需要销毁后端的 Pod / Deployment。你只需要用 Go 脚手架重新遍历一遍已有的实例库，用新逻辑生成相应的 `Ingress` / `TCPIngress` 即可让网关重新接管流量。

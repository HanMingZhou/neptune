# APISIX 网关

## 目录

- [一、整体架构](#一整体架构)
- [二、HTTP 流量](#二http-流量)
- [三、SSH 流量](#三ssh-流量)
- [四、快速开始](#四快速开始)
- [五、故障排查](#五故障排查)
- [六、本目录文件说明](#六本目录文件说明)

---

## 一、整体架构

APISIX 在本系统中承担**统一网关**的角色，处理所有进入集群的流量：

```
                  用户请求
          ┌──────────┴──────────┐
          │ HTTP (:80)          │ SSH (:22)
          ▼                     ▼
┌──────────────────────────────────────────────┐
│                APISIX Gateway                │
│                                              │
│  HTTP Proxy (:9080)     Stream Proxy (:9100) │
│  - Notebook (Jupyter)   - SSH 流量转发        │
│  - TensorBoard                               │
│  - 认证 (forward-auth)                        │
└──────────────────────────────────────────────┘
          │                     │
          ▼                     ▼
  Notebook Service         SSHPiper
  (JupyterLab / TB)      (SSH 路由分发)
                                │
                                ▼
                          Notebook Pod
                           (sshd:22)
```

### 核心组件

| 组件 | 作用 |
|------|------|
| **APISIX Gateway** | API 网关，处理 HTTP 和 TCP 流量 |
| **APISIX Ingress Controller** | 监听 K8s CRD，自动同步路由到 Gateway |
| **ApisixRoute CRD** | 声明式定义路由规则（HTTP 路由和 Stream 路由） |

---

## 二、HTTP 流量

### 2.1 Notebook 访问流程

```
浏览器: http://ai.local/notebook/caixukun/notebook-xxx/lab
                         │
                         ▼
                   APISIX Gateway
                         │
           ┌─────────────┴─────────────┐
           ▼                            ▼
   ① 路由匹配                      匹配失败 → 404
   /notebook/caixukun/notebook-xxx/*
           │
           ▼
   ② forward-auth 插件
   调用认证服务验证身份
           │
     ┌─────┴─────┐
     ▼           ▼
   200 OK    401/403
   (有权限)   (拒绝访问)
     │
     ▼
   ③ 转发请求到后端
   notebook-xxx.caixukun.svc:80
           │
           ▼
   Notebook Pod (Jupyter:8888)
```

> **关键点**：APISIX **不做路径重写**，原样透传请求路径给 Jupyter。因为 Jupyter 启动时配置了 `base_url=/notebook/caixukun/notebook-xxx/`，它自己知道如何处理带前缀的路径。

### 2.2 TensorBoard 访问流程

```
浏览器: http://ai.local/tensorboard/caixukun/notebook-xxx/
                         │
                         ▼
                   APISIX Gateway
                         │
                         ▼
   ① 路由匹配
   /tensorboard/caixukun/notebook-xxx/*
                         │
                         ▼
   ② proxy-rewrite 插件（路径重写）
   /tensorboard/cai xu kun/notebook-xxx/(.*) → /$1
                         │
                         ▼
   ③ 转发到后端
   notebook-xxx-tb.caixukun.svc:80
                         │
                         ▼
   TensorBoard Pod (监听 /)
```

> **关键点**：TensorBoard **没有** `base_url` 机制，只能处理根路径 `/`，因此 APISIX **必须 rewrite** 把 `/tensorboard/caixukun/notebook-xxx/` 前缀剥掉。这与 Jupyter 的处理方式**恰好相反**。

### 2.3 认证机制

使用 APISIX 的 **forward-auth** 插件，将认证逻辑委托给后端平台服务：

1. APISIX 将请求的 `Cookie`、`x-token`、`Authorization` 等 Header 转发给认证接口
2. 认证服务从中提取 Token，验证 JWT 有效性，并检查用户对该 Notebook 的访问权限
3. 认证成功：返回 200，APISIX 放行请求并通过 `serverless-post-function` 插件将 Token 写入 Cookie
4. 认证失败：返回 401/403，APISIX 直接拒绝请求

> ⚠️ **配置注意**：`auth-uri` 必须包含完整路径，不能省略 `/api/v1/` 前缀。

---

## 三、SSH 流量

SSH 使用 APISIX 的 **Stream Proxy** 功能处理 TCP 流量。

### 3.1 架构概述

APISIX 在 SSH 链路中只扮演**入口**角色 —— 把集群外的 TCP 流量引进来，转发给 SSHPiper。按用户名分发到具体 Pod 的工作由 SSHPiper 负责。

```
                  1 个 Stream Route              N 个 Pipe（每个 Notebook 一个）
                 ┌──────────────┐          ┌──────────────────────────────────┐
                 │              │          │                                  │
外部 :22 ──TCP→ APISIX ──TCP→ SSHPiper ──┬──. Pipe: user=caixukun-nb-aaa → Pod-A    │
                 │              │          ├── Pipe: user=caixukun-nb-bbb → Pod-B  │
                 │              │          └── Pipe: user=dev-nb-ccc → Pod-C  │
                 └──────────────┘          └──────────────────────────────────┘
```

| 资源 | 数量 | 作用 |
|------|------|------|
| **Stream Route** | 整个集群 **1 个** | APISIX 外部端口 22 (内部 9100) → SSHPiper Service 端口 22 (内部 2222) |
| **Pipe** | 每个 Notebook **1 个** | SSHPiper 内部路由：按用户名分发到对应 Pod |

Stream Route 是全局共享的，所有 Notebook 的 SSH 流量都走同一个入口。系统在创建第一个需要 SSH 的 Notebook 时自动创建此路由（幂等，已存在则跳过）。

> 关于 SSHPiper 的 Pipe 路由规则和认证机制，请参阅 [SSHPiper 文档](../sshpiper/readme.md)。

### 3.2 Stream Proxy 配置

**1. 启用 Stream Proxy（APISIX ConfigMap）**

```yaml
# kubectl edit configmap apisix -n apisix
apisix:
  proxy_mode: http&stream    # 同时启用 HTTP 和 Stream
  stream_proxy:
    only: false
    tcp:
      - 9100                 # Stream 内部监听端口
```

**2. 暴露 SSH 端口（Gateway Service）**

```yaml
# kubectl edit svc apisix-gateway -n apisix
spec:
  ports:
    - name: http
      port: 80
      targetPort: 9080
    - name: ssh
      port: 22               # 对外 SSH 端口
      targetPort: 9100       # 转发到 Stream Proxy
```

---

## 四、快速开始

### 4.1 安装 APISIX

```bash
# 添加 Helm 仓库
helm repo add apisix https://apache.github.io/apisix-helm-chart
helm repo update

# 安装 (针对 Ingress Controller 2.0+ 版本)
helm install apisix apisix/apisix \
  --create-namespace \
  --namespace apisix \
  --set ingress-controller.enabled=true \
  --set ingress-controller.gatewayProxy.createDefault=true \
  --set ingress-controller.gatewayProxy.provider.controlPlane.service.name=apisix-admin \
  --set ingress-controller.gatewayProxy.provider.controlPlane.service.port=9180 \
  --set ingress-controller.gatewayProxy.provider.controlPlane.auth.adminKey.value=edd1c9f034335f136f87ad84b625c8f1 \
  --set ingress-controller.config.apisix.serviceNamespace=apisix \
  --set apisix.proxy_mode=http\&stream \
  --set "apisix.stream_proxy.tcp[0]=9100"

# 重要：修复 proxy_mode 可能未生效的问题（因命令行 & 字符转义问题）
# 如果发现 SSH(Stream) 路由无法添加，请手动检查并修改 ConfigMap
kubectl get configmap apisix -n apisix -o jsonpath='{.data.config\.yaml}' | sed 's/proxy_mode: http/proxy_mode: http\&stream/' | kubectl create configmap apisix -n apisix --from-file=config.yaml=/dev/stdin --dry-run=client -o yaml | kubectl apply -f -
kubectl rollout restart deployment apisix -n apisix

# 添加 SSH 端口
kubectl patch svc apisix-gateway -n apisix --type='json' \
  -p='[{"op":"add","path":"/spec/ports/-","value":{"name":"ssh","port":22,"targetPort":9100,"protocol":"TCP"}}]'
```

### 4.2 验证安装

```bash
# 检查 Pod
kubectl get pods -n apisix
# 应该看到: apisix-xxx, apisix-ingress-controller-xxx, apisix-etcd-xxx

# 检查 Service
kubectl get svc -n apisix
```

### 4.3 本地开发测试

```bash
# HTTP 访问
kubectl port-forward svc/apisix-gateway -n apisix 8888:80 &
open http://localhost:8888/notebook/caixukun/notebook-xxx/lab

# SSH 访问
kubectl port-forward svc/apisix-gateway -n apisix 2222:22 &
ssh caixukun-notebook-xxx@localhost -p 2222
```

---

## 五、故障排查

### 5.1 HTTP 404

**症状**：访问 Notebook 返回 `404 page not found`

```bash
# 1. 检查 ApisixRoute 是否已创建
kubectl get apisixroute -n <namespace>

# 2. 检查路由是否已同步到 APISIX
kubectl port-forward svc/apisix-admin -n apisix 9180:9180 &
curl -s http://127.0.0.1:9180/apisix/admin/routes \
  -H "X-API-KEY: edd1c9f034335f136f87ad84b625c8f1" | jq '.list[].value.name'

# 3. 检查 Ingress Controller 日志
kubectl logs -n apisix deployment/apisix-ingress-controller -c manager --tail=50
# 常见错误: "no GatewayProxy configs provided"
# 原因：Ingress Controller 2.0+ 需要 GatewayProxy 资源来配置 Admin API 连接。
# 修复：确保安装时设置了 gatewayProxy.createDefault=true，或手动创建 GatewayProxy 资源。

### 5.3 路由未同步

如果 ApisixRoute 状态正常但 Admin API 查不到路由：
1. 检查 GatewayProxy 资源：`kubectl get gatewayproxy -n apisix`
2. 确认 `apisix-ingress-config` 中的配置。
3. 检查 `proxy_mode` 是否包含 `stream`：`kubectl exec deployment/apisix -n apisix -- cat /usr/local/apisix/conf/config.yaml | grep proxy_mode`
   如果显示 `proxy_mode: http`，请手动修复 ConfigMap 并重启 APISIX。

### 5.4 SSH 连接失败

```bash
# 1. 检查 Stream Route
kubectl get apisixroute -n kubeflow

# 2. 检查 SSHPiper Pod
kubectl get pods -n kubeflow -l app=sshpiper

# 3. 检查 Pipe CRD
kubectl get pipe -n kubeflow

# 4. 检查 APISIX 是否监听 9100 端口
kubectl exec -n apisix deployment/apisix -- netstat -tlnp | grep 9100

# 5. 检查 SSH Service
kubectl get svc -n <namespace> | grep ssh
```

---

## 六、本目录文件说明

| 文件 | 用途 | 使用方式 |
|------|------|----------|
| `readme.md` | 本文档 | 阅读 |
| `apisix-ingress-config-fix.yaml` | 修复 IngressClass 配置 | `kubectl apply -f` |
| `apisix-stream-proxy.yaml` | Stream Proxy 配置参考 | **仅参考，勿直接 apply** |
| `apisix-config-template.yaml` | APISIX 完整配置模板 | **仅参考，勿直接 apply** |

> ⚠️ ConfigMap 类的文件会完全覆盖现有配置，请通过 `kubectl edit` 或 `kubectl patch` 做增量修改。

---

## 参考资料

- [APISIX Ingress Controller 文档](https://apisix.apache.org/docs/ingress-controller/)
- [APISIX 官方文档](https://apisix.apache.org/docs/apisix/)
- [forward-auth 插件](https://apisix.apache.org/docs/apisix/plugins/forward-auth/)
- [Stream Proxy 配置](https://apisix.apache.org/docs/apisix/stream-proxy/)

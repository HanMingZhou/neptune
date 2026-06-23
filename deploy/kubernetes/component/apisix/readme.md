# APISIX 统一流量网关

本网关在 Neptune 平台中承担统一流量入口的角色，处理集群外部打入的所有 HTTP、WebSocket（L7）流量和 SSH TCP（L4）流量。

---

## 📖 目录

- [一、整体架构与双模机制](#一整体架构与双模机制)
- [二、核心组件与作用](#二核心组件与作用)
- [三、端口对应与流转链路](#三端口对应与流转链路)
- [四、后端配置参数详解 (config.dev.yaml / ConfigMap)](#四后端配置参数详解-configdevyaml--configmap)
- [五、部署流程与机制说明](#五部署流程与机制说明)
- [六、故障排查指南](#六故障排查指南)
- [七、本目录文件说明](#七本目录文件说明)

---

## 一、整体架构与双模机制

APISIX 在集群边缘运行，同时启用 **HTTP** 与 **Stream (TCP/UDP)** 双模代理：

```text
                            外部客户端请求
                     ┌─────────────┴─────────────┐
                     │ HTTP (Port 80)            │ SSH (Port 22 / NodePort)
                     ▼                           ▼
┌───────────────────────────────────────────────────────────────────────┐
│                          APISIX Gateway                               │
│                                                                       │
│  [HTTP/L7 引擎] (监听 :9080)             [Stream/L4 引擎] (监听 :9100)   │
│  - Neptune Web 前端                     - SSH TCP 隧道入口            │
│  - Neptune Server 后端 API             - 转发至 SSHPiper 服务         │
│  - Notebook/Jupyter (forward-auth 鉴权)                              │
│  - TensorBoard (proxy-rewrite 重写)                                   │
┌───────────────────────────────────────────────────────────────────────┐
                     │                           │
         ┌───────────┴───────────┐               ▼
         ▼                       ▼            SSHPiper
    Neptune Web           Neptune Server         │
    (Web 静态资源)          (Go API/Auth)         ▼
         │                                 Notebook Pod (sshd)
         ▼
    Notebook Pod
    (Jupyter:8888)
```

### 1. HTTP/HTTPS (L7) 模式
*   处理 Web 前端、后端 API 访问。
*   提供租户鉴权功能（配合 `forward-auth` 插件调用 Neptune Server 后端进行 JWT/API-Key 校验）。
*   执行 URL 路径改写（配合 `proxy-rewrite` 插件进行路径剥离，方便不带 BaseURL 的服务使用）。

### 2. Stream (L4) 模式
*   直接在 TCP 层面进行字节流路由（不解析 HTTP/HTTPS 协议），用于代理 SSH 流量。
*   APISIX 仅作为流量的中转站，将外部发往网关 SSH 端口的 TCP 流，透明转发到 `sshpiper` 服务上，具体的 SSH 密钥校验、用户名解析由 `sshpiper` 完成。

---

## 二、核心组件与作用

在 Kubernetes 部署中，网关由以下几个关键组件协同工作：

| 组件 | 角色类型 | 作用描述 |
|:---|:---|:---|
| **APISIX Pod (Deployment)** | 数据面 (Data Plane) | 真正运行网关路由、执行插件、转发流量的容器进程（基于 OpenResty/Nginx）。 |
| **apisix-gateway (Service)** | 流量入口 (Entrypoint) | K8s Service 资源，对外暴露网关的端口（如 `80`, `22`），并映射到容器内端口。 |
| **Ingress Controller (Deployment)** | 控制面 (Control Plane) | 监听 Kubernetes API，将 `ApisixRoute` 等声明式 YAML 翻译成 APISIX 能读懂的 JSON 规则，写入 APISIX 内存。 |
| **ApisixRoute (CRD)** | 路由声明配置 (Rule) | 声明式路由定义，包括匹配条件（Host/Path/IngressPort）、目标后端及绑定的插件列表。 |

### ApisixRoute 的分类与同步机制
1.  **静态路由 (全局共享)**：例如 SSH 的流路由 `stream-sshpiper`，在集群部署时通过 YAML 一次性部署。
2.  **动态路由 (租户隔离)**：例如 Notebook 的 HTTP 访问路由，在用户启动实例时由 **Neptune 后端代码**调用 Kubernetes Client 自动创建，销毁时自动删除。

---

## 三、端口对应与流转链路

理解端口的流转是理解网关转发、解决“无法连接”的关键。下面整理了端口速查以及完整的流量流转链路。

### 1. 端口速查表

| 端口类型 | 物理/默认端口 | 链路作用与说明 |
|:---|:---|:---|
| **外部访问端口** | `32485` (HTTP) <br> `30177` (SSH) | **NodePort 模式下，用户直接在浏览器或 SSH 客户端输入的物理端口。** |
| **Service 暴露端口** | `80` (HTTP) <br> `22` (SSH) | `apisix-gateway` Service 对外声明的端口号。 |
| **Pod 实际监听端口** | `9080` (HTTP) <br> `9100` (SSH TCP) | **APISIX 容器内进程真正绑定监听的端口（TargetPort）。** |
| **SSHPiper 服务端口** | `22` (Service) <br> `2222` (Pod) | SSHPiper 对外暴露的端口以及其容器内部的监听端口。 |
| **Notebook 服务端口** | `80` (Jupyter HTTP) <br> `22` (sshd TCP) | K8s 内网的服务端口，以及具体 Pod 容器内部的端口。 |

### 2. 流量流转链路

#### 🌐 HTTP / 访问 Notebook 链路
```text
[浏览器] http://10.255.141.8:32485/notebook/caixukun/notebook-1/lab
  │
  ▼ (NodePort 映射)
[K8s Node] 监听端口 32485 
  │
  ▼ (Service 端口重定向)
[Service] apisix-gateway (Port: 80 -> targetPort: 9080)
  │
  ▼ (进入 Pod)
[APISIX Pod] 监听 9080 端口 (OpenResty)
  │
  ├─► [执行插件] forward-auth -> 调用 Neptune 后端进行 JWT/Cookie 权限验证 (200 放行 / 401 拦截)
  │
  ▼ (按 ApisixRoute 规则转发)
[K8s Service] notebook-1.caixukun.svc:80
  │
  ▼
[Notebook Pod] 监听 8888 端口的 JupyterLab
```

#### 🔑 SSH 流量链路
```text
[SSH 客户端] ssh -p 30177 caixukun-notebook-1@10.255.141.8
  │
  ▼ (NodePort 映射)
[K8s Node] 监听端口 30177
  │
  ▼ (Service 端口重定向)
[Service] apisix-gateway (Port: 22 -> targetPort: 9100)
  │
  ▼ (进入 Pod)
[APISIX Pod] 监听 9100 端口 (Stream TCP)
  │
  ▼ (按 ApisixRoute 规则匹配 ingressPort: 9100，无需域名，原样透传)
[K8s Service] sshpiper.kubeflow.svc:22 (对应 Pod 容器端口 2222)
  │
  ▼
[SSHPiper Pod] 接收到连接，识别登录用户名 `caixukun-notebook-1` 
  │
  ├─► 匹配 Pipe (CRD) 规则获取目标后端: notebook-1-ssh.caixukun.svc.cluster.local:22
  │
  ▼ (代理 SSH 握手并转发)
[K8s Service] notebook-1-ssh.caixukun.svc:22
  │
  ▼
[Notebook Pod] 监听 22 端口的 sshd 进程，完成登录
```

---

## 四、后端配置参数详解 (config.dev.yaml / ConfigMap)

平台后端的配置文件中有几个与网关及端口直接相关的参数。在不同部署环境下，必须根据网关的实际暴露端口修改这些配置：

```yaml
apisix:
    enabled: true
    base-domain: "localhost"
    auth-enabled: true
    auth-uri: "http://localhost:8001/aiInfra/api/v1/apisix/auth"
    stream-port: 9100

sshpiper:
    host: "127.0.0.1"
    port: 22
```

### 1. `apisix.stream-port` (默认 9100)
*   **使用位置**：Go 后端在通过 APISIX Ingress 客户端创建 TCP 流路由（`stream-sshpiper`）时，会将此端口填入 `stream-port` 字段：
    ```go
    stream-port: stream-port // 对应 apisix.stream-port
    ```
*   **为什么是这个值**：它是 APISIX Pod 内部流代理监听的端口。无论外部暴露端口怎么变，它**必须与 APISIX 容器内配置的 Stream 端口（9100）保持完全一致**。

### 2. `sshpiper.port` (默认 22)
*   **使用位置**：Go 后端在获取 Notebook 详情、生成“一键复制 SSH 登录命令”时，拼接命令行端口：
    ```go
    item.SSHCommand = fmt.Sprintf("ssh -p %d %s@%s", sshPort, sshUser, sshHost)
    ```
*   **配置原则**：**该参数纯粹是为了给前端控制台生成“一键复制 SSH 命令”的文本展示，对集群内网流量路由完全没有影响。** 它必须与**外部用户实际连接网关 SSH 的物理端口**一致。
    *   *LoadBalancer/域名映射场景*：如果外部防火墙将 `22` 映射到了网关，则填 `22`。
    *   *NodePort 直连接入场景*：必须修改为实际暴露的 NodePort 端口，例如 **`30177`**。如果不修改，用户在前端复制的连接命令端口为 22，将无法连通。

### 3. `sshpiper.host` 与 `apisix.base-domain`
*   **`sshpiper.host`**：SSH 客户端登录的 IP 或域名。**同样仅用于前端一键复制命令的文本展示，不影响内网路由。** 如果为空，则自动降级回退使用 `apisix.base-domain`。
*   **`apisix.base-domain`**：APISIX 对外服务的总入口域名或 IP。Go 后端在动态创建 `ApisixRoute` (HTTP) 时会读取它，并作为 HTTP 规则的 `Host` 头进行限制匹配（防未绑定域名泛解析），同时为用户拼接前端访问链接。

---

## 五、部署流程与机制说明

`deploy_all.sh` 脚本在部署 APISIX 时，会在 Kubernetes 中进行以下初始化：

1.  **开启流代理功能**：
    向 Helm 传入选项，在 APISIX 配置中注入 `proxy_mode: http&stream` 及流 TCP 监听列表：
    ```bash
    --set 'apisix.proxy_mode=http&stream' \
    --set "apisix.stream_proxy.tcp[0]=9100"
```
2.  **清理历史冲突**：
    如果集群里遗留了旧版 `GatewayProxy apisix` 与当前 Helm 生成的 `apisix-config` 共存导致 Admin API 无法同步，脚本会自动删掉旧版冲突资源。
3.  **修复配置格式**：
    由于 Helm 在处理特殊合并语法时可能导致 `proxy_mode` YAML 格式出错崩溃，脚本通过 `fix_apisix_proxy_mode_config` 自动获取 ConfigMap、校验格式，并重启 APISIX Deployment。
4.  **向 Service 注入 SSH 端口规则**：
    APISIX 官方 Helm 默认不暴露 `22`，脚本通过 `kubectl patch` 为 `apisix-gateway` Service 注入端口定义，将外部端口 `22` 的 TCP 流量绑定至 Pod 容器端口 `9100` (`targetPort: 9100`)。

---

## 六、故障排查指南

### 1. HTTP 访问资源（Notebook 等）报 404
*   **排查 1**：K8s 中 `ApisixRoute` 资源是否创建成功？
    ```bash
    kubectl get apisixroute -n <namespace>
    ```
*   **排查 2**：查看 Ingress Controller 容器日志，检查是否提示 `no GatewayProxy configs provided` 或 `gateway proxy configuration conflict`。
    如果存在冲突，需清理多余的 GatewayProxy：
    ```bash
    kubectl get gatewayproxy -n apisix
    kubectl delete gatewayproxy apisix -n apisix # 清理旧版默认资源
    ```
*   **排查 3**：后端平台配置文件中的 `apisix.auth-uri` 是否可以从 APISIX Ingress Controller 内部解析访问？
    （注：此地址不能填 `localhost`，必须是 APISIX 能解析的 K8s ClusterIP，如 `http://neptune-server.neptune.svc.cluster.local:8888/...`）。

### 2. SSH 连不上 (排查链路)
当 SSH 连接报 Timeout 或 Refused 时，请按流量经过的各层逐步排查：

*   **第 1 步 (Service层)**：检查 `apisix-gateway` Service 中是否包含 TCP `22 -> 9100` 的端口段：
    ```bash
    kubectl get svc -n apisix apisix-gateway -o yaml
    ```
*   **第 2 步 (Pod监听)**：确认 APISIX 进程是否真的在容器内监听 `9100`：
    ```bash
    kubectl exec -n apisix deployment/apisix -- netstat -tlnp | grep 9100
    ```
*   **第 3 步 (路由层)**：检查 `apisix-admin` 运行时流路由列表是否为空（不为空说明控制面同步正常）：
    ```bash
    kubectl port-forward -n apisix svc/apisix-admin 9180:9180
    curl -s http://127.0.0.1:9180/apisix/admin/stream_routes -H "X-API-KEY: edd1c9f034335f136f87ad84b625c8f1"
    ```
*   **第 4 步 (SSHPiper层)**：查看 SSHPiper 容器的实时日志，观察是否捕获到了连接握手，是否成功读取了 Pipe：
    ```bash
    kubectl logs -n kubeflow deployment/sshpiper --tail=100 -f
    ```
*   **第 5 步 (目标Pod)**：如果 SSHPiper 报 `connection refused to backend`，检查目标 Notebook 的 SSH Service 是否就绪，以及容器内的 sshd 服务是否在运行。

---

## 七、本目录文件说明

*   `readme.md`：本说明文档。
*   `apisix-ingress-config-fix.yaml`：修复 IngressClass 与 Controller 连接的配置。
*   `apisix-stream-proxy.yaml`：APISIX TCP 流路由配置的模版（声明监听 9100 转发至 sshpiper:22）。
*   `neptune-platform-route.yaml`：Neptune 平台控制台 Web/API 的统一 HTTP 入口路由配置。
*   `apisix-config-template.yaml`：运行参数配置，供参考。

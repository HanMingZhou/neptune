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

APISIX 在本系统中承担统一网关的角色，处理所有进入集群的流量：

```text
                  用户请求
          ┌──────────┴──────────┐
          │ HTTP (:80)          │ SSH (:22)
          ▼                     ▼
┌──────────────────────────────────────────────┐
│                APISIX Gateway                │
│                                              │
│  HTTP Proxy (:9080)     Stream Proxy (:9100) │
│  - Notebook (Jupyter)   - SSH 流量入口        │
│  - TensorBoard                               │
│  - 认证 (forward-auth)                       │
└──────────────────────────────────────────────┘
          │                     │
          ▼                     ▼
  Notebook Service         SSHPiper Service
  TensorBoard Service           │
                                ▼
                          Notebook SSH Service
                                │
                                ▼
                           Notebook Pod
                            (sshd:22)
```

### 核心组件

| 组件 | 作用 |
|------|------|
| APISIX Gateway | API 网关，处理 HTTP 和 TCP 流量 |
| APISIX Ingress Controller | 监听 K8s CRD，自动同步路由到 Gateway |
| ApisixRoute CRD | 声明式定义路由规则，包括 HTTP 路由和 Stream 路由 |
| GatewayProxy | 告诉 Ingress Controller / ADC 应该连接哪个 APISIX Admin API |

---

## 二、HTTP 流量

### 2.1 Notebook 访问流程

```text
浏览器: http://ai.local/notebook/caixukun/notebook-xxx/lab
                         │
                         ▼
                   APISIX Gateway
                         │
           ┌─────────────┴─────────────┐
           ▼                            ▼
   ① 路由匹配                      匹配失败 -> 404
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

关键点：APISIX 不做路径重写，原样透传请求路径给 Jupyter。因为 Jupyter 启动时已经配置了 `base_url=/notebook/<namespace>/<instance>/`。

### 2.2 TensorBoard 访问流程

```text
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
   ② proxy-rewrite 插件
   /tensorboard/caixukun/notebook-xxx/(.*) -> /$1
                         │
                         ▼
   ③ 转发到后端
   notebook-xxx-tb.caixukun.svc:80
                         │
                         ▼
   TensorBoard Pod (监听 /)
```

关键点：TensorBoard 没有 `base_url` 机制，所以 APISIX 必须先做 rewrite，把前缀剥掉。

### 2.3 认证机制

使用 APISIX 的 `forward-auth` 插件，将认证逻辑委托给后端平台服务：

1. APISIX 将 `Cookie`、`x-token`、`Authorization` 等 Header 转发给认证接口。
2. 认证服务验证 JWT，并检查用户是否有权访问目标 Notebook / TensorBoard。
3. 认证成功返回 200，APISIX 放行请求。
4. 认证失败返回 401/403，APISIX 直接拒绝请求。

注意：`auth-uri` 必须包含完整路径，不能省略 `/api/v1/` 前缀。

---

## 三、SSH 流量

SSH 使用 APISIX 的 Stream Proxy 功能处理 TCP 流量。

### 3.1 先记住职责边界

APISIX 在 SSH 链路里只负责两件事：

1. 接住集群外用户打进来的 SSH TCP 流量。
2. 把这股流量转发给 `sshpiper` 这个 Kubernetes Service。

APISIX 不负责：

- 按用户名决定目标 Notebook
- 校验 SSH 公钥
- 透传密码
- 直接连接具体 Notebook Pod

这些都由 SSHPiper 处理。所以可以简单记成：

- APISIX 解决“流量怎么进集群”
- SSHPiper 解决“进来以后该转给谁”

还有一条在排障时很容易忽略的“控制面链路”：

```text
代码创建 ApisixRoute CR
        │
        ▼
Ingress Controller / ADC 读取 CR
        │
        ▼
根据 GatewayProxy 找到 apisix-admin:9180
        │
        ▼
把配置写入 APISIX 运行时 (/apisix/admin/stream_routes)
```

也就是说，`kubectl get apisixroute` 能看到资源，只能说明 Kubernetes 里有这条 CR；还要再看 APISIX Admin API 里的 `/stream_routes`，才能确认数据面真的生效。

### 3.2 端口速查表

最容易让人混乱的是 `22 / 9100 / 22 / 2222 / 22`。这些端口并不冲突，它们属于不同层：

| 所在层 | 资源 | 端口 | 含义 |
|------|------|------|------|
| 集群外 | 用户访问地址 | `22` | 用户执行 `ssh user@host -p 22` 时看到的端口 |
| K8s Service | `apisix-gateway` | `22 -> 9100` | Gateway Service 对外暴露 `22`，转发到 APISIX Pod 内部的 stream proxy `9100` |
| APISIX Pod | Stream Proxy | `9100` | APISIX 内部真正监听 SSH TCP 流量的端口；`ApisixRoute.spec.stream.match.ingressPort` 必须填它 |
| K8s Service | `sshpiper.kubeflow.svc` | `22 -> 2222` | APISIX 把流量转给 `sshpiper:22`，再由 Service 转到 sshpiper Pod 的 `2222` |
| K8s Service | `notebook-xxx-ssh.<ns>.svc` | `22 -> 22` | SSHPiper 根据 Pipe 规则，把流量再转给目标 Notebook 的 SSH Service |
| Notebook Pod | `sshd` | `22` | 容器内部最终处理登录的 sshd |

因此两个很容易混淆的字段要这样理解：

- `ingressPort=9100`
  这是 APISIX 内部 stream proxy 的监听端口
- `servicePort=22`
  这是 APISIX 要访问的后端 Service 端口

它们不是同一层，所以不应该填成一样。

### 3.3 完整链路图

```text
用户 ssh 到外部地址
ssh zzz-notebook-eeb47f@ai.example.com -p 22
        │
        ▼
① apisix-gateway Service
   port: 22
   targetPort: 9100
        │
        ▼
② APISIX Pod stream proxy
   listen: 9100
   ApisixRoute.stream.match.ingressPort = 9100
   backend = sshpiper:22
        │
        ▼
③ sshpiper Service (kubeflow)
   port: 22
   targetPort: 2222
        │
        ▼
④ sshpiper Pod
   listen: 2222
   根据 Pipe 的 from.username 查找目标后端
        │
        ▼
⑤ notebook-eeb47f-ssh Service (zzz)
   port: 22
   targetPort: 22
        │
        ▼
⑥ notebook Pod
   sshd listen: 22
```

从这张图里可以直接得出：

- APISIX Stream Route 只到 `sshpiper:22`
- Pipe 只在 SSHPiper 这一层生效
- Notebook 的 `-ssh` Service 和 Pod 内 `sshd` 是否正常，决定最后能不能真正登录

### 3.4 基于实际资源的链路样例

下面这组资源来自一次真实排障，适合和上面的抽象链路图对着看：

```bash
kubectl get svc -n apisix apisix-gateway
```

```text
NAME             TYPE       CLUSTER-IP      PORT(S)
apisix-gateway   NodePort   10.233.27.27   80:32485/TCP,22:30177/TCP
```

```bash
kubectl get svc -n kubeflow sshpiper
```

```text
NAME       TYPE        CLUSTER-IP    PORT(S)
sshpiper   ClusterIP   10.233.6.39   22/TCP
```

```bash
kubectl get pipes.sshpiper.com -n kubeflow
```

```text
NAME                       FROMUSER              TOUSER   TOHOST
pipe-zzz-notebook-eeb47f   zzz-notebook-eeb47f   root     notebook-eeb47f-ssh.zzz.svc.cluster.local:22
```

```bash
kubectl get gatewayproxy -n apisix apisix-config -o yaml
```

```yaml
spec:
  provider:
    controlPlane:
      service:
        name: apisix-admin
        port: 9180
```

把这些资源串起来后，真实业务链路就是：

```text
Windows:
  ssh -p 30177 zzz-notebook-eeb47f@10.255.141.8
        │
        ▼
apisix-gateway NodePort 30177
        │
        ▼
apisix-gateway Service port 22
        │
        ▼
APISIX Pod stream proxy 9100
        │
        ▼
kubeflow/sshpiper Service 10.233.6.39:22
        │
        ▼
sshpiper Pod :2222
        │
        ▼
Pipe: pipe-zzz-notebook-eeb47f
  from.username = zzz-notebook-eeb47f
  to.host = notebook-eeb47f-ssh.zzz.svc.cluster.local:22
  to.user = root
        │
        ▼
zzz/notebook-eeb47f-ssh:22
        │
        ▼
zzz/notebook-eeb47f-0:sshd
```

这里有两个非常关键的观察点：

- `Pipe` 在 `kubeflow` namespace，不在 `zzz`，也不在 `apisix`
- `GatewayProxy` 里只有 `9180`，这是控制面 Admin API 端口，不是业务 SSH 流量端口

### 3.5 需要创建哪些资源

| 资源 | 数量 | 所在 namespace | 作用 |
|------|------|----------------|------|
| APISIX Stream Route | 整个集群通常 1 个 | `kubeflow` | 把 APISIX 内部 `9100` 的 TCP 流量转给 `sshpiper:22` |
| Pipe | 每个 Notebook 1 个 | `kubeflow` | 按用户名把 SSH 连接路由到对应 Notebook 的 `-ssh` Service |
| Notebook SSH Service | 每个 Notebook 1 个 | 用户 namespace | 给 SSHPiper 提供稳定的后端 DNS 地址 |

`Stream Route` 是全局共享的；`Pipe` 是每个 Notebook 单独一条。

关于 Pipe 的匹配规则、公钥认证和密码透传，请参阅 [SSHPiper 文档](../sshpiper/readme.md)。

### 3.6 Stream Proxy 配置

1. 启用 Stream Proxy。

```yaml
# kubectl edit configmap apisix -n apisix
apisix:
  proxy_mode: http&stream
  stream_proxy:
    only: false
    tcp:
      - 9100
```

2. 暴露 SSH 端口。

```yaml
# kubectl edit svc apisix-gateway -n apisix
spec:
  ports:
    - name: http
      port: 80
      targetPort: 9080
    - name: ssh
      port: 22
      targetPort: 9100
```

3. 创建 Stream Route。

```yaml
apiVersion: apisix.apache.org/v2
kind: ApisixRoute
metadata:
  name: stream-sshpiper
  namespace: kubeflow
spec:
  ingressClassName: apisix
  stream:
    - name: tcp-access
      protocol: TCP
      match:
        ingressPort: 9100
      backend:
        serviceName: sshpiper
        servicePort: 22
```

再次强调：

- `ingressPort: 9100` 对应 APISIX Pod 内部 stream 监听端口
- `servicePort: 22` 对应后端 `sshpiper` Service 暴露的端口
- `2222` 是 sshpiper Pod 的容器端口，不是 APISIX route backend 里要填的值

### 3.7 为什么 GatewayProxy 里没有 9100

很多人第一次看到下面这个资源时会困惑：

```yaml
spec:
  provider:
    controlPlane:
      service:
        name: apisix-admin
        port: 9180
```

看起来只有 `9180`，为什么没有 `9100`？

原因是 `GatewayProxy` 只负责告诉 Ingress Controller / ADC：

- APISIX Admin API 在哪里
- 用什么方式认证
- 往哪里写入运行时配置

所以它属于“控制面”资源，只会出现：

- `apisix-admin`
- `9180`
- `adminKey`

而 `9100` 属于“数据面”资源，它只会出现在下面这些地方：

- APISIX `config.yaml` 的 `stream_proxy.tcp`
- `apisix-gateway` Service 的 `targetPort`
- `ApisixRoute.spec.stream.match.ingressPort`
- APISIX Admin API 的 `/stream_routes`

---

## 四、快速开始

### 4.1 安装 APISIX

```bash
# 推荐直接使用组件目录脚本
cd ..
./deploy_all.sh
```
建议升级到这组版本
按官方 apisix Helm chart 2.13.0 对应的版本，建议统一成这一套：
apache/apisix:3.15.0-ubuntu
apache/apisix-ingress-controller:2.0.1
ghcr.io/api7/adc:0.24.2

如果只想单独部署 APISIX，对应的核心 Helm 参数是：

```bash
helm repo add apisix https://apache.github.io/apisix-helm-chart --force-update
helm repo update
kubectl delete gatewayproxy apisix -n apisix --ignore-not-found

helm upgrade --install apisix apisix/apisix \
  -n apisix \
  --version 2.13.0 \
  --reuse-values \
  --set image.repository=apache/apisix \
  --set image.tag=3.15.0-ubuntu \
  --set ingress-controller.enabled=true \
  --set ingress-controller.image.repository=apache/apisix-ingress-controller \
  --set ingress-controller.image.tag=2.0.1 \
  --set ingress-controller.adcContainer.image.repository=ghcr.io/api7/adc \
  --set ingress-controller.adcContainer.image.tag=0.24.2 \
  --set ingress-controller.gatewayProxy.createDefault=true \
  --set ingress-controller.gatewayProxy.provider.controlPlane.service.name=apisix-admin \
  --set ingress-controller.gatewayProxy.provider.controlPlane.service.port=9180 \
  --set ingress-controller.gatewayProxy.provider.controlPlane.auth.adminKey.value=edd1c9f034335f136f87ad84b625c8f1 \
  --set ingress-controller.config.apisix.serviceNamespace=apisix \
  --set 'apisix.proxy_mode=http&stream' \
  --set 'apisix.stream_proxy.tcp[0]=9100'
```

脚本还会额外做这些兼容处理：

- 自动清理旧版默认 `GatewayProxy apisix`
- 检查并修复 `proxy_mode` 配置
- 给 `apisix-gateway` 自动补 SSH 端口 `22 -> 9100`

### 4.2 验证安装

```bash
# 检查 Pod
kubectl get pods -n apisix

# 检查 Service
kubectl get svc -n apisix
```

### 4.3 本地开发测试

```bash
# HTTP 访问
kubectl port-forward svc/apisix-gateway -n apisix 8888:80

# SSH 访问
kubectl port-forward svc/apisix-gateway -n apisix 2222:22
ssh caixukun-notebook-xxx@localhost -p 2222
```

---

## 五、故障排查

### 5.1 HTTP 404

症状：访问 Notebook 返回 `404 page not found`。

```bash
# 1. 检查 ApisixRoute 是否已创建
kubectl get apisixroute -n <namespace>

# 2. 检查路由是否已同步到 APISIX
kubectl port-forward svc/apisix-admin -n apisix 9180:9180
curl -s http://127.0.0.1:9180/apisix/admin/routes \
  -H "X-API-KEY: edd1c9f034335f136f87ad84b625c8f1" | jq '.list[].value.name'

# 3. 检查 Ingress Controller 日志
kubectl logs -n apisix deployment/apisix-ingress-controller -c manager --tail=50
```

常见错误：`no GatewayProxy configs provided`。  
原因：Ingress Controller 2.0+ 需要 `GatewayProxy` 资源来配置 Admin API 连接。

### 5.2 Ingress Controller 提示 no GatewayProxy configs provided

```bash
kubectl get gatewayproxy -n apisix
kubectl logs -n apisix deployment/apisix-ingress-controller -c manager --tail=50
```

如果集群里有旧版 `gatewayproxy/apisix` 与 Helm 新版 `gatewayproxy/apisix-config` 同时存在，建议先清理旧版默认资源：

```bash
kubectl delete gatewayproxy apisix -n apisix
```

### 5.3 路由未同步

如果 ApisixRoute 状态正常但 Admin API 查不到路由：

1. 检查 `GatewayProxy`：`kubectl get gatewayproxy -n apisix`
2. 检查 `apisix-ingress-config`
3. 检查 `proxy_mode` 是否包含 `stream`

```bash
kubectl exec deployment/apisix -n apisix -- cat /usr/local/apisix/conf/config.yaml | grep proxy_mode
```

如果显示 `proxy_mode: http`，说明 APISIX 还没启用 stream 模式。

对于 SSH stream，建议一定要把“CR 层”和“运行时层”分开检查：

```bash
# K8s CR 层
kubectl get apisixroute stream-sshpiper -n kubeflow -o yaml

# APISIX 运行时层
kubectl port-forward -n apisix svc/apisix-admin 9180:9180
curl http://127.0.0.1:9180/apisix/admin/stream_routes \
  -H "X-API-KEY: edd1c9f034335f136f87ad84b625c8f1"
```

如果前者有 `stream-sshpiper`，后者却是：

```json
{"total":0,"list":[]}
```

说明问题不在 `Pipe` 或 Notebook 后端，而在：

- `ApisixRoute.spec.stream`
- Ingress Controller / ADC
- APISIX Admin API

这条同步链路之间。

### 5.4 GatewayProxy 冲突

症状：Helm upgrade APISIX 时出现：

```text
gateway proxy configuration conflict
```

通常表示集群里同时存在两个 `GatewayProxy`：

- `apisix/apisix`
- `apisix/apisix-config`

排查命令：

```bash
kubectl get gatewayproxy -n apisix
kubectl get gatewayproxy apisix -n apisix -o yaml
kubectl get gatewayproxy apisix-config -n apisix -o yaml
```

如果 `apisix` 是旧版默认资源，可以直接删除：

```bash
kubectl delete gatewayproxy apisix -n apisix
```

### 5.5 APISIX Pod 启动即崩

症状：`apisix` Pod 反复 `CrashLoopBackOff`，日志中出现：

```text
did not find expected key
```

通常是 `config.yaml` 的 YAML 格式被改坏了，比如把下面两行误拼成一行：

```yaml
proxy_mode: "http&stream"
stream_proxy:
```

排查命令：

```bash
kubectl get configmap apisix -n apisix -o jsonpath='{.data.config\.yaml}' | nl -ba | sed -n '64,70p'
```

### 5.6 SSH 连接失败

SSH 建议按链路顺序排查：

```bash
# 1. 检查 APISIX Gateway Service 是否为 22 -> 9100
kubectl get svc apisix-gateway -n apisix -o yaml

# 2. 检查 APISIX 是否真的监听 9100
kubectl exec -n apisix deployment/apisix -- netstat -tlnp | grep 9100

# 3. 检查 GatewayProxy 是否正确指向 apisix-admin:9180
kubectl get gatewayproxy apisix-config -n apisix -o yaml

# 4. 检查 Stream Route CR 是否存在，且 ingressPort 是否为 9100
kubectl get apisixroute stream-sshpiper -n kubeflow -o yaml

# 5. 检查 APISIX 运行时里是否真的有 stream_routes
kubectl port-forward -n apisix svc/apisix-admin 9180:9180
curl http://127.0.0.1:9180/apisix/admin/stream_routes \
  -H "X-API-KEY: edd1c9f034335f136f87ad84b625c8f1"

# 6. 检查 SSHPiper Pod / Service
kubectl get pods -n kubeflow -l app=sshpiper
kubectl get svc sshpiper -n kubeflow -o yaml
kubectl get endpoints sshpiper -n kubeflow -o wide

# 7. 检查 Pipe 是否创建在 kubeflow
kubectl get pipe -n kubeflow
kubectl get pipe -n apisix
kubectl get pipe -n <user-namespace>

# 8. 检查具体 Notebook 的 SSH Service
kubectl get svc -n <namespace> <notebook-name>-ssh -o yaml
kubectl get endpoints -n <namespace> <notebook-name>-ssh -o wide

# 9. 查看 sshpiper 实时日志
kubectl logs -n kubeflow deployment/sshpiper --tail=100 -f

# 10. 如果前面都正常，再进入 Notebook Pod 检查 sshd
kubectl exec -n <namespace> <pod-name> -- sh -c "ss -lntp | grep ':22' || ps -ef | grep [s]shd"
```

看命令输出时，建议按下面方式判断：

- `kubectl get pipe -n apisix` 或 `kubectl get pipe -n zzz` 为空，通常是正常的，因为 Pipe 设计上就在 `kubeflow`
- `stream-sshpiper` CR 存在，但 `/stream_routes` 为空，说明 CR 没同步进 APISIX runtime
- `/stream_routes` 有数据，但 `sshpiper` 日志没有连接记录，优先看 APISIX 到 `sshpiper` 的连通性
- `sshpiper` 日志已经出现 `ssh connection pipe created ...`，说明 APISIX -> SSHPiper 这一段已经通了，问题应继续往 Pipe / Notebook SSH Service / sshd 查

---

## 六、本目录文件说明

| 文件 | 用途 | 使用方式 |
|------|------|----------|
| `readme.md` | 本文档 | 阅读 |
| `apisix-ingress-config-fix.yaml` | 修复 IngressClass 配置 | `kubectl apply -f` |
| `apisix-stream-proxy.yaml` | Stream Proxy 配置参考 | 仅参考，不要直接 `kubectl apply` |
| `apisix-config-template.yaml` | APISIX 完整配置模板 | 仅参考，不要直接 `kubectl apply` |

注意：ConfigMap 类文件会完全覆盖现有配置，请通过 `kubectl edit` 或 `kubectl patch` 做增量修改。

---

## 参考资料

- [APISIX Ingress Controller 文档](https://apisix.apache.org/docs/ingress-controller/)
- [APISIX 官方文档](https://apisix.apache.org/docs/apisix/)
- [forward-auth 插件](https://apisix.apache.org/docs/apisix/plugins/forward-auth/)
- [Stream Proxy 配置](https://apisix.apache.org/docs/apisix/stream-proxy/)

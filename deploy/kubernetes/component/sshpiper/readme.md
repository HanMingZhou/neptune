# SSHPiper 部署与 SSH 访问架构

## 目录

- [一、SSHPiper 是什么](#一sshpiper-是什么)
- [二、SSH 连接全链路](#二ssh-连接全链路)
- [三、Pipe 路由规则](#三pipe-路由规则)
- [四、认证模式](#四认证模式)
- [五、创建 Notebook 时的 SSH 资源](#五创建-notebook-时的-ssh-资源)
- [六、部署步骤](#六部署步骤)
- [七、故障排查](#七故障排查)

---

## 一、SSHPiper 是什么

SSHPiper 是一个 SSH 反向代理，功能类似 Nginx，但专门处理 SSH 协议。它根据 SSH 登录用户名，把连接路由到不同的后端 Pod。

在本系统里，用户不能直接从集群外访问某个 Notebook Pod，因此需要两段能力：

1. 让 SSH 流量先进入集群
2. 进入集群后，再根据用户名找到正确的 Notebook

这两段能力分别由不同组件负责：

- APISIX 负责把 SSH 流量引进集群
- SSHPiper 负责按用户名把流量分发到正确 Notebook

所以可以把 SSHPiper 理解成 SSH 方向的“按用户名路由器”。

---

## 二、SSH 连接全链路

用户执行 `ssh hmz-notebook-xxx@ai.local -p 22` 后，请求实际上会经过多层端口映射。把这些层拆开看，就不会再被 `22 / 9100 / 22 / 2222 / 22` 绕晕了。

### 2.1 端口速查表

| 层级 | 资源 | 端口 | 说明 |
|------|------|------|------|
| 用户访问 | 外部地址 | `22` | 用户看到的 SSH 端口 |
| APISIX Gateway Service | `apisix-gateway` | `22 -> 9100` | 对外暴露 `22`，转发到 APISIX Pod 内部 stream proxy `9100` |
| APISIX Pod | stream proxy | `9100` | APISIX 内部监听 TCP SSH 流量的端口，`ApisixRoute.spec.stream.match.ingressPort` 必须填它 |
| SSHPiper Service | `sshpiper.kubeflow.svc` | `22 -> 2222` | APISIX backend 连的是 Service 端口 `22`，再由 Service 转到 Pod 的 `2222` |
| Notebook SSH Service | `notebook-xxx-ssh.<ns>.svc` | `22 -> 22` | SSHPiper 找到目标 Notebook 后，转给这个稳定的 Service 地址 |
| Notebook Pod | `sshd` | `22` | 容器内部最终处理登录的 sshd |

两个最容易混淆的字段要这样理解：

- `ApisixRoute.spec.stream.match.ingressPort = 9100`
  这是 APISIX 内部监听端口
- `ApisixRoute.spec.stream.backend.servicePort = 22`
  这是后端 `sshpiper` Service 暴露的端口

`2222` 只是 sshpiper Pod 的容器端口，只属于 `sshpiper` Service 的 `targetPort`，不应该直接填到 APISIX route 的 backend 里。

### 2.2 全链路图

```text
用户电脑                          K8s 集群
--------                          --------

ssh hmz-notebook-xxx@ai.local:22
        │
        ▼
┌──────────────────────────────┐
│  ① apisix-gateway Service    │  22 -> 9100
└──────────────┬───────────────┘
               │
               ▼
┌──────────────────────────────┐
│  ② APISIX stream proxy       │  listen 9100
│     (Pod 内部)                │  Stream Route: 9100 -> sshpiper:22
└──────────────┬───────────────┘
               │
               ▼
┌──────────────────────────────┐
│  ③ sshpiper Service          │  22 -> 2222
│     (kubeflow ns)            │
└──────────────┬───────────────┘
               │
               ▼
┌──────────────────────────────┐
│  ④ SSHPiper Pod              │  listen 2222
│                               │  Pipe: 按用户名匹配后端
└──────────────┬───────────────┘
               │
               ▼
┌──────────────────────────────┐
│  ⑤ Notebook SSH Service      │  22 -> 22
│     (用户 namespace)          │
└──────────────┬───────────────┘
               │
               ▼
┌──────────────────────────────┐
│  ⑥ Notebook Pod              │  sshd listen 22
└──────────────────────────────┘
```

### 2.3 每层分别解决什么问题

| 层 | 解决的问题 |
|---|---|
| ① `apisix-gateway` Service | 把集群外用户访问的 `22` 转给 APISIX 内部 `9100` |
| ② APISIX Stream Route | 在 APISIX 内部 `9100` 上匹配路由，并把流量转给 `sshpiper:22` |
| ③ `sshpiper` Service | 把 `sshpiper:22` 转给 sshpiper Pod 的 `2222` |
| ④ SSHPiper Pipe | 按 SSH 用户名把连接路由到正确 Notebook 的 `-ssh` Service |
| ⑤ Notebook SSH Service | 给 SSHPiper 提供稳定的后端 DNS 地址 |
| ⑥ Notebook Pod 内 sshd | 真正做最终登录认证和会话处理 |

Kubeflow Notebook 控制器默认只创建 Jupyter 用的 Service，不会暴露 22，所以系统需要额外创建一个 `*-ssh` Service。

### 2.4 基于实际资源的链路样例

下面这组命令和输出来自一次真实排障，适合在你忘记资源落位时直接照着查：

```bash
kubectl get svc -n kubeflow
```

```text
NAME         TYPE        CLUSTER-IP    PORT(S)
sshpiper     ClusterIP   10.233.6.39   22/TCP
```

```bash
kubectl get pipes.sshpiper.com -n kubeflow
```

```text
NAME                       FROMUSER              TOUSER   TOHOST
pipe-zzz-notebook-eeb47f   zzz-notebook-eeb47f   root     notebook-eeb47f-ssh.zzz.svc.cluster.local:22
```

```bash
kubectl get pipes.sshpiper.com -n apisix
kubectl get pipes.sshpiper.com -n zzz
```

```text
No resources found
```

这两条空结果通常也是正常的，因为 Pipe 的设计位置就是 `kubeflow`，不是 `apisix`，也不是具体业务 namespace。

再把这几条资源和一次真实访问串起来看：

```text
Windows:
  ssh -p 30177 zzz-notebook-eeb47f@10.255.141.8
        │
        ▼
APISIX 把流量转给 kubeflow/sshpiper:22
        │
        ▼
sshpiper Service 10.233.6.39:22
        │
        ▼
sshpiper Pod :2222
        │
        ▼
匹配 Pipe:
  from.user = zzz-notebook-eeb47f
  to.user   = root
  to.host   = notebook-eeb47f-ssh.zzz.svc.cluster.local:22
        │
        ▼
zzz/notebook-eeb47f-ssh:22
        │
        ▼
zzz/notebook-eeb47f-0:sshd
```

### 2.5 全局资源和每个 Notebook 资源

这是最容易混掉的一层，建议直接记成下面这张表：

| 资源 | 是否全局唯一 | Namespace | 作用 |
|------|--------------|-----------|------|
| `sshpiper` Deployment | 是 | `kubeflow` | 统一接收所有进入集群的 SSH 连接 |
| `sshpiper` Service | 是 | `kubeflow` | 给 APISIX 一个固定的后端入口 |
| `stream-sshpiper` ApisixRoute | 是 | `kubeflow` | 把 APISIX `9100` 的 TCP 流量转给 `sshpiper:22` |
| `Pipe` | 否，每个 Notebook 一条 | `kubeflow` | 按用户名把连接路由到正确的 Notebook SSH Service |
| `notebook-xxx-ssh` Service | 否，每个 Notebook 一条 | 用户 namespace | 给 SSHPiper 提供稳定的后端地址 |

这里最关键的理解是：

- `sshpiper` Service 不是给某个 Notebook 单独准备的
- 它是整个集群的统一 SSH 汇聚入口
- 每个 Notebook 真正不同的部分是 `Pipe` 和 `notebook-xxx-ssh` Service

---

## 三、Pipe 路由规则

Pipe 是 SSHPiper 的 CRD，每个需要 SSH 访问的 Notebook 对应一个 Pipe。

注意：Pipe 创建在 `kubeflow` namespace，而不是业务 namespace。  
所以排查时应该优先用：

```bash
kubectl get pipe -A
kubectl get pipe -n kubeflow
```

### 3.1 From 和 To

一个 Pipe 有两部分：

- `From`：定义“什么条件的连接匹配这个 Pipe”
- `To`：定义“匹配后转发到哪里”

```text
From:                              To:
  用户名: hmz-notebook-xxx    ->     目标: notebook-xxx-ssh.hmz.svc:22
  验证方式: 用户公钥 / 无             连接方式: 平台私钥 / 密码透传
```

### 3.2 SSHPiper 的匹配流程

1. 用户执行 `ssh hmz-notebook-xxx@ai.local`
2. SSHPiper 收到连接，提取用户名 `hmz-notebook-xxx`
3. 在 `kubeflow` namespace 下遍历 Pipe
4. 找到 `From.Username == "hmz-notebook-xxx"` 的那条规则
5. 根据 `From` 配置决定如何认证
6. 根据 `To` 配置把连接转给后端 `notebook-xxx-ssh.<ns>.svc:22`

---

## 四、认证模式

系统支持两种 SSH 认证方式，由 Pipe 的配置决定行为。

### 4.1 公钥认证：双密钥系统

公钥模式下，SSH 连接被拆成两段，每段使用不同的密钥对：

```text
用户电脑                  SSHPiper                    Notebook Pod
--------                  --------                    -------------

持有: 用户私钥             持有: 用户公钥               持有: 平台公钥
                                平台私钥

第一段: 用户 -> SSHPiper
  用户用自己的私钥签名 --> SSHPiper 用用户公钥验证

第二段: SSHPiper -> Pod
  SSHPiper 用平台私钥签名 --> Pod 用平台公钥验证
```

为什么需要双密钥：

- SSHPiper 是中间代理，它把一条 SSH 连接拆成了两段
- SSHPiper 不应该持有用户私钥
- 所以平台需要再生成一对“平台密钥”，专门用于 SSHPiper -> Pod 这一段

| 密钥 | 来源 | 私钥位置 | 公钥位置 | 用途 |
|------|------|----------|----------|------|
| 用户密钥 | 用户上传 | 用户电脑 | Pipe 的 `From` + Pod 的 `authorized_keys` | 验证用户身份 |
| 平台密钥 | 创建 Notebook 时自动生成 | `kubeflow` 的 Secret | Pod 的 `authorized_keys` | 让 SSHPiper 连接 Pod |

### 4.2 密码认证：透传模式

密码模式下，SSHPiper 不自己验证密码，而是直接把用户输入的密码透传给后端 Pod 的 sshd：

```text
用户输入密码
      │
      ▼
SSHPiper（不验证，直接透传）
      │
      ▼
Pod 的 sshd 验证密码 -> 通过 / 拒绝
```

实现方式是“什么都不配”：

- `From` 不设置 `AuthorizedKeysData` / `HtpasswdData`
- `To` 不设置 `PrivateKeySecret` / `PasswordSecret`

这样 SSHPiper 会自动走透传逻辑。

好处是：用户在容器里执行 `passwd` 修改密码后立即生效，不需要再同步外部配置。

### 4.3 Pipe 字段怎么理解

| From 字段 | 设置了 | 没设置 |
|-----------|--------|--------|
| `AuthorizedKeysData` | SSHPiper 用此公钥验证用户 | SSHPiper 不自己验证，放行 |
| `HtpasswdData` | SSHPiper 用 htpasswd 验证密码 | 不做密码验证 |

| To 字段 | 设置了 | 没设置 |
|---------|--------|--------|
| `PrivateKeySecret` | SSHPiper 用此私钥连后端 | 透传用户凭证给后端 |
| `PasswordSecret` | SSHPiper 用固定密码连后端 | 透传用户凭证给后端 |

---

## 五、创建 Notebook 时的 SSH 资源

创建一个需要 SSH 访问的 Notebook 时，系统会自动创建以下资源：

### 5.1 资源清单

| 资源 | Namespace | 作用 | 删除/停止时 |
|------|-----------|------|------------|
| SSH Service | 用户 namespace | 暴露 Pod 的 22 端口，供 SSHPiper 连接 | 删除 |
| Pipe | `kubeflow` | SSHPiper 路由规则 | 删除 |
| 公钥 Secret | 用户 namespace | 挂载到 Pod 的 `authorized_keys` | 删除时删除，停止时保留 |
| 私钥 Secret | `kubeflow` | SSHPiper 连 Pod 用的平台私钥，被 Pipe 引用 | 删除时删除，停止时保留 |
| 密码 Secret | 用户 namespace | 挂载到 Pod，启动时设置 root 密码 | 删除时删除，停止时保留 |
| Stream Route | `kubeflow` | APISIX 内部 `9100` -> `sshpiper:22`，全局共享 | 保留 |

注意最后一行：`Stream Route` 虽然是在创建 Notebook 时触发代码去 ensure，但它的语义是“全局共享资源”，不是每个 Notebook 单独一条。

### 5.2 代码创建了什么，没创建什么

结合本仓库代码，创建 Notebook SSH 能力时，逻辑上会做三类事情：

1. 在用户 namespace 创建 `notebook-xxx-ssh` Service
2. 在 `kubeflow` 创建 `Pipe`
3. 在 `kubeflow` ensure 一条全局 `stream-sshpiper` ApisixRoute

但要注意：

- 代码创建的是 Kubernetes 里的 `ApisixRoute CR`
- 真正让 APISIX 生效，还需要 Ingress Controller / ADC 把它同步到 APISIX Admin API 的 `/stream_routes`

所以排障时一定要同时看：

```bash
# K8s CR 层
kubectl get apisixroute stream-sshpiper -n kubeflow -o yaml

# APISIX 运行时层
kubectl port-forward -n apisix svc/apisix-admin 9180:9180
curl http://127.0.0.1:9180/apisix/admin/stream_routes \
  -H "X-API-KEY: edd1c9f034335f136f87ad84b625c8f1"
```

### 5.3 Secret 的具体作用

```text
公钥 Secret（用户 ns）
  内容: 用户公钥 + 平台公钥
  挂载: Pod 的 /root/.ssh/authorized_keys
  使用者: Pod 里的 sshd

私钥 Secret（kubeflow ns）
  内容: 平台私钥
  引用方: Pipe 的 To.PrivateKeySecret
  使用者: SSHPiper

密码 Secret（用户 ns，可选）
  内容: 随机生成的密码
  挂载: Pod 启动时 postStart 读取并设置 root 密码
  说明: SSHPiper 本身不验证密码，只做透传
```

### 5.4 公钥传递流程

```text
① 用户上传 SSH 公钥到平台
           │
           ▼
② 平台生成一对密钥（平台公钥 + 平台私钥）
           │
     ┌─────┴─────┐
     ▼           ▼
③ 创建公钥      ④ 创建私钥
   Secret          Secret
   (用户 ns)       (kubeflow ns)
     │               │
     ▼               ▼
⑤ Pod 挂载公钥    ⑥ Pipe 引用私钥
     │               │
     └──────┬────────┘
            ▼
⑦ 用户连接 SSH，SSHPiper -> Pod 完成双段认证
```

---

## 六、部署步骤

### 6.1 生成 Host Key

SSHPiper 作为 SSH 服务器需要一个 Host Key：

```bash
ssh-keygen -f ssh_host_ed25519_key -t ed25519 -N ''
```

### 6.2 创建 SSHPiper 服务端私钥 Secret

`deploy.yml` 只引用 `kubeflow/sshpiper-server-key`，不会自动创建这个 Secret。生成好 Host Key 后执行：

```bash
kubectl create secret generic sshpiper-server-key \
  -n kubeflow \
  --from-file=server_key=ssh_host_ed25519_key
```

如果需要轮换：

```bash
kubectl delete secret sshpiper-server-key -n kubeflow
kubectl create secret generic sshpiper-server-key \
  -n kubeflow \
  --from-file=server_key=ssh_host_ed25519_key
```

### 6.3 安装 SSHPiper CRD

```bash
kubectl apply -f https://raw.githubusercontent.com/tg123/sshpiper/7ce7b52e6a71f167ee78fd439a19d016e610d1d2/plugin/kubernetes/crd.yaml
```

### 6.4 部署 SSHPiper

```bash
kubectl apply -f deploy.yml
```

必须注意的端口映射坑：

- `sshpiperd` 官方镜像默认监听容器内 `2222`
- 不是 `22`

所以 `deploy.yml` 里必须是：

```yaml
service:
  port: 22
  targetPort: 2222
```

如果这里写错，APISIX 发往 `sshpiper:22` 的流量就会被拒绝。

`2222` 的定义点就在本目录的部署清单里：

- `deploy.yml` 中 `Service.spec.ports[].targetPort: 2222`
- `deploy.yml` 中 `Deployment.spec.template.spec.containers[].ports[].containerPort: 2222`

这也是为什么代码里 APISIX route backend 应该写 `servicePort: 22`，而不是 `2222`：

- `22` 是 `sshpiper` Service 暴露给集群内其他组件访问的端口
- `2222` 只是 sshpiper Pod 内部真实监听端口

如果平台配置里还看到 `sshpiper.port: 2222`，要再额外区分一下：那通常只是服务端给前端拼接 SSH 命令时使用的配置项，不等于当前网关对外真正暴露给用户的端口。

### 6.5 验证

```bash
# 检查 Pod
kubectl get pods -n kubeflow -l app=sshpiper

# 检查 Service
kubectl get svc sshpiper -n kubeflow -o yaml

# 检查 Pipe（创建 Notebook 后才会有）
kubectl get pipe -A

# 检查 Stream Route
kubectl get apisixroute stream-sshpiper -n kubeflow -o yaml

# 检查 sshpiper Service 的 endpoints
kubectl get endpoints sshpiper -n kubeflow -o wide
```

---

## 七、故障排查

### 7.1 SSH 连接超时

```bash
# 检查 APISIX Gateway Service 是否为 22 -> 9100
kubectl get svc apisix-gateway -n apisix -o yaml

# 检查 APISIX 是否监听 9100
kubectl exec -n apisix deployment/apisix -- netstat -tlnp | grep 9100

# 检查 Stream Route 是否为 ingressPort: 9100 -> sshpiper:22
kubectl get apisixroute stream-sshpiper -n kubeflow -o yaml

# 检查 APISIX 运行时是否真的有 stream route
kubectl port-forward -n apisix svc/apisix-admin 9180:9180
curl http://127.0.0.1:9180/apisix/admin/stream_routes \
  -H "X-API-KEY: edd1c9f034335f136f87ad84b625c8f1"

# 检查 sshpiper Service 是否指向 Pod:2222
kubectl get svc sshpiper -n kubeflow -o yaml
kubectl get endpoints sshpiper -n kubeflow -o wide
```

如果 `stream-sshpiper` 这个 CR 已经存在，但 `/stream_routes` 仍然为空，说明问题不在 `Pipe`，而在 `ApisixRoute.stream -> ingress-controller/adc -> APISIX runtime` 这条同步链路。

### 7.2 SSH 认证失败

```bash
# 检查 Pipe 是否创建在 kubeflow
kubectl get pipe -A

# 查看 Pipe 详情
kubectl get pipe <pipe-name> -n kubeflow -o yaml

# 检查公钥 Secret
kubectl get secret -n <namespace> | grep ssh

# 检查平台私钥 Secret
kubectl get secret -n kubeflow | grep ssh
```

如果你看到：

```bash
kubectl get pipes.sshpiper.com -n apisix
kubectl get pipes.sshpiper.com -n <user-namespace>
```

返回为空，不必先怀疑 Pipe 丢了。先去 `kubeflow` namespace 看，因为代码和部署约定都是把 Pipe 建在那里。

### 7.3 SSHPiper 能连上但 Pod 拒绝

```bash
# 检查 SSH Service 是否存在
kubectl get svc -n <namespace> | grep ssh

# 检查 SSH Service 是否真的有 endpoints
kubectl get endpoints -n <namespace> <notebook-name>-ssh -o yaml

# 检查 Pod 内 sshd 是否在运行
kubectl exec -n <namespace> <pod-name> -- sh -c "ss -lntp | grep ':22' || ps aux | grep [s]shd"

# 检查 authorized_keys 是否正确挂载
kubectl exec -n <namespace> <pod-name> -- cat /root/.ssh/authorized_keys

# 检查 postStart 日志
kubectl exec -n <namespace> <pod-name> -- sh -c "cat /tmp/neptune-poststart.log 2>/dev/null; cat /tmp/neptune-ssh-init.log 2>/dev/null"
```

### 7.4 SSHPiper 日志

```bash
kubectl logs -n kubeflow deployment/sshpiper --tail=50
```

当 SSH 请求真正进到 SSHPiper 时，通常能看到类似这样的日志：

```text
ssh connection pipe created 10.233.x.x:xxxxx (username [zzz-notebook-eeb47f]) -> 10.233.x.x:22 (username [root])
```

如果已经看到这类日志，说明下面这几跳已经没问题了：

- APISIX -> `sshpiper` Service
- `sshpiper` Service -> sshpiper Pod
- Pipe 匹配到了正确用户名

这时应继续往后查：

- `notebook-xxx-ssh` Service 是否有 endpoints
- Notebook Pod 里的 `sshd` 是否正常监听 `22`
- 公钥 / 密码是否正确

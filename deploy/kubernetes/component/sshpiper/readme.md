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

SSHPiper 是一个 **SSH 反向代理**，功能类似 Nginx 但专门处理 SSH 协议。它根据 SSH 登录用户名将连接路由到不同的后端 Pod。

在本系统中，一个集群可能有几百个 Notebook 容器，用户无法直接从外部连到集群内部的 Pod。SSHPiper 解决的问题是：**通过一个统一入口（APISIX:22），根据用户名把 SSH 连接分发到正确的 Notebook Pod**。

---

## 二、SSH 连接全链路

用户执行 `ssh hmz-notebook-xxx@ai.local -p 22` 后，请求经过三层转发：

```
用户电脑                          K8s 集群
────────                          ────────

ssh hmz-nb-xxx@ai.local:22
        │
        ▼
┌─────────────────────┐
│  ① APISIX Gateway   │  Stream Route: 端口 22 的 TCP 流量 → SSHPiper
│     (端口 22)        │  （全局唯一，所有 Notebook 共享）
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│  ② SSHPiper         │  Pipe: 按用户名 "hmz-nb-xxx" 查找路由规则
│  (kubeflow ns)       │  → 转发到 notebook-xxx-ssh.hmz.svc:22
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│  ③ SSH Service      │  Service: 将流量导向 Notebook Pod 的 22 端口
│  (用户 namespace)    │  selector: statefulset=notebook-xxx
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│  ④ Notebook Pod     │  容器里的 sshd 处理登录
│     (sshd:22)        │
└─────────────────────┘
```

### 为什么需要三层？

| 层 | 解决的问题 |
|---|---|
| ① APISIX Stream Route | 把集群外的 SSH 流量引进来（TCP 9100 入口转发给 SSHPiper Service 的 22 端口） |
| ② SSHPiper Pipe | 按用户名路由到正确的 Pod（监听内部 2222 端口，收到流量后动态分发） |
| ③ SSH Service | 在 K8s 内部提供通往 Pod 22 端口的稳定 DNS 地址 |

> Kubeflow Notebook 控制器自动创建的 Service 只暴露 8888 端口（给 Jupyter 用），不暴露 22，所以需要额外创建一个 SSH Service。

---

## 三、Pipe 路由规则

Pipe 是 SSHPiper 的 CRD（Custom Resource），每个需要 SSH 访问的 Notebook 对应一个 Pipe。

### From 和 To

一个 Pipe 有两部分：

- **From（入口）**：定义"什么条件的连接匹配这个 Pipe"
- **To（出口）**：定义"匹配后转发到哪里"

```
From:                              To:
  用户名: hmz-notebook-xxx    →      目标: notebook-xxx-ssh.hmz.svc:22
  验证方式: 用户公钥 / 无             连接方式: 平台私钥 / 密码透传
```

### SSHPiper 的匹配流程

1. 用户执行 `ssh hmz-notebook-xxx@ai.local`
2. SSHPiper 收到连接，提取用户名 `hmz-notebook-xxx`
3. 遍历 kubeflow namespace 下所有 Pipe，找到 `From.Username == "hmz-notebook-xxx"` 的那个
4. 根据 From 的配置验证用户身份
5. 根据 To 的配置连接后端 Pod

---

## 四、认证模式

系统支持两种 SSH 认证方式，由 Pipe 的配置决定行为。

### 4.1 公钥认证（双密钥系统）

公钥模式下，SSH 连接被拆成两段，每段使用不同的密钥对：

```
用户电脑                  SSHPiper                    Notebook Pod
────────                  ────────                    ──────────────

持有: 用户私钥             持有: 用户公钥               持有: 平台公钥
                                平台私钥

第一段: 用户 → SSHPiper
  用户用自己的私钥签名 ──→ SSHPiper 用用户公钥验证 ✅

第二段: SSHPiper → Pod
  SSHPiper 用平台私钥签名 ──→ Pod 用平台公钥验证 ✅
```

**为什么需要"双密钥"？**

SSHPiper 是中间代理，它把一条 SSH 连接拆成了两段。SSHPiper 没有用户的私钥（也不应该有），所以它没法拿用户的私钥去连后端。因此需要**平台自己生成一对密钥**，专门用于 SSHPiper→Pod 这一段的认证。

| 密钥 | 来源 | 私钥位置 | 公钥位置 | 用途 |
|------|------|----------|----------|------|
| **用户密钥** | 用户生成并上传 | 用户电脑 | Pipe 的 From 配置 + Pod 的 authorized_keys | 验证用户身份 |
| **平台密钥** | 创建 Notebook 时自动生成 | kubeflow 的 Secret → Pipe 的 To 引用 | Pod 的 authorized_keys | SSHPiper 连接 Pod |

### 4.2 密码认证（透传模式）

密码模式下，SSHPiper **不做任何验证**，直接把用户输入的密码原样转发给后端 Pod 的 sshd。

```
用户输入密码 123456
        │
        ▼
SSHPiper（不验证，直接透传）
        │
        ▼
Pod 的 sshd 验证密码 → 通过/拒绝
```

**密码透传的实现方式**：这是 SSHPiper 的**默认行为** —— 当 Pipe 的 From 中不设置 `AuthorizedKeysData` 或 `HtpasswdData`，To 中不设置 `PrivateKeySecret` 或 `PasswordSecret` 时，SSHPiper 自动透传凭证。**不做配置本身就是实现**。

这种模式的好处：用户在容器内执行 `passwd` 修改密码后立即生效，不需要更新任何外部配置。

### 4.3 Pipe 配置对照表

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

创建一个需要 SSH 访问的 Notebook 时，系统自动创建以下资源：

### 5.1 资源清单

| 资源 | Namespace | 作用 | 删除/停止时 |
|------|-----------|------|------------|
| **SSH Service** | 用户 ns | 暴露 Pod 的 22 端口，供 SSHPiper 连接 | 删除 |
| **Pipe** | kubeflow | SSHPiper 路由规则 | 删除 |
| **公钥 Secret** | 用户 ns | 挂载到 Pod 的 authorized_keys（含用户公钥 + 平台公钥） | 删除时删除，停止时保留 |
| **私钥 Secret** | kubeflow | SSHPiper 连 Pod 用的平台私钥，被 Pipe 引用 | 删除时删除，停止时保留 |
| **密码 Secret** | 用户 ns | 挂载到 Pod，启动时设置 root 密码（仅密码模式） | 删除时删除，停止时保留 |
| **Stream Route** | kubeflow | APISIX:22 → SSHPiper:22（全局共享，仅首次创建） | 保留 |

### 5.2 Secret 的具体作用

```
┌──────────────────────────────────────────────────────────┐
│                    公钥 Secret（用户 ns）                   │
│  内容: 用户上传的公钥 + 平台生成的公钥                       │
│  挂载: Pod 的 /root/.ssh/authorized_keys                  │
│  谁用: Pod 里的 sshd，用来验证 SSH 连接                    │
└──────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────┐
│                    私钥 Secret（kubeflow ns）               │
│  内容: 平台生成的私钥                                       │
│  引用方: Pipe 的 To.PrivateKeySecret 字段                  │
│  谁用: SSHPiper，用此私钥连接后端 Pod                       │
│  为什么在 kubeflow: K8s Secret 只能被同 ns 的资源引用       │
└──────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────┐
│                    密码 Secret（用户 ns，可选）              │
│  内容: 自动生成的随机密码                                   │
│  挂载: Pod 启动时 postStart 读取并设置 root 密码            │
│  与 SSHPiper 无关: SSHPiper 密码模式只做透传                │
└──────────────────────────────────────────────────────────┘
```

### 5.3 公钥传递到容器的流程

```
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
   内容:            内容:
   用户公钥         平台私钥
   + 平台公钥
     │
     ▼
⑤ Pod 启 动，Volume 挂载 Secret
     │
     ▼
⑥ postStart 脚本：
   复制公钥 → ~/.ssh/authorized_keys
   启动/重启 sshd
     │
     ▼
⑦ 用户可以 SSH 连接
```

---

## 六、部署步骤

### 6.1 生成 Host Key

SSHPiper 作为 SSH 服务器需要一个 Host Key（类似 HTTPS 的服务器证书）：

```bash
# 生成 Host Key
ssh-keygen -f ssh_host_ed25519_key -t ed25519 -N ''
```

### 6.2 创建 SSHPiper 服务端私钥 Secret

`deploy.yml` 只引用 `kubeflow/sshpiper-server-key`，不会在仓库中创建私钥 Secret。生成好 Host Key 后执行：

```bash
kubectl create secret generic sshpiper-server-key \
  -n kubeflow \
  --from-file=server_key=ssh_host_ed25519_key
```

如果你已经在别处生成了服务端私钥，也可以把上面的 `ssh_host_ed25519_key` 替换成你的实际文件路径。

如果需要轮换，可以先删后建：

```bash
kubectl delete secret sshpiper-server-key -n kubeflow
kubectl create secret generic sshpiper-server-key \
  --from-file=server_key=ssh_host_ed25519_key \
  -n kubeflow
```

建议使用专门给当前集群生成的服务端私钥，不要继续复用已经提交过仓库或用于其他环境的旧密钥。

如果你使用 [deploy_all.sh](/Users/jerrytom/go/src/test/kubeflow/neptune/deploy/kubernetes/component/deploy_all.sh)，脚本会先检查 `kubeflow/sshpiper-server-key` 是否已存在；若不存在且本地有 `./sshpiper/ssh_host_ed25519_key`，会自动创建该 Secret。生产环境建议改用你自己生成的私钥，并在执行时指定：

```bash
SSHPIPER_SERVER_KEY_FILE=/path/to/ssh_host_ed25519_key ./deploy_all.sh
```

### 6.3 安装 SSHPiper CRD

```bash
kubectl apply -f https://raw.githubusercontent.com/tg123/sshpiper/7ce7b52e6a71f167ee78fd439a19d016e610d1d2/plugin/kubernetes/crd.yaml
```

### 6.4 部署 SSHPiper

```bash
kubectl apply -f deploy.yml
```

> **⚠️ 必须注意的端口映射坑**：
> `sshpiperd` 官方镜像默认监听容器内的 **`2222`** 端口，而不是 `22` 端口！
> 因此，在 `deploy.yml` 的 Service 配置中，必须写明 `port: 22`映射到 `targetPort: 2222`，同时 Deployment 的 `containerPort` 必须声明为 `2222`，否则 APISIX 发往 `sshpiper:22` 的流量将遭到拒绝 (`Connection refused`)。

### 6.5 验证

```bash
# 检查 Pod
kubectl get pods -n kubeflow -l app=sshpiper

# 检查 Service
kubectl get svc -n kubeflow sshpiper

# 检查 Pipe（创建 Notebook 后才会有）
kubectl get pipe -n kubeflow
```

---

## 七、故障排查

### 7.1 SSH 连接超时

```bash
# 检查 APISIX Stream Route 是否存在
kubectl get apisixroute -n kubeflow | grep stream

# 检查 APISIX 是否监听 9100 端口
kubectl exec -n apisix deployment/apisix -- netstat -tlnp | grep 9100
```

### 7.2 SSH 认证失败

```bash
# 检查 Pipe 是否创建
kubectl get pipe -n kubeflow

# 查看 Pipe 详情（确认 From 和 To 配置是否正确）
kubectl get pipe <pipe-name> -n kubeflow -o yaml

# 检查公钥 Secret 是否存在
kubectl get secret -n <namespace> | grep ssh

# 检查私钥 Secret 是否存在（在 kubeflow namespace）
kubectl get secret -n kubeflow | grep ssh
```

### 7.3 SSHPiper 能连上但 Pod 拒绝

```bash
# 检查 SSH Service 是否存在
kubectl get svc -n <namespace> | grep ssh

# 检查 Pod 内 sshd 是否在运行
kubectl exec -n <namespace> <pod-name> -- ps aux | grep sshd

# 检查 authorized_keys 是否正确挂载
kubectl exec -n <namespace> <pod-name> -- cat /root/.ssh/authorized_keys
```

### 7.4 SSHPiper 日志

```bash
kubectl logs -n kubeflow deployment/sshpiper --tail=50
```

# SSHPiper 部署与 SSH 访问架构

本篇文档详细说明了 Neptune 机器学习平台中基于 **APISIX + SSHPiper** 架构实现的容器实例（Notebook）动态 SSH 访问链路、双段密钥验证机制，以及跨命名空间路由机制。

---

## 📖 目录

- [一、整体访问架构](#一整体访问架构)
- [二、端口对应与流转链路](#二端口对应与流转链路)
- [三、双段认证机制（核心设计）](#三双段认证机制核心设计)
- [四、跨命名空间 (Cross-Namespace) 路由原理](#四跨命名空间-cross-namespace-路由原理)
- [五、创建实例时的 SSH 资源落位](#五创建实例时的-ssh-资源落位)
- [六、部署步骤](#六部署步骤)
- [七、故障排查与诊断指南](#七故障排查与诊断指南)

---

## 一、整体访问架构

在 Neptune 系统中，用户无法直接从集群外通过 SSH 连接至特定的 Notebook 容器。这涉及到两个关键的限制与设计考量：

1. **网络层限制（内网隔离与端口保护）**：
   - K8s Pod 拥有的是私有内网 IP（如 `10.233.x.x`），集群外部根本无法直接路由寻址。
   - 如果为每个实例单独暴露 NodePort 端口（如 `30001`、`30002`...）或分配独立的公网负载均衡 IP，当用户量庞大时，会消耗海量的公网端口与 IP 资源，不仅管理极其复杂，且将海量端口直接暴露在公网上会带来巨大安全漏洞。
   - **解决方案**：引入统一网关。系统在集群边界仅暴露 **1个** 统一的公网 NodePort 端口（如 `30177`），所有的 SSH 流量都通过此单一入口进入，并通过 **SSH 登录用户名**（如 `namespace-instancename`）来进行动态路由和分发。

2. **协议层限制（SSH 端到端加密）**：
   - SSH 协议是高强度、端到端加密的网络协议。普通的网络层网关（如 APISIX）在传输层只能看到加密后的 TCP 数据流，无法在密文中直接读取握手包中的登录用户名。
   - **解决方案**：引入 **SSHPiper**。它作为一个 SSH 反向代理，能“中断”前半段客户端连接，读取并解析出用户名，根据匹配的路由规则重新发起对内网 Pod 的后半段连接，并将两段会话进行内存流桥接。

### 1. 流量链路流转图

```text
用户电脑 (客户端)                 K8s 集群边界                  K8s 集群内部
----------------                 ------------                  ------------

ssh user-xxx@node-ip -p 30177
       │
       ▼ (统一入口 NodePort: 30177)
┌──────────────────────────────┐
│     apisix-gateway Service   │  对外暴露端口，转发流量至 APISIX 容器端口 9100 (22 -> 9100)
└──────────────┬───────────────┘
               │
               ▼
┌──────────────────────────────┐
│     APISIX Stream Proxy      │  匹配 StreamRoute (L4 路由规则)，原样透明转发至 sshpiper Service:22
└──────────────┬───────────────┘
               │
               ▼
┌──────────────────────────────┐
│     sshpiper Service         │  内网服务转发，将连接分发至 SSHPiper Pod 容器的 2222 端口
└──────────────┬───────────────┘
               │
               ▼
┌──────────────────────────────┐
│     SSHPiper Pod             │  【控制大脑】：拦截握手提取用户名，查询 Pipe 规则，代理后段连接
└──────────────┬───────────────┘
               │
               ▼ (跨命名空间 FQDN 路由)
┌──────────────────────────────┐
│     Notebook SSH Service     │  用户命名空间下的 Service (notebook-xxx-ssh:22 -> Pod:22)
└──────────────┬───────────────┘
               │
               ▼
┌──────────────────────────────┐
│     Notebook Pod (sshd)      │  真正运行用户容器环境的 Pod，完成验证后打通交互式 Shell
└──────────────────────────────┘
```

### 2. 流量流转时序图

```text
用户客户端 (User PC)              APISIX 网关              SSHPiper 代理              Notebook Pod (sshd)
       │                              │                         │                          │
       │─── 1. 发起 SSH 物理连接 ────►│                         │                          │
       │    (ssh -p 30177)            │                         │                          │
       │                              │─── 2. 透明转发 TCP 流 ──►│                          │
       │                              │    (Service:22->Pod:2222)│                          │
       │                              │                         │ [解析用户名]              │
       │                              │                         │ [读取对应 Pipe 资源]      │
       │                              │                         │                          │
       │◄── 3. 发起密钥质询 (Challenge)─────────────────────────│                          │
       │─── 4. 返回【用户私钥】签名 ────────────────────────────►│                          │
       │                              │                         │ [使用用户公钥验签通过]    │
       │                              │                         │                          │
       │                              │                         │── 5. 发起新 SSH 连接 ───►│
       │                              │                         │   (携带【平台私钥】签名) │
       │                              │                         │                          │ [sshd 读取平台公钥]
       │                              │                         │                          │ [验签通过，放行连接]
       │                              │                         │◄─ 6. 回应连接成功 ───────│
       │                              │                         │                          │
       │◄───────────────────────────── 7. 桥接并建立双端会话 ──────────────────────────────►│
       │ (用户成功登录交互式终端)
```

#### 详细步骤说明：

1. **第一阶段：建立前段 TCP & SSHPiper 拦截用户名**
   - **步骤 1**：用户在本地终端执行登录命令 `ssh -p 30177 <username>@<node-ip>`，流量到达 APISIX 网关暴露的 `30177` 物理端口。
   - **步骤 2**：APISIX 网关接收到连接，根据绑定的 TCP 路由规则，直接将原始加密 TCP 数据流透明透传至 `sshpiper` 内网 Service。Service 负载分发至 SSHPiper Pod 容器的 `2222` 端口。
   - **拦截提取**：SSHPiper 接管连接握手，解密读取到用户登录用户名（例如 `zzz-notebook-1`），并在 `kubeflow` 命名空间下匹配对应的 `Pipe` 资源。

2. **第二阶段：前段身份验证（用户公钥验签）**
   - **步骤 3**：SSHPiper 作为 SSH 服务器，向用户客户端发起 SSH 密钥质询。
   - **步骤 4**：用户客户端收到质询后，使用保存在本地电脑上的**用户私钥**对质询进行加密签名并返回给 SSHPiper。
   - **身份确认**：SSHPiper 读取 `Pipe.spec.from.authorizedKeys`（用户公钥），对用户发回的签名进行解密与验签。如果通过，证明客户端确实拥有该 Notebook 对应的用户私钥，前段认证成功。

3. **第三阶段：后段免密代理（平台私钥验签）**
   - **步骤 5**：前段认证成功后，SSHPiper 在内存中向后端的 Kubernetes 服务 FQDN（例如 `notebook-1-ssh.zzz.svc.cluster.local:22`）发起一条全新的 SSH 连接。它读取 `Pipe.spec.to.privateKeySecret`（保存在 `kubeflow` 下的平台私钥），使用**平台私钥**对后段握手数据包进行签名，发送给 Notebook Pod。
   - **代理免密**：Notebook Pod 容器内的 `sshd` 进程接收连接，读取挂载到容器内的 `/root/.ssh/authorized_keys` 配置文件（该文件包含了**平台公钥**）。`sshd` 使用该公钥对 SSHPiper 提供的私钥签名进行解密验签。

4. **第四阶段：双端会话对接**
   - **步骤 6-7**：后半段连接成功打通后，SSHPiper 将前段与后段在内存中完全“桥接”连通。用户便能顺畅地登录到 Pod 的交互式终端环境中。

---

## 二、端口对应与流转链路

理解端口在各层组件中的翻译和映射是排查连接故障（如 Timeout / Refused）的关键。

| 链路环节 | 组件与资源名称 | 端口对应 | 转换角色与设计目的 |
|:---|:---|:---|:---|
| **外部客户端** | 用户连接端 | `30177` (NodePort) | 公网暴露的唯一物理端口。用户在命令行中指定的 `-p 30177`。 |
| **入口 Service** | `apisix-gateway` (apisix 命名空间) | `22 -> 9100` | 将公网的 `30177` 流量，在 K8s 边界重定向为网关内网 Service 的 `22` 端口，并转发至 APISIX 容器实际监听的 `9100` 端口。 |
| **APISIX 路由匹配** | Stream Proxy (APISIX Pod) | `9100` | APISIX Stream 引擎绑定此端口，匹配到流路由配置（`ingressPort: 9100`）后，原样透传给 `sshpiper` 服务。 |
| **代理 Service** | `sshpiper` (kubeflow 命名空间) | `22 -> 2222` | 内网服务中转。将 APISIX 传来的端口为 `22` 的连接，分发至 SSHPiper Pod 的真实监听端口 `2222`。 |
| **业务 Service** | `notebook-xxx-ssh` (用户命名空间) | `22 -> 22` | 每个 Notebook 专属 Service，负责接收 SSHPiper 的后半段连接，将其导入 Pod 内监听 `22` 端口的 `sshd` 进程。 |

---

## 三、双段认证机制（核心设计）

Neptune 机器学习平台基于 **“双密钥系统”** 实现了一种既安全又透明的 SSH 反向代理机制。

```text
               【前段连接：验证用户】                  【后段连接：免密代理】
  用户客户端 ────────────────────────► SSHPiper ────────────────────────► Notebook Pod
 (持有: 用户私钥)                    (持有: 用户公钥)                    (持有: 平台公钥)
                                     (持有: 平台私钥)
```

### 1. 密钥分工与落位

| 密钥类别 | 密钥角色 | 产生与存储机制 | 作用与验证链路 |
|:---|:---|:---|:---|
| **用户密钥对** | **用户公钥** | 用户自行生成并上传至个人中心。存储在 `kubeflow` 的 `Pipe` 资源中，同时也被挂载到 Notebook Pod 的 `/root/.ssh/authorized_keys` 中。 | 在 **前半段（用户 ➔ SSHPiper）** 验证连接者的真实身份。 |
| **平台/Pipe密钥对** | **平台私钥 & 公钥** | 实例创建时，由 Neptune Go 后端动态调用 OpenSSL 生成。其中**私钥**存储为 `kubeflow` 下的 Secret `pipe-<instance>-ssh-private-key` 供 SSHPiper 读取；**公钥**与用户公钥拼接后挂载到 Pod 中。 | 在 **后半段（SSHPiper ➔ Notebook Pod）** 代理用户进行免密登录。 |

### 2. 核心设计深度问答

#### Q1：为什么不能让用户的“用户私钥”一路传给 Pod 验证，非要引入“平台密钥对”？
- **SSH 协议设计限制**：SSH 协议的握手是强加密的。如果只想做简单的网络透传，APISIX stream route 就足够了，但那无法实现根据用户名分发。要读取用户名，SSHPiper 必须“截断”连接，作为代理商**终止第一段连接的握手**，获取到用户名后再去**重新发起第二段连接**。
- **用户私钥的绝对安全性**：在发起第二段连接时，Pod 的 `sshd` 必然要求对握手包进行签名。然而，出于安全和隐私审计原则，**用户的“私钥”只能存放在用户自己的本地电脑上**，中间代理 SSHPiper 无法、也绝不能获取到用户的私钥。
- **完美的解耦方案**：既然 SSHPiper 无法获取用户私钥，它就无法代替用户向 Pod 发送签名。因此，平台必须产生另一套**“平台密钥对”**。前半段用“用户密钥对”来验明正身，确认用户是该 Notebook 的拥有者；后半段用“平台密钥对”来实现代理登录，完美解决了协议限制。

#### Q2：既然 SSHPiper 已经通过平台私钥代理登录了，Pod 内为什么还要挂载“用户公钥”？
虽然对于正常的 SSHPiper 路由而言，Pod 只需要平台公钥即可登录，但在 Pod 内部同时挂载用户公钥具有极高的架构和运维价值：
1. **允许绕过网关的“直接调试模式”（Kubectl 端口转发）**：
   如果有一天平台网关或 SSHPiper 发生了故障，管理员或用户可以通过 K8s 原生命令直接建立端口转发：
   ```bash
   kubectl port-forward pod/notebook-xxx-0 -n namespace 2222:22
   ```
   在此直连场景下，TCP 流量完全绕过了网关代理，客户端直接与 Pod 内部 `sshd` 通信。此时 Pod 内部必须有用户的公钥，用户才能通过本地私钥直接登录进去。
2. **容器自主性与未来架构兼容 (Container Autonomy)**：
   从所有权语义来看，Notebook 实例属于该用户。将用户公钥写入 Pod 的 `authorized_keys` 符合容器自主的标准。未来如果平台进行技术升级，取消代理中转而改用直接分配公网 IP，当前的容器镜像与配置无需做任何调整。
3. **实现细节**：
   Neptune 后端服务在创建 SSH Secret 时，会将 **用户公钥** 与 **平台公钥** 通过换行符（`
`）拼接在一起：
   ```text
   <用户公钥>
   <平台/Pipe公钥>
   ```
   最终将其整体以 `authorized_keys` 挂载到 Pod 内部。这兼顾了**网关路由**与**直连调试**两种场景。

#### Q3：既然 APISIX 支持 Stream Proxy (L4 代理)，为什么不能直接配置路由直连 Pod，非要多引入一个 SSHPiper 代理？
- **L4 与 L7 路由能力的物理差异**：APISIX 的 Stream Proxy（流代理）工作在 L4 传输层，它做路由决策时仅能识别 IP、端口以及 TLS SNI（对于未解密的 TLS 连接）。而普通的 SSH 客户端发起的是强加密的纯 TCP 数据流，既没有 SNI 信息，其载荷里的 SSH 登录用户名（如 `zzz-notebook-1`）也是被高强度加密封装的。
- **单一公网端口的寻址难题**：由于平台在集群边界只为所有用户暴露了唯一一个 SSH 入口端口（`30177`），所有的连接流量对于 APISIX 来说是完全一模一样的（源 IP 不同，但目的端口都是网关的 9100 端口）。APISIX 无法在不解密的情况下读取 SSH 用户名，因而无法判定应该把当前 TCP 流路由到哪一个 Notebook Pod。
- **SSHPiper 的核心定位**：SSHPiper 是一个“能听懂 SSH 协议”的 L7/应用层反向代理。它拦截并解开前半段连接握手，安全提取用户名（如 `zzz-notebook-1`）后，通过 K8s API 查询对应的 `Pipe` 资源，得知后端 Notebook FQDN 地址后，再使用平台私钥向目标 Namespace 下的 Notebook Pod 发起后半段连接。
- **若不使用 SSHPiper 的代价**：如果硬要去掉 SSHPiper、只用 APISIX 路由直连，要么需要“一实例一端口”（即为每个用户分配一个独立的 NodePort 端口，消耗海量公网端口资源，且极大增加端口暴露的安全风险）；要么需要用户在本地 SSH 配置中增加非常复杂的代理命令（ProxyCommand，如 SSH over WebSockets/TLS），严重损害用户体验。

#### 💡 形象的比喻：保税仓代收包裹
这就像您从海外网购了一件保密商品，但出于安全和海关规定，国外快递不能直接送到您家：
1. **网络不可达（网络限制）**：国外快递员（用户）根本不知道您家的内网地址（Pod IP），必须统一寄到**国家中转保税仓（SSHPiper）**。
2. **包裹需拆检（协议限制）**：为了确认商品合法性（确认 SSH 登录用户名），保税仓（SSHPiper）必须拆开外包装。拆开之后，原来的海外防伪标签（用户私钥签名）就失效了。
3. **重新封条与钥匙（双密钥）**：保税仓检查完后，必须把商品重新打包寄送给您。但是保税仓没有您私人的海外防伪标签和印章。因此，保税仓和您家约定了另一套专属封条（平台密钥）：保税仓盖上“保税仓专用章”（平台私钥），您家门口的收件人（Pod 内部 sshd）手持“保税仓公章图样”（平台公钥）验封放行。

---

## 四、跨命名空间 (Cross-Namespace) 路由原理

在 Neptune 平台中，**整个 K8s 集群仅部署一个全局共享的 SSHPiper 实例**（运行在 `kubeflow` 命名空间下）。
当它接收到外部多租户的连接时，需要把连接分发到各自用户空间（如 `caixukun`, `zzz` 等）的 Pod 内。

### 1. 跨命名空间通信的实现
*   在 `Pipe` 定义中，Go 后端会把 `to.host` 属性配置为目标服务在对应命名空间下的 **FQDN（完全限定域名）**：
    ```yaml
    # notebook.go 1234行
    TargetHost: fmt.Sprintf("%s.%s.svc.cluster.local:22", sshServiceName, nbRef.Namespace)
    # 示例: notebook-eeb47f-ssh.zzz.svc.cluster.local:22
    ```
*   由于 K8s 的 DNS 解析和容器网络默认支持跨命名空间寻址，处于 `kubeflow` 空间下的 SSHPiper 可以直接通过域名连接到 `zzz` 空间下的 Service。
*   这避免了在每个 Namespace 下都重复部署一套 SSHPiper，实现了**单端口多路复用**，极大地节省了系统计算资源。

---

## 五、创建实例时的 SSH 资源落位

创建支持 SSH 的 Notebook 时，Go 后端通过 `createNotebookSSHSecrets` 和 `createNotebookSSHResources` 自动在不同命名空间下落位资源，生命周期由平台统一控制：

```text
                  [Go 后端控制逻辑]
            ┌─────────────┴─────────────┐
            ▼                           ▼
    【kubeflow namespace】       【用户 namespace】
    - 创建 Pipe 规则            - 创建 Service: notebook-xxx-ssh:22
    - 创建 Secret:              - 创建 Secret: notebook-xxx-ssh-key
      pipe-notebook-xxx-          (挂载为 authorized_keys,
      ssh-private-key             包含 用户公钥+平台公钥)
      (包含平台私钥)
```

---

## 六、部署步骤

### 1. 生成服务端 Host Key
SSHPiper 在接收外部握手时，本身作为 SSH 服务器需要一对主机钥匙：
```bash
ssh-keygen -f ssh_host_ed25519_key -t ed25519 -N ''
```

### 2. 创建 SSHPiper Server Key Secret
在部署 SSHPiper 前，在 `kubeflow` 命名空间创建对应的主机 Secret：
```bash
kubectl create secret generic sshpiper-server-key \
  -n kubeflow \
  --from-file=server_key=ssh_host_ed25519_key
```

### 3. 安装 CRD 与部署清单
应用本目录下的 `sshpiper_crd.yaml` 注册 Pipe 资源：
```bash
kubectl apply -f sshpiper_crd.yaml
```
然后部署 SSHPiper 数据面组件：
```bash
kubectl apply -f deploy.yml
```

---

## 七、故障排查与诊断指南

### 1. 外部 SSH 连接提示 Timeout / Refused
*   **第 1 步：检查外部网关端口**
    检查 K8s 节点安全组，并确保 `apisix-gateway` 服务成功暴露了 NodePort `30177`。
*   **第 2 步：检查 APISIX Stream 监听**
    确认 APISIX Ingress Controller 成功同步了流路由：
    ```bash
    kubectl get apisixroute stream-sshpiper -n kubeflow -o yaml
    ```
    并检查 APISIX 运行配置是否正常开启了流代理（可通过 netstat 查看 Pod 内部 `9100` 端口是否处于 `LISTEN` 状态）。
*   **第 3 步：检查 SSHPiper 接收端**
    确保 `sshpiper` Service 中的端口转换正确（`Port: 22` 必须对应 `targetPort: 2222`）。如果对应错配为 22，APISIX 转发流量时将被拒绝连接。

### 2. 用户连接立刻被拒绝，提示 Permission Denied (publickey)
*   **第 1 步：检查 Pipe 状态**
    检查 Pipe 是否在 `kubeflow` 空间正确落位，并且 `From.Username` 是否与用户的登录命令行严格对应：
    ```bash
    kubectl get pipe -n kubeflow -o yaml
    ```
*   **第 2 步：查看 SSHPiper 日志（排查分界点）**
    ```bash
    kubectl logs -n kubeflow deployment/sshpiper --tail=50
    ```
    *   *有连接建立日志*（`ssh connection pipe created...`）：说明 APISIX 到 SSHPiper 这一段完全畅通，前段匹配成功。故障在后半段。请去检查 Notebook 空间的 `-ssh` Service 状态、Pod 内部 `/root/.ssh/authorized_keys` 是否正确挂载了平台公钥，以及 Pod 容器内的 `sshd` 是否在运行。
    *   *无连接日志*：说明握手在外侧被拒绝，请检查上传的公钥格式是否合法，或客户端是否错误指定了私钥。

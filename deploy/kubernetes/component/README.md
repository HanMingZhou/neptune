# Kubernetes 组件部署说明

## 目录说明

本目录用于部署 Neptune 依赖的基础组件，当前默认覆盖：

- `Volcano`
- `Kubeflow` 最小控制器
- `APISIX`
- `SSHPiper`

相关脚本：

- [deploy_all.sh](/Users/jerrytom/go/src/test/kubeflow/neptune/deploy/kubernetes/component/deploy_all.sh)
- [uninstall_all.sh](/Users/jerrytom/go/src/test/kubeflow/neptune/deploy/kubernetes/component/uninstall_all.sh)

相关子文档：

- [APISIX](/Users/jerrytom/go/src/test/kubeflow/neptune/deploy/kubernetes/component/apisix/readme.md)
- [SSHPiper](/Users/jerrytom/go/src/test/kubeflow/neptune/deploy/kubernetes/component/sshpiper/readme.md)
- [Kubeflow 最小控制器安装说明](/Users/jerrytom/go/src/test/kubeflow/neptune/deploy/kubernetes/component/kubeflow/deploy.md)

## 默认安装内容

`deploy_all.sh` 当前按“最小控制器安装”执行，Kubeflow 只安装：

- `cert-manager`
- `kubeflow-issuer`
- `istio-crds`
- `notebook-controller`
- `tensorboard-controller`

不会安装完整 Kubeflow 组件栈。

## 一键部署

```bash
cd deploy/kubernetes/component
./deploy_all.sh
```

执行脚本后会逐项询问是否继续部署，可输入 `y` 或 `n`：

```text
是否部署 Volcano ? [y/n]:
是否部署 Kubeflow ? [y/n]:
是否部署 APISIX ? [y/n]:
是否部署 SSHPiper ? [y/n]:
```

脚本默认行为：

- `Volcano` 和 `APISIX` 使用 Helm 安装或升级
- `Kubeflow manifests` 固定到 `v1.10.0`
- 自动兼容 `apps`/`applications` 两种 manifests 目录结构
- 自动兼容不同版本的 `istio-crds` 路径
- `cert-manager` 安装后会等待 webhook CA bundle 就绪，再继续安装 `kubeflow-issuer`
- 如果 `kubeflow/sshpiper-server-key` 不存在，且本地存在 `./sshpiper/ssh_host_ed25519_key`，会自动创建 Secret
- `APISIX` Chart 会先走本地缓存，远程拉取失败时自动重试
- 检测到旧版 `GatewayProxy apisix` 与 Helm 新版 `apisix-config` 冲突时，会自动清理默认旧资源
- 交互式执行时会逐项询问是否部署各组件，可输入 `y` 或 `n`
- 支持通过环境变量选择只部署部分组件，例如只部署 `APISIX` 和 `SSHPiper`
- 如果检测到 `Volcano` 已经存在但不是当前 Helm Release `volcano/volcano-system` 管理，会自动跳过该步骤，避免 Helm ownership 冲突

## 一键卸载

```bash
cd deploy/kubernetes/component
./uninstall_all.sh
```

默认卸载内容：

- `SSHPiper`
- `APISIX`
- `Kubeflow` 最小控制器
- `Volcano`

默认保留内容：

- `kubeflow` namespace
- `apisix`、`cert-manager`、`volcano-system` namespace
- 本地 `kubeflow/manifests` 目录

## 常用环境变量

部署脚本常用变量：

- `VOLCANO_CHART_FILE=/path/to/volcano-1.14.1.tgz`
- `APISIX_CHART_FILE=/path/to/apisix-2.13.0.tgz`
- `HELM_RETRY_COUNT=5`
- `HELM_RETRY_DELAY_SECONDS=10`
- `SSHPIPER_SERVER_KEY_FILE=/path/to/ssh_host_ed25519_key`
- `APISIX_ADMIN_KEY=<your-admin-key>`
- `DEPLOY_VOLCANO=false`
- `DEPLOY_KUBEFLOW=false`
- `DEPLOY_APISIX=false`
- `DEPLOY_SSHPIPER=false`

卸载脚本常用变量：

- `DELETE_COMPONENT_NAMESPACES=true`
- `DELETE_KUBEFLOW_NAMESPACE=true`
- `REMOVE_SSHPIPER_SECRET=true`
- `REMOVE_LOCAL_MANIFESTS=true`

示例：

```bash
SSHPIPER_SERVER_KEY_FILE=/path/to/ssh_host_ed25519_key ./deploy_all.sh

DEPLOY_VOLCANO=false \
DEPLOY_KUBEFLOW=false \
APISIX_CHART_FILE=/path/to/apisix-2.13.0.tgz \
SSHPIPER_CRD_URL=/path/to/sshpiper_crd.yaml \
./deploy_all.sh

DELETE_COMPONENT_NAMESPACES=true \
DELETE_KUBEFLOW_NAMESPACE=true \
REMOVE_LOCAL_MANIFESTS=true \
./uninstall_all.sh
```

## 目录内文件说明

| 路径 | 说明 |
|------|------|
| `deploy_all.sh` | 组件一键部署脚本 |
| `uninstall_all.sh` | 组件一键卸载脚本 |
| `apisix/` | APISIX 相关说明和参考模板 |
| `sshpiper/` | SSHPiper 资源清单与文档 |
| `kubeflow/` | Kubeflow 最小控制器安装说明和本地 manifests 目录 |
| `volcano/` | Volcano 相关说明 |
| `kong/` | Kong 迁移预研文档，当前不是默认方案 |

## 注意事项

- `deploy_all.sh` 只负责组件层，不会自动部署 `deploy/kubernetes/server` 和 `deploy/kubernetes/web`
- 仓库内的 `sshpiper/ssh_host_ed25519_key` 只建议用于开发测试环境，生产环境应使用自定义私钥
- 如遇到 Helm 仓库临时 `502/503`，可先在浏览器下载 `.tgz`，再通过 `VOLCANO_CHART_FILE` 或 `APISIX_CHART_FILE` 指定本地文件
- 如果组件曾经由旧脚本或手工命令部署过，建议优先执行一次清理，再重跑最新脚本；如果只是补装剩余组件，优先使用 `DEPLOY_*` 环境变量跳过已存在组件

<p align="center">
  <img src="docs/readme/images/platform.png" width="200" alt="Neptune ML Platform" />
</p>

<h1 align="center">Neptune · 机器学习平台</h1>

<p align="center">
  基于 Kubernetes 的一站式 AI 基础设施管理平台，提供 GPU 容器实例、分布式训练、模型推理和全生命周期资源管理。
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24-00ADD8?logo=go" />
  <img src="https://img.shields.io/badge/Vue-3-4FC08D?logo=vuedotjs" />
  <img src="https://img.shields.io/badge/Kubernetes-Ready-326CE5?logo=kubernetes" />
  <img src="https://img.shields.io/badge/Kubeflow-Integrated-FF6F00" />
  <img src="https://img.shields.io/badge/License-Proprietary-red" />
</p>

---

## 📖 目录

- [功能概述](#-功能概述)
- [技术架构](#-技术架构)
- [项目结构](#-项目结构)
- [环境要求](#-环境要求)
- [快速启动](#-快速启动)
- [界面预览](#-界面预览)
- [部署方式](#-部署方式)
- [配置说明](#-配置说明)

---

## ✨ 功能概述

### 🖥️ 容器实例（Notebook）

- 基于 Kubeflow Notebook Controller 的 GPU 交互式计算环境
- 支持 JupyterLab、VS Code Server 等多种 IDE
- 内置 SSH 连接（通过 SSHPiper 动态路由）
- TensorBoard 可视化训练监控
- 实例生命周期管理（创建、启动、停止、删除）

### 🏋️ 分布式训练

- 基于 Volcano 的分布式训练任务调度
- 支持 PyTorch DDP、MPI 等主流框架
- Master/Worker 多节点配置
- 实时日志流和任务状态监控

### 🚀 模型推理

- 在线推理服务部署与管理
- API Key 访问控制和流量策略
- 服务自动扩缩容

### 📦 资源管理

- **文件存储**：基于 Kubernetes PVC 的持久化存储，支持容器/训练/推理间共享
- **SSH 密钥**：密钥对管理，支持密钥登录容器实例
- **镜像管理**：基础镜像、社区镜像和自定义镜像

### 💰 费用与账单

- 按量计费 / 包时计费 多种计费模式
- 实时余额和消费趋势监控
- 交易记录和发票管理

### ⚙️ 系统管理

- 多集群管理（K8s 集群动态注册）
- 基于 Casbin 的 RBAC 权限控制
- 角色管理、菜单管理、API 管理
- 操作审计和访问日志
- 产品和定价管理
- 中英文双语支持

---

## 🏗️ 技术架构

```
+------------------------------------------------------------------+
|                    Browser / SSH Client                          |
+----------------------+-------------------------------+------------+
                       |                               |
                 HTTP / WS (:80)                 SSH (:22)
                       |                               |
                       v                               v
                 +--------------------+         +--------------+
                 |   APISIX Gateway   |         | APISIX Stream|
                 |    (统一入口)      |         |  (K8s 模式)  |
                 +---------+----------+         +------+-------+
                           |                           |
      +--------------------+--------------------+      |
      |                                         |      v
      | /                                       |  +----------+
      v                                         v  | SSHPiper |
+---------------------+                +-----------------------+
|    Neptune Web      |                |    Neptune Server     |
|   (Vue3 + Nginx)    |                |      (Go + Gin)       |
+---------------------+                +----+-------------+-----+
                                            |             |
                                            v             v
                                        +-------+     +---------+
                                        | MySQL |     |  Redis  |
                                        +-------+     +---------+
                                                 \
                                                  v
                                      +---------------------------+
                                      |      Kubernetes API       |
                                      | Kubeflow / Volcano / PVC  |
                                      +---------------------------+
```

- HTTP 统一入口：`APISIX` 按路径转发，`/` 到前端，`/aiInfra/*` 到后端。
- WebSocket（日志流/终端）同样走 APISIX 的 HTTP 网关链路。
- SSH 流量链路（APISIX Stream -> SSHPiper）在 Kubernetes 部署中启用，Compose 主要覆盖 HTTP 网关场景。

### 技术栈

| 层级 | 技术 |
|------|------|
| **前端** | Vue 3 + Vite + Element Plus + UnoCSS + ECharts |
| **后端** | Go 1.24 + Gin + GORM + Casbin + Viper |
| **数据库** | MySQL 8.0 + Redis 7 |
| **容器编排** | Kubernetes + Kubeflow + Volcano |
| **网关** | APISIX（HTTP 反向代理 + TCP Stream） |
| **SSH** | SSHPiper（基于用户名的动态 SSH 路由） |

---

## 📂 项目结构

```
neptune/
├── server/                  # 后端服务 (Go)
│   ├── api/v1/              # API 接口层
│   ├── config/              # 配置结构体
│   ├── core/                # 核心组件初始化 (Zap, Viper, Server)
│   ├── global/              # 全局变量
│   ├── initialize/          # 初始化 (Router, DB, Redis, K8s)
│   ├── middleware/          # Gin 中间件
│   ├── model/               # 数据模型
│   ├── router/              # 路由注册
│   ├── service/             # 业务逻辑
│   ├── mcp/                 # MCP 协议服务
│   ├── Dockerfile           # 后端镜像构建
│   ├── config.yaml          # 生产配置
│   └── config.dev.yaml      # 开发配置
│
├── web/                     # 前端 (Vue 3)
│   ├── src/
│   │   ├── api/             # API 请求封装
│   │   ├── components/      # 公共组件
│   │   ├── i18n/locales/    # 国际化 (zh-CN / en-US)
│   │   ├── pinia/           # 状态管理
│   │   ├── view/            # 页面视图
│   │   └── router/          # 前端路由
│   ├── Dockerfile           # 前端镜像构建
│   └── nginx.deploy.conf    # 部署态前端 Nginx 配置（docker compose/容器直连调试）
│
├── deploy/                  # 部署配置
│   ├── docker-compose/      # Docker Compose 部署（含 APISIX 网关）
│   │   ├── docker-compose.yaml
│   │   ├── config.yaml
│   │   └── apisix/
│   │       ├── config.yaml  # APISIX 运行配置（standalone）
│   │       └── apisix.yaml  # APISIX 路由配置
│   └── kubernetes/          # Kubernetes 部署
│       ├── server/          # 后端 K8s 资源
│       ├── web/             # 前端 K8s 资源
│       └── component/       # 依赖组件一键部署
│           ├── kubeflow/    # Notebook/TB Controller
│           ├── volcano/     # 训练任务调度器
│           ├── apisix/      # API 网关
│           └── sshpiper/    # SSH 路由
│
└── docs/                    # 文档
```

---

## 📋 环境要求

### 开发环境

| 依赖 | 版本 |
|------|------|
| Go | >= 1.24 |
| Node.js | >= 20 |
| MySQL | >= 8.0 |
| Redis | >= 6.0 |

### 生产环境

| 依赖 | 版本 |
|------|------|
| Kubernetes | >= 1.28 |
| Helm | >= 3.x |
| Docker | >= 24.x |

---

## 🚀 快速启动

### 方式一：本地开发

**1. 启动后端**

```bash
cd server

# 复制开发配置（首次）
# 修改 config.dev.yaml 中的 MySQL 和 Redis 连接信息

# 启动服务（默认端口 8001）
go run main.go
```

**2. 启动前端**

```bash
cd web

# 安装依赖
npm install

# 启动开发服务器（默认端口 5173）
npm run dev
```

访问 http://localhost:5173

### 方式二：Docker Compose

```bash
# 在项目根目录执行
docker compose -f deploy/docker-compose/docker-compose.yaml up -d --build
```

服务启动后：
- **统一网关入口（APISIX）**: http://localhost
- **前端页面**: http://localhost
- **前端页面（直连调试）**: http://localhost:8080
- **后端 API（经 APISIX）**: http://localhost/aiInfra
- **后端 API（直连排障）**: http://localhost:8888/aiInfra

---

## 🖼️ 界面预览

### 容器实例（Notebook）
![容器实例示例 1](docs/readme/images/notebook-example-1.png)
![容器实例示例 2](docs/readme/images/notebook-example-2.png)

### 训练任务
![训练任务示例 1](docs/readme/images/training-example-1.png)
![训练任务示例 2](docs/readme/images/training-example-2.png)

### 推理服务

- 当前文档中暂无推理服务截图，后续可继续补充到 `./images/` 目录。

## 📦 部署方式

### Docker Compose 部署

适用于**单机部署、测试环境**。

```bash
docker compose -f deploy/docker-compose/docker-compose.yaml up -d --build
```

包含服务：APISIX、Neptune Web、Neptune Server、MySQL、Redis

路由关系：
- `http://localhost/` -> `neptune-web:8080`
- `http://localhost:8080/` -> `neptune-web:8080`（直连调试）
- `http://localhost/aiInfra/*` -> `neptune-server:8888`

> ⚠️ Docker Compose 模式不包含 Kubernetes 相关功能（Notebook、Training 等需要 K8s 集群支持）

### Kubernetes 部署

适用于**生产环境**。

`deploy_all.sh` 只负责部署基础组件，不会自动部署 Neptune 平台自己的前后端服务和平台入口路由。

**1. 部署依赖组件**

```bash
# 一键部署 Volcano、Kubeflow、APISIX、SSHPiper
bash deploy/kubernetes/component/deploy_all.sh
```

**2. 创建 Neptune 命名空间**

```bash
kubectl apply -f deploy/kubernetes/namespace.yaml
```

**3. 部署后端**

```bash
kubectl apply -f deploy/kubernetes/server/
```

**4. 部署前端**

```bash
kubectl apply -f deploy/kubernetes/web/
```

**5. 部署平台 APISIX 入口路由**

```bash
kubectl apply -f deploy/kubernetes/component/apisix/neptune-platform-route.yaml
```

> 📝 部署前请修改 Deployment 中的镜像地址和 Ingress 域名

---

## ⚙️ 配置说明

| 配置文件 | 用途 |
|----------|------|
| `server/config.dev.yaml` | 本地开发配置 |
| `server/config.yaml` | 生产配置模板 |
| `deploy/docker-compose/config.yaml` | Docker Compose 环境配置 |
| `deploy/docker-compose/apisix/config.yaml` | Compose 下 APISIX 运行配置（standalone） |
| `deploy/docker-compose/apisix/apisix.yaml` | Compose 下 APISIX 路由配置 |
| `deploy/kubernetes/server/gva-server-configmap.yaml` | K8s 环境配置 |
| `deploy/kubernetes/web/neptune-web-nginx-configmap.yaml` | K8s 前端 Nginx 配置 |

> 本地开发默认走 Vite dev server，不使用 Nginx。

### 关键配置项

```yaml
system:
  addr: 8888                # 服务监听端口
  router-prefix: "aiInfra"  # 路由全局前缀

mysql:
  path: "127.0.0.1"         # 数据库地址
  db-name: "aiInfra"        # 数据库名

redis:
  addr: "127.0.0.1:6379"    # Redis 地址

apisix:
  enabled: true             # 启用 APISIX 网关
  base-domain: "localhost"  # 网关访问域名
  auth-uri: "http://neptune-server:8888/aiInfra/api/v1/apisix/auth"
  http-port: 80

sshpiper:
  host: "127.0.0.1"         # SSHPiper 地址
  port: 22                  # SSHPiper 端口
```

---

## 📄 License

Proprietary - All Rights Reserved

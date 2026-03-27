<p align="center">
  <img src="docs/platform.png" width="200" alt="Neptune ML Platform" />
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
|                    Browser / SSH Client                           |
+---------------+----------------------------------+---------------+
                |                                  |
          HTTP (:80)                          SSH (:22)
                |                                  |
                v                                  v
+---------------------+              +--------------------+
|    Neptune Web      |              |   APISIX Gateway   |
|   (Vue3 + Nginx)    |              |   HTTP + Stream    |
+---------+-----------+              +------+-------+-----+
          |                                 |       |
          | /aiInfra                        |       |
          v                                 |       v
+-----------------------+                   |  +----------+
|    Neptune Server     |<------------------+  | SSHPiper |
|     (Go + Gin)        |                      +----+-----+
|                       |                           |
|  - Notebook API       |                           v
|  - Training API       |                    +--------------+
|  - Inference API      |                    | Notebook Pod |
|  - APISIX Route API   |                    |  (sshd:22)   |
|  - System Admin API   |                    +--------------+
+---+-------------+-----+
    |             |
    v             v
+-------+  +---------+   +------------------------------+
| MySQL |  |  Redis  |   |       Kubernetes API         |
+-------+  +---------+   |  - Kubeflow (Notebook/TB)    |
                          |  - Volcano  (Training)       |
                          |  - PVC      (Storage)        |
                          +------------------------------+
```

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
│   └── nginx.conf           # Nginx 反向代理配置
│
├── deploy/                  # 部署配置
│   ├── docker/              # All-in-One Docker 构建
│   ├── docker-compose/      # Docker Compose 部署
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
docker-compose -f deploy/docker-compose/docker-compose.yaml up -d
```

服务启动后：
- **前端**: http://localhost
- **后端 API**: http://localhost:8888/aiInfra

---

## 📦 部署方式

### Docker Compose 部署

适用于**单机部署、测试环境**。

```bash
docker-compose -f deploy/docker-compose/docker-compose.yaml up -d
```

包含服务：Neptune Server、Neptune Web、MySQL、Redis

> ⚠️ Docker Compose 模式不包含 Kubernetes 相关功能（Notebook、Training 等需要 K8s 集群支持）

### Kubernetes 部署

适用于**生产环境**。

**1. 部署依赖组件**

```bash
# 一键部署 Volcano、Kubeflow、APISIX、SSHPiper
bash deploy/kubernetes/component/deploy_all.sh
```

**2. 部署 Neptune 服务**

```bash
# 部署后端
kubectl apply -f deploy/kubernetes/server/

# 部署前端
kubectl apply -f deploy/kubernetes/web/
```

> 📝 部署前请修改 Deployment 中的镜像地址和 Ingress 域名

---

## ⚙️ 配置说明

| 配置文件 | 用途 |
|----------|------|
| `server/config.dev.yaml` | 本地开发配置 |
| `server/config.yaml` | 生产配置模板 |
| `deploy/docker-compose/config.yaml` | Docker Compose 环境配置 |
| `deploy/kubernetes/server/gva-server-configmap.yaml` | K8s 环境配置 |

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
  auth-uri: "..."           # 认证回调地址

sshpiper:
  host: "127.0.0.1"         # SSHPiper 地址
  port: 22                  # SSHPiper 端口
```

---

## 📄 License

Proprietary - All Rights Reserved

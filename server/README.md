## server项目结构

```shell
├── api
│   └── v1
├── config
├── core
├── docs
├── global
├── initialize
│   └── internal
├── middleware
├── model
│   ├── request
│   └── response
├── packfile
├── resource
│   ├── excel
│   ├── page
│   └── template
├── router
├── service
├── source
└── utils
    ├── timer
    └── upload
```

| 文件夹       | 说明                    | 描述                        |
| ------------ | ----------------------- | --------------------------- |
| `api`        | api层                   | api层 |
| `--v1`       | v1版本接口              | v1版本接口                  |
| `config`     | 配置包                  | config.yaml对应的配置结构体 |
| `core`       | 核心文件                | 核心组件(zap, viper, server)的初始化 |
| `docs`       | swagger文档目录         | swagger文档目录 |
| `global`     | 全局对象                | 全局对象 |
| `initialize` | 初始化 | router,redis,gorm,validator, timer的初始化 |
| `--internal` | 初始化内部函数 | gorm 的 longger 自定义,在此文件夹的函数只能由 `initialize` 层进行调用 |
| `middleware` | 中间件层 | 用于存放 `gin` 中间件代码 |
| `model`      | 模型层                  | 模型对应数据表              |
| `--request`  | 入参结构体              | 接收前端发送到后端的数据。  |
| `--response` | 出参结构体              | 返回给前端的数据结构体      |
| `packfile`   | 静态文件打包            | 静态文件打包 |
| `resource`   | 静态资源文件夹          | 负责存放静态文件                |
| `--excel` | excel导入导出默认路径 | excel导入导出默认路径 |
| `--page` | 表单生成器 | 表单生成器 打包后的dist |
| `--template` | 模板 | 模板文件夹,存放的是代码生成器的模板 |
| `router`     | 路由层                  | 路由层 |
| `service`    | service层               | 存放业务逻辑问题 |
| `source` | source层 | 存放初始化数据的函数 |
| `utils`      | 工具包                  | 工具函数封装            |
| `--timer` | timer | 定时器接口封装 |
| `--upload`      | oss                  | oss接口封装        |

## 镜像构建与运行

后端镜像入口文件是 [Dockerfile](/Users/jerrytom/go/src/test/kubeflow/neptune/server/Dockerfile)。

当前镜像设计要点：

- 使用多阶段构建，构建阶段基于 `golang:1.24-alpine`
- 运行阶段基于 `alpine:3.20`
- 默认启动二进制为 `/app/server`
- 默认配置文件路径为 `/app/config.yaml`
- 通过 `GVA_CONFIG=/app/config.yaml` 指定配置文件
- 容器以非 root 用户 `app:app` 运行，uid/gid 为 `10001`
- 默认暴露端口 `8888`
- 预创建目录：
  - `/app/log`
  - `/app/uploads/file`
  - `/app/resource/excel`

示例构建命令：

```bash
docker build -t neptune/server:latest -f server/Dockerfile server
```

## Kubernetes 部署说明

Kubernetes 相关清单位于：

- [gva-server-configmap.yaml](/Users/jerrytom/go/src/test/kubeflow/neptune/deploy/kubernetes/server/gva-server-configmap.yaml)
- [gva-server-deployment.yaml](/Users/jerrytom/go/src/test/kubeflow/neptune/deploy/kubernetes/server/gva-server-deployment.yaml)
- [gva-server-service.yaml](/Users/jerrytom/go/src/test/kubeflow/neptune/deploy/kubernetes/server/gva-server-service.yaml)

当前部署约定：

- namespace 为 `neptune`
- Service 名称为 `neptune-server`
- 容器端口为 `8888`
- 健康检查路径为 `/aiInfra/health`
- 路由前缀统一为 `aiInfra`
- 推荐通过 `ConfigMap` 挂载 `config.yaml`

当前 K8s 配置中的关键项：

- `system.router-prefix: "aiInfra"`
- `mysql.path: neptune-mysql`
- `redis.addr: neptune-redis:6379`
- `apisix.auth-uri: http://neptune-server.neptune.svc.cluster.local:8888/aiInfra/api/v1/apisix/auth`
- `sshpiper.host: sshpiper.kubeflow.svc.cluster.local`

## 部署注意事项

- 当前上传目录和日志目录都在容器文件系统内，生产环境建议再挂载 `PVC`
- 如果通过 `ConfigMap` 挂载 `/app/config.yaml`，需要保证镜像内的 `GVA_CONFIG` 仍然指向该路径
- 如果调整了 `router-prefix`，需要同步更新前端 `VITE_BASE_API`、Nginx 反向代理路径和 K8s 健康检查路径
- 后端镜像已按非 root 运行方式处理，K8s `securityContext` 需要与 uid/gid `10001` 对齐

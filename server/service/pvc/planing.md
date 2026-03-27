# 开发规范
## 项目结构
```bash
项目结构分为api router model service，所有的router url以api/v1为起始，然后加业务名称，除了必须用get外，其余全部用post。
所有的service 入参的req和response结构体都要在model层的requset和response创建，不要直接返回db的model字段。
db model和其他model用id关联，model创建完后需要在/server/initialize/gorm.go中通过gorm的迁移方法自动迁移到数据库中
工具放到对应的utils包中
如无特殊要求，不要总结并创建文档
所有的列表都需要有刷新按钮，并且5s/10s刷新一次请求。这个按钮全平台列表统一。
```


# 菜单视图
```bash
平台是按照用户角色来展示菜单的，创建完菜单后需要在平台上的超级管理员菜单下的菜单管理模块添加菜单，然后在角色管理绑定对应的菜单权限。
服务启动时，会给超级管理员添加一些默认的菜单和api权限，代码在neptune/server/source/system/menu.go，neptune/server/source/system/api.go。
```

# 开发功能
```bash
机器学习平台 - 训练任务模块概要设计 
1. 设计原则 (Design Principles)
在进入具体设计前，确立四条不可撼动的核心原则：
 - 强类型约束 (Strict Typing): 
 - 严禁使用字符串拼接 (fmt.Sprintf) 或简单的 yaml.Marshal 生成 K8s 资源。
 - 必须使用 Go Structs (volcano.sh/apis + k8s.io/api) 配合 Typed Client。
 - 数据库禁用json字段类型存储。
设计方案可以根据开发实际过程调整，但是四条核心原则不动！
收益: 编译期查错、防止注入攻击、完美的 IDE 补全、平滑的版本升级。
K8s 为事实来源 (Single Source of Truth): 数据库仅作为用户意图的持久化存储和历史归档。

2. 核心业务架构图
Jupyter Notebook Pod
   ↕ 共享 PVC
Training Job (Volcano)
用户：在 Jupyter 里写代码，提交训练任务
训练：使用同一个代码 / 数据 PVC
连接桥梁：共享存储 (RWX PVC)

3. 功能规划
1️⃣ 文件存储是什么？
3.1 文件存储（File Storage）是用户级的持久共享存储，用于在 Notebook 与训练任务和推理服务之间共享代码、数据与输出结果。
它的特点是：
独立于计算生命周期
可被多个 Pod 同时挂载（RWX）
由平台统一管理
不随训练任务销毁
2️⃣ 文件存储不是什么（非常重要）
明确排除，避免后期被“滥用”：
❌ 不关心底层 PV 实现（NFS / CephFS / EFS）

3.2 整体架构设计（核心关系）
┌────────────────────────┐
│   File Storage (RWX)   │  ←—— 用户创建
│      PVC (Namespace)  │
└──────────┬────────────┘
           │
   ┌───────┴────────┐
   │                │
┌──▼───┐        ┌───▼────────┐
│Jupyter│        │ Training Job│
│ Pod   │        │ (Volcano)  │
└───────┘        └────────────┘


关键点：
PVC 是一等资源
Notebook / Training / Inference 都只是“消费者”
生命周期完全解耦
3.3 前端产品设计
1️⃣ 菜单结构
一级菜单：文件存储
  └── 二级菜单：存储列表
没有其他子菜单，刻意保持极简。
2️⃣ 列表页设计（不暴露技术细节）
表格字段（产品语言）
字段名	说明
存储名称	用户自定义
容量	如 50 GB
使用状态	未使用 / 使用中
创建时间	
操作	扩容 / 删除
🚫 不展示：
PVC 名称
StorageClass
AccessMode
Namespace
VolumeHandle
3️⃣ 创建存储（极简表单）
选择地区后，创建弹窗字段
存储名称
容量
数值 + 单位（GB）
（无其他选项）
💡 所有技术参数由平台兜底：
RWX
默认 StorageClass
Namespace = 用户

4️⃣ 扩容操作
交互规则
仅允许 增大
不允许缩容
扩容即时生效（PVC resize）
前端文案示例：
当前容量：50GB
扩容至：[ 100 ] GB
5️⃣ 删除操作（带“是否被使用”校验）
删除前校验逻辑
若 无 Pod 挂载
允许删除
若 被使用
禁止删除
给出友好提示
前端提示文案示例：
❌ 当前存储正在被以下资源使用：
Notebook: notebook-01
Training Job: resnet-train-03
请先停止相关任务。

3.4 后端资源模型设计（关键）
见：server/model/pvc/volume.go，可根据实际开发情况调整

3.5 Notebook / Training / Inference 如何使用文件存储
1️⃣ Jupyter 挂载方式（固定路径）
系统盘为平台自动创建默认都是10GB,目前都挂载到了/home/jovyan/workspace
数据盘为用户创建的pvc，默挂载到 /home/jovyan/neptune-fs （pvc菜单页面给出挂载到默认路径的提示）

2️⃣ Training Job 挂载方式
在训练创建表单中：
选择文件存储
指定容器内路径（默认同上）
用户感知的是“选择一个存储”，不是 PVC
```





create table notebooks
(
    id                   bigint unsigned auto_increment
        primary key,
    created_at           datetime(3)      null,
    updated_at           datetime(3)      null,
    deleted_at           datetime(3)      null,
    display_name         varchar(191)     null comment 'Notebook名称',
    instance_name        varchar(191)     null comment '实例名称',
    namespace            varchar(191)     null comment '命名空间',
    image                varchar(191)     null comment '镜像地址',
    image_id             bigint unsigned  null comment '镜像ID',
    cpu                  bigint           null comment 'CPU资源限制',
    memory               bigint           null comment '内存资源限制',
    gpu                  bigint           null comment 'GPU资源限制',
    gpu_model            varchar(191)     null comment 'GPU型号',
    v_gpu_number         bigint default 0 null comment 'vGPU数量',
    v_gpu_memory         bigint default 0 null comment 'vGPU显存 (GB)',
    v_gpu_cores          bigint default 0 null comment 'vGPU核心数',
    storage_size         bigint           null comment '存储大小',
    status               varchar(191)     null comment '状态',
    user_id              bigint unsigned  null comment '用户ID',
    cluster_id           bigint unsigned  null comment '集群ID',
    product_id           bigint unsigned  null comment '产品ID',
    pay_type             bigint           null comment '付费类型(1-按量 2-包日 3-包周 4-包月)',
    order_id             bigint unsigned  null comment '订单ID',
    ssh_key_id           bigint unsigned  null comment 'SSH密钥ID',
    enable_ssh_password  tinyint(1)       null comment '是否启用SSH密码登录',
    ssh_password         varchar(191)     null comment 'SSH登录密码',
    enable_tensorboard   tinyint(1)       null comment '是否启用Tensorboard',
    tensorboard_log_path varchar(191)     null comment 'Tensorboard日志路径',
    tensorboard_id       bigint unsigned  null comment 'Tensorboard ID'
);

create index idx_notebooks_deleted_at
    on notebooks (deleted_at);

create index idx_notebooks_image_id
    on notebooks (image_id);

create index idx_notebooks_order_id
    on notebooks (order_id);

create index idx_notebooks_product_id
    on notebooks (product_id);

create index idx_notebooks_tensorboard_id
    on notebooks (tensorboard_id);

create index idx_notebooks_user_id
    on notebooks (user_id);

INSERT INTO aiInfra.notebooks (id, created_at, updated_at, deleted_at, display_name, instance_name, namespace, image, image_id, cpu, memory, gpu, gpu_model, v_gpu_number, v_gpu_memory, v_gpu_cores, storage_size, status, user_id, cluster_id, product_id, pay_type, order_id, ssh_key_id, enable_ssh_password, ssh_password, enable_tensorboard, tensorboard_log_path, tensorboard_id) VALUES (1, '2026-02-24 20:37:05.256', '2026-02-24 20:38:38.643', '2026-02-24 20:38:38.656', '1', 'notebook-16cff2', 'hmz', 'kubeflownotebookswg/jupyter-scipy:v1.8.0', 12, 1, 1, 0, '', 0, 0, 0, 10, 'DELETING', 1, 1, 3, 1, 33, 1, 1, 'bTQpAM4m', 1, 'logs', 151);
INSERT INTO aiInfra.notebooks (id, created_at, updated_at, deleted_at, display_name, instance_name, namespace, image, image_id, cpu, memory, gpu, gpu_model, v_gpu_number, v_gpu_memory, v_gpu_cores, storage_size, status, user_id, cluster_id, product_id, pay_type, order_id, ssh_key_id, enable_ssh_password, ssh_password, enable_tensorboard, tensorboard_log_path, tensorboard_id) VALUES (2, '2026-02-24 20:39:19.939', '2026-03-25 23:44:02.124', null, 'container', 'notebook-837ab9', 'hmz', 'kubeflownotebookswg/jupyter-scipy:v1.8.0', 12, 1, 1, 0, '', 0, 0, 0, 10, 'STOPPED', 1, 1, 3, 1, 48, 1, 1, 'oA8XamHZ', 1, 'logs', 159);
INSERT INTO aiInfra.notebooks (id, created_at, updated_at, deleted_at, display_name, instance_name, namespace, image, image_id, cpu, memory, gpu, gpu_model, v_gpu_number, v_gpu_memory, v_gpu_cores, storage_size, status, user_id, cluster_id, product_id, pay_type, order_id, ssh_key_id, enable_ssh_password, ssh_password, enable_tensorboard, tensorboard_log_path, tensorboard_id) VALUES (3, '2026-03-25 23:45:15.294', '2026-03-28 19:00:03.783', null, 'cd', 'notebook-34bd53', 'hmz', 'kubeflownotebookswg/jupyter-scipy:v1.8.0', 12, 1, 1, 0, '', 0, 0, 0, 10, 'STOPPED', 1, 1, 3, 1, 52, 1, 1, 'MgLmPCYt', 1, '/logs', 161);

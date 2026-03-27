create table products
(
    id               bigint unsigned auto_increment
        primary key,
    created_at       datetime(3)      null,
    updated_at       datetime(3)      null,
    deleted_at       datetime(3)      null,
    name             varchar(100)     null comment '产品名称',
    description      varchar(500)     null comment '产品描述',
    cluster_id       bigint           null comment '集群ID',
    area             varchar(50)      null comment '地域区域',
    node_name        varchar(100)     null comment 'K8s节点名称',
    node_type        varchar(50)      null comment '节点类型',
    cpu_model        varchar(100)     null comment 'CPU型号',
    cpu              mediumint        null comment '每个实例的CPU核数',
    memory           mediumint        null comment '每个实例的内存大小(GB)',
    gpu_model        varchar(50)      null comment 'GPU型号(GPU产品必填,与vGPU/CPU-only互斥)',
    gpu_count        bigint default 0 null comment '每个实例的GPU卡数(GPU产品必填)',
    gpu_memory       mediumint        null comment '单卡GPU显存(GB)',
    used_gpu         bigint default 0 null comment '已使用GPU/实例数',
    v_gpu_count      bigint default 0 null comment 'vGPU数量',
    v_gpu_memory     bigint default 0 null comment '每个实例的vGPU显存(vGPU产品)',
    price_hourly     decimal(20, 6)   null comment '每小时价格(单个实例)',
    price_daily      decimal(20, 6)   null comment '包日价格(单个实例)',
    price_weekly     decimal(20, 6)   null comment '包周价格(单个实例)',
    price_monthly    decimal(20, 6)   null comment '包月价格(单个实例)',
    driver_version   varchar(50)      null comment '显卡驱动版本',
    cuda_version     varchar(50)      null comment 'CUDA版本',
    system_disk      bigint           null comment '系统盘大小(GB)',
    data_disk        bigint           null comment '数据盘大小(GB)',
    status           bigint default 1 null comment '状态(1-上架 0-下架)',
    sort_order       bigint default 0 null comment '排序',
    product_type     bigint default 1 null comment '产品类型(1-计算资源 2-文件存储)',
    storage_class    varchar(100)     null comment 'K8s StorageClass名称',
    storage_price_gb decimal(20, 6)   null comment '每GB每日价格',
    v_gpu_cores      bigint default 0 null comment '每个实例的vGPU核心数(vGPU产品)',
    max_instances    bigint default 0 null comment '最大实例数(自动计算:节点资源/产品规格)',
    used_capacity    bigint default 0 null comment '已使用实例数',
    v_gpu_number     bigint default 0 null comment '每个实例的vGPU数量(vGPU产品,与GPU/CPU-only互斥)',
    version          bigint default 0 null comment '乐观锁版本号',
    constraint idx_product_unique
        unique (deleted_at, product_type, name, cluster_id, node_name, storage_class)
);

create index idx_products_area
    on products (area);

create index idx_products_cluster_id
    on products (cluster_id);

create index idx_products_deleted_at
    on products (deleted_at);

INSERT INTO aiInfra.products (id, created_at, updated_at, deleted_at, name, description, cluster_id, area, node_name, node_type, cpu_model, cpu, memory, gpu_model, gpu_count, gpu_memory, used_gpu, v_gpu_count, v_gpu_memory, price_hourly, price_daily, price_weekly, price_monthly, driver_version, cuda_version, system_disk, data_disk, status, sort_order, product_type, storage_class, storage_price_gb, v_gpu_cores, max_instances, used_capacity, v_gpu_number, version) VALUES (3, '2025-12-28 23:30:04.434', '2026-02-25 00:07:12.063', null, 'CPU 1核', '', 1, '北京', 'minikube', '', 'amd64', 1, 1, '', 0, 0, 2, 0, 0, 1.000000, 10.000000, 100.000000, 1000.000000, '', '', 0, 0, 1, 0, 1, '', 0.000000, 0, 7, 0, 0, 126);
INSERT INTO aiInfra.products (id, created_at, updated_at, deleted_at, name, description, cluster_id, area, node_name, node_type, cpu_model, cpu, memory, gpu_model, gpu_count, gpu_memory, used_gpu, v_gpu_count, v_gpu_memory, price_hourly, price_daily, price_weekly, price_monthly, driver_version, cuda_version, system_disk, data_disk, status, sort_order, product_type, storage_class, storage_price_gb, v_gpu_cores, max_instances, used_capacity, v_gpu_number, version) VALUES (4, '2025-12-28 23:49:19.940', '2026-02-20 16:44:00.802', null, 'standard标准存储', '', 1, '北京', '', '', '', 0, 0, '', 0, 0, 0, 0, 0, 0.000000, 0.000000, 0.000000, 0.000000, '', '', 0, 0, 1, 0, 2, 'standard', 0.000100, 0, 0, 0, 0, 0);
INSERT INTO aiInfra.products (id, created_at, updated_at, deleted_at, name, description, cluster_id, area, node_name, node_type, cpu_model, cpu, memory, gpu_model, gpu_count, gpu_memory, used_gpu, v_gpu_count, v_gpu_memory, price_hourly, price_daily, price_weekly, price_monthly, driver_version, cuda_version, system_disk, data_disk, status, sort_order, product_type, storage_class, storage_price_gb, v_gpu_cores, max_instances, used_capacity, v_gpu_number, version) VALUES (5, '2026-01-05 22:11:07.027', '2026-02-23 21:17:43.933', null, 'CPU 2核', '', 1, '北京', 'minikube', '', 'amd64', 2, 2, '', 0, 0, 2, 0, 0, 1.000000, 1.000000, 1.000000, 1.000000, '', '', 0, 0, 1, 0, 1, '', 0.000000, 0, 5, 0, 0, 0);
INSERT INTO aiInfra.products (id, created_at, updated_at, deleted_at, name, description, cluster_id, area, node_name, node_type, cpu_model, cpu, memory, gpu_model, gpu_count, gpu_memory, used_gpu, v_gpu_count, v_gpu_memory, price_hourly, price_daily, price_weekly, price_monthly, driver_version, cuda_version, system_disk, data_disk, status, sort_order, product_type, storage_class, storage_price_gb, v_gpu_cores, max_instances, used_capacity, v_gpu_number, version) VALUES (6, '2026-02-23 22:31:08.892', '2026-02-23 22:31:08.892', '2026-02-23 22:31:20.773', 'CPU 12核', '', 1, '上海', 'minikube', '', 'amd64', 12, 28, '', 0, 0, 0, 0, 0, 0.000000, 0.000000, 0.000000, 0.000000, '', '', 0, 0, 1, 0, 1, '', 0.000000, 0, 1, 0, 0, 0);
INSERT INTO aiInfra.products (id, created_at, updated_at, deleted_at, name, description, cluster_id, area, node_name, node_type, cpu_model, cpu, memory, gpu_model, gpu_count, gpu_memory, used_gpu, v_gpu_count, v_gpu_memory, price_hourly, price_daily, price_weekly, price_monthly, driver_version, cuda_version, system_disk, data_disk, status, sort_order, product_type, storage_class, storage_price_gb, v_gpu_cores, max_instances, used_capacity, v_gpu_number, version) VALUES (7, '2026-02-23 23:02:12.948', '2026-02-25 00:01:51.703', null, 'CPU 12核', '推理服务', 1, '上海', 'minikube', '', 'amd64', 4, 4, '', 0, 0, 0, 0, 0, 1.000000, 2.000000, 3.000000, 4.000000, '', '', 0, 0, 1, 0, 1, '', 0.000000, 0, 3, 0, 0, 38);

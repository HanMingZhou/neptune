create table volumes
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)      null,
    updated_at datetime(3)      null,
    deleted_at datetime(3)      null,
    name       varchar(191)     null comment '数据盘名称',
    namespace  varchar(191)     null comment '命名空间',
    size       bigint           null comment '大小',
    status     varchar(191)     null comment '状态(Ready/Bound/Error)',
    pvc_name   varchar(191)     null comment 'K8s PVC名称',
    type       bigint default 1 null comment '类型(1:dataset, 2:model)',
    user_id    bigint unsigned  null comment '用户ID',
    cluster_id bigint unsigned  null comment '集群ID',
    product_id bigint unsigned  null comment '产品ID',
    constraint idx_volume_name_ns_deleted
        unique (deleted_at, name, namespace)
);

create index idx_volumes_cluster_id
    on volumes (cluster_id);

create index idx_volumes_product_id
    on volumes (product_id);

create index idx_volumes_user_id
    on volumes (user_id);

INSERT INTO aiInfra.volumes (id, created_at, updated_at, deleted_at, name, namespace, size, status, pvc_name, type, user_id, cluster_id, product_id) VALUES (2, '2026-02-20 12:30:34.741', '2026-02-20 12:30:34.741', null, 'hmz', 'hmz', 10, 'READY', 'vol-1771561834-ufyl', 1, 1, 1, 4);

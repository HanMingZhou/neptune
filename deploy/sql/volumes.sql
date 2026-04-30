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
    area       varchar(50)      null comment '展示地域标签',
    constraint idx_volume_name_ns_deleted
        unique (deleted_at, name, namespace)
);

create index idx_volumes_cluster_id
    on volumes (cluster_id);

create index idx_volumes_product_id
    on volumes (product_id);

create index idx_volumes_user_id
    on volumes (user_id);

INSERT INTO aiInfra.volumes (id, created_at, updated_at, deleted_at, name, namespace, size, status, pvc_name, type, user_id, cluster_id, product_id, area) VALUES (2, '2026-02-20 12:30:34.741', '2026-02-20 12:30:34.741', null, 'hmz', 'hmz', 10, 'READY', 'vol-1771561834-ufyl', 1, 1, 1, 4, null);
INSERT INTO aiInfra.volumes (id, created_at, updated_at, deleted_at, name, namespace, size, status, pvc_name, type, user_id, cluster_id, product_id, area) VALUES (3, '2026-04-01 16:40:25.412', '2026-04-01 16:40:25.412', '2026-04-03 10:16:25.227', '10.255.141.8', 'hmz', 50, 'READY', 'vol-1775032825-tkpw', 1, 1, 1, 4, null);
INSERT INTO aiInfra.volumes (id, created_at, updated_at, deleted_at, name, namespace, size, status, pvc_name, type, user_id, cluster_id, product_id, area) VALUES (4, '2026-04-03 10:17:38.997', '2026-04-16 14:21:51.021', null, '111', 'zzz', 50, 'READY', 'vol-1775182658-plxm', 1, 1, 1, 4, null);
INSERT INTO aiInfra.volumes (id, created_at, updated_at, deleted_at, name, namespace, size, status, pvc_name, type, user_id, cluster_id, product_id, area) VALUES (5, '2026-04-16 14:31:07.435', '2026-04-16 19:52:18.381', null, '模型', 'zzz', 600, 'READY', 'vol-1776321067-hajo', 1, 1, 1, 4, null);

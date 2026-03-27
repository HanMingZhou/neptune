create table inference_mounts
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)             null,
    updated_at datetime(3)             null,
    deleted_at datetime(3)             null,
    service_id bigint unsigned         not null comment '关联的推理服务ID',
    mount_type varchar(20)             null comment '挂载类型(MODEL/DATA)',
    pvc_id     bigint unsigned         null comment 'PVC ID',
    pvc_name   varchar(255)            null comment 'PVC名称',
    sub_path   varchar(255) default '' null comment 'PVC内子路径',
    mount_path varchar(512)            not null comment '容器内挂载路径',
    read_only  tinyint(1)   default 1  null comment '是否只读'
);

create index idx_inference_mounts_deleted_at
    on inference_mounts (deleted_at);

create index idx_inference_mounts_service_id
    on inference_mounts (service_id);

INSERT INTO aiInfra.inference_mounts (id, created_at, updated_at, deleted_at, service_id, mount_type, pvc_id, pvc_name, sub_path, mount_path, read_only) VALUES (1, '2026-02-24 03:21:49.363', '2026-02-24 03:21:49.363', '2026-02-24 03:54:28.252', 1, '', 2, 'vol-1771561834-ufyl', '', '/model/Qwen2.5-0.5B-Instruct', 1);

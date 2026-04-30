create table inference_envs
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)     null,
    updated_at datetime(3)     null,
    deleted_at datetime(3)     null,
    service_id bigint unsigned not null comment '关联的推理服务ID',
    name       varchar(255)    not null comment '环境变量名',
    value      text            null comment '环境变量值'
);

create index idx_inference_envs_deleted_at
    on inference_envs (deleted_at);

create index idx_inference_envs_service_id
    on inference_envs (service_id);

INSERT INTO aiInfra.inference_envs (id, created_at, updated_at, deleted_at, service_id, name, value) VALUES (1, '2026-02-24 03:21:49.372', '2026-02-24 03:21:49.372', '2026-02-24 03:54:28.256', 1, 'VLLM_TARGET_DEVICE', 'cpu');

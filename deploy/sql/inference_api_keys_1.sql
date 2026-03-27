create table inference_api_keys
(
    id           bigint unsigned auto_increment
        primary key,
    created_at   datetime(3)                       null,
    updated_at   datetime(3)                       null,
    deleted_at   datetime(3)                       null,
    service_id   bigint unsigned                   not null comment '关联的推理服务ID',
    api_key      varchar(64)                       not null comment 'API Key',
    name         varchar(100)                      null comment 'Key名称',
    description  varchar(500)                      null comment '描述',
    status       varchar(20)  default 'active'     null comment '状态(active/disabled)',
    scopes       varchar(100) default 'read,write' null comment '权限范围(read,write)',
    rate_limit   bigint       default 0            null comment '每分钟请求限制(0表示不限制)',
    last_used_at datetime(3)                       null comment '最后使用时间',
    expired_at   datetime(3)                       null comment '过期时间',
    user_id      bigint unsigned                   not null comment '创建者用户ID',
    constraint idx_inference_api_keys_api_key
        unique (api_key)
);

create index idx_inference_api_keys_deleted_at
    on inference_api_keys (deleted_at);

create index idx_inference_api_keys_service_id
    on inference_api_keys (service_id);

create index idx_inference_api_keys_user_id
    on inference_api_keys (user_id);

INSERT INTO aiInfra.inference_api_keys (id, created_at, updated_at, deleted_at, service_id, api_key, name, description, status, scopes, rate_limit, last_used_at, expired_at, user_id) VALUES (1, '2026-02-24 23:22:12.928', '2026-02-24 23:53:18.849', null, 33, 'sk-552d315096e53a161086e8306426c757deeb1e676798c212', 'inference', '', 'active', 'read,write', 0, '2026-02-24 23:53:18.849', null, 1);

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

INSERT INTO aiInfra.inference_api_keys (id, created_at, updated_at, deleted_at, service_id, api_key, name, description, status, scopes, rate_limit, last_used_at, expired_at, user_id) VALUES (1, '2026-02-24 23:22:12.928', '2026-02-24 23:53:18.849', '2026-04-16 11:49:49.176', 33, 'sk-552d315096e53a161086e8306426c757deeb1e676798c212', 'inference', '', 'active', 'read,write', 0, '2026-02-24 23:53:18.849', null, 1);
INSERT INTO aiInfra.inference_api_keys (id, created_at, updated_at, deleted_at, service_id, api_key, name, description, status, scopes, rate_limit, last_used_at, expired_at, user_id) VALUES (2, '2026-04-20 14:15:48.321', '2026-04-22 15:11:47.901', '2026-04-23 22:17:48.540', 49, 'sk-f6c6cfc193699e933b9c789b5af852bab7b7f093f2d5b36e', 'hi', '', 'active', 'read,write', 0, '2026-04-22 15:11:47.900', null, 1);
INSERT INTO aiInfra.inference_api_keys (id, created_at, updated_at, deleted_at, service_id, api_key, name, description, status, scopes, rate_limit, last_used_at, expired_at, user_id) VALUES (3, '2026-04-22 16:11:56.340', '2026-04-22 16:13:29.821', null, 50, 'sk-96f06546e8c58d293d4638d6a07234b3a17f31a47e103961', '2', '', 'active', 'read,write', 0, '2026-04-22 16:13:29.821', null, 1);
INSERT INTO aiInfra.inference_api_keys (id, created_at, updated_at, deleted_at, service_id, api_key, name, description, status, scopes, rate_limit, last_used_at, expired_at, user_id) VALUES (4, '2026-04-22 21:05:55.936', '2026-04-22 21:18:51.364', '2026-04-23 22:17:27.288', 51, 'sk-734384f5ec6768eaee7aecdcfcb8e843ca9fbd3a92de39ad', '111', '', 'active', 'read,write', 0, '2026-04-22 21:18:51.364', null, 1);
INSERT INTO aiInfra.inference_api_keys (id, created_at, updated_at, deleted_at, service_id, api_key, name, description, status, scopes, rate_limit, last_used_at, expired_at, user_id) VALUES (5, '2026-04-23 21:16:26.363', '2026-04-23 21:21:12.044', '2026-04-23 22:17:23.242', 65, 'sk-0f8375b960ad0b87c60ca4406ffdd726ebc70a215b7ec5ee', '111', '', 'active', 'read,write', 0, '2026-04-23 21:21:12.044', null, 1);
INSERT INTO aiInfra.inference_api_keys (id, created_at, updated_at, deleted_at, service_id, api_key, name, description, status, scopes, rate_limit, last_used_at, expired_at, user_id) VALUES (6, '2026-04-23 21:57:32.681', '2026-04-23 22:03:35.613', null, 66, 'sk-5f6f69a7348d200eba6bd22a58b0ee6b1daa5e850f95bf8b', '111', '', 'active', 'read,write', 0, '2026-04-23 22:03:35.612', null, 1);
INSERT INTO aiInfra.inference_api_keys (id, created_at, updated_at, deleted_at, service_id, api_key, name, description, status, scopes, rate_limit, last_used_at, expired_at, user_id) VALUES (7, '2026-04-23 23:07:20.271', '2026-04-23 23:07:20.271', '2026-04-23 23:07:56.676', 71, 'sk-8a108832a8bf03096af3455f5a14ef7ed3bc70bc6e607fd7', '1', '', 'active', 'read,write', 0, null, null, 1);
INSERT INTO aiInfra.inference_api_keys (id, created_at, updated_at, deleted_at, service_id, api_key, name, description, status, scopes, rate_limit, last_used_at, expired_at, user_id) VALUES (8, '2026-04-23 23:07:28.624', '2026-04-23 23:07:28.624', '2026-04-23 23:07:54.811', 71, 'sk-2c5e513360b651d5dbb0f79087a7cfcd04de50dbe6b6306d', '2', '', 'active', 'read,write', 0, null, null, 1);
INSERT INTO aiInfra.inference_api_keys (id, created_at, updated_at, deleted_at, service_id, api_key, name, description, status, scopes, rate_limit, last_used_at, expired_at, user_id) VALUES (9, '2026-04-23 23:09:20.229', '2026-04-23 23:11:00.308', null, 71, 'sk-15d6cfdd2b269e19e0f4348f2d77aa1fdcba8cbd65153b74', '1', '', 'active', 'read,write', 0, '2026-04-23 23:11:00.308', null, 1);
INSERT INTO aiInfra.inference_api_keys (id, created_at, updated_at, deleted_at, service_id, api_key, name, description, status, scopes, rate_limit, last_used_at, expired_at, user_id) VALUES (10, '2026-04-24 00:53:05.883', '2026-04-24 00:53:05.883', '2026-04-27 10:24:41.142', 73, 'sk-875b4b1d8d4f64fc6c38f66e2f0f0195afda63a79a059018', '666', '', 'active', 'read,write', 0, null, null, 1);
INSERT INTO aiInfra.inference_api_keys (id, created_at, updated_at, deleted_at, service_id, api_key, name, description, status, scopes, rate_limit, last_used_at, expired_at, user_id) VALUES (11, '2026-04-27 11:14:15.204', '2026-04-27 11:22:30.050', null, 75, 'sk-3ce056a2b7e8e40d77d70142ded7ddd3693d9b19ae8f38ab', 'hi', '', 'active', 'read,write', 0, '2026-04-27 11:22:30.049', null, 1);

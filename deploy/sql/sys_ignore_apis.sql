create table sys_ignore_apis
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)                 null,
    updated_at datetime(3)                 null,
    deleted_at datetime(3)                 null,
    path       varchar(191)                null comment 'api路径',
    method     varchar(191) default 'POST' null comment '方法'
);

create index idx_sys_ignore_apis_deleted_at
    on sys_ignore_apis (deleted_at);

INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (1, '2025-12-22 23:16:24.595', '2025-12-22 23:16:24.595', null, '/swagger/*any', 'GET');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (2, '2025-12-22 23:16:24.595', '2025-12-22 23:16:24.595', null, '/api/freshCasbin', 'GET');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (3, '2025-12-22 23:16:24.595', '2025-12-22 23:16:24.595', null, '/uploads/file/*filepath', 'GET');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (4, '2025-12-22 23:16:24.595', '2025-12-22 23:16:24.595', null, '/health', 'GET');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (5, '2025-12-22 23:16:24.595', '2025-12-22 23:16:24.595', null, '/uploads/file/*filepath', 'HEAD');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (6, '2025-12-22 23:16:24.595', '2025-12-22 23:16:24.595', null, '/autoCode/llmAuto', 'POST');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (7, '2025-12-22 23:16:24.595', '2025-12-22 23:16:24.595', null, '/system/reloadSystem', 'POST');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (8, '2025-12-22 23:16:24.595', '2025-12-22 23:16:24.595', null, '/base/login', 'POST');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (9, '2025-12-22 23:16:24.595', '2025-12-22 23:16:24.595', null, '/base/captcha', 'POST');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (10, '2025-12-22 23:16:24.595', '2025-12-22 23:16:24.595', null, '/init/initdb', 'POST');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (11, '2025-12-22 23:16:24.595', '2025-12-22 23:16:24.595', null, '/init/checkdb', 'POST');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (12, '2025-12-22 23:16:24.595', '2025-12-22 23:16:24.595', null, '/info/getInfoDataSource', 'GET');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (13, '2025-12-22 23:16:24.595', '2025-12-22 23:16:24.595', null, '/info/getInfoPublic', 'GET');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (14, '2026-02-20 00:59:48.494', '2026-02-20 00:59:48.494', null, '/api/v1/freshCasbin', 'GET');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (15, '2026-02-20 00:59:48.494', '2026-02-20 00:59:48.494', null, '/api/v1/autocode/llm/add', 'POST');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (16, '2026-02-20 00:59:48.494', '2026-02-20 00:59:48.494', null, '/api/v1/system/reload', 'POST');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (17, '2026-02-20 00:59:48.494', '2026-02-20 00:59:48.494', null, '/api/v1/base/login', 'POST');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (18, '2026-02-20 00:59:48.494', '2026-02-20 00:59:48.494', null, '/api/v1/base/captcha', 'POST');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (19, '2026-02-20 00:59:48.494', '2026-02-20 00:59:48.494', null, '/api/v1/init/add', 'POST');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (20, '2026-02-20 00:59:48.494', '2026-02-20 00:59:48.494', null, '/api/v1/init/check', 'POST');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (21, '2026-02-20 00:59:48.494', '2026-02-20 00:59:48.494', null, '/api/v1/info/get/datasource', 'GET');
INSERT INTO aiInfra.sys_ignore_apis (id, created_at, updated_at, deleted_at, path, method) VALUES (22, '2026-02-20 00:59:48.494', '2026-02-20 00:59:48.494', null, '/api/v1/info/get/public', 'GET');

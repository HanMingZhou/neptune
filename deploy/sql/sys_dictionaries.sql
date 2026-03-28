create table sys_dictionaries
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)     null,
    updated_at datetime(3)     null,
    deleted_at datetime(3)     null,
    name       varchar(191)    null comment '字典名（中）',
    type       varchar(191)    null comment '字典名（英）',
    status     tinyint(1)      null comment '状态',
    `desc`     varchar(191)    null comment '描述',
    parent_id  bigint unsigned null comment '父级字典ID'
);

create index idx_sys_dictionaries_deleted_at
    on sys_dictionaries (deleted_at);

INSERT INTO aiInfra.sys_dictionaries (id, created_at, updated_at, deleted_at, name, type, status, `desc`, parent_id) VALUES (1, '2025-12-22 23:16:24.727', '2026-03-28 19:39:41.876', null, '性别', 'gender', 1, '性别字典', null);
INSERT INTO aiInfra.sys_dictionaries (id, created_at, updated_at, deleted_at, name, type, status, `desc`, parent_id) VALUES (2, '2025-12-22 23:16:24.727', '2026-03-28 19:39:41.903', null, '数据库int类型', 'int', 1, 'int类型对应的数据库类型', null);
INSERT INTO aiInfra.sys_dictionaries (id, created_at, updated_at, deleted_at, name, type, status, `desc`, parent_id) VALUES (3, '2025-12-22 23:16:24.727', '2026-03-28 19:39:41.927', null, '数据库时间日期类型', 'time.Time', 1, '数据库时间日期类型', null);
INSERT INTO aiInfra.sys_dictionaries (id, created_at, updated_at, deleted_at, name, type, status, `desc`, parent_id) VALUES (4, '2025-12-22 23:16:24.727', '2026-03-28 19:39:41.949', null, '数据库浮点型', 'float64', 1, '数据库浮点型', null);
INSERT INTO aiInfra.sys_dictionaries (id, created_at, updated_at, deleted_at, name, type, status, `desc`, parent_id) VALUES (5, '2025-12-22 23:16:24.727', '2026-03-28 19:39:41.972', null, '数据库字符串', 'string', 1, '数据库字符串', null);
INSERT INTO aiInfra.sys_dictionaries (id, created_at, updated_at, deleted_at, name, type, status, `desc`, parent_id) VALUES (6, '2025-12-22 23:16:24.727', '2026-03-28 19:39:42.003', null, '数据库bool类型', 'bool', 1, '数据库bool类型', null);

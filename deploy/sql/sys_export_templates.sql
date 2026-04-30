create table sys_export_templates
(
    id            bigint unsigned auto_increment
        primary key,
    created_at    datetime(3)  null,
    updated_at    datetime(3)  null,
    deleted_at    datetime(3)  null,
    db_name       varchar(191) null comment '数据库名称',
    name          varchar(191) null comment '模板名称',
    table_name    varchar(191) null comment '表名称',
    template_id   varchar(191) null comment '模板标识',
    template_info text         null,
    `limit`       bigint       null comment '导出限制',
    `order`       varchar(191) null comment '排序'
);

create index idx_sys_export_templates_deleted_at
    on sys_export_templates (deleted_at);

INSERT INTO aiInfra.sys_export_templates (id, created_at, updated_at, deleted_at, db_name, name, table_name, template_id, template_info, `limit`, `order`) VALUES (1, '2025-12-22 23:16:25.312', '2025-12-22 23:16:25.312', null, '', 'api', 'sys_apis', 'api', '{
"path":"路径",
"method":"方法（大写）",
"description":"方法介绍",
"api_group":"方法分组"
}', null, '');

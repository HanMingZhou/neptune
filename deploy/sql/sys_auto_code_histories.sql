create table sys_auto_code_histories
(
    id                 bigint unsigned auto_increment
        primary key,
    created_at         datetime(3)     null,
    updated_at         datetime(3)     null,
    deleted_at         datetime(3)     null,
    table_name         varchar(191)    null comment '表名',
    package            varchar(191)    null comment '模块名/插件名',
    request            text            null comment '前端传入的结构化信息',
    struct_name        varchar(191)    null comment '结构体名称',
    abbreviation       varchar(191)    null comment '结构体名称缩写',
    business_db        varchar(191)    null comment '业务库',
    description        varchar(191)    null comment 'Struct中文名称',
    templates          text            null comment '模板信息',
    Injections         text            null comment '注入路径',
    flag               bigint          null comment '[0:创建,1:回滚]',
    api_ids            varchar(191)    null comment 'api表注册内容',
    menu_id            bigint unsigned null comment '菜单ID',
    export_template_id bigint unsigned null comment '导出模板ID',
    package_id         bigint unsigned null comment '包ID'
);

create index idx_sys_auto_code_histories_deleted_at
    on sys_auto_code_histories (deleted_at);


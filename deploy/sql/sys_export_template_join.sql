create table sys_export_template_join
(
    id          bigint unsigned auto_increment
        primary key,
    created_at  datetime(3)  null,
    updated_at  datetime(3)  null,
    deleted_at  datetime(3)  null,
    template_id varchar(191) null comment '模板标识',
    joins       varchar(191) null comment '关联',
    `table`     varchar(191) null comment '关联表',
    `on`        varchar(191) null comment '关联条件'
);

create index idx_sys_export_template_join_deleted_at
    on sys_export_template_join (deleted_at);


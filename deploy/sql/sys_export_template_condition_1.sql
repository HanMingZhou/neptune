create table sys_export_template_condition
(
    id          bigint unsigned auto_increment
        primary key,
    created_at  datetime(3)  null,
    updated_at  datetime(3)  null,
    deleted_at  datetime(3)  null,
    template_id varchar(191) null comment '模板标识',
    `from`      varchar(191) null comment '条件取的key',
    `column`    varchar(191) null comment '作为查询条件的字段',
    operator    varchar(191) null comment '操作符'
);

create index idx_sys_export_template_condition_deleted_at
    on sys_export_template_condition (deleted_at);


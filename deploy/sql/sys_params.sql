create table sys_params
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)  null,
    updated_at datetime(3)  null,
    deleted_at datetime(3)  null,
    name       varchar(191) null comment '参数名称',
    `key`      varchar(191) null comment '参数键',
    value      varchar(191) null comment '参数值',
    `desc`     varchar(191) null comment '参数说明'
);

create index idx_sys_params_deleted_at
    on sys_params (deleted_at);


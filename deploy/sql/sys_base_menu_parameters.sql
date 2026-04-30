create table sys_base_menu_parameters
(
    id               bigint unsigned auto_increment
        primary key,
    created_at       datetime(3)     null,
    updated_at       datetime(3)     null,
    deleted_at       datetime(3)     null,
    sys_base_menu_id bigint unsigned null,
    type             varchar(191)    null comment '地址栏携带参数为params还是query',
    `key`            varchar(191)    null comment '地址栏携带参数的key',
    value            varchar(191)    null comment '地址栏携带参数的值'
);

create index idx_sys_base_menu_parameters_deleted_at
    on sys_base_menu_parameters (deleted_at);


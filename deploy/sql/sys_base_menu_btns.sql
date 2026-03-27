create table sys_base_menu_btns
(
    id               bigint unsigned auto_increment
        primary key,
    created_at       datetime(3)     null,
    updated_at       datetime(3)     null,
    deleted_at       datetime(3)     null,
    name             varchar(191)    null comment '按钮关键key',
    `desc`           varchar(191)    null,
    sys_base_menu_id bigint unsigned null comment '菜单ID'
);

create index idx_sys_base_menu_btns_deleted_at
    on sys_base_menu_btns (deleted_at);


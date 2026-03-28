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

INSERT INTO aiInfra.sys_base_menu_btns (id, created_at, updated_at, deleted_at, name, `desc`, sys_base_menu_id) VALUES (1, '2026-03-28 19:38:13.307', '2026-03-28 19:38:13.307', null, 'recharge_system', '平台代充值', 268);
INSERT INTO aiInfra.sys_base_menu_btns (id, created_at, updated_at, deleted_at, name, `desc`, sys_base_menu_id) VALUES (2, '2026-03-28 19:38:13.307', '2026-03-28 19:38:13.307', null, 'recharge_alipay', '支付宝充值', 268);
INSERT INTO aiInfra.sys_base_menu_btns (id, created_at, updated_at, deleted_at, name, `desc`, sys_base_menu_id) VALUES (3, '2026-03-28 19:38:13.307', '2026-03-28 19:38:13.307', null, 'recharge_wechat', '微信充值', 268);

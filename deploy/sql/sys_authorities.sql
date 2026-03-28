create table sys_authorities
(
    created_at     datetime(3)                      null,
    updated_at     datetime(3)                      null,
    deleted_at     datetime(3)                      null,
    authority_id   bigint unsigned auto_increment comment '角色ID'
        primary key,
    authority_name varchar(191)                     null comment '角色名',
    parent_id      bigint unsigned                  null comment '父角色ID',
    default_router varchar(191) default 'dashboard' null comment '默认菜单',
    constraint uni_sys_authorities_authority_id
        unique (authority_id)
);

INSERT INTO aiInfra.sys_authorities (created_at, updated_at, deleted_at, authority_id, authority_name, parent_id, default_router) VALUES ('2025-12-22 23:16:24.645', '2026-03-28 19:39:42.193', null, 888, '管理员', 0, 'dashboard');
INSERT INTO aiInfra.sys_authorities (created_at, updated_at, deleted_at, authority_id, authority_name, parent_id, default_router) VALUES ('2025-12-22 23:16:24.645', '2026-03-28 19:39:42.212', null, 8881, '普通用户', 888, 'dashboard');
INSERT INTO aiInfra.sys_authorities (created_at, updated_at, deleted_at, authority_id, authority_name, parent_id, default_router) VALUES ('2025-12-22 23:16:24.645', '2026-03-28 19:39:42.202', null, 9528, '测试角色', 0, 'dashboard');

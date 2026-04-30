create table sys_users
(
    id             bigint unsigned auto_increment
        primary key,
    created_at     datetime(3)                                                              null,
    updated_at     datetime(3)                                                              null,
    deleted_at     datetime(3)                                                              null,
    uuid           varchar(191)                                                             null comment '用户UUID',
    username       varchar(191)                                                             null comment '用户登录名',
    password       varchar(191)                                                             null comment '用户登录密码',
    nick_name      varchar(191)    default '系统用户'                                       null comment '用户昵称',
    header_img     varchar(191)    default 'https://qmplusimg.henrongyi.top/gva_header.jpg' null comment '用户头像',
    authority_id   bigint unsigned default '888'                                            null comment '用户角色ID',
    phone          varchar(191)                                                             null comment '用户手机号',
    email          varchar(191)                                                             null comment '用户邮箱',
    enable         bigint          default 1                                                null comment '用户是否被冻结 1正常 2冻结',
    origin_setting text                                                                     null comment '配置',
    namespace      varchar(191)                                                             null comment '用户命名空间'
);

create index idx_sys_users_deleted_at
    on sys_users (deleted_at);

create index idx_sys_users_username
    on sys_users (username);

create index idx_sys_users_uuid
    on sys_users (uuid);

INSERT INTO aiInfra.sys_users (id, created_at, updated_at, deleted_at, uuid, username, password, nick_name, header_img, authority_id, phone, email, enable, origin_setting, namespace) VALUES (1, '2025-12-22 23:16:25.260', '2026-04-15 15:11:43.266', null, 'c3e9f113-f2f8-4856-8fc1-080860514559', 'admin', '$2a$10$uurhikLo8KvtQLll5790zuprD0cGwuU6uhPIhs8XBiO93tkncOW1i', 'Mr.奇淼', 'https://qmplusimg.henrongyi.top/gva_header.jpg', 888, '18888888888', '333333333@qq.com', 1, null, 'zzz');
INSERT INTO aiInfra.sys_users (id, created_at, updated_at, deleted_at, uuid, username, password, nick_name, header_img, authority_id, phone, email, enable, origin_setting, namespace) VALUES (2, '2025-12-22 23:16:25.260', '2026-04-14 21:47:55.617', null, '35dfd37f-f0e8-4ff3-bcb3-523866a25d89', 'a303176530', '$2a$10$DgEMCdIXOfpneUNQAnKq..g6rGmMCG8mqXjdW0WhsTrZIqNkufX7K', '用户1', 'https://qmplusimg.henrongyi.top/1572075907logo.png', 9528, '17611111111', '333333333@qq.com', 1, null, 'hmz');
INSERT INTO aiInfra.sys_users (id, created_at, updated_at, deleted_at, uuid, username, password, nick_name, header_img, authority_id, phone, email, enable, origin_setting, namespace) VALUES (4, '2026-02-25 01:09:53.003', '2026-02-25 01:09:53.003', null, '66fca86e-e6b6-4bd0-8fbb-0f1f2a7ab6e4', 'maqiliang', '$2a$10$7Pncjq1Ygpci6WwCYoe00eCMxr1xnDt3ImDnhExCggRKe4c.ddwzO', 'maqiliang', 'https://qmplusimg.henrongyi.top/gva_header.jpg', 8881, '', 'maqiliang@gmail.com', 1, null, '2c8c2a');
INSERT INTO aiInfra.sys_users (id, created_at, updated_at, deleted_at, uuid, username, password, nick_name, header_img, authority_id, phone, email, enable, origin_setting, namespace) VALUES (5, '2026-04-16 21:43:13.700', '2026-04-16 21:43:13.700', null, 'bb00a55a-b038-4e62-9bff-3a9f9d48b53d', '66666', '$2a$10$vXwZOHOyEWsTpZ0iD0XCqOhx6VUdrFyb5xkmVWR9GEb8tIW7TgvzC', '66666', 'https://qmplusimg.henrongyi.top/gva_header.jpg', 8881, '', '66666@gmail.com', 1, null, 'd145a8');
INSERT INTO aiInfra.sys_users (id, created_at, updated_at, deleted_at, uuid, username, password, nick_name, header_img, authority_id, phone, email, enable, origin_setting, namespace) VALUES (6, '2026-04-30 10:38:57.859', '2026-04-30 10:38:57.859', null, '853113e1-7b74-4d45-b847-db5beaae3d51', 'jackma', '$2a$10$vbINhWke9ymjJ4gHQUJRWuWzOGgf7uYrf3BbE0lQoywpkU.7n.fIm', 'jackma', 'https://qmplusimg.henrongyi.top/gva_header.jpg', 8881, '', 'jackma@qq.com', 1, null, 'ceb852');

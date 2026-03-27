create table sys_account_securities
(
    id                bigint unsigned auto_increment
        primary key,
    created_at        datetime(3)          null,
    updated_at        datetime(3)          null,
    deleted_at        datetime(3)          null,
    user_id           bigint unsigned      null comment '用户ID',
    mfa_enabled       tinyint(1) default 0 null comment '是否开启MFA',
    mfa_secret        varchar(191)         null comment 'MFA密钥',
    github_id         varchar(191)         null comment 'Github ID',
    github_username   varchar(191)         null comment 'Github用户名',
    access_key_id     varchar(191)         null comment 'AccessKey ID',
    access_key_secret varchar(191)         null comment 'AccessKey Secret',
    constraint idx_sys_account_securities_access_key_id
        unique (access_key_id),
    constraint idx_sys_account_securities_user_id
        unique (user_id)
);

create index idx_sys_account_securities_deleted_at
    on sys_account_securities (deleted_at);

create index idx_sys_account_securities_github_id
    on sys_account_securities (github_id);

INSERT INTO aiInfra.sys_account_securities (id, created_at, updated_at, deleted_at, user_id, mfa_enabled, mfa_secret, github_id, github_username, access_key_id, access_key_secret) VALUES (1, '2026-02-16 11:57:50.263', '2026-02-20 14:49:35.130', null, 1, 0, 'YNke2CARJYpgE6zk9y2Pntb1F0S9xaWc', '', '', '3fyJL1PIIPU3UFXo', 'fFA20jpMEITEIGoRKm423c5MHiP5i4KL');

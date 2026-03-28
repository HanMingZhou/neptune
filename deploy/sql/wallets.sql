create table wallets
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)                     null,
    updated_at datetime(3)                     null,
    deleted_at datetime(3)                     null,
    user_id    bigint unsigned                 null comment '用户ID',
    balance    decimal(20, 6) default 0.000000 null comment '可用余额',
    frozen     decimal(20, 6) default 0.000000 null comment '锁定余额',
    version    bigint         default 0        null comment '乐观锁版本号',
    constraint idx_wallets_user_id
        unique (user_id)
);

create index idx_wallets_deleted_at
    on wallets (deleted_at);

INSERT INTO aiInfra.wallets (id, created_at, updated_at, deleted_at, user_id, balance, frozen, version) VALUES (1, '2026-02-23 12:27:44.944', '2026-03-28 19:45:10.428', null, 1, 1000000558.310969, 0.000000, 59);
INSERT INTO aiInfra.wallets (id, created_at, updated_at, deleted_at, user_id, balance, frozen, version) VALUES (2, '2026-02-25 00:48:42.732', '2026-02-25 00:48:42.735', null, 2, 1000000.000000, 0.000000, 0);
INSERT INTO aiInfra.wallets (id, created_at, updated_at, deleted_at, user_id, balance, frozen, version) VALUES (3, '2026-02-25 01:09:53.027', '2026-02-25 01:09:53.027', null, 4, 0.000000, 0.000000, 0);

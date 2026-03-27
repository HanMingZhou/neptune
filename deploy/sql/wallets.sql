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

INSERT INTO aiInfra.wallets (id, created_at, updated_at, deleted_at, user_id, balance, frozen, version) VALUES (1, '2026-02-23 12:27:44.944', '2026-02-25 00:07:12.049', null, 1, 999999958.697080, 0.000000, 50);

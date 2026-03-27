create table order_invoice_addresses
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)     null,
    updated_at datetime(3)     null,
    deleted_at datetime(3)     null,
    user_id    bigint unsigned null comment '用户ID',
    consignee  varchar(50)     null comment '收货人',
    phone      varchar(20)     null comment '联系电话',
    province   varchar(50)     null comment '省',
    city       varchar(50)     null comment '市',
    district   varchar(50)     null comment '区',
    detail     varchar(200)    null comment '详细地址'
);

create index idx_order_invoice_addresses_deleted_at
    on order_invoice_addresses (deleted_at);

create index idx_order_invoice_addresses_user_id
    on order_invoice_addresses (user_id);


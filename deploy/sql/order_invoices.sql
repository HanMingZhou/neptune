create table order_invoices
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)     null,
    updated_at datetime(3)     null,
    deleted_at datetime(3)     null,
    request_id varchar(64)     null comment '申请单号',
    user_id    bigint unsigned null comment '用户ID',
    amount     decimal(20, 6)  null comment '发票金额',
    title      varchar(200)    null comment '发票抬头',
    type       varchar(50)     null comment '发票类型',
    status     varchar(20)     null comment '状态(Processing/Sent)',
    code       varchar(50)     null comment '发票代码(电子发票)',
    number     varchar(50)     null comment '发票号码',
    file_url   varchar(500)    null comment '电子发票下载地址',
    constraint idx_order_invoices_request_id
        unique (request_id)
);

create index idx_order_invoices_deleted_at
    on order_invoices (deleted_at);

create index idx_order_invoices_user_id
    on order_invoices (user_id);


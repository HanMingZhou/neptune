create table order_invoice_titles
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)     null,
    updated_at datetime(3)     null,
    deleted_at datetime(3)     null,
    user_id    bigint unsigned null comment '用户ID',
    title      varchar(200)    null comment '抬头名称',
    tax_no     varchar(50)     null comment '纳税人识别号',
    type       varchar(20)     null comment '类型(Personal/Enterprise)'
);

create index idx_order_invoice_titles_deleted_at
    on order_invoice_titles (deleted_at);

create index idx_order_invoice_titles_user_id
    on order_invoice_titles (user_id);


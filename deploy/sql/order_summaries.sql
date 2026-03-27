create table order_summaries
(
    id            bigint unsigned auto_increment
        primary key,
    created_at    datetime(3)                     null,
    updated_at    datetime(3)                     null,
    deleted_at    datetime(3)                     null,
    user_id       bigint unsigned                 null comment '用户ID',
    order_date    varchar(10)                     null comment '账单日期(YYYY-MM-DD)',
    order_month   varchar(7)                      null comment '账单月份(YYYY-MM)',
    total_amount  decimal(20, 6) default 0.000000 null comment '总金额',
    paid_amount   decimal(20, 6) default 0.000000 null comment '已支付金额',
    unpaid_amount decimal(20, 6) default 0.000000 null comment '未支付金额',
    refund_amount decimal(20, 6) default 0.000000 null comment '退款金额'
);

create index idx_order_summaries_deleted_at
    on order_summaries (deleted_at);

create index idx_order_summaries_order_date
    on order_summaries (order_date);

create index idx_order_summaries_order_month
    on order_summaries (order_month);

create index idx_order_summaries_user_id
    on order_summaries (user_id);


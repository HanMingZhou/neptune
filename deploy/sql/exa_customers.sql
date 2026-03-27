create table exa_customers
(
    id                    bigint unsigned auto_increment
        primary key,
    created_at            datetime(3)     null,
    updated_at            datetime(3)     null,
    deleted_at            datetime(3)     null,
    customer_name         varchar(191)    null comment '客户名',
    customer_phone_data   varchar(191)    null comment '客户手机号',
    sys_user_id           bigint unsigned null comment '管理ID',
    sys_user_authority_id bigint unsigned null comment '管理角色ID'
);

create index idx_exa_customers_deleted_at
    on exa_customers (deleted_at);


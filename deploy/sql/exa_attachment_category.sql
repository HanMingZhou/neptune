create table exa_attachment_category
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)      null,
    updated_at datetime(3)      null,
    deleted_at datetime(3)      null,
    name       varchar(255)     null comment '分类名称',
    pid        bigint default 0 null comment '父节点ID'
);

create index idx_exa_attachment_category_deleted_at
    on exa_attachment_category (deleted_at);


create table sys_versions
(
    id           bigint unsigned auto_increment
        primary key,
    created_at   datetime(3)  null,
    updated_at   datetime(3)  null,
    deleted_at   datetime(3)  null,
    version_name varchar(255) null comment '版本名称',
    version_code varchar(100) null comment '版本号',
    description  varchar(500) null comment '版本描述',
    version_data text         null comment '版本数据JSON'
);

create index idx_sys_versions_deleted_at
    on sys_versions (deleted_at);


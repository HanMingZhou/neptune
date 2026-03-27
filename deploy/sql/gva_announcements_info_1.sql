create table gva_announcements_info
(
    id          bigint unsigned auto_increment
        primary key,
    created_at  datetime(3)  null,
    updated_at  datetime(3)  null,
    deleted_at  datetime(3)  null,
    title       varchar(191) null comment '公告标题',
    content     text         null comment '公告内容',
    user_id     bigint       null comment '发布者',
    attachments json         null comment '相关附件'
);

create index idx_gva_announcements_info_deleted_at
    on gva_announcements_info (deleted_at);


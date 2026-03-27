create table exa_files
(
    id          bigint unsigned auto_increment
        primary key,
    created_at  datetime(3)  null,
    updated_at  datetime(3)  null,
    deleted_at  datetime(3)  null,
    file_name   varchar(191) null,
    file_md5    varchar(191) null,
    file_path   varchar(191) null,
    chunk_total bigint       null,
    is_finish   tinyint(1)   null
);

create index idx_exa_files_deleted_at
    on exa_files (deleted_at);


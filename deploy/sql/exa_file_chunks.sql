create table exa_file_chunks
(
    id                bigint unsigned auto_increment
        primary key,
    created_at        datetime(3)     null,
    updated_at        datetime(3)     null,
    deleted_at        datetime(3)     null,
    exa_file_id       bigint unsigned null,
    file_chunk_number bigint          null,
    file_chunk_path   varchar(191)    null
);

create index idx_exa_file_chunks_deleted_at
    on exa_file_chunks (deleted_at);


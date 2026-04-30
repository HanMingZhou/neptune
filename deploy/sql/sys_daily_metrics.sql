create table sys_daily_metrics
(
    id            bigint unsigned auto_increment
        primary key,
    created_at    datetime(3)  null,
    updated_at    datetime(3)  null,
    deleted_at    datetime(3)  null,
    date          varchar(191) null comment '日期(YYYY-MM-DD)',
    gpu_usage     double       null comment 'GPU使用率',
    running_tasks bigint       null comment '运行中任务数',
    storage_usage double       null comment '存储使用量(GB)',
    total_cost    double       null comment '总花费'
);

create index idx_sys_daily_metrics_date
    on sys_daily_metrics (date);

create index idx_sys_daily_metrics_deleted_at
    on sys_daily_metrics (deleted_at);


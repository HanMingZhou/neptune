create table training_job_envs
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)  null,
    updated_at datetime(3)  null,
    deleted_at datetime(3)  null,
    job_id     bigint       not null comment '关联的训练任务ID',
    name       varchar(255) not null comment '环境变量名',
    value      text         null comment '环境变量值'
);

create index idx_training_job_envs_deleted_at
    on training_job_envs (deleted_at);

create index idx_training_job_envs_job_id
    on training_job_envs (job_id);

INSERT INTO aiInfra.training_job_envs (id, created_at, updated_at, deleted_at, job_id, name, value) VALUES (1, '2026-02-21 22:22:47.127', '2026-02-21 22:22:47.127', null, 56, '--log-dir', '/mpi/logs');
INSERT INTO aiInfra.training_job_envs (id, created_at, updated_at, deleted_at, job_id, name, value) VALUES (2, '2026-04-03 21:41:34.731', '2026-04-03 21:41:34.731', null, 85, '--log-dir', ' /home/notebook/neptune/logs');
INSERT INTO aiInfra.training_job_envs (id, created_at, updated_at, deleted_at, job_id, name, value) VALUES (3, '2026-04-03 21:44:21.242', '2026-04-03 21:44:21.242', null, 86, '-log-dir', '/train/logs');
INSERT INTO aiInfra.training_job_envs (id, created_at, updated_at, deleted_at, job_id, name, value) VALUES (4, '2026-04-03 21:46:27.231', '2026-04-03 21:46:27.231', null, 87, '--log-dir', '/train/logs');

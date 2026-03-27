create table sys_account_access_logs
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)     null,
    updated_at datetime(3)     null,
    deleted_at datetime(3)     null,
    user_id    bigint unsigned null comment '用户ID',
    ip         varchar(191)    null comment '登录IP',
    location   varchar(191)    null comment 'IP归属地',
    device     varchar(191)    null comment '登录设备',
    browser    varchar(191)    null comment '浏览器',
    os         varchar(191)    null comment '操作系统',
    method     varchar(191)    null comment '登录方式(密码/Github/Wechat等)',
    status     varchar(191)    null comment '状态',
    reason     varchar(191)    null comment '失败原因',
    log_type   varchar(191)    null comment '日志类型',
    login_time datetime(3)     null comment '登录时间'
);

create index idx_sys_account_access_logs_deleted_at
    on sys_account_access_logs (deleted_at);

create index idx_sys_account_access_logs_user_id
    on sys_account_access_logs (user_id);

INSERT INTO aiInfra.sys_account_access_logs (id, created_at, updated_at, deleted_at, user_id, ip, location, device, browser, os, method, status, reason, log_type, login_time) VALUES (1, '2026-02-20 02:05:46.682', '2026-02-20 02:05:46.682', null, 1, '127.0.0.1', 'Unknown', 'Chrome / Mac OS', 'Chrome', 'Mac OS', 'Password', 'Success', '', 'Login', '2026-02-20 02:05:46.680');
INSERT INTO aiInfra.sys_account_access_logs (id, created_at, updated_at, deleted_at, user_id, ip, location, device, browser, os, method, status, reason, log_type, login_time) VALUES (2, '2026-02-20 02:20:05.642', '2026-02-20 02:20:05.642', null, 1, '127.0.0.1', 'Unknown', 'Chrome / Mac OS', 'Chrome', 'Mac OS', 'Password', 'Success', '', 'Login', '2026-02-20 02:20:05.638');
INSERT INTO aiInfra.sys_account_access_logs (id, created_at, updated_at, deleted_at, user_id, ip, location, device, browser, os, method, status, reason, log_type, login_time) VALUES (3, '2026-02-20 10:42:57.106', '2026-02-20 10:42:57.106', null, 1, '127.0.0.1', 'Unknown', 'Chrome / Mac OS', 'Chrome', 'Mac OS', 'Password', 'Success', '', 'Login', '2026-02-20 10:42:57.105');
INSERT INTO aiInfra.sys_account_access_logs (id, created_at, updated_at, deleted_at, user_id, ip, location, device, browser, os, method, status, reason, log_type, login_time) VALUES (4, '2026-02-20 10:58:49.617', '2026-02-20 10:58:49.617', null, 1, '127.0.0.1', 'Unknown', 'Chrome / Mac OS', 'Chrome', 'Mac OS', 'Password', 'Success', '', 'Login', '2026-02-20 10:58:49.617');
INSERT INTO aiInfra.sys_account_access_logs (id, created_at, updated_at, deleted_at, user_id, ip, location, device, browser, os, method, status, reason, log_type, login_time) VALUES (5, '2026-02-20 12:30:25.869', '2026-02-20 12:30:25.869', null, 1, '127.0.0.1', 'Unknown', 'Chrome / Mac OS', 'Chrome', 'Mac OS', 'Password', 'Success', '', 'Login', '2026-02-20 12:30:25.868');
INSERT INTO aiInfra.sys_account_access_logs (id, created_at, updated_at, deleted_at, user_id, ip, location, device, browser, os, method, status, reason, log_type, login_time) VALUES (6, '2026-02-20 13:01:33.739', '2026-02-20 13:01:33.739', null, 1, '127.0.0.1', 'Unknown', 'Chrome / Mac OS', 'Chrome', 'Mac OS', 'Password', 'Success', '', 'Login', '2026-02-20 13:01:33.737');
INSERT INTO aiInfra.sys_account_access_logs (id, created_at, updated_at, deleted_at, user_id, ip, location, device, browser, os, method, status, reason, log_type, login_time) VALUES (7, '2026-02-20 14:18:49.670', '2026-02-20 14:18:49.670', null, 1, '127.0.0.1', 'Unknown', 'Chrome / Mac OS', 'Chrome', 'Mac OS', 'Password', 'Success', '', 'Login', '2026-02-20 14:18:49.666');
INSERT INTO aiInfra.sys_account_access_logs (id, created_at, updated_at, deleted_at, user_id, ip, location, device, browser, os, method, status, reason, log_type, login_time) VALUES (8, '2026-02-20 15:53:13.807', '2026-02-20 15:53:13.807', null, 1, '::1', 'Unknown', 'Chrome / Mac OS', 'Chrome', 'Mac OS', 'Password', 'Success', '', 'Login', '2026-02-20 15:53:13.805');
INSERT INTO aiInfra.sys_account_access_logs (id, created_at, updated_at, deleted_at, user_id, ip, location, device, browser, os, method, status, reason, log_type, login_time) VALUES (9, '2026-02-21 00:03:18.588', '2026-02-21 00:03:18.588', null, 1, '::1', 'Unknown', 'Chrome / Mac OS', 'Chrome', 'Mac OS', 'Password', 'Success', '', 'Login', '2026-02-21 00:03:18.587');
INSERT INTO aiInfra.sys_account_access_logs (id, created_at, updated_at, deleted_at, user_id, ip, location, device, browser, os, method, status, reason, log_type, login_time) VALUES (10, '2026-02-21 16:31:20.053', '2026-02-21 16:31:20.053', null, 1, '::1', 'Unknown', 'Chrome / Mac OS', 'Chrome', 'Mac OS', 'Password', 'Success', '', 'Login', '2026-02-21 16:31:20.048');
INSERT INTO aiInfra.sys_account_access_logs (id, created_at, updated_at, deleted_at, user_id, ip, location, device, browser, os, method, status, reason, log_type, login_time) VALUES (11, '2026-02-22 10:15:23.941', '2026-02-22 10:15:23.941', null, 1, '::1', 'Unknown', 'Chrome / Mac OS', 'Chrome', 'Mac OS', 'Password', 'Success', '', 'Login', '2026-02-22 10:15:23.939');
INSERT INTO aiInfra.sys_account_access_logs (id, created_at, updated_at, deleted_at, user_id, ip, location, device, browser, os, method, status, reason, log_type, login_time) VALUES (12, '2026-02-23 11:13:20.450', '2026-02-23 11:13:20.450', null, 1, '::1', 'Unknown', 'Chrome / Mac OS', 'Chrome', 'Mac OS', 'Password', 'Success', '', 'Login', '2026-02-23 11:13:20.448');
INSERT INTO aiInfra.sys_account_access_logs (id, created_at, updated_at, deleted_at, user_id, ip, location, device, browser, os, method, status, reason, log_type, login_time) VALUES (13, '2026-02-23 12:27:17.733', '2026-02-23 12:27:17.733', null, 1, '::1', 'Unknown', 'Chrome / Mac OS', 'Chrome', 'Mac OS', 'Password', 'Success', '', 'Login', '2026-02-23 12:27:17.731');
INSERT INTO aiInfra.sys_account_access_logs (id, created_at, updated_at, deleted_at, user_id, ip, location, device, browser, os, method, status, reason, log_type, login_time) VALUES (14, '2026-02-23 21:06:57.379', '2026-02-23 21:06:57.379', null, 1, '::1', 'Unknown', 'Chrome / Mac OS', 'Chrome', 'Mac OS', 'Password', 'Success', '', 'Login', '2026-02-23 21:06:57.379');
INSERT INTO aiInfra.sys_account_access_logs (id, created_at, updated_at, deleted_at, user_id, ip, location, device, browser, os, method, status, reason, log_type, login_time) VALUES (15, '2026-02-24 00:58:19.487', '2026-02-24 00:58:19.487', null, 1, '::1', 'Unknown', 'Chrome / Mac OS', 'Chrome', 'Mac OS', 'Password', 'Success', '', 'Login', '2026-02-24 00:58:19.485');
INSERT INTO aiInfra.sys_account_access_logs (id, created_at, updated_at, deleted_at, user_id, ip, location, device, browser, os, method, status, reason, log_type, login_time) VALUES (16, '2026-02-24 10:40:27.402', '2026-02-24 10:40:27.402', null, 1, '::1', 'Unknown', 'Chrome / Mac OS', 'Chrome', 'Mac OS', 'Password', 'Success', '', 'Login', '2026-02-24 10:40:27.401');
INSERT INTO aiInfra.sys_account_access_logs (id, created_at, updated_at, deleted_at, user_id, ip, location, device, browser, os, method, status, reason, log_type, login_time) VALUES (17, '2026-02-25 00:05:08.455', '2026-02-25 00:05:08.455', null, 1, '::1', 'Unknown', 'Chrome / Mac OS', 'Chrome', 'Mac OS', 'Password', 'Success', '', 'Login', '2026-02-25 00:05:08.454');
INSERT INTO aiInfra.sys_account_access_logs (id, created_at, updated_at, deleted_at, user_id, ip, location, device, browser, os, method, status, reason, log_type, login_time) VALUES (18, '2026-02-25 00:15:18.886', '2026-02-25 00:15:18.886', null, 1, '::1', 'Unknown', 'Chrome / Mac OS', 'Chrome', 'Mac OS', 'Password', 'Success', '', 'Login', '2026-02-25 00:15:18.884');

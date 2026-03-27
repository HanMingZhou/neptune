create table sys_auto_code_packages
(
    id           bigint unsigned auto_increment
        primary key,
    created_at   datetime(3)  null,
    updated_at   datetime(3)  null,
    deleted_at   datetime(3)  null,
    `desc`       varchar(191) null comment '描述',
    label        varchar(191) null comment '展示名',
    template     varchar(191) null comment '模版',
    package_name varchar(191) null comment '包名',
    module       varchar(191) null
);

create index idx_sys_auto_code_packages_deleted_at
    on sys_auto_code_packages (deleted_at);

INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (1, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取apisix包', 'apisix包', 'package', 'apisix', 'github.com/flipped-aurora/gin-vue-admin/server');
INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (2, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取billing包', 'billing包', 'package', 'billing', 'github.com/flipped-aurora/gin-vue-admin/server');
INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (3, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取cms包', 'cms包', 'package', 'cms', 'github.com/flipped-aurora/gin-vue-admin/server');
INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (4, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取example包', 'example包', 'package', 'example', 'github.com/flipped-aurora/gin-vue-admin/server');
INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (5, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取image包', 'image包', 'package', 'image', 'github.com/flipped-aurora/gin-vue-admin/server');
INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (6, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取inference包', 'inference包', 'package', 'inference', 'github.com/flipped-aurora/gin-vue-admin/server');
INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (7, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取notebook包', 'notebook包', 'package', 'notebook', 'github.com/flipped-aurora/gin-vue-admin/server');
INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (8, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取piper包', 'piper包', 'package', 'piper', 'github.com/flipped-aurora/gin-vue-admin/server');
INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (9, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取podgroup包', 'podgroup包', 'package', 'podgroup', 'github.com/flipped-aurora/gin-vue-admin/server');
INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (10, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取product包', 'product包', 'package', 'product', 'github.com/flipped-aurora/gin-vue-admin/server');
INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (11, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取pvc包', 'pvc包', 'package', 'pvc', 'github.com/flipped-aurora/gin-vue-admin/server');
INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (12, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取secret包', 'secret包', 'package', 'secret', 'github.com/flipped-aurora/gin-vue-admin/server');
INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (13, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取sshkey包', 'sshkey包', 'package', 'sshkey', 'github.com/flipped-aurora/gin-vue-admin/server');
INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (14, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取system包', 'system包', 'package', 'system', 'github.com/flipped-aurora/gin-vue-admin/server');
INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (15, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取tensorboard包', 'tensorboard包', 'package', 'tensorboard', 'github.com/flipped-aurora/gin-vue-admin/server');
INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (16, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取terminal包', 'terminal包', 'package', 'terminal', 'github.com/flipped-aurora/gin-vue-admin/server');
INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (17, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取training包', 'training包', 'package', 'training', 'github.com/flipped-aurora/gin-vue-admin/server');
INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (18, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取announcement插件，使用前请确认是否为v2版本插件', 'announcement插件', 'plugin', 'announcement', 'github.com/flipped-aurora/gin-vue-admin/server');
INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (19, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取，但是缺少 initialize、plugin 结构，不建议自动化和mcp使用', 'email插件', 'plugin', 'email', 'github.com/flipped-aurora/gin-vue-admin/server');
INSERT INTO aiInfra.sys_auto_code_packages (id, created_at, updated_at, deleted_at, `desc`, label, template, package_name, module) VALUES (20, '2026-01-31 21:43:25.492', '2026-01-31 21:43:25.492', null, '系统自动读取，但是缺少 api、config、initialize、plugin、router、service 结构，不建议自动化和mcp使用', 'plugin-tool插件', 'plugin', 'plugin-tool', 'github.com/flipped-aurora/gin-vue-admin/server');

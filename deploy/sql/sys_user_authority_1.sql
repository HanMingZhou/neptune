create table sys_user_authority
(
    sys_user_id                bigint unsigned not null,
    sys_authority_authority_id bigint unsigned not null comment '角色ID',
    primary key (sys_user_id, sys_authority_authority_id)
);

INSERT INTO aiInfra.sys_user_authority (sys_user_id, sys_authority_authority_id) VALUES (1, 888);
INSERT INTO aiInfra.sys_user_authority (sys_user_id, sys_authority_authority_id) VALUES (1, 8881);
INSERT INTO aiInfra.sys_user_authority (sys_user_id, sys_authority_authority_id) VALUES (1, 9528);
INSERT INTO aiInfra.sys_user_authority (sys_user_id, sys_authority_authority_id) VALUES (2, 888);

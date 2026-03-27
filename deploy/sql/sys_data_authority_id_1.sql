create table sys_data_authority_id
(
    sys_authority_authority_id     bigint unsigned not null comment '角色ID',
    data_authority_id_authority_id bigint unsigned not null comment '角色ID',
    primary key (sys_authority_authority_id, data_authority_id_authority_id)
);

INSERT INTO aiInfra.sys_data_authority_id (sys_authority_authority_id, data_authority_id_authority_id) VALUES (888, 888);
INSERT INTO aiInfra.sys_data_authority_id (sys_authority_authority_id, data_authority_id_authority_id) VALUES (888, 8881);
INSERT INTO aiInfra.sys_data_authority_id (sys_authority_authority_id, data_authority_id_authority_id) VALUES (888, 9528);
INSERT INTO aiInfra.sys_data_authority_id (sys_authority_authority_id, data_authority_id_authority_id) VALUES (9528, 8881);
INSERT INTO aiInfra.sys_data_authority_id (sys_authority_authority_id, data_authority_id_authority_id) VALUES (9528, 9528);

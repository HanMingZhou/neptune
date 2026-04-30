create table inference_api_key_policies
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)                      null,
    updated_at datetime(3)                      null,
    deleted_at datetime(3)                      null,
    api_key_id bigint unsigned                  not null comment '关联的API Key ID',
    service_id bigint unsigned                  not null comment '关联的服务ID(0表示所有服务)',
    actions    varchar(100) default 'inference' null comment '允许的操作(inference,read,write)'
);

create index idx_inference_api_key_policies_api_key_id
    on inference_api_key_policies (api_key_id);

create index idx_inference_api_key_policies_deleted_at
    on inference_api_key_policies (deleted_at);

create index idx_inference_api_key_policies_service_id
    on inference_api_key_policies (service_id);


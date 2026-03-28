create table resource_allocations
(
    id                  bigint unsigned auto_increment
        primary key,
    created_at          datetime(3)                    null,
    updated_at          datetime(3)                    null,
    deleted_at          datetime(3)                    null,
    instance_type       varchar(32)                    not null comment '实例类型(training/inference)',
    instance_id         bigint unsigned                not null comment '实例ID',
    cluster_id          bigint unsigned                not null comment '集群ID',
    template_product_id bigint unsigned                not null comment '模板商品ID',
    product_id          bigint unsigned                not null comment '实际扣减的节点产品ID',
    node_name           varchar(100)                   not null comment '节点名称',
    schedule_strategy   varchar(20) default 'BALANCED' not null comment '调度策略',
    replica_index       bigint                         not null comment '副本索引',
    task_role           varchar(32)                    null comment '任务角色(master/worker/head/standalone)',
    reserved_count      bigint      default 1          not null comment '占用数量'
);

create index idx_allocation_owner
    on resource_allocations (instance_type, instance_id);

create index idx_resource_allocations_cluster_id
    on resource_allocations (cluster_id);

create index idx_resource_allocations_deleted_at
    on resource_allocations (deleted_at);

create index idx_resource_allocations_node_name
    on resource_allocations (node_name);

create index idx_resource_allocations_product_id
    on resource_allocations (product_id);

INSERT INTO aiInfra.resource_allocations (id, created_at, updated_at, deleted_at, instance_type, instance_id, cluster_id, template_product_id, product_id, node_name, schedule_strategy, replica_index, task_role, reserved_count) VALUES (1, '2026-03-28 19:06:57.206', '2026-03-28 19:06:57.206', null, 'inference', 33, 1, 7, 7, 'minikube', 'BALANCED', 0, 'standalone', 1);

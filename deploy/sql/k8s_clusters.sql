create table k8s_clusters
(
    id            bigint unsigned auto_increment
        primary key,
    created_at    datetime(3)      null,
    updated_at    datetime(3)      null,
    deleted_at    datetime(3)      null,
    name          varchar(100)     null comment '集群名称',
    area          varchar(50)      null comment '地域区域',
    description   varchar(500)     null comment '描述',
    kubeconfig    text             null comment 'kubeconfig内容',
    api_server    varchar(200)     null comment 'API Server地址',
    status        bigint default 1 null comment '状态(1-启用 0-停用)',
    harbor_addr   varchar(191)     null comment 'Harbor地址',
    storage_class varchar(100)     null comment 'K8s StorageClass名称',
    constraint idx_k8s_clusters_name
        unique (name)
);

create index idx_k8s_clusters_deleted_at
    on k8s_clusters (deleted_at);

INSERT INTO aiInfra.k8s_clusters (id, created_at, updated_at, deleted_at, name, area, description, kubeconfig, api_server, status, harbor_addr, storage_class) VALUES (1, '2025-12-23 00:33:07.000', '2026-02-24 10:40:55.181', null, 'minikube', '上海', 'local', 'apiVersion: v1
clusters:
- cluster:
    certificate-authority: /Users/jerrytom/.minikube/ca.crt
    extensions:
    - extension:
        last-update: Tue, 24 Feb 2026 10:37:22 CST
        provider: minikube.sigs.k8s.io
        version: v1.36.0
      name: cluster_info
    server: https://127.0.0.1:62586
  name: minikube
contexts:
- context:
    cluster: minikube
    extensions:
    - extension:
        last-update: Tue, 24 Feb 2026 10:37:22 CST
        provider: minikube.sigs.k8s.io
        version: v1.36.0
      name: context_info
    namespace: default
    user: minikube
  name: minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate: /Users/jerrytom/.minikube/profiles/minikube/client.crt
    client-key: /Users/jerrytom/.minikube/profiles/minikube/client.key', 'https://127.0.0.1:62586', 1, 'https://hub.docker.com/', 'standard');

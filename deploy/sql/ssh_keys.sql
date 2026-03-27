create table ssh_keys
(
    id          bigint unsigned auto_increment
        primary key,
    created_at  datetime(3)          null,
    updated_at  datetime(3)          null,
    deleted_at  datetime(3)          null,
    name        varchar(100)         null comment '密钥名称',
    user_id     bigint               null comment '用户ID',
    public_key  text                 null comment '公钥内容',
    fingerprint varchar(100)         null comment '密钥指纹',
    is_default  tinyint(1) default 0 null comment '是否默认密钥'
);

create index idx_ssh_keys_deleted_at
    on ssh_keys (deleted_at);

create index idx_ssh_keys_user_id
    on ssh_keys (user_id);

INSERT INTO aiInfra.ssh_keys (id, created_at, updated_at, deleted_at, name, user_id, public_key, fingerprint, is_default) VALUES (1, '2026-02-15 16:12:28.271', '2026-02-15 16:30:52.803', null, '1', 1, 'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDXI+lCfkSJbtMQ/hAT4V5YvqBOkq26P8gZRM5N5VII4U3qJN6AWkSq2mZ0ibGmZc4i6s6zRrRjHY9fJKA0ABV/xCCzjCQNSnr4mQuZlOXRNzbfLK7HPZOGMhdUIhPN8AHEHNba/WM/dAw6s0ZiXakZQWjOJhS43WgQ0/w4wccnIRW5Y4IhJHS2N9VbY2grvepeTkVYeH2nLeh4OIyAh/n2lyelv7uOe7kUQaZIiT8tx2Hc2jcSynGjRElxZD78w3CE27b+4I78ALw64ukK+jPR/pVv0giYTAEgWi0o9zusJgtYU7cnyfYhlH2gKwvso7pW2y/WDOFHKAUzbo/xo6JYKpEHkCGTTeznP1MUl5ssYEF7yfEUSlvHj939wtpEUUmqN1MdqBRQv/Mlu+ZY+CP//18qZl3FoTo4vsAG61Tqjqo2/TWN0kwhmaMwfyHcTE5AoQSLTnDq05YqTEG+VHOCdty8NAmZaiFa/AxyUPHk4Ja3WJd2T0zzfuEgNHNz6VKvt/JLV1gi965q5RE+07uagCcxeolt0CZRAQXP7E1axbmidtoHGufaEinXIhsgcoqVxsiibOtd7B/25fkMrz46KHbDijoVzZQhi5YFnqkolF9bmpHkQFfiodBxj24bifypxtAgz8IW8gTpgET4CB02/ad3iqOuZRCYSfEios6fXQ== jerrytom@jerryMacBook-Pro.local', 'ac:a3:b2:9e:4e:c4:4e:d0:8a:56:d3:f4:19:3a:14:1f', 1);

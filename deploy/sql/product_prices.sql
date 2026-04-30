create table product_prices
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)     null,
    updated_at datetime(3)     null,
    product_id bigint unsigned null comment '产品ID',
    price_type bigint          null comment '价格类型(1-小时 2-天 3-周 4-月)',
    price      decimal(20, 6)  null comment '价格',
    constraint idx_product_price_unique
        unique (product_id, price_type)
);

INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (85, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 3, 1, 1.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (86, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 3, 2, 10.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (87, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 3, 3, 100.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (88, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 3, 4, 1000.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (89, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 5, 1, 1.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (90, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 5, 2, 1.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (91, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 5, 3, 1.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (92, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 5, 4, 1.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (93, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 6, 1, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (94, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 6, 2, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (95, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 6, 3, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (96, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 6, 4, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (97, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 7, 1, 1.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (98, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 7, 2, 2.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (99, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 7, 3, 3.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (100, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 7, 4, 4.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (101, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 8, 1, 1.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (102, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 8, 2, 2.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (103, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 8, 3, 3.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (104, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 8, 4, 4.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (105, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 10, 1, 1.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (106, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 10, 2, 200.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (107, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 10, 3, 300.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (108, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 10, 4, 400.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (109, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 11, 1, 1.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (110, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 11, 2, 2.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (111, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 11, 3, 3.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (112, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 11, 4, 4.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (113, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 12, 1, 10.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (114, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 12, 2, 20.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (115, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 12, 3, 30.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (116, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 12, 4, 40.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (117, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 13, 1, 1.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (118, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 13, 2, 2.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (119, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 13, 3, 3.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (120, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 13, 4, 400.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (121, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 14, 1, 1.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (122, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 14, 2, 2.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (123, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 14, 3, 3.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (124, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 14, 4, 4.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (125, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 15, 1, 1.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (126, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 15, 2, 2.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (127, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 15, 3, 3.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (128, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 15, 4, 4.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (129, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 16, 1, 1.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (130, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 16, 2, 2.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (131, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 16, 3, 3.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (132, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 16, 4, 4.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (133, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 17, 1, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (134, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 17, 2, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (135, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 17, 3, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (136, '2026-04-17 10:47:36.643', '2026-04-29 23:46:46.070', 17, 4, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (589, '2026-04-22 21:57:35.962', '2026-04-29 23:46:46.070', 21, 1, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (590, '2026-04-22 21:57:35.962', '2026-04-29 23:46:46.070', 21, 2, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (591, '2026-04-22 21:57:35.962', '2026-04-29 23:46:46.070', 21, 3, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (592, '2026-04-22 21:57:35.962', '2026-04-29 23:46:46.070', 21, 4, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (689, '2026-04-23 10:22:57.230', '2026-04-29 23:46:46.070', 22, 1, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (690, '2026-04-23 10:22:57.230', '2026-04-29 23:46:46.070', 22, 2, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (691, '2026-04-23 10:22:57.230', '2026-04-29 23:46:46.070', 22, 3, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (692, '2026-04-23 10:22:57.230', '2026-04-29 23:46:46.070', 22, 4, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (693, '2026-04-23 10:22:57.230', '2026-04-29 23:46:46.070', 23, 1, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (694, '2026-04-23 10:22:57.230', '2026-04-29 23:46:46.070', 23, 2, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (695, '2026-04-23 10:22:57.230', '2026-04-29 23:46:46.070', 23, 3, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (696, '2026-04-23 10:22:57.230', '2026-04-29 23:46:46.070', 23, 4, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (821, '2026-04-23 10:54:23.458', '2026-04-29 23:46:46.070', 24, 1, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (822, '2026-04-23 10:54:23.458', '2026-04-29 23:46:46.070', 24, 2, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (823, '2026-04-23 10:54:23.458', '2026-04-29 23:46:46.070', 24, 3, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (824, '2026-04-23 10:54:23.458', '2026-04-29 23:46:46.070', 24, 4, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (825, '2026-04-23 10:54:23.458', '2026-04-29 23:46:46.070', 25, 1, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (826, '2026-04-23 10:54:23.458', '2026-04-29 23:46:46.070', 25, 2, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (827, '2026-04-23 10:54:23.458', '2026-04-29 23:46:46.070', 25, 3, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (828, '2026-04-23 10:54:23.458', '2026-04-29 23:46:46.070', 25, 4, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1201, '2026-04-23 19:45:25.659', '2026-04-29 23:47:16.721', 31, 1, 1.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1202, '2026-04-23 19:45:25.659', '2026-04-29 23:47:16.721', 31, 2, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1203, '2026-04-23 19:45:25.659', '2026-04-29 23:47:16.721', 31, 3, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1204, '2026-04-23 19:45:25.659', '2026-04-29 23:47:16.721', 31, 4, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1205, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 18, 1, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1206, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 18, 2, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1207, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 18, 3, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1208, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 18, 4, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1209, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 19, 1, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1210, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 19, 2, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1211, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 19, 3, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1212, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 19, 4, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1213, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 20, 1, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1214, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 20, 2, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1215, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 20, 3, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1216, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 20, 4, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1217, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 26, 1, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1218, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 26, 2, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1219, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 26, 3, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1220, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 26, 4, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1221, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 27, 1, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1222, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 27, 2, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1223, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 27, 3, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1224, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 27, 4, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1225, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 28, 1, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1226, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 28, 2, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1227, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 28, 3, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1228, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 28, 4, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1229, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 29, 1, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1230, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 29, 2, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1231, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 29, 3, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (1232, '2026-04-23 20:06:58.585', '2026-04-29 23:46:46.070', 29, 4, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (2725, '2026-04-29 22:21:59.717', '2026-04-29 23:47:12.196', 32, 1, 1.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (2726, '2026-04-29 22:21:59.717', '2026-04-29 23:47:12.196', 32, 2, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (2727, '2026-04-29 22:21:59.717', '2026-04-29 23:47:12.196', 32, 3, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (2728, '2026-04-29 22:21:59.717', '2026-04-29 23:47:12.196', 32, 4, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (2733, '2026-04-29 22:40:55.179', '2026-04-29 23:46:46.070', 30, 1, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (2734, '2026-04-29 22:40:55.179', '2026-04-29 23:46:46.070', 30, 2, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (2735, '2026-04-29 22:40:55.179', '2026-04-29 23:46:46.070', 30, 3, 0.000000);
INSERT INTO aiInfra.product_prices (id, created_at, updated_at, product_id, price_type, price) VALUES (2736, '2026-04-29 22:40:55.179', '2026-04-29 23:46:46.070', 30, 4, 0.000000);

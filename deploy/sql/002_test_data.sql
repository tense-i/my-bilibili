-- MyBilibili 测试数据
USE mybilibili;

-- 插入测试用户
INSERT INTO `user_info` (`uid`, `username`, `nickname`, `avatar`) VALUES
(1, 'user001', '测试用户1', 'https://example.com/avatar1.jpg'),
(2, 'user002', '测试用户2', 'https://example.com/avatar2.jpg'),
(3, 'user003', '测试用户3', 'https://example.com/avatar3.jpg');

-- 插入测试视频信息（50个测试视频）
INSERT INTO `video_info` (`vid`, `title`, `cover`, `author_id`, `author_name`, `region_id`, `duration`, `desc`, `pub_time`) VALUES
-- 最近24小时内的新视频（会获得1.5倍提权）
(1001, '【新视频】Go语言入门教程', 'https://example.com/cover1.jpg', 1, '测试用户1', 1, 600, '适合初学者的Go语言教程', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 12 HOUR))),
(1002, '【新视频】微服务架构实战', 'https://example.com/cover2.jpg', 2, '测试用户2', 1, 1200, '从零开始构建微服务系统', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 6 HOUR))),
(1003, '【新视频】Docker容器化部署', 'https://example.com/cover3.jpg', 3, '测试用户3', 1, 900, 'Docker实战教程', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 18 HOUR))),
(1004, '【新视频】Kubernetes集群管理', 'https://example.com/cover4.jpg', 1, '测试用户1', 1, 1500, 'K8s从入门到精通', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 20 HOUR))),
(1005, '【新视频】gRPC微服务通信', 'https://example.com/cover5.jpg', 2, '测试用户2', 1, 800, 'gRPC实战案例', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 10 HOUR))),

-- 普通视频（2-7天前）
(1006, 'MySQL数据库优化技巧', 'https://example.com/cover6.jpg', 3, '测试用户3', 2, 1000, 'MySQL性能优化全攻略', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 2 DAY))),
(1007, 'Redis缓存设计模式', 'https://example.com/cover7.jpg', 1, '测试用户1', 2, 700, 'Redis最佳实践', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 3 DAY))),
(1008, '分布式事务解决方案', 'https://example.com/cover8.jpg', 2, '测试用户2', 1, 1100, '分布式系统事务处理', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 4 DAY))),
(1009, '消息队列Kafka实战', 'https://example.com/cover9.jpg', 3, '测试用户3', 1, 950, 'Kafka从入门到实战', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 5 DAY))),
(1010, 'Elasticsearch搜索引擎', 'https://example.com/cover10.jpg', 1, '测试用户1', 2, 1300, 'ES全文检索实战', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 6 DAY))),

-- 更多测试视频（7-30天前）
(1011, 'Go并发编程详解', 'https://example.com/cover11.jpg', 2, '测试用户2', 1, 850, 'Goroutine和Channel', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 10 DAY))),
(1012, '前端React框架入门', 'https://example.com/cover12.jpg', 3, '测试用户3', 3, 950, 'React从零开始', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 15 DAY))),
(1013, 'Vue3组件化开发', 'https://example.com/cover13.jpg', 1, '测试用户1', 3, 800, 'Vue3实战教程', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 20 DAY))),
(1014, 'Python数据分析实战', 'https://example.com/cover14.jpg', 2, '测试用户2', 4, 1100, 'Pandas和NumPy应用', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 25 DAY))),
(1015, '机器学习算法详解', 'https://example.com/cover15.jpg', 3, '测试用户3', 4, 1400, '常见ML算法讲解', UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 30 DAY)));

-- 插入测试视频统计数据
INSERT INTO `video_stat` (`vid`, `view`, `like_count`, `coin`, `fav`, `share`, `reply`, `danmaku`) VALUES
-- 新视频（数据较少但会有提权）
(1001, 5000, 300, 150, 200, 50, 100, 80),
(1002, 8000, 500, 250, 350, 80, 150, 120),
(1003, 6500, 400, 200, 280, 60, 120, 90),
(1004, 7200, 450, 220, 320, 70, 140, 110),
(1005, 5800, 350, 180, 240, 55, 110, 85),

-- 普通视频（数据较多但无提权）
(1006, 50000, 3000, 1500, 2500, 500, 1200, 800),
(1007, 65000, 4000, 2000, 3200, 650, 1500, 1000),
(1008, 45000, 2800, 1400, 2300, 450, 1100, 750),
(1009, 58000, 3500, 1800, 2800, 580, 1400, 900),
(1010, 72000, 4500, 2300, 3500, 720, 1600, 1100),

-- 更多视频
(1011, 38000, 2300, 1200, 1900, 380, 900, 600),
(1012, 42000, 2600, 1300, 2100, 420, 1000, 700),
(1013, 35000, 2100, 1100, 1800, 350, 850, 550),
(1014, 48000, 2900, 1500, 2400, 480, 1150, 800),
(1015, 55000, 3300, 1700, 2700, 550, 1300, 900);

-- 插入 academy_archive 记录（用于热度计算）
-- 注意：初始 hot 值为 0，会由 hotrank-job 定时任务计算更新
INSERT INTO `academy_archive` (`oid`, `business`, `region_id`, `pub_time`, `state`) 
SELECT `vid`, 1, `region_id`, `pub_time`, 0 FROM `video_info`;

-- 验证数据
SELECT '=== 数据统计 ===' as '';
SELECT COUNT(*) as '用户总数' FROM `user_info`;
SELECT COUNT(*) as '视频总数' FROM `video_info`;
SELECT COUNT(*) as '统计数据总数' FROM `video_stat`;
SELECT COUNT(*) as 'Academy记录总数' FROM `academy_archive`;

SELECT '=== 新视频列表（24小时内）===' as '';
SELECT vid, title, FROM_UNIXTIME(pub_time) as pub_datetime 
FROM `video_info` 
WHERE pub_time >= UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 1 DAY))
ORDER BY pub_time DESC;


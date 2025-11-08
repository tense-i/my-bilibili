-- 推荐系统测试数据

-- 清空现有测试数据
TRUNCATE TABLE video_info;
TRUNCATE TABLE video_tag;
TRUNCATE TABLE user_behavior;
TRUNCATE TABLE user_follow;
TRUNCATE TABLE user_blacklist;

-- ===== 1. 插入更多视频数据 =====
INSERT INTO `video_info` (`avid`, `cid`, `mid`, `title`, `cover`, `duration`, `pub_time`, `zone_id`, `state`, 
  `play_hive`, `likes_hive`, `fav_hive`, `share_hive`, `coin_hive`, `reply_hive`,
  `play_month`, `likes_month`, `share_month`, `reply_month`, `play_month_finish`)
VALUES
  -- 动画区视频
  (100001, 100001, 1001, '【MMD】初音未来 - 千本樱', 'http://example.com/cover1.jpg', 240, UNIX_TIMESTAMP(NOW() - INTERVAL 1 DAY), 20, 5,
   150000, 8000, 3000, 1000, 2000, 500,
   50000, 2500, 300, 150, 35000),
  (100002, 100002, 1001, '【手书】VOCALOID - 洛天依', 'http://example.com/cover2.jpg', 180, UNIX_TIMESTAMP(NOW() - INTERVAL 2 DAY), 20, 4,
   80000, 4500, 2000, 500, 1000, 300,
   30000, 1500, 200, 100, 20000),
  (100004, 100004, 1001, '【MMD】初音未来 - 甩葱歌', 'http://example.com/cover4.jpg', 200, UNIX_TIMESTAMP(NOW() - INTERVAL 5 DAY), 20, 4,
   120000, 6000, 2500, 800, 1500, 400,
   40000, 2000, 250, 120, 28000),
  (100005, 100005, 1003, '【AMV】命运石之门混剪', 'http://example.com/cover5.jpg', 220, UNIX_TIMESTAMP(NOW() - INTERVAL 3 DAY), 24, 3,
   60000, 3500, 1500, 400, 800, 200,
   25000, 1200, 150, 80, 18000),
  (100006, 100006, 1003, '【MAD】进击的巨人燃向剪辑', 'http://example.com/cover6.jpg', 280, UNIX_TIMESTAMP(NOW() - INTERVAL 4 DAY), 24, 4,
   95000, 5000, 2200, 600, 1200, 350,
   35000, 1800, 200, 100, 25000),
  
  -- 游戏区视频
  (100003, 100003, 1002, '【我的世界】超大型城堡建筑教程', 'http://example.com/cover3.jpg', 600, UNIX_TIMESTAMP(NOW() - INTERVAL 3 DAY), 17, 4,
   200000, 12000, 5000, 2000, 3000, 800,
   80000, 4000, 500, 300, 60000),
  (100007, 100007, 1002, '【我的世界】红石音乐演奏装置', 'http://example.com/cover7.jpg', 450, UNIX_TIMESTAMP(NOW() - INTERVAL 6 DAY), 17, 3,
   85000, 4800, 2100, 700, 1300, 320,
   32000, 1600, 180, 90, 22000),
  (100008, 100008, 1002, '【我的世界】自动化农场设计', 'http://example.com/cover8.jpg', 380, UNIX_TIMESTAMP(NOW() - INTERVAL 7 DAY), 17, 3,
   72000, 4200, 1900, 600, 1100, 280,
   28000, 1400, 160, 85, 20000),
  (100009, 100009, 1004, '【原神】新手攻略完全指南', 'http://example.com/cover9.jpg', 420, UNIX_TIMESTAMP(NOW() - INTERVAL 2 DAY), 17, 4,
   180000, 10000, 4500, 1800, 2800, 700,
   75000, 3800, 450, 280, 55000),
  (100010, 100010, 1004, '【原神】五星角色抽卡分析', 'http://example.com/cover10.jpg', 340, UNIX_TIMESTAMP(NOW() - INTERVAL 1 DAY), 17, 3,
   110000, 6500, 2800, 900, 1600, 450,
   45000, 2300, 280, 150, 32000),
  
  -- 音乐区视频
  (100011, 100011, 1005, '【钢琴】周杰伦经典串烧', 'http://example.com/cover11.jpg', 300, UNIX_TIMESTAMP(NOW() - INTERVAL 4 DAY), 28, 5,
   250000, 15000, 6000, 2500, 3500, 900,
   90000, 4500, 550, 320, 65000),
  (100012, 100012, 1005, '【吉他】民谣弹唱合集', 'http://example.com/cover12.jpg', 280, UNIX_TIMESTAMP(NOW() - INTERVAL 5 DAY), 28, 4,
   130000, 7500, 3200, 1200, 1800, 520,
   50000, 2600, 300, 180, 38000),
  
  -- 生活区视频
  (100013, 100013, 1006, '【美食】正宗四川麻辣火锅', 'http://example.com/cover13.jpg', 480, UNIX_TIMESTAMP(NOW() - INTERVAL 2 DAY), 76, 4,
   160000, 9000, 3800, 1500, 2200, 650,
   60000, 3000, 350, 220, 45000),
  (100014, 100014, 1006, '【旅行】日本京都樱花季', 'http://example.com/cover14.jpg', 520, UNIX_TIMESTAMP(NOW() - INTERVAL 3 DAY), 21, 4,
   140000, 8200, 3500, 1300, 2000, 580,
   55000, 2800, 320, 200, 42000),
  
  -- 知识区视频
  (100015, 100015, 1007, '【科普】量子力学入门', 'http://example.com/cover15.jpg', 720, UNIX_TIMESTAMP(NOW() - INTERVAL 1 DAY), 36, 5,
   220000, 13000, 5500, 2200, 3200, 850,
   85000, 4200, 500, 300, 62000),
  (100016, 100016, 1007, '【历史】大航海时代解析', 'http://example.com/cover16.jpg', 650, UNIX_TIMESTAMP(NOW() - INTERVAL 2 DAY), 37, 4,
   170000, 9500, 4000, 1600, 2400, 700,
   68000, 3400, 400, 240, 50000),
  
  -- 时尚区视频
  (100017, 100017, 1008, '【穿搭】夏季清爽搭配指南', 'http://example.com/cover17.jpg', 360, UNIX_TIMESTAMP(NOW() - INTERVAL 3 DAY), 155, 3,
   98000, 5800, 2400, 800, 1400, 420,
   38000, 1900, 220, 130, 28000),
  (100018, 100018, 1008, '【美妆】日系妆容教程', 'http://example.com/cover18.jpg', 320, UNIX_TIMESTAMP(NOW() - INTERVAL 4 DAY), 155, 4,
   115000, 6800, 2900, 950, 1650, 480,
   45000, 2250, 260, 155, 33000),
  
  -- 鬼畜区视频
  (100019, 100019, 1009, '【鬼畜】经典鬼畜合集', 'http://example.com/cover19.jpg', 240, UNIX_TIMESTAMP(NOW() - INTERVAL 1 DAY), 119, 4,
   280000, 16000, 6500, 2800, 3800, 950,
   95000, 4800, 580, 340, 70000),
  (100020, 100020, 1009, '【鬼畜】名场面混剪', 'http://example.com/cover20.jpg', 200, UNIX_TIMESTAMP(NOW() - INTERVAL 2 DAY), 119, 3,
   190000, 11000, 4500, 1900, 2600, 720,
   72000, 3600, 420, 260, 53000);

-- ===== 2. 插入视频标签数据 =====
INSERT INTO `video_tag` (`avid`, `tag_id`, `tag_name`, `tag_type`)
VALUES
  -- 动画区标签
  (100001, 1, 'MMD', 2),
  (100001, 2, '初音未来', 2),
  (100001, 3, '千本樱', 2),
  (100001, 101, 'VOCALOID', 2),
  
  (100002, 1, 'MMD', 2),
  (100002, 4, '洛天依', 2),
  (100002, 101, 'VOCALOID', 2),
  
  (100004, 1, 'MMD', 2),
  (100004, 2, '初音未来', 2),
  (100004, 102, '甩葱歌', 2),
  
  (100005, 103, 'AMV', 2),
  (100005, 104, '命运石之门', 2),
  (100005, 105, '动漫混剪', 2),
  
  (100006, 106, 'MAD', 2),
  (100006, 107, '进击的巨人', 2),
  (100006, 108, '燃向', 2),
  
  -- 游戏区标签
  (100003, 5, '我的世界', 2),
  (100003, 6, '建筑', 2),
  (100003, 109, '教程', 2),
  
  (100007, 5, '我的世界', 2),
  (100007, 110, '红石', 2),
  (100007, 111, '音乐', 2),
  
  (100008, 5, '我的世界', 2),
  (100008, 112, '自动化', 2),
  (100008, 113, '农场', 2),
  
  (100009, 114, '原神', 2),
  (100009, 115, '攻略', 2),
  (100009, 116, '新手向', 2),
  
  (100010, 114, '原神', 2),
  (100010, 117, '抽卡', 2),
  (100010, 118, '分析', 2),
  
  -- 音乐区标签
  (100011, 119, '钢琴', 2),
  (100011, 120, '周杰伦', 2),
  (100011, 121, '串烧', 2),
  
  (100012, 122, '吉他', 2),
  (100012, 123, '民谣', 2),
  (100012, 124, '弹唱', 2),
  
  -- 生活区标签
  (100013, 125, '美食', 2),
  (100013, 126, '火锅', 2),
  (100013, 127, '四川', 2),
  
  (100014, 128, '旅行', 2),
  (100014, 129, '日本', 2),
  (100014, 130, '京都', 2),
  
  -- 知识区标签
  (100015, 131, '科普', 2),
  (100015, 132, '量子力学', 2),
  (100015, 133, '物理', 2),
  
  (100016, 134, '历史', 2),
  (100016, 135, '大航海时代', 2),
  (100016, 136, '地理', 2),
  
  -- 时尚区标签
  (100017, 137, '穿搭', 2),
  (100017, 138, '夏季', 2),
  (100017, 139, '搭配', 2),
  
  (100018, 140, '美妆', 2),
  (100018, 141, '日系', 2),
  (100018, 142, '教程', 2),
  
  -- 鬼畜区标签
  (100019, 143, '鬼畜', 2),
  (100019, 144, '合集', 2),
  (100019, 145, '经典', 2),
  
  (100020, 143, '鬼畜', 2),
  (100020, 146, '混剪', 2),
  (100020, 147, '名场面', 2);

-- ===== 3. 插入用户行为数据 =====
INSERT INTO `user_behavior` (`mid`, `avid`, `behavior_type`, `duration`, `finish_rate`, `ctime`)
VALUES
  -- 用户 1001 的行为（喜欢动画）
  (1001, 100001, 1, 240, 100.00, UNIX_TIMESTAMP(NOW() - INTERVAL 2 HOUR)),
  (1001, 100001, 2, NULL, NULL, UNIX_TIMESTAMP(NOW() - INTERVAL 2 HOUR)),
  (1001, 100002, 1, 180, 100.00, UNIX_TIMESTAMP(NOW() - INTERVAL 5 HOUR)),
  (1001, 100002, 2, NULL, NULL, UNIX_TIMESTAMP(NOW() - INTERVAL 5 HOUR)),
  (1001, 100004, 1, 200, 100.00, UNIX_TIMESTAMP(NOW() - INTERVAL 1 DAY)),
  (1001, 100005, 1, 150, 68.18, UNIX_TIMESTAMP(NOW() - INTERVAL 6 HOUR)),
  
  -- 用户 1002 的行为（喜欢游戏）
  (1002, 100003, 1, 600, 100.00, UNIX_TIMESTAMP(NOW() - INTERVAL 3 HOUR)),
  (1002, 100003, 2, NULL, NULL, UNIX_TIMESTAMP(NOW() - INTERVAL 3 HOUR)),
  (1002, 100007, 1, 450, 100.00, UNIX_TIMESTAMP(NOW() - INTERVAL 8 HOUR)),
  (1002, 100008, 1, 380, 100.00, UNIX_TIMESTAMP(NOW() - INTERVAL 1 DAY)),
  (1002, 100009, 1, 420, 100.00, UNIX_TIMESTAMP(NOW() - INTERVAL 4 HOUR)),
  (1002, 100009, 2, NULL, NULL, UNIX_TIMESTAMP(NOW() - INTERVAL 4 HOUR)),
  
  -- 用户 1003 的行为（喜欢音乐）
  (1003, 100011, 1, 300, 100.00, UNIX_TIMESTAMP(NOW() - INTERVAL 2 HOUR)),
  (1003, 100011, 2, NULL, NULL, UNIX_TIMESTAMP(NOW() - INTERVAL 2 HOUR)),
  (1003, 100012, 1, 280, 100.00, UNIX_TIMESTAMP(NOW() - INTERVAL 6 HOUR)),
  (1003, 100012, 3, NULL, NULL, UNIX_TIMESTAMP(NOW() - INTERVAL 6 HOUR)),
  
  -- 用户 1004 的行为（喜欢知识）
  (1004, 100015, 1, 720, 100.00, UNIX_TIMESTAMP(NOW() - INTERVAL 1 HOUR)),
  (1004, 100015, 2, NULL, NULL, UNIX_TIMESTAMP(NOW() - INTERVAL 1 HOUR)),
  (1004, 100015, 3, NULL, NULL, UNIX_TIMESTAMP(NOW() - INTERVAL 1 HOUR)),
  (1004, 100016, 1, 650, 100.00, UNIX_TIMESTAMP(NOW() - INTERVAL 5 HOUR)),
  (1004, 100016, 2, NULL, NULL, UNIX_TIMESTAMP(NOW() - INTERVAL 5 HOUR)),
  
  -- 用户 1005 的行为（兴趣广泛）
  (1005, 100001, 1, 240, 100.00, UNIX_TIMESTAMP(NOW() - INTERVAL 3 HOUR)),
  (1005, 100003, 1, 600, 100.00, UNIX_TIMESTAMP(NOW() - INTERVAL 7 HOUR)),
  (1005, 100011, 1, 300, 100.00, UNIX_TIMESTAMP(NOW() - INTERVAL 10 HOUR)),
  (1005, 100013, 1, 480, 100.00, UNIX_TIMESTAMP(NOW() - INTERVAL 12 HOUR)),
  (1005, 100019, 1, 240, 100.00, UNIX_TIMESTAMP(NOW() - INTERVAL 1 HOUR)),
  (1005, 100019, 2, NULL, NULL, UNIX_TIMESTAMP(NOW() - INTERVAL 1 HOUR));

-- ===== 4. 插入用户关注数据 =====
INSERT INTO `user_follow` (`mid`, `up_mid`, `status`, `ctime`)
VALUES
  -- 用户 1001 关注
  (1001, 1001, 1, UNIX_TIMESTAMP(NOW() - INTERVAL 10 DAY)),
  (1001, 1003, 1, UNIX_TIMESTAMP(NOW() - INTERVAL 8 DAY)),
  (1001, 1005, 1, UNIX_TIMESTAMP(NOW() - INTERVAL 5 DAY)),
  
  -- 用户 1002 关注
  (1002, 1002, 1, UNIX_TIMESTAMP(NOW() - INTERVAL 15 DAY)),
  (1002, 1004, 1, UNIX_TIMESTAMP(NOW() - INTERVAL 7 DAY)),
  
  -- 用户 1003 关注
  (1003, 1005, 1, UNIX_TIMESTAMP(NOW() - INTERVAL 12 DAY)),
  (1003, 1006, 1, UNIX_TIMESTAMP(NOW() - INTERVAL 6 DAY)),
  
  -- 用户 1004 关注
  (1004, 1007, 1, UNIX_TIMESTAMP(NOW() - INTERVAL 20 DAY)),
  (1004, 1006, 1, UNIX_TIMESTAMP(NOW() - INTERVAL 9 DAY)),
  
  -- 用户 1005 关注（广泛）
  (1005, 1001, 1, UNIX_TIMESTAMP(NOW() - INTERVAL 25 DAY)),
  (1005, 1002, 1, UNIX_TIMESTAMP(NOW() - INTERVAL 20 DAY)),
  (1005, 1005, 1, UNIX_TIMESTAMP(NOW() - INTERVAL 15 DAY)),
  (1005, 1006, 1, UNIX_TIMESTAMP(NOW() - INTERVAL 10 DAY)),
  (1005, 1009, 1, UNIX_TIMESTAMP(NOW() - INTERVAL 5 DAY));

-- ===== 5. 插入黑名单数据（示例）=====
INSERT INTO `user_blacklist` (`mid`, `up_mid`)
VALUES
  (1001, 1009),  -- 用户 1001 拉黑了 UP主 1009
  (1002, 1008);  -- 用户 1002 拉黑了 UP主 1008

-- 查询统计
SELECT '====== 数据统计 ======' AS info;
SELECT '视频总数' AS type, COUNT(*) AS count FROM video_info
UNION ALL
SELECT '标签总数', COUNT(*) FROM video_tag
UNION ALL
SELECT '用户行为总数', COUNT(*) FROM user_behavior
UNION ALL
SELECT '用户关注总数', COUNT(*) FROM user_follow
UNION ALL
SELECT '黑名单总数', COUNT(*) FROM user_blacklist;

SELECT '====== 按分区统计视频 ======' AS info;
SELECT zone_id, COUNT(*) AS video_count FROM video_info GROUP BY zone_id ORDER BY video_count DESC;

SELECT '====== 按状态统计视频 ======' AS info;
SELECT 
  CASE state
    WHEN 1 THEN '正常'
    WHEN 3 THEN '回查可放出'
    WHEN 4 THEN '优质'
    WHEN 5 THEN '精选'
  END AS state_name,
  COUNT(*) AS count
FROM video_info
GROUP BY state;


-- 推荐系统相关表

-- 视频信息表
CREATE TABLE IF NOT EXISTS `video_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `avid` bigint(20) NOT NULL COMMENT '视频AV号',
  `cid` bigint(20) NOT NULL COMMENT '分P的CID',
  `mid` bigint(20) NOT NULL COMMENT 'UP主MID',
  `title` varchar(255) NOT NULL COMMENT '视频标题',
  `cover` varchar(512) DEFAULT NULL COMMENT '封面URL',
  `duration` int(11) NOT NULL COMMENT '视频时长(秒)',
  `pub_time` bigint(20) NOT NULL COMMENT '发布时间戳',
  `zone_id` int(11) NOT NULL COMMENT '分区ID',
  `state` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态: 1-正常 3-回查可放出 4-优质 5-精选',
  
  -- 统计数据（全站）
  `play_hive` bigint(20) DEFAULT '0' COMMENT 'B站全站播放量',
  `likes_hive` bigint(20) DEFAULT '0' COMMENT 'B站全站点赞数',
  `fav_hive` bigint(20) DEFAULT '0' COMMENT 'B站全站收藏数',
  `share_hive` bigint(20) DEFAULT '0' COMMENT 'B站全站分享数',
  `coin_hive` bigint(20) DEFAULT '0' COMMENT 'B站全站投币数',
  `reply_hive` bigint(20) DEFAULT '0' COMMENT 'B站全站评论数',
  
  -- 统计数据（月度）
  `play_month` bigint(20) DEFAULT '0' COMMENT '近30天播放量',
  `likes_month` bigint(20) DEFAULT '0' COMMENT '近30天点赞数',
  `share_month` bigint(20) DEFAULT '0' COMMENT '近30天分享数',
  `reply_month` bigint(20) DEFAULT '0' COMMENT '近30天评论数',
  `play_month_finish` bigint(20) DEFAULT '0' COMMENT '近30天完播量',
  
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_avid` (`avid`),
  KEY `idx_mid` (`mid`),
  KEY `idx_zone` (`zone_id`),
  KEY `idx_state` (`state`),
  KEY `idx_pubtime` (`pub_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频信息表';

-- 视频标签表
CREATE TABLE IF NOT EXISTS `video_tag` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `avid` bigint(20) NOT NULL COMMENT '视频AV号',
  `tag_id` int(11) NOT NULL COMMENT '标签ID',
  `tag_name` varchar(64) NOT NULL COMMENT '标签名称',
  `tag_type` tinyint(4) DEFAULT '1' COMMENT '标签类型: 1-分类标签 2-内容标签',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_avid` (`avid`),
  KEY `idx_tag_id` (`tag_id`),
  KEY `idx_tag_name` (`tag_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频标签表';

-- 用户行为表
CREATE TABLE IF NOT EXISTS `user_behavior` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `mid` bigint(20) NOT NULL COMMENT '用户MID',
  `avid` bigint(20) NOT NULL COMMENT '视频AV号',
  `behavior_type` tinyint(4) NOT NULL COMMENT '行为类型: 1-播放 2-点赞 3-收藏 4-分享 5-关注',
  `duration` int(11) DEFAULT NULL COMMENT '观看时长(秒)',
  `finish_rate` decimal(5,2) DEFAULT NULL COMMENT '完播率',
  `ctime` bigint(20) NOT NULL COMMENT '行为时间戳',
  PRIMARY KEY (`id`),
  KEY `idx_mid_type` (`mid`, `behavior_type`),
  KEY `idx_avid` (`avid`),
  KEY `idx_ctime` (`ctime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户行为表';

-- 用户关注表
CREATE TABLE IF NOT EXISTS `user_follow` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `mid` bigint(20) NOT NULL COMMENT '用户MID',
  `up_mid` bigint(20) NOT NULL COMMENT 'UP主MID',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态: 1-关注 0-取消关注',
  `ctime` bigint(20) NOT NULL COMMENT '关注时间戳',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_mid_up` (`mid`, `up_mid`),
  KEY `idx_up_mid` (`up_mid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户关注表';

-- 用户黑名单表
CREATE TABLE IF NOT EXISTS `user_blacklist` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `mid` bigint(20) NOT NULL COMMENT '用户MID',
  `up_mid` bigint(20) NOT NULL COMMENT 'UP主MID',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_mid_up` (`mid`, `up_mid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户黑名单表';

-- 插入测试数据
INSERT INTO `video_info` (`avid`, `cid`, `mid`, `title`, `cover`, `duration`, `pub_time`, `zone_id`, `state`, 
  `play_hive`, `likes_hive`, `fav_hive`, `share_hive`, `coin_hive`, `reply_hive`,
  `play_month`, `likes_month`, `share_month`, `reply_month`, `play_month_finish`)
VALUES
  (100001, 100001, 1001, '【MMD】初音未来 - 千本樱', 'http://example.com/cover1.jpg', 240, UNIX_TIMESTAMP(NOW() - INTERVAL 1 DAY), 20, 5,
   150000, 8000, 3000, 1000, 2000, 500,
   50000, 2500, 300, 150, 35000),
  (100002, 100002, 1001, '【手书】VOCALOID - 洛天依', 'http://example.com/cover2.jpg', 180, UNIX_TIMESTAMP(NOW() - INTERVAL 2 DAY), 20, 4,
   80000, 4500, 2000, 500, 1000, 300,
   30000, 1500, 200, 100, 20000),
  (100003, 100003, 1002, '【游戏】我的世界建筑教程', 'http://example.com/cover3.jpg', 600, UNIX_TIMESTAMP(NOW() - INTERVAL 3 DAY), 17, 4,
   200000, 12000, 5000, 2000, 3000, 800,
   80000, 4000, 500, 300, 60000);

INSERT INTO `video_tag` (`avid`, `tag_id`, `tag_name`, `tag_type`)
VALUES
  (100001, 1, 'MMD', 2),
  (100001, 2, '初音未来', 2),
  (100001, 3, '千本樱', 2),
  (100002, 1, 'MMD', 2),
  (100002, 4, '洛天依', 2),
  (100003, 5, '我的世界', 2),
  (100003, 6, '建筑', 2);

-- 初始化 Redis 召回索引的说明
-- 以下命令需要在 Redis 中执行：
-- 
-- 热门视频索引:
-- ZADD recall:hot:default 95.6 100001 94.2 100002 93.8 100003
-- 
-- 精选视频索引:
-- LPUSH recall:selection 100001 100002
--
-- 标签索引示例:
-- ZADD recall:tag:1 95.0 100001 94.0 100002
--
-- UP主视频索引示例:
-- ZADD recall:up:1001 1736294400 100001 1736208000 100002


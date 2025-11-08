-- MyBilibili 数据库初始化脚本
-- 完全参考主项目 Bilibili 的表结构设计

CREATE DATABASE IF NOT EXISTS mybilibili DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE mybilibili;

-- ========================================
-- 视频相关表
-- ========================================

-- 视频基本信息表（参考主项目）
CREATE TABLE IF NOT EXISTS `video_info` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `vid` bigint NOT NULL COMMENT '视频ID',
  `title` varchar(255) NOT NULL DEFAULT '' COMMENT '标题',
  `cover` varchar(512) NOT NULL DEFAULT '' COMMENT '封面URL',
  `author_id` bigint NOT NULL DEFAULT 0 COMMENT '作者ID',
  `author_name` varchar(64) NOT NULL DEFAULT '' COMMENT '作者名称',
  `region_id` int NOT NULL DEFAULT 0 COMMENT '分区ID',
  `duration` int NOT NULL DEFAULT 0 COMMENT '时长（秒）',
  `desc` text COMMENT '简介',
  `pub_time` bigint NOT NULL DEFAULT 0 COMMENT '发布时间（Unix时间戳）',
  `state` tinyint NOT NULL DEFAULT 0 COMMENT '状态：0-正常，1-删除',
  `ctime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_vid` (`vid`),
  KEY `idx_author` (`author_id`, `state`),
  KEY `idx_region` (`region_id`, `state`),
  KEY `idx_pub_time` (`pub_time`, `state`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频基本信息表';

-- 视频统计表（参考主项目）
CREATE TABLE IF NOT EXISTS `video_stat` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `vid` bigint NOT NULL COMMENT '视频ID',
  `view` bigint NOT NULL DEFAULT 0 COMMENT '播放数',
  `like_count` bigint NOT NULL DEFAULT 0 COMMENT '点赞数',
  `coin` bigint NOT NULL DEFAULT 0 COMMENT '硬币数',
  `fav` bigint NOT NULL DEFAULT 0 COMMENT '收藏数',
  `share` bigint NOT NULL DEFAULT 0 COMMENT '分享数',
  `reply` bigint NOT NULL DEFAULT 0 COMMENT '评论数',
  `danmaku` bigint NOT NULL DEFAULT 0 COMMENT '弹幕数',
  `ctime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_vid` (`vid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频统计表';

-- ========================================
-- 热门排行榜相关表（完全参考主项目的 academy_archive 表）
-- ========================================

-- 热度记录表（参考主项目 app/admin/main/creative 的 academy_archive 表）
CREATE TABLE IF NOT EXISTS `academy_archive` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增ID（用于游标分页）',
  `oid` bigint NOT NULL COMMENT '对象ID（视频ID或专栏ID）',
  `hot` bigint NOT NULL DEFAULT 0 COMMENT '热度值',
  `business` tinyint NOT NULL DEFAULT 1 COMMENT '业务类型：1-视频，2-专栏',
  `region_id` int NOT NULL DEFAULT 0 COMMENT '分区ID',
  `pub_time` bigint NOT NULL DEFAULT 0 COMMENT '发布时间（Unix时间戳）',
  `state` tinyint NOT NULL DEFAULT 0 COMMENT '状态：0-正常，1-删除',
  `comment` varchar(255) DEFAULT '' COMMENT '备注',
  `ctime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_oid` (`oid`),
  KEY `idx_business_state_id` (`business`, `state`, `id`),
  KEY `idx_hot_global` (`hot`, `state`, `business`),
  KEY `idx_hot_region` (`region_id`, `hot`, `state`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='热度排行榜表';

-- ========================================
-- 用户相关表（简化版）
-- ========================================

-- 用户基本信息表
CREATE TABLE IF NOT EXISTS `user_info` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `uid` bigint NOT NULL COMMENT '用户ID',
  `username` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `nickname` varchar(64) NOT NULL DEFAULT '' COMMENT '昵称',
  `avatar` varchar(512) NOT NULL DEFAULT '' COMMENT '头像URL',
  `state` tinyint NOT NULL DEFAULT 0 COMMENT '状态：0-正常，1-封禁',
  `ctime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_uid` (`uid`),
  KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户基本信息表';


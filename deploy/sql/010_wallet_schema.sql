-- =====================================================
-- MyBilibili 虚拟钱包系统数据库初始化脚本
-- 版本: v2.0.0
-- 创建时间: 2025-11-09
-- 说明: 严格按照bilibili主项目设计
-- =====================================================

USE mybilibili;

-- =====================================================
-- 1. 用户钱包表（分10张表）
-- =====================================================

CREATE TABLE IF NOT EXISTS `user_wallet_0` (
  `uid` bigint(20) NOT NULL COMMENT '用户ID',
  `gold` bigint(20) NOT NULL DEFAULT '0' COMMENT '金瓜子（Android/PC/H5）',
  `iap_gold` bigint(20) NOT NULL DEFAULT '0' COMMENT 'IAP金瓜子（iOS专用）',
  `silver` bigint(20) NOT NULL DEFAULT '0' COMMENT '银瓜子（免费货币）',
  `gold_recharge_cnt` bigint(20) NOT NULL DEFAULT '0' COMMENT '累计充值金瓜子',
  `gold_pay_cnt` bigint(20) NOT NULL DEFAULT '0' COMMENT '累计消费金瓜子',
  `silver_pay_cnt` bigint(20) NOT NULL DEFAULT '0' COMMENT '累计消费银瓜子',
  `cost_base` bigint(20) NOT NULL DEFAULT '0' COMMENT '消费基数',
  
  -- 快照字段（用于每日对账）
  `snapshot_time` datetime DEFAULT NULL COMMENT '快照时间（每日0点更新）',
  `snapshot_gold` bigint(20) DEFAULT '0' COMMENT '快照金瓜子',
  `snapshot_iap_gold` bigint(20) DEFAULT '0' COMMENT '快照IAP金瓜子',
  `snapshot_silver` bigint(20) DEFAULT '0' COMMENT '快照银瓜子',
  
  `reserved1` bigint(20) DEFAULT '0' COMMENT '预留字段1',
  `reserved2` varchar(255) DEFAULT '' COMMENT '预留字段2',
  `ctime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户钱包表（分表0）';

-- 创建其他9张分表
CREATE TABLE IF NOT EXISTS `user_wallet_1` LIKE `user_wallet_0`;
CREATE TABLE IF NOT EXISTS `user_wallet_2` LIKE `user_wallet_0`;
CREATE TABLE IF NOT EXISTS `user_wallet_3` LIKE `user_wallet_0`;
CREATE TABLE IF NOT EXISTS `user_wallet_4` LIKE `user_wallet_0`;
CREATE TABLE IF NOT EXISTS `user_wallet_5` LIKE `user_wallet_0`;
CREATE TABLE IF NOT EXISTS `user_wallet_6` LIKE `user_wallet_0`;
CREATE TABLE IF NOT EXISTS `user_wallet_7` LIKE `user_wallet_0`;
CREATE TABLE IF NOT EXISTS `user_wallet_8` LIKE `user_wallet_0`;
CREATE TABLE IF NOT EXISTS `user_wallet_9` LIKE `user_wallet_0`;

-- 修改其他表的注释
ALTER TABLE `user_wallet_1` COMMENT='用户钱包表（分表1）';
ALTER TABLE `user_wallet_2` COMMENT='用户钱包表（分表2）';
ALTER TABLE `user_wallet_3` COMMENT='用户钱包表（分表3）';
ALTER TABLE `user_wallet_4` COMMENT='用户钱包表（分表4）';
ALTER TABLE `user_wallet_5` COMMENT='用户钱包表（分表5）';
ALTER TABLE `user_wallet_6` COMMENT='用户钱包表（分表6）';
ALTER TABLE `user_wallet_7` COMMENT='用户钱包表（分表7）';
ALTER TABLE `user_wallet_8` COMMENT='用户钱包表（分表8）';
ALTER TABLE `user_wallet_9` COMMENT='用户钱包表（分表9）';

-- =====================================================
-- 2. 流水记录表（分10张表）
-- =====================================================

CREATE TABLE IF NOT EXISTS `coin_stream_record_0` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL COMMENT '用户ID',
  `transaction_id` varchar(64) NOT NULL COMMENT '交易ID（幂等性保证）',
  `extend_tid` varchar(64) DEFAULT '' COMMENT '扩展交易ID',
  `coin_type` int(11) NOT NULL COMMENT '币种类型：1=gold 2=iap_gold 3=silver',
  `delta_coin_num` bigint(20) NOT NULL COMMENT '变化金额（正数=增加，负数=减少）',
  `org_coin_num` bigint(20) NOT NULL COMMENT '变化前余额',
  `op_result` int(11) NOT NULL COMMENT '操作结果：1=增加成功 2=减少成功 -1=增加失败 -2=减少失败',
  `op_reason` int(11) NOT NULL DEFAULT '0' COMMENT '失败原因：0=成功 1=余额不足 2=参数错误 3=锁失败',
  `op_type` int(11) NOT NULL COMMENT '操作类型：1=充值 2=消费 3=兑换',
  `op_time` datetime NOT NULL COMMENT '操作时间',
  
  -- 业务字段
  `biz_code` varchar(64) DEFAULT '' COMMENT '业务代码',
  `area` bigint(20) DEFAULT '0' COMMENT '业务分区',
  `source` varchar(64) DEFAULT '' COMMENT '来源',
  `metadata` varchar(1024) DEFAULT '' COMMENT '元数据（JSON格式）',
  `biz_source` varchar(64) DEFAULT '' COMMENT '业务来源',
  `platform` int(11) DEFAULT '0' COMMENT '平台：1=ios 2=android 3=pc 4=h5',
  
  `reserved1` bigint(20) DEFAULT '0' COMMENT '预留字段1',
  `version` bigint(20) DEFAULT '0' COMMENT '版本号',
  
  PRIMARY KEY (`id`),
  KEY `idx_uid` (`uid`),
  KEY `idx_tid` (`transaction_id`),
  KEY `idx_op_time` (`op_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='流水记录表（分表0）';

-- 创建其他9张分表
CREATE TABLE IF NOT EXISTS `coin_stream_record_1` LIKE `coin_stream_record_0`;
CREATE TABLE IF NOT EXISTS `coin_stream_record_2` LIKE `coin_stream_record_0`;
CREATE TABLE IF NOT EXISTS `coin_stream_record_3` LIKE `coin_stream_record_0`;
CREATE TABLE IF NOT EXISTS `coin_stream_record_4` LIKE `coin_stream_record_0`;
CREATE TABLE IF NOT EXISTS `coin_stream_record_5` LIKE `coin_stream_record_0`;
CREATE TABLE IF NOT EXISTS `coin_stream_record_6` LIKE `coin_stream_record_0`;
CREATE TABLE IF NOT EXISTS `coin_stream_record_7` LIKE `coin_stream_record_0`;
CREATE TABLE IF NOT EXISTS `coin_stream_record_8` LIKE `coin_stream_record_0`;
CREATE TABLE IF NOT EXISTS `coin_stream_record_9` LIKE `coin_stream_record_0`;

ALTER TABLE `coin_stream_record_1` COMMENT='流水记录表（分表1）';
ALTER TABLE `coin_stream_record_2` COMMENT='流水记录表（分表2）';
ALTER TABLE `coin_stream_record_3` COMMENT='流水记录表（分表3）';
ALTER TABLE `coin_stream_record_4` COMMENT='流水记录表（分表4）';
ALTER TABLE `coin_stream_record_5` COMMENT='流水记录表（分表5）';
ALTER TABLE `coin_stream_record_6` COMMENT='流水记录表（分表6）';
ALTER TABLE `coin_stream_record_7` COMMENT='流水记录表（分表7）';
ALTER TABLE `coin_stream_record_8` COMMENT='流水记录表（分表8）';
ALTER TABLE `coin_stream_record_9` COMMENT='流水记录表（分表9）';

-- =====================================================
-- 3. 兑换记录表
-- =====================================================

CREATE TABLE IF NOT EXISTS `coin_exchange_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL COMMENT '用户ID',
  `transaction_id` varchar(64) NOT NULL COMMENT '交易ID',
  `extend_tid` varchar(64) DEFAULT '' COMMENT '扩展交易ID',
  
  -- 源币种信息
  `src_coin_type` int(11) NOT NULL COMMENT '源币种类型：1=gold 2=iap_gold 3=silver',
  `src_coin_num` bigint(20) NOT NULL COMMENT '源币种数量',
  
  -- 目标币种信息
  `dest_coin_type` int(11) NOT NULL COMMENT '目标币种类型：1=gold 2=iap_gold 3=silver',
  `dest_coin_num` bigint(20) NOT NULL COMMENT '目标币种数量',
  
  `exchange_rate` decimal(10,4) DEFAULT '1.0000' COMMENT '兑换比例',
  `status` int(11) NOT NULL DEFAULT '1' COMMENT '状态：1=成功 0=失败',
  
  `ctime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  
  PRIMARY KEY (`id`),
  KEY `idx_uid` (`uid`),
  KEY `idx_tid` (`transaction_id`),
  KEY `idx_ctime` (`ctime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='兑换记录表';

-- =====================================================
-- 4. 插入测试数据
-- =====================================================

-- 插入测试用户钱包（uid=1001，分表1）
INSERT INTO `user_wallet_1` (`uid`, `gold`, `iap_gold`, `silver`, `gold_recharge_cnt`, `gold_pay_cnt`, `silver_pay_cnt`)
VALUES (1001, 1000, 0, 500, 1000, 0, 0)
ON DUPLICATE KEY UPDATE `uid`=`uid`;

-- 插入测试用户钱包（uid=1002，分表2）
INSERT INTO `user_wallet_2` (`uid`, `gold`, `iap_gold`, `silver`, `gold_recharge_cnt`, `gold_pay_cnt`, `silver_pay_cnt`)
VALUES (1002, 2000, 500, 1000, 2500, 0, 0)
ON DUPLICATE KEY UPDATE `uid`=`uid`;

-- 插入测试用户钱包（uid=1003，分表3）
INSERT INTO `user_wallet_3` (`uid`, `gold`, `iap_gold`, `silver`, `gold_recharge_cnt`, `gold_pay_cnt`, `silver_pay_cnt`)
VALUES (1003, 500, 0, 200, 500, 0, 0)
ON DUPLICATE KEY UPDATE `uid`=`uid`;

-- =====================================================
-- 完成
-- =====================================================

SELECT 'Wallet schema initialized successfully!' AS status;

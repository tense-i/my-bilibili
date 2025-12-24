-- ==================== Coupon 优惠券系统表结构 ====================
-- 完全参照 openbilibili 主项目

-- ==================== 批次信息表 ====================
CREATE TABLE IF NOT EXISTS `coupon_batch_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `app_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '应用ID',
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '批次名称',
  `batch_token` varchar(64) NOT NULL DEFAULT '' COMMENT '批次Token',
  `max_count` bigint(20) NOT NULL DEFAULT '-1' COMMENT '最大发放数量,-1不限制',
  `current_count` bigint(20) NOT NULL DEFAULT '0' COMMENT '当前发放数量',
  `start_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '开始时间',
  `expire_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '过期时间',
  `expire_day` bigint(20) NOT NULL DEFAULT '-1' COMMENT '有效天数,-1使用expire_time',
  `limit_count` bigint(20) NOT NULL DEFAULT '-1' COMMENT '每人限领数量,-1不限制',
  `full_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '满额条件',
  `amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '优惠金额',
  `state` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态:0正常,1冻结',
  `coupon_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '券类型:1观影券,2漫画券,3代金券',
  `platform_limit` varchar(64) NOT NULL DEFAULT '' COMMENT '平台限制,逗号分隔',
  `product_limit_month` tinyint(4) NOT NULL DEFAULT '0' COMMENT '商品月份限制:0不限,1月度,3季度,12年度',
  `product_limit_renewal` tinyint(4) NOT NULL DEFAULT '0' COMMENT '续费限制:0不限,1自动续期,2非自动续期',
  `ver` bigint(20) NOT NULL DEFAULT '0' COMMENT '版本号',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_batch_token` (`batch_token`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='优惠券批次信息表';


-- ==================== 优惠券订单表 ====================
CREATE TABLE IF NOT EXISTS `coupon_order` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `order_no` varchar(64) NOT NULL DEFAULT '' COMMENT '订单号',
  `mid` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `count` bigint(20) NOT NULL DEFAULT '0' COMMENT '使用数量',
  `state` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态:0待支付,1支付中,2支付成功,3支付失败',
  `coupon_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '券类型',
  `third_trade_no` varchar(64) NOT NULL DEFAULT '' COMMENT '第三方交易号',
  `remark` varchar(128) NOT NULL DEFAULT '' COMMENT '备注',
  `tips` varchar(128) NOT NULL DEFAULT '' COMMENT '提示',
  `use_ver` bigint(20) NOT NULL DEFAULT '0' COMMENT '使用版本',
  `ver` bigint(20) NOT NULL DEFAULT '0' COMMENT '版本号',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `ix_mid_state_type` (`mid`,`state`,`coupon_type`),
  KEY `ix_third_trade_no` (`third_trade_no`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='优惠券订单表';

-- ==================== 优惠券订单日志表 ====================
CREATE TABLE IF NOT EXISTS `coupon_order_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `order_no` varchar(64) NOT NULL DEFAULT '' COMMENT '订单号',
  `mid` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `state` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `ix_order_no` (`order_no`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='优惠券订单日志表';

-- ==================== 兑换码表 ====================
CREATE TABLE IF NOT EXISTS `coupon_code` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `batch_token` varchar(64) NOT NULL DEFAULT '' COMMENT '批次Token',
  `code` varchar(32) NOT NULL DEFAULT '' COMMENT '兑换码',
  `state` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态:1未使用,2已使用,3已冻结',
  `mid` bigint(20) NOT NULL DEFAULT '0' COMMENT '使用者ID',
  `coupon_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '券类型',
  `coupon_token` varchar(64) NOT NULL DEFAULT '' COMMENT '兑换后的券Token',
  `ver` bigint(20) NOT NULL DEFAULT '0' COMMENT '版本号',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `ix_batch_token` (`batch_token`),
  KEY `ix_mid` (`mid`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='兑换码表';

-- ==================== 领取日志表 ====================
CREATE TABLE IF NOT EXISTS `coupon_receive_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `appkey` varchar(32) NOT NULL DEFAULT '' COMMENT '应用Key',
  `order_no` varchar(64) NOT NULL DEFAULT '' COMMENT '订单号',
  `mid` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `coupon_token` varchar(64) NOT NULL DEFAULT '' COMMENT '券Token',
  `coupon_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '券类型',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_appkey_order_type` (`appkey`,`order_no`,`coupon_type`),
  KEY `ix_mid` (`mid`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='领取日志表';

-- ==================== 用户卡片表（新年活动）====================
CREATE TABLE IF NOT EXISTS `coupon_user_card` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `mid` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `card_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '卡片类型:0月度,1季度,2年度',
  `state` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态:0未开,1已开,2已使用',
  `batch_token` varchar(64) NOT NULL DEFAULT '' COMMENT '批次Token',
  `coupon_token` varchar(64) NOT NULL DEFAULT '' COMMENT '券Token',
  `act_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '活动ID',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_mid_act_card` (`mid`,`act_id`,`card_type`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户卡片表';

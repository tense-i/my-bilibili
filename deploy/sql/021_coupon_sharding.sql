-- ==================== Coupon 分表创建脚本 ====================
-- 观影券信息表：分100表 (coupon_info_00 ~ coupon_info_99)
-- 观影券变更日志表：分100表 (coupon_change_log_00 ~ coupon_change_log_99)
-- 代金券信息表：分10表 (coupon_allowance_info_00 ~ coupon_allowance_info_09)
-- 代金券变更日志表：分10表 (coupon_allowance_change_log_00 ~ coupon_allowance_change_log_09)
-- 漫画券余额表：分10表 (coupon_balance_info_00 ~ coupon_balance_info_09)
-- 漫画券余额变更日志表：分10表 (coupon_balance_change_log_00 ~ coupon_balance_change_log_09)

-- ==================== 观影券信息表模板 ====================
-- 分表键：mid % 100
DELIMITER //
DROP PROCEDURE IF EXISTS create_coupon_info_tables//
CREATE PROCEDURE create_coupon_info_tables()
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE table_name VARCHAR(64);
    DECLARE create_sql TEXT;
    
    WHILE i < 100 DO
        SET table_name = CONCAT('coupon_info_', LPAD(i, 2, '0'));
        SET create_sql = CONCAT('
            CREATE TABLE IF NOT EXISTS `', table_name, '` (
              `id` bigint(20) NOT NULL AUTO_INCREMENT,
              `coupon_token` varchar(64) NOT NULL DEFAULT '''' COMMENT ''券Token'',
              `mid` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''用户ID'',
              `state` tinyint(4) NOT NULL DEFAULT ''0'' COMMENT ''状态:0未使用,1使用中,2已使用,3已过期,4已冻结'',
              `start_time` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''开始时间'',
              `expire_time` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''过期时间'',
              `origin` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''来源'',
              `coupon_type` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''券类型'',
              `order_no` varchar(64) NOT NULL DEFAULT '''' COMMENT ''订单号'',
              `oid` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''关联ID'',
              `remark` varchar(128) NOT NULL DEFAULT '''' COMMENT ''备注'',
              `batch_token` varchar(64) NOT NULL DEFAULT '''' COMMENT ''批次Token'',
              `use_ver` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''使用版本'',
              `ver` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''版本号'',
              `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT ''创建时间'',
              `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT ''修改时间'',
              PRIMARY KEY (`id`),
              UNIQUE KEY `uk_coupon_token` (`coupon_token`),
              KEY `ix_mid_state_type` (`mid`,`state`,`coupon_type`),
              KEY `ix_order_no` (`order_no`),
              KEY `ix_mtime` (`mtime`)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT=''观影券信息表''
        ');
        
        SET @sql = create_sql;
        PREPARE stmt FROM @sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
        
        SET i = i + 1;
    END WHILE;
END//
DELIMITER ;

CALL create_coupon_info_tables();
DROP PROCEDURE IF EXISTS create_coupon_info_tables;


-- ==================== 观影券变更日志表 ====================
-- 分表键：mid % 100
DELIMITER //
DROP PROCEDURE IF EXISTS create_coupon_change_log_tables//
CREATE PROCEDURE create_coupon_change_log_tables()
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE table_name VARCHAR(64);
    DECLARE create_sql TEXT;
    
    WHILE i < 100 DO
        SET table_name = CONCAT('coupon_change_log_', LPAD(i, 2, '0'));
        SET create_sql = CONCAT('
            CREATE TABLE IF NOT EXISTS `', table_name, '` (
              `id` bigint(20) NOT NULL AUTO_INCREMENT,
              `coupon_token` varchar(64) NOT NULL DEFAULT '''' COMMENT ''券Token'',
              `mid` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''用户ID'',
              `state` tinyint(4) NOT NULL DEFAULT ''0'' COMMENT ''状态'',
              `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT ''创建时间'',
              `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT ''修改时间'',
              PRIMARY KEY (`id`),
              KEY `ix_coupon_token` (`coupon_token`),
              KEY `ix_mid` (`mid`),
              KEY `ix_mtime` (`mtime`)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT=''观影券变更日志表''
        ');
        
        SET @sql = create_sql;
        PREPARE stmt FROM @sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
        
        SET i = i + 1;
    END WHILE;
END//
DELIMITER ;

CALL create_coupon_change_log_tables();
DROP PROCEDURE IF EXISTS create_coupon_change_log_tables;

-- ==================== 代金券信息表 ====================
-- 分表键：mid % 10
DELIMITER //
DROP PROCEDURE IF EXISTS create_coupon_allowance_info_tables//
CREATE PROCEDURE create_coupon_allowance_info_tables()
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE table_name VARCHAR(64);
    DECLARE create_sql TEXT;
    
    WHILE i < 10 DO
        SET table_name = CONCAT('coupon_allowance_info_0', i);
        SET create_sql = CONCAT('
            CREATE TABLE IF NOT EXISTS `', table_name, '` (
              `id` bigint(20) NOT NULL AUTO_INCREMENT,
              `coupon_token` varchar(64) NOT NULL DEFAULT '''' COMMENT ''券Token'',
              `mid` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''用户ID'',
              `state` tinyint(4) NOT NULL DEFAULT ''0'' COMMENT ''状态:0未使用,1使用中,2已使用'',
              `start_time` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''开始时间'',
              `expire_time` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''过期时间'',
              `origin` tinyint(4) NOT NULL DEFAULT ''0'' COMMENT ''来源:1系统发放,2业务领取,3新年活动,4兑换码'',
              `ver` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''版本号'',
              `batch_token` varchar(64) NOT NULL DEFAULT '''' COMMENT ''批次Token'',
              `order_no` varchar(64) NOT NULL DEFAULT '''' COMMENT ''订单号'',
              `amount` decimal(10,2) NOT NULL DEFAULT ''0.00'' COMMENT ''优惠金额'',
              `full_amount` decimal(10,2) NOT NULL DEFAULT ''0.00'' COMMENT ''满额条件'',
              `app_id` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''应用ID'',
              `remark` varchar(128) NOT NULL DEFAULT '''' COMMENT ''备注'',
              `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT ''创建时间'',
              `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT ''修改时间'',
              PRIMARY KEY (`id`),
              UNIQUE KEY `uk_coupon_token` (`coupon_token`),
              KEY `ix_mid_state` (`mid`,`state`),
              KEY `ix_order_no` (`order_no`),
              KEY `ix_mid_batch` (`mid`,`batch_token`),
              KEY `ix_mtime` (`mtime`)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT=''代金券信息表''
        ');
        
        SET @sql = create_sql;
        PREPARE stmt FROM @sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
        
        SET i = i + 1;
    END WHILE;
END//
DELIMITER ;

CALL create_coupon_allowance_info_tables();
DROP PROCEDURE IF EXISTS create_coupon_allowance_info_tables;


-- ==================== 代金券变更日志表 ====================
-- 分表键：mid % 10
DELIMITER //
DROP PROCEDURE IF EXISTS create_coupon_allowance_change_log_tables//
CREATE PROCEDURE create_coupon_allowance_change_log_tables()
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE table_name VARCHAR(64);
    DECLARE create_sql TEXT;
    
    WHILE i < 10 DO
        SET table_name = CONCAT('coupon_allowance_change_log_0', i);
        SET create_sql = CONCAT('
            CREATE TABLE IF NOT EXISTS `', table_name, '` (
              `id` bigint(20) NOT NULL AUTO_INCREMENT,
              `coupon_token` varchar(64) NOT NULL DEFAULT '''' COMMENT ''券Token'',
              `order_no` varchar(64) NOT NULL DEFAULT '''' COMMENT ''订单号'',
              `mid` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''用户ID'',
              `state` tinyint(4) NOT NULL DEFAULT ''0'' COMMENT ''状态'',
              `change_type` tinyint(4) NOT NULL DEFAULT ''0'' COMMENT ''变更类型:1发放,2消费,3取消,4消费成功,5消费失败,6领取'',
              `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT ''创建时间'',
              `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT ''修改时间'',
              PRIMARY KEY (`id`),
              KEY `ix_coupon_token` (`coupon_token`),
              KEY `ix_mid` (`mid`),
              KEY `ix_mtime` (`mtime`)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT=''代金券变更日志表''
        ');
        
        SET @sql = create_sql;
        PREPARE stmt FROM @sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
        
        SET i = i + 1;
    END WHILE;
END//
DELIMITER ;

CALL create_coupon_allowance_change_log_tables();
DROP PROCEDURE IF EXISTS create_coupon_allowance_change_log_tables;

-- ==================== 漫画券余额表 ====================
-- 分表键：mid % 10
DELIMITER //
DROP PROCEDURE IF EXISTS create_coupon_balance_info_tables//
CREATE PROCEDURE create_coupon_balance_info_tables()
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE table_name VARCHAR(64);
    DECLARE create_sql TEXT;
    
    WHILE i < 10 DO
        SET table_name = CONCAT('coupon_balance_info_0', i);
        SET create_sql = CONCAT('
            CREATE TABLE IF NOT EXISTS `', table_name, '` (
              `id` bigint(20) NOT NULL AUTO_INCREMENT,
              `batch_token` varchar(64) NOT NULL DEFAULT '''' COMMENT ''批次Token'',
              `mid` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''用户ID'',
              `balance` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''余额'',
              `start_time` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''开始时间'',
              `expire_time` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''过期时间'',
              `origin` tinyint(4) NOT NULL DEFAULT ''0'' COMMENT ''来源'',
              `coupon_type` tinyint(4) NOT NULL DEFAULT ''0'' COMMENT ''券类型'',
              `ver` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''版本号'',
              `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT ''创建时间'',
              `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT ''修改时间'',
              PRIMARY KEY (`id`),
              UNIQUE KEY `uk_mid_batch` (`mid`,`batch_token`),
              KEY `ix_mid_type` (`mid`,`coupon_type`),
              KEY `ix_mtime` (`mtime`)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT=''漫画券余额表''
        ');
        
        SET @sql = create_sql;
        PREPARE stmt FROM @sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
        
        SET i = i + 1;
    END WHILE;
END//
DELIMITER ;

CALL create_coupon_balance_info_tables();
DROP PROCEDURE IF EXISTS create_coupon_balance_info_tables;

-- ==================== 漫画券余额变更日志表 ====================
-- 分表键：mid % 10
DELIMITER //
DROP PROCEDURE IF EXISTS create_coupon_balance_change_log_tables//
CREATE PROCEDURE create_coupon_balance_change_log_tables()
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE table_name VARCHAR(64);
    DECLARE create_sql TEXT;
    
    WHILE i < 10 DO
        SET table_name = CONCAT('coupon_balance_change_log_0', i);
        SET create_sql = CONCAT('
            CREATE TABLE IF NOT EXISTS `', table_name, '` (
              `id` bigint(20) NOT NULL AUTO_INCREMENT,
              `order_no` varchar(64) NOT NULL DEFAULT '''' COMMENT ''订单号'',
              `mid` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''用户ID'',
              `batch_token` varchar(64) NOT NULL DEFAULT '''' COMMENT ''批次Token'',
              `balance` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''变更后余额'',
              `change_balance` bigint(20) NOT NULL DEFAULT ''0'' COMMENT ''变更数量'',
              `change_type` tinyint(4) NOT NULL DEFAULT ''0'' COMMENT ''变更类型:1VIP发放,2系统发放,3消费,4消费失败退回'',
              `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT ''创建时间'',
              `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT ''修改时间'',
              PRIMARY KEY (`id`),
              KEY `ix_mid_batch` (`mid`,`batch_token`),
              KEY `ix_order_no` (`order_no`),
              KEY `ix_mtime` (`mtime`)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT=''漫画券余额变更日志表''
        ');
        
        SET @sql = create_sql;
        PREPARE stmt FROM @sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
        
        SET i = i + 1;
    END WHILE;
END//
DELIMITER ;

CALL create_coupon_balance_change_log_tables();
DROP PROCEDURE IF EXISTS create_coupon_balance_change_log_tables;

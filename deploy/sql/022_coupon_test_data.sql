-- ==================== Coupon 测试数据 ====================

-- 1. 批次信息测试数据
INSERT INTO `coupon_batch_info` (`app_id`, `name`, `batch_token`, `max_count`, `current_count`, `start_time`, `expire_time`, `expire_day`, `limit_count`, `full_amount`, `amount`, `state`, `coupon_type`, `platform_limit`, `product_limit_month`, `product_limit_renewal`, `ver`) VALUES
-- 代金券批次（满25减5）
(1, '大会员代金券-满25减5', 'batch_allowance_001', 10000, 100, UNIX_TIMESTAMP('2024-01-01'), UNIX_TIMESTAMP('2025-12-31'), -1, 3, 25.00, 5.00, 0, 3, '', 0, 0, 1),
-- 代金券批次（满50减10）
(1, '大会员代金券-满50减10', 'batch_allowance_002', 5000, 50, UNIX_TIMESTAMP('2024-01-01'), UNIX_TIMESTAMP('2025-12-31'), -1, 2, 50.00, 10.00, 0, 3, '3', 0, 0, 1),
-- 代金券批次（满100减20，仅PC）
(1, '大会员代金券-满100减20', 'batch_allowance_003', 2000, 20, UNIX_TIMESTAMP('2024-01-01'), UNIX_TIMESTAMP('2025-12-31'), -1, 1, 100.00, 20.00, 0, 3, '3', 1, 0, 1),
-- 观影券批次
(1, '观影券-月度会员', 'batch_video_001', -1, 1000, UNIX_TIMESTAMP('2024-01-01'), UNIX_TIMESTAMP('2025-12-31'), 30, -1, 0.00, 0.00, 0, 1, '', 0, 0, 1),
-- 漫画券批次
(1, '漫画券-新用户', 'batch_cartoon_001', 5000, 200, UNIX_TIMESTAMP('2024-01-01'), UNIX_TIMESTAMP('2025-12-31'), 7, 1, 0.00, 0.00, 0, 2, '', 0, 0, 1),
-- 冻结的批次
(1, '已冻结批次', 'batch_frozen_001', 1000, 0, UNIX_TIMESTAMP('2024-01-01'), UNIX_TIMESTAMP('2025-12-31'), -1, 1, 30.00, 5.00, 1, 3, '', 0, 0, 1);

-- 2. 代金券测试数据（用户 mid=1001，分表 coupon_allowance_info_01）
INSERT INTO `coupon_allowance_info_01` (`coupon_token`, `mid`, `state`, `start_time`, `expire_time`, `origin`, `ver`, `batch_token`, `order_no`, `amount`, `full_amount`, `app_id`, `remark`) VALUES
-- 未使用的代金券
('token_allowance_001', 1001, 0, UNIX_TIMESTAMP('2024-01-01'), UNIX_TIMESTAMP('2025-12-31'), 1, 1, 'batch_allowance_001', '', 5.00, 25.00, 1, '系统发放'),
('token_allowance_002', 1001, 0, UNIX_TIMESTAMP('2024-01-01'), UNIX_TIMESTAMP('2025-12-31'), 1, 1, 'batch_allowance_002', '', 10.00, 50.00, 1, '系统发放'),
('token_allowance_003', 1001, 0, UNIX_TIMESTAMP('2024-01-01'), UNIX_TIMESTAMP('2025-12-31'), 1, 1, 'batch_allowance_003', '', 20.00, 100.00, 1, '系统发放'),
-- 使用中的代金券
('token_allowance_004', 1001, 1, UNIX_TIMESTAMP('2024-01-01'), UNIX_TIMESTAMP('2025-12-31'), 1, 1, 'batch_allowance_001', 'order_001', 5.00, 25.00, 1, '使用中'),
-- 已使用的代金券
('token_allowance_005', 1001, 2, UNIX_TIMESTAMP('2024-01-01'), UNIX_TIMESTAMP('2025-06-30'), 1, 1, 'batch_allowance_001', 'order_002', 5.00, 25.00, 1, '已使用'),
-- 已过期的代金券
('token_allowance_006', 1001, 0, UNIX_TIMESTAMP('2024-01-01'), UNIX_TIMESTAMP('2024-06-30'), 1, 1, 'batch_allowance_001', '', 5.00, 25.00, 1, '已过期');

-- 3. 代金券测试数据（用户 mid=1002，分表 coupon_allowance_info_02）
INSERT INTO `coupon_allowance_info_02` (`coupon_token`, `mid`, `state`, `start_time`, `expire_time`, `origin`, `ver`, `batch_token`, `order_no`, `amount`, `full_amount`, `app_id`, `remark`) VALUES
('token_allowance_101', 1002, 0, UNIX_TIMESTAMP('2024-01-01'), UNIX_TIMESTAMP('2025-12-31'), 2, 1, 'batch_allowance_001', '', 5.00, 25.00, 1, '业务领取'),
('token_allowance_102', 1002, 0, UNIX_TIMESTAMP('2024-01-01'), UNIX_TIMESTAMP('2025-12-31'), 4, 1, 'batch_allowance_002', '', 10.00, 50.00, 1, '兑换码');

-- 4. 观影券测试数据（用户 mid=1001，分表 coupon_info_01）
INSERT INTO `coupon_info_01` (`coupon_token`, `mid`, `state`, `start_time`, `expire_time`, `origin`, `coupon_type`, `order_no`, `oid`, `remark`, `batch_token`, `use_ver`, `ver`) VALUES
('token_video_001', 1001, 0, UNIX_TIMESTAMP('2024-01-01'), UNIX_TIMESTAMP('2025-12-31'), 1, 1, '', 0, '观影券', 'batch_video_001', 0, 1),
('token_video_002', 1001, 0, UNIX_TIMESTAMP('2024-01-01'), UNIX_TIMESTAMP('2025-12-31'), 1, 1, '', 0, '观影券', 'batch_video_001', 0, 1),
('token_video_003', 1001, 2, UNIX_TIMESTAMP('2024-01-01'), UNIX_TIMESTAMP('2025-06-30'), 1, 1, 'order_video_001', 12345, '已使用', 'batch_video_001', 1, 2);

-- 5. 兑换码测试数据
INSERT INTO `coupon_code` (`batch_token`, `code`, `state`, `mid`, `coupon_type`, `coupon_token`, `ver`) VALUES
('batch_allowance_001', 'CODE001ABC', 1, 0, 3, '', 1),
('batch_allowance_001', 'CODE002DEF', 1, 0, 3, '', 1),
('batch_allowance_002', 'CODE003GHI', 1, 0, 3, '', 1),
('batch_allowance_001', 'CODE004JKL', 2, 1001, 3, 'token_allowance_code_001', 2),
('batch_allowance_001', 'CODE005MNO', 3, 0, 3, '', 1);

-- 6. 领取日志测试数据
INSERT INTO `coupon_receive_log` (`appkey`, `order_no`, `mid`, `coupon_token`, `coupon_type`) VALUES
('app_vip', 'receive_order_001', 1001, 'token_allowance_001', 3),
('app_vip', 'receive_order_002', 1001, 'token_allowance_002', 3),
('app_vip', 'receive_order_003', 1002, 'token_allowance_101', 3);


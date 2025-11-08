#!/bin/bash
# Redis 测试数据初始化脚本

REDIS_HOST="127.0.0.1"
REDIS_PORT="63790"
REDIS_PASS="redis123456"

echo "========================================"
echo "开始初始化 Redis 推荐系统测试数据"
echo "========================================"

# 1. 热门视频索引 (RECALL:HOT:INDEX, RECALL:HOT_DEFAULT:0)
echo "✓ 初始化热门视频索引..."
redis-cli -h $REDIS_HOST -p $REDIS_PORT -a $REDIS_PASS << EOF
DEL RECALL:HOT:INDEX
DEL RECALL:HOT_DEFAULT:0
ZADD RECALL:HOT:INDEX 98.5 100001
ZADD RECALL:HOT:INDEX 97.2 100011
ZADD RECALL:HOT:INDEX 96.8 100019
ZADD RECALL:HOT:INDEX 96.5 100003
ZADD RECALL:HOT:INDEX 96.0 100015
ZADD RECALL:HOT:INDEX 95.5 100009
ZADD RECALL:HOT:INDEX 94.8 100013
ZADD RECALL:HOT:INDEX 94.3 100006
ZADD RECALL:HOT:INDEX 93.8 100010
ZADD RECALL:HOT:INDEX 93.2 100016
ZADD RECALL:HOT:INDEX 92.5 100012
ZADD RECALL:HOT:INDEX 91.8 100018
ZADD RECALL:HOT:INDEX 91.2 100007
ZADD RECALL:HOT:INDEX 90.5 100014
ZADD RECALL:HOT:INDEX 90.0 100020
ZADD RECALL:HOT:INDEX 89.8 100004
ZADD RECALL:HOT:INDEX 89.2 100017
ZADD RECALL:HOT:INDEX 88.5 100008
ZADD RECALL:HOT:INDEX 87.8 100005
ZADD RECALL:HOT:INDEX 87.0 100002

# 复制一份到 HOT_DEFAULT
ZUNIONSTORE RECALL:HOT_DEFAULT:0 1 RECALL:HOT:INDEX
EOF

# 2. 精选视频索引 (recall:selection)
echo "初始化精选视频索引..."
redis-cli -h $REDIS_HOST -p $REDIS_PORT -a $REDIS_PASS << EOF
DEL recall:selection
LPUSH recall:selection 100001 100002 100003
EOF

# 3. I2I 相似视频索引
echo "初始化 I2I 相似视频索引..."
redis-cli -h $REDIS_HOST -p $REDIS_PORT -a $REDIS_PASS << EOF
DEL recall:i2i:100001
ZADD recall:i2i:100001 0.95 100002
ZADD recall:i2i:100001 0.89 100003
ZADD recall:i2i:100001 0.85 100004

DEL recall:i2i:100002
ZADD recall:i2i:100002 0.95 100001
ZADD recall:i2i:100002 0.88 100005
ZADD recall:i2i:100002 0.82 100006

DEL recall:i2i:100003
ZADD recall:i2i:100003 0.89 100001
ZADD recall:i2i:100003 0.87 100007
ZADD recall:i2i:100003 0.84 100008
EOF

# 4. 标签热门视频索引
echo "初始化标签视频索引..."
redis-cli -h $REDIS_HOST -p $REDIS_PORT -a $REDIS_PASS << EOF
# MMD 标签
DEL recall:tag:1
ZADD recall:tag:1 95.0 100001
ZADD recall:tag:1 94.0 100002
ZADD recall:tag:1 93.0 100004

# 初音未来 标签
DEL recall:tag:2
ZADD recall:tag:2 96.0 100001
ZADD recall:tag:2 92.0 100005

# 洛天依 标签
DEL recall:tag:4
ZADD recall:tag:4 95.0 100002
ZADD recall:tag:4 91.0 100006

# 我的世界 标签
DEL recall:tag:5
ZADD recall:tag:5 97.0 100003
ZADD recall:tag:5 93.0 100007
ZADD recall:tag:5 90.0 100008

# 建筑 标签
DEL recall:tag:6
ZADD recall:tag:6 96.0 100003
ZADD recall:tag:6 92.0 100009

# 动画 标签
DEL recall:tag:动画
ZADD recall:tag:动画 95.0 100001
ZADD recall:tag:动画 94.0 100002
ZADD recall:tag:动画 90.0 100004

# 游戏 标签
DEL recall:tag:游戏
ZADD recall:tag:游戏 97.0 100003
ZADD recall:tag:游戏 93.0 100007
ZADD recall:tag:游戏 91.0 100008
EOF

# 5. UP主视频索引
echo "初始化 UP 主视频索引..."
redis-cli -h $REDIS_HOST -p $REDIS_PORT -a $REDIS_PASS << EOF
# UP主 1001
DEL recall:up:1001
ZADD recall:up:1001 $(date +%s) 100001
ZADD recall:up:1001 $(($(date +%s) - 86400)) 100002
ZADD recall:up:1001 $(($(date +%s) - 172800)) 100004

# UP主 1002
DEL recall:up:1002
ZADD recall:up:1002 $(date +%s) 100003
ZADD recall:up:1002 $(($(date +%s) - 86400)) 100007
ZADD recall:up:1002 $(($(date +%s) - 172800)) 100008

# UP主 1003
DEL recall:up:1003
ZADD recall:up:1003 $(date +%s) 100005
ZADD recall:up:1003 $(($(date +%s) - 86400)) 100006
ZADD recall:up:1003 $(($(date +%s) - 172800)) 100009
EOF

# 6. 用户画像缓存（示例）
echo "初始化用户画像缓存..."
redis-cli -h $REDIS_HOST -p $REDIS_PORT -a $REDIS_PASS << EOF
# 用户 1001 的画像
DEL user:profile:1001
HSET user:profile:1001 tags '{"MMD":0.9,"VOCALOID":0.8,"动画":0.7}'
HSET user:profile:1001 zones '{"动画":0.9,"音乐":0.7}'
HSET user:profile:1001 last_update $(date +%s)
EXPIRE user:profile:1001 3600

# 用户 1002 的画像
DEL user:profile:1002
HSET user:profile:1002 tags '{"游戏":0.9,"我的世界":0.8,"建筑":0.7}'
HSET user:profile:1002 zones '{"游戏":0.9,"生活":0.6}'
HSET user:profile:1002 last_update $(date +%s)
EXPIRE user:profile:1002 3600
EOF

# 7. 用户行为数据（点赞、正反馈）
echo "初始化用户行为数据..."
DATE=$(date +%Y%m%d)
redis-cli -h $REDIS_HOST -p $REDIS_PORT -a $REDIS_PASS << EOF
# 用户 1001 点赞记录
DEL user:action:1001:like:$DATE
ZADD user:action:1001:like:$DATE $(date +%s) 100001
ZADD user:action:1001:like:$DATE $(($(date +%s) - 3600)) 100002
EXPIRE user:action:1001:like:$DATE 604800

# 用户 1001 正反馈记录
DEL user:action:1001:pos:$DATE
ZADD user:action:1001:pos:$DATE $(date +%s) 100001
ZADD user:action:1001:pos:$DATE $(($(date +%s) - 3600)) 100002
ZADD user:action:1001:pos:$DATE $(($(date +%s) - 7200)) 100003
EXPIRE user:action:1001:pos:$DATE 604800

# 用户 1002 点赞记录
DEL user:action:1002:like:$DATE
ZADD user:action:1002:like:$DATE $(date +%s) 100003
ZADD user:action:1002:like:$DATE $(($(date +%s) - 3600)) 100007
EXPIRE user:action:1002:like:$DATE 604800

# 用户 1002 正反馈记录
DEL user:action:1002:pos:$DATE
ZADD user:action:1002:pos:$DATE $(date +%s) 100003
ZADD user:action:1002:pos:$DATE $(($(date +%s) - 3600)) 100007
ZADD user:action:1002:pos:$DATE $(($(date +%s) - 7200)) 100008
EXPIRE user:action:1002:pos:$DATE 604800
EOF

echo "✅ Redis 测试数据初始化完成！"
echo ""
echo "数据统计："
redis-cli -h $REDIS_HOST -p $REDIS_PORT -a $REDIS_PASS << EOF
ECHO "热门视频数量:"
ZCARD recall:hot:default
ECHO "精选视频数量:"
LLEN recall:selection
ECHO "标签索引数量:"
KEYS recall:tag:*
ECHO "UP主索引数量:"
KEYS recall:up:*
EOF


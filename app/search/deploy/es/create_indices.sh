#!/bin/bash

# Elasticsearch 索引创建脚本
# 用法: ./create_indices.sh [ES_HOST]
# 默认 ES_HOST: http://localhost:9200

ES_HOST=${1:-"http://localhost:9200"}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=========================================="
echo "Elasticsearch 索引创建脚本"
echo "ES Host: $ES_HOST"
echo "=========================================="

# 检查 ES 是否可用
echo "检查 Elasticsearch 连接..."
if ! curl -s "$ES_HOST" > /dev/null; then
    echo "错误: 无法连接到 Elasticsearch ($ES_HOST)"
    exit 1
fi
echo "Elasticsearch 连接成功!"
echo ""

# 创建 PGC Media 索引
echo "创建 pgc_media 索引..."
curl -X PUT "$ES_HOST/pgc_media" \
    -H 'Content-Type: application/json' \
    -d @"$SCRIPT_DIR/pgc_media_mapping.json" \
    2>/dev/null
echo ""

# 创建弹幕搜索索引 (1000个分片: dm_search_000 到 dm_search_999)
echo "创建弹幕搜索索引 (dm_search_000 - dm_search_999)..."
for i in $(seq 0 999); do
    index_name=$(printf "dm_search_%03d" $i)
    curl -X PUT "$ES_HOST/$index_name" \
        -H 'Content-Type: application/json' \
        -d @"$SCRIPT_DIR/dm_search_mapping.json" \
        2>/dev/null
    
    # 每100个打印进度
    if [ $((i % 100)) -eq 0 ]; then
        echo "  已创建: $index_name"
    fi
done
echo "弹幕搜索索引创建完成!"
echo ""

# 创建评论记录索引 (100个分片: replyrecord_00 到 replyrecord_99)
echo "创建评论记录索引 (replyrecord_00 - replyrecord_99)..."
for i in $(seq 0 99); do
    index_name=$(printf "replyrecord_%02d" $i)
    curl -X PUT "$ES_HOST/$index_name" \
        -H 'Content-Type: application/json' \
        -d @"$SCRIPT_DIR/replyrecord_mapping.json" \
        2>/dev/null
    
    # 每10个打印进度
    if [ $((i % 10)) -eq 0 ]; then
        echo "  已创建: $index_name"
    fi
done
echo "评论记录索引创建完成!"
echo ""

# 创建当前月份的弹幕日期索引
current_month=$(date +"%Y_%m")
echo "创建当前月份弹幕日期索引 (dm_date_$current_month)..."
curl -X PUT "$ES_HOST/dm_date_$current_month" \
    -H 'Content-Type: application/json' \
    -d @"$SCRIPT_DIR/dm_date_mapping.json" \
    2>/dev/null
echo ""

echo "=========================================="
echo "所有索引创建完成!"
echo "=========================================="

# 显示索引列表
echo ""
echo "当前索引列表:"
curl -s "$ES_HOST/_cat/indices?v" | head -20
echo "..."

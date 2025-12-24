#!/bin/bash

# Elasticsearch 索引创建脚本 (开发环境精简版)
# 只创建少量索引用于开发测试
# 用法: ./create_indices_dev.sh [ES_HOST]

ES_HOST=${1:-"http://localhost:9200"}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=========================================="
echo "Elasticsearch 索引创建脚本 (开发环境)"
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
    -d @"$SCRIPT_DIR/pgc_media_mapping.json"
echo ""

# 创建弹幕搜索索引 (开发环境只创建10个: dm_search_000 到 dm_search_009)
echo "创建弹幕搜索索引 (dm_search_000 - dm_search_009)..."
for i in $(seq 0 9); do
    index_name=$(printf "dm_search_%03d" $i)
    echo "  创建: $index_name"
    curl -s -X PUT "$ES_HOST/$index_name" \
        -H 'Content-Type: application/json' \
        -d @"$SCRIPT_DIR/dm_search_mapping.json" > /dev/null
done
echo ""

# 创建评论记录索引 (开发环境只创建10个: replyrecord_00 到 replyrecord_09)
echo "创建评论记录索引 (replyrecord_00 - replyrecord_09)..."
for i in $(seq 0 9); do
    index_name=$(printf "replyrecord_%02d" $i)
    echo "  创建: $index_name"
    curl -s -X PUT "$ES_HOST/$index_name" \
        -H 'Content-Type: application/json' \
        -d @"$SCRIPT_DIR/replyrecord_mapping.json" > /dev/null
done
echo ""

# 创建当前月份的弹幕日期索引
current_month=$(date +"%Y_%m")
echo "创建弹幕日期索引 (dm_date_$current_month)..."
curl -X PUT "$ES_HOST/dm_date_$current_month" \
    -H 'Content-Type: application/json' \
    -d @"$SCRIPT_DIR/dm_date_mapping.json"
echo ""

echo "=========================================="
echo "开发环境索引创建完成!"
echo "=========================================="

# 显示索引列表
echo ""
echo "当前索引列表:"
curl -s "$ES_HOST/_cat/indices?v"

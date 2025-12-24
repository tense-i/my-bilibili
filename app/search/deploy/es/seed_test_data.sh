#!/bin/bash
# seed_test_data.sh - 为搜索服务的 ES 索引写入测试数据
# 用法: ./seed_test_data.sh [ES_HOST]
# 示例: ./seed_test_data.sh http://localhost:9200

set -e

ES_HOST="${1:-http://localhost:9200}"

echo "=========================================="
echo "ES 测试数据写入脚本"
echo "ES Host: $ES_HOST"
echo "=========================================="

# 检查 ES 是否可用
check_es() {
    echo "检查 ES 连接..."
    if ! curl -s "$ES_HOST" > /dev/null; then
        echo "错误: 无法连接到 ES: $ES_HOST"
        exit 1
    fi
    echo "ES 连接正常"
}

# 弹幕内容模板
DM_CONTENTS=(
    "前方高能预警"
    "太好看了"
    "泪目了"
    "哈哈哈哈哈"
    "这也太强了吧"
    "awsl"
    "名场面"
    "经典永流传"
    "爷青回"
    "下次一定"
    "我来组成头部"
    "开幕雷击"
    "标准结局"
    "这波啊这波是"
    "有内味了"
    "好家伙"
    "绝绝子"
    "yyds"
    "破防了"
    "DNA动了"
)

# PGC 番剧标题模板
PGC_TITLES=(
    "进击的巨人 最终季"
    "鬼灭之刃 无限列车篇"
    "咒术回战"
    "间谍过家家"
    "电锯人"
    "我的英雄学院 第六季"
    "JOJO的奇妙冒险 石之海"
    "辉夜大小姐想让我告白"
    "孤独摇滚"
    "葬送的芙莉莲"
    "药屋少女的呢喃"
    "迷宫饭"
    "怪兽8号"
    "排球少年 垃圾场的决战"
    "蓝色监狱"
    "天国大魔境"
    "地狱乐"
    "我推的孩子"
    "无职转生"
    "86-不存在的战区-"
)

# 评论内容模板
REPLY_CONTENTS=(
    "这个视频太棒了！"
    "UP主辛苦了"
    "三连支持"
    "学到了"
    "感谢分享"
    "催更催更"
    "这期质量很高"
    "看完了，意犹未尽"
    "建议收藏"
    "已投币"
    "前排占座"
    "来晚了"
    "每日打卡"
    "这就是我想看的"
    "太强了"
    "涨知识了"
    "期待下期"
    "收藏从未停止"
    "这波操作666"
    "爱了爱了"
)

# 生成随机时间戳 (最近30天内)
random_timestamp() {
    local now=$(date +%s)
    local days_ago=$((RANDOM % 30))
    local hours=$((RANDOM % 24))
    local minutes=$((RANDOM % 60))
    local seconds=$((RANDOM % 60))
    echo $((now - days_ago * 86400 - hours * 3600 - minutes * 60 - seconds))
}

# 格式化时间戳为日期字符串
format_datetime() {
    date -r "$1" "+%Y-%m-%d %H:%M:%S" 2>/dev/null || date -d "@$1" "+%Y-%m-%d %H:%M:%S"
}

# 写入弹幕搜索测试数据 (dm_search_000)
seed_dm_search() {
    echo ""
    echo "=========================================="
    echo "写入弹幕搜索测试数据 (dm_search_000)"
    echo "=========================================="
    
    local index_name="dm_search_000"
    local oid=1000  # oid % 1000 = 0，所以写入 dm_search_000
    
    for i in $(seq 1 40); do
        local id=$((10000 + i))
        local mid=$((1000 + RANDOM % 100))
        local content="${DM_CONTENTS[$((RANDOM % ${#DM_CONTENTS[@]}))]}"
        local mode=$((RANDOM % 3 + 1))  # 1-3
        local pool=$((RANDOM % 3))       # 0-2
        local progress=$((RANDOM % 300000))  # 0-300秒
        local state=0
        local type=1
        local fontsize=25
        local color=$((RANDOM % 16777216))
        local ts=$(random_timestamp)
        local ctime=$(format_datetime $ts)
        
        local doc=$(cat <<EOF
{
    "id": $id,
    "oid": $oid,
    "oidstr": "$oid",
    "mid": $mid,
    "content": "$content",
    "mode": $mode,
    "pool": $pool,
    "progress": $progress,
    "state": $state,
    "type": $type,
    "attr": 0,
    "attr_format": 0,
    "fontsize": $fontsize,
    "color": $color,
    "ctime": "$ctime",
    "mtime": "$ctime"
}
EOF
)
        
        curl -s -X POST "$ES_HOST/$index_name/_doc/$id" \
            -H "Content-Type: application/json" \
            -d "$doc" > /dev/null
        
        echo "  写入弹幕 $i/40: id=$id, content=$content"
    done
    
    echo "弹幕搜索数据写入完成"
}

# 写入弹幕日期测试数据 (dm_date_2024_12)
seed_dm_date() {
    echo ""
    echo "=========================================="
    echo "写入弹幕日期测试数据 (dm_date_2024_12)"
    echo "=========================================="
    
    local index_name="dm_date_2024_12"
    
    for i in $(seq 1 40); do
        local id=$((20000 + i))
        local oid=$((1000 + i))
        local day=$((i % 31 + 1))
        local day_str=$(printf "%02d" $day)
        local month="2024-12"
        local date="2024-12-$day_str"
        local total=$((RANDOM % 10000 + 100))
        local ctime="2024-12-$day_str 00:00:00"
        
        local doc=$(cat <<EOF
{
    "id": $id,
    "oid": $oid,
    "month": "$month",
    "date": "$date",
    "total": $total,
    "ctime": "$ctime"
}
EOF
)
        
        curl -s -X POST "$ES_HOST/$index_name/_doc/$id" \
            -H "Content-Type: application/json" \
            -d "$doc" > /dev/null
        
        echo "  写入弹幕日期 $i/40: oid=$oid, date=$date, total=$total"
    done
    
    echo "弹幕日期数据写入完成"
}

# 写入 PGC 番剧测试数据 (pgc_media)
seed_pgc_media() {
    echo ""
    echo "=========================================="
    echo "写入 PGC 番剧测试数据 (pgc_media)"
    echo "=========================================="
    
    local index_name="pgc_media"
    
    for i in $(seq 1 40); do
        local media_id=$((28220000 + i))
        local season_id=$((39000 + i))
        local title="${PGC_TITLES[$((i % ${#PGC_TITLES[@]}))]}"
        local season_type=$((RANDOM % 4 + 1))  # 1-4
        local style_id=$((RANDOM % 100 + 1))
        local status=$((RANDOM % 3))  # 0-2
        local producer_id=$((RANDOM % 1000 + 1))
        local area_id=$((RANDOM % 5 + 1))
        local score=$(echo "scale=1; $((RANDOM % 30 + 70)) / 10" | bc)
        local season_version=$((RANDOM % 3 + 1))
        local season_status=$((RANDOM % 5))
        local season_month=$((RANDOM % 12 + 1))
        local dm_count=$((RANDOM % 5000000 + 100000))
        local play_count=$((RANDOM % 100000000 + 1000000))
        local fav_count=$((RANDOM % 5000000 + 100000))
        local ts=$(random_timestamp)
        local ctime=$(format_datetime $ts)
        local release_date="2024-$((RANDOM % 12 + 1))-$((RANDOM % 28 + 1))"
        
        local doc=$(cat <<EOF
{
    "media_id": $media_id,
    "season_id": $season_id,
    "title": "$title",
    "season_type": $season_type,
    "style_id": $style_id,
    "status": $status,
    "release_date": "$release_date",
    "producer_id": $producer_id,
    "is_deleted": 0,
    "area_id": "$area_id",
    "score": $score,
    "is_finish": "1",
    "season_version": $season_version,
    "season_status": $season_status,
    "pub_time": "$ctime",
    "season_month": $season_month,
    "latest_time": "$ctime",
    "copyright_info": "bilibili",
    "dm_count": $dm_count,
    "play_count": $play_count,
    "fav_count": $fav_count,
    "ctime": "$ctime",
    "mtime": "$ctime"
}
EOF
)
        
        curl -s -X POST "$ES_HOST/$index_name/_doc/$media_id" \
            -H "Content-Type: application/json" \
            -d "$doc" > /dev/null
        
        echo "  写入番剧 $i/40: media_id=$media_id, title=$title"
    done
    
    echo "PGC 番剧数据写入完成"
}

# 写入评论记录测试数据 (replyrecord_00)
seed_replyrecord() {
    echo ""
    echo "=========================================="
    echo "写入评论记录测试数据 (replyrecord_00)"
    echo "=========================================="
    
    local index_name="replyrecord_00"
    local mid=100  # mid % 100 = 0，所以写入 replyrecord_00
    
    for i in $(seq 1 40); do
        local id=$((30000 + i))
        local oid=$((1000 + RANDOM % 1000))
        local type=$((RANDOM % 3 + 1))  # 1-3
        local state=0
        local content="${REPLY_CONTENTS[$((RANDOM % ${#REPLY_CONTENTS[@]}))]}"
        local like=$((RANDOM % 10000))
        local hate=$((RANDOM % 100))
        local rcount=$((RANDOM % 50))
        local floor=$((i))
        local ts=$(random_timestamp)
        local ctime=$(format_datetime $ts)
        
        local doc=$(cat <<EOF
{
    "id": $id,
    "oid": $oid,
    "mid": $mid,
    "type": $type,
    "state": $state,
    "content": "$content",
    "like": $like,
    "hate": $hate,
    "rcount": $rcount,
    "floor": $floor,
    "ctime": "$ctime",
    "mtime": "$ctime"
}
EOF
)
        
        local doc_id="${id}_${oid}"
        
        curl -s -X POST "$ES_HOST/$index_name/_doc/$doc_id" \
            -H "Content-Type: application/json" \
            -d "$doc" > /dev/null
        
        echo "  写入评论 $i/40: id=$id, content=$content"
    done
    
    echo "评论记录数据写入完成"
}

# 刷新索引
refresh_indices() {
    echo ""
    echo "=========================================="
    echo "刷新索引"
    echo "=========================================="
    
    curl -s -X POST "$ES_HOST/dm_search_000/_refresh" > /dev/null
    curl -s -X POST "$ES_HOST/dm_date_2024_12/_refresh" > /dev/null
    curl -s -X POST "$ES_HOST/pgc_media/_refresh" > /dev/null
    curl -s -X POST "$ES_HOST/replyrecord_00/_refresh" > /dev/null
    
    echo "索引刷新完成"
}

# 验证数据
verify_data() {
    echo ""
    echo "=========================================="
    echo "验证数据"
    echo "=========================================="
    
    echo ""
    echo "dm_search_000 文档数:"
    curl -s "$ES_HOST/dm_search_000/_count" | grep -o '"count":[0-9]*' | cut -d: -f2
    
    echo ""
    echo "dm_date_2024_12 文档数:"
    curl -s "$ES_HOST/dm_date_2024_12/_count" | grep -o '"count":[0-9]*' | cut -d: -f2
    
    echo ""
    echo "pgc_media 文档数:"
    curl -s "$ES_HOST/pgc_media/_count" | grep -o '"count":[0-9]*' | cut -d: -f2
    
    echo ""
    echo "replyrecord_00 文档数:"
    curl -s "$ES_HOST/replyrecord_00/_count" | grep -o '"count":[0-9]*' | cut -d: -f2
}

# 主函数
main() {
    check_es
    seed_dm_search
    seed_dm_date
    seed_pgc_media
    seed_replyrecord
    refresh_indices
    verify_data
    
    echo ""
    echo "=========================================="
    echo "所有测试数据写入完成！"
    echo "=========================================="
    echo ""
    echo "测试搜索示例:"
    echo ""
    echo "1. 弹幕搜索:"
    echo "   curl '$ES_HOST/dm_search_000/_search?q=content:高能'"
    echo ""
    echo "2. 弹幕日期搜索:"
    echo "   curl '$ES_HOST/dm_date_2024_12/_search?q=month:2024-12'"
    echo ""
    echo "3. PGC 番剧搜索:"
    echo "   curl '$ES_HOST/pgc_media/_search?q=title:进击'"
    echo ""
    echo "4. 评论记录搜索:"
    echo "   curl '$ES_HOST/replyrecord_00/_search?q=mid:100'"
}

main

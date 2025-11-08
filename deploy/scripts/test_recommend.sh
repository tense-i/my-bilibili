#!/bin/bash
# 推荐系统测试脚本

echo "========================================"
echo "推荐系统功能测试"
echo "========================================"

# 测试用户
TEST_USERS=(1001 1002 1003 1004 1005)

echo ""
echo "=== 1. 测试 Recall RPC (召回服务) ==="
echo "端口: 9006"
echo ""

# 使用 grpcurl 测试召回服务
for mid in "${TEST_USERS[@]}"; do
    echo "测试用户 $mid 的召回结果:"
    grpcurl -plaintext -d "{\"mid\": $mid, \"total_limit\": 20, \"infos\": [{\"name\": \"HotRecall\", \"tag\": \"recall:hot:default\", \"limit\": 10, \"priority\": 10}]}" \
        127.0.0.1:9006 recall.Recall/Recall 2>/dev/null | jq -r '.list[:5] | .[] | "  - AVID: \(.avid), 分数: \(.score), 召回类型: \(.recall_type)"'
    echo ""
done

echo ""
echo "=== 2. 测试 Recommend RPC (推荐服务) ==="
echo "端口: 9005"
echo ""

# 使用 grpcurl 测试推荐服务
for mid in "${TEST_USERS[@]}"; do
    echo "测试用户 $mid 的推荐结果:"
    grpcurl -plaintext -d "{\"mid\": $mid, \"limit\": 5}" \
        127.0.0.1:9005 recommend.Recommend/GetRecommendList 2>/dev/null | jq -r '.list[:5] | .[] | "  - AVID: \(.avid), 标题: \(.title), 分数: \(.score), 理由: \(.reason)"'
    echo ""
done

echo "========================================"
echo "测试完成！"
echo "========================================"


#!/bin/bash

# 充值功能测试脚本

echo "========================================="
echo "  虚拟钱包系统 - 充值功能测试"
echo "========================================="
echo ""

# 生成唯一transaction_id
TID="tid_$(date +%s)_$$"

echo "测试参数："
echo "  用户ID: 1001"
echo "  币种类型: gold"
echo "  充值金额: 100"
echo "  TransactionID: $TID"
echo "  平台: android"
echo ""

echo "发送充值请求..."
response=$(curl -s -X POST http://localhost:8004/api/wallet/v1/recharge \
  -H "Content-Type: application/json" \
  -d "{
    \"uid\": 1001,
    \"coin_type\": \"gold\",
    \"coin_num\": 100,
    \"transaction_id\": \"$TID\",
    \"platform\": \"android\"
  }")

echo "响应结果："
echo "$response" | python3 -m json.tool 2>/dev/null || echo "$response"
echo ""

# 查询最新余额
echo "查询当前余额..."
detail=$(curl -s "http://localhost:8004/api/wallet/v1/detail?uid=1001&platform=android")
echo "$detail" | python3 -m json.tool 2>/dev/null || echo "$detail"
echo ""

echo "========================================="
echo "测试完成！"
echo "========================================="

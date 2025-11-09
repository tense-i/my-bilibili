#!/bin/bash

# 虚拟钱包系统 - 全功能测试脚本

set -e

API_URL="http://127.0.0.1:8004/api/wallet/v1"
TEST_UID=1002
TIMESTAMP=$(date +%s%N)

echo "========================================="
echo "  虚拟钱包系统 - 全功能测试"
echo "========================================="
echo ""
echo "测试用户: UID=$TEST_UID"
echo "时间戳: $TIMESTAMP"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试计数器
PASS_COUNT=0
FAIL_COUNT=0

# 测试函数
test_api() {
    local test_name=$1
    local method=$2
    local url=$3
    local data=$4
    local expected=$5
    
    echo -n "测试: $test_name ... "
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s "$url" --noproxy '*')
    else
        response=$(curl -s -X POST "$url" -H "Content-Type: application/json" -d "$data" --noproxy '*')
    fi
    
    if echo "$response" | grep -q "$expected"; then
        echo -e "${GREEN}✓ PASS${NC}"
        ((PASS_COUNT++))
        return 0
    else
        echo -e "${RED}✗ FAIL${NC}"
        echo "  响应: $response"
        ((FAIL_COUNT++))
        return 1
    fi
}

echo "========================================="
echo "阶段一：初始状态查询"
echo "========================================="

# 1. 查询初始余额
test_api "查询初始余额" "GET" "$API_URL/detail?uid=$TEST_UID&platform=android" "" "uid"

echo ""
echo "========================================="
echo "阶段二：充值功能测试"
echo "========================================="

# 2. 充值100金瓜子
TID_RECHARGE1="tid_recharge_${TIMESTAMP}_1"
test_api "充值100金瓜子" "POST" "$API_URL/recharge" \
    "{\"uid\":$TEST_UID,\"coin_type\":\"gold\",\"coin_num\":100,\"transaction_id\":\"$TID_RECHARGE1\",\"platform\":\"android\"}" \
    "gold"

# 3. 充值200金瓜子
TID_RECHARGE2="tid_recharge_${TIMESTAMP}_2"
test_api "充值200金瓜子" "POST" "$API_URL/recharge" \
    "{\"uid\":$TEST_UID,\"coin_type\":\"gold\",\"coin_num\":200,\"transaction_id\":\"$TID_RECHARGE2\",\"platform\":\"android\"}" \
    "gold"

# 4. 重复充值（应该失败）
sleep 1
test_api "重复充值（防重）" "POST" "$API_URL/recharge" \
    "{\"uid\":$TEST_UID,\"coin_type\":\"gold\",\"coin_num\":100,\"transaction_id\":\"$TID_RECHARGE1\",\"platform\":\"android\"}" \
    "处理中"

echo ""
echo "========================================="
echo "阶段三：消费功能测试"
echo "========================================="

# 5. 消费50金瓜子
TID_PAY1="tid_pay_${TIMESTAMP}_1"
test_api "消费50金瓜子" "POST" "$API_URL/pay" \
    "{\"uid\":$TEST_UID,\"coin_type\":\"gold\",\"coin_num\":50,\"transaction_id\":\"$TID_PAY1\",\"platform\":\"android\"}" \
    "gold"

# 6. 消费超额（应该失败）
TID_PAY2="tid_pay_${TIMESTAMP}_2"
test_api "消费超额（防超支）" "POST" "$API_URL/pay" \
    "{\"uid\":$TEST_UID,\"coin_type\":\"gold\",\"coin_num\":99999,\"transaction_id\":\"$TID_PAY2\",\"platform\":\"android\"}" \
    "余额不足"

echo ""
echo "========================================="
echo "阶段四：兑换功能测试"
echo "========================================="

# 7. 兑换100金瓜子->银瓜子
TID_EXCHANGE1="tid_exchange_${TIMESTAMP}_1"
test_api "兑换100金瓜子→银瓜子" "POST" "$API_URL/exchange" \
    "{\"uid\":$TEST_UID,\"src_coin_type\":\"gold\",\"src_coin_num\":100,\"dest_coin_type\":\"silver\",\"dest_coin_num\":100,\"transaction_id\":\"$TID_EXCHANGE1\",\"platform\":\"android\"}" \
    "silver"

# 8. 错误兑换比例（应该失败）
TID_EXCHANGE2="tid_exchange_${TIMESTAMP}_2"
test_api "错误兑换比例" "POST" "$API_URL/exchange" \
    "{\"uid\":$TEST_UID,\"src_coin_type\":\"gold\",\"src_coin_num\":100,\"dest_coin_type\":\"silver\",\"dest_coin_num\":200,\"transaction_id\":\"$TID_EXCHANGE2\",\"platform\":\"android\"}" \
    "比例错误"

echo ""
echo "========================================="
echo "阶段五：查询功能测试"
echo "========================================="

# 9. 查询最终余额
test_api "查询最终余额" "GET" "$API_URL/detail?uid=$TEST_UID&platform=android" "" "gold"

# 10. 查询流水列表
test_api "查询流水列表" "GET" "$API_URL/stream?uid=$TEST_UID&offset=0&limit=20" "" "list"

echo ""
echo "========================================="
echo "数据一致性验证"
echo "========================================="

echo "查询数据库余额..."
docker exec mybilibili-mysql mysql -uroot -proot123456 mybilibili \
    -e "SELECT uid, gold, iap_gold, silver, gold_recharge_cnt, gold_pay_cnt FROM user_wallet_2 WHERE uid=$TEST_UID" 2>&1 | grep -v Warning || echo "用户尚未初始化"

echo ""
echo "查询流水记录数..."
docker exec mybilibili-mysql mysql -uroot -proot123456 mybilibili \
    -e "SELECT COUNT(*) as total FROM coin_stream_record_0 WHERE uid=$TEST_UID 
        UNION ALL SELECT COUNT(*) FROM coin_stream_record_1 WHERE uid=$TEST_UID
        UNION ALL SELECT COUNT(*) FROM coin_stream_record_2 WHERE uid=$TEST_UID
        UNION ALL SELECT COUNT(*) FROM coin_stream_record_3 WHERE uid=$TEST_UID
        UNION ALL SELECT COUNT(*) FROM coin_stream_record_4 WHERE uid=$TEST_UID
        UNION ALL SELECT COUNT(*) FROM coin_stream_record_5 WHERE uid=$TEST_UID
        UNION ALL SELECT COUNT(*) FROM coin_stream_record_6 WHERE uid=$TEST_UID
        UNION ALL SELECT COUNT(*) FROM coin_stream_record_7 WHERE uid=$TEST_UID
        UNION ALL SELECT COUNT(*) FROM coin_stream_record_8 WHERE uid=$TEST_UID
        UNION ALL SELECT COUNT(*) FROM coin_stream_record_9 WHERE uid=$TEST_UID" 2>&1 | grep -v Warning

echo ""
echo "========================================="
echo "测试结果汇总"
echo "========================================="
echo -e "${GREEN}通过: $PASS_COUNT${NC}"
echo -e "${RED}失败: $FAIL_COUNT${NC}"
echo "总计: $((PASS_COUNT + FAIL_COUNT))"
echo ""

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}✓ 所有测试通过！${NC}"
    exit 0
else
    echo -e "${RED}✗ 部分测试失败${NC}"
    exit 1
fi

# 虚拟钱包系统（Wallet Service）

## 项目概述

MyBilibili虚拟钱包系统，实现了双币种账户体系、分布式事务保障和实时清算功能。

### 核心特性

- ✅ **双币种账户体系**：金瓜子（gold/iap_gold）和银瓜子（silver）
- ✅ **分布式事务保障**：TransactionID锁+用户锁+FOR UPDATE
- ✅ **完整流水记录**：成功和失败都记录
- ✅ **数据分表**：user_wallet和coin_stream_record分10张表

### 技术栈

- go-zero v1.6+
- MySQL 8.0（分表）
- Redis 7（分布式锁+缓存）
- etcd（服务发现）
- Prometheus（监控）
- Jaeger（链路追踪）

---

## 快速启动

### 1. 初始化数据库

```bash
# 确保MySQL运行
docker-compose ps mysql

# 执行初始化脚本
mysql -h127.0.0.1 -P33060 -uroot -proot123456 < ../../deploy/sql/010_wallet_schema.sql

# 验证
mysql -h127.0.0.1 -P33060 -uroot -proot123456 mybilibili -e "SHOW TABLES LIKE '%wallet%'"
```

### 2. 启动RPC服务

```bash
cd cmd/rpc
go run wallet.go -f etc/wallet.yaml
```

**RPC服务地址**：`127.0.0.1:9004`  
**监控地址**：`http://127.0.0.1:9084/metrics`

### 3. 启动API服务

```bash
cd cmd/api
go run wallet.go -f etc/wallet.yaml
```

**API服务地址**：`http://127.0.0.1:8004`  
**监控地址**：`http://127.0.0.1:9094/metrics`

---

## API 接口

### 1. 充值接口

```bash
curl -X POST http://localhost:8004/api/wallet/v1/recharge \
  -H "Content-Type: application/json" \
  -d '{
    "uid": 1001,
    "coin_type": "gold",
    "coin_num": 100,
    "transaction_id": "tid_'$(date +%s)'",
    "platform": "android"
  }'
```

**响应示例**：
```json
{
  "gold": 1100,
  "iap_gold": 0,
  "silver": 500
}
```

### 2. 查询余额

```bash
curl "http://localhost:8004/api/wallet/v1/detail?uid=1001&platform=android"
```

### 3. 消费接口

```bash
curl -X POST http://localhost:8004/api/wallet/v1/pay \
  -H "Content-Type: application/json" \
  -d '{
    "uid": 1001,
    "coin_type": "gold",
    "coin_num": 50,
    "transaction_id": "tid_pay_'$(date +%s)'",
    "platform": "android"
  }'
```

### 4. 兑换接口

```bash
curl -X POST http://localhost:8004/api/wallet/v1/exchange \
  -H "Content-Type: application/json" \
  -d '{
    "uid": 1001,
    "src_coin_type": "gold",
    "src_coin_num": 100,
    "dest_coin_type": "silver",
    "dest_coin_num": 100,
    "transaction_id": "tid_exchange_'$(date +%s)'",
    "platform": "android"
  }'
```

### 5. 查询流水

```bash
curl "http://localhost:8004/api/wallet/v1/stream?uid=1001&offset=0&limit=10"
```

---

## 项目结构

```
wallet/
├── cmd/
│   ├── api/                    # HTTP API服务
│   │   ├── desc/wallet.api    # API定义
│   │   ├── etc/wallet.yaml    # 配置
│   │   └── internal/
│   │       ├── handler/       # HTTP Handler
│   │       ├── logic/         # 业务逻辑
│   │       └── svc/           # 服务上下文
│   │
│   └── rpc/                    # gRPC服务
│       ├── wallet.proto        # Proto定义
│       ├── etc/wallet.yaml     # 配置
│       └── internal/
│           ├── logic/          # 核心业务逻辑⭐
│           ├── server/         # gRPC服务器
│           └── svc/            # 服务上下文
│
└── model/                      # 数据模型
    ├── types.go                # 常量和工具
    ├── user_wallet_model.go    # 钱包Model
    ├── coin_stream_record_model.go    # 流水Model
    └── coin_exchange_record_model.go  # 兑换Model
```

---

## 核心实现

### 充值流程（Recharge）

```
1. 参数校验
2. TransactionID锁（Redis，300s）→ 防重复提交
3. 用户锁（Redis，600s）→ 防并发修改
4. 开启MySQL事务
5. SELECT FOR UPDATE（悲观锁）
6. UPDATE余额 + INSERT流水
7. COMMIT
8. 返回新余额
```

### 双重锁机制

```go
// 1. TransactionID锁（幂等性）
tidLockKey := fmt.Sprintf("wallet:lock:tid:%s", transactionId)
redis.SetnxEx(tidLockKey, "locked", 300)

// 2. 用户锁（防超支）
userLockKey := fmt.Sprintf("wallet:lock:user:%d", uid)
redis.Setnx(userLockKey, lockValue)
redis.Expire(userLockKey, 600)

// 3. FOR UPDATE（终极保障）
SELECT * FROM user_wallet_X WHERE uid=? FOR UPDATE
```

### 分表策略

```go
// user_wallet表：按uid取模
tableIndex := uid % 10  // user_wallet_0 ~ user_wallet_9

// coin_stream_record表：按transaction_id hash
h := fnv.New32a()
h.Write([]byte(transactionId))
tableIndex := h.Sum32() % 10
```

---

## 测试数据

系统已预置测试用户：

| UID  | Gold | IapGold | Silver | 表 |
|------|------|---------|--------|-----|
| 1001 | 1000 | 0       | 500    | user_wallet_1 |
| 1002 | 2000 | 500     | 1000   | user_wallet_2 |
| 1003 | 500  | 0       | 200    | user_wallet_3 |

---

## 监控

### Prometheus指标

```bash
# RPC服务监控
curl http://localhost:9084/metrics

# API服务监控
curl http://localhost:9094/metrics
```

### Jaeger链路追踪

访问：`http://localhost:16686`

---

## 开发状态

### ✅ 已完成（阶段一+阶段二）

- [x] 项目结构创建
- [x] 数据库表设计
- [x] Protobuf接口定义
- [x] Model层实现
- [x] RPC Service Context
- [x] **充值功能完整实现**
- [x] API Service Context
- [x] API Handler和Logic

### 🚧 进行中（阶段三）

- [ ] 消费功能实现
- [ ] 余额检查和防超支
- [ ] 单元测试

### 📋 待实现

- [ ] 兑换功能（阶段四）
- [ ] 查询功能（阶段五）
- [ ] Kafka消息发布（阶段六）
- [ ] 快照机制（阶段六）

---

## 常见问题

### Q1: 为什么使用双重锁？

**TransactionID锁**：防止重复提交（幂等性）  
**用户锁**：防止并发修改余额（防超支）  
**FOR UPDATE**：数据库层面最终保障

### Q2: iOS和Android为什么要分开？

Apple政策要求iOS应用内虚拟货币必须通过IAP购买，所以使用`iap_gold`单独存储。

### Q3: 如何验证功能？

```bash
# 1. 充值
./test_recharge.sh

# 2. 查询余额
./test_detail.sh

# 3. 检查流水
mysql -h127.0.0.1 -P33060 -uroot -proot123456 mybilibili \
  -e "SELECT * FROM coin_stream_record_1 WHERE uid=1001 ORDER BY id DESC LIMIT 5"
```

---

##文档

- [设计方案总览](../../设计与修复方案/v2.0.0-虚拟钱包系统-总览.md)
- [第1部分：概述与架构](../../设计与修复方案/v2.0.0-虚拟钱包系统-第1部分-概述与架构.md)
- [第2部分：接口与实现](../../设计与修复方案/v2.0.0-虚拟钱包系统-第2部分-接口与实现.md)
- [第3部分：实施计划](../../设计与修复方案/v2.0.0-虚拟钱包系统-第3部分-实施计划.md)

---

**版本**: v2.0.0  
**状态**: 阶段二完成（充值功能）  
**下一步**: 实现消费功能

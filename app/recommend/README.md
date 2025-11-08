# 推荐系统模块

## 目录结构

```
app/recommend/
├── cmd/
│   └── rpc/
│       ├── etc/
│       │   └── recommend.yaml          # 推荐服务配置
│       ├── internal/
│       │   ├── config/
│       │   │   └── config.go           # 配置结构体
│       │   ├── dao/
│       │   │   └── dao.go              # 数据访问层
│       │   ├── logic/
│       │   │   ├── filter/             # 过滤器
│       │   │   │   └── filter.go
│       │   │   ├── postprocess/        # 后处理
│       │   │   │   └── postprocess.go
│       │   │   ├── rank/               # 排序
│       │   │   │   ├── rank.go
│       │   │   │   └── xgboost_model.go
│       │   │   └── get_recommend_list_logic.go
│       │   └── svc/
│       │       └── service_context.go
│       ├── models/
│       │   └── 0.0.13/                 # XGBoost 模型文件
│       │       ├── config.json
│       │       ├── feature_metadata.json
│       │       ├── feature_names.txt
│       │       ├── model.json
│       │       └── model.proto         # TreeLite Protobuf 格式
│       ├── model/
│       │   └── model.go                # 数据模型
│       ├── recommend.go                # 服务入口
│       └── recommend.proto             # gRPC 接口定义
└── README.md                           # 本文档

app/recall/
└── cmd/
    └── rpc/
        ├── etc/
        │   └── recall.yaml             # 召回服务配置
        ├── internal/
        │   ├── config/
        │   │   └── config.go
        │   ├── dao/
        │   │   └── dao.go
        │   ├── logic/
        │   │   ├── recall_logic.go     # 召回主逻辑
        │   │   └── video_index_logic.go # 视频索引
        │   └── svc/
        │       └── service_context.go
        ├── model/
        │   └── model.go
        ├── recall.go
        └── recall.proto
```

## 功能概述

推荐系统模块包括两个核心服务：

### 1. Recommend-RPC（推荐服务）

**职责**: 提供个性化推荐接口，协调召回、排序、过滤、后处理等环节。

**核心流程**:
1. 加载用户画像
2. 调用召回服务获取候选集
3. 过滤（去重、黑名单、时长限制）
4. 排序（XGBoost 模型或规则排序）
5. 后处理（打散、推荐理由）
6. 存储推荐记录

### 2. Recall-RPC（召回服务）

**职责**: 提供多路召回策略，从海量视频中快速筛选候选集。

**召回策略**:
- **热门召回**: 全站热门视频
- **精选召回**: 运营精选内容
- **I2I 召回**: 基于视频相似度（Like I2I, Pos I2I）
- **标签召回**: 基于用户标签偏好
- **UP主召回**: 用户关注的UP主视频
- **用户画像召回**: 基于用户历史行为的个性化召回

## 技术亮点

### 1. XGBoost 模型集成

- 使用 TreeLite Protobuf 格式存储模型
- 完整的特征工程（65维特征向量）
- 支持规则排序降级

### 2. 三层架构

```
┌──────────────────────────────────────────┐
│  Creative-API (HTTP接口层)               │
└──────────────┬───────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────┐
│  Recommend-RPC (推荐协调层)              │
│  - 用户画像                              │
│  - 过滤器链                              │
│  - 排序模型                              │
│  - 后处理器                              │
└──────────────┬───────────────────────────┘
               │
               ▼
┌──────────────────────────────────────────┐
│  Recall-RPC (召回层)                     │
│  - 多路召回策略                          │
│  - Bloom Filter 去重                     │
│  - 视频索引服务                          │
└──────────────────────────────────────────┘
```

### 3. 数据层设计

**MySQL 表**:
- `video_info`: 视频基础信息和统计数据
- `video_tag`: 视频标签
- `user_behavior`: 用户行为记录
- `user_follow`: 用户关注关系
- `user_blacklist`: 用户黑名单

**Redis 缓存**:
- 用户画像缓存
- 推荐历史记录
- 召回索引（热门、标签、I2I、UP主）
- Bloom Filter（已观看去重）

## 配置说明

### Recommend-RPC 配置

```yaml
Name: recommend.rpc
ListenOn: 127.0.0.1:9005

# etcd 服务注册
Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: recommend.rpc

# MySQL 配置
MySQL:
  DataSource: root:root123456@tcp(127.0.0.1:3306)/mybilibili?charset=utf8mb4&parseTime=true

# Redis 配置
Redis:
  - Host: 127.0.0.1:6379
    Pass: redis123456
    Type: node

# 召回服务配置
RecallRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: recall.rpc

# 排序模型配置
RankModel:
  ModelDir: app/recommend/cmd/rpc/models/0.0.13

# 业务配置
Business:
  RecallLimit: 500        # 召回总数
  RankEnable: true        # 启用模型排序
  BloomFilterEnable: true # 启用去重
  DurationMin: 60         # 最小时长
  DurationMax: 3600       # 最大时长
```

### Recall-RPC 配置

```yaml
Name: recall.rpc
ListenOn: 127.0.0.1:9004

# etcd 服务注册
Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: recall.rpc

# 召回策略配置
Strategies:
  HotRecall:
    Enable: true
    Limit: 100
    Weight: 1.0
  SelectionRecall:
    Enable: true
    Limit: 50
    Weight: 1.5
  # ... 其他策略
```

## 启动步骤

### 1. 初始化数据

```bash
# 创建数据库表
mysql -h 127.0.0.1 -P 3306 -uroot -proot123456 mybilibili < deploy/sql/03_recommend.sql

# 初始化测试数据
mysql -h 127.0.0.1 -P 3306 -uroot -proot123456 mybilibili < deploy/sql/04_recommend_test_data.sql

# 初始化 Redis 数据
bash deploy/scripts/init_redis_data.sh
```

### 2. 启动服务

```bash
# 启动召回服务
cd app/recall/cmd/rpc
go run recall.go -f etc/recall.yaml

# 启动推荐服务
cd app/recommend/cmd/rpc
go run recommend.go -f etc/recommend.yaml

# 重启 creative-api（包含推荐接口）
cd app/api/creative
go run creative.go -f etc/creative.yaml
```

### 3. 测试接口

```bash
# 获取推荐列表
curl -X GET "http://localhost:8888/api/creative/v1/recommend/list?mid=1000&limit=10&debug=true" | jq .

# 响应示例
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "avid": 10001,
        "title": "深入理解Go语言并发编程",
        "cover": "http://example.com/cover/10001.jpg",
        "duration": 1800,
        "pub_time": 1699500000,
        "zone_id": 36,
        "zone_name": "科技",
        "up_mid": 2001,
        "play": 150000,
        "like": 8500,
        "score": 0.9523,
        "reason": "SelectionRecall",
        "tags": ["Go语言", "并发编程", "教程"]
      }
    ],
    "total": 20,
    "has_more": true,
    "debug_info": {
      "recall_count": "500",
      "filter_count": "420",
      "rank_type": "xgboost"
    }
  }
}
```

## XGBoost 模型

### 特征向量（65维）

1. **分区特征（35维 One-Hot）**
2. **召回特征（10维 One-Hot）**
3. **状态特征（5维 One-Hot）**
4. **全站统计（6维）**: 播放、点赞、收藏、评论、分享、投币
5. **月度统计（5维）**: 播放、点赞、评论、分享、完播
6. **交叉特征（4维）**: 标签匹配度

### 模型文件

- `config.json`: 模型配置
- `feature_metadata.json`: 特征元数据
- `feature_names.txt`: 特征名称列表
- `model.json`: XGBoost 原始模型（JSON格式）
- `model.proto`: TreeLite Protobuf 格式（Go 使用）

### 特征工程

特征归一化采用对数归一化：

```go
func logNormalize(value, max float64) float64 {
    if value <= 0 {
        return 0.0
    }
    return math.Log10(math.Min(value+1.0, max)) / math.Log10(max)
}
```

## 性能优化

### 1. 缓存策略

- 用户画像缓存 1 小时
- 推荐历史缓存 7 天
- 召回索引实时更新

### 2. 并发优化

- 多路召回并发执行
- 批量查询视频信息
- Redis Pipeline 批量操作

### 3. 降级方案

- 模型加载失败 → 规则排序
- 召回服务超时 → 使用热门视频
- Redis 不可用 → 直接查询数据库

## 监控指标

### 业务指标

- 推荐列表请求量（QPS）
- 召回数量分布
- 排序耗时
- CTR（点击率）
- 覆盖率

### 技术指标

- RPC 调用延迟（P50, P99）
- 数据库查询耗时
- Redis 命中率
- 错误率

## 常见问题

### Q1: 推荐结果为空？

**排查**:
1. 检查用户画像是否正常加载
2. 检查召回服务是否返回数据
3. 检查过滤器是否过度过滤
4. 查看 debug_info 中的召回和过滤数量

### Q2: 模型加载失败？

**排查**:
1. 检查 `models/0.0.13/` 目录是否存在
2. 检查 `model.proto` 文件是否完整
3. 查看服务启动日志中的模型加载信息

### Q3: Redis 数据丢失？

**排查**:
1. 重新运行 `init_redis_data.sh` 初始化数据
2. 检查 Redis 持久化配置
3. 检查数据过期时间设置

## 未来优化方向

1. **召回优化**: 加入深度学习召回（双塔模型）
2. **排序优化**: 在线学习、多模型融合
3. **实时性**: 实时用户画像更新
4. **A/B测试**: 推荐策略实验平台
5. **冷启动**: 新用户推荐策略
6. **多样性**: 打散策略优化

## 参考文档

- [设计方案](../../设计与修复方案/v1.0.0.3.md)
- [XGBoost 官方文档](https://xgboost.readthedocs.io/)
- [TreeLite 文档](https://treelite.readthedocs.io/)
- [Go-Zero 文档](https://go-zero.dev/)


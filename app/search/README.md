# Search Service (搜索服务)

基于 go-zero 框架重构的 B站搜索服务，参照 openbilibili/app/service/main/search 实现。

## 功能特性

- 弹幕搜索 (DM Search)
- 弹幕日期搜索 (DM Date Search)
- 弹幕历史搜索 (DM History Search)
- PGC番剧搜索 (PGC Media Search)
- 评论记录搜索 (Reply Record Search)
- 索引更新 (Index Update)

## 技术栈

- Go 1.21+
- go-zero 框架
- Elasticsearch 7.x
- Redis (可选，缓存热门搜索词)

## 项目结构

```
search/
├── cmd/
│   ├── api/           # HTTP API 服务
│   │   ├── desc/      # API 定义文件
│   │   ├── etc/       # 配置文件
│   │   └── internal/  # 内部实现
│   └── rpc/           # gRPC 服务
│       ├── etc/       # 配置文件
│       ├── internal/  # 内部实现
│       └── search/    # protobuf 生成代码
├── model/             # 数据模型 (ES客户端)
└── deploy/
    └── es/            # Elasticsearch 索引配置
```

## 快速开始

### 1. 启动 Elasticsearch

```bash
# 使用 Docker 启动 ES
docker run -d --name elasticsearch \
  -p 9200:9200 -p 9300:9300 \
  -e "discovery.type=single-node" \
  -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" \
  elasticsearch:7.17.0
```

### 2. 创建 ES 索引

```bash
# 开发环境 (创建少量索引)
cd deploy/es
./create_indices_dev.sh http://localhost:9200

# 生产环境 (创建全部索引)
./create_indices.sh http://localhost:9200
```

### 3. 启动 RPC 服务

```bash
cd cmd/rpc
go run search.go -f etc/search.yaml
```

### 4. 启动 API 服务

```bash
cd cmd/api
go run search.go -f etc/search.yaml
```

## Elasticsearch 索引

| 索引名 | 说明 | 分片策略 |
|--------|------|----------|
| `dm_search_{000-999}` | 弹幕搜索索引 | 按 oid % 1000 分片 |
| `dm_date_{YYYY_MM}` | 弹幕日期索引 | 按月份分片 |
| `pgc_media` | PGC番剧索引 | 单索引 |
| `replyrecord_{00-99}` | 评论记录索引 | 按 mid % 100 分片 |

### ES 集群配置

| 集群名 | 用途 |
|--------|------|
| `dmExternal` | 弹幕搜索 (dm_search, dm_date) |
| `externalPublic` | PGC番剧搜索 (pgc_media) |
| `replyExternal` | 评论记录搜索 (replyrecord) |

## API 接口

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 弹幕搜索 | GET | /x/internal/search/dm | 搜索弹幕内容 |
| 弹幕日期搜索 | GET | /x/internal/search/dm/date | 按日期搜索弹幕 |
| 弹幕历史搜索 | GET | /x/internal/search/dmhistory | 搜索弹幕历史 |
| PGC搜索 | GET | /x/internal/search/pgc | 搜索番剧内容 |
| 评论搜索 | GET | /x/internal/search/reply | 搜索评论记录 |
| 弹幕更新 | POST | /x/internal/search/dm/update | 更新弹幕索引 |
| PGC更新 | POST | /x/internal/search/pgc/update | 更新PGC索引 |
| 评论更新 | POST | /x/internal/search/reply/update | 更新评论索引 |

## 配置说明

### RPC 配置 (etc/search.yaml)

```yaml
Name: search.rpc
ListenOn: 0.0.0.0:8080

ES:
  dmExternal:
    - "http://localhost:9200"
  externalPublic:
    - "http://localhost:9200"
  replyExternal:
    - "http://localhost:9200"
```

### API 配置 (etc/search.yaml)

```yaml
Name: search-api
Host: 0.0.0.0
Port: 8888

SearchRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: search.rpc
```

## 数据流说明

```
业务模块 (dm/reply/pgc)
    ↓ 写入
MySQL 数据库
    ↓ 数据同步 (Canal/Kafka)
Elasticsearch 集群
    ↑ 搜索查询
Search 服务
```

Search 服务不直接管理原始数据，只负责：
1. 从 ES 搜索数据
2. 更新 ES 索引

原始数据存储在各业务模块的 MySQL 中，通过数据同步机制同步到 ES。

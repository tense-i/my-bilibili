# Elasticsearch 索引配置

本目录包含 search 模块所需的 Elasticsearch 索引映射定义。

## 索引列表

| 索引名 | 说明 | 分片策略 |
|--------|------|----------|
| `dm_search_{000-999}` | 弹幕搜索索引 | 按 oid % 1000 分片 |
| `dm_date_{YYYY_MM}` | 弹幕日期索引 | 按月份分片 |
| `pgc_media` | PGC番剧索引 | 单索引 |
| `replyrecord_{00-99}` | 评论记录索引 | 按 mid % 100 分片 |

## 创建索引

### 方式一：使用脚本批量创建

```bash
# 创建弹幕搜索索引 (1000个分片)
./create_dm_search_indices.sh

# 创建评论记录索引 (100个分片)
./create_replyrecord_indices.sh

# 创建 PGC 索引
curl -X PUT "localhost:9200/pgc_media" -H 'Content-Type: application/json' -d @pgc_media_mapping.json
```

### 方式二：手动创建单个索引

```bash
# 创建单个弹幕索引
curl -X PUT "localhost:9200/dm_search_000" -H 'Content-Type: application/json' -d @dm_search_mapping.json

# 创建 PGC 索引
curl -X PUT "localhost:9200/pgc_media" -H 'Content-Type: application/json' -d @pgc_media_mapping.json

# 创建评论记录索引
curl -X PUT "localhost:9200/replyrecord_00" -H 'Content-Type: application/json' -d @replyrecord_mapping.json
```

## ES 集群配置

search 服务需要连接以下 ES 集群：

| 集群名 | 用途 |
|--------|------|
| `dmExternal` | 弹幕搜索 |
| `externalPublic` | PGC番剧搜索 |
| `replyExternal` | 评论记录搜索 |

在 `search.yaml` 配置文件中配置：

```yaml
ES:
  dmExternal:
    - "http://localhost:9200"
  externalPublic:
    - "http://localhost:9200"
  replyExternal:
    - "http://localhost:9200"
```

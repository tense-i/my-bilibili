# Grafana 监控配置说明

## 目录结构

```
grafana/
├── README.md                              # 本文档
├── provisioning/                          # 自动配置目录
│   ├── datasources/                       # 数据源配置
│   │   └── prometheus.yml                 # Prometheus 数据源
│   └── dashboards/                        # 仪表板配置
│       └── default.yml                    # 默认仪表板提供者
└── dashboards/                            # 仪表板 JSON（可选）
    ├── mysql-overview.json                # MySQL 监控仪表板
    ├── redis-overview.json                # Redis 监控仪表板
    └── mybilibili-services.json           # 微服务监控仪表板
```

## 监控指标说明

### 1. MySQL 监控指标

#### 连接指标
- `mysql_global_status_threads_connected` - 当前连接数
- `mysql_global_status_threads_running` - 正在运行的线程数
- `mysql_global_status_max_used_connections` - 最大使用连接数
- `mysql_global_variables_max_connections` - 最大连接数限制

#### 性能指标
- `mysql_global_status_queries` - 查询总数
- `mysql_global_status_slow_queries` - 慢查询数
- `mysql_global_status_questions` - 问题数
- `rate(mysql_global_status_queries[5m])` - QPS（每秒查询数）

#### InnoDB 指标
- `mysql_global_status_innodb_buffer_pool_read_requests` - 缓冲池读取请求
- `mysql_global_status_innodb_buffer_pool_reads` - 缓冲池读取
- `mysql_global_status_innodb_buffer_pool_pages_data` - 数据页数
- `mysql_global_status_innodb_buffer_pool_pages_free` - 空闲页数

#### 表锁指标
- `mysql_global_status_table_locks_waited` - 表锁等待数
- `mysql_global_status_table_locks_immediate` - 表锁立即获取数

### 2. Redis 监控指标

#### 基本信息
- `redis_up` - Redis 是否在线 (1=在线, 0=离线)
- `redis_uptime_in_seconds` - 运行时间（秒）
- `redis_version` - Redis 版本

#### 内存指标
- `redis_memory_used_bytes` - 已使用内存（字节）
- `redis_memory_max_bytes` - 最大内存（字节）
- `redis_memory_used_rss_bytes` - RSS 内存
- `redis_memory_fragmentation_ratio` - 内存碎片率

#### 连接指标
- `redis_connected_clients` - 当前连接客户端数
- `redis_connected_slaves` - 连接的从节点数
- `redis_blocked_clients` - 阻塞的客户端数

#### 性能指标
- `rate(redis_commands_processed_total[5m])` - QPS（每秒命令数）
- `redis_keyspace_hits_total` - 命中次数
- `redis_keyspace_misses_total` - 未命中次数
- `redis_keyspace_hits_total / (redis_keyspace_hits_total + redis_keyspace_misses_total)` - 命中率

#### 持久化指标
- `redis_rdb_last_save_timestamp_seconds` - 最后一次 RDB 保存时间
- `redis_aof_enabled` - AOF 是否启用
- `redis_rdb_changes_since_last_save` - 自上次保存以来的变更数

### 3. 微服务监控指标

#### Go 运行时指标
- `go_goroutines` - Goroutine 数量
- `go_memstats_alloc_bytes` - 分配的内存
- `go_memstats_heap_inuse_bytes` - 堆使用内存
- `go_gc_duration_seconds` - GC 耗时

#### HTTP 指标（API 服务）
- `http_server_requests_code_total` - HTTP 请求总数（按状态码）
- `http_server_requests_duration_ms` - HTTP 请求耗时
- `rate(http_server_requests_code_total[5m])` - HTTP QPS

#### RPC 指标（RPC 服务）
- `rpc_server_requests_code_total` - RPC 请求总数
- `rpc_server_requests_duration_ms` - RPC 请求耗时
- `rate(rpc_server_requests_code_total[5m])` - RPC QPS

### 4. 系统监控指标（Node Exporter）

#### CPU 指标
- `node_cpu_seconds_total` - CPU 使用时间
- `rate(node_cpu_seconds_total{mode="idle"}[5m])` - CPU 空闲率
- `100 - (rate(node_cpu_seconds_total{mode="idle"}[5m]) * 100)` - CPU 使用率

#### 内存指标
- `node_memory_MemTotal_bytes` - 总内存
- `node_memory_MemAvailable_bytes` - 可用内存
- `node_memory_MemFree_bytes` - 空闲内存
- `(node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes * 100` - 内存使用率

#### 磁盘指标
- `node_filesystem_size_bytes` - 文件系统大小
- `node_filesystem_avail_bytes` - 可用空间
- `node_disk_read_bytes_total` - 磁盘读取字节数
- `node_disk_written_bytes_total` - 磁盘写入字节数

#### 网络指标
- `node_network_receive_bytes_total` - 网络接收字节数
- `node_network_transmit_bytes_total` - 网络发送字节数

## 推荐的 Grafana 仪表板

### 官方仪表板 ID（可直接导入）

1. **MySQL Overview**
   - Dashboard ID: `7362`
   - URL: https://grafana.com/grafana/dashboards/7362
   - 说明：MySQL 数据库综合监控

2. **Redis Dashboard**
   - Dashboard ID: `11835`
   - URL: https://grafana.com/grafana/dashboards/11835
   - 说明：Redis 缓存综合监控

3. **Node Exporter Full**
   - Dashboard ID: `1860`
   - URL: https://grafana.com/grafana/dashboards/1860
   - 说明：系统资源综合监控

4. **Go Processes**
   - Dashboard ID: `6671`
   - URL: https://grafana.com/grafana/dashboards/6671
   - 说明：Go 应用程序监控

## 导入仪表板步骤

### 方式一：通过 Dashboard ID 导入（推荐）

1. 访问 Grafana：http://localhost:3000
2. 登录（用户名：admin，密码：admin123456）
3. 点击左侧菜单 "+" → "Import"
4. 输入 Dashboard ID（如 `7362`）
5. 点击 "Load"
6. 选择数据源："Prometheus"
7. 点击 "Import"

### 方式二：通过 JSON 文件导入

1. 访问 Grafana：http://localhost:3000
2. 点击左侧菜单 "+" → "Import"
3. 点击 "Upload JSON file"
4. 选择对应的 JSON 文件
5. 选择数据源："Prometheus"
6. 点击 "Import"

## 常用 Prometheus 查询示例

### MySQL 查询

```promql
# MySQL 是否在线
up{job="mysql"}

# MySQL QPS
rate(mysql_global_status_queries[5m])

# MySQL 连接数
mysql_global_status_threads_connected

# MySQL 慢查询数
rate(mysql_global_status_slow_queries[5m])

# InnoDB 缓冲池命中率
100 * mysql_global_status_innodb_buffer_pool_read_requests / 
(mysql_global_status_innodb_buffer_pool_read_requests + mysql_global_status_innodb_buffer_pool_reads)
```

### Redis 查询

```promql
# Redis 是否在线
redis_up

# Redis QPS
rate(redis_commands_processed_total[5m])

# Redis 内存使用率
redis_memory_used_bytes / redis_memory_max_bytes * 100

# Redis 命中率
redis_keyspace_hits_total / 
(redis_keyspace_hits_total + redis_keyspace_misses_total) * 100

# Redis 连接数
redis_connected_clients
```

### 微服务查询

```promql
# 服务是否在线
up{app="mybilibili"}

# HTTP QPS
rate(http_server_requests_code_total[5m])

# HTTP 平均响应时间
rate(http_server_requests_duration_ms_sum[5m]) / 
rate(http_server_requests_duration_ms_count[5m])

# HTTP 成功率
sum(rate(http_server_requests_code_total{code=~"2.."}[5m])) / 
sum(rate(http_server_requests_code_total[5m])) * 100

# RPC QPS
rate(rpc_server_requests_code_total[5m])
```

## 告警规则示例

创建 `deploy/prometheus/rules/alerts.yml`：

```yaml
groups:
  - name: mybilibili_alerts
    interval: 30s
    rules:
      # MySQL 告警
      - alert: MySQLDown
        expr: up{job="mysql"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "MySQL is down"
          description: "MySQL instance {{ $labels.instance }} is down"

      - alert: MySQLTooManyConnections
        expr: mysql_global_status_threads_connected / mysql_global_variables_max_connections > 0.8
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "MySQL connections too high"
          description: "MySQL connections usage is {{ $value | humanizePercentage }}"

      # Redis 告警
      - alert: RedisDown
        expr: redis_up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Redis is down"
          description: "Redis instance {{ $labels.instance }} is down"

      - alert: RedisMemoryHigh
        expr: redis_memory_used_bytes / redis_memory_max_bytes > 0.9
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Redis memory usage too high"
          description: "Redis memory usage is {{ $value | humanizePercentage }}"

      # 服务告警
      - alert: ServiceDown
        expr: up{app="mybilibili"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Service {{ $labels.service }} is down"
          description: "Service {{ $labels.service }} has been down for more than 1 minute"

      - alert: HighErrorRate
        expr: |
          sum(rate(http_server_requests_code_total{code=~"5.."}[5m])) by (service) /
          sum(rate(http_server_requests_code_total[5m])) by (service) > 0.05
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High error rate on {{ $labels.service }}"
          description: "Error rate is {{ $value | humanizePercentage }}"
```

## 访问地址

- **Grafana UI**: http://localhost:3000
  - 用户名：`admin`
  - 密码：`admin123456`

- **Prometheus UI**: http://localhost:9090
  - 查询指标
  - 查看 Targets 状态

- **Exporters**:
  - MySQL Exporter: http://localhost:9104/metrics
  - Redis Exporter: http://localhost:9121/metrics
  - Node Exporter: http://localhost:9100/metrics

## 验证监控是否正常

### 1. 检查 Prometheus Targets

访问 http://localhost:9090/targets

确认所有 targets 状态为 "UP"：
- ✅ prometheus (1/1 up)
- ✅ mysql (1/1 up)
- ✅ redis (1/1 up)
- ✅ node (1/1 up)
- ✅ video-rpc (1/1 up)
- ✅ hotrank-rpc (1/1 up)
- ✅ creative-api (1/1 up)

### 2. 查询测试指标

在 Prometheus UI 中测试查询：

```promql
# 测试 MySQL 指标
mysql_up

# 测试 Redis 指标
redis_up

# 测试服务指标
up{app="mybilibili"}
```

### 3. 检查 Grafana 数据源

1. 登录 Grafana
2. Configuration → Data Sources
3. 确认 Prometheus 数据源状态为绿色 "✓"

## 故障排查

### MySQL Exporter 无法连接

**问题**：MySQL Exporter 显示 DOWN

**解决方案**：
```bash
# 检查 MySQL 是否运行
docker logs mybilibili-mysql

# 检查 MySQL Exporter 日志
docker logs mybilibili-mysql-exporter

# 检查连接字符串
# DATA_SOURCE_NAME=root:root123456@(mysql:3306)/
```

### Redis Exporter 无法连接

**问题**：Redis Exporter 显示 DOWN

**解决方案**：
```bash
# 检查 Redis 是否运行
docker logs mybilibili-redis

# 检查 Redis Exporter 日志
docker logs mybilibili-redis-exporter

# 测试 Redis 连接
docker exec mybilibili-redis redis-cli -a redis123456 ping
```

### 微服务指标无法采集

**问题**：微服务 targets 显示 DOWN

**解决方案**：
```bash
# 确认服务是否运行
ps aux | grep -E "(video-rpc|hotrank-rpc|creative-api)"

# 测试 Prometheus 端点
curl http://127.0.0.1:9081/metrics  # video-rpc
curl http://127.0.0.1:9093/metrics  # hotrank-rpc
curl http://127.0.0.1:9091/metrics  # creative-api
```

## 参考资料

- [Prometheus 官方文档](https://prometheus.io/docs/)
- [Grafana 官方文档](https://grafana.com/docs/)
- [MySQL Exporter](https://github.com/prometheus/mysqld_exporter)
- [Redis Exporter](https://github.com/oliver006/redis_exporter)
- [Node Exporter](https://github.com/prometheus/node_exporter)
- [Grafana Dashboards](https://grafana.com/grafana/dashboards/)

---

**最后更新**：2025-11-08
**作者**：mybilibili 团队


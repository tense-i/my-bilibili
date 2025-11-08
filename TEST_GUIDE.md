# MyBilibili 热门视频排行榜系统测试指南

## 测试目标

验证热门视频排行榜系统的完整功能：
- ✅ 视频服务（video-rpc）
- ✅ 热度计算任务（hotrank-job）
- ✅ 热门排行榜服务（hotrank-rpc）
- ✅ API 网关（creative-api）

## 前置准备

### 1. 启动基础设施

```bash
# 启动 MySQL、Redis、etcd、Jaeger
cd deploy
docker-compose up -d mysql redis etcd jaeger
```

### 2. 初始化数据库

```bash
# 连接到 MySQL
mysql -h127.0.0.1 -P3306 -uroot -p123456

# 执行初始化脚本
source deploy/sql/001_init.sql
source deploy/sql/002_test_data.sql
```

### 3. 插入测试数据到 academy_archive 表

```sql
-- 插入热度榜数据（video_info 表中的视频）
INSERT INTO academy_archive (oid, uid, business, region_id, hot, state, ctime, mtime) VALUES
(1001, 1, 1, 1, 0, 0, NOW(), NOW()),
(1002, 2, 1, 2, 0, 0, NOW(), NOW()),
(1003, 3, 1, 1, 0, 0, NOW(), NOW()),
(1004, 4, 1, 3, 0, 0, NOW(), NOW()),
(1005, 5, 1, 2, 0, 0, NOW(), NOW());
```

## 测试流程

### 阶段1：启动 video-rpc 服务

```bash
# 终端1
cd app/video/cmd/rpc
go run video.go -f etc/video.yaml
```

**预期输出：**
```
Starting rpc server at 0.0.0.0:9002...
```

**测试验证：**
```bash
# 检查 etcd 服务注册
etcdctl get --prefix video.rpc
```

---

### 阶段2：启动 hotrank-job 热度计算任务

```bash
# 终端2
cd app/hotrank/cmd/job
go run job.go -f etc/hotrank-job.yaml
```

**预期输出：**
```
Starting hotrank-job...
hotrank-job started successfully, press Ctrl+C to exit
FlushHot started for business: 1
FlushHot success: processed 5 videos, last_id=5
```

**测试验证：**
```sql
-- 查看热度值是否已计算
SELECT id, oid, hot, mtime FROM academy_archive ORDER BY hot DESC;
```

**预期结果：**
- `hot` 字段应该有值（根据公式计算）
- 热度值从高到低排列
- `mtime` 更新为最新时间

---

### 阶段3：启动 hotrank-rpc 服务

```bash
# 终端3
cd app/hotrank/cmd/rpc
go run hotrank.go -f etc/hotrank.yaml
```

**预期输出：**
```
Starting rpc server at 0.0.0.0:9003...
```

**测试验证：**
```bash
# 检查 etcd 服务注册
etcdctl get --prefix hotrank.rpc
```

---

### 阶段4：启动 creative-api 网关

```bash
# 终端4
cd app/api/creative
go run creative.go -f etc/creative-api.yaml
```

**预期输出：**
```
Starting server at 0.0.0.0:8001...
```

---

## API 测试

### 1. 测试全站热门排行榜

```bash
curl "http://localhost:8001/api/creative/v1/hotrank/list?offset=0&limit=10" | jq
```

**预期响应：**
```json
{
  "list": [
    {
      "oid": 1002,
      "hot": 15250,
      "rank": 1,
      "title": "深入理解Go并发编程",
      "cover": "https://example.com/cover2.jpg",
      "author": "李四",
      "view": 50000,
      "like": 3000
    },
    {
      "oid": 1001,
      "hot": 12875,
      "rank": 2,
      "title": "Go语言入门教程",
      "cover": "https://example.com/cover1.jpg",
      "author": "张三",
      "view": 10000,
      "like": 500
    }
  ],
  "total": 5
}
```

### 2. 测试分区热门排行榜

```bash
# 查询分区1的热门视频
curl "http://localhost:8001/api/creative/v1/hotrank/region?region_id=1&offset=0&limit=10" | jq
```

**预期响应：**
```json
{
  "list": [
    {
      "oid": 1001,
      "hot": 12875,
      "rank": 1,
      "title": "Go语言入门教程",
      "cover": "https://example.com/cover1.jpg",
      "author": "张三",
      "view": 10000,
      "like": 500
    },
    {
      "oid": 1003,
      "hot": 8750,
      "rank": 2,
      "title": "Rust入门与实践",
      "cover": "https://example.com/cover3.jpg",
      "author": "王五",
      "view": 30000,
      "like": 2000
    }
  ],
  "total": 2
}
```

### 3. 测试视频详情

```bash
curl "http://localhost:8001/api/creative/v1/video/1001" | jq
```

**预期响应：**
```json
{
  "video": {
    "vid": 1001,
    "title": "Go语言入门教程",
    "cover": "https://example.com/cover1.jpg",
    "author_id": 1,
    "author_name": "张三",
    "region_id": 1,
    "pub_time": 1699401600,
    "duration": 1200,
    "desc": "适合初学者的Go语言教程",
    "view": 10000,
    "like": 500,
    "coin": 200,
    "fav": 300,
    "share": 100,
    "reply": 50,
    "danmaku": 1000,
    "hot": 12875
  }
}
```

## 验证热度计算公式

### 热度计算公式
```
hot = coin×0.4 + fav×0.3 + danmaku×0.4 + reply×0.4 + view×0.25 + like×0.4 + share×0.6
新视频（24小时内）：hot × 1.5
```

### 手动验证示例（vid=1001）

```
视频数据：
- view: 10000
- like: 500
- coin: 200
- fav: 300
- share: 100
- reply: 50
- danmaku: 1000

计算：
hot = 200×0.4 + 300×0.3 + 1000×0.4 + 50×0.4 + 10000×0.25 + 500×0.4 + 100×0.6
    = 80 + 90 + 400 + 20 + 2500 + 200 + 60
    = 3350

如果是新视频（发布时间 < 24小时）：
hot = 3350 × 1.5 = 5025
```

## 压力测试

### 1. 测试热度计算性能

```bash
# 批量插入1000条测试数据
for i in {1..1000}; do
  mysql -h127.0.0.1 -P3306 -uroot -p123456 mybilibili -e \
  "INSERT INTO academy_archive (oid, uid, business, region_id, hot, state, ctime, mtime) 
   VALUES ($((2000+i)), 1, 1, 1, 0, 0, NOW(), NOW());"
done
```

观察 hotrank-job 的处理速度和日志输出。

### 2. 测试 API 并发

```bash
# 使用 ab 工具测试
ab -n 1000 -c 10 "http://localhost:8001/api/creative/v1/hotrank/list?offset=0&limit=10"
```

**预期指标：**
- QPS > 500
- 平均响应时间 < 50ms
- 成功率 = 100%

## 监控和追踪

### 1. Prometheus 监控

- video-rpc: http://localhost:9092/metrics
- hotrank-rpc: http://localhost:9093/metrics
- creative-api: http://localhost:9091/metrics

### 2. Jaeger 分布式追踪

访问：http://localhost:16686

查看请求链路：
- creative-api → hotrank-rpc → MySQL
- creative-api → video-rpc → MySQL

## 常见问题排查

### 1. hotrank-job 无法连接 video-rpc

**症状：** `VideoRpc.Archives BatchGetVideoInfo error`

**解决：**
- 检查 etcd 服务是否正常：`etcdctl get --prefix video.rpc`
- 检查 video-rpc 是否启动
- 检查网络连接

### 2. 热度值为0或不更新

**症状：** `academy_archive` 表的 `hot` 字段全是0

**解决：**
- 检查 `video_info` 和 `video_stat` 是否有数据
- 检查 `oid` 是否匹配 `vid`
- 查看 hotrank-job 日志

### 3. API 返回空列表

**症状：** `/api/creative/v1/hotrank/list` 返回 `{"list": [], "total": 0}`

**解决：**
- 检查 `academy_archive` 表是否有数据
- 检查 `state` 字段是否为0（0表示正常）
- 检查 hotrank-rpc 是否启动

## 测试结果检查清单

- [ ] video-rpc 成功启动并注册到 etcd
- [ ] hotrank-job 成功连接 video-rpc 并计算热度
- [ ] academy_archive 表的 hot 字段有正确的值
- [ ] hotrank-rpc 成功启动并注册到 etcd
- [ ] creative-api 成功启动
- [ ] 全站排行榜 API 返回正确数据
- [ ] 分区排行榜 API 返回正确数据
- [ ] 视频详情 API 返回正确数据
- [ ] 热度计算公式验证正确
- [ ] Prometheus 指标正常
- [ ] Jaeger 追踪链路完整

## 下一步优化

1. **性能优化**
   - Redis 缓存热门排行榜
   - 视频详情缓存
   - 批量查询优化

2. **功能扩展**
   - 用户服务（user-rpc）
   - 评论服务（comment-rpc）
   - 弹幕服务（danmaku-rpc）

3. **运维优化**
   - Kubernetes 部署
   - 自动扩缩容
   - 监控告警


# MyBilibili - 热门视频排行榜系统

基于 go-zero 框架构建的微服务系统，实现热门视频排行榜功能。

## 项目简介

MyBilibili 是一个模仿 Bilibili 主项目的热门视频排行榜系统，采用微服务架构设计：

- **框架**：go-zero
- **数据库**：MySQL 8.0
- **缓存**：Redis 7
- **服务发现**：etcd
- **链路追踪**：Jaeger
- **监控**：Prometheus + Grafana

## 核心功能

### 1. 热度计算算法
- **多维度加权**：硬币×0.4 + 收藏×0.3 + 弹幕×0.4 + 评论×0.4 + 播放×0.25 + 点赞×0.4 + 分享×0.6
- **新视频提权**：24小时内发布的视频热度×1.5
- **防刷机制**：降低播放量权重，提高互动指标权重

### 2. 服务架构
```
├── video-rpc         # 视频服务（视频信息、统计数据）
├── hotrank-job       # 热度计算定时任务⭐核心
├── hotrank-rpc       # 热门排行榜RPC服务
└── creative-api      # HTTP API网关
```

## 快速开始

### 1. 环境准备

**依赖**：
- Go 1.21+
- Docker & Docker Compose
- goctl（go-zero 代码生成工具）

**安装 goctl**：
```bash
go install github.com/zeromicro/go-zero/tools/goctl@latest
```

### 2. 启动基础服务

```bash
cd deploy
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f mysql
```

### 3. 初始化数据库

```bash
# 连接 MySQL（密码：root123456）
mysql -h127.0.0.1 -P33060 -uroot -p

# 执行初始化脚本
source deploy/sql/001_init.sql
source deploy/sql/002_test_data.sql

# 验证数据
USE mybilibili;
SELECT COUNT(*) FROM video_info;
SELECT COUNT(*) FROM academy_archive;
```

### 4. 启动服务

```bash
# 1. 启动视频服务
cd app/video/cmd/rpc
go run video.go -f etc/video.yaml

# 2. 启动热度计算任务
cd app/hotrank/cmd/job
go run hotrankjob.go -f etc/hotrank-job.yaml

# 3. 启动排行榜RPC服务
cd app/hotrank/cmd/rpc
go run hotrank.go -f etc/hotrank.yaml

# 4. 启动API网关
cd app/creative/cmd/api
go run creative.go -f etc/creative-api.yaml
```

## 项目结构

```
mybilibili/
├── app/                          # 应用服务
│   ├── video/                    # 视频服务
│   │   └── cmd/rpc/             # RPC 服务
│   ├── hotrank/                  # 热门排行榜
│   │   ├── cmd/job/             # 定时任务
│   │   └── cmd/rpc/             # RPC 服务
│   └── creative/                 # 创作中心
│       └── cmd/api/             # HTTP API
│
├── common/                       # 公共组件
│   ├── model/                   # 数据模型
│   ├── xerr/                    # 错误码
│   ├── tool/                    # 工具函数
│   └── result/                  # 响应封装
│
├── deploy/                       # 部署配置
│   ├── docker-compose.yml       # Docker编排
│   ├── sql/                     # SQL脚本
│   └── prometheus/              # 监控配置
│
└── doc/                         # 文档
```

## 核心设计

### 1. 热度计算流程（参考主项目）

```go
// 定时任务循环
for {
    // 1. 游标分页查询（每批30条）
    arcs := aca.Archives(id, business, 30)
    
    // 2. 批量获取视频数据
    stats := videoRpc.BatchGetVideoStat(vids)
    infos := videoRpc.BatchGetVideoInfo(vids)
    
    // 3. 计算热度值
    for each video {
        hot = coin*0.4 + fav*0.3 + danmaku*0.4 + 
              reply*0.4 + view*0.25 + like*0.4 + share*0.6
        if isNew (24h) {
            hot *= 1.5  // 新视频提权
        }
    }
    
    // 4. CASE WHEN 批量更新
    UPDATE academy_archive SET hot = CASE oid
        WHEN vid1 THEN hot1
        WHEN vid2 THEN hot2
        ...
    END WHERE oid IN (vid1, vid2, ...)
    
    // 5. 循环完成后休眠1小时
    if no more data {
        sleep(1 hour)
        continue from beginning
    }
}
```

### 2. 数据库表设计（参考主项目）

#### academy_archive 表
```sql
CREATE TABLE `academy_archive` (
  `id` bigint AUTO_INCREMENT,        -- 游标分页用
  `oid` bigint,                      -- 视频ID
  `hot` bigint DEFAULT 0,            -- 热度值
  `business` tinyint DEFAULT 1,      -- 业务类型
  `pub_time` bigint,                 -- 发布时间
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_oid` (`oid`),
  KEY `idx_business_state_id` (`business`, `state`, `id`),
  KEY `idx_hot_global` (`hot`, `state`)
);
```

## 监控和追踪

### Jaeger UI
- URL: http://localhost:16686
- 查看服务链路追踪

### Prometheus
- URL: http://localhost:9090
- 查看指标数据

### Grafana
- URL: http://localhost:3000
- 用户名：admin
- 密码：admin123456

## API 文档

### 1. 获取全站热门排行榜
```bash
GET /api/v1/academy/rank/global?page=1&page_size=20
```

### 2. 获取分区热门排行榜
```bash
GET /api/v1/academy/rank/region?region_id=1&page=1&page_size=20
```

## 开发指南

### 生成代码

```bash
# 生成 RPC 代码
cd app/video/cmd/rpc
goctl rpc protoc video.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style go_zero

# 生成 API 代码
cd app/creative/cmd/api
goctl api go -api creative.api -dir . --style go_zero

# 生成 Model 代码
goctl model mysql datasource -url="root:root123456@tcp(127.0.0.1:33060)/mybilibili" \
  -table="video_info,video_stat" -dir=./common/model --style go_zero
```

### 测试

```bash
# 运行单元测试
go test ./...

# 运行集成测试
go test -tags=integration ./...
```

## 参考资料

- [go-zero 官方文档](https://go-zero.dev)
- [go-zero-looklook 实战项目](https://github.com/Mikaelemmmm/go-zero-looklook)
- [Bilibili 主项目架构](../bilibili)

## 许可证

MIT License


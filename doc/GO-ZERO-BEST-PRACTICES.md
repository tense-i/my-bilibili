# Go-Zero 微服务开发最佳实践

> 结合 MyBilibili 热门视频排行榜系统和 go-zero-looklook 电商项目的实战经验总结

**作者**：mybilibili 团队  
**版本**：v1.0  
**更新时间**：2025-11-08

---

## 目录

- [1. 项目架构设计](#1-项目架构设计)
- [2. 代码组织结构](#2-代码组织结构)
- [3. 配置管理](#3-配置管理)
- [4. 服务发现与注册](#4-服务发现与注册)
- [5. 数据库与缓存](#5-数据库与缓存)
- [6. 错误处理](#6-错误处理)
- [7. RPC 服务开发](#7-rpc-服务开发)
- [8. API 网关开发](#8-api-网关开发)
- [9. 定时任务与消息队列](#9-定时任务与消息队列)
- [10. 链路追踪与监控](#10-链路追踪与监控)
- [11. 性能优化](#11-性能优化)
- [12. 部署策略](#12-部署策略)
- [13. 常见问题与解决方案](#13-常见问题与解决方案)

---

## 1. 项目架构设计

### 1.1 微服务分层架构

**推荐分层**（从上到下）：

```
┌────────────────────────────────────┐
│    Nginx / API Gateway             │  ← 统一网关入口
├────────────────────────────────────┤
│    API 服务层（HTTP）              │  ← 聚合服务、协议转换
├────────────────────────────────────┤
│    RPC 服务层（gRPC）              │  ← 核心业务逻辑
├────────────────────────────────────┤
│    Model 层（数据访问）            │  ← 数据库/缓存访问
├────────────────────────────────────┤
│    MySQL / Redis / Kafka           │  ← 基础设施
└────────────────────────────────────┘
```

**实践案例**：

#### MyBilibili 架构
```
creative-api (HTTP)
    ├── hotrank-rpc (排行榜服务)
    │   └── academy_archive 表
    └── video-rpc (视频服务)
        ├── video_info 表
        └── video_stat 表

hotrank-job (定时任务)
    ├── 调用 video-rpc
    └── 更新 academy_archive
```

#### go-zero-looklook 架构
```
nginx gateway
    ├── usercenter-api (用户中心)
    │   └── usercenter-rpc
    ├── travel-api (旅游服务)
    │   └── travel-rpc
    ├── order-api (订单服务)
    │   ├── order-rpc
    │   └── payment-rpc
    └── payment-api (支付服务)
        └── payment-rpc
```

### 1.2 服务划分原则

**✅ 推荐做法**：

1. **按业务领域划分**
   - 用户服务（user）
   - 视频服务（video）
   - 订单服务（order）
   - 支付服务（payment）

2. **单一职责原则**
   - 每个服务只负责一个业务领域
   - 避免服务之间过度耦合

3. **数据独立性**
   - 每个服务拥有自己的数据库
   - 避免跨服务直接访问数据库

**❌ 不推荐做法**：

- ❌ 服务划分过细（过度设计）
- ❌ 服务之间循环依赖
- ❌ 共享数据库表

---

## 2. 代码组织结构

### 2.1 标准项目结构

```
project/
├── app/                          # 业务服务目录
│   ├── service1/                # 服务1
│   │   ├── cmd/                 # 启动入口
│   │   │   ├── api/            # HTTP 服务
│   │   │   │   ├── desc/       # API 定义文件（.api）
│   │   │   │   ├── etc/        # 配置文件
│   │   │   │   ├── internal/   # 内部实现
│   │   │   │   │   ├── config/     # 配置结构
│   │   │   │   │   ├── handler/    # 路由处理
│   │   │   │   │   ├── logic/      # 业务逻辑 ⭐核心
│   │   │   │   │   ├── svc/        # 服务上下文
│   │   │   │   │   └── types/      # 类型定义
│   │   │   │   └── service.go  # 启动文件
│   │   │   ├── rpc/            # gRPC 服务
│   │   │   │   ├── etc/        # 配置文件
│   │   │   │   ├── internal/   # 内部实现
│   │   │   │   │   ├── config/     # 配置结构
│   │   │   │   │   ├── logic/      # 业务逻辑 ⭐核心
│   │   │   │   │   ├── server/     # gRPC 服务器
│   │   │   │   │   └── svc/        # 服务上下文
│   │   │   │   ├── pb/         # Protobuf 生成文件
│   │   │   │   ├── service.proto    # Protobuf 定义
│   │   │   │   ├── service_client/  # 客户端封装
│   │   │   │   └── service.go       # 启动文件
│   │   │   └── mq/             # 消息队列消费者（可选）
│   │   └── model/              # 数据模型
│   └── service2/
├── common/                       # 公共组件
│   ├── xerr/                    # 错误码定义
│   ├── result/                  # 统一响应
│   ├── tool/                    # 工具函数
│   ├── middleware/              # 中间件
│   └── interceptor/             # 拦截器
├── deploy/                       # 部署相关
│   ├── docker-compose.yml       # 本地开发环境
│   ├── k8s/                     # K8s 配置
│   └── sql/                     # 数据库脚本
├── doc/                         # 文档
├── go.mod                       # Go 模块
├── go.sum
└── README.md
```

### 2.2 目录命名规范

| 目录 | 说明 | 示例 |
|-----|------|------|
| `cmd` | 命令行入口 | `app/video/cmd/rpc/video.go` |
| `internal` | 内部实现（不对外暴露） | `internal/logic/` |
| `api` | HTTP API 服务 | `app/creative/cmd/api/` |
| `rpc` | gRPC 服务 | `app/video/cmd/rpc/` |
| `model` | 数据模型 | `app/video/model/` |
| `pb` | Protobuf 生成文件 | `rpc/pb/` |

---

## 3. 配置管理

### 3.1 配置文件结构

**RPC 服务配置示例**（video.yaml）：

```yaml
Name: video.rpc
ListenOn: 127.0.0.1:9001
Mode: dev                        # dev/test/prod

# Etcd 服务注册（生产环境）
Etcd:
  Hosts:
    - 127.0.0.1:23790
  Key: video.rpc
  # User: ""                     # 可选：认证
  # Pass: ""

# MySQL 配置
Mysql:
  DataSource: root:password@tcp(127.0.0.1:3306)/database?charset=utf8mb4&parseTime=true&loc=Local

# Redis 缓存配置
CacheRedis:
  - Host: 127.0.0.1:6379
    Pass: ""
    Type: node                   # node/cluster

# 日志配置
Log:
  ServiceName: video-rpc
  Mode: console                  # console/file/volume
  Level: info                    # debug/info/warn/error
  Encoding: plain                # json/plain

# Prometheus 监控
Prometheus:
  Host: 0.0.0.0
  Port: 9081
  Path: /metrics

# Telemetry 链路追踪
Telemetry:
  Name: video-rpc
  Endpoint: http://127.0.0.1:14268/api/traces
  Sampler: 1.0                   # 采样率
  Batcher: jaeger                # jaeger/zipkin
```

**API 服务配置示例**（creative-api.yaml）：

```yaml
Name: creative-api
Host: 0.0.0.0
Port: 8001
Mode: dev

# 日志配置
Log:
  ServiceName: creative-api
  Mode: console
  Level: info

# Prometheus 监控
Prometheus:
  Host: 0.0.0.0
  Port: 9091
  Path: /metrics

# Telemetry 配置
Telemetry:
  Name: creative-api
  Endpoint: http://127.0.0.1:14268/api/traces
  Sampler: 1.0
  Batcher: jaeger

# Video RPC 服务配置
VideoRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:23790
    Key: video.rpc
  Timeout: 10000                 # 超时时间（毫秒）
  NonBlock: true                 # 非阻塞模式

# Hotrank RPC 服务配置
HotrankRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:23790
    Key: hotrank.rpc
  Timeout: 10000
```

### 3.2 配置管理最佳实践

**✅ 推荐做法**：

1. **环境分离**
   ```
   etc/
   ├── service-dev.yaml      # 开发环境
   ├── service-test.yaml     # 测试环境
   └── service-prod.yaml     # 生产环境
   ```

2. **敏感信息管理**
   - 使用环境变量
   - K8s ConfigMap + Secret
   - 不要将密码提交到 Git

3. **配置热更新**（可选）
   ```go
   // 使用 go-zero 的配置热更新
   conf.MustLoad(*configFile, &c, conf.UseEnv())
   ```

**❌ 不推荐做法**：

- ❌ 硬编码配置信息
- ❌ 所有环境共用一个配置文件
- ❌ 将生产环境密码提交到 Git

---

## 4. 服务发现与注册

### 4.1 服务注册方式

**方式一：使用 Etcd（开发/测试环境）**

```yaml
# RPC 配置
Etcd:
  Hosts:
    - 127.0.0.1:23790
  Key: video.rpc
```

**方式二：直连方式（Docker Compose）**

```yaml
# go-zero-looklook 推荐方式
VideoRpc:
  Target: dns:///video-rpc:9001  # 直接指定地址
  # 不使用 Etcd
```

**方式三：K8s Service（生产环境推荐）** ⭐

```yaml
# 不需要 Etcd！K8s 自带服务发现
VideoRpc:
  Target: dns:///video-rpc-svc.default.svc.cluster.local:9001
```

### 4.2 最佳实践对比

| 方式 | 优点 | 缺点 | 适用场景 |
|-----|------|------|---------|
| Etcd | 动态服务发现、支持多实例 | 需要额外维护 Etcd | 传统部署 |
| 直连 | 简单、无额外依赖 | 无法动态扩容 | Docker Compose 开发环境 |
| K8s Service | 原生支持、高可用 | 依赖 K8s | **生产环境推荐** ⭐ |

**实践经验**（来自 go-zero-looklook）：

```
开发环境: Docker Compose + 直连
测试环境: K8s + Service
生产环境: K8s + Service

不再需要 Etcd、Nacos、Consul！
```

### 4.3 监听地址配置技巧

**⚠️ 重要：ListenOn 配置**

```yaml
# ❌ 错误配置（会导致 etcd 注册 0.0.0.0，客户端无法连接）
ListenOn: 0.0.0.0:9001

# ✅ 正确配置
ListenOn: 127.0.0.1:9001     # 本地开发
ListenOn: 0.0.0.0:9001       # K8s 环境（不使用 etcd）
```

**问题案例**（MyBilibili 实际遇到）：

```
问题：hotrank-job 无法连接 video-rpc
原因：video-rpc 监听 0.0.0.0:9001，注册到 etcd 的地址也是 0.0.0.0:9001
解决：改为 127.0.0.1:9001
```

---

## 5. 数据库与缓存

### 5.1 Model 层设计

**使用 goctl 生成 Model**：

```bash
# 方式一：从数据库生成
goctl model mysql datasource \
  -url="root:password@tcp(127.0.0.1:3306)/database" \
  -table="video_info,video_stat" \
  -dir=./model \
  --style go_zero \
  -c  # 启用缓存

# 方式二：从 DDL 生成
goctl model mysql ddl \
  -src="./sql/video_info.sql" \
  -dir=./model \
  --style go_zero \
  -c
```

**Model 层标准结构**：

```go
// video_info_model.go
package model

import (
    "context"
    "database/sql"
    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ VideoInfoModel = (*customVideoInfoModel)(nil)

type (
    // VideoInfoModel 接口定义
    VideoInfoModel interface {
        videoInfoModel
        // 自定义方法
        FindByVids(ctx context.Context, vids []int64) ([]*VideoInfo, error)
    }

    customVideoInfoModel struct {
        *defaultVideoInfoModel
    }

    // VideoInfo 结构体
    VideoInfo struct {
        Id         int64          `db:"id"`
        Vid        int64          `db:"vid"`
        Title      string         `db:"title"`
        Cover      string         `db:"cover"`
        AuthorId   int64          `db:"author_id"`
        AuthorName string         `db:"author_name"`
        RegionId   int64          `db:"region_id"`
        PubTime    int64          `db:"pub_time"`
        Duration   int32          `db:"duration"`
        Desc       sql.NullString `db:"desc"`
        State      int32          `db:"state"`
        Ctime      string         `db:"ctime"`
        Mtime      string         `db:"mtime"`
    }
)

// NewVideoInfoModel 构造函数
func NewVideoInfoModel(conn sqlx.SqlConn, c cache.CacheConf) VideoInfoModel {
    return &customVideoInfoModel{
        defaultVideoInfoModel: newVideoInfoModel(conn, c),
    }
}

// FindByVids 批量查询（自定义方法）⭐
func (m *customVideoInfoModel) FindByVids(ctx context.Context, vids []int64) ([]*VideoInfo, error) {
    if len(vids) == 0 {
        return []*VideoInfo{}, nil
    }

    query := fmt.Sprintf("select %s from %s where `vid` in (?) and `state` = ?",
        videoInfoRows, m.table)
    
    var resp []*VideoInfo
    err := m.QueryRowsNoCacheCtx(ctx, &resp, query, 
        sqlx.In(vids), 
        VideoStateNormal)
    
    return resp, err
}
```

### 5.2 缓存策略

**go-zero 缓存机制**：

```go
// 1. 自动缓存（goctl 生成 -c 参数）
type Config struct {
    zrpc.RpcServerConf
    
    // Redis 缓存配置
    CacheRedis cache.CacheConf
}

// 2. 缓存 Key 命名规范
// cache:video:info:{{vid}}
// cache:video:stat:{{vid}}
// cache:hotrank:list:{{offset}}:{{limit}}

// 3. 缓存失效策略
// - 更新数据时自动失效
// - TTL 过期
// - 手动删除
```

**缓存最佳实践**：

```go
// ✅ 推荐：使用缓存的场景
// 1. 查询频繁、更新不频繁的数据（视频信息）
// 2. 计算复杂的数据（热门排行榜）
// 3. 外部 API 调用结果

// ❌ 不推荐：不适合缓存的场景
// 1. 实时性要求极高的数据（库存、余额）
// 2. 大对象数据（视频文件）
// 3. 频繁更新的数据
```

### 5.3 数据库优化技巧

**1. 批量操作优化**（MyBilibili 实践）：

```go
// ❌ 差：循环更新
for _, arc := range archives {
    dao.UpdateHot(arc.OID, arc.Hot)  // N 次数据库操作
}

// ✅ 好：CASE WHEN 批量更新
sql := `UPDATE academy_archive 
SET hot = CASE oid 
    WHEN ? THEN ?
    WHEN ? THEN ?
    ...
END, mtime = NOW()
WHERE oid IN (?)`

dao.BatchUpdateHot(oidsAndHots)  // 1 次数据库操作
```

**2. 游标分页**（避免深分页）：

```go
// ❌ 差：OFFSET 分页（深分页性能差）
SELECT * FROM video_info 
WHERE state = 0 
ORDER BY id 
LIMIT 1000000, 30;  // 扫描 100 万+ 行

// ✅ 好：游标分页（只扫描需要的行）
SELECT * FROM video_info 
WHERE state = 0 AND id > 12345  // 上次最后一条的 ID
ORDER BY id 
LIMIT 30;  // 只扫描 30 行
```

**3. 索引设计**：

```sql
-- 复合索引（参考主项目 Bilibili）
CREATE INDEX idx_business_state_id ON academy_archive(business, state, id);
CREATE INDEX idx_hot_global ON academy_archive(hot, state);

-- 索引使用原则：
-- 1. WHERE 条件列
-- 2. ORDER BY 排序列  
-- 3. 高频查询字段
```

---

## 6. 错误处理

### 6.1 统一错误码设计

```go
// common/xerr/errCode.go
package xerr

const (
    OK uint32 = 0
    
    // 全局错误码 1-1000
    SERVER_COMMON_ERROR uint32 = 1001
    REQUEST_PARAM_ERROR uint32 = 1002
    DB_ERROR            uint32 = 1003
    
    // 视频服务错误码 10001-11000
    VIDEO_NOT_FOUND      uint32 = 10001
    VIDEO_STAT_NOT_FOUND uint32 = 10002
    
    // 热门排行榜错误码 11001-12000
    HOTRANK_NOT_FOUND uint32 = 11001
    
    // 用户服务错误码 12001-13000
    USER_NOT_FOUND uint32 = 12001
)

// 错误码映射
var mapErrMsg = map[uint32]string{
    OK:                  "SUCCESS",
    SERVER_COMMON_ERROR: "服务器繁忙，请稍后重试",
    REQUEST_PARAM_ERROR: "参数错误",
    DB_ERROR:            "数据库错误",
    
    VIDEO_NOT_FOUND:     "视频不存在",
    HOTRANK_NOT_FOUND:   "排行榜数据不存在",
    USER_NOT_FOUND:      "用户不存在",
}

// NewErrCode 创建错误
func NewErrCode(code uint32) error {
    return &CodeError{
        errCode: code,
        errMsg:  mapErrMsg[code],
    }
}

// NewErrMsg 自定义错误消息
func NewErrMsg(msg string) error {
    return &CodeError{
        errCode: SERVER_COMMON_ERROR,
        errMsg:  msg,
    }
}
```

### 6.2 统一响应封装

```go
// common/result/httpResult.go
package result

import (
    "net/http"
    "github.com/zeromicro/go-zero/rest/httpx"
    "mybilibili/common/xerr"
)

// Response 统一响应结构
type Response struct {
    Code uint32      `json:"code"`
    Msg  string      `json:"msg"`
    Data interface{} `json:"data,omitempty"`
}

// HttpResult 响应处理
func HttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
    if err == nil {
        // 成功响应
        httpx.WriteJson(w, http.StatusOK, Response{
            Code: xerr.OK,
            Msg:  "SUCCESS",
            Data: resp,
        })
        return
    }
    
    // 错误响应
    errCode := xerr.SERVER_COMMON_ERROR
    errMsg := "服务器繁忙"
    
    if e, ok := err.(*xerr.CodeError); ok {
        errCode = e.GetErrCode()
        errMsg = e.GetErrMsg()
    }
    
    httpx.WriteJson(w, http.StatusOK, Response{
        Code: errCode,
        Msg:  errMsg,
        Data: nil,
    })
}
```

### 6.3 错误处理实践

**API 层错误处理**：

```go
// handler/video/get_video_detail_handler.go
func GetVideoDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req types.GetVideoDetailReq
        if err := httpx.Parse(r, &req); err != nil {
            result.HttpResult(r, w, nil, xerr.NewErrCode(xerr.REQUEST_PARAM_ERROR))
            return
        }
        
        l := logic.NewGetVideoDetailLogic(r.Context(), svcCtx)
        resp, err := l.GetVideoDetail(&req)
        result.HttpResult(r, w, resp, err)  // 统一处理
    }
}
```

**Logic 层错误处理**：

```go
// logic/video/get_video_detail_logic.go
func (l *GetVideoDetailLogic) GetVideoDetail(req *types.GetVideoDetailReq) (*types.GetVideoDetailResp, error) {
    // 1. 参数校验
    if req.Vid <= 0 {
        return nil, xerr.NewErrCode(xerr.REQUEST_PARAM_ERROR)
    }
    
    // 2. 调用 RPC
    videoResp, err := l.svcCtx.VideoRpc.GetVideoInfo(l.ctx, &video.GetVideoInfoReq{
        Vid: req.Vid,
    })
    if err != nil {
        // RPC 错误处理
        if status.Code(err) == codes.NotFound {
            return nil, xerr.NewErrCode(xerr.VIDEO_NOT_FOUND)
        }
        return nil, err
    }
    
    // 3. 组装响应
    return &types.GetVideoDetailResp{
        Video: convertToVideoDetail(videoResp),
    }, nil
}
```

**RPC 层错误处理**：

```go
// logic/get_video_info_logic.go
func (l *GetVideoInfoLogic) GetVideoInfo(in *video.GetVideoInfoReq) (*video.GetVideoInfoResp, error) {
    // 1. 查询数据库
    info, err := l.svcCtx.VideoInfoModel.FindOne(l.ctx, in.Vid)
    if err != nil {
        if err == model.ErrNotFound {
            // 返回 gRPC 标准错误码
            return nil, status.Error(codes.NotFound, "video not found")
        }
        l.Errorf("FindOne error: %v", err)
        return nil, status.Error(codes.Internal, "internal error")
    }
    
    // 2. 返回结果
    return &video.GetVideoInfoResp{
        Info: convertToProto(info),
    }, nil
}
```

---

## 7. RPC 服务开发

### 7.1 Protobuf 定义规范

```protobuf
// video.proto
syntax = "proto3";

package video;
option go_package = "./video";

// ==================== 服务定义 ====================
service Video {
  // 获取视频信息
  rpc GetVideoInfo(GetVideoInfoReq) returns(GetVideoInfoResp);
  
  // 批量获取视频信息（⭐高频接口优化）
  rpc BatchGetVideoInfo(BatchGetVideoInfoReq) returns(BatchGetVideoInfoResp);
  
  // 获取视频列表（游标分页）
  rpc GetVideoList(GetVideoListReq) returns(GetVideoListResp);
  
  // 获取视频统计数据
  rpc GetVideoStat(GetVideoStatReq) returns(GetVideoStatResp);
  
  // 批量获取视频统计数据
  rpc BatchGetVideoStat(BatchGetVideoStatReq) returns(BatchGetVideoStatResp);
}

// ==================== 请求/响应定义 ====================

// 获取视频信息请求
message GetVideoInfoReq {
  int64 vid = 1;  // 视频ID
}

message GetVideoInfoResp {
  VideoInfo info = 1;
}

// 批量获取视频信息请求
message BatchGetVideoInfoReq {
  repeated int64 vids = 1;  // 视频ID列表
}

message BatchGetVideoInfoResp {
  map<int64, VideoInfo> infos = 1;  // 使用 map 便于查找
}

// ==================== 数据模型定义 ====================

// 视频信息
message VideoInfo {
  int64 vid = 1;          // 视频ID
  string title = 2;       // 标题
  string cover = 3;       // 封面
  int64 author_id = 4;    // 作者ID
  string author_name = 5; // 作者名称
  int64 region_id = 6;    // 分区ID
  int64 pub_time = 7;     // 发布时间（时间戳）
  int32 duration = 8;     // 时长（秒）
  string desc = 9;        // 简介
  int32 state = 10;       // 状态
}

// 视频统计
message VideoStat {
  int64 vid = 1;        // 视频ID
  int64 view = 2;       // 播放量
  int64 like = 3;       // 点赞数
  int64 coin = 4;       // 投币数
  int64 fav = 5;        // 收藏数
  int64 share = 6;      // 分享数
  int64 reply = 7;      // 评论数
  int64 danmaku = 8;    // 弹幕数
}
```

### 7.2 生成 RPC 代码

```bash
# 1. 生成 gRPC 代码
cd app/video/cmd/rpc
goctl rpc protoc video.proto \
  --go_out=. \
  --go-grpc_out=. \
  --zrpc_out=. \
  --style go_zero

# 生成的文件：
# ├── pb/
# │   ├── video.pb.go           # Protobuf 生成
# │   └── video_grpc.pb.go      # gRPC 生成
# ├── video_client/
# │   └── video.go              # 客户端封装
# ├── internal/
# │   ├── config/config.go
# │   ├── server/video_server.go
# │   ├── svc/service_context.go
# │   └── logic/                # 业务逻辑（手动实现）
# └── video.go                  # 启动文件
```

### 7.3 Service Context 模式 ⭐

**核心理念**：依赖注入，统一管理服务依赖

```go
// internal/svc/service_context.go
package svc

import (
    "mybilibili/app/video/cmd/rpc/internal/config"
    "mybilibili/app/video/model"
    
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
    Config config.Config
    
    // Model 层依赖
    VideoInfoModel model.VideoInfoModel
    VideoStatModel model.VideoStatModel
    
    // 其他 RPC 客户端（如果需要）
    // UserRpc usercenter.UserCenter
}

func NewServiceContext(c config.Config) *ServiceContext {
    // 初始化数据库连接
    conn := sqlx.NewMysql(c.Mysql.DataSource)
    
    return &ServiceContext{
        Config: c,
        
        // 初始化 Model
        VideoInfoModel: model.NewVideoInfoModel(conn, c.CacheRedis),
        VideoStatModel: model.NewVideoStatModel(conn, c.CacheRedis),
    }
}
```

**优势**：

1. ✅ 统一依赖管理
2. ✅ 易于单元测试（Mock）
3. ✅ 清晰的依赖关系
4. ✅ 避免全局变量

### 7.4 Logic 层实现

```go
// internal/logic/batch_get_video_info_logic.go
package logic

import (
    "context"
    "mybilibili/app/video/cmd/rpc/internal/svc"
    "mybilibili/app/video/cmd/rpc/video"
    
    "github.com/zeromicro/go-zero/core/logx"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type BatchGetVideoInfoLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewBatchGetVideoInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetVideoInfoLogic {
    return &BatchGetVideoInfoLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

// BatchGetVideoInfo 批量获取视频信息（⭐hotrank-job 会调用）
func (l *BatchGetVideoInfoLogic) BatchGetVideoInfo(in *video.BatchGetVideoInfoReq) (*video.BatchGetVideoInfoResp, error) {
    // 1. 参数校验
    if len(in.Vids) == 0 {
        return &video.BatchGetVideoInfoResp{
            Infos: make(map[int64]*video.VideoInfo),
        }, nil
    }
    
    // 2. 批量查询数据库（自定义 Model 方法）
    infos, err := l.svcCtx.VideoInfoModel.FindByVids(l.ctx, in.Vids)
    if err != nil {
        l.Errorf("BatchGetVideoInfo FindByVids error: %v", err)
        return nil, status.Error(codes.Internal, "database error")
    }
    
    // 3. 转换为 Proto 结构
    result := make(map[int64]*video.VideoInfo)
    for _, info := range infos {
        result[info.Vid] = &video.VideoInfo{
            Vid:        info.Vid,
            Title:      info.Title,
            Cover:      info.Cover,
            AuthorId:   info.AuthorId,
            AuthorName: info.AuthorName,
            RegionId:   info.RegionId,
            PubTime:    info.PubTime,
            Duration:   info.Duration,
            Desc:       info.Desc.String,
            State:      info.State,
        }
    }
    
    // 4. 日志记录（生产环境建议改为 Debug 级别）
    l.Infof("BatchGetVideoInfo success, count: %d", len(result))
    
    return &video.BatchGetVideoInfoResp{
        Infos: result,
    }, nil
}
```

---

## 8. API 网关开发

### 8.1 API 定义规范

```go
// creative.api
syntax = "v1"

info (
    title:   "创作者 API"
    desc:    "提供热门视频排行榜、视频详情等接口"
    author:  "mybilibili"
    version: "v1.0"
)

// ==================== 数据结构定义 ====================

// 热门排行榜项
type HotRankItem {
    Oid    int64  `json:"oid"`     // 视频ID
    Hot    int64  `json:"hot"`     // 热度值
    Rank   int64  `json:"rank"`    // 排名
    Title  string `json:"title"`   // 视频标题
    Cover  string `json:"cover"`   // 封面
    Author string `json:"author"`  // 作者
    View   int64  `json:"view"`    // 播放量
    Like   int64  `json:"like"`    // 点赞数
}

// 请求/响应定义
type (
    GetHotRankListReq {
        Offset int64 `form:"offset,optional,default=0"`  // 偏移量
        Limit  int64 `form:"limit,optional,default=50"`  // 返回条数
    }
    
    GetHotRankListResp {
        List  []HotRankItem `json:"list"`   // 排行榜列表
        Total int64         `json:"total"`  // 总数
    }
)

// ==================== 路由定义 ====================

@server (
    prefix: /api/creative/v1    // 路由前缀
    group:  hotrank             // 分组
)
service creative {
    @doc "获取全站热门排行榜"
    @handler GetHotRankList
    get /hotrank/list (GetHotRankListReq) returns (GetHotRankListResp)
    
    @doc "获取分区热门排行榜"
    @handler GetRegionHotRankList
    get /hotrank/region (GetRegionHotRankListReq) returns (GetRegionHotRankListResp)
}

@server (
    prefix: /api/creative/v1
    group:  video
)
service creative {
    @doc "获取视频详情"
    @handler GetVideoDetail
    get /video/:vid (GetVideoDetailReq) returns (GetVideoDetailResp)
}
```

### 8.2 生成 API 代码

```bash
cd app/api/creative
goctl api go \
  -api creative.api \
  -dir . \
  --style go_zero

# 生成的文件：
# ├── etc/creative-api.yaml
# ├── internal/
# │   ├── config/config.go
# │   ├── handler/
# │   │   ├── routes.go
# │   │   ├── hotrank/
# │   │   │   ├── get_hot_rank_list_handler.go
# │   │   │   └── get_region_hot_rank_list_handler.go
# │   │   └── video/
# │   │       └── get_video_detail_handler.go
# │   ├── logic/
# │   │   ├── hotrank/
# │   │   │   ├── get_hot_rank_list_logic.go
# │   │   │   └── get_region_hot_rank_list_logic.go
# │   │   └── video/
# │   │       └── get_video_detail_logic.go
# │   ├── svc/service_context.go
# │   └── types/types.go
# └── creative.go
```

### 8.3 API Service Context

```go
// internal/svc/service_context.go
package svc

import (
    "mybilibili/app/api/creative/internal/config"
    "mybilibili/app/video/cmd/rpc/video_client"
    "mybilibili/app/hotrank/cmd/rpc/hotrank_client"
    
    "github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
    Config config.Config
    
    // RPC 客户端依赖
    VideoRpc   video_client.Video
    HotrankRpc hotrank_client.Hotrank
}

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config: c,
        
        // 初始化 RPC 客户端
        VideoRpc:   video_client.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
        HotrankRpc: hotrank_client.NewHotrank(zrpc.MustNewClient(c.HotrankRpc)),
    }
}
```

### 8.4 API Logic 实现（聚合服务）

```go
// internal/logic/hotrank/get_hot_rank_list_logic.go
package hotrank

import (
    "context"
    "mybilibili/app/api/creative/internal/svc"
    "mybilibili/app/api/creative/internal/types"
    "mybilibili/common/xerr"
    
    "github.com/zeromicro/go-zero/core/logx"
    "google.golang.org/grpc/status"
)

type GetHotRankListLogic struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewGetHotRankListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHotRankListLogic {
    return &GetHotRankListLogic{
        Logger: logx.WithContext(ctx),
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

func (l *GetHotRankListLogic) GetHotRankList(req *types.GetHotRankListReq) (*types.GetHotRankListResp, error) {
    // 1. 调用 hotrank-rpc 获取排行榜
    rankResp, err := l.svcCtx.HotrankRpc.GetHotRankList(l.ctx, &hotrank.GetHotRankListReq{
        Offset: req.Offset,
        Limit:  req.Limit,
    })
    if err != nil {
        l.Errorf("HotrankRpc.GetHotRankList error: %v", err)
        return nil, xerr.NewErrMsg("获取排行榜失败")
    }
    
    // 2. 提取所有视频 ID
    vids := make([]int64, 0, len(rankResp.List))
    for _, item := range rankResp.List {
        vids = append(vids, item.Oid)
    }
    
    // 3. 并发调用多个 RPC（性能优化）⭐
    videoInfoChan := make(chan map[int64]*video.VideoInfo, 1)
    videoStatChan := make(chan map[int64]*video.VideoStat, 1)
    errChan := make(chan error, 2)
    
    // 3.1 获取视频信息
    go func() {
        resp, err := l.svcCtx.VideoRpc.BatchGetVideoInfo(l.ctx, &video.BatchGetVideoInfoReq{
            Vids: vids,
        })
        if err != nil {
            errChan <- err
            return
        }
        videoInfoChan <- resp.Infos
    }()
    
    // 3.2 获取视频统计
    go func() {
        resp, err := l.svcCtx.VideoRpc.BatchGetVideoStat(l.ctx, &video.BatchGetVideoStatReq{
            Vids: vids,
        })
        if err != nil {
            errChan <- err
            return
        }
        videoStatChan <- resp.Stats
    }()
    
    // 3.3 等待所有 goroutine 完成
    var videoInfoMap map[int64]*video.VideoInfo
    var videoStatMap map[int64]*video.VideoStat
    
    for i := 0; i < 2; i++ {
        select {
        case err := <-errChan:
            l.Errorf("RPC error: %v", err)
            return nil, xerr.NewErrMsg("获取视频信息失败")
        case videoInfoMap = <-videoInfoChan:
        case videoStatMap = <-videoStatChan:
        }
    }
    
    // 4. 组装响应数据
    result := make([]types.HotRankItem, 0, len(rankResp.List))
    for _, item := range rankResp.List {
        info := videoInfoMap[item.Oid]
        stat := videoStatMap[item.Oid]
        
        if info == nil || stat == nil {
            continue
        }
        
        result = append(result, types.HotRankItem{
            Oid:    item.Oid,
            Hot:    item.Hot,
            Rank:   item.Rank,
            Title:  info.Title,
            Cover:  info.Cover,
            Author: info.AuthorName,
            View:   stat.View,
            Like:   stat.Like,
        })
    }
    
    return &types.GetHotRankListResp{
        List:  result,
        Total: rankResp.Total,
    }, nil
}
```

**性能优化要点**：

1. ✅ 批量 RPC 调用（而非循环调用）
2. ✅ 并发调用多个 RPC（goroutine）
3. ✅ 使用 map 查找（O(1) 复杂度）

---

## 9. 定时任务与消息队列

### 9.1 定时任务实现（Asynq）

**推荐使用**：Asynq（基于 Redis）

```go
// hotrank-job/job.go
package main

import (
    "flag"
    "fmt"
    
    "mybilibili/app/hotrank/cmd/job/internal/config"
    "mybilibili/app/hotrank/cmd/job/internal/svc"
    
    "github.com/zeromicro/go-zero/core/conf"
    "github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/hotrank-job.yaml", "the config file")

func main() {
    flag.Parse()
    
    var c config.Config
    conf.MustLoad(*configFile, &c)
    
    // 初始化服务上下文
    ctx := svc.NewServiceContext(c)
    
    // 启动热度计算任务
    logx.Info("Starting hotrank-job...")
    
    // 阻塞主 goroutine
    logx.Info("hotrank-job started successfully")
    fmt.Println("hotrank-job started successfully, press Ctrl+C to exit")
    select {}
}
```

**Service Context（定时任务版本）**：

```go
// internal/svc/service_context.go
package svc

import (
    "mybilibili/app/hotrank/cmd/job/internal/config"
    "mybilibili/app/hotrank/cmd/job/internal/service"
    "mybilibili/app/video/cmd/rpc/video_client"
    
    "github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
    Config  config.Config
    Service *service.Service
}

func NewServiceContext(c config.Config) *ServiceContext {
    // 初始化 Video RPC 客户端
    videoRpc := video_client.NewVideo(zrpc.MustNewClient(c.VideoRpc))
    
    // 初始化服务（会自动启动热度计算任务）
    svc := service.New(c, videoRpc)
    
    return &ServiceContext{
        Config:  c,
        Service: svc,
    }
}
```

**热度计算核心逻辑**（参考主项目 Bilibili）：

```go
// internal/service/academy.go
package service

import (
    "context"
    "fmt"
    "time"
    
    "github.com/zeromicro/go-zero/core/logx"
)

const (
    batchSize = 30  // 每批处理 30 个视频
)

// FlushHot 热度计算主循环（参考主项目）
func (s *Service) FlushHot(business int64) {
    var lastID int64 = 0
    
    for {
        // 1. 游标分页查询（ID > lastID，避免深分页）
        archives, err := s.academyDao.Archives(lastID, business, batchSize)
        if err != nil {
            logx.Errorf("FlushHot Archives error: %v", err)
            time.Sleep(time.Second * 5)
            continue
        }
        
        // 2. 如果没有更多数据，休眠 1 小时后从头开始
        if len(archives) == 0 {
            logx.Info("FlushHot: no more records, sleep 1 hour and restart from beginning")
            time.Sleep(time.Hour)
            lastID = 0
            continue
        }
        
        // 3. 计算热度并更新
        if err := s.computeHotByOIDs(archives); err != nil {
            logx.Errorf("s.computeHotByOIDs error(%v)", err)
        }
        
        // 4. 更新游标
        lastID = archives[len(archives)-1].ID
        logx.Infof("FlushHot success: processed %d videos, last_id=%d", len(archives), lastID)
        
        // 5. 防止过快循环，休眠 5 秒
        time.Sleep(time.Second * 5)
    }
}

// computeHotByOIDs 计算热度值（完全参考主项目）
func (s *Service) computeHotByOIDs(archives []*model.AcademyArchive) error {
    if len(archives) == 0 {
        return nil
    }
    
    // 1. 提取所有视频 ID
    vids := make([]int64, 0, len(archives))
    for _, arc := range archives {
        vids = append(vids, arc.OID)
    }
    
    // 2. 批量获取视频信息
    videoInfoMap, err := s.videoDao.Archives(context.Background(), vids)
    if err != nil {
        return err
    }
    
    // 3. 批量获取视频统计数据
    videoStatMap, err := s.videoDao.Stats(context.Background(), vids)
    if err != nil {
        return err
    }
    
    // 4. 计算每个视频的热度值
    updates := make(map[int64]int64)  // oid -> hot
    for _, arc := range archives {
        info := videoInfoMap[arc.OID]
        stat := videoStatMap[arc.OID]
        
        if info == nil || stat == nil {
            continue
        }
        
        // 热度计算公式（参考主项目）
        hot := s.countArcHot(stat, info.PubTime)
        updates[arc.OID] = hot
    }
    
    // 5. CASE WHEN 批量更新（一条 SQL 更新多条记录）⭐
    if len(updates) > 0 {
        if err := s.academyDao.UPHotByAIDs(context.Background(), updates); err != nil {
            return err
        }
    }
    
    return nil
}

// countArcHot 热度计算算法（完全参考主项目 Bilibili）
func (s *Service) countArcHot(stat *video.VideoStat, pubTime int64) int64 {
    // 基础热度公式
    hot := float64(stat.Coin)*0.4 + 
           float64(stat.Fav)*0.3 + 
           float64(stat.Danmaku)*0.4 + 
           float64(stat.Reply)*0.4 + 
           float64(stat.View)*0.25 + 
           float64(stat.Like)*0.4 + 
           float64(stat.Share)*0.6
    
    // 新视频提权（24小时内发布的视频热度×1.5）
    now := time.Now().Unix()
    if now-pubTime < 86400 {  // 24小时 = 86400秒
        hot = hot * 1.5
    }
    
    return int64(hot)
}
```

**批量更新 DAO 实现**：

```go
// internal/dao/academy.go
package dao

import (
    "context"
    "fmt"
    "strings"
    "time"
)

// UPHotByAIDs CASE WHEN 批量更新（参考主项目）⭐
func (d *AcademyDao) UPHotByAIDs(ctx context.Context, updates map[int64]int64) error {
    if len(updates) == 0 {
        return nil
    }
    
    // 构建 CASE WHEN 语句
    var caseSQL strings.Builder
    var oids []int64
    var args []interface{}
    
    for oid, hot := range updates {
        caseSQL.WriteString(fmt.Sprintf(" WHEN %d THEN %d", oid, hot))
        oids = append(oids, oid)
    }
    
    // 拼接完整 SQL
    sql := fmt.Sprintf(`
        UPDATE academy_archive 
        SET hot = CASE oid%s END, 
            mtime = ? 
        WHERE oid IN (%s)`,
        caseSQL.String(),
        strings.Repeat(",?", len(oids))[1:],  // 去掉第一个逗号
    )
    
    // 准备参数
    args = append(args, time.Now().Format("2006-01-02 15:04:05"))
    for _, oid := range oids {
        args = append(args, oid)
    }
    
    // 执行 SQL
    result, err := d.conn.ExecCtx(ctx, sql, args...)
    if err != nil {
        return err
    }
    
    affected, _ := result.RowsAffected()
    logx.Infof("AcademyDao.UPHotByAIDs success: updated %d records", affected)
    
    return nil
}
```

### 9.2 消息队列（Kafka + go-queue）

**go-zero-looklook 实践**：

```go
// 1. 配置 Kafka
type Config struct {
    zrpc.RpcServerConf
    
    // Kafka 配置
    PaymentUpdateStatusConf struct {
        Brokers []string
        Topic   string
    }
}

// 2. 发布消息
func (l *UpdateTradeStateLogic) publishPaymentSuccess(orderSn string) error {
    msg := &kqueue.Message{
        Key:   orderSn,
        Value: fmt.Sprintf(`{"order_sn":"%s","trade_state":"SUCCESS"}`, orderSn),
    }
    
    return l.svcCtx.KqPusher.Push(msg)
}

// 3. 消费消息
type PaymentSuccessHandler struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func (h *PaymentSuccessHandler) Consume(key, val string) error {
    logx.Infof("PaymentSuccess consume, key: %s, val: %s", key, val)
    
    // 解析消息
    var msg PaymentSuccessMsg
    if err := json.Unmarshal([]byte(val), &msg); err != nil {
        return err
    }
    
    // 处理业务逻辑
    return h.updateOrderStatus(msg.OrderSn, msg.TradeState)
}
```

---

## 10. 链路追踪与监控

### 10.1 Jaeger 链路追踪

**配置方式**：

```yaml
# RPC/API 配置
Telemetry:
  Name: video-rpc
  Endpoint: http://127.0.0.1:14268/api/traces
  Sampler: 1.0        # 采样率（1.0 = 100%）
  Batcher: jaeger     # jaeger/zipkin
```

**自动追踪**：
- ✅ HTTP 请求
- ✅ gRPC 调用
- ✅ 数据库查询
- ✅ Redis 操作

**查看链路**：
```
Jaeger UI: http://localhost:16686
```

### 10.2 Prometheus 监控

**配置方式**：

```yaml
# RPC/API 配置
Prometheus:
  Host: 0.0.0.0
  Port: 9081          # 每个服务不同端口
  Path: /metrics
```

**自动暴露指标**：
- `http_server_requests_code_total` - HTTP 请求统计
- `http_server_requests_duration_ms` - 请求耗时
- `rpc_server_requests_code_total` - RPC 请求统计
- `rpc_server_requests_duration_ms` - RPC 耗时

**Prometheus 配置**：

```yaml
# deploy/prometheus/prometheus.yml
scrape_configs:
  - job_name: 'video-rpc'
    static_configs:
      - targets: ['video-rpc:9081']
  
  - job_name: 'hotrank-rpc'
    static_configs:
      - targets: ['hotrank-rpc:9093']
  
  - job_name: 'creative-api'
    static_configs:
      - targets: ['creative-api:9091']
```

### 10.3 Grafana 可视化

```
访问地址: http://localhost:3000
默认账号: admin
默认密码: admin123456

推荐导入 go-zero 官方仪表板：
https://grafana.com/grafana/dashboards/
```

---

## 11. 性能优化

### 11.1 数据库优化

**1. 批量操作**

```go
// ❌ 差：N+1 查询
for _, vid := range vids {
    info, _ := model.FindOne(vid)  // N 次数据库查询
}

// ✅ 好：一次批量查询
infos, _ := model.FindByVids(vids)  // 1 次数据库查询
```

**2. 索引优化**

```sql
-- 查询条件索引
CREATE INDEX idx_state_id ON video_info(state, id);

-- 复合索引（覆盖索引）
CREATE INDEX idx_business_state_id ON academy_archive(business, state, id);

-- 排序索引
CREATE INDEX idx_hot_state ON academy_archive(hot DESC, state);
```

**3. 游标分页 vs OFFSET 分页**

```sql
-- ❌ OFFSET 分页（深分页性能差）
SELECT * FROM video_info 
WHERE state = 0 
LIMIT 100000, 30;  -- 扫描 100030 行

-- ✅ 游标分页（只扫描需要的行）
SELECT * FROM video_info 
WHERE state = 0 AND id > 100000  -- 游标
LIMIT 30;  -- 只扫描 30 行
```

### 11.2 缓存优化

**1. 缓存击穿防护**

```go
// go-zero 自带 singleflight 机制
// 多个并发请求只会触发一次数据库查询
info, err := model.FindOne(ctx, vid)
```

**2. 缓存预热**

```go
// 启动时预加载热门数据
func (s *Service) PreloadCache() {
    // 预加载 TOP100 排行榜
    s.GetHotRankList(0, 100)
}
```

**3. 缓存更新策略**

```go
// 写入时更新缓存
func (m *VideoInfoModel) Update(ctx context.Context, data *VideoInfo) error {
    // 1. 更新数据库
    err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) error {
        sql := fmt.Sprintf("update %s set ... where vid = ?", m.table)
        _, err := conn.ExecCtx(ctx, sql, ...)
        return err
    }, m.cacheKeys(data.Vid)...)  // 2. 自动删除缓存
    
    return err
}
```

### 11.3 RPC 调用优化

**1. 批量调用**

```go
// ❌ 差：循环调用 RPC
for _, vid := range vids {
    info, _ := videoRpc.GetVideoInfo(ctx, &video.GetVideoInfoReq{Vid: vid})
}

// ✅ 好：批量调用
infos, _ := videoRpc.BatchGetVideoInfo(ctx, &video.BatchGetVideoInfoReq{Vids: vids})
```

**2. 并发调用**

```go
// 并发调用多个 RPC
var wg sync.WaitGroup
var infoMap map[int64]*video.VideoInfo
var statMap map[int64]*video.VideoStat

wg.Add(2)

go func() {
    defer wg.Done()
    resp, _ := videoRpc.BatchGetVideoInfo(ctx, &video.BatchGetVideoInfoReq{Vids: vids})
    infoMap = resp.Infos
}()

go func() {
    defer wg.Done()
    resp, _ := videoRpc.BatchGetVideoStat(ctx, &video.BatchGetVideoStatReq{Vids: vids})
    statMap = resp.Stats
}()

wg.Wait()
```

**3. 超时控制**

```yaml
# RPC 客户端配置
VideoRpc:
  Timeout: 10000      # 10 秒超时
  NonBlock: true      # 非阻塞模式
```

---

## 12. 部署策略

### 12.1 开发环境（Docker Compose）

**docker-compose.yml**：

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: mybilibili-mysql
    ports:
      - "33060:3306"
    environment:
      MYSQL_ROOT_PASSWORD: root123456
      MYSQL_DATABASE: mybilibili
    volumes:
      - ./data/mysql:/var/lib/mysql
      - ./sql:/docker-entrypoint-initdb.d
    networks:
      - mybilibili-net

  redis:
    image: redis:7-alpine
    container_name: mybilibili-redis
    ports:
      - "63790:6379"
    command: redis-server --requirepass redis123456
    networks:
      - mybilibili-net

  etcd:
    image: quay.io/coreos/etcd:v3.5.9
    container_name: mybilibili-etcd
    ports:
      - "23790:2379"
    environment:
      - ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379
      - ETCD_ADVERTISE_CLIENT_URLS=http://0.0.0.0:2379
    networks:
      - mybilibili-net

  jaeger:
    image: jaegertracing/all-in-one:latest
    container_name: mybilibili-jaeger
    ports:
      - "16686:16686"  # UI
      - "14268:14268"  # Collector HTTP
    networks:
      - mybilibili-net

  prometheus:
    image: prom/prometheus:latest
    container_name: mybilibili-prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
    networks:
      - mybilibili-net

  grafana:
    image: grafana/grafana:latest
    container_name: mybilibili-grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin123456
    networks:
      - mybilibili-net

networks:
  mybilibili-net:
    driver: bridge
```

**启动命令**：

```bash
# 启动基础设施
docker-compose up -d

# 初始化数据库
docker exec -i mybilibili-mysql mysql -uroot -proot123456 < sql/001_init.sql
docker exec -i mybilibili-mysql mysql -uroot -proot123456 mybilibili < sql/002_test_data.sql

# 启动微服务
cd app/video/cmd/rpc && ./video-rpc -f etc/video.yaml &
cd app/hotrank/cmd/rpc && ./hotrank-rpc -f etc/hotrank.yaml &
cd app/hotrank/cmd/job && ./hotrank-job -f etc/hotrank-job.yaml &
cd app/api/creative && ./creative-api -f etc/creative-api.yaml &
```

### 12.2 生产环境（K8s）⭐

**推荐架构**（来自 go-zero-looklook）：

```
┌─────────────┐
│   SLB/ALB   │  ← 负载均衡
└──────┬──────┘
       │
┌──────▼──────┐
│    Nginx    │  ← 网关
└──────┬──────┘
       │
┌──────▼──────────────────────────┐
│        K8s Cluster               │
│  ┌────────┐  ┌────────┐         │
│  │ API Pod│  │ API Pod│         │
│  └───┬────┘  └───┬────┘         │
│      │           │               │
│  ┌───▼────┐  ┌──▼─────┐         │
│  │RPC Pod │  │RPC Pod │         │
│  └───┬────┘  └───┬────┘         │
│      │           │               │
│  ┌───▼───────────▼────┐         │
│  │ MySQL / Redis      │         │
│  └────────────────────┘         │
└──────────────────────────────────┘
```

**K8s Deployment 示例**：

```yaml
# video-rpc-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: video-rpc
  namespace: mybilibili
spec:
  replicas: 3  # 3 个副本
  selector:
    matchLabels:
      app: video-rpc
  template:
    metadata:
      labels:
        app: video-rpc
    spec:
      containers:
      - name: video-rpc
        image: mybilibili/video-rpc:v1.0.0
        ports:
        - containerPort: 9001
        env:
        - name: MYSQL_HOST
          value: "mysql-svc"
        - name: REDIS_HOST
          value: "redis-svc"
        resources:
          requests:
            cpu: "100m"
            memory: "128Mi"
          limits:
            cpu: "500m"
            memory: "512Mi"
        livenessProbe:
          tcpSocket:
            port: 9001
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          tcpSocket:
            port: 9001
          initialDelaySeconds: 5
          periodSeconds: 5

---
# video-rpc-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: video-rpc-svc
  namespace: mybilibili
spec:
  selector:
    app: video-rpc
  ports:
  - port: 9001
    targetPort: 9001
  type: ClusterIP
```

**ConfigMap 管理配置**：

```yaml
# video-rpc-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: video-rpc-config
  namespace: mybilibili
data:
  video.yaml: |
    Name: video.rpc
    ListenOn: 0.0.0.0:9001
    Mode: prod
    
    Mysql:
      DataSource: ${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:3306)/mybilibili
    
    CacheRedis:
      - Host: ${REDIS_HOST}:6379
        Pass: ${REDIS_PASSWORD}
        Type: node
```

**部署流程**（Jenkins + GitLab + Harbor + K8s）：

```
1. 开发推送代码到 GitLab
2. GitLab Webhook 触发 Jenkins
3. Jenkins 拉取代码 + 配置
4. Docker 构建镜像
5. 推送镜像到 Harbor
6. kubectl apply 部署到 K8s
7. K8s 滚动更新
```

### 12.3 部署最佳实践

**✅ 推荐做法**：

1. **环境隔离**
   - 开发环境：Docker Compose
   - 测试环境：K8s
   - 生产环境：K8s + 多集群

2. **配置管理**
   - K8s ConfigMap（普通配置）
   - K8s Secret（敏感信息）
   - 不使用配置中心（简化架构）

3. **服务发现**
   - 不需要 Etcd/Nacos/Consul
   - 使用 K8s Service 原生服务发现

4. **监控告警**
   - Prometheus 采集指标
   - Grafana 可视化
   - AlertManager 告警

5. **日志收集**
   - Filebeat → Kafka → ES
   - 或使用 ELK/EFK Stack

**❌ 不推荐做法**：

- ❌ 所有环境共用配置
- ❌ 手动部署（容易出错）
- ❌ 没有监控告警
- ❌ 没有日志收集

---

## 13. 常见问题与解决方案

### 13.1 Etcd 连接问题

**问题**：
```
rpc error: code = DeadlineExceeded desc = context deadline exceeded
error reading server preface: EOF
```

**原因**：
- 服务监听 `0.0.0.0:9001`
- 注册到 etcd 的地址也是 `0.0.0.0:9001`
- 客户端无法连接到 `0.0.0.0`

**解决方案**：

```yaml
# ✅ 方案一：使用 127.0.0.1（本地开发）
ListenOn: 127.0.0.1:9001

# ✅ 方案二：使用 K8s Service（生产环境推荐）
# 不需要 etcd！
VideoRpc:
  Target: dns:///video-rpc-svc:9001
```

### 13.2 Go 编译问题

**问题**：
```
package go/build/constraint is not in std
```

**原因**：Go 标准库文件缺失或版本不兼容

**解决方案**：

```bash
# 方案一：使用已编译的二进制文件
./video-rpc -f etc/video.yaml

# 方案二：重新编译
go build -o video-rpc video.go

# 方案三：检查 Go 版本
go version  # 推荐 1.21+
```

### 13.3 API 404 错误

**问题**：
```
404 page not found
```

**原因**：路由配置不正确

**解决方案**：

```bash
# 检查 API 定义文件
# creative.api
@server (
    prefix: /api/creative/v1  # 注意路由前缀
    group:  hotrank
)

# 正确的请求 URL
curl "http://127.0.0.1:8001/api/creative/v1/hotrank/list?offset=0&limit=10"

# ❌ 错误的 URL
curl "http://127.0.0.1:8001/api/v1/academy/rank/global"
```

### 13.4 代理问题

**问题**：
```
curl: (7) Failed to connect to 127.0.0.1 port 8001
```

**原因**：系统代理设置影响本地访问

**解决方案**：

```bash
# 临时取消代理
unset http_proxy https_proxy HTTP_PROXY HTTPS_PROXY ALL_PROXY all_proxy

# 然后访问本地服务
curl "http://127.0.0.1:8001/api/creative/v1/hotrank/list"
```

### 13.5 热度计算不生效

**问题**：热度值一直为 0

**排查步骤**：

```bash
# 1. 检查 hotrank-job 是否运行
ps aux | grep hotrank-job

# 2. 查看日志
tail -f /path/to/hotrank-job.log

# 3. 检查数据库
mysql> SELECT id, oid, hot FROM academy_archive ORDER BY hot DESC LIMIT 5;

# 4. 手动触发计算（重启 job）
kill <pid>
./hotrank-job -f etc/hotrank-job.yaml
```

---

## 附录

### A. 项目对比总结

| 特性 | MyBilibili | go-zero-looklook |
|-----|-----------|------------------|
| 业务领域 | 视频排行榜 | 电商旅游 |
| 服务数量 | 3 个 RPC + 1 个 API | 5+ 个领域服务 |
| 技术栈 | go-zero + MySQL + Redis | go-zero 全家桶 |
| 服务发现 | Etcd / K8s | 直连 / K8s |
| 消息队列 | 无 | Kafka + go-queue |
| 分布式事务 | 无 | DTM |
| 定时任务 | 自实现 | Asynq |
| 部署方式 | Docker Compose / K8s | K8s |

### B. 核心设计模式

1. **Service Context 模式**：统一依赖管理
2. **批量操作模式**：减少数据库/RPC 调用
3. **游标分页模式**：避免深分页性能问题
4. **CASE WHEN 批量更新**：一条 SQL 更新多行
5. **并发 RPC 调用**：提升聚合服务性能

### C. 学习资源

- **go-zero 官方文档**：https://go-zero.dev
- **go-zero-looklook**：https://github.com/Mikaelemmmm/go-zero-looklook
- **MyBilibili**：本项目
- **go-zero 社区**：微信群、Discord

---

## 总结

本文结合 **MyBilibili** 和 **go-zero-looklook** 两个实战项目，总结了 go-zero 微服务开发的最佳实践。

**核心要点**：

1. ✅ 清晰的分层架构（API → RPC → Model）
2. ✅ Service Context 依赖注入模式
3. ✅ 统一的错误处理和响应封装
4. ✅ 批量操作优化性能
5. ✅ K8s 原生服务发现（不需要 Etcd）
6. ✅ 完善的监控和链路追踪
7. ✅ 自动化部署流程

**建议学习路径**：

```
1. 学习 go-zero 基础（官方文档）
2. 运行 MyBilibili（理解基础架构）
3. 研究 go-zero-looklook（学习高级特性）
4. 实战自己的项目
```

希望本文能帮助你更好地使用 go-zero 开发微服务！🚀

---

**版权声明**：本文档为开源文档，欢迎转载和分享，请注明出处。

**贡献者**：mybilibili 团队、go-zero 社区

**最后更新**：2025-11-08


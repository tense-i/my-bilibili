# MyBilibili AI 开发助手指南

> **版本**: v1.3.0  
> **更新时间**: 2025-12-24  
> **项目**: MyBilibili - 基于 go-zero 的微服务重构项目
> **重要更新**: 新增 Swagger 文档生成和 gRPC 反射调试配置

---

## 📋 项目概述

MyBilibili 是从 openbilibili 主项目模仿实现的微服务系统，使用 go-zero 框架重构。

### 项目结构
```
mybilibili/
├── app/                          # 业务服务
│   ├── api/                      # HTTP 网关
│   ├── creative/                 # 创作中心 API
│   ├── hotrank/                  # 热门排行榜（已重构）
│   ├── recall/                   # 召回服务
│   ├── recommend/                # 推荐服务（已重构）
│   ├── video/                    # 视频服务（已重构）
│   └── wallet/                   # 虚拟钱包（已重构）⭐参考模板
├── common/                       # 公共组件
│   ├── model/                    # 公共数据模型
│   ├── result/                   # 统一响应
│   ├── tool/                     # 工具函数
│   └── xerr/                     # 错误码
├── deploy/                       # 部署配置
│   ├── sql/                      # SQL 脚本
│   └── docker-compose.yml
├── doc/                          # 文档
└── 设计与修复方案/                # 设计文档
```

### 参考项目
- **主项目**: `openbilibili/` - 业务逻辑和表结构的唯一参考来源
- **重构模板**: `mybilibili/app/wallet/` - go-zero 最佳实践参考

---

## 🎯 核心开发准则

### 准则一：严格遵循 go-zero 最佳实践

**微服务分层架构**：
```
API 服务层（HTTP）
    ↓ gRPC
RPC 服务层（核心业务）
    ↓
Model 层（数据访问）
    ↓
MySQL / Redis / Kafka
```

**标准目录结构**（以 wallet 为模板）：
```
app/{module}/
├── cmd/
│   ├── api/                    # HTTP API 服务
│   │   ├── desc/{module}.api   # API 定义文件
│   │   ├── etc/{module}.yaml   # 配置文件
│   │   └── internal/
│   │       ├── config/         # 配置结构
│   │       ├── handler/        # HTTP Handler
│   │       ├── logic/          # 业务逻辑
│   │       ├── svc/            # 服务上下文
│   │       └── types/          # 类型定义
│   └── rpc/                    # gRPC 服务
│       ├── {module}.proto      # Protobuf 定义
│       ├── etc/{module}.yaml   # 配置文件
│       └── internal/
│           ├── config/         # 配置结构
│           ├── logic/          # 核心业务逻辑 ⭐
│           ├── server/         # gRPC 服务器
│           └── svc/            # 服务上下文
└── model/                      # 数据模型
```

### 准则二：业务实现完全参照 openbilibili

**⚠️ 重要：不要私自改动业务逻辑**

1. **查找主项目对应模块**：
   ```
   openbilibili/app/service/main/{module}/    # 主站服务
   openbilibili/app/service/live/{module}/    # 直播服务
   openbilibili/app/interface/main/{module}/  # 接口层
   ```

2. **参考内容**：
   - 业务流程和算法
   - 数据校验规则
   - 错误处理逻辑
   - 常量和枚举定义

3. **示例**：钱包模块参考
   ```
   主项目: openbilibili/app/service/live/wallet/
   重构后: mybilibili/app/wallet/
   ```

### 准则三：表结构完全参照主项目

**数据库设计原则**：
- 表名、字段名、类型完全一致
- 索引设计参照主项目
- 分表策略保持一致

**SQL 脚本位置**：`deploy/sql/`

**命名规范**：
```
001_init.sql              # 基础表
002_test_data.sql         # 测试数据
010_wallet_schema.sql     # 钱包模块
03_recommend.sql          # 推荐模块
```

### 准则四：开发前输出技术设计文档

**每次重构新模块前，必须先输出设计文档**，包含：

1. **概述与架构**
   - 项目目标和要求
   - 微服务划分
   - 目录结构设计
   - 核心设计（锁机制、分表策略等）

2. **接口与实现**
   - Protobuf 接口定义
   - API 接口定义
   - 核心 Logic 实现代码
   - Model 层实现

3. **实施计划**
   - 分阶段实施步骤
   - 测试方案
   - 部署方案

**设计文档位置**：`设计与修复方案/`

**命名规范**：
```
v{版本}-{模块名}-总览.md
v{版本}-{模块名}-第1部分-概述与架构.md
v{版本}-{模块名}-第2部分-接口与实现.md
v{版本}-{模块名}-第3部分-实施计划.md
```

### 准则五：技术栈参考已重构模块

**核心技术栈**：
- **框架**: go-zero v1.6+
- **数据库**: MySQL 8.0
- **缓存**: Redis 7
- **服务发现**: etcd
- **消息队列**: Kafka
- **链路追踪**: Jaeger
- **监控**: Prometheus + Grafana

**参考模块**：
- `app/wallet/` - 虚拟钱包（分布式事务、双重锁、分表）
- `app/recommend/` - 推荐系统（召回、排序、特征工程）
- `app/hotrank/` - 热门排行（定时任务、批量更新）

### 准则六：必须使用 /Users/zh/go/bin/goctl 工具生成代码 ⭐⭐⭐

**� 强制要求：所有代码必须通过 /Users/zh/go/bin/goctl 工具生成，禁止手写框架代码！**

**为什么必须使用 /Users/zh/go/bin/goctl**：
1. 保证代码结构与 go-zero 框架完全兼容
2. 避免手写代码导致的格式、命名不一致
3. 自动生成 handler、logic、types、server 等样板代码
4. 确保 gRPC 和 HTTP 接口定义的一致性

**/Users/zh/go/bin/goctl 代码生成命令**：

```bash
# ==================== RPC 服务代码生成 ====================
# 进入 RPC 目录
cd app/{module}/cmd/rpc

# 生成 gRPC 代码（必须使用此命令）
/Users/zh/go/bin/goctl rpc protoc {module}.proto \
    --go_out=. \
    --go-grpc_out=. \
    --zrpc_out=. \
    --style go_zero

# ==================== API 服务代码生成 ====================
# 进入 API 目录
cd app/{module}/cmd/api

# 生成 HTTP API 代码（必须使用此命令）
/Users/zh/go/bin/goctl api go \
    -api desc/{module}.api \
    -dir . \
    --style go_zero

# ==================== Model 代码生成（可选）====================
# 从 SQL 文件生成 Model
/Users/zh/go/bin/goctl model mysql ddl \
    -src {table}.sql \
    -dir ./model \
    -c \
    --style go_zero

# 从数据库连接生成 Model
/Users/zh/go/bin/goctl model mysql datasource \
    -url "root:password@tcp(127.0.0.1:3306)/dbname" \
    -table "{table_name}" \
    -dir ./model \
    -c \
    --style go_zero
```

**开发工作流程**：

```
1. 编写 .proto 文件（定义 RPC 接口）
       ↓
2. 执行 /Users/zh/go/bin/goctl rpc protoc 生成 RPC 代码
       ↓
3. 编写 .api 文件（定义 HTTP 接口）
       ↓
4. 执行 /Users/zh/go/bin/goctl api go 生成 API 代码
       ↓
5. 只修改 logic 层实现业务逻辑
       ↓
6. 只修改 svc 层添加依赖注入
```

**⚠️ 禁止手写的文件**：
- `internal/handler/*.go` - 必须由 /Users/zh/go/bin/goctl 生成
- `internal/types/types.go` - 必须由 /Users/zh/go/bin/goctl 生成
- `internal/server/*_server.go` - 必须由 /Users/zh/go/bin/goctl 生成
- `{module}.go` (main 入口) - 必须由 /Users/zh/go/bin/goctl 生成
- `*_grpc.pb.go` / `*.pb.go` - 必须由 /Users/zh/go/bin/goctl 生成

**✅ 允许手动修改的文件**：
- `internal/logic/*.go` - 业务逻辑实现
- `internal/svc/service_context.go` - 服务上下文
- `internal/config/config.go` - 配置结构扩展
- `etc/*.yaml` - 配置文件
- `model/*.go` - 数据模型（或用 /Users/zh/go/bin/goctl 生成后修改）

**接口变更时的操作**：

```bash
# 当 .proto 文件有变更时，重新生成
cd app/{module}/cmd/rpc
/Users/zh/go/bin/goctl rpc protoc {module}.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style go_zero

# 当 .api 文件有变更时，重新生成
cd app/{module}/cmd/api
/Users/zh/go/bin/goctl api go -api desc/{module}.api -dir . --style go_zero

# 注意：重新生成会覆盖 handler/types/server，但不会覆盖 logic 和 svc
```

---

## 📝 开发流程

### 阶段一：需求分析与设计

```markdown
1. 确认要重构的模块
2. 在 openbilibili 中找到对应模块
3. 分析业务逻辑和数据结构
4. 输出技术设计文档（必须）
5. 等待确认后开始实施
```

### 阶段二：基础架构搭建

```bash
# 1. 创建目录结构
mkdir -p app/{module}/cmd/api/desc
mkdir -p app/{module}/cmd/api/etc
mkdir -p app/{module}/cmd/rpc/etc
mkdir -p app/{module}/model

# 2. 编写 Protobuf 定义
# app/{module}/cmd/rpc/{module}.proto

# 3. 生成 RPC 代码
cd app/{module}/cmd/rpc
/Users/zh/go/bin//Users/zh/go/bin/goctl rpc protoc {module}.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style go_zero

# 4. 编写 API 定义
# app/{module}/cmd/api/desc/{module}.api

# 5. 生成 API 代码
cd app/{module}/cmd/api
/Users/zh/go/bin//Users/zh/go/bin/goctl api go -api desc/{module}.api -dir . --style go_zero

# 6. 创建数据库表
# deploy/sql/0XX_{module}_schema.sql
```

### 阶段三：核心功能实现

```markdown
1. 实现 Model 层（参照主项目表结构）
2. 实现 RPC Service Context
3. 实现 RPC Logic（核心业务逻辑）
4. 实现 API Service Context
5. 实现 API Logic（聚合调用）
6. 编写配置文件
```

### 阶段四：测试与验证

```bash
# 1. 启动基础服务
cd deploy && docker-compose up -d

# 2. 初始化数据库
mysql -h127.0.0.1 -P33060 -uroot -proot123456 < deploy/sql/0XX_{module}_schema.sql

# 3. 启动 RPC 服务
cd app/{module}/cmd/rpc && go run {module}.go -f etc/{module}.yaml

# 4. 启动 API 服务
cd app/{module}/cmd/api && go run {module}.go -f etc/{module}.yaml

# 5. 测试接口
curl http://localhost:800X/api/{module}/v1/xxx
```

### 阶段五：API 文档与调试配置 ⭐

完成功能开发后，必须配置 Swagger 文档和 gRPC 反射，方便使用 Apifox 等工具进行接口调试。

#### 5.1 生成 Swagger 文档（HTTP API）

```bash
# 1. 安装 goctl-swagger 插件（首次使用）
go install github.com/zeromicro/goctl-swagger@latest

# 2. 进入 API 目录，创建 doc 目录
cd app/{module}/cmd/api
mkdir -p doc

# 3. 生成 Swagger JSON 文档
/Users/zh/go/bin/goctl api plugin \
    -plugin /Users/zh/go/bin/goctl-swagger="swagger -filename {module}.json" \
    -api desc/{module}.api \
    -dir doc
```

#### 5.2 暴露 Swagger 文档 URL

修改 API 服务入口文件 `{module}.go`，添加静态文件路由：

```go
// 在 handler.RegisterHandlers(server, ctx) 之后添加
import "net/http"

// 添加 swagger.json 静态文件路由，供 Apifox 通过 URL 导入
server.AddRoute(rest.Route{
    Method: http.MethodGet,
    Path:   "/doc/{module}/swagger.json",
    Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "doc/{module}.json")
    }),
})
```

**Apifox 导入方式**：
- URL 导入地址：`http://127.0.0.1:{port}/doc/{module}/swagger.json`
- 示例：`http://127.0.0.1:8085/doc/coupon/swagger.json`

#### 5.3 启用 gRPC 反射（RPC 服务）

**步骤1**：确认 RPC 入口文件已包含反射代码

```go
// app/{module}/cmd/rpc/{module}.go
import (
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    "github.com/zeromicro/go-zero/core/service"
)

func main() {
    // ...
    s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
        {module}.Register{Module}Server(grpcServer, server.New{Module}Server(ctx))
        
        // 开发/测试环境启用 gRPC 反射
        if c.Mode == service.DevMode || c.Mode == service.TestMode {
            reflection.Register(grpcServer)
        }
    })
    // ...
}
```

**步骤2**：在配置文件中设置 Mode 为 dev

```yaml
# app/{module}/cmd/rpc/etc/{module}.yaml
Name: {module}.rpc
ListenOn: 127.0.0.1:90XX
Mode: dev  # 启用开发模式，开启 gRPC 反射

Etcd:
  Hosts:
    - 127.0.0.1:23790
  Key: {module}.rpc
# ...
```

**步骤3**：验证 gRPC 反射是否启用

```bash
# 安装 grpcurl（首次使用）
brew install grpcurl

# 列出所有服务
grpcurl -plaintext 127.0.0.1:90XX list

# 列出服务的所有方法
grpcurl -plaintext 127.0.0.1:90XX list {module}.{Module}

# 测试调用方法
grpcurl -plaintext -d '{"mid": 1001}' 127.0.0.1:90XX {module}.{Module}/MethodName
```

**Apifox gRPC 调试**：
1. 新建 gRPC 项目 → 导入 → 选择"服务器反射"
2. 输入地址：`127.0.0.1:90XX`
3. 点击"获取服务"自动发现所有方法

#### 5.4 服务端口规划

| 模块 | API 端口 | RPC 端口 | Swagger URL |
|------|---------|---------|-------------|
| wallet | 8004 | 9004 | `/doc/wallet/swagger.json` |
| coupon | 8085 | 9085 | `/doc/coupon/swagger.json` |
| search | 8006 | 9006 | `/doc/search/swagger.json` |
| video | 8001 | 9001 | `/doc/video/swagger.json` |

---

## 🔧 代码规范

### Protobuf 定义规范

```protobuf
syntax = "proto3";

package {module};
option go_package = "./{module}";

// ==================== 服务定义 ====================
service {Module} {
  // 方法注释
  rpc MethodName(MethodReq) returns(MethodResp);
}

// ==================== 请求/响应定义 ====================
message MethodReq {
  int64 uid = 1;                    // 字段注释
  string field = 2;                 // 字段注释
}

message MethodResp {
  int64 code = 1;                   // 响应码
  string msg = 2;                   // 响应消息
}
```

### API 定义规范

```go
syntax = "v1"

info (
    title:   "{模块名} API"
    desc:    "{模块描述}"
    author:  "mybilibili"
    version: "v1.0.0"
)

// ==================== 数据结构定义 ====================
type MethodReq {
    Uid    int64  `json:"uid"`
    Field  string `json:"field"`
}

type MethodResp {
    Data interface{} `json:"data"`
}

// ==================== 路由定义 ====================
@server (
    prefix: /api/{module}/v1
    group:  {group}
)
service {module} {
    @doc "接口描述"
    @handler MethodName
    post /path (MethodReq) returns (MethodResp)
}
```

### Service Context 规范

```go
package svc

import (
    "mybilibili/app/{module}/cmd/rpc/internal/config"
    "mybilibili/app/{module}/model"
    
    "github.com/zeromicro/go-zero/core/stores/redis"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
    Config config.Config
    
    // MySQL 连接
    DB sqlx.SqlConn
    
    // Redis 连接
    Redis *redis.Redis
    
    // Model 层
    XxxModel model.XxxModel
}

func NewServiceContext(c config.Config) *ServiceContext {
    db := sqlx.NewMysql(c.Mysql.DataSource)
    
    rds := redis.New(c.RedisConf.Host, func(r *redis.Redis) {
        r.Type = c.RedisConf.Type
        r.Pass = c.RedisConf.Pass
    })
    
    return &ServiceContext{
        Config:   c,
        DB:       db,
        Redis:    rds,
        XxxModel: model.NewXxxModel(db, c.CacheRedis),
    }
}
```

### Logic 实现规范

```go
package logic

import (
    "context"
    
    "mybilibili/app/{module}/cmd/rpc/internal/svc"
    "mybilibili/app/{module}/cmd/rpc/{module}"
    "mybilibili/common/xerr"
    
    "github.com/zeromicro/go-zero/core/logx"
)

type XxxLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewXxxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *XxxLogic {
    return &XxxLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *XxxLogic) Xxx(in *{module}.XxxReq) (*{module}.XxxResp, error) {
    // 1. 参数校验
    if in.Uid <= 0 {
        return nil, xerr.NewErrCode(xerr.REQUEST_PARAM_ERROR)
    }
    
    // 2. 业务逻辑（参照主项目实现）
    // ...
    
    // 3. 返回结果
    return &{module}.XxxResp{}, nil
}
```

---

## 📚 参考资源

### 已重构模块（可作为模板）

| 模块 | 路径 | 特点 |
|-----|------|------|
| wallet | `app/wallet/` | 分布式事务、双重锁、分表、Kafka |
| recommend | `app/recommend/` | 召回、排序、特征工程 |
| hotrank | `app/hotrank/` | 定时任务、批量更新、游标分页 |
| video | `app/video/` | 基础 CRUD、缓存策略 |

### 主项目模块对照

| 主项目路径 | 说明 |
|-----------|------|
| `openbilibili/app/service/main/` | 主站核心服务 |
| `openbilibili/app/service/live/` | 直播相关服务 |
| `openbilibili/app/interface/main/` | 接口层 |
| `openbilibili/app/admin/` | 管理后台 |
| `openbilibili/app/job/` | 定时任务 |

### 文档资源

- `doc/GO-ZERO-BEST-PRACTICES.md` - go-zero 最佳实践
- `设计与修复方案/` - 各模块设计文档
- `gozero重构准则.md` - 重构准则

---

## ⚠️ 注意事项

### 必须遵守

1. **不要私自修改业务逻辑** - 完全参照 openbilibili
2. **不要简化方案** - 保持与主项目一致的复杂度
3. **先输出设计文档** - 每次重构前必须先设计
4. **表结构完全一致** - 不要修改字段名和类型
5. **必须使用 /Users/zh/go/bin/goctl 生成代码** ⭐ - 禁止手写 handler/types/server 等框架代码
6. **接口变更必须重新生成** - 修改 .proto 或 .api 后必须执行 /Users/zh/go/bin/goctl 命令

### 常见问题

1. **服务发现问题**
   - 开发环境使用 `127.0.0.1` 而非 `0.0.0.0`
   - 检查 etcd 注册地址

2. **RPC 调用失败**
   - 检查服务是否启动
   - 检查配置文件中的 RPC 地址

3. **数据库连接问题**
   - 确认 MySQL 服务运行
   - 检查连接字符串格式

4. **/Users/zh/go/bin/goctl 生成代码问题**
   - 确保 /Users/zh/go/bin/goctl 版本与 go-zero 版本匹配：`/Users/zh/go/bin/goctl --version`
   - proto 文件语法错误：先用 `protoc` 验证
   - api 文件语法错误：检查 type 定义和路由格式
   - 生成失败时检查 `go_package` 配置是否正确

5. **代码被覆盖问题**
   - /Users/zh/go/bin/goctl 重新生成会覆盖 handler/types/server
   - logic 和 svc 不会被覆盖，业务代码安全
   - 如需保留自定义 handler，考虑使用中间件

---

## 🚀 快速开始模板

当需要重构新模块时，请按以下步骤进行：

```markdown
## 第一步：确认模块信息

请告诉我：
1. 要重构的模块名称
2. 主项目中对应的路径
3. 核心功能列表

## 第二步：我将输出技术设计文档

包含：
- 概述与架构
- 接口定义
- 实施计划

## 第三步：确认后开始实施

按照设计文档逐步实现：
1. 基础架构
2. 核心功能
3. 测试验证
```

---

**准备就绪，开始重构！** 🎉

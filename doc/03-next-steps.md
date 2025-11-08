# MyBilibili 下一步操作指南

## 📊 当前进度：35% 完成

### ✅ 已完成
- ✅ 项目初始化（100%）
- ✅ 公共组件（100%）
- ✅ 数据库 Model（100%）- 手动创建
- ✅ Proto 文件定义（100%）

### 🔨 立即可执行的步骤

#### 步骤1：安装必要工具（5分钟）

```bash
# 安装 goctl
GO111MODULE=on GOPROXY=https://goproxy.cn/,direct go install github.com/zeromicro/go-zero/tools/goctl@latest

# 安装 protoc
brew install protobuf

# 安装 proto-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 验证安装
goctl --version
protoc --version
```

#### 步骤2：生成 video-rpc 代码（1分钟）

```bash
cd /Users/zh/project/goproj/bilibili/mybilibili/app/video/cmd/rpc

# 生成 RPC 代码
goctl rpc protoc video.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style go_zero

# 查看生成的文件
tree -L 3
```

**生成的文件结构**：
```
app/video/cmd/rpc/
├── internal/
│   ├── config/           # ✅ 已手动创建
│   ├── logic/            # 将生成：业务逻辑
│   │   ├── getVideoInfoLogic.go
│   │   ├── batchGetVideoInfoLogic.go
│   │   ├── getVideoListLogic.go
│   │   ├── getVideoStatLogic.go
│   │   └── batchGetVideoStatLogic.go
│   ├── server/           # 将生成：gRPC server
│   │   └── videoServer.go
│   └── svc/              # 将生成：Service Context
│       └── serviceContext.go
├── pb/                   # 将生成：protobuf 代码
│   ├── video.pb.go
│   └── video_grpc.pb.go
├── videoclient/          # 将生成：RPC 客户端
│   └── video.go
└── video.go              # 将生成：启动入口
```

#### 步骤3：实现 video-rpc 业务逻辑（30分钟）

生成代码后，需要实现以下5个 logic 文件：

**3.1 实现 GetVideoInfoLogic**

编辑 `internal/logic/getVideoInfoLogic.go`:

```go
func (l *GetVideoInfoLogic) GetVideoInfo(in *video.GetVideoInfoReq) (*video.GetVideoInfoResp, error) {
    // 从数据库查询视频信息
    info, err := l.svcCtx.VideoInfoModel.FindOne(l.ctx, in.Vid)
    if err != nil {
        if err == model.ErrNotFound {
            return nil, xerr.NewCodeErrorWithMsg(xerr.VIDEO_NOT_FOUND)
        }
        return nil, err
    }
    
    // 转换为 proto 结构
    return &video.GetVideoInfoResp{
        Info: &video.VideoInfo{
            Vid:        info.Vid,
            Title:      info.Title,
            Cover:      info.Cover,
            AuthorId:   info.AuthorId,
            AuthorName: info.AuthorName,
            RegionId:   int64(info.RegionId),
            PubTime:    info.PubTime,
            Duration:   int32(info.Duration),
            Desc:       info.Desc,
            State:      int32(info.State),
        },
    }, nil
}
```

**3.2 实现 BatchGetVideoInfoLogic**（⭐重要，hotrank-job 会调用）

```go
func (l *BatchGetVideoInfoLogic) BatchGetVideoInfo(in *video.BatchGetVideoInfoReq) (*video.BatchGetVideoInfoResp, error) {
    // 批量查询视频信息
    infos, err := l.svcCtx.VideoInfoModel.FindByVids(l.ctx, in.Vids)
    if err != nil {
        return nil, err
    }
    
    // 转换为 map
    result := make(map[int64]*video.VideoInfo)
    for _, info := range infos {
        result[info.Vid] = &video.VideoInfo{
            Vid:        info.Vid,
            Title:      info.Title,
            Cover:      info.Cover,
            AuthorId:   info.AuthorId,
            AuthorName: info.AuthorName,
            RegionId:   int64(info.RegionId),
            PubTime:    info.PubTime,
            Duration:   int32(info.Duration),
            Desc:       info.Desc,
            State:      int32(info.State),
        }
    }
    
    return &video.BatchGetVideoInfoResp{Infos: result}, nil
}
```

**3.3 实现 BatchGetVideoStatLogic**（⭐重要，hotrank-job 会调用）

```go
func (l *BatchGetVideoStatLogic) BatchGetVideoStat(in *video.BatchGetVideoStatReq) (*video.BatchGetVideoStatResp, error) {
    // 批量查询统计数据
    stats, err := l.svcCtx.VideoStatModel.FindByVids(l.ctx, in.Vids)
    if err != nil {
        return nil, err
    }
    
    // 转换为 map
    result := make(map[int64]*video.VideoStat)
    for _, stat := range stats {
        result[stat.Vid] = &video.VideoStat{
            Vid:     stat.Vid,
            View:    stat.View,
            Like:    stat.LikeCount,
            Coin:    stat.Coin,
            Fav:     stat.Fav,
            Share:   stat.Share,
            Reply:   stat.Reply,
            Danmaku: stat.Danmaku,
        }
    }
    
    return &video.BatchGetVideoStatResp{Stats: result}, nil
}
```

**3.4 修改 serviceContext.go**

```go
package svc

import (
    "mybilibili/app/video/cmd/rpc/internal/config"
    "mybilibili/common/model"
    
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
    Config         config.Config
    VideoInfoModel model.VideoInfoModel
    VideoStatModel model.VideoStatModel
}

func NewServiceContext(c config.Config) *ServiceContext {
    // 初始化数据库连接
    conn := sqlx.NewMysql(c.Mysql.DataSource)
    
    return &ServiceContext{
        Config:         c,
        VideoInfoModel: model.NewVideoInfoModel(conn, c.CacheRedis),
        VideoStatModel: model.NewVideoStatModel(conn, c.CacheRedis),
    }
}
```

#### 步骤4：启动并测试 video-rpc（5分钟）

```bash
# 1. 启动基础服务
cd deploy
docker-compose up -d

# 2. 初始化数据库
mysql -h127.0.0.1 -P33060 -uroot -proot123456 < sql/001_init.sql
mysql -h127.0.0.1 -P33060 -uroot -proot123456 mybilibili < sql/002_test_data.sql

# 3. 启动 video-rpc
cd ../app/video/cmd/rpc
go run video.go -f etc/video.yaml

# 查看日志，应该看到：
# Starting rpc server at 0.0.0.0:9001...
```

#### 步骤5：开发 hotrank-job（核心，1-2小时）

这是**最核心**的部分，完全参考主项目实现：

**5.1 创建目录结构**

```bash
cd /Users/zh/project/goproj/bilibili/mybilibili/app/hotrank/cmd/job
mkdir -p {internal/{config,dao,model,service,svc},etc}
```

**5.2 创建配置文件** `etc/hotrank-job.yaml`

```yaml
Name: hotrank-job

# 热度计算开关
HotSwitch: true

# MySQL 配置
Mysql:
  DataSource: root:root123456@tcp(127.0.0.1:33060)/mybilibili?charset=utf8mb4&parseTime=true

# Video RPC 配置
VideoRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:23790
    Key: video.rpc

Log:
  ServiceName: hotrank-job
  Mode: console
  Level: info

Telemetry:
  Name: hotrank-job
  Endpoint: http://localhost:14268/api/traces
  Sampler: 1.0
  Batcher: jaeger
```

**5.3 实现核心逻辑**（我已经在设计方案中提供了完整代码）

主要文件：
- `internal/dao/academy.go` - 游标分页、CASE WHEN批量更新
- `internal/dao/video.go` - RPC调用封装
- `internal/service/academy.go` - 热度计算核心逻辑
- `hotrankjob.go` - 启动入口

**关键代码片段**（完全参考主项目）：

```go
// countArcHot 热度计算公式
func countArcHot(stat *pb.VideoStat, ptime int64) int64 {
    hot := float64(stat.Coin)*0.4 +
           float64(stat.Fav)*0.3 +
           float64(stat.Danmaku)*0.4 +
           float64(stat.Reply)*0.4 +
           float64(stat.View)*0.25 +
           float64(stat.Like)*0.4 +
           float64(stat.Share)*0.6
    
    // 24小时内发布的新视频提权
    if ptime >= time.Now().AddDate(0, 0, -1).Unix() && 
       ptime <= time.Now().Unix() {
        hot *= 1.5
    }
    
    return int64(math.Floor(hot))
}
```

#### 步骤6：验证热度计算（5分钟）

```bash
# 1. 启动 hotrank-job
cd app/hotrank/cmd/job
go run hotrankjob.go -f etc/hotrank-job.yaml

# 2. 查看日志，应该看到：
# FlushHot success: processed 15 videos, last_id=15

# 3. 验证数据库
mysql -h127.0.0.1 -P33060 -uroot -proot123456 -e "
USE mybilibili;
SELECT oid, hot, FROM_UNIXTIME(pub_time) as pub_time 
FROM academy_archive 
ORDER BY hot DESC 
LIMIT 10;
"

# 应该看到热度值已更新，新视频（24小时内）热度较高
```

### 🎯 核心文件清单

您需要重点关注以下文件：

| 文件 | 状态 | 说明 |
|-----|------|-----|
| `app/video/cmd/rpc/video.proto` | ✅ 已创建 | RPC 接口定义 |
| `app/video/cmd/rpc/etc/video.yaml` | ✅ 已创建 | RPC 配置 |
| `common/model/videoInfoModel.go` | ✅ 已创建 | 视频信息Model |
| `common/model/videoStatModel.go` | ✅ 已创建 | 视频统计Model |
| `app/hotrank/cmd/job/internal/service/academy.go` | ❌ 待创建 | **核心**：热度计算逻辑 |
| `app/hotrank/cmd/job/internal/dao/academy.go` | ❌ 待创建 | 游标分页、批量更新 |
| `app/hotrank/cmd/job/internal/dao/video.go` | ❌ 待创建 | RPC调用封装 |

### 💡 开发提示

1. **先完成 video-rpc**：hotrank-job 依赖它
2. **重点实现批量接口**：`BatchGetVideoInfo` 和 `BatchGetVideoStat`
3. **完全参考主项目**：设计方案中的代码可以直接复制使用
4. **测试驱动**：每完成一个服务，立即测试

### 📞 需要帮助？

如果遇到问题，可以：
1. 查看 `doc/02-progress-summary.md` 了解当前进度
2. 参考 `设计方案/v0.0.2-基于主项目优化.md` 中的完整代码
3. 查看主项目 Bilibili 的对应实现

---

**预计完成时间**：2-3小时
**核心开发时间**：hotrank-job 约1-2小时


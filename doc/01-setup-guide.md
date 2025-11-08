# MyBilibili 开发环境搭建指南

## 当前进度

✅ 阶段1：项目初始化完成
✅ 阶段2：公共组件开发完成
🔄 阶段3：视频服务开发（进行中）

## 需要安装的工具

### 1. 安装 goctl 工具

```bash
# 方式1：使用 go install（需要 Go 1.16+）
GO111MODULE=on GOPROXY=https://goproxy.cn/,direct go install github.com/zeromicro/go-zero/tools/goctl@latest

# 方式2：下载二进制文件
# Mac Intel
wget https://github.com/zeromicro/go-zero/releases/download/tools/goctl/v1.6.1/goctl-v1.6.1-darwin-amd64.tar.gz
tar -xzf goctl-v1.6.1-darwin-amd64.tar.gz
sudo mv goctl /usr/local/bin/

# Mac Apple Silicon
wget https://github.com/zeromicro/go-zero/releases/download/tools/goctl/v1.6.1/goctl-v1.6.1-darwin-arm64.tar.gz
tar -xzf goctl-v1.6.1-darwin-arm64.tar.gz
sudo mv goctl /usr/local/bin/

# 验证安装
goctl --version
```

### 2. 安装 protoc（Protocol Buffers 编译器）

```bash
# Mac
brew install protobuf

# 验证安装
protoc --version
```

### 3. 安装 protoc-gen-go 和 protoc-gen-go-grpc

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 接下来的步骤

安装完 goctl 后，请执行以下命令继续开发：

### 1. 生成 video-rpc 代码

```bash
cd /Users/zh/project/goproj/bilibili/mybilibili/app/video/cmd/rpc
goctl rpc protoc video.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style go_zero
```

### 2. 生成数据库 Model 代码

```bash
cd /Users/zh/project/goproj/bilibili/mybilibili
goctl model mysql datasource \
  -url="root:root123456@tcp(127.0.0.1:33060)/mybilibili" \
  -table="video_info,video_stat,academy_archive" \
  -dir=./common/model \
  --style go_zero \
  -c
```

### 3. 启动基础服务

```bash
cd deploy
docker-compose up -d

# 等待服务启动
sleep 10

# 查看服务状态
docker-compose ps
```

### 4. 初始化数据库

```bash
# 连接 MySQL（密码：root123456）
mysql -h127.0.0.1 -P33060 -uroot -p

# 或者使用命令行直接执行
mysql -h127.0.0.1 -P33060 -uroot -proot123456 < deploy/sql/001_init.sql
mysql -h127.0.0.1 -P33060 -uroot -proot123456 mybilibili < deploy/sql/002_test_data.sql
```

## 项目已完成的部分

### 1. 项目结构 ✅
```
mybilibili/
├── app/
│   ├── video/cmd/rpc/        # 视频服务（proto文件已创建）
│   ├── hotrank/cmd/{rpc,job} # 热门排行榜
│   └── creative/cmd/api/      # API网关
├── common/
│   ├── xerr/                  # 错误码定义 ✅
│   ├── tool/                  # 工具函数 ✅
│   └── result/                # 响应封装 ✅
└── deploy/
    ├── docker-compose.yml     # Docker编排 ✅
    └── sql/                   # SQL脚本 ✅
```

### 2. 基础组件 ✅
- ✅ 错误码定义（common/xerr）
- ✅ 工具函数（common/tool/xstr.go）
- ✅ 响应封装（common/result）

### 3. Proto 文件 ✅
- ✅ video.proto（视频服务接口定义）

### 4. 配置文件 ✅
- ✅ video.yaml（视频服务配置）
- ✅ docker-compose.yml（基础服务）
- ✅ prometheus.yml（监控配置）

### 5. 数据库脚本 ✅
- ✅ 001_init.sql（表结构）
- ✅ 002_test_data.sql（测试数据）

## 下一步开发计划

### 方案A：使用 goctl 生成代码（推荐）

如果已安装 goctl，执行：

```bash
# 1. 生成 video-rpc 代码
cd app/video/cmd/rpc
goctl rpc protoc video.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style go_zero

# 2. 生成数据库 Model
cd ../../../../
goctl model mysql datasource \
  -url="root:root123456@tcp(127.0.0.1:33060)/mybilibili" \
  -table="video_info,video_stat" \
  -dir=./common/model \
  --style go_zero \
  -c

# 3. 实现业务逻辑
# - 编辑 app/video/cmd/rpc/internal/logic/*.go
# - 实现 GetVideoInfo、BatchGetVideoInfo 等方法
```

### 方案B：手动创建代码

如果无法安装 goctl，我可以手动创建所有必要的代码文件。

---

**请告诉我：**
1. 您是否可以安装 goctl 工具？
2. 如果可以，请运行上述安装命令后告诉我
3. 如果不能，我将继续手动创建所有代码


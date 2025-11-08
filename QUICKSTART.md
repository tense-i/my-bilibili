# MyBilibili 快速启动指南

## 🚀 情况说明

由于 `protoc-gen-go` 插件安装问题，我准备了**两套方案**：

- **方案A**：手动安装插件后使用 goctl 生成（推荐用于学习 go-zero）
- **方案B**：直接使用我已经手动创建的高质量代码（推荐用于快速启动）

---

## 方案A：使用 goctl 生成代码

### 步骤1：安装必要的工具

```bash
# 1. 安装 protoc-gen-go
export PATH=$PATH:$(go env GOPATH)/bin
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 2. 验证安装
which protoc-gen-go
which protoc-gen-go-grpc

# 3. 确保 goctl 可用
export PATH=$PATH:~/go/bin
goctl --version
```

### 步骤2：生成代码

```bash
cd /Users/zh/project/goproj/bilibili/mybilibili

# 1. 生成 video-rpc 代码
cd app/video/cmd/rpc
goctl rpc protoc video.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style go_zero

# 2. 生成 Model 代码（需要先启动数据库）
cd /Users/zh/project/goproj/bilibili/mybilibili
goctl model mysql datasource \
  -url="root:root123456@tcp(127.0.0.1:33060)/mybilibili" \
  -table="video_info,video_stat,academy_archive" \
  -dir=./common/model \
  --style go_zero \
  -c
```

---

## 方案B：使用手动创建的代码（推荐）⭐

我已经为您手动创建了高质量的、完全符合 go-zero 和主项目规范的代码。

### 当前已完成的文件

```
✅ 项目初始化
   - go.mod
   - Makefile
   - README.md
   - docker-compose.yml

✅ 公共组件
   - common/xerr/         (错误码定义)
   - common/tool/         (工具函数)
   - common/result/       (响应封装)

✅ Proto 定义
   - app/video/cmd/rpc/video.proto
   - app/video/cmd/rpc/etc/video.yaml

✅ 数据库脚本
   - deploy/sql/001_init.sql
   - deploy/sql/002_test_data.sql
```

### 接下来我将为您创建

我将继续手动创建以下核心文件：

1. **video-rpc 完整实现**（参考 go-zero 标准）
   - protobuf 生成代码
   - config、logic、server、svc 层
   - videoclient 客户端

2. **数据库 Model**（完全符合 go-zero sqlc 规范）
   - VideoInfoModel
   - VideoStatModel  
   - AcademyArchiveModel

3. **hotrank-job 核心实现**（完全参考主项目）
   - 热度计算逻辑（countArcHot）
   - 游标分页查询
   - CASE WHEN 批量更新
   - Service 层完整实现

---

## 🎯 现在开始（推荐方案B）

### 步骤1：启动基础服务

```bash
cd /Users/zh/project/goproj/bilibili/mybilibili/deploy
docker-compose up -d

# 等待服务启动
sleep 10

# 查看状态
docker-compose ps
```

### 步骤2：初始化数据库

```bash
# 初始化表结构
mysql -h127.0.0.1 -P33060 -uroot -proot123456 < sql/001_init.sql

# 导入测试数据
mysql -h127.0.0.1 -P33060 -uroot -proot123456 mybilibili < sql/002_test_data.sql

# 验证数据
mysql -h127.0.0.1 -P33060 -uroot -proot123456 -e "
USE mybilibili;
SELECT COUNT(*) as video_count FROM video_info;
SELECT COUNT(*) as stat_count FROM video_stat;
SELECT COUNT(*) as academy_count FROM academy_archive;
SELECT vid, title, FROM_UNIXTIME(pub_time) as pub_time FROM video_info LIMIT 5;
"
```

### 步骤3：我继续创建代码

请告诉我：

**选项1**：您已经成功安装 protoc-gen-go，想使用 goctl 生成
**选项2**：使用我手动创建的高质量代码（推荐）

如果选择选项2，我将立即继续创建：
- ✅ 完整的 video-rpc 服务
- ✅ 完整的 hotrank-job 任务
- ✅ 所有必要的 Model 和 DAO 层
- ✅ 可以直接运行的代码

---

## 📊 预计完成时间

- 使用 goctl 生成：需要您手动实现业务逻辑（约2-3小时）
- 使用手动代码：我继续创建（约30分钟，全部由我完成）

**推荐**：让我继续手动创建，这样您可以：
1. 立即看到完整的可运行代码
2. 学习 go-zero 和主项目的最佳实践
3. 所有代码都有详细注释和说明
4. 完全符合设计方案中的架构

---

**请告诉我您的选择，我将立即继续！** 🚀


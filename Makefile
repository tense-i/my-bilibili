.PHONY: help init docker-up docker-down db-init gen-video-rpc gen-hotrank-rpc gen-api gen-model run-video run-hotrank-job run-hotrank-rpc run-api test clean

# 默认目标
help:
	@echo "MyBilibili Makefile"
	@echo ""
	@echo "使用方法:"
	@echo "  make init          - 初始化项目（安装依赖）"
	@echo "  make docker-up     - 启动基础服务（MySQL、Redis、etcd、Jaeger）"
	@echo "  make docker-down   - 停止基础服务"
	@echo "  make db-init       - 初始化数据库"
	@echo "  make gen-all       - 生成所有代码"
	@echo "  make run-all       - 启动所有服务"
	@echo "  make test          - 运行测试"
	@echo "  make clean         - 清理生成的文件"

# 初始化项目
init:
	@echo "==> 安装 goctl"
	go install github.com/zeromicro/go-zero/tools/goctl@latest
	@echo "==> 下载依赖"
	go mod download
	@echo "==> 初始化完成"

# 启动基础服务
docker-up:
	@echo "==> 启动基础服务"
	cd deploy && docker-compose up -d
	@echo "==> 等待服务启动..."
	sleep 10
	@echo "==> 服务已启动"
	cd deploy && docker-compose ps

# 停止基础服务
docker-down:
	@echo "==> 停止基础服务"
	cd deploy && docker-compose down

# 初始化数据库
db-init:
	@echo "==> 初始化数据库"
	mysql -h127.0.0.1 -P33060 -uroot -proot123456 < deploy/sql/001_init.sql
	@echo "==> 导入测试数据"
	mysql -h127.0.0.1 -P33060 -uroot -proot123456 mybilibili < deploy/sql/002_test_data.sql
	@echo "==> 数据库初始化完成"

# 生成所有代码
gen-all: gen-video-rpc gen-model

# 生成 video-rpc 代码
gen-video-rpc:
	@echo "==> 生成 video-rpc 代码"
	cd app/video/cmd/rpc && goctl rpc protoc video.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style go_zero
	@echo "==> video-rpc 代码生成完成"

# 生成 hotrank-rpc 代码
gen-hotrank-rpc:
	@echo "==> 生成 hotrank-rpc 代码"
	cd app/hotrank/cmd/rpc && goctl rpc protoc hotrank.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style go_zero
	@echo "==> hotrank-rpc 代码生成完成"

# 生成 API 代码
gen-api:
	@echo "==> 生成 creative-api 代码"
	cd app/creative/cmd/api && goctl api go -api creative.api -dir . --style go_zero
	@echo "==> creative-api 代码生成完成"

# 生成 Model 代码
gen-model:
	@echo "==> 生成 Model 代码"
	goctl model mysql datasource \
		-url="root:root123456@tcp(127.0.0.1:33060)/mybilibili" \
		-table="video_info,video_stat,academy_archive" \
		-dir=./common/model \
		--style go_zero
	@echo "==> Model 代码生成完成"

# 启动所有服务
run-all:
	@echo "==> 启动所有服务（后台运行）"
	@make run-video &
	@sleep 5
	@make run-hotrank-job &
	@sleep 2
	@make run-hotrank-rpc &
	@sleep 2
	@make run-api &
	@echo "==> 所有服务已启动"

# 启动 video-rpc
run-video:
	@echo "==> 启动 video-rpc"
	cd app/video/cmd/rpc && go run video.go -f etc/video.yaml

# 启动 hotrank-job
run-hotrank-job:
	@echo "==> 启动 hotrank-job"
	cd app/hotrank/cmd/job && go run hotrankjob.go -f etc/hotrank-job.yaml

# 启动 hotrank-rpc
run-hotrank-rpc:
	@echo "==> 启动 hotrank-rpc"
	cd app/hotrank/cmd/rpc && go run hotrank.go -f etc/hotrank.yaml

# 启动 creative-api
run-api:
	@echo "==> 启动 creative-api"
	cd app/creative/cmd/api && go run creative.go -f etc/creative-api.yaml

# 测试
test:
	@echo "==> 运行测试"
	go test -v ./...

# 清理
clean:
	@echo "==> 清理生成的文件"
	find . -name "*.pb.go" -delete
	find . -name "*_gen.go" -delete
	@echo "==> 清理完成"

# 格式化代码
fmt:
	@echo "==> 格式化代码"
	gofmt -s -w .
	@echo "==> 格式化完成"

# 代码检查
lint:
	@echo "==> 代码检查"
	golangci-lint run
	@echo "==> 检查完成"


.PHONY: proto api rpc build clean

# 生成 protobuf 代码
proto:
	cd cmd/rpc && protoc --go_out=. --go-grpc_out=. search.proto

# 生成 API 代码 (可选，使用 goctl)
api:
	cd cmd/api && goctl api go -api desc/search.api -dir . -style go_zero

# 构建 RPC 服务
rpc:
	cd cmd/rpc && go build -o ../../bin/search-rpc search.go

# 构建 API 服务
api-build:
	cd cmd/api && go build -o ../../bin/search-api search.go

# 构建所有
build: proto rpc api-build

# 清理
clean:
	rm -rf bin/
	rm -f cmd/rpc/search/*.pb.go

# 运行 RPC 服务
run-rpc:
	cd cmd/rpc && go run search.go -f etc/search.yaml

# 运行 API 服务
run-api:
	cd cmd/api && go run search.go -f etc/search.yaml

# 测试
test:
	go test -v ./...

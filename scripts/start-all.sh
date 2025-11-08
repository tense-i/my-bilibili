#!/bin/bash

# MyBilibili 一键启动脚本
# 用法: ./scripts/start-all.sh

set -e

PROJECT_ROOT=$(cd "$(dirname "$0")/.." && pwd)
cd "$PROJECT_ROOT"

echo "========================================="
echo "  MyBilibili 热门视频排行榜系统启动"
echo "========================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查命令是否存在
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# 检查端口是否被占用
check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null ; then
        return 0
    else
        return 1
    fi
}

# 等待端口就绪
wait_for_port() {
    local port=$1
    local service=$2
    local max_wait=30
    local count=0
    
    echo -n "  等待 $service 启动 (端口 $port)..."
    while ! check_port $port; do
        sleep 1
        count=$((count + 1))
        if [ $count -ge $max_wait ]; then
            echo -e "${RED} 超时！${NC}"
            return 1
        fi
    done
    echo -e "${GREEN} ✓${NC}"
    return 0
}

# 1. 检查依赖
echo ""
echo "1. 检查依赖..."

if ! command_exists go; then
    echo -e "${RED}错误: 未安装 Go${NC}"
    exit 1
fi

if ! command_exists docker; then
    echo -e "${RED}错误: 未安装 Docker${NC}"
    exit 1
fi

if ! command_exists docker-compose; then
    echo -e "${RED}错误: 未安装 docker-compose${NC}"
    exit 1
fi

echo -e "${GREEN}  所有依赖已就绪${NC}"

# 2. 启动基础设施
echo ""
echo "2. 启动基础设施 (MySQL, Redis, etcd, Jaeger)..."
cd "$PROJECT_ROOT/deploy"

docker-compose up -d mysql redis etcd jaeger

# 等待服务就绪
wait_for_port 3306 "MySQL" || exit 1
wait_for_port 6379 "Redis" || exit 1
wait_for_port 2379 "etcd" || exit 1
wait_for_port 16686 "Jaeger UI" || exit 1

echo -e "${GREEN}  基础设施启动完成${NC}"

# 3. 初始化数据库（如果需要）
echo ""
echo "3. 检查数据库..."

# 检查数据库是否已初始化
DB_EXISTS=$(docker exec mybilibili-mysql mysql -uroot -p123456 -e "SHOW DATABASES LIKE 'mybilibili';" 2>/dev/null | grep mybilibili || echo "")

if [ -z "$DB_EXISTS" ]; then
    echo "  初始化数据库..."
    docker exec -i mybilibili-mysql mysql -uroot -p123456 < "$PROJECT_ROOT/deploy/sql/001_init.sql"
    docker exec -i mybilibili-mysql mysql -uroot -p123456 mybilibili < "$PROJECT_ROOT/deploy/sql/002_test_data.sql"
    
    # 插入 academy_archive 测试数据
    docker exec mybilibili-mysql mysql -uroot -p123456 mybilibili -e "
        INSERT IGNORE INTO academy_archive (oid, uid, business, region_id, hot, state, ctime, mtime) VALUES
        (1001, 1, 1, 1, 0, 0, NOW(), NOW()),
        (1002, 2, 1, 2, 0, 0, NOW(), NOW()),
        (1003, 3, 1, 1, 0, 0, NOW(), NOW()),
        (1004, 4, 1, 3, 0, 0, NOW(), NOW()),
        (1005, 5, 1, 2, 0, 0, NOW(), NOW());
    "
    
    echo -e "${GREEN}  数据库初始化完成${NC}"
else
    echo -e "${GREEN}  数据库已存在${NC}"
fi

# 4. 编译所有服务
echo ""
echo "4. 编译服务..."

cd "$PROJECT_ROOT"

# 编译 video-rpc
echo "  编译 video-rpc..."
cd "$PROJECT_ROOT/app/video/cmd/rpc"
go build -o video-rpc video.go

# 编译 hotrank-rpc
echo "  编译 hotrank-rpc..."
cd "$PROJECT_ROOT/app/hotrank/cmd/rpc"
go build -o hotrank-rpc hotrank.go

# 编译 hotrank-job
echo "  编译 hotrank-job..."
cd "$PROJECT_ROOT/app/hotrank/cmd/job"
go build -o hotrank-job job.go

# 编译 creative-api
echo "  编译 creative-api..."
cd "$PROJECT_ROOT/app/api/creative"
go build -o creative-api creative.go

echo -e "${GREEN}  编译完成${NC}"

# 5. 启动服务
echo ""
echo "5. 启动服务..."

# 创建日志目录
mkdir -p "$PROJECT_ROOT/logs"

# 启动 video-rpc
echo "  启动 video-rpc (端口 9002)..."
cd "$PROJECT_ROOT/app/video/cmd/rpc"
nohup ./video-rpc -f etc/video.yaml > "$PROJECT_ROOT/logs/video-rpc.log" 2>&1 &
VIDEO_RPC_PID=$!
echo $VIDEO_RPC_PID > "$PROJECT_ROOT/logs/video-rpc.pid"
wait_for_port 9002 "video-rpc" || exit 1

# 启动 hotrank-rpc
echo "  启动 hotrank-rpc (端口 9003)..."
cd "$PROJECT_ROOT/app/hotrank/cmd/rpc"
nohup ./hotrank-rpc -f etc/hotrank.yaml > "$PROJECT_ROOT/logs/hotrank-rpc.log" 2>&1 &
HOTRANK_RPC_PID=$!
echo $HOTRANK_RPC_PID > "$PROJECT_ROOT/logs/hotrank-rpc.pid"
wait_for_port 9003 "hotrank-rpc" || exit 1

# 启动 hotrank-job
echo "  启动 hotrank-job..."
cd "$PROJECT_ROOT/app/hotrank/cmd/job"
nohup ./hotrank-job -f etc/hotrank-job.yaml > "$PROJECT_ROOT/logs/hotrank-job.log" 2>&1 &
HOTRANK_JOB_PID=$!
echo $HOTRANK_JOB_PID > "$PROJECT_ROOT/logs/hotrank-job.pid"
sleep 3  # 等待任务启动

# 启动 creative-api
echo "  启动 creative-api (端口 8001)..."
cd "$PROJECT_ROOT/app/api/creative"
nohup ./creative-api -f etc/creative-api.yaml > "$PROJECT_ROOT/logs/creative-api.log" 2>&1 &
CREATIVE_API_PID=$!
echo $CREATIVE_API_PID > "$PROJECT_ROOT/logs/creative-api.pid"
wait_for_port 8001 "creative-api" || exit 1

echo -e "${GREEN}  所有服务启动完成${NC}"

# 6. 显示服务状态
echo ""
echo "========================================="
echo "  服务启动完成！"
echo "========================================="
echo ""
echo "服务列表："
echo "  • video-rpc      : http://localhost:9002 (PID: $VIDEO_RPC_PID)"
echo "  • hotrank-rpc    : http://localhost:9003 (PID: $HOTRANK_RPC_PID)"
echo "  • hotrank-job    : 后台运行 (PID: $HOTRANK_JOB_PID)"
echo "  • creative-api   : http://localhost:8001 (PID: $CREATIVE_API_PID)"
echo ""
echo "基础设施："
echo "  • MySQL          : localhost:3306 (用户: root, 密码: 123456)"
echo "  • Redis          : localhost:6379"
echo "  • etcd           : localhost:2379"
echo "  • Jaeger UI      : http://localhost:16686"
echo ""
echo "监控指标："
echo "  • creative-api   : http://localhost:9091/metrics"
echo "  • video-rpc      : http://localhost:9092/metrics"
echo "  • hotrank-rpc    : http://localhost:9093/metrics"
echo ""
echo "API 测试："
echo "  • 全站排行榜     : curl http://localhost:8001/api/creative/v1/hotrank/list"
echo "  • 分区排行榜     : curl http://localhost:8001/api/creative/v1/hotrank/region?region_id=1"
echo "  • 视频详情       : curl http://localhost:8001/api/creative/v1/video/1001"
echo ""
echo "日志文件："
echo "  • 所有日志在     : $PROJECT_ROOT/logs/"
echo ""
echo "停止所有服务："
echo "  • 运行脚本       : ./scripts/stop-all.sh"
echo ""
echo -e "${YELLOW}提示: 等待 10 秒后，hotrank-job 会开始计算热度值${NC}"
echo "========================================="


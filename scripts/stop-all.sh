#!/bin/bash

# MyBilibili 停止脚本
# 用法: ./scripts/stop-all.sh

set -e

PROJECT_ROOT=$(cd "$(dirname "$0")/.." && pwd)
cd "$PROJECT_ROOT"

echo "========================================="
echo "  MyBilibili 停止所有服务"
echo "========================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# 停止进程
stop_service() {
    local name=$1
    local pid_file="$PROJECT_ROOT/logs/${name}.pid"
    
    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        if ps -p $pid > /dev/null 2>&1; then
            echo -n "  停止 $name (PID: $pid)..."
            kill $pid
            sleep 1
            
            # 如果还在运行，强制kill
            if ps -p $pid > /dev/null 2>&1; then
                kill -9 $pid
            fi
            
            rm -f "$pid_file"
            echo -e "${GREEN} ✓${NC}"
        else
            echo "  $name 未运行"
            rm -f "$pid_file"
        fi
    else
        echo "  $name 未运行"
    fi
}

# 1. 停止应用服务
echo ""
echo "1. 停止应用服务..."

stop_service "creative-api"
stop_service "hotrank-job"
stop_service "hotrank-rpc"
stop_service "video-rpc"

# 2. 停止基础设施 (可选)
echo ""
read -p "是否停止基础设施 (MySQL, Redis, etcd, Jaeger)? [y/N] " -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "  停止基础设施..."
    cd "$PROJECT_ROOT/deploy"
    docker-compose down
    echo -e "${GREEN}  基础设施已停止${NC}"
else
    echo "  保留基础设施运行"
fi

echo ""
echo "========================================="
echo -e "${GREEN}  所有服务已停止${NC}"
echo "========================================="



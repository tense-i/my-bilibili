package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"mybilibili/app/hotrank/cmd/job/internal/config"
	"mybilibili/app/hotrank/cmd/job/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/hotrank-job.yaml", "the config file")

func main() {
	flag.Parse()

	// 加载配置
	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 初始化日志
	logx.MustSetup(logx.LogConf{
		ServiceName: c.Log.ServiceName,
		Mode:        c.Log.Mode,
		Level:       c.Log.Level,
	})
	defer logx.Close()

	logx.Info("Starting hotrank-job...")

	// 初始化服务上下文（会自动启动热度计算任务）
	ctx := svc.NewServiceContext(c)

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	logx.Info("hotrank-job started successfully")
	fmt.Println("hotrank-job started successfully, press Ctrl+C to exit")

	<-quit

	logx.Info("Shutting down hotrank-job...")

	// 优雅关闭
	ctx.Service.Close()

	logx.Info("hotrank-job stopped")
}

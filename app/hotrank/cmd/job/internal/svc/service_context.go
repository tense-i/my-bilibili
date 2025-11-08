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
	// 初始化 Video RPC 客户端（添加重试配置）
	videoRpcConf := c.VideoRpc
	videoRpcConf.Timeout = 10000 // 10秒超时
	videoRpc := video_client.NewVideo(zrpc.MustNewClient(videoRpcConf))

	// 初始化服务（会自动启动热度计算任务）
	svc := service.New(c, videoRpc)

	return &ServiceContext{
		Config:  c,
		Service: svc,
	}
}

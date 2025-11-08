// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"mybilibili/app/api/creative/internal/config"
	"mybilibili/app/hotrank/cmd/rpc/hotrank_client"
	"mybilibili/app/recommend/cmd/rpc/recommend_client"
	"mybilibili/app/video/cmd/rpc/video_client"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	VideoRpc     video_client.Video
	HotrankRpc   hotrank_client.Hotrank
	RecommendRpc recommend_client.Recommend
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		VideoRpc:     video_client.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
		HotrankRpc:   hotrank_client.NewHotrank(zrpc.MustNewClient(c.HotrankRpc)),
		RecommendRpc: recommend_client.NewRecommend(zrpc.MustNewClient(c.RecommendRpc)),
	}
}

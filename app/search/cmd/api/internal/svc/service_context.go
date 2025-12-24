package svc

import (
	"mybilibili/app/search/cmd/api/internal/config"
	"mybilibili/app/search/cmd/rpc/search_client"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	SearchRpc search_client.Search
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		SearchRpc: search_client.NewSearch(zrpc.MustNewClient(c.SearchRpc)),
	}
}

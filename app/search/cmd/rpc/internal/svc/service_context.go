package svc

import (
	"mybilibili/app/search/cmd/rpc/internal/config"
	"mybilibili/app/search/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config config.Config

	// Elasticsearch 客户端
	ESClient *model.ESClient

	// Redis 客户端
	Redis *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 ES 集群
	clusters := map[string][]string{
		"dmExternal":     c.Elasticsearch.DmExternal.Addresses,
		"replyExternal":  c.Elasticsearch.ReplyExternal.Addresses,
		"externalPublic": c.Elasticsearch.ExternalPublic.Addresses,
	}

	esClient, err := model.NewESClient(clusters)
	if err != nil {
		logx.Errorf("failed to create ES client: %v", err)
	}

	// 初始化 Redis
	rds := redis.New(c.SearchRedis.Host, func(r *redis.Redis) {
		r.Type = c.SearchRedis.Type
		r.Pass = c.SearchRedis.Pass
	})

	return &ServiceContext{
		Config:   c,
		ESClient: esClient,
		Redis:    rds,
	}
}

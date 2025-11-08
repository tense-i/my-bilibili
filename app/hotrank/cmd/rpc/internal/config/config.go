package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	// MySQL 配置
	Mysql struct {
		DataSource string
	}

	// Redis 缓存配置
	CacheRedis cache.CacheConf
}

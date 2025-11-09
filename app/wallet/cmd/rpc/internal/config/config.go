package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	// MySQL配置
	Mysql struct {
		DataSource string
	}

	// Redis配置（分布式锁）
	RedisConf struct {
		Host string
		Type string
		Pass string
	}

	// 缓存配置
	CacheRedis cache.CacheConf
}

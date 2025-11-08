package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	// MySQL 配置
	Mysql struct {
		DataSource string
	}

	// Redis 缓存配置
	CacheRedis cache.CacheConf

	// 热度计算开关
	HotSwitch bool

	// Video RPC 服务配置
	VideoRpc zrpc.RpcClientConf

	// 日志配置
	Log struct {
		ServiceName string
		Mode        string
		Level       string
	}
}


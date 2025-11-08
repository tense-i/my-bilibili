package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	// MySQL 配置
	MySQL struct {
		DataSource string
	}

	// Redis 配置
	CacheRedis cache.CacheConf

	// 依赖服务 - 召回服务
	RecallRpc zrpc.RpcClientConf

	// 依赖服务 - 视频服务
	VideoRpc zrpc.RpcClientConf

	// 模型配置
	RankModel struct {
		ModelDir     string
		ModelVersion string
		Enable       bool
	}

	// 业务配置
	Business BusinessConfig
}

type BusinessConfig struct {
	// 召回配置
	RecallLimit   int32
	RecallTimeout int64

	// 排序配置
	RankEnable   bool
	RankFallback bool

	// 过滤配置
	BloomFilterEnable bool
	DurationMin       int32
	DurationMax       int32

	// 后处理配置
	ScatterUPMinDistance  int
	ScatterTagMinDistance int

	// 缓存配置
	CacheTTL int
}

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

	// 召回策略配置
	RecallStrategies RecallStrategiesConfig

	// 布隆过滤器配置
	BloomFilter BloomFilterConfig
}

// 召回策略配置
type RecallStrategiesConfig struct {
	HotRecall       RecallStrategyConfig
	SelectionRecall RecallStrategyConfig
	LikeI2IRecall   RecallStrategyConfig
	TagRecall       RecallStrategyConfig
	FollowRecall    RecallStrategyConfig
}

// 单个召回策略配置
type RecallStrategyConfig struct {
	Enable           bool
	Limit            int32
	TopN             int
	LimitPerItem     int32
	LimitPerTag      int32
	LimitPerUP       int32
	Priority         int32
	RedisKey         string
	RedisKeyTemplate string
}

// 布隆过滤器配置
type BloomFilterConfig struct {
	ExpectedInsertions int
	FalsePositiveRate  float64
	RedisKeyTemplate   string
}

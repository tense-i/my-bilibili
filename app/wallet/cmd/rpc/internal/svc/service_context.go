package svc

import (
	"mybilibili/app/wallet/cmd/rpc/internal/config"
	"mybilibili/app/wallet/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	// MySQL连接
	DB sqlx.SqlConn

	// Redis连接（分布式锁）
	Redis *redis.Redis

	// Model层
	UserWalletModel         model.UserWalletModel
	CoinStreamRecordModel   model.CoinStreamRecordModel
	CoinExchangeRecordModel model.CoinExchangeRecordModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化MySQL连接
	db := sqlx.NewMysql(c.Mysql.DataSource)

	// 初始化Redis连接
	rds := redis.New(c.RedisConf.Host, func(r *redis.Redis) {
		r.Type = c.RedisConf.Type
		r.Pass = c.RedisConf.Pass
	})

	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  rds,

		// 初始化Model
		UserWalletModel:         model.NewUserWalletModel(db, c.CacheRedis),
		CoinStreamRecordModel:   model.NewCoinStreamRecordModel(db),
		CoinExchangeRecordModel: model.NewCoinExchangeRecordModel(db),
	}
}

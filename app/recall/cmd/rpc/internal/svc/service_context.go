package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"mybilibili/app/recall/cmd/rpc/internal/config"
	"mybilibili/app/recall/cmd/rpc/internal/dao"
)

type ServiceContext struct {
	Config config.Config
	Dao    *dao.Dao
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 MySQL
	sqlConn := sqlx.NewMysql(c.MySQL.DataSource)

	// 初始化 Redis
	rds := redis.MustNewRedis(c.CacheRedis[0].RedisConf)

	// 初始化 DAO
	d := dao.NewDao(sqlConn, rds)

	return &ServiceContext{
		Config: c,
		Dao:    d,
	}
}

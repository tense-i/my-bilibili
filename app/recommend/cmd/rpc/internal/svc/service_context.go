package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"

	"mybilibili/app/recall/cmd/rpc/recall_client"
	"mybilibili/app/recommend/cmd/rpc/internal/config"
	"mybilibili/app/recommend/cmd/rpc/internal/dao"
	"mybilibili/app/recommend/cmd/rpc/internal/logic/rank"
)

type ServiceContext struct {
	Config config.Config
	Dao    *dao.Dao

	// RPC 客户端
	RecallRpc recall_client.Recall

	// 业务配置
	BusinessConfig config.BusinessConfig

	// 排序模型
	RankModel *rank.RankModelManager
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 MySQL
	sqlConn := sqlx.NewMysql(c.MySQL.DataSource)

	// 初始化 Redis
	rds := redis.MustNewRedis(c.CacheRedis[0].RedisConf)

	// 初始化 DAO
	d := dao.NewDao(sqlConn, rds)

	// 初始化召回服务客户端
	recallRpc := recall_client.NewRecall(zrpc.MustNewClient(c.RecallRpc))

	// 初始化排序模型
	rankModel := rank.NewRankModelManager(c.RankModel.ModelDir)

	return &ServiceContext{
		Config:         c,
		Dao:            d,
		RecallRpc:      recallRpc,
		BusinessConfig: c.Business,
		RankModel:      rankModel,
	}
}

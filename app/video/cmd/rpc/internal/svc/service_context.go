package svc

import (
	"mybilibili/app/video/cmd/rpc/internal/config"
	"mybilibili/common/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config         config.Config
	VideoInfoModel model.VideoInfoModel
	VideoStatModel model.VideoStatModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:         c,
		VideoInfoModel: model.NewVideoInfoModel(conn, c.CacheRedis),
		VideoStatModel: model.NewVideoStatModel(conn, c.CacheRedis),
	}
}

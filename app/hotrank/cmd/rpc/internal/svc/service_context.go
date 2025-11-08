package svc

import (
	"mybilibili/app/hotrank/cmd/rpc/internal/config"
	"mybilibili/common/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config              config.Config
	AcademyArchiveModel model.AcademyArchiveModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:              c,
		AcademyArchiveModel: model.NewAcademyArchiveModel(conn, c.CacheRedis),
	}
}

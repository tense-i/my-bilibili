package svc

import (
	"context"
	"sync"
	"time"

	"mybilibili/app/coupon/cmd/rpc/internal/config"
	"mybilibili/app/coupon/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	DB     sqlx.SqlConn
	Redis  *redis.Redis

	// Models
	CouponBatchInfoModel          model.CouponBatchInfoModel
	CouponAllowanceInfoModel      model.CouponAllowanceInfoModel
	CouponAllowanceChangeLogModel model.CouponAllowanceChangeLogModel
	CouponCodeModel               model.CouponCodeModel
	CouponReceiveLogModel         model.CouponReceiveLogModel
	CouponInfoModel               model.CouponInfoModel

	// 批次信息缓存（内存）
	BatchInfoCache map[string]*model.CouponBatchInfo
	batchMu        sync.RWMutex
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := sqlx.NewMysql(c.Mysql.DataSource)
	rds := redis.MustNewRedis(redis.RedisConf{
		Host: c.RedisConf.Host,
		Type: c.RedisConf.Type,
		Pass: c.RedisConf.Pass,
	})

	svc := &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  rds,

		CouponBatchInfoModel:          model.NewCouponBatchInfoModel(db, c.CacheRedis),
		CouponAllowanceInfoModel:      model.NewCouponAllowanceInfoModel(db, c.CacheRedis),
		CouponAllowanceChangeLogModel: model.NewCouponAllowanceChangeLogModel(db),
		CouponCodeModel:               model.NewCouponCodeModel(db, c.CacheRedis),
		CouponReceiveLogModel:         model.NewCouponReceiveLogModel(db),
		CouponInfoModel:               model.NewCouponInfoModel(db, c.CacheRedis),

		BatchInfoCache: make(map[string]*model.CouponBatchInfo),
	}

	// 启动批次信息加载
	svc.loadBatchInfo()
	go svc.loadBatchInfoLoop()

	return svc
}

// loadBatchInfo 加载批次信息到内存
func (s *ServiceContext) loadBatchInfo() {
	ctx := context.Background()
	list, err := s.CouponBatchInfoModel.FindAll(ctx)
	if err != nil {
		logx.Errorf("loadBatchInfo error: %v", err)
		return
	}

	s.batchMu.Lock()
	defer s.batchMu.Unlock()

	tmp := make(map[string]*model.CouponBatchInfo, len(list))
	for _, v := range list {
		tmp[v.BatchToken] = v
	}
	s.BatchInfoCache = tmp
	logx.Infof("loadBatchInfo success, count: %d", len(tmp))
}

// loadBatchInfoLoop 定时刷新批次信息
func (s *ServiceContext) loadBatchInfoLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.loadBatchInfo()
	}
}

// GetBatchInfo 获取批次信息
func (s *ServiceContext) GetBatchInfo(batchToken string) *model.CouponBatchInfo {
	s.batchMu.RLock()
	defer s.batchMu.RUnlock()
	return s.BatchInfoCache[batchToken]
}

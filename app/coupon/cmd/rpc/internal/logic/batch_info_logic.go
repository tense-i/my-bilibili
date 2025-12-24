package logic

import (
	"context"
	"errors"

	"mybilibili/app/coupon/cmd/rpc/coupon"
	"mybilibili/app/coupon/cmd/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchInfoLogic {
	return &BatchInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BatchInfo 批次信息
func (l *BatchInfoLogic) BatchInfo(in *coupon.BatchInfoReq) (*coupon.BatchInfoResp, error) {
	// 优先从内存缓存获取
	batch := l.svcCtx.GetBatchInfo(in.BatchToken)
	if batch != nil {
		return &coupon.BatchInfoResp{
			Id:           batch.Id,
			AppId:        batch.AppId,
			Name:         batch.Name,
			BatchToken:   batch.BatchToken,
			MaxCount:     batch.MaxCount,
			CurrentCount: batch.CurrentCount,
			StartTime:    batch.StartTime,
			ExpireTime:   batch.ExpireTime,
			ExpireDay:    batch.ExpireDay,
			LimitCount:   batch.LimitCount,
			FullAmount:   batch.FullAmount,
			Amount:       batch.Amount,
			State:        int32(batch.State),
			CouponType:   int32(batch.CouponType),
		}, nil
	}

	// 从数据库查询
	batchInfo, err := l.svcCtx.CouponBatchInfoModel.FindByBatchToken(l.ctx, in.BatchToken)
	if err != nil {
		l.Errorf("FindByBatchToken error: %v", err)
		return nil, errors.New("系统错误")
	}
	if batchInfo == nil {
		return nil, errors.New("批次不存在")
	}

	return &coupon.BatchInfoResp{
		Id:           batchInfo.Id,
		AppId:        batchInfo.AppId,
		Name:         batchInfo.Name,
		BatchToken:   batchInfo.BatchToken,
		MaxCount:     batchInfo.MaxCount,
		CurrentCount: batchInfo.CurrentCount,
		StartTime:    batchInfo.StartTime,
		ExpireTime:   batchInfo.ExpireTime,
		ExpireDay:    batchInfo.ExpireDay,
		LimitCount:   batchInfo.LimitCount,
		FullAmount:   batchInfo.FullAmount,
		Amount:       batchInfo.Amount,
		State:        int32(batchInfo.State),
		CouponType:   int32(batchInfo.CouponType),
	}, nil
}

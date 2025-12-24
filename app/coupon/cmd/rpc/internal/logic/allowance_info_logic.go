package logic

import (
	"context"
	"errors"

	"mybilibili/app/coupon/cmd/rpc/coupon"
	"mybilibili/app/coupon/cmd/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AllowanceInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAllowanceInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllowanceInfoLogic {
	return &AllowanceInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AllowanceInfo 代金券详情
func (l *AllowanceInfoLogic) AllowanceInfo(in *coupon.AllowanceInfoReq) (*coupon.AllowanceInfoResp, error) {
	cp, err := l.svcCtx.CouponAllowanceInfoModel.FindByToken(l.ctx, in.Mid, in.CouponToken)
	if err != nil {
		l.Errorf("FindByToken error: %v", err)
		return nil, errors.New("系统错误")
	}
	if cp == nil {
		return nil, errors.New("优惠券不存在")
	}

	return &coupon.AllowanceInfoResp{
		Id:          cp.Id,
		CouponToken: cp.CouponToken,
		Mid:         cp.Mid,
		State:       cp.State,
		StartTime:   cp.StartTime,
		ExpireTime:  cp.ExpireTime,
		Amount:      cp.Amount,
		FullAmount:  cp.FullAmount,
		BatchToken:  cp.BatchToken,
	}, nil
}

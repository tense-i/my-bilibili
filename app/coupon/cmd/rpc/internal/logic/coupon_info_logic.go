package logic

import (
	"context"
	"errors"

	"mybilibili/app/coupon/cmd/rpc/coupon"
	"mybilibili/app/coupon/cmd/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CouponInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCouponInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponInfoLogic {
	return &CouponInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CouponInfo 观影券详情
func (l *CouponInfoLogic) CouponInfo(in *coupon.CouponInfoReq) (*coupon.CouponInfoResp, error) {
	cp, err := l.svcCtx.CouponInfoModel.FindByToken(l.ctx, in.Mid, in.CouponToken)
	if err != nil {
		l.Errorf("FindByToken error: %v", err)
		return nil, errors.New("系统错误")
	}
	if cp == nil {
		return nil, errors.New("优惠券不存在")
	}

	return &coupon.CouponInfoResp{
		Info: &coupon.CouponInfo{
			Id:          cp.Id,
			CouponToken: cp.CouponToken,
			Mid:         cp.Mid,
			State:       cp.State,
			StartTime:   cp.StartTime,
			ExpireTime:  cp.ExpireTime,
			Origin:      int64(cp.Origin),
			CouponType:  int64(cp.CouponType),
			OrderNo:     cp.OrderNO,
			Oid:         cp.Oid,
			Remark:      cp.Remark,
		},
	}, nil
}

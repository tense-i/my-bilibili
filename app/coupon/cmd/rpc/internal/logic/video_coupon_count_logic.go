package logic

import (
	"context"

	"mybilibili/app/coupon/cmd/rpc/coupon"
	"mybilibili/app/coupon/cmd/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoCouponCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVideoCouponCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoCouponCountLogic {
	return &VideoCouponCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// VideoCouponCount 观影券数量
func (l *VideoCouponCountLogic) VideoCouponCount(in *coupon.VideoCouponCountReq) (*coupon.VideoCouponCountResp, error) {
	count, err := l.svcCtx.CouponInfoModel.CountByMidAndType(l.ctx, in.Mid, int8(in.CouponType))
	if err != nil {
		l.Errorf("CountByMidAndType error: %v", err)
		return &coupon.VideoCouponCountResp{Count: 0}, nil
	}

	return &coupon.VideoCouponCountResp{Count: count}, nil
}

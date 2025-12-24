// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package allowance

import (
	"context"

	"mybilibili/app/coupon/cmd/api/internal/svc"
	"mybilibili/app/coupon/cmd/api/internal/types"
	"mybilibili/app/coupon/cmd/rpc/coupon"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelAllowanceCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 取消使用代金券
func NewCancelAllowanceCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelAllowanceCouponLogic {
	return &CancelAllowanceCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelAllowanceCouponLogic) CancelAllowanceCoupon(req *types.CancelAllowanceCouponReq) error {
	_, err := l.svcCtx.CouponRpc.CancelAllowanceCoupon(l.ctx, &coupon.CancelAllowanceCouponReq{
		Mid:         req.Mid,
		CouponToken: req.CouponToken,
	})
	return err
}

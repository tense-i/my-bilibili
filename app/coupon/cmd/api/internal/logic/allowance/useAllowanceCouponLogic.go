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

type UseAllowanceCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 使用代金券
func NewUseAllowanceCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UseAllowanceCouponLogic {
	return &UseAllowanceCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UseAllowanceCouponLogic) UseAllowanceCoupon(req *types.UseAllowanceCouponReq) error {
	_, err := l.svcCtx.CouponRpc.UseAllowanceCoupon(l.ctx, &coupon.UseAllowanceCouponReq{
		Mid:            req.Mid,
		CouponToken:    req.CouponToken,
		Remark:         req.Remark,
		OrderNo:        req.OrderNo,
		Price:          req.Price,
		Platform:       req.Platform,
		ProdLimMonth:   req.ProdLimMonth,
		ProdLimRenewal: req.ProdLimRenewal,
	})
	return err
}

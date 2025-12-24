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

type ReceiveAllowanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 领取代金券
func NewReceiveAllowanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReceiveAllowanceLogic {
	return &ReceiveAllowanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReceiveAllowanceLogic) ReceiveAllowance(req *types.ReceiveAllowanceReq) (resp *types.ReceiveAllowanceResp, err error) {
	rpcResp, err := l.svcCtx.CouponRpc.ReceiveAllowance(l.ctx, &coupon.ReceiveAllowanceReq{
		Mid:        req.Mid,
		BatchToken: req.BatchToken,
		OrderNo:    req.OrderNo,
		Appkey:     req.Appkey,
	})
	if err != nil {
		return nil, err
	}

	return &types.ReceiveAllowanceResp{
		CouponToken: rpcResp.CouponToken,
	}, nil
}

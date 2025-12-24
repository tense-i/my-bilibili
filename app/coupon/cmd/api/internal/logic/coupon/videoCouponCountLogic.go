// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package coupon

import (
	"context"

	"mybilibili/app/coupon/cmd/api/internal/svc"
	"mybilibili/app/coupon/cmd/api/internal/types"
	"mybilibili/app/coupon/cmd/rpc/coupon"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoCouponCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 观影券数量
func NewVideoCouponCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoCouponCountLogic {
	return &VideoCouponCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoCouponCountLogic) VideoCouponCount(req *types.VideoCouponCountReq) (resp *types.VideoCouponCountResp, err error) {
	rpcResp, err := l.svcCtx.CouponRpc.VideoCouponCount(l.ctx, &coupon.VideoCouponCountReq{
		Mid:        req.Mid,
		CouponType: req.CouponType,
	})
	if err != nil {
		return nil, err
	}

	return &types.VideoCouponCountResp{
		Count: rpcResp.Count,
	}, nil
}

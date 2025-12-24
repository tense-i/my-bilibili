// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package code

import (
	"context"

	"mybilibili/app/coupon/cmd/api/internal/svc"
	"mybilibili/app/coupon/cmd/api/internal/types"
	"mybilibili/app/coupon/cmd/rpc/coupon"

	"github.com/zeromicro/go-zero/core/logx"
)

type UseCouponCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 使用兑换码
func NewUseCouponCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UseCouponCodeLogic {
	return &UseCouponCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UseCouponCodeLogic) UseCouponCode(req *types.UseCouponCodeReq) (resp *types.UseCouponCodeResp, err error) {
	// TODO: 从请求头获取 IP
	ip := "127.0.0.1"

	rpcResp, err := l.svcCtx.CouponRpc.UseCouponCode(l.ctx, &coupon.UseCouponCodeReq{
		Token:  req.Token,
		Code:   req.Code,
		Verify: req.Verify,
		Ip:     ip,
		Mid:    req.Mid,
	})
	if err != nil {
		return nil, err
	}

	return &types.UseCouponCodeResp{
		CouponToken:          rpcResp.CouponToken,
		CouponAmount:         rpcResp.CouponAmount,
		FullAmount:           rpcResp.FullAmount,
		PlatformLimitExplain: rpcResp.PlatformLimitExplain,
		ProductLimitMonth:    rpcResp.ProductLimitMonth,
		ProductLimitRenewal:  rpcResp.ProductLimitRenewal,
	}, nil
}

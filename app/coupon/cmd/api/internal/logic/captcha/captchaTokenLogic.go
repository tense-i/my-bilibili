// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package captcha

import (
	"context"

	"mybilibili/app/coupon/cmd/api/internal/svc"
	"mybilibili/app/coupon/cmd/api/internal/types"
	"mybilibili/app/coupon/cmd/rpc/coupon"

	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取验证码Token
func NewCaptchaTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaTokenLogic {
	return &CaptchaTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CaptchaTokenLogic) CaptchaToken(req *types.CaptchaTokenReq) (resp *types.CaptchaTokenResp, err error) {
	rpcResp, err := l.svcCtx.CouponRpc.CaptchaToken(l.ctx, &coupon.CaptchaTokenReq{})
	if err != nil {
		return nil, err
	}

	return &types.CaptchaTokenResp{
		Token: rpcResp.Token,
		Url:   rpcResp.Url,
	}, nil
}

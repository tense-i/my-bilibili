package logic

import (
	"context"

	"mybilibili/app/coupon/cmd/rpc/coupon"
	"mybilibili/app/coupon/cmd/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCaptchaTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaTokenLogic {
	return &CaptchaTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CaptchaToken 获取验证码Token
// 注：实际生产环境需要对接验证码服务，这里返回模拟数据
func (l *CaptchaTokenLogic) CaptchaToken(in *coupon.CaptchaTokenReq) (*coupon.CaptchaTokenResp, error) {
	// TODO: 对接实际的验证码服务
	// 这里返回模拟数据，实际需要调用验证码服务获取token和url
	return &coupon.CaptchaTokenResp{
		Token: "mock_captcha_token",
		Url:   "https://captcha.bilibili.com/verify",
	}, nil
}

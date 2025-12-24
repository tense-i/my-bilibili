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

type UsableAllowanceCouponV2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取可用代金券V2
func NewUsableAllowanceCouponV2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *UsableAllowanceCouponV2Logic {
	return &UsableAllowanceCouponV2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UsableAllowanceCouponV2Logic) UsableAllowanceCouponV2(req *types.UsableAllowanceCouponV2Req) (resp *types.UsableAllowanceCouponV2Resp, err error) {
	// 转换请求参数
	priceInfos := make([]*coupon.PriceInfo, 0, len(req.PriceInfo))
	for _, p := range req.PriceInfo {
		priceInfos = append(priceInfos, &coupon.PriceInfo{
			Price:          p.Price,
			Plat:           p.Plat,
			ProdLimMonth:   p.ProdLimMonth,
			ProdLimRenewal: p.ProdLimRenewal,
		})
	}

	// 调用 RPC
	rpcResp, err := l.svcCtx.CouponRpc.UsableAllowanceCouponV2(l.ctx, &coupon.UsableAllowanceCouponV2Req{
		Mid:       req.Mid,
		PriceInfo: priceInfos,
	})
	if err != nil {
		return nil, err
	}

	// 转换响应
	resp = &types.UsableAllowanceCouponV2Resp{
		CouponTip: rpcResp.CouponTip,
	}

	if rpcResp.CouponInfo != nil {
		resp.CouponInfo = &types.CouponAllowancePanelInfo{
			CouponToken:         rpcResp.CouponInfo.CouponToken,
			CouponAmount:        rpcResp.CouponInfo.CouponAmount,
			State:               rpcResp.CouponInfo.State,
			FullLimitExplain:    rpcResp.CouponInfo.FullLimitExplain,
			ScopeExplain:        rpcResp.CouponInfo.ScopeExplain,
			FullAmount:          rpcResp.CouponInfo.FullAmount,
			CouponDiscountPrice: rpcResp.CouponInfo.CouponDiscountPrice,
			StartTime:           rpcResp.CouponInfo.StartTime,
			ExpireTime:          rpcResp.CouponInfo.ExpireTime,
			Selected:            rpcResp.CouponInfo.Selected,
			DisablesExplains:    rpcResp.CouponInfo.DisablesExplains,
			OrderNo:             rpcResp.CouponInfo.OrderNo,
			Name:                rpcResp.CouponInfo.Name,
			Usable:              rpcResp.CouponInfo.Usable,
		}
	}

	return resp, nil
}

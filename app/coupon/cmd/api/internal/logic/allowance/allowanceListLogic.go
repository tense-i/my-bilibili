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

type AllowanceListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 代金券列表
func NewAllowanceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllowanceListLogic {
	return &AllowanceListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllowanceListLogic) AllowanceList(req *types.AllowanceListReq) (resp *types.AllowanceListResp, err error) {
	rpcResp, err := l.svcCtx.CouponRpc.AllowanceList(l.ctx, &coupon.AllowanceListReq{
		Mid:   req.Mid,
		State: req.State,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*types.CouponAllowancePanelInfo, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, &types.CouponAllowancePanelInfo{
			CouponToken:         item.CouponToken,
			CouponAmount:        item.CouponAmount,
			State:               item.State,
			FullLimitExplain:    item.FullLimitExplain,
			ScopeExplain:        item.ScopeExplain,
			FullAmount:          item.FullAmount,
			CouponDiscountPrice: item.CouponDiscountPrice,
			StartTime:           item.StartTime,
			ExpireTime:          item.ExpireTime,
			Selected:            item.Selected,
			DisablesExplains:    item.DisablesExplains,
			OrderNo:             item.OrderNo,
			Name:                item.Name,
			Usable:              item.Usable,
		})
	}

	return &types.AllowanceListResp{List: list}, nil
}

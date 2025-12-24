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

type CouponPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 观影券分页列表
func NewCouponPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponPageLogic {
	return &CouponPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CouponPageLogic) CouponPage(req *types.CouponPageReq) (resp *types.CouponPageResp, err error) {
	rpcResp, err := l.svcCtx.CouponRpc.CouponPage(l.ctx, &coupon.CouponPageReq{
		Mid:   req.Mid,
		State: req.State,
		Pn:    req.Pn,
		Ps:    req.Ps,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*types.CouponPageItem, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, &types.CouponPageItem{
			Id:    item.Id,
			Title: item.Title,
			Time:  item.Time,
			RefId: item.RefId,
		})
	}

	return &types.CouponPageResp{
		Count: rpcResp.Count,
		List:  list,
	}, nil
}

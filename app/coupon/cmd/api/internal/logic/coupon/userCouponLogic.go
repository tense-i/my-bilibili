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

type UserCouponLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户观影券列表
func NewUserCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCouponLogic {
	return &UserCouponLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserCouponLogic) UserCoupon(req *types.UserCouponReq) (resp *types.UserCouponResp, err error) {
	rpcResp, err := l.svcCtx.CouponRpc.UserCoupon(l.ctx, &coupon.UserCouponReq{
		Mid:        req.Mid,
		CouponType: req.CouponType,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*types.CouponInfoItem, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, &types.CouponInfoItem{
			Id:          item.Id,
			CouponToken: item.CouponToken,
			State:       item.State,
			StartTime:   item.StartTime,
			ExpireTime:  item.ExpireTime,
		})
	}

	return &types.UserCouponResp{List: list}, nil
}

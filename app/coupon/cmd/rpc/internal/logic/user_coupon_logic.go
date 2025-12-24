package logic

import (
	"context"

	"mybilibili/app/coupon/cmd/rpc/coupon"
	"mybilibili/app/coupon/cmd/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCouponLogic {
	return &UserCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserCoupon 用户观影券列表
func (l *UserCouponLogic) UserCoupon(in *coupon.UserCouponReq) (*coupon.UserCouponResp, error) {
	list, err := l.svcCtx.CouponInfoModel.FindByMidAndType(l.ctx, in.Mid, int8(in.CouponType))
	if err != nil {
		l.Errorf("FindByMidAndType error: %v", err)
		return &coupon.UserCouponResp{List: []*coupon.CouponInfo{}}, nil
	}

	result := make([]*coupon.CouponInfo, 0, len(list))
	for _, v := range list {
		result = append(result, &coupon.CouponInfo{
			Id:          v.Id,
			CouponToken: v.CouponToken,
			Mid:         v.Mid,
			State:       v.State,
			StartTime:   v.StartTime,
			ExpireTime:  v.ExpireTime,
			Origin:      int64(v.Origin),
			CouponType:  int64(v.CouponType),
			OrderNo:     v.OrderNO,
			Oid:         v.Oid,
			Remark:      v.Remark,
		})
	}

	return &coupon.UserCouponResp{List: result}, nil
}

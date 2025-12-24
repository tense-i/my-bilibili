package logic

import (
	"context"
	"errors"
	"time"

	"mybilibili/app/coupon/cmd/rpc/coupon"
	"mybilibili/app/coupon/cmd/rpc/internal/svc"
	"mybilibili/app/coupon/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UseCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUseCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UseCouponLogic {
	return &UseCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UseCoupon 使用观影券
func (l *UseCouponLogic) UseCoupon(in *coupon.UseCouponReq) (*coupon.UseCouponResp, error) {
	// 1. 获取用户可用的观影券
	list, err := l.svcCtx.CouponInfoModel.FindByMidAndType(l.ctx, in.Mid, int8(in.CouponType))
	if err != nil {
		l.Errorf("FindByMidAndType error: %v", err)
		return nil, errors.New("系统错误")
	}

	if len(list) == 0 {
		return &coupon.UseCouponResp{Ret: int32(model.UseFaild)}, nil
	}

	// 2. 选择一张可用的券（按过期时间排序，优先使用快过期的）
	now := time.Now().Unix()
	var selectedCoupon *model.CouponInfo
	for _, cp := range list {
		if cp.State == model.NotUsed && cp.ExpireTime > now && cp.StartTime <= now {
			if selectedCoupon == nil || cp.ExpireTime < selectedCoupon.ExpireTime {
				selectedCoupon = cp
			}
		}
	}

	if selectedCoupon == nil {
		return &coupon.UseCouponResp{Ret: int32(model.UseFaild)}, nil
	}

	// 3. 如果指定了 use_ver，检查版本
	if in.UseVer > 0 && selectedCoupon.Ver != in.UseVer {
		return &coupon.UseCouponResp{Ret: int32(model.UseFaild)}, nil
	}

	// 4. 更新券状态为已使用
	aff, err := l.svcCtx.CouponInfoModel.UpdateState(
		l.ctx,
		selectedCoupon.Id,
		selectedCoupon.Mid,
		model.Used,
		in.OrderNo,
		in.Oid,
		in.Remark,
		selectedCoupon.Ver,
	)
	if err != nil {
		l.Errorf("UpdateState error: %v", err)
		return nil, errors.New("系统错误")
	}
	if aff != 1 {
		return &coupon.UseCouponResp{Ret: int32(model.UseFaild)}, nil
	}

	return &coupon.UseCouponResp{
		Ret:         int32(model.UseSuccess),
		CouponToken: selectedCoupon.CouponToken,
	}, nil
}

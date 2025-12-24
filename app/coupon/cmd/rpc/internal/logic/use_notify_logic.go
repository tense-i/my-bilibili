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

type UseNotifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUseNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UseNotifyLogic {
	return &UseNotifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UseNotify 同步检查代金券是否可用（用于订单确认时）
func (l *UseNotifyLogic) UseNotify(in *coupon.UseNotifyReq) (*coupon.UseNotifyResp, error) {
	// 1. 根据订单号查询券
	cp, err := l.svcCtx.CouponAllowanceInfoModel.FindByOrderNo(l.ctx, in.Mid, in.OrderNo)
	if err != nil {
		l.Errorf("FindByOrderNo error: %v", err)
		return nil, errors.New("系统错误")
	}
	if cp == nil || cp.Mid != in.Mid {
		return nil, errors.New("优惠券使用失败")
	}

	// 2. 如果已经是已使用状态，直接返回
	if cp.State == model.Used {
		return &coupon.UseNotifyResp{
			CouponToken: cp.CouponToken,
			Amount:      cp.Amount,
			FullAmount:  cp.FullAmount,
		}, nil
	}

	// 3. 检查是否为使用中状态
	if cp.State != model.InUse {
		return nil, errors.New("优惠券状态异常")
	}

	// 4. 更新为已使用状态
	aff, err := l.svcCtx.CouponAllowanceInfoModel.UpdateStateToUsed(l.ctx, cp.Id, cp.Mid, in.OrderNo, cp.Ver)
	if err != nil {
		l.Errorf("UpdateStateToUsed error: %v", err)
		return nil, errors.New("系统错误")
	}
	if aff != 1 {
		return nil, errors.New("优惠券已被使用")
	}

	// 5. 记录变更日志
	changeLog := &model.CouponAllowanceChangeLog{
		CouponToken: cp.CouponToken,
		OrderNo:     in.OrderNo,
		Mid:         cp.Mid,
		State:       int8(model.Used),
		ChangeType:  model.AllowanceConsumeSuccess,
		Ctime:       time.Now(),
	}
	_, err = l.svcCtx.CouponAllowanceChangeLogModel.Insert(l.ctx, changeLog)
	if err != nil {
		l.Errorf("Insert change log error: %v", err)
	}

	return &coupon.UseNotifyResp{
		CouponToken: cp.CouponToken,
		Amount:      cp.Amount,
		FullAmount:  cp.FullAmount,
	}, nil
}

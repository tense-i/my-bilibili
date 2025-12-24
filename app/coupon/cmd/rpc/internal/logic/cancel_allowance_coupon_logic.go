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

type CancelAllowanceCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelAllowanceCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelAllowanceCouponLogic {
	return &CancelAllowanceCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CancelAllowanceCoupon 取消使用代金券
func (l *CancelAllowanceCouponLogic) CancelAllowanceCoupon(in *coupon.CancelAllowanceCouponReq) (*coupon.CancelAllowanceCouponResp, error) {
	// 1. 查询券信息
	cp, err := l.svcCtx.CouponAllowanceInfoModel.FindByToken(l.ctx, in.Mid, in.CouponToken)
	if err != nil {
		l.Errorf("FindByToken error: %v", err)
		return nil, errors.New("系统错误")
	}
	if cp == nil {
		return nil, errors.New("优惠券不存在")
	}

	// 2. 检查状态是否为使用中
	if cp.State != model.InUse {
		return nil, errors.New("优惠券状态不可取消")
	}

	// 3. 更新状态为未使用
	aff, err := l.svcCtx.CouponAllowanceInfoModel.UpdateStateToNotUsed(l.ctx, cp.Id, cp.Mid, cp.Ver)
	if err != nil {
		l.Errorf("UpdateStateToNotUsed error: %v", err)
		return nil, errors.New("系统错误")
	}
	if aff != 1 {
		return nil, errors.New("优惠券状态已变更")
	}

	// 4. 插入变更日志
	changeLog := &model.CouponAllowanceChangeLog{
		CouponToken: cp.CouponToken,
		OrderNo:     "",
		Mid:         cp.Mid,
		State:       int8(model.NotUsed),
		ChangeType:  model.AllowanceCancel,
		Ctime:       time.Now(),
	}
	_, err = l.svcCtx.CouponAllowanceChangeLogModel.Insert(l.ctx, changeLog)
	if err != nil {
		l.Errorf("Insert change log error: %v", err)
	}

	return &coupon.CancelAllowanceCouponResp{}, nil
}

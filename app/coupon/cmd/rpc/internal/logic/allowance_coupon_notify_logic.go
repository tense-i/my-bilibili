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

type AllowanceCouponNotifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAllowanceCouponNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllowanceCouponNotifyLogic {
	return &AllowanceCouponNotifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AllowanceCouponNotify 代金券使用通知（支付结果回调）
func (l *AllowanceCouponNotifyLogic) AllowanceCouponNotify(in *coupon.AllowanceCouponNotifyReq) (*coupon.AllowanceCouponNotifyResp, error) {
	// 1. 根据订单号查询券
	cp, err := l.svcCtx.CouponAllowanceInfoModel.FindByOrderNo(l.ctx, in.Mid, in.OrderNo)
	if err != nil {
		l.Errorf("FindByOrderNo error: %v", err)
		return nil, errors.New("系统错误")
	}
	if cp == nil {
		return nil, errors.New("优惠券不存在")
	}

	// 2. 检查状态是否为使用中
	if cp.State != model.InUse {
		return nil, errors.New("优惠券状态不可变更")
	}

	// 3. 根据支付状态决定券的最终状态
	var (
		state      int32
		changeType int8
		orderNo    string
		remark     string
	)

	switch int8(in.PayState) {
	case model.AllowanceUseFaild:
		// 支付失败，券退回未使用状态
		state = model.NotUsed
		changeType = model.AllowanceConsumeFaild
		orderNo = ""
		remark = ""
	case model.AllowanceUseSuccess:
		// 支付成功，券标记为已使用
		state = model.Used
		changeType = model.AllowanceConsumeSuccess
		orderNo = cp.OrderNO
		remark = cp.Remark
	default:
		return nil, errors.New("无效的支付状态")
	}

	// 4. 更新券状态
	aff, err := l.svcCtx.CouponAllowanceInfoModel.UpdateState(l.ctx, cp.Id, cp.Mid, state, orderNo, remark, cp.Ver)
	if err != nil {
		l.Errorf("UpdateState error: %v", err)
		return nil, errors.New("系统错误")
	}
	if aff != 1 {
		return nil, errors.New("优惠券状态已变更")
	}

	// 5. 插入变更日志
	changeLog := &model.CouponAllowanceChangeLog{
		CouponToken: cp.CouponToken,
		OrderNo:     orderNo,
		Mid:         cp.Mid,
		State:       int8(state),
		ChangeType:  changeType,
		Ctime:       time.Now(),
	}
	_, err = l.svcCtx.CouponAllowanceChangeLogModel.Insert(l.ctx, changeLog)
	if err != nil {
		l.Errorf("Insert change log error: %v", err)
	}

	return &coupon.AllowanceCouponNotifyResp{}, nil
}

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

type UseAllowanceCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUseAllowanceCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UseAllowanceCouponLogic {
	return &UseAllowanceCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UseAllowanceCoupon 使用代金券
func (l *UseAllowanceCouponLogic) UseAllowanceCoupon(in *coupon.UseAllowanceCouponReq) (*coupon.UseAllowanceCouponResp, error) {
	// 1. 验证券是否可用
	cp, err := l.judgeCouponUsable(
		in.Mid,
		in.Price,
		in.CouponToken,
		model.PlatformByName[in.Platform],
		int8(in.ProdLimMonth),
		int8(in.ProdLimRenewal),
	)
	if err != nil {
		return nil, err
	}

	// 2. 检查订单是否已绑定券
	exist, err := l.svcCtx.CouponAllowanceInfoModel.FindByOrderNo(l.ctx, in.Mid, in.OrderNo)
	if err != nil {
		l.Errorf("FindByOrderNo error: %v", err)
		return nil, errors.New("系统错误")
	}
	if exist != nil {
		return nil, errors.New("订单已绑定优惠券")
	}

	// 3. 更新券状态为使用中
	aff, err := l.svcCtx.CouponAllowanceInfoModel.UpdateState(l.ctx, cp.Id, cp.Mid, model.InUse, in.OrderNo, in.Remark, cp.Ver)
	if err != nil {
		l.Errorf("UpdateState error: %v", err)
		return nil, errors.New("系统错误")
	}
	if aff != 1 {
		return nil, errors.New("优惠券状态已变更")
	}

	// 4. 插入变更日志
	changeLog := &model.CouponAllowanceChangeLog{
		CouponToken: cp.CouponToken,
		OrderNo:     in.OrderNo,
		Mid:         cp.Mid,
		State:       int8(model.InUse),
		ChangeType:  model.AllowanceConsume,
		Ctime:       time.Now(),
	}
	_, err = l.svcCtx.CouponAllowanceChangeLogModel.Insert(l.ctx, changeLog)
	if err != nil {
		l.Errorf("Insert change log error: %v", err)
	}

	return &coupon.UseAllowanceCouponResp{}, nil
}

// judgeCouponUsable 验证券是否可用
func (l *UseAllowanceCouponLogic) judgeCouponUsable(mid int64, price float64, couponToken string, plat int, prodLimMonth, prodLimRenewal int8) (*model.CouponAllowanceInfo, error) {
	now := time.Now().Unix()

	// 1. 查询券信息
	cp, err := l.svcCtx.CouponAllowanceInfoModel.FindByToken(l.ctx, mid, couponToken)
	if err != nil {
		l.Errorf("FindByToken error: %v", err)
		return nil, errors.New("系统错误")
	}
	if cp == nil {
		return nil, errors.New("优惠券不存在")
	}

	// 2. 获取批次信息
	batch := l.svcCtx.GetBatchInfo(cp.BatchToken)
	if batch == nil {
		return nil, errors.New("批次不存在")
	}

	// 3. 批次状态检查
	if batch.State != model.BatchStateNormal {
		return nil, errors.New("优惠券已被冻结")
	}

	// 4. 券状态检查
	if cp.State != model.NotUsed {
		return nil, errors.New("优惠券已被使用")
	}

	// 5. 有效期检查
	if cp.StartTime > now || now > cp.ExpireTime {
		return nil, errors.New("优惠券不在有效期内")
	}

	// 6. 满额检查
	if cp.FullAmount > price {
		return nil, errors.New("未达到满额条件")
	}

	// 7. 平台限制检查
	if !platformLimit(batch.PlatformLimit, plat) {
		return nil, errors.New("当前平台不可使用")
	}

	// 8. 商品限制检查
	if !productLimit(batch.ProductLimitMonth, batch.ProductLimitRenewal, prodLimMonth, prodLimRenewal) {
		return nil, errors.New("当前商品不可使用")
	}

	return cp, nil
}

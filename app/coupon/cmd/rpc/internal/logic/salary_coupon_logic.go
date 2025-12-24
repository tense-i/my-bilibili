package logic

import (
	"context"
	"errors"
	"time"

	"mybilibili/app/coupon/cmd/rpc/coupon"
	"mybilibili/app/coupon/cmd/rpc/internal/svc"
	"mybilibili/app/coupon/model"
	"mybilibili/common/tool"

	"github.com/zeromicro/go-zero/core/logx"
)

type SalaryCouponLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSalaryCouponLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SalaryCouponLogic {
	return &SalaryCouponLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SalaryCoupon 发放优惠券
func (l *SalaryCouponLogic) SalaryCoupon(in *coupon.SalaryCouponReq) (*coupon.SalaryCouponResp, error) {
	// 1. 参数校验
	if in.Count <= 0 || in.Count > model.MaxSalaryCount {
		return nil, errors.New("发放数量不合法")
	}

	// 2. 获取批次信息
	batch := l.svcCtx.GetBatchInfo(in.BatchToken)
	if batch == nil {
		return nil, errors.New("批次不存在")
	}

	// 3. 检查批次状态
	if batch.State == model.BatchStateBlock {
		return nil, errors.New("批次已冻结")
	}

	// 4. 检查用户领取上限
	if batch.LimitCount >= 0 {
		count, err := l.svcCtx.CouponAllowanceInfoModel.CountByBatchToken(l.ctx, in.Mid, in.BatchToken)
		if err != nil {
			l.Errorf("CountByBatchToken error: %v", err)
			return nil, errors.New("系统错误")
		}
		if count+int64(in.Count) > batch.LimitCount {
			return nil, errors.New("已达到领取上限")
		}
	}

	// 5. 计算有效期
	var startTime, expireTime int64
	if batch.ExpireDay >= 0 {
		now, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
		startTime = now.Unix()
		expireTime = now.AddDate(0, 0, int(batch.ExpireDay+1)).Add(-1 * time.Second).Unix()
	} else {
		if batch.ExpireTime < time.Now().Unix() {
			return nil, errors.New("批次已过期")
		}
		startTime = batch.StartTime
		expireTime = batch.ExpireTime
	}

	// 6. 批量创建券
	coupons := make([]*model.CouponAllowanceInfo, in.Count)
	for i := int32(0); i < in.Count; i++ {
		coupons[i] = &model.CouponAllowanceInfo{
			CouponToken: tool.GenerateToken(),
			Mid:         in.Mid,
			State:       model.NotUsed,
			StartTime:   startTime,
			ExpireTime:  expireTime,
			Origin:      int8(in.Origin),
			BatchToken:  batch.BatchToken,
			Amount:      batch.Amount,
			FullAmount:  batch.FullAmount,
			AppId:       in.AppId,
			Ctime:       time.Now(),
		}
	}

	// 7. 批量插入
	_, err := l.svcCtx.CouponAllowanceInfoModel.BatchInsert(l.ctx, in.Mid, coupons)
	if err != nil {
		l.Errorf("BatchInsert error: %v", err)
		return nil, errors.New("系统错误")
	}

	// 8. 更新批次发放数量
	_, err = l.svcCtx.CouponBatchInfoModel.UpdateCurrentCount(l.ctx, in.BatchToken, int(in.Count))
	if err != nil {
		l.Errorf("UpdateCurrentCount error: %v", err)
	}

	return &coupon.SalaryCouponResp{}, nil
}

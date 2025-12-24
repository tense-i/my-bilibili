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

type ReceiveAllowanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReceiveAllowanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReceiveAllowanceLogic {
	return &ReceiveAllowanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ReceiveAllowance 领取代金券
func (l *ReceiveAllowanceLogic) ReceiveAllowance(in *coupon.ReceiveAllowanceReq) (*coupon.ReceiveAllowanceResp, error) {
	// 1. 检查是否已领取（幂等性）
	rlog, err := l.svcCtx.CouponReceiveLogModel.FindByAppkeyAndOrderNo(l.ctx, in.Appkey, in.OrderNo, model.CouponAllowance)
	if err != nil {
		l.Errorf("FindByAppkeyAndOrderNo error: %v", err)
		return nil, errors.New("系统错误")
	}
	if rlog != nil {
		// 已领取，返回之前的券token
		return &coupon.ReceiveAllowanceResp{CouponToken: rlog.CouponToken}, nil
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

	// 4. 检查批次发放上限
	if batch.MaxCount >= 0 && batch.CurrentCount >= batch.MaxCount {
		return nil, errors.New("批次已发放完毕")
	}

	// 5. 检查用户领取上限
	if batch.LimitCount >= 0 {
		count, err := l.svcCtx.CouponAllowanceInfoModel.CountByBatchToken(l.ctx, in.Mid, in.BatchToken)
		if err != nil {
			l.Errorf("CountByBatchToken error: %v", err)
			return nil, errors.New("系统错误")
		}
		if count >= batch.LimitCount {
			return nil, errors.New("已达到领取上限")
		}
	}

	// 6. 计算有效期
	var startTime, expireTime int64
	if batch.ExpireDay >= 0 {
		now, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
		startTime = time.Now().Unix()
		expireTime = now.AddDate(0, 0, int(batch.ExpireDay+1)).Add(-1 * time.Second).Unix()
	} else {
		if batch.ExpireTime < time.Now().Unix() {
			return nil, errors.New("批次已过期")
		}
		startTime = batch.StartTime
		expireTime = batch.ExpireTime
	}

	// 7. 生成券token
	couponToken := tool.GenerateToken()

	// 8. 创建券
	cpInfo := &model.CouponAllowanceInfo{
		CouponToken: couponToken,
		Mid:         in.Mid,
		State:       model.NotUsed,
		StartTime:   startTime,
		ExpireTime:  expireTime,
		Origin:      model.AllowanceBusinessReceive,
		BatchToken:  batch.BatchToken,
		Amount:      batch.Amount,
		FullAmount:  batch.FullAmount,
		AppId:       batch.AppId,
		Ctime:       time.Now(),
	}

	_, err = l.svcCtx.CouponAllowanceInfoModel.Insert(l.ctx, cpInfo)
	if err != nil {
		l.Errorf("Insert coupon error: %v", err)
		return nil, errors.New("系统错误")
	}

	// 9. 更新批次发放数量
	_, err = l.svcCtx.CouponBatchInfoModel.UpdateCurrentCount(l.ctx, in.BatchToken, 1)
	if err != nil {
		l.Errorf("UpdateCurrentCount error: %v", err)
	}

	// 10. 记录领取日志
	receiveLog := &model.CouponReceiveLog{
		Appkey:      in.Appkey,
		OrderNo:     in.OrderNo,
		Mid:         in.Mid,
		CouponToken: couponToken,
		CouponType:  model.CouponAllowance,
	}
	_, err = l.svcCtx.CouponReceiveLogModel.Insert(l.ctx, receiveLog)
	if err != nil {
		l.Errorf("Insert receive log error: %v", err)
	}

	// 11. 记录变更日志
	changeLog := &model.CouponAllowanceChangeLog{
		CouponToken: couponToken,
		Mid:         in.Mid,
		State:       int8(model.NotUsed),
		ChangeType:  model.AllowanceReceive,
		Ctime:       time.Now(),
	}
	_, err = l.svcCtx.CouponAllowanceChangeLogModel.Insert(l.ctx, changeLog)
	if err != nil {
		l.Errorf("Insert change log error: %v", err)
	}

	return &coupon.ReceiveAllowanceResp{CouponToken: couponToken}, nil
}

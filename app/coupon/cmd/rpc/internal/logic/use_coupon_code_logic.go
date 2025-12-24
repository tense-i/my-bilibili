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

type UseCouponCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUseCouponCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UseCouponCodeLogic {
	return &UseCouponCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UseCouponCode 使用兑换码
func (l *UseCouponCodeLogic) UseCouponCode(in *coupon.UseCouponCodeReq) (*coupon.UseCouponCodeResp, error) {
	// 1. 查询兑换码
	code, err := l.svcCtx.CouponCodeModel.FindByCode(l.ctx, in.Code)
	if err != nil {
		l.Errorf("FindByCode error: %v", err)
		return nil, errors.New("系统错误")
	}
	if code == nil {
		return nil, errors.New("兑换码不存在")
	}

	// 2. 检查兑换码状态
	if code.State == model.CodeStateUsed {
		return nil, errors.New("兑换码已被使用")
	}
	if code.State == model.CodeStateBlock {
		return nil, errors.New("兑换码已被冻结")
	}

	// 3. 获取批次信息
	batch := l.svcCtx.GetBatchInfo(code.BatchToken)
	if batch == nil {
		return nil, errors.New("批次不存在")
	}

	// 4. 检查批次状态
	if batch.State == model.BatchStateBlock {
		return nil, errors.New("批次已冻结")
	}

	// 5. 检查批次有效期
	if batch.ExpireDay == -1 && batch.ExpireTime < time.Now().Unix() {
		return nil, errors.New("批次已过期")
	}

	// 6. 检查用户领取上限
	if batch.LimitCount != -1 {
		count, err := l.svcCtx.CouponCodeModel.CountByMidAndBatchToken(l.ctx, in.Mid, batch.BatchToken)
		if err != nil {
			l.Errorf("CountByMidAndBatchToken error: %v", err)
			return nil, errors.New("系统错误")
		}
		if count >= batch.LimitCount {
			return nil, errors.New("已达到领取上限")
		}
	}

	// 7. 检查批次发放上限
	if batch.MaxCount != -1 && batch.CurrentCount >= batch.MaxCount {
		return nil, errors.New("批次已发放完毕")
	}

	// 8. 计算有效期
	var startTime, expireTime int64
	if batch.ExpireDay >= 0 {
		now, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
		startTime = time.Now().Unix()
		expireTime = now.AddDate(0, 0, int(batch.ExpireDay+1)).Add(-1 * time.Second).Unix()
	} else {
		startTime = batch.StartTime
		expireTime = batch.ExpireTime
	}

	// 9. 生成券token
	couponToken := tool.GenerateToken()

	// 10. 创建券
	cpInfo := &model.CouponAllowanceInfo{
		CouponToken: couponToken,
		Mid:         in.Mid,
		State:       model.NotUsed,
		StartTime:   startTime,
		ExpireTime:  expireTime,
		Origin:      model.AllowanceCodeOpen,
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

	// 11. 更新兑换码状态
	aff, err := l.svcCtx.CouponCodeModel.UpdateState(l.ctx, in.Code, model.CodeStateUsed, in.Mid, couponToken, code.Ver)
	if err != nil {
		l.Errorf("UpdateState error: %v", err)
		return nil, errors.New("系统错误")
	}
	if aff != 1 {
		return nil, errors.New("兑换码已被使用")
	}

	// 12. 更新批次发放数量
	_, err = l.svcCtx.CouponBatchInfoModel.UpdateCurrentCount(l.ctx, batch.BatchToken, 1)
	if err != nil {
		l.Errorf("UpdateCurrentCount error: %v", err)
	}

	return &coupon.UseCouponCodeResp{
		CouponToken:          couponToken,
		CouponAmount:         batch.Amount,
		FullAmount:           batch.FullAmount,
		ProductLimitMonth:    int32(batch.ProductLimitMonth),
		ProductLimitRenewal:  int32(batch.ProductLimitRenewal),
		PlatformLimitExplain: l.getPlatformLimitExplain(batch.PlatformLimit),
	}, nil
}

// getPlatformLimitExplain 获取平台限制说明
func (l *UseCouponCodeLogic) getPlatformLimitExplain(platformLimit string) string {
	if platformLimit == "" {
		return "全平台通用"
	}
	return "限指定平台使用"
}

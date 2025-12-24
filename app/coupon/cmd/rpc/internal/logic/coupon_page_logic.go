package logic

import (
	"context"

	"mybilibili/app/coupon/cmd/rpc/coupon"
	"mybilibili/app/coupon/cmd/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CouponPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCouponPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponPageLogic {
	return &CouponPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CouponPage 观影券分页列表
func (l *CouponPageLogic) CouponPage(in *coupon.CouponPageReq) (*coupon.CouponPageResp, error) {
	// 默认分页参数
	pn := in.Pn
	ps := in.Ps
	if pn <= 0 {
		pn = 1
	}
	if ps <= 0 {
		ps = 20
	}
	if ps > 100 {
		ps = 100
	}

	// 获取总数
	count, err := l.svcCtx.CouponInfoModel.CountByMidAndState(l.ctx, in.Mid, in.State)
	if err != nil {
		l.Errorf("CountByMidAndState error: %v", err)
		return &coupon.CouponPageResp{Count: 0, List: []*coupon.CouponPageItem{}}, nil
	}

	if count == 0 {
		return &coupon.CouponPageResp{Count: 0, List: []*coupon.CouponPageItem{}}, nil
	}

	// 获取列表
	list, err := l.svcCtx.CouponInfoModel.FindByMidAndState(l.ctx, in.Mid, in.State, pn, ps)
	if err != nil {
		l.Errorf("FindByMidAndState error: %v", err)
		return &coupon.CouponPageResp{Count: count, List: []*coupon.CouponPageItem{}}, nil
	}

	result := make([]*coupon.CouponPageItem, 0, len(list))
	for _, v := range list {
		result = append(result, &coupon.CouponPageItem{
			Id:    v.Id,
			Title: l.getCouponTitle(v.CouponType),
			Time:  v.Ctime.Unix(),
			RefId: v.Oid,
			Tips:  l.getCouponTips(v.State, v.ExpireTime),
			Count: 1,
		})
	}

	return &coupon.CouponPageResp{
		Count: count,
		List:  result,
	}, nil
}

// getCouponTitle 获取券标题
func (l *CouponPageLogic) getCouponTitle(couponType int8) string {
	switch couponType {
	case 1:
		return "观影券"
	case 2:
		return "漫画券"
	default:
		return "优惠券"
	}
}

// getCouponTips 获取券提示
func (l *CouponPageLogic) getCouponTips(state int32, expireTime int64) string {
	switch state {
	case 0:
		return "未使用"
	case 2:
		return "已使用"
	default:
		return ""
	}
}

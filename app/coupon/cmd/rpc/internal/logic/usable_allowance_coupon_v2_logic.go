package logic

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"mybilibili/app/coupon/cmd/rpc/coupon"
	"mybilibili/app/coupon/cmd/rpc/internal/svc"
	"mybilibili/app/coupon/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UsableAllowanceCouponV2Logic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUsableAllowanceCouponV2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *UsableAllowanceCouponV2Logic {
	return &UsableAllowanceCouponV2Logic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UsableAllowanceCouponV2 获取可用代金券V2
func (l *UsableAllowanceCouponV2Logic) UsableAllowanceCouponV2(in *coupon.UsableAllowanceCouponV2Req) (*coupon.UsableAllowanceCouponV2Resp, error) {
	resp := &coupon.UsableAllowanceCouponV2Resp{
		CouponTip: model.CouponTipNotUse,
	}

	if len(in.PriceInfo) == 0 {
		return resp, nil
	}

	now := time.Now().Unix()

	// 1. 获取用户未使用的代金券
	notUsedList, err := l.svcCtx.CouponAllowanceInfoModel.FindByMidAndState(l.ctx, in.Mid, model.NotUsed, now)
	if err != nil {
		l.Errorf("FindByMidAndState NotUsed error: %v", err)
		return resp, nil
	}

	// 2. 获取用户使用中的代金券
	inUseList, err := l.svcCtx.CouponAllowanceInfoModel.FindByMidAndState(l.ctx, in.Mid, model.InUse, now)
	if err != nil {
		l.Errorf("FindByMidAndState InUse error: %v", err)
		return resp, nil
	}

	all := append(notUsedList, inUseList...)
	if len(all) == 0 {
		return resp, nil
	}

	selectedPrice := in.PriceInfo[0].Price
	var maxAmount float64
	var bestCoupon *coupon.CouponAllowancePanelInfo
	var selectedCoupon *coupon.CouponAllowancePanelInfo

	// 3. 遍历每个价格，找出最优券
	for _, priceInfo := range in.PriceInfo {
		usables, _, using := l.filterUsableCoupons(
			priceInfo.Price,
			all,
			int(priceInfo.Plat),
			int8(priceInfo.ProdLimMonth),
			int8(priceInfo.ProdLimRenewal),
		)

		availables := append(usables, using...)
		if len(availables) == 0 {
			continue
		}

		// 按金额排序
		sort.Slice(availables, func(i, j int) bool {
			if availables[i].CouponAmount == availables[j].CouponAmount {
				return availables[i].Usable > availables[j].Usable
			}
			return availables[i].CouponAmount > availables[j].CouponAmount
		})

		if availables[0].CouponAmount > maxAmount {
			maxAmount = availables[0].CouponAmount
			bestCoupon = availables[0]
		}

		if priceInfo.Price == selectedPrice {
			selectedCoupon = availables[0]
			break
		}
	}

	if bestCoupon == nil {
		return resp, nil
	}

	// 4. 设置返回提示
	switch {
	case bestCoupon.Usable == int32(model.AllowanceUsable):
		resp.CouponTip = fmt.Sprintf(model.CouponTipUse, bestCoupon.CouponAmount)
		resp.CouponInfo = bestCoupon
	case bestCoupon.State == model.InUse:
		resp.CouponTip = model.CouponTipInUse
	case selectedCoupon == nil:
		resp.CouponTip = model.CouponTipChooseOther
	}

	return resp, nil
}

// filterUsableCoupons 筛选可用券
func (l *UsableAllowanceCouponV2Logic) filterUsableCoupons(
	price float64,
	all []*model.CouponAllowanceInfo,
	plat int,
	prodLimMonth, prodLimRenewal int8,
) (usables, disables, using []*coupon.CouponAllowancePanelInfo) {
	now := time.Now().Unix()
	usables = make([]*coupon.CouponAllowancePanelInfo, 0)
	disables = make([]*coupon.CouponAllowancePanelInfo, 0)
	using = make([]*coupon.CouponAllowancePanelInfo, 0)

	for _, c := range all {
		// 获取批次信息
		batch := l.svcCtx.GetBatchInfo(c.BatchToken)
		if batch == nil {
			continue
		}

		if batch.State == model.BatchStateBlock {
			continue
		}

		ok := true
		explains := make([]string, 0)

		// 平台限制检查
		if !platformLimit(batch.PlatformLimit, plat) {
			explains = append(explains, model.CouponPlatformExplain)
			ok = false
		}

		// 商品限制检查
		if !productLimit(batch.ProductLimitMonth, batch.ProductLimitRenewal, prodLimMonth, prodLimRenewal) {
			explains = append(explains, model.CouponProductExplain)
			ok = false
		}

		// 满额检查
		if c.FullAmount > price {
			explains = append(explains, model.CouponFullAmountDissatisfy)
			ok = false
		}

		// 开始时间检查
		if c.StartTime > now {
			explains = append(explains, model.CouponNotInUsableTime)
			ok = false
		}

		// 使用中状态处理
		state := c.State
		if c.State == model.InUse {
			if len(explains) > 0 {
				state = model.NotUsed
			}
			if ok {
				using = append(using, convertToPanelInfo(c, explains, price, model.AllowanceDisables))
			}
			ok = false
		}

		if ok {
			usables = append(usables, convertToPanelInfo(c, explains, price, model.AllowanceUsable))
		} else {
			info := convertToPanelInfo(c, explains, price, model.AllowanceDisables)
			info.State = state
			disables = append(disables, info)
		}
	}

	// 可用券排序：按过期时间升序，金额降序
	if len(usables) > 0 {
		sort.Slice(usables, func(i, j int) bool {
			return usables[i].ExpireTime < usables[j].ExpireTime
		})
		sort.Slice(usables, func(i, j int) bool {
			return usables[i].CouponAmount > usables[j].CouponAmount
		})
		usables[0].Selected = 1
	}

	return
}

// platformLimit 平台限制验证
func platformLimit(pstr string, plat int) bool {
	if len(pstr) == 0 {
		return true
	}
	ps := strings.Split(pstr, ",")
	for _, v := range ps {
		if v == fmt.Sprintf("%d", plat) {
			return true
		}
	}
	return false
}

// productLimit 商品限制验证
func productLimit(bplm, bplr, plm, plr int8) bool {
	if bplm == model.ProdLimMonthNone && bplr == model.ProdLimRenewalAll {
		return true
	}
	if bplm == model.ProdLimMonthNone && bplr == plr {
		return true
	}
	if bplr == model.ProdLimRenewalAll && bplm == plm {
		return true
	}
	if bplm == plm && bplr == plr {
		return true
	}
	return false
}

// convertToPanelInfo 转换为面板信息
func convertToPanelInfo(c *model.CouponAllowanceInfo, explains []string, price float64, usable int8) *coupon.CouponAllowancePanelInfo {
	return &coupon.CouponAllowancePanelInfo{
		CouponToken:         c.CouponToken,
		CouponAmount:        c.Amount,
		State:               c.State,
		FullLimitExplain:    fmt.Sprintf("满%.0f元可用", c.FullAmount),
		FullAmount:          c.FullAmount,
		StartTime:           c.StartTime,
		ExpireTime:          c.ExpireTime,
		CouponDiscountPrice: price - c.Amount,
		DisablesExplains:    strings.Join(explains, ","),
		OrderNo:             c.OrderNO,
		Name:                "大会员代金券",
		Usable:              int32(usable),
	}
}

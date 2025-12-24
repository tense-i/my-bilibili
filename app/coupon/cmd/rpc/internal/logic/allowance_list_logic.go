package logic

import (
	"context"
	"fmt"
	"time"

	"mybilibili/app/coupon/cmd/rpc/coupon"
	"mybilibili/app/coupon/cmd/rpc/internal/svc"
	"mybilibili/app/coupon/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AllowanceListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAllowanceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllowanceListLogic {
	return &AllowanceListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AllowanceList 代金券列表
func (l *AllowanceListLogic) AllowanceList(in *coupon.AllowanceListReq) (*coupon.AllowanceListResp, error) {
	now := time.Now().Unix()
	stime := time.Now().AddDate(0, -3, 0) // 只查询3个月内的券

	list, err := l.svcCtx.CouponAllowanceInfoModel.FindList(l.ctx, in.Mid, in.State, now, stime)
	if err != nil {
		l.Errorf("FindList error: %v", err)
		return &coupon.AllowanceListResp{List: []*coupon.CouponAllowancePanelInfo{}}, nil
	}

	result := make([]*coupon.CouponAllowancePanelInfo, 0, len(list))
	for _, v := range list {
		state := v.State
		// 如果查询的是过期券，强制设置状态为过期
		if in.State == model.Expire {
			state = model.Expire
		}

		info := &coupon.CouponAllowancePanelInfo{
			CouponToken:      v.CouponToken,
			CouponAmount:     v.Amount,
			State:            state,
			FullLimitExplain: fmt.Sprintf("满%.0f元可用", v.FullAmount),
			FullAmount:       v.FullAmount,
			StartTime:        v.StartTime,
			ExpireTime:       v.ExpireTime,
			OrderNo:          v.OrderNO,
			Name:             "大会员代金券",
			Usable:           int32(model.AllowanceDisables),
		}

		// 获取批次信息补充说明
		batch := l.svcCtx.GetBatchInfo(v.BatchToken)
		if batch != nil {
			info.ScopeExplain = l.getScopeExplain(batch)
		}

		result = append(result, info)
	}

	return &coupon.AllowanceListResp{List: result}, nil
}

// getScopeExplain 获取使用范围说明
func (l *AllowanceListLogic) getScopeExplain(batch *model.CouponBatchInfo) string {
	// 根据平台限制和商品限制生成说明
	if batch.PlatformLimit == "" && batch.ProductLimitMonth == model.ProdLimMonthNone {
		return "全平台通用"
	}
	return "限指定商品使用"
}

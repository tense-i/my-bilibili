package search

import (
	"context"
	"strings"

	"mybilibili/app/search/cmd/api/internal/svc"
	"mybilibili/app/search/cmd/api/internal/types"
	"mybilibili/app/search/cmd/rpc/search"

	"github.com/zeromicro/go-zero/core/logx"
)

type DmDateSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 弹幕日期搜索
func NewDmDateSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DmDateSearchLogic {
	return &DmDateSearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DmDateSearchLogic) DmDateSearch(req *types.DmDateSearchReq) (*types.SearchResp, error) {
	// 解析 kw_fields
	var kwFields []string
	if req.KwFields != "" {
		kwFields = strings.Split(req.KwFields, ",")
	}

	// 解析 order 和 sort
	var order, sort []string
	if req.Order != "" {
		order = strings.Split(req.Order, ",")
	}
	if req.Sort != "" {
		sort = strings.Split(req.Sort, ",")
	}

	// 调用 RPC
	resp, err := l.svcCtx.SearchRpc.DmDateSearch(l.ctx, &search.DmDateSearchReq{
		Bsp: &search.BasicSearchParams{
			AppId:    req.AppId,
			Kw:       req.Kw,
			KwFields: kwFields,
			Order:    order,
			Sort:     sort,
			Pn:       req.Pn,
			Ps:       req.Ps,
			Debug:    req.Debug,
		},
		Oid:       req.Oid,
		Month:     req.Month,
		MonthFrom: req.MonthFrom,
		MonthTo:   req.MonthTo,
	})
	if err != nil {
		return nil, err
	}

	// 转换结果
	result := make([]string, len(resp.Result))
	for i, r := range resp.Result {
		result[i] = string(r)
	}

	return &types.SearchResp{
		Order:  resp.Order,
		Sort:   resp.Sort,
		Result: result,
		Page: types.Page{
			Pn:    resp.Page.Pn,
			Ps:    resp.Page.Ps,
			Total: resp.Page.Total,
		},
		Debug: resp.Debug,
	}, nil
}

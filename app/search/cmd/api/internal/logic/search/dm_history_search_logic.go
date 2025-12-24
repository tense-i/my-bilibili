package search

import (
	"context"
	"strconv"
	"strings"

	"mybilibili/app/search/cmd/api/internal/svc"
	"mybilibili/app/search/cmd/api/internal/types"
	"mybilibili/app/search/cmd/rpc/search"

	"github.com/zeromicro/go-zero/core/logx"
)

type DmHistorySearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDmHistorySearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DmHistorySearchLogic {
	return &DmHistorySearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DmHistorySearchLogic) DmHistorySearch(req *types.DmHistorySearchReq) (*types.SearchResp, error) {
	// 解析 states
	var states []int64
	if req.States != "" {
		for _, s := range strings.Split(req.States, ",") {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}
			if v, err := strconv.ParseInt(s, 10, 64); err == nil {
				states = append(states, v)
			}
		}
	}

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
	resp, err := l.svcCtx.SearchRpc.DmHistorySearch(l.ctx, &search.DmHistorySearchReq{
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
		States:    states,
		CtimeFrom: req.CtimeFrom,
		CtimeTo:   req.CtimeTo,
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

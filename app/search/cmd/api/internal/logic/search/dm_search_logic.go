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

type DmSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDmSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DmSearchLogic {
	return &DmSearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DmSearchLogic) DmSearch(req *types.DmSearchReq) (*types.SearchResp, error) {
	// 解析参数
	var kwFields, order, sort []string
	var states, attrFormat []int32

	if req.KwFields != "" {
		kwFields = strings.Split(req.KwFields, ",")
	}
	if req.Order != "" {
		order = strings.Split(req.Order, ",")
	}
	if req.Sort != "" {
		sort = strings.Split(req.Sort, ",")
	}
	if req.States != "" {
		states = parseIntSlice(req.States)
	}
	if req.AttrFormat != "" {
		attrFormat = parseIntSlice(req.AttrFormat)
	}

	// 调用 RPC
	resp, err := l.svcCtx.SearchRpc.DmSearch(l.ctx, &search.DmSearchReq{
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
		Oid:        req.Oid,
		Mid:        req.Mid,
		Mode:       req.Mode,
		Pool:       req.Pool,
		Progress:   req.Progress,
		States:     states,
		Type:       req.Type,
		AttrFormat: attrFormat,
		CtimeFrom:  req.CtimeFrom,
		CtimeTo:    req.CtimeTo,
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

// parseIntSlice 解析逗号分隔的整数字符串
func parseIntSlice(s string) []int32 {
	parts := strings.Split(s, ",")
	result := make([]int32, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if v, err := strconv.ParseInt(p, 10, 32); err == nil {
			result = append(result, int32(v))
		}
	}
	return result
}

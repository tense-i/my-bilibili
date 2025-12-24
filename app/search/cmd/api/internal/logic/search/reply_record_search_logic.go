package search

import (
	"context"
	"strings"

	"mybilibili/app/search/cmd/api/internal/svc"
	"mybilibili/app/search/cmd/api/internal/types"
	"mybilibili/app/search/cmd/rpc/search"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyRecordSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReplyRecordSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyRecordSearchLogic {
	return &ReplyRecordSearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReplyRecordSearchLogic) ReplyRecordSearch(req *types.ReplyRecordSearchReq) (*types.SearchResp, error) {
	// 解析参数
	var order, sort []string
	if req.Order != "" {
		order = strings.Split(req.Order, ",")
	}
	if req.Sort != "" {
		sort = strings.Split(req.Sort, ",")
	}

	// 调用 RPC
	resp, err := l.svcCtx.SearchRpc.ReplyRecordSearch(l.ctx, &search.ReplyRecordSearchReq{
		Bsp: &search.BasicSearchParams{
			AppId: req.AppId,
			Kw:    req.Kw,
			Order: order,
			Sort:  sort,
			Pn:    req.Pn,
			Ps:    req.Ps,
			Debug: req.Debug,
		},
		Mid:       req.Mid,
		Types:     parseInt64Slice(req.Types),
		States:    parseInt64Slice(req.States),
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

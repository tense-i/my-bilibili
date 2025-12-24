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

type PgcMediaSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPgcMediaSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PgcMediaSearchLogic {
	return &PgcMediaSearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PgcMediaSearchLogic) PgcMediaSearch(req *types.PgcMediaSearchReq) (*types.SearchResp, error) {
	// 解析参数
	var order, sort []string
	if req.Order != "" {
		order = strings.Split(req.Order, ",")
	}
	if req.Sort != "" {
		sort = strings.Split(req.Sort, ",")
	}

	// 调用 RPC
	resp, err := l.svcCtx.SearchRpc.PgcMediaSearch(l.ctx, &search.PgcMediaSearchReq{
		Bsp: &search.BasicSearchParams{
			AppId: req.AppId,
			Kw:    req.Kw,
			Order: order,
			Sort:  sort,
			Pn:    req.Pn,
			Ps:    req.Ps,
			Debug: req.Debug,
		},
		MediaIds:        parseInt64Slice(req.MediaIds),
		SeasonIds:       parseInt64Slice(req.SeasonIds),
		SeasonTypes:     parseInt64Slice(req.SeasonTypes),
		StyleIds:        parseInt64Slice(req.StyleIds),
		Status:          req.Status,
		ReleaseDateFrom: req.ReleaseDateFrom,
		ReleaseDateTo:   req.ReleaseDateTo,
		SeasonIdFrom:    req.SeasonIdFrom,
		SeasonIdTo:      req.SeasonIdTo,
		ProducerIds:     parseInt64Slice(req.ProducerIds),
		IsDeleted:       req.IsDeleted,
		AreaIds:         parseStringSlice(req.AreaIds),
		ScoreFrom:       req.ScoreFrom,
		ScoreTo:         req.ScoreTo,
		IsFinish:        req.IsFinish,
		SeasonVersions:  parseInt64Slice(req.SeasonVersions),
		SeasonStatuses:  parseInt64Slice(req.SeasonStatuses),
		PubTimeFrom:     req.PubTimeFrom,
		PubTimeTo:       req.PubTimeTo,
		SeasonMonths:    parseInt64Slice(req.SeasonMonths),
		LatestTimeFrom:  req.LatestTimeFrom,
		LatestTimeTo:    req.LatestTimeTo,
		CopyrightInfos:  parseStringSlice(req.CopyrightInfos),
		CtimeFrom:       req.CtimeFrom,
		CtimeTo:         req.CtimeTo,
		MtimeFrom:       req.MtimeFrom,
		MtimeTo:         req.MtimeTo,
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

// parseInt64Slice 解析逗号分隔的整数字符串为 int64 切片
func parseInt64Slice(s string) []int64 {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]int64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if v, err := strconv.ParseInt(p, 10, 64); err == nil {
			result = append(result, v)
		}
	}
	return result
}

// parseStringSlice 解析逗号分隔的字符串
func parseStringSlice(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

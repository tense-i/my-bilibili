package logic

import (
	"context"
	"fmt"

	"mybilibili/app/search/cmd/rpc/internal/svc"
	"mybilibili/app/search/cmd/rpc/search"
	"mybilibili/app/search/model"

	"github.com/olivere/elastic/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

type DmSearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDmSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DmSearchLogic {
	return &DmSearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DmSearchLogic) DmSearch(in *search.DmSearchReq) (*search.SearchResp, error) {
	// 构建索引名
	indexName := fmt.Sprintf("dm_search_%03d", in.Oid%1000)

	// 构建查询
	query := elastic.NewBoolQuery()

	// 关键词搜索
	if in.Bsp != nil && in.Bsp.Kw != "" && len(in.Bsp.KwFields) > 0 {
		query = query.Must(elastic.NewRegexpQuery(in.Bsp.KwFields[0], ".*"+in.Bsp.Kw+".*"))
	}

	// 过滤条件
	if in.Oid > 0 {
		query = query.Filter(elastic.NewTermQuery("oid", in.Oid))
	}
	if in.Mid > 0 {
		query = query.Filter(elastic.NewTermQuery("mid", in.Mid))
	}
	if in.Mode >= 0 {
		query = query.Filter(elastic.NewTermQuery("mode", in.Mode))
	}
	if in.Pool >= 0 {
		query = query.Filter(elastic.NewTermQuery("pool", in.Pool))
	}
	if in.Progress >= 0 {
		query = query.Filter(elastic.NewTermQuery("progress", in.Progress))
	}
	if len(in.States) > 0 {
		states := make([]interface{}, len(in.States))
		for i, s := range in.States {
			states[i] = s
		}
		query = query.Filter(elastic.NewTermsQuery("state", states...))
	}
	if in.Type >= 0 {
		query = query.Filter(elastic.NewTermQuery("type", in.Type))
	}
	if len(in.AttrFormat) > 0 {
		attrs := make([]interface{}, len(in.AttrFormat))
		for i, a := range in.AttrFormat {
			attrs[i] = a
		}
		query = query.Filter(elastic.NewTermsQuery("attr_format", attrs...))
	}
	if in.CtimeFrom != "" {
		query = query.Filter(elastic.NewRangeQuery("ctime").Gte(in.CtimeFrom))
	}
	if in.CtimeTo != "" {
		query = query.Filter(elastic.NewRangeQuery("ctime").Lte(in.CtimeTo))
	}

	// 构建搜索参数
	params := &model.BasicSearchParams{
		Pn:     in.Bsp.Pn,
		Ps:     in.Bsp.Ps,
		Order:  in.Bsp.Order,
		Sort:   in.Bsp.Sort,
		Debug:  in.Bsp.Debug,
		KW:     in.Bsp.Kw,
		Source: []string{"id"},
	}

	// 执行搜索
	result, err := l.svcCtx.ESClient.Search(l.ctx, "dmExternal", indexName, query, params)
	if err != nil {
		l.Errorf("DmSearch failed: %v", err)
		return nil, err
	}

	// 转换结果
	resultBytes := make([][]byte, len(result.Result))
	for i, r := range result.Result {
		resultBytes[i] = r
	}

	return &search.SearchResp{
		Order:  result.Order,
		Sort:   result.Sort,
		Result: resultBytes,
		Page: &search.Page{
			Pn:    result.Page.Pn,
			Ps:    result.Page.Ps,
			Total: result.Page.Total,
		},
		Debug: result.Debug,
	}, nil
}

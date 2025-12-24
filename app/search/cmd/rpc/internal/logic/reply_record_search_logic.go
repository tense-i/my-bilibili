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

type ReplyRecordSearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyRecordSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyRecordSearchLogic {
	return &ReplyRecordSearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReplyRecordSearchLogic) ReplyRecordSearch(in *search.ReplyRecordSearchReq) (*search.SearchResp, error) {
	// mid 是必须的
	if in.Mid <= 0 {
		return nil, fmt.Errorf("mid is required")
	}

	// 构建索引名
	indexName := fmt.Sprintf("replyrecord_%02d", in.Mid%100)

	// 构建查询
	query := elastic.NewBoolQuery()
	query = query.Must(elastic.NewTermQuery("mid", in.Mid))

	// 过滤条件
	if len(in.Types) > 0 {
		types := toInterfaceSlice(in.Types)
		query = query.Must(elastic.NewTermsQuery("type", types...))
	}
	if len(in.States) > 0 {
		states := toInterfaceSlice(in.States)
		query = query.Must(elastic.NewTermsQuery("state", states...))
	}
	if in.CtimeFrom != "" {
		query = query.Must(elastic.NewRangeQuery("ctime").Gte(in.CtimeFrom))
	}
	if in.CtimeTo != "" {
		query = query.Must(elastic.NewRangeQuery("ctime").Lte(in.CtimeTo))
	}

	// 构建搜索参数
	params := &model.BasicSearchParams{
		Pn:    in.Bsp.Pn,
		Ps:    in.Bsp.Ps,
		Order: in.Bsp.Order,
		Sort:  in.Bsp.Sort,
		Debug: in.Bsp.Debug,
		KW:    in.Bsp.Kw,
	}

	// 执行搜索
	result, err := l.svcCtx.ESClient.Search(l.ctx, "replyExternal", indexName, query, params)
	if err != nil {
		l.Errorf("ReplyRecordSearch failed: %v", err)
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

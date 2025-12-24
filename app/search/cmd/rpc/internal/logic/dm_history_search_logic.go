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

type DmHistorySearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDmHistorySearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DmHistorySearchLogic {
	return &DmHistorySearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DmHistorySearch 弹幕历史搜索
// 参照 openbilibili/app/service/main/search/dao/dm_history.go
func (l *DmHistorySearchLogic) DmHistorySearch(in *search.DmHistorySearchReq) (*search.SearchResp, error) {
	// 构建索引名 - 使用 dm_search 索引（与主项目一致）
	indexName := fmt.Sprintf("dm_search_%03d", in.Oid%1000)

	// 构建查询
	query := elastic.NewBoolQuery()

	// 关键词搜索
	if in.Bsp != nil && in.Bsp.Kw != "" && len(in.Bsp.KwFields) > 0 {
		query = query.Must(elastic.NewMultiMatchQuery(in.Bsp.Kw, in.Bsp.KwFields...).Type("best_fields").TieBreaker(0.6))
	}

	// oid 过滤 - 使用 oidstr 字段（与主项目一致）
	if in.Oid > 0 {
		query = query.Filter(elastic.NewTermQuery("oidstr", in.Oid))
	}

	// 状态过滤
	if len(in.States) > 0 {
		states := make([]interface{}, len(in.States))
		for i, s := range in.States {
			states[i] = s
		}
		query = query.Filter(elastic.NewTermsQuery("state", states...))
	}

	// 时间范围过滤
	if in.CtimeFrom != "" {
		query = query.Filter(elastic.NewRangeQuery("ctime").Gte(in.CtimeFrom))
	}
	if in.CtimeTo != "" {
		query = query.Filter(elastic.NewRangeQuery("ctime").Lte(in.CtimeTo))
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
	result, err := l.svcCtx.ESClient.Search(l.ctx, "dmExternal", indexName, query, params)
	if err != nil {
		l.Errorf("DmHistorySearch failed: %v", err)
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

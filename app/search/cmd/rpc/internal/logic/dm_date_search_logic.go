package logic

import (
	"context"
	"strings"

	"mybilibili/app/search/cmd/rpc/internal/svc"
	"mybilibili/app/search/cmd/rpc/search"
	"mybilibili/app/search/model"

	"github.com/olivere/elastic/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

type DmDateSearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDmDateSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DmDateSearchLogic {
	return &DmDateSearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DmDateSearch 弹幕日期搜索
// 参照 openbilibili/app/service/main/search/dao/dm_date.go
func (l *DmDateSearchLogic) DmDateSearch(in *search.DmDateSearchReq) (*search.SearchResp, error) {
	// 构建索引名 - 格式: dm_date_2024_01
	indexName := "dm_date_" + strings.Replace(in.Month, "-", "_", -1)

	// 构建查询
	query := elastic.NewBoolQuery()

	// 关键词搜索
	if in.Bsp != nil && in.Bsp.Kw != "" && len(in.Bsp.KwFields) > 0 {
		query = query.Must(elastic.NewRegexpQuery(in.Bsp.KwFields[0], ".*"+in.Bsp.Kw+".*"))
	}

	// oid 过滤
	if in.Oid > 0 {
		query = query.Filter(elastic.NewTermQuery("oid", in.Oid))
	}

	// 月份过滤
	if in.Month != "" {
		query = query.Filter(elastic.NewTermQuery("month", in.Month))
	}

	// 月份范围过滤
	if in.MonthFrom != "" {
		query = query.Filter(elastic.NewRangeQuery("month").Gte(in.MonthFrom))
	}
	if in.MonthTo != "" {
		query = query.Filter(elastic.NewRangeQuery("month").Lte(in.MonthTo))
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
		l.Errorf("DmDateSearch failed: %v", err)
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

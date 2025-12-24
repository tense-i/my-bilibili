package logic

import (
	"context"

	"mybilibili/app/search/cmd/rpc/internal/svc"
	"mybilibili/app/search/cmd/rpc/search"
	"mybilibili/app/search/model"

	"github.com/olivere/elastic/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

type PgcMediaSearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPgcMediaSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PgcMediaSearchLogic {
	return &PgcMediaSearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PgcMediaSearchLogic) PgcMediaSearch(in *search.PgcMediaSearchReq) (*search.SearchResp, error) {
	query := elastic.NewBoolQuery()

	// 关键词搜索
	if in.Bsp != nil && in.Bsp.Kw != "" {
		query = query.Must(elastic.NewMultiMatchQuery(in.Bsp.Kw, "title").Type("best_fields").TieBreaker(0.3))
	}

	// 过滤条件
	if len(in.MediaIds) > 0 {
		ids := toInterfaceSlice(in.MediaIds)
		query = query.Filter(elastic.NewTermsQuery("media_id", ids...))
	}
	if len(in.SeasonIds) > 0 {
		ids := toInterfaceSlice(in.SeasonIds)
		query = query.Filter(elastic.NewTermsQuery("season_id", ids...))
	}
	if len(in.SeasonTypes) > 0 {
		types := toInterfaceSlice(in.SeasonTypes)
		query = query.Filter(elastic.NewTermsQuery("season_type", types...))
	}
	if len(in.StyleIds) > 0 {
		ids := toInterfaceSlice(in.StyleIds)
		query = query.Filter(elastic.NewTermsQuery("style_id", ids...))
	}
	if in.Status > -1000 {
		query = query.Filter(elastic.NewTermQuery("status", in.Status))
	}
	if in.ReleaseDateFrom != "" {
		query = query.Filter(elastic.NewRangeQuery("release_date").Gte(in.ReleaseDateFrom))
	}
	if in.ReleaseDateTo != "" {
		query = query.Filter(elastic.NewRangeQuery("release_date").Lte(in.ReleaseDateTo))
	}
	if in.SeasonIdFrom > 0 {
		query = query.Filter(elastic.NewRangeQuery("season_id").Gte(in.SeasonIdFrom))
	}
	if in.SeasonIdTo > 0 {
		query = query.Filter(elastic.NewRangeQuery("season_id").Lte(in.SeasonIdTo))
	}
	if len(in.ProducerIds) > 0 {
		ids := toInterfaceSlice(in.ProducerIds)
		query = query.Filter(elastic.NewTermsQuery("producer_id", ids...))
	}
	if in.IsDeleted == 0 {
		query = query.MustNot(elastic.NewTermQuery("is_deleted", 1))
	}
	if len(in.AreaIds) > 0 {
		areas := toStringInterfaceSlice(in.AreaIds)
		query = query.Filter(elastic.NewTermsQuery("area_id", areas...))
	}
	if in.ScoreFrom > 0 {
		query = query.Filter(elastic.NewRangeQuery("score_from").Gte(in.ScoreFrom))
	}
	if in.ScoreTo > 0 {
		query = query.Filter(elastic.NewRangeQuery("score_to").Lte(in.ScoreTo))
	}
	if in.IsFinish != "" {
		query = query.Filter(elastic.NewTermQuery("is_finish", in.IsFinish))
	}
	if len(in.SeasonVersions) > 0 {
		versions := toInterfaceSlice(in.SeasonVersions)
		query = query.Filter(elastic.NewTermsQuery("season_version", versions...))
	}
	if len(in.SeasonStatuses) > 0 {
		statuses := toInterfaceSlice(in.SeasonStatuses)
		query = query.Filter(elastic.NewTermsQuery("season_status", statuses...))
	}
	if in.PubTimeFrom != "" {
		query = query.Filter(elastic.NewRangeQuery("pub_time").Gte(in.PubTimeFrom))
	}
	if in.PubTimeTo != "" {
		query = query.Filter(elastic.NewRangeQuery("pub_time").Lte(in.PubTimeTo))
	}
	if len(in.SeasonMonths) > 0 {
		months := toInterfaceSlice(in.SeasonMonths)
		query = query.Filter(elastic.NewTermsQuery("season_month", months...))
	}
	if in.LatestTimeFrom != "" {
		query = query.Filter(elastic.NewRangeQuery("latest_time").Gte(in.LatestTimeFrom))
	}
	if in.LatestTimeTo != "" {
		query = query.Filter(elastic.NewRangeQuery("latest_time").Lte(in.LatestTimeTo))
	}
	if len(in.CopyrightInfos) > 0 {
		infos := toStringInterfaceSlice(in.CopyrightInfos)
		query = query.Filter(elastic.NewTermsQuery("copyright_info", infos...))
	}
	if in.CtimeFrom != "" {
		query = query.Filter(elastic.NewRangeQuery("ctime").Gte(in.CtimeFrom))
	}
	if in.CtimeTo != "" {
		query = query.Filter(elastic.NewRangeQuery("ctime").Lte(in.CtimeTo))
	}
	if in.MtimeFrom != "" {
		query = query.Filter(elastic.NewRangeQuery("mtime").Gte(in.MtimeFrom))
	}
	if in.MtimeTo != "" {
		query = query.Filter(elastic.NewRangeQuery("mtime").Lte(in.MtimeTo))
	}

	// 构建搜索参数
	params := &model.BasicSearchParams{
		Pn:     in.Bsp.Pn,
		Ps:     in.Bsp.Ps,
		Order:  in.Bsp.Order,
		Sort:   in.Bsp.Sort,
		Debug:  in.Bsp.Debug,
		KW:     in.Bsp.Kw,
		Source: []string{"media_id", "season_id", "season_type", "dm_count", "play_count", "fav_count", "score", "latest_time", "pub_time", "release_date"},
	}

	// 执行搜索
	result, err := l.svcCtx.ESClient.Search(l.ctx, "externalPublic", "pgc_media", query, params)
	if err != nil {
		l.Errorf("PgcMediaSearch failed: %v", err)
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

// toInterfaceSlice 将 int64 切片转换为 interface{} 切片
func toInterfaceSlice(slice []int64) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}

// toStringInterfaceSlice 将 string 切片转换为 interface{} 切片
func toStringInterfaceSlice(slice []string) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}

package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ VideoStatModel = (*customVideoStatModel)(nil)

type (
	// VideoStatModel is an interface to be customized, add more methods here,
	// and implement the added methods in customVideoStatModel.
	VideoStatModel interface {
		videoStatModel
		// 批量查询视频统计数据（hotrank-job 需要）
		FindByVids(ctx context.Context, vids []int64) ([]*VideoStat, error)
	}

	customVideoStatModel struct {
		*defaultVideoStatModel
	}
)

// NewVideoStatModel returns a model for the database table.
func NewVideoStatModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) VideoStatModel {
	return &customVideoStatModel{
		defaultVideoStatModel: newVideoStatModel(conn, c, opts...),
	}
}

// FindByVids 批量查询视频统计数据
func (m *customVideoStatModel) FindByVids(ctx context.Context, vids []int64) ([]*VideoStat, error) {
	if len(vids) == 0 {
		return []*VideoStat{}, nil
	}

	// 构建 SQL 的 IN 语句
	placeholders := make([]string, len(vids))
	args := make([]interface{}, len(vids))
	for i, vid := range vids {
		placeholders[i] = "?"
		args[i] = vid
	}

	query := fmt.Sprintf("select %s from %s where `vid` in (%s)",
		videoStatRows, m.table, strings.Join(placeholders, ","))

	var resp []*VideoStat
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

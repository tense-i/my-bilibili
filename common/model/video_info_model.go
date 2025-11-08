package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ VideoInfoModel = (*customVideoInfoModel)(nil)

type (
	// VideoInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customVideoInfoModel.
	VideoInfoModel interface {
		videoInfoModel
		// 批量查询视频信息（hotrank-job 需要）
		FindByVids(ctx context.Context, vids []int64) ([]*VideoInfo, error)
		// 游标分页查询
		FindListByLastVid(ctx context.Context, lastVid int64, limit int) ([]*VideoInfo, error)
	}

	customVideoInfoModel struct {
		*defaultVideoInfoModel
	}
)

// NewVideoInfoModel returns a model for the database table.
func NewVideoInfoModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) VideoInfoModel {
	return &customVideoInfoModel{
		defaultVideoInfoModel: newVideoInfoModel(conn, c, opts...),
	}
}

// FindByVids 批量查询视频信息
func (m *customVideoInfoModel) FindByVids(ctx context.Context, vids []int64) ([]*VideoInfo, error) {
	if len(vids) == 0 {
		return []*VideoInfo{}, nil
	}

	// 构建 SQL 的 IN 语句
	placeholders := make([]string, len(vids))
	args := make([]interface{}, len(vids))
	for i, vid := range vids {
		placeholders[i] = "?"
		args[i] = vid
	}

	query := fmt.Sprintf("select %s from %s where `vid` in (%s) and `state` = 0",
		videoInfoRows, m.table, strings.Join(placeholders, ","))

	var resp []*VideoInfo
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// FindListByLastVid 游标分页查询（参考主项目）
func (m *customVideoInfoModel) FindListByLastVid(ctx context.Context, lastVid int64, limit int) ([]*VideoInfo, error) {
	query := fmt.Sprintf("select %s from %s where `vid` > ? and `state` = 0 order by `vid` asc limit ?",
		videoInfoRows, m.table)

	var resp []*VideoInfo
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, lastVid, limit)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

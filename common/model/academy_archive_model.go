package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AcademyArchiveModel = (*customAcademyArchiveModel)(nil)

type (
	// AcademyArchiveModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAcademyArchiveModel.
	AcademyArchiveModel interface {
		academyArchiveModel
		// 全站热门排行榜
		FindHotRankList(ctx context.Context, business int, offset, limit int64) ([]*AcademyArchive, error)
		// 分区热门排行榜
		FindRegionHotRankList(ctx context.Context, regionId int64, business int, offset, limit int64) ([]*AcademyArchive, error)
		// 统计全站热门数量
		CountHotRank(ctx context.Context, business int) (int64, error)
		// 统计分区热门数量
		CountRegionHotRank(ctx context.Context, regionId int64, business int) (int64, error)
		// 根据OID查询热度值
		FindHotByOID(ctx context.Context, oid int64, business int) (int64, error)
	}

	customAcademyArchiveModel struct {
		*defaultAcademyArchiveModel
	}
)

// NewAcademyArchiveModel returns a model for the database table.
func NewAcademyArchiveModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AcademyArchiveModel {
	return &customAcademyArchiveModel{
		defaultAcademyArchiveModel: newAcademyArchiveModel(conn, c, opts...),
	}
}

// FindHotRankList 全站热门排行榜（按热度降序）
func (m *customAcademyArchiveModel) FindHotRankList(ctx context.Context, business int, offset, limit int64) ([]*AcademyArchive, error) {
	query := fmt.Sprintf("select %s from %s where `state` = 0 and `business` = ? order by `hot` desc, `id` desc limit ?,?",
		academyArchiveRows, m.table)

	var resp []*AcademyArchive
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, business, offset, limit)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// FindRegionHotRankList 分区热门排行榜（按热度降序）
func (m *customAcademyArchiveModel) FindRegionHotRankList(ctx context.Context, regionId int64, business int, offset, limit int64) ([]*AcademyArchive, error) {
	query := fmt.Sprintf("select %s from %s where `state` = 0 and `business` = ? and `region_id` = ? order by `hot` desc, `id` desc limit ?,?",
		academyArchiveRows, m.table)

	var resp []*AcademyArchive
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, business, regionId, offset, limit)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CountHotRank 统计全站热门数量
func (m *customAcademyArchiveModel) CountHotRank(ctx context.Context, business int) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where `state` = 0 and `business` = ?", m.table)

	var count int64
	err := m.QueryRowNoCacheCtx(ctx, &count, query, business)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CountRegionHotRank 统计分区热门数量
func (m *customAcademyArchiveModel) CountRegionHotRank(ctx context.Context, regionId int64, business int) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where `state` = 0 and `business` = ? and `region_id` = ?", m.table)

	var count int64
	err := m.QueryRowNoCacheCtx(ctx, &count, query, business, regionId)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// FindHotByOID 根据OID查询热度值
func (m *customAcademyArchiveModel) FindHotByOID(ctx context.Context, oid int64, business int) (int64, error) {
	query := fmt.Sprintf("select `hot` from %s where `oid` = ? and `business` = ? and `state` = 0", m.table)

	var hot int64
	err := m.QueryRowNoCacheCtx(ctx, &hot, query, oid, business)
	if err != nil {
		return 0, err
	}

	return hot, nil
}

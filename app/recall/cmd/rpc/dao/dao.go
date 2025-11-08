package dao

import (
	"context"
	"fmt"
	"strconv"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"mybilibili/app/recall/cmd/rpc/model"
)

type Dao struct {
	db    sqlx.SqlConn
	redis *redis.Redis
}

func New(db sqlx.SqlConn, rds *redis.Redis) *Dao {
	return &Dao{
		db:    db,
		redis: rds,
	}
}

// RecallFromRedis 从 Redis 召回视频
func (d *Dao) RecallFromRedis(ctx context.Context, key string, limit int32) ([]*model.RecallItem, error) {
	// 使用 ZREVRANGE 获取 Top N
	pairs, err := d.redis.ZrevrangeWithScoresCtx(ctx, key, 0, int64(limit-1))
	if err != nil {
		return nil, err
	}

	items := make([]*model.RecallItem, 0, len(pairs))
	for _, pair := range pairs {
		avidStr := pair.Key
		score := pair.Score

		avid, err := strconv.ParseInt(avidStr, 10, 64)
		if err != nil {
			continue
		}

		items = append(items, &model.RecallItem{
			AVID:  avid,
			Score: float64(score),
		})
	}

	return items, nil
}

// RecallFromRedisList 从 Redis List 召回
func (d *Dao) RecallFromRedisList(ctx context.Context, key string, limit int32) ([]*model.RecallItem, error) {
	avidsStr, err := d.redis.LrangeCtx(ctx, key, 0, int64(limit-1))
	if err != nil {
		return nil, err
	}

	items := make([]*model.RecallItem, 0, len(avidsStr))
	for _, avidStr := range avidsStr {
		avid, err := strconv.ParseInt(avidStr, 10, 64)
		if err != nil {
			continue
		}

		items = append(items, &model.RecallItem{
			AVID: avid,
		})
	}

	return items, nil
}

// GetVideoIndex 获取视频索引信息
func (d *Dao) GetVideoIndex(ctx context.Context, avids []int64) ([]*model.VideoIndex, error) {
	if len(avids) == 0 {
		return nil, nil
	}

	indexes := make([]*model.VideoIndex, 0, len(avids))

	for _, avid := range avids {
		index, err := d.getVideoIndexByAVID(ctx, avid)
		if err != nil || index == nil {
			continue
		}
		indexes = append(indexes, index)
	}

	return indexes, nil
}

// getVideoIndexByAVID 获取单个视频索引
func (d *Dao) getVideoIndexByAVID(ctx context.Context, avid int64) (*model.VideoIndex, error) {
	query := `
		SELECT avid, mid, title, zone_id, duration, pub_time, state
		FROM video_info
		WHERE avid = ?
	`

	index := &model.VideoIndex{}
	err := d.db.QueryRowCtx(ctx, query, avid).Scan(
		&index.AVID, &index.MID, &index.Title,
		&index.ZoneID, &index.Duration, &index.PubTime, &index.State,
	)

	if err != nil {
		return nil, err
	}

	// 获取标签
	tags, err := d.getVideoTags(ctx, avid)
	if err == nil {
		index.Tags = tags
	}

	return index, nil
}

// getVideoTags 获取视频标签
func (d *Dao) getVideoTags(ctx context.Context, avid int64) ([]model.Tag, error) {
	query := `SELECT tag_id, tag_name, tag_type FROM video_tag WHERE avid = ?`

	rows, err := d.db.QueryRowsCtx(ctx, query, avid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]model.Tag, 0)
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(&tag.TagID, &tag.TagName, &tag.TagType); err != nil {
			continue
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// CheckBloomFilter 检查布隆过滤器
func (d *Dao) CheckBloomFilter(ctx context.Context, mid int64, avid int64, date string) (bool, error) {
	key := fmt.Sprintf("bloomfilter:%d:%s", mid, date)
	return d.redis.SismemberCtx(ctx, key, avid)
}

// GetUserActionVideos 获取用户行为视频列表
func (d *Dao) GetUserActionVideos(ctx context.Context, mid int64, actionType string, date string, limit int) ([]int64, error) {
	key := fmt.Sprintf("user:action:%d:%s:%s", mid, actionType, date)

	pairs, err := d.redis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, 0, 9999999999, 0, int64(limit))
	if err != nil {
		return nil, err
	}

	avids := make([]int64, 0, len(pairs))
	for _, pair := range pairs {
		avid, err := strconv.ParseInt(pair.Key, 10, 64)
		if err != nil {
			continue
		}
		avids = append(avids, avid)
	}

	return avids, nil
}

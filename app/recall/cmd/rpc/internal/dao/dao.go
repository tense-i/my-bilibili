package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"mybilibili/app/recall/cmd/rpc/model"
)

type Dao struct {
	db    sqlx.SqlConn
	redis *redis.Redis
}

func NewDao(db sqlx.SqlConn, rds *redis.Redis) *Dao {
	return &Dao{
		db:    db,
		redis: rds,
	}
}

// ==================== 召回索引相关 ====================

// GetHotVideos 获取热门视频索引
func (d *Dao) GetHotVideos(ctx context.Context, limit int64, offset int64) ([]int64, error) {
	key := model.RedisKeyHotIndex
	vals, err := d.redis.ZrevrangeCtx(ctx, key, offset, offset+limit-1)
	if err != nil {
		return nil, err
	}

	avids := make([]int64, 0, len(vals))
	for _, val := range vals {
		if avid, err := strconv.ParseInt(val, 10, 64); err == nil {
			avids = append(avids, avid)
		}
	}

	return avids, nil
}

// GetSelectionVideos 获取精选视频索引
func (d *Dao) GetSelectionVideos(ctx context.Context, limit int64) ([]int64, error) {
	key := model.RedisKeySelectionIndex
	vals, err := d.redis.SmembersCtx(ctx, key)
	if err != nil {
		return nil, err
	}

	avids := make([]int64, 0, len(vals))
	for _, val := range vals {
		if avid, err := strconv.ParseInt(val, 10, 64); err == nil {
			avids = append(avids, avid)
			if len(avids) >= int(limit) {
				break
			}
		}
	}

	return avids, nil
}

// GetI2IVideos 获取 I2I 召回索引
func (d *Dao) GetI2IVideos(ctx context.Context, avid int64, limit int64) ([]int64, error) {
	key := fmt.Sprintf(model.RedisKeyI2IIndex, avid)
	vals, err := d.redis.ZrevrangeCtx(ctx, key, 0, limit-1)
	if err != nil {
		return nil, err
	}

	avids := make([]int64, 0, len(vals))
	for _, val := range vals {
		if recallAVID, err := strconv.ParseInt(val, 10, 64); err == nil {
			avids = append(avids, recallAVID)
		}
	}

	return avids, nil
}

// GetTagVideos 获取标签召回索引
func (d *Dao) GetTagVideos(ctx context.Context, tag string, limit int64) ([]int64, error) {
	key := fmt.Sprintf(model.RedisKeyTagIndex, tag)
	vals, err := d.redis.ZrevrangeCtx(ctx, key, 0, limit-1)
	if err != nil {
		return nil, err
	}

	avids := make([]int64, 0, len(vals))
	for _, val := range vals {
		if avid, err := strconv.ParseInt(val, 10, 64); err == nil {
			avids = append(avids, avid)
		}
	}

	return avids, nil
}

// GetUPVideos 获取 UP 主视频索引
func (d *Dao) GetUPVideos(ctx context.Context, upMid int64, limit int64) ([]int64, error) {
	key := fmt.Sprintf(model.RedisKeyUPIndex, upMid)
	vals, err := d.redis.ZrevrangeCtx(ctx, key, 0, limit-1)
	if err != nil {
		return nil, err
	}

	avids := make([]int64, 0, len(vals))
	for _, val := range vals {
		if avid, err := strconv.ParseInt(val, 10, 64); err == nil {
			avids = append(avids, avid)
		}
	}

	return avids, nil
}

// ==================== 用户行为相关 ====================

// GetUserRecentLikes 获取用户最近点赞的视频
func (d *Dao) GetUserRecentLikes(ctx context.Context, mid int64, limit int64) ([]int64, error) {
	query := `
		SELECT avid
		FROM user_behavior
		WHERE mid = ? AND action = 2 AND action_time > ?
		ORDER BY action_time DESC
		LIMIT ?
	`
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Unix()

	var avids []int64
	err := d.db.QueryRowsCtx(ctx, &avids, query, mid, thirtyDaysAgo, limit)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return avids, nil
}

// GetUserRecentPlays 获取用户最近播放的视频
func (d *Dao) GetUserRecentPlays(ctx context.Context, mid int64, limit int64) ([]int64, error) {
	query := `
		SELECT avid
		FROM user_behavior
		WHERE mid = ? AND action = 1 AND action_time > ?
		ORDER BY action_time DESC
		LIMIT ?
	`
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Unix()

	var avids []int64
	err := d.db.QueryRowsCtx(ctx, &avids, query, mid, thirtyDaysAgo, limit)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return avids, nil
}

// GetUserRecentPosVideos 获取用户最近正反馈的视频（观看时长>=15秒或完播率>=95%）
func (d *Dao) GetUserRecentPosVideos(ctx context.Context, mid int64, limit int64) ([]int64, error) {
	query := `
		SELECT avid
		FROM user_behavior
		WHERE mid = ? 
			AND behavior_type = 1 
			AND (duration >= 15 OR finish_rate >= 0.95)
			AND ctime > ?
		ORDER BY ctime DESC
		LIMIT ?
	`
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Unix()

	var avids []int64
	err := d.db.QueryRowsCtx(ctx, &avids, query, mid, sevenDaysAgo, limit)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return avids, nil
}

// GetUserRecentNegVideos 获取用户最近负反馈的视频（观看时长<5秒）
func (d *Dao) GetUserRecentNegVideos(ctx context.Context, mid int64, limit int64) ([]int64, error) {
	query := `
		SELECT avid
		FROM user_behavior
		WHERE mid = ? 
			AND behavior_type = 1 
			AND duration < 5
			AND ctime > ?
		ORDER BY ctime DESC
		LIMIT ?
	`
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Unix()

	var avids []int64
	err := d.db.QueryRowsCtx(ctx, &avids, query, mid, sevenDaysAgo, limit)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return avids, nil
}

// GetRandomVideos 获取随机视频（新发布的视频）
func (d *Dao) GetRandomVideos(ctx context.Context, tag string, limit int64) ([]int64, error) {
	// 先尝试从 Redis 获取新发布视频索引
	key := fmt.Sprintf(model.RedisKeyTagNewPubIndex, tag)
	vals, err := d.redis.ZrevrangeCtx(ctx, key, 0, limit-1)
	if err == nil && len(vals) > 0 {
		avids := make([]int64, 0, len(vals))
		for _, val := range vals {
			if avid, err := strconv.ParseInt(val, 10, 64); err == nil {
				avids = append(avids, avid)
			}
		}
		return avids, nil
	}

	// 如果 Redis 没有，从数据库查询最新视频
	query := `
		SELECT v.avid
		FROM video_info v
		INNER JOIN video_tag vt ON v.avid = vt.avid
		WHERE vt.tag_name = ? 
			AND v.state IN (1, 3, 4, 5)
			AND v.pub_time > ?
		ORDER BY v.pub_time DESC
		LIMIT ?
	`
	threeDaysAgo := time.Now().AddDate(0, 0, -3).Unix()

	var avids []int64
	err = d.db.QueryRowsCtx(ctx, &avids, query, tag, threeDaysAgo, limit)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return avids, nil
}

// GetUserFollowUPs 获取用户关注的 UP 主
func (d *Dao) GetUserFollowUPs(ctx context.Context, mid int64) ([]int64, error) {
	query := `
		SELECT up_mid
		FROM user_follow
		WHERE mid = ? AND status = 1
	`

	var upMids []int64
	err := d.db.QueryRowsCtx(ctx, &upMids, query, mid)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return upMids, nil
}

// ==================== 视频信息相关 ====================

// GetVideosBasicInfo 批量获取视频基础信息
func (d *Dao) GetVideosBasicInfo(ctx context.Context, avids []int64) ([]*model.VideoInfo, error) {
	if len(avids) == 0 {
		return []*model.VideoInfo{}, nil
	}

	// 构建 IN 查询
	placeholders := make([]string, len(avids))
	args := make([]interface{}, len(avids))
	for i, avid := range avids {
		placeholders[i] = "?"
		args[i] = avid
	}

	query := fmt.Sprintf(`
		SELECT 
			avid, mid, title, zone_id, duration, pub_time, state,
			play_hive, likes_hive, fav_hive, reply_hive, share_hive, coin_hive,
			play_month, likes_month, reply_month, share_month, play_month_finish
		FROM video_info
		WHERE avid IN (%s)
	`, strings.Join(placeholders, ","))

	var videos []*model.VideoInfo
	err := d.db.QueryRowsCtx(ctx, &videos, query, args...)
	if err != nil {
		return nil, err
	}

	// 批量获取视频标签
	tagsQuery := fmt.Sprintf(`
		SELECT avid, tag_id, tag_name
		FROM video_tag
		WHERE avid IN (%s)
	`, strings.Join(placeholders, ","))

	var tags []model.VideoTag
	err = d.db.QueryRowsCtx(ctx, &tags, tagsQuery, args...)
	if err != nil && err != sql.ErrNoRows {
		logx.Errorf("查询视频标签失败: %v", err)
	}

	// 将标签分组到对应的视频
	tagsMap := make(map[int64][]model.VideoTag)
	for _, tag := range tags {
		tagsMap[tag.AVID] = append(tagsMap[tag.AVID], tag)
	}

	for _, video := range videos {
		if videoTags, ok := tagsMap[video.AVID]; ok {
			video.Tags = videoTags
		}
	}

	return videos, nil
}

// ==================== Bloom Filter 相关 ====================

// CheckBloomFilter 检查 Bloom Filter
func (d *Dao) CheckBloomFilter(ctx context.Context, mid int64, avid int64) (bool, error) {
	key := fmt.Sprintf(model.RedisKeyBloomFilter, mid)
	member := strconv.FormatInt(avid, 10)

	// 使用 Redis Bitmap 实现简单的 Bloom Filter
	// 这里简化为使用 Set 来实现
	exists, err := d.redis.SismemberCtx(ctx, key, member)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// AddToBloomFilter 添加到 Bloom Filter
func (d *Dao) AddToBloomFilter(ctx context.Context, mid int64, avids []int64) error {
	if len(avids) == 0 {
		return nil
	}

	key := fmt.Sprintf(model.RedisKeyBloomFilter, mid)
	members := make([]interface{}, len(avids))
	for i, avid := range avids {
		members[i] = strconv.FormatInt(avid, 10)
	}

	_, err := d.redis.SaddCtx(ctx, key, members...)
	if err != nil {
		return err
	}

	// 设置过期时间（30天）
	d.redis.ExpireCtx(ctx, key, 30*24*3600)

	return nil
}

package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"mybilibili/app/recommend/cmd/rpc/model"
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

// ==================== 用户画像相关 ====================

// GetUserProfile 获取用户画像
func (d *Dao) GetUserProfile(ctx context.Context, mid int64) (*model.UserProfile, error) {
	// 1. 尝试从 Redis 获取
	key := fmt.Sprintf(model.RedisKeyUserProfile, mid)
	val, err := d.redis.GetCtx(ctx, key)
	if err == nil && val != "" {
		profile := &model.UserProfile{}
		if err := json.Unmarshal([]byte(val), profile); err == nil {
			return profile, nil
		}
	}

	// 2. 从数据库构建用户画像
	profile := &model.UserProfile{
		MID:       mid,
		Tags:      make(map[string]float64),
		Zones:     make(map[int32]float64),
		UPs:       make(map[int64]float64),
		UpdatedAt: time.Now().Unix(),
	}

	// 获取用户行为统计
	query := `
		SELECT 
			COUNT(*) as total_actions,
			SUM(CASE WHEN behavior_type = 1 THEN 1 ELSE 0 END) as play_count,
			SUM(CASE WHEN behavior_type = 2 THEN 1 ELSE 0 END) as like_count,
			SUM(CASE WHEN behavior_type = 3 THEN 1 ELSE 0 END) as coin_count,
			SUM(CASE WHEN behavior_type = 4 THEN 1 ELSE 0 END) as fav_count,
			SUM(CASE WHEN behavior_type = 5 THEN 1 ELSE 0 END) as share_count
		FROM user_behavior
		WHERE mid = ? AND ctime > ?
	`
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Unix()

	var stats struct {
		TotalActions int64
		PlayCount    int64
		LikeCount    int64
		CoinCount    int64
		FavCount     int64
		ShareCount   int64
	}

	err = d.db.QueryRowCtx(ctx, &stats, query, mid, thirtyDaysAgo)
	if err != nil && err != sql.ErrNoRows {
		logx.Errorf("查询用户行为统计失败: %v", err)
	}

	profile.TotalActions = stats.TotalActions
	profile.PlayCount = stats.PlayCount
	profile.LikeCount = stats.LikeCount
	profile.CoinCount = stats.CoinCount
	profile.FavCount = stats.FavCount
	profile.ShareCount = stats.ShareCount

	// 获取用户标签偏好（从点赞/收藏的视频标签统计）
	tagsQuery := `
		SELECT vt.tag_name, COUNT(*) as cnt
		FROM user_behavior ub
		INNER JOIN video_tag vt ON ub.avid = vt.avid
		WHERE ub.mid = ? AND ub.behavior_type IN (2,4) AND ub.ctime > ?
		GROUP BY vt.tag_name
		ORDER BY cnt DESC
		LIMIT 50
	`

	var tagStats []struct {
		TagName string
		Cnt     int64
	}

	err = d.db.QueryRowsCtx(ctx, &tagStats, tagsQuery, mid, thirtyDaysAgo)
	if err != nil && err != sql.ErrNoRows {
		logx.Errorf("查询用户标签偏好失败: %v", err)
	}

	maxTagCount := int64(1)
	for _, tag := range tagStats {
		if tag.Cnt > maxTagCount {
			maxTagCount = tag.Cnt
		}
	}

	for _, tag := range tagStats {
		profile.Tags[tag.TagName] = float64(tag.Cnt) / float64(maxTagCount)
	}

	// 获取用户分区偏好
	zonesQuery := `
		SELECT vi.zone_id, COUNT(*) as cnt
		FROM user_behavior ub
		INNER JOIN video_info vi ON ub.avid = vi.avid
		WHERE ub.mid = ? AND ub.behavior_type IN (1,2,4) AND ub.ctime > ?
		GROUP BY vi.zone_id
		ORDER BY cnt DESC
		LIMIT 20
	`

	var zoneStats []struct {
		ZoneID int32
		Cnt    int64
	}

	err = d.db.QueryRowsCtx(ctx, &zoneStats, zonesQuery, mid, thirtyDaysAgo)
	if err != nil && err != sql.ErrNoRows {
		logx.Errorf("查询用户分区偏好失败: %v", err)
	}

	maxZoneCount := int64(1)
	for _, zone := range zoneStats {
		if zone.Cnt > maxZoneCount {
			maxZoneCount = zone.Cnt
		}
	}

	for _, zone := range zoneStats {
		profile.Zones[zone.ZoneID] = float64(zone.Cnt) / float64(maxZoneCount)
	}

	// 获取用户关注的 UP 主
	upsQuery := `
		SELECT up_mid
		FROM user_follow
		WHERE mid = ? AND status = 1
	`

	var upMids []int64
	err = d.db.QueryRowsCtx(ctx, &upMids, upsQuery, mid)
	if err != nil && err != sql.ErrNoRows {
		logx.Errorf("查询用户关注UP主失败: %v", err)
	}

	for _, upMid := range upMids {
		profile.UPs[upMid] = 1.0
	}

	// 3. 缓存到 Redis（1小时过期）
	profileData, _ := json.Marshal(profile)
	d.redis.SetexCtx(ctx, key, string(profileData), 3600)

	return profile, nil
}

// ==================== 推荐记录相关 ====================

// SaveRecommendRecord 保存推荐记录
func (d *Dao) SaveRecommendRecord(ctx context.Context, mid int64, avids []int64) error {
	if len(avids) == 0 {
		return nil
	}

	// 保存到 Redis ZSet（用于去重）
	key := fmt.Sprintf(model.RedisKeyRecommendHistory, mid)
	now := time.Now().Unix()

	pairs := make([]redis.Pair, 0, len(avids))
	for _, avid := range avids {
		pairs = append(pairs, redis.Pair{
			Key:   strconv.FormatInt(avid, 10),
			Score: now,
		})
	}

	_, err := d.redis.ZaddsCtx(ctx, key, pairs...)
	if err != nil {
		logx.Errorf("保存推荐记录到Redis失败: %v", err)
		return err
	}

	// 设置7天过期
	d.redis.ExpireCtx(ctx, key, 7*24*3600)

	// 保持最近1000条记录
	d.redis.ZremrangebyrankCtx(ctx, key, 0, -1001)

	return nil
}

// GetRecommendHistory 获取推荐历史
func (d *Dao) GetRecommendHistory(ctx context.Context, mid int64, limit int64) ([]int64, error) {
	key := fmt.Sprintf(model.RedisKeyRecommendHistory, mid)

	// 获取最近的推荐记录
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

// ==================== 视频信息相关 ====================

// getInt64OrZero 从 sql.NullInt64 获取值，NULL 时返回 0
func getInt64OrZero(n sql.NullInt64) int64 {
	if n.Valid {
		return n.Int64
	}
	return 0
}

// videoInfoDB 数据库扫描专用结构体（使用 sql.NullXxx 处理可空字段）
type videoInfoDB struct {
	AVID            int64          `db:"avid"`
	CID             int64          `db:"cid"`
	UPMID           int64          `db:"mid"`
	Title           string         `db:"title"`
	Cover           sql.NullString `db:"cover"`
	Duration        int32          `db:"duration"`
	ZoneID          int32          `db:"zone_id"`
	PubTime         int64          `db:"pub_time"`
	State           int8           `db:"state"`
	PlayHive        sql.NullInt64  `db:"play_hive"`
	LikesHive       sql.NullInt64  `db:"likes_hive"`
	FavHive         sql.NullInt64  `db:"fav_hive"`
	ReplyHive       sql.NullInt64  `db:"reply_hive"`
	ShareHive       sql.NullInt64  `db:"share_hive"`
	CoinHive        sql.NullInt64  `db:"coin_hive"`
	PlayMonth       sql.NullInt64  `db:"play_month"`
	LikesMonth      sql.NullInt64  `db:"likes_month"`
	ReplyMonth      sql.NullInt64  `db:"reply_month"`
	ShareMonth      sql.NullInt64  `db:"share_month"`
	PlayMonthFinish sql.NullInt64  `db:"play_month_finish"`
}

// GetVideosInfo 批量获取视频信息
func (d *Dao) GetVideosInfo(ctx context.Context, avids []int64) (map[int64]*model.RecommendRecord, error) {
	if len(avids) == 0 {
		return make(map[int64]*model.RecommendRecord), nil
	}

	// 使用专用数据库结构体查询（使用 sql.NullXxx 处理可空字段）
	query := `
		SELECT 
			avid, cid, mid, title, cover, duration, zone_id, pub_time, state,
			play_hive, likes_hive, fav_hive, reply_hive, share_hive, coin_hive,
			play_month, likes_month, reply_month, share_month, play_month_finish
		FROM video_info
		WHERE avid = ?
	`

	result := make(map[int64]*model.RecommendRecord, len(avids))
	for _, avid := range avids {
		var dbVideo videoInfoDB
		err := d.db.QueryRowCtx(ctx, &dbVideo, query, avid)
		if err != nil && err != sql.ErrNoRows {
			logx.Errorf("查询视频信息失败: avid=%d, err=%v", avid, err)
			continue
		}
		if err == sql.ErrNoRows {
			continue
		}

		// 映射到 model.RecommendRecord（处理 NULL 值）
		cover := ""
		if dbVideo.Cover.Valid {
			cover = dbVideo.Cover.String
		}

		record := &model.RecommendRecord{
			AVID:            dbVideo.AVID,
			CID:             dbVideo.CID,
			UPMID:           dbVideo.UPMID,
			Title:           dbVideo.Title,
			Cover:           cover,
			Duration:        dbVideo.Duration,
			ZoneID:          dbVideo.ZoneID,
			PubTime:         dbVideo.PubTime,
			State:           dbVideo.State,
			PlayHive:        getInt64OrZero(dbVideo.PlayHive),
			LikesHive:       getInt64OrZero(dbVideo.LikesHive),
			FavHive:         getInt64OrZero(dbVideo.FavHive),
			ReplyHive:       getInt64OrZero(dbVideo.ReplyHive),
			ShareHive:       getInt64OrZero(dbVideo.ShareHive),
			CoinHive:        getInt64OrZero(dbVideo.CoinHive),
			PlayMonth:       getInt64OrZero(dbVideo.PlayMonth),
			LikesMonth:      getInt64OrZero(dbVideo.LikesMonth),
			ReplyMonth:      getInt64OrZero(dbVideo.ReplyMonth),
			ShareMonth:      getInt64OrZero(dbVideo.ShareMonth),
			PlayMonthFinish: getInt64OrZero(dbVideo.PlayMonthFinish),
			Tags:            make([]string, 0),
			TagIDs:          make([]int64, 0),
			Extra:           make(map[string]string),
		}
		result[record.AVID] = record
	}

	// 批量获取视频标签
	if len(result) > 0 {
		placeholders := make([]string, 0, len(avids))
		args := make([]interface{}, 0, len(avids))
		for avid := range result {
			placeholders = append(placeholders, "?")
			args = append(args, avid)
		}

		tagsQuery := fmt.Sprintf(`
			SELECT avid, tag_name
			FROM video_tag
			WHERE avid IN (%s)
		`, strings.Join(placeholders, ","))

		var tags []struct {
			AVID    int64  `db:"avid"`
			TagName string `db:"tag_name"`
		}

		err := d.db.QueryRowsCtx(ctx, &tags, tagsQuery, args...)
		if err != nil && err != sql.ErrNoRows {
			logx.Errorf("查询视频标签失败: %v", err)
		}

		for _, tag := range tags {
			if record, ok := result[tag.AVID]; ok {
				record.Tags = append(record.Tags, tag.TagName)
			}
		}
	}

	return result, nil
}

// ==================== 黑名单相关 ====================

// GetUserBlacklist 获取用户黑名单
func (d *Dao) GetUserBlacklist(ctx context.Context, mid int64) (map[int64]bool, error) {
	query := `
		SELECT up_mid
		FROM user_blacklist
		WHERE mid = ?
	`

	var targetIDs []int64
	err := d.db.QueryRowsCtx(ctx, &targetIDs, query, mid)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	blacklist := make(map[int64]bool, len(targetIDs))
	for _, id := range targetIDs {
		blacklist[id] = true
	}

	return blacklist, nil
}

// ==================== Bloom Filter 相关 ====================

// CheckBloomFilter 检查 Bloom Filter（检查用户是否看过该视频）
func (d *Dao) CheckBloomFilter(ctx context.Context, mid int64, avid int64) (bool, error) {
	// 使用推荐历史记录来检查
	history, err := d.GetRecommendHistory(ctx, mid, 1000)
	if err != nil {
		return false, err
	}

	// 检查 avid 是否在历史记录中
	for _, hAVID := range history {
		if hAVID == avid {
			return true, nil
		}
	}

	return false, nil
}

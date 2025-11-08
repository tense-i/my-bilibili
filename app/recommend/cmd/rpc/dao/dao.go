package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"mybilibili/app/recommend/cmd/rpc/model"
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

// LoadUserProfile 加载用户画像
func (d *Dao) LoadUserProfile(ctx context.Context, mid int64) (*model.UserProfile, error) {
	profile := &model.UserProfile{
		MID:         mid,
		Tags:        make(map[string]float64),
		Zones:       make(map[string]float64),
		PrefUps:     make(map[int64]int64),
		FollowUps:   make(map[int64]int64),
		BlackUps:    make(map[int64]bool),
		LikeVideos:  make(map[int64]int64),
		PosVideos:   make(map[int64]int64),
		NegVideos:   make(map[int64]int64),
		LikeTagIDs:  make(map[int64]int64),
		PosTagIDs:   make(map[int64]int64),
		NegTagIDs:   make(map[int64]int64),
		LikeUPs:     make(map[int64]int64),
		LastRecords: make([]int64, 0),
	}

	// 1. 从 Redis 加载用户画像缓存
	key := fmt.Sprintf("user:profile:%d", mid)
	exists, err := d.redis.ExistsCtx(ctx, key)
	if err == nil && exists {
		// TODO: 从 Redis 加载画像数据
		// tagsJson, _ := d.redis.HgetCtx(ctx, key, "tags")
		// json.Unmarshal([]byte(tagsJson), &profile.Tags)
	}

	// 2. 加载实时行为数据（点赞）
	date := time.Now().Format("20060102")
	likeKey := fmt.Sprintf("user:action:%d:like:%s", mid, date)
	profile.LikeVideos, _ = d.loadUserAction(ctx, likeKey)

	// 3. 加载正反馈数据
	posKey := fmt.Sprintf("user:action:%d:pos:%s", mid, date)
	profile.PosVideos, _ = d.loadUserAction(ctx, posKey)

	// 4. 加载负反馈数据
	negKey := fmt.Sprintf("user:action:%d:neg:%s", mid, date)
	profile.NegVideos, _ = d.loadUserAction(ctx, negKey)

	return profile, nil
}

// loadUserAction 加载用户行为数据
func (d *Dao) loadUserAction(ctx context.Context, key string) (map[int64]int64, error) {
	result := make(map[int64]int64)

	pairs, err := d.redis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, 0, time.Now().Unix(), 0, 100)
	if err != nil {
		return result, err
	}

	for _, pair := range pairs {
		avid := pair.Key
		timestamp := int64(pair.Score)
		avidInt, _ := fmt.Sscanf(avid, "%d", &avidInt)
		result[int64(avidInt)] = timestamp
	}

	return result, nil
}

// GetUserFollow 获取用户关注列表
func (d *Dao) GetUserFollow(ctx context.Context, mid int64, profile *model.UserProfile) error {
	query := `SELECT up_mid, ctime FROM user_follow WHERE mid = ? AND status = 1 ORDER BY ctime DESC LIMIT 100`

	rows, err := d.db.QueryRowsCtx(ctx, query, mid)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var upMID, ctime int64
		if err := rows.Scan(&upMID, &ctime); err != nil {
			continue
		}
		profile.FollowUps[upMID] = ctime
	}

	return nil
}

// GetUserBlack 获取用户黑名单
func (d *Dao) GetUserBlack(ctx context.Context, mid int64, profile *model.UserProfile) error {
	query := `SELECT up_mid FROM user_blacklist WHERE mid = ?`

	rows, err := d.db.QueryRowsCtx(ctx, query, mid)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var upMID int64
		if err := rows.Scan(&upMID); err != nil {
			continue
		}
		profile.BlackUps[upMID] = true
	}

	return nil
}

// GetVideoInfo 获取视频信息
func (d *Dao) GetVideoInfo(ctx context.Context, avid int64) (*model.RecommendRecord, error) {
	query := `
		SELECT avid, cid, mid, title, cover, duration, pub_time, zone_id, state,
			   play_hive, likes_hive, fav_hive, share_hive, reply_hive, coin_hive,
			   play_month, likes_month, share_month, reply_month, play_month_finish
		FROM video_info
		WHERE avid = ?
	`

	record := &model.RecommendRecord{}
	err := d.db.QueryRowCtx(ctx, query, avid).Scan(
		&record.AVID, &record.CID, &record.UPMID, &record.Title, &record.Cover,
		&record.Duration, &record.PubTime, &record.ZoneID, &record.State,
		&record.PlayHive, &record.LikesHive, &record.FavHive, &record.ShareHive,
		&record.ReplyHive, &record.CoinHive,
		&record.PlayMonth, &record.LikesMonth, &record.ShareMonth,
		&record.ReplyMonth, &record.PlayMonthFinish,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return record, err
}

// BatchGetVideoInfo 批量获取视频信息
func (d *Dao) BatchGetVideoInfo(ctx context.Context, avids []int64) ([]*model.RecommendRecord, error) {
	if len(avids) == 0 {
		return nil, nil
	}

	// 构建 IN 查询
	query := `
		SELECT avid, cid, mid, title, cover, duration, pub_time, zone_id, state,
			   play_hive, likes_hive, fav_hive, share_hive, reply_hive, coin_hive,
			   play_month, likes_month, share_month, reply_month, play_month_finish
		FROM video_info
		WHERE avid IN (?` + fmt.Sprintf(", ?") + `)`

	// TODO: 使用正确的批量查询方式
	records := make([]*model.RecommendRecord, 0, len(avids))

	for _, avid := range avids {
		record, err := d.GetVideoInfo(ctx, avid)
		if err != nil || record == nil {
			continue
		}
		records = append(records, record)
	}

	return records, nil
}

// GetVideoTags 获取视频标签
func (d *Dao) GetVideoTags(ctx context.Context, avid int64) ([]string, []int64, error) {
	query := `SELECT tag_id, tag_name FROM video_tag WHERE avid = ?`

	rows, err := d.db.QueryRowsCtx(ctx, query, avid)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	tags := make([]string, 0)
	tagIDs := make([]int64, 0)

	for rows.Next() {
		var tagID int64
		var tagName string
		if err := rows.Scan(&tagID, &tagName); err != nil {
			continue
		}
		tags = append(tags, tagName)
		tagIDs = append(tagIDs, tagID)
	}

	return tags, tagIDs, nil
}

// StoreRecommendResult 存储推荐结果（用于去重）
func (d *Dao) StoreRecommendResult(ctx context.Context, mid int64, avids []int64) error {
	if len(avids) == 0 {
		return nil
	}

	// 存储到布隆过滤器
	date := time.Now().Format("20060102")
	bfKey := fmt.Sprintf("bloomfilter:%d:%s", mid, date)

	for _, avid := range avids {
		// TODO: 使用实际的布隆过滤器库
		d.redis.SaddCtx(ctx, bfKey, avid)
	}

	// 设置30天过期
	d.redis.ExpireCtx(ctx, bfKey, 30*24*3600)

	return nil
}

// CheckBloomFilter 检查视频是否在布隆过滤器中
func (d *Dao) CheckBloomFilter(ctx context.Context, mid int64, avid int64) (bool, error) {
	date := time.Now().Format("20060102")
	bfKey := fmt.Sprintf("bloomfilter:%d:%s", mid, date)

	// TODO: 使用实际的布隆过滤器库
	return d.redis.SismemberCtx(ctx, bfKey, avid)
}

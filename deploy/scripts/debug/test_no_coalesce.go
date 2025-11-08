package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type videoInfoDB struct {
	AVID            int64           `db:"avid"`
	CID             int64           `db:"cid"`
	UPMID           int64           `db:"mid"`
	Title           string          `db:"title"`
	Cover           sql.NullString  `db:"cover"`
	Duration        int32           `db:"duration"`
	ZoneID          int32           `db:"zone_id"`
	PubTime         int64           `db:"pub_time"`
	State           int8            `db:"state"`
	PlayHive        sql.NullInt64   `db:"play_hive"`
	LikesHive       sql.NullInt64   `db:"likes_hive"`
	FavHive         sql.NullInt64   `db:"fav_hive"`
	ReplyHive       sql.NullInt64   `db:"reply_hive"`
	ShareHive       sql.NullInt64   `db:"share_hive"`
	CoinHive        sql.NullInt64   `db:"coin_hive"`
	PlayMonth       sql.NullInt64   `db:"play_month"`
	LikesMonth      sql.NullInt64   `db:"likes_month"`
	ReplyMonth      sql.NullInt64   `db:"reply_month"`
	ShareMonth      sql.NullInt64   `db:"share_month"`
	PlayMonthFinish sql.NullInt64   `db:"play_month_finish"`
}

func main() {
	dsn := "root:root123456@tcp(127.0.0.1:33060)/mybilibili?charset=utf8mb4&parseTime=true&loc=Local"
	db := sqlx.NewMysql(dsn)

	query := `
		SELECT 
			avid, cid, mid, title, cover, duration, zone_id, pub_time, state,
			play_hive, likes_hive, fav_hive, reply_hive, share_hive, coin_hive,
			play_month, likes_month, reply_month, share_month, play_month_finish
		FROM video_info
		WHERE avid = ?
	`

	fmt.Println("========== 测试不使用COALESCE，使用sql.NullInt64 ==========")
	
	var video videoInfoDB
	err := db.QueryRowCtx(context.Background(), &video, query, 100001)
	if err != nil {
		log.Fatalf("❌ 查询失败: %v\n", err)
	}

	fmt.Printf("✓ AVID=%d 查询成功:\n", video.AVID)
	fmt.Printf("  Title: %s\n", video.Title)
	fmt.Printf("  Cover.Valid: %v, Cover.String: %s\n", video.Cover.Valid, video.Cover.String)
	fmt.Printf("  PlayHive.Valid: %v, PlayHive.Int64: %d\n", video.PlayHive.Valid, video.PlayHive.Int64)
	fmt.Printf("  LikesHive.Valid: %v, LikesHive.Int64: %d\n", video.LikesHive.Valid, video.LikesHive.Int64)
	fmt.Printf("  PlayMonth.Valid: %v, PlayMonth.Int64: %d\n", video.PlayMonth.Valid, video.PlayMonth.Int64)
}

package main

import (
	"fmt"
	"log"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// 完整的测试结构体（20个数据库字段）
type FullTestRecord struct {
	AVID            int64  `db:"avid"`
	CID             int64  `db:"cid"`
	UPMID           int64  `db:"mid"`
	Title           string `db:"title"`
	Cover           string `db:"cover"`
	Duration        int32  `db:"duration"`
	ZoneID          int32  `db:"zone_id"`
	PubTime         int64  `db:"pub_time"`
	State           int8   `db:"state"`
	PlayHive        int64  `db:"play_hive"`
	LikesHive       int64  `db:"likes_hive"`
	FavHive         int64  `db:"fav_hive"`
	ReplyHive       int64  `db:"reply_hive"`
	ShareHive       int64  `db:"share_hive"`
	CoinHive        int64  `db:"coin_hive"`
	PlayMonth       int64  `db:"play_month"`
	LikesMonth      int64  `db:"likes_month"`
	ReplyMonth      int64  `db:"reply_month"`
	ShareMonth      int64  `db:"share_month"`
	PlayMonthFinish int64  `db:"play_month_finish"`
}

func main() {
	// 连接数据库
	dsn := "root:root123456@tcp(127.0.0.1:33060)/mybilibili?charset=utf8mb4&parseTime=true&loc=Local"
	db := sqlx.NewMysql(dsn)

	// 查询一个视频（完整20个字段）
	query := `
		SELECT 
			avid, cid, mid, title, COALESCE(cover, ''), duration, zone_id, pub_time, state,
			COALESCE(play_hive, 0), COALESCE(likes_hive, 0), COALESCE(fav_hive, 0), 
			COALESCE(reply_hive, 0), COALESCE(share_hive, 0), COALESCE(coin_hive, 0),
			COALESCE(play_month, 0), COALESCE(likes_month, 0), COALESCE(reply_month, 0), 
			COALESCE(share_month, 0), COALESCE(play_month_finish, 0)
		FROM video_info
		WHERE avid = 100001
	`

	fmt.Println("========== 测试完整20个字段 ==========")

	var video FullTestRecord
	err := db.QueryRow(&video, query)
	if err != nil {
		log.Fatalf("❌ 查询失败: %v", err)
	}

	fmt.Println("✓ 查询成功!")
	fmt.Printf("AVID: %d\n", video.AVID)
	fmt.Printf("CID: %d\n", video.CID)
	fmt.Printf("UPMID: %d\n", video.UPMID)
	fmt.Printf("Title: %s\n", video.Title)
	fmt.Printf("Cover: %s\n", video.Cover)
	fmt.Printf("Duration: %d\n", video.Duration)
	fmt.Printf("ZoneID: %d\n", video.ZoneID)
	fmt.Printf("PubTime: %d\n", video.PubTime)
	fmt.Printf("State: %d (int8)\n", video.State)
	fmt.Printf("PlayHive: %d\n", video.PlayHive)
	fmt.Printf("LikesHive: %d\n", video.LikesHive)
	fmt.Printf("PlayMonth: %d\n", video.PlayMonth)
	fmt.Printf("LikesMonth: %d\n", video.LikesMonth)
	fmt.Println("\n✓ 完整测试成功！")
}

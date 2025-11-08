package main

import (
	"fmt"
	"log"

	"mybilibili/app/recommend/cmd/rpc/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func main() {
	// 连接数据库
	dsn := "root:root123456@tcp(127.0.0.1:33060)/mybilibili?charset=utf8mb4&parseTime=true&loc=Local"
	db := sqlx.NewMysql(dsn)

	// 查询一个视频
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

	var video model.RecommendRecord
	err := db.QueryRow(&video, query)
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}

	fmt.Printf("查询成功!\n")
	fmt.Printf("AVID: %d\n", video.AVID)
	fmt.Printf("CID: %d\n", video.CID)
	fmt.Printf("Title: %s\n", video.Title)
	fmt.Printf("UPMID: %d\n", video.UPMID)
	fmt.Printf("Duration: %d\n", video.Duration)
	fmt.Printf("PlayHive: %d\n", video.PlayHive)
	fmt.Printf("PlayMonth: %d\n", video.PlayMonth)
}

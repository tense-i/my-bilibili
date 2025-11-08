package main

import (
	"context"
	"fmt"
	"log"

	"mybilibili/app/recommend/cmd/rpc/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func main() {
	// 连接数据库
	dsn := "root:root123456@tcp(127.0.0.1:33060)/mybilibili?charset=utf8mb4&parseTime=true&loc=Local"
	db := sqlx.NewMysql(dsn)

	// 批量查询（模拟服务中的查询）
	query := `
		SELECT 
			avid, cid, mid, title, COALESCE(cover, ''), duration, zone_id, pub_time, state,
			COALESCE(play_hive, 0), COALESCE(likes_hive, 0), COALESCE(fav_hive, 0), 
			COALESCE(reply_hive, 0), COALESCE(share_hive, 0), COALESCE(coin_hive, 0),
			COALESCE(play_month, 0), COALESCE(likes_month, 0), COALESCE(reply_month, 0), 
			COALESCE(share_month, 0), COALESCE(play_month_finish, 0)
		FROM video_info
		WHERE avid IN (100001, 100002, 100003)
	`

	fmt.Println("========== 测试批量查询 (QueryRowsCtx) ==========")

	var videos []model.RecommendRecord
	err := db.QueryRowsCtx(context.Background(), &videos, query)
	if err != nil {
		log.Fatalf("❌ 批量查询失败: %v", err)
	}

	fmt.Printf("✓ 批量查询成功! 查询到 %d 条记录\n\n", len(videos))

	for i, video := range videos {
		fmt.Printf("记录 %d:\n", i+1)
		fmt.Printf("  AVID: %d\n", video.AVID)
		fmt.Printf("  CID: %d\n", video.CID)
		fmt.Printf("  UPMID: %d\n", video.UPMID)
		fmt.Printf("  Title: %s\n", video.Title)
		fmt.Printf("  State: %d (int8)\n", video.State)
		fmt.Printf("  PlayHive: %d\n", video.PlayHive)
		fmt.Printf("  LikesHive: %d\n", video.LikesHive)
		fmt.Printf("  PlayMonth: %d\n\n", video.PlayMonth)
	}

	fmt.Println("✓ 测试成功！")
}

package main

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"mybilibili/app/recommend/cmd/rpc/model"
)

func main() {
	dsn := "root:root123456@tcp(127.0.0.1:33060)/mybilibili?charset=utf8mb4&parseTime=true&loc=Local"
	db := sqlx.NewMysql(dsn)

	query := `
		SELECT 
			avid, cid, mid, title, COALESCE(cover, ''), duration, zone_id, pub_time, state,
			COALESCE(play_hive, 0), COALESCE(likes_hive, 0), COALESCE(fav_hive, 0), 
			COALESCE(reply_hive, 0), COALESCE(share_hive, 0), COALESCE(coin_hive, 0),
			COALESCE(play_month, 0), COALESCE(likes_month, 0), COALESCE(reply_month, 0), 
			COALESCE(share_month, 0), COALESCE(play_month_finish, 0)
		FROM video_info
		WHERE avid = ?
	`

	fmt.Println("========== 测试 QueryRow ===========")
	var video1 model.RecommendRecord
	err := db.QueryRow(&video1, query, 100001)
	if err != nil {
		fmt.Printf("❌ QueryRow 失败: %v\n", err)
	} else {
		fmt.Printf("✓ QueryRow 成功: AVID=%d, Title=%s\n", video1.AVID, video1.Title)
	}

	fmt.Println("\n========== 测试 QueryRowCtx ===========")
	var video2 model.RecommendRecord
	err = db.QueryRowCtx(context.Background(), &video2, query, 100001)
	if err != nil {
		fmt.Printf("❌ QueryRowCtx 失败: %v\n", err)
	} else {
		fmt.Printf("✓ QueryRowCtx 成功: AVID=%d, Title=%s\n", video2.AVID, video2.Title)
	}
}

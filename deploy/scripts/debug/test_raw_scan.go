package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 连接数据库
	dsn := "root:root123456@tcp(127.0.0.1:33060)/mybilibili?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer db.Close()

	// 查询一个视频（使用COALESCE处理NULL）
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

	var (
		avid, cid, mid, pubTime                                        int64
		playHive, likesHive, favHive, replyHive, shareHive, coinHive   int64
		playMonth, likesMonth, replyMonth, shareMonth, playMonthFinish int64
		title, cover                                                   string
		duration, zoneID                                               int32
		state                                                          int8
	)

	err = db.QueryRow(query).Scan(
		&avid, &cid, &mid, &title, &cover, &duration, &zoneID, &pubTime, &state,
		&playHive, &likesHive, &favHive, &replyHive, &shareHive, &coinHive,
		&playMonth, &likesMonth, &replyMonth, &shareMonth, &playMonthFinish,
	)

	if err != nil {
		log.Fatalf("扫描数据失败: %v", err)
	}

	fmt.Println("========== 原生SQL扫描成功 ==========")
	fmt.Printf("AVID: %d\n", avid)
	fmt.Printf("CID: %d\n", cid)
	fmt.Printf("MID: %d\n", mid)
	fmt.Printf("Title: %s\n", title)
	fmt.Printf("Cover: %s\n", cover)
	fmt.Printf("Duration: %d\n", duration)
	fmt.Printf("ZoneID: %d\n", zoneID)
	fmt.Printf("PubTime: %d\n", pubTime)
	fmt.Printf("State: %d (int8)\n", state)
	fmt.Printf("PlayHive: %d\n", playHive)
	fmt.Printf("PlayMonth: %d\n", playMonth)
	fmt.Println("\n✓ 原生SQL测试成功！")
	fmt.Println("========================================")
	fmt.Println("下一步：测试 sqlx 结构体扫描")
}

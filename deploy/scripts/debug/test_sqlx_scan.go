package main

import (
	"fmt"
	"log"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// 简化的测试结构体
type TestRecord struct {
	AVID      int64  `db:"avid"`
	CID       int64  `db:"cid"`
	UPMID     int64  `db:"mid"`
	Title     string `db:"title"`
	Cover     string `db:"cover"`
	Duration  int32  `db:"duration"`
	ZoneID    int32  `db:"zone_id"`
	PubTime   int64  `db:"pub_time"`
	State     int8   `db:"state"`
	PlayHive  int64  `db:"play_hive"`
	PlayMonth int64  `db:"play_month"`
}

func main() {
	// 连接数据库
	dsn := "root:root123456@tcp(127.0.0.1:33060)/mybilibili?charset=utf8mb4&parseTime=true&loc=Local"
	db := sqlx.NewMysql(dsn)

	// 查询一个视频
	query := `
		SELECT 
			avid, cid, mid, title, COALESCE(cover, ''), duration, zone_id, pub_time, state,
			COALESCE(play_hive, 0), COALESCE(play_month, 0)
		FROM video_info
		WHERE avid = 100001
	`

	fmt.Println("========== 测试 go-zero sqlx 结构体扫描 ==========")
	fmt.Println("查询SQL:")
	fmt.Println(query)
	fmt.Println()

	var video TestRecord
	err := db.QueryRow(&video, query)
	if err != nil {
		log.Fatalf("查询失败: %v\n\n建议：go-zero的sqlx可能不支持db标签映射，需要字段顺序完全匹配", err)
	}

	fmt.Println("✓ 查询成功!")
	fmt.Printf("AVID: %d\n", video.AVID)
	fmt.Printf("CID: %d\n", video.CID)
	fmt.Printf("UPMID: %d\n", video.UPMID)
	fmt.Printf("Title: %s\n", video.Title)
	fmt.Printf("State: %d (int8)\n", video.State)
	fmt.Printf("PlayHive: %d\n", video.PlayHive)
}

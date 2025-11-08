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

	// 查询一个视频
	query := `
		SELECT 
			avid, cid, mid, title, cover, duration, zone_id, pub_time, state,
			play_hive, likes_hive, fav_hive, reply_hive, share_hive, coin_hive,
			play_month, likes_month, reply_month, share_month, play_month_finish
		FROM video_info
		WHERE avid = 100001
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}
	defer rows.Close()

	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		log.Fatalf("获取列信息失败: %v", err)
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		log.Fatalf("获取列类型失败: %v", err)
	}

	fmt.Println("========== 数据库字段信息 ==========")
	fmt.Printf("字段数量: %d\n\n", len(columns))

	for i, col := range columns {
		colType := columnTypes[i]
		scanType := colType.ScanType()
		dbType := colType.DatabaseTypeName()
		nullable, _ := colType.Nullable()
		length, hasLength := colType.Length()
		precision, scale, hasPrecision := colType.DecimalSize()

		fmt.Printf("%d. %s\n", i+1, col)
		fmt.Printf("   数据库类型: %s\n", dbType)
		fmt.Printf("   Go扫描类型: %v\n", scanType)
		fmt.Printf("   可为空: %v\n", nullable)
		if hasLength {
			fmt.Printf("   长度: %d\n", length)
		}
		if hasPrecision {
			fmt.Printf("   精度: %d, 小数位: %d\n", precision, scale)
		}
		fmt.Println()
	}

	// 尝试扫描数据
	if rows.Next() {
		var (
			avid, cid, mid, pubTime                                        int64
			playHive, likesHive, favHive, replyHive, shareHive, coinHive   int64
			playMonth, likesMonth, replyMonth, shareMonth, playMonthFinish int64
			title, cover                                                   string
			duration, zoneID, state                                        int32
		)

		err := rows.Scan(
			&avid, &cid, &mid, &title, &cover, &duration, &zoneID, &pubTime, &state,
			&playHive, &likesHive, &favHive, &replyHive, &shareHive, &coinHive,
			&playMonth, &likesMonth, &replyMonth, &shareMonth, &playMonthFinish,
		)

		if err != nil {
			log.Fatalf("扫描数据失败: %v", err)
		}

		fmt.Println("========== 扫描结果 ==========")
		fmt.Printf("AVID: %d\n", avid)
		fmt.Printf("CID: %d\n", cid)
		fmt.Printf("MID: %d\n", mid)
		fmt.Printf("Title: %s\n", title)
		fmt.Printf("Cover: %s\n", cover)
		fmt.Printf("Duration: %d\n", duration)
		fmt.Printf("ZoneID: %d\n", zoneID)
		fmt.Printf("PubTime: %d\n", pubTime)
		fmt.Printf("State: %d\n", state)
		fmt.Printf("PlayHive: %d\n", playHive)
		fmt.Printf("PlayMonth: %d\n", playMonth)
		fmt.Println("\n✓ 数据扫描成功！")
	}
}

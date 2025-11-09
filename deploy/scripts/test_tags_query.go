package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func main() {
	dsn := "root:root123456@tcp(127.0.0.1:33060)/mybilibili?charset=utf8mb4&parseTime=true&loc=Local"
	db := sqlx.NewMysql(dsn)

	avids := []int64{100005, 100007}
	placeholders := make([]string, 0, len(avids))
	args := make([]interface{}, 0, len(avids))
	for _, avid := range avids {
		placeholders = append(placeholders, "?")
		args = append(args, avid)
	}

	tagsQuery := fmt.Sprintf(`
		SELECT avid, tag_name
		FROM video_tag
		WHERE avid IN (%s)
	`, strings.Join(placeholders, ","))

	fmt.Printf("Query: %s\n", tagsQuery)
	fmt.Printf("Args: %v\n\n", args)

	var tags []struct {
		AVID    int64  `db:"avid"`
		TagName string `db:"tag_name"`
	}

	err := db.QueryRowsCtx(context.Background(), &tags, tagsQuery, args...)
	if err != nil {
		fmt.Printf("❌ 查询失败: %v\n", err)
		if err == sql.ErrNoRows {
			fmt.Println("(没有找到记录)")
		}
		return
	}

	fmt.Printf("✓ 查询成功，找到 %d 条标签:\n", len(tags))
	for _, tag := range tags {
		fmt.Printf("  AVID=%d, TagName=%s\n", tag.AVID, tag.TagName)
	}
}

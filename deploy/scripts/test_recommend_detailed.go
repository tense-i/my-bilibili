package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"mybilibili/app/recommend/cmd/rpc/recommend"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接到推荐服务
	conn, err := grpc.Dial("127.0.0.1:9005", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	client := recommend.NewRecommendClient(conn)

	// 测试用户
	testUsers := []int64{1001, 1002}

	fmt.Println("========================================")
	fmt.Println("详细测试推荐服务")
	fmt.Println("========================================\n")

	for _, mid := range testUsers {
		fmt.Printf("=== 测试用户 %d ===\n", mid)

		// 构建推荐请求
		req := &recommend.RecommendRequest{
			Mid:   mid,
			Limit: 3,
			Debug: true,
		}

		// 调用推荐服务
		resp, err := client.GetRecommendList(context.Background(), req)
		if err != nil {
			fmt.Printf("  错误: %v\n\n", err)
			continue
		}

		// 输出结果（JSON格式，避免编码问题）
		fmt.Printf("  推荐数量: %d\n", len(resp.List))
		fmt.Printf("  还有更多: %v\n", resp.HasMore)

		if len(resp.List) > 0 {
			fmt.Println("\n  推荐列表（JSON格式）:")
			for i, item := range resp.List {
				itemJSON, _ := json.MarshalIndent(item, "    ", "  ")
				fmt.Printf("    %d. %s\n\n", i+1, string(itemJSON))
			}
		}

		// 验证字段
		fmt.Println("  字段验证:")
		for i, item := range resp.List {
			fmt.Printf("    %d. AVID=%d\n", i+1, item.Avid)
			fmt.Printf("       Title长度=%d, Title非空=%v\n", len(item.Title), item.Title != "")
			fmt.Printf("       UPName=%s, UPMID=%d\n", item.UpName, item.UpMid)
			fmt.Printf("       ZoneName=%s, ZoneID=%d\n", item.ZoneName, item.ZoneId)
			fmt.Printf("       Tags数量=%d, Tags=%v\n", len(item.Tags), item.Tags)
			fmt.Printf("       Score=%.4f, Reason=%s\n", item.Score, item.Reason)
			fmt.Printf("       Play=%d, Like=%d\n", item.Play, item.Like)
		}

		// 输出调试信息
		if len(resp.DebugInfo) > 0 {
			fmt.Println("\n  调试信息:")
			debugJSON, _ := json.MarshalIndent(resp.DebugInfo, "    ", "  ")
			fmt.Printf("    %s\n", string(debugJSON))
		}

		fmt.Println()
	}

	fmt.Println("========================================")
	fmt.Println("测试完成！")
	fmt.Println("========================================")
}


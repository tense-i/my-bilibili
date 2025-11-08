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
	testUsers := []int64{1001, 1002, 1003}

	fmt.Println("========================================")
	fmt.Println("测试推荐服务 (Recommend RPC)")
	fmt.Println("========================================\n")

	for _, mid := range testUsers {
		fmt.Printf("=== 测试用户 %d ===\n", mid)

		// 构建推荐请求
		req := &recommend.RecommendRequest{
			Mid:   mid,
			Limit: 5,
			Debug: true,
		}

		// 调用推荐服务
		resp, err := client.GetRecommendList(context.Background(), req)
		if err != nil {
			fmt.Printf("  错误: %v\n\n", err)
			continue
		}

		// 输出结果
		fmt.Printf("  推荐数量: %d\n", len(resp.List))
		fmt.Printf("  还有更多: %v\n", resp.HasMore)

		if len(resp.List) > 0 {
			fmt.Println("  推荐列表:")
			for i, item := range resp.List {
				fmt.Printf("    %d. AVID=%d, 标题=%s\n", i+1, item.Avid, item.Title)
				fmt.Printf("       UP主=%s, 播放=%d, 点赞=%d\n", item.UpName, item.Play, item.Like)
				fmt.Printf("       分数=%.2f, 推荐理由=%s\n", item.Score, item.Reason)
				fmt.Printf("       标签=%v\n", item.Tags)
				fmt.Println()
			}
		}

		// 输出调试信息
		if len(resp.DebugInfo) > 0 {
			fmt.Println("  调试信息:")
			debugJSON, _ := json.MarshalIndent(resp.DebugInfo, "    ", "  ")
			fmt.Printf("    %s\n", string(debugJSON))
		}

		fmt.Println()
	}

	fmt.Println("========================================")
	fmt.Println("测试完成！")
	fmt.Println("========================================")
}

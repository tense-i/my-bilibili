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
	req := &recommend.RecommendRequest{
		Mid:   1001,
		Limit: 10,
		Debug: true,
	}

	// 调用推荐服务
	resp, err := client.GetRecommendList(context.Background(), req)
	if err != nil {
		log.Fatalf("调用失败: %v", err)
	}

	fmt.Println("========================================")
	fmt.Println("特征值和分数分析")
	fmt.Println("========================================\n")

	// 按分数分组
	scoreGroups := make(map[float64][]*recommend.RecommendItem)
	for _, record := range resp.List {
		// 将分数四舍五入到4位小数，便于分组
		roundedScore := float64(int(record.Score*10000)) / 10000
		scoreGroups[roundedScore] = append(scoreGroups[roundedScore], record)
	}

	fmt.Printf("不同分数组数量: %d\n\n", len(scoreGroups))

	for score, records := range scoreGroups {
		fmt.Printf("=== 分数组: %.4f (数量: %d) ===\n", score, len(records))
		for i, record := range records {
			if i >= 3 { // 只显示前3个
				fmt.Printf("  ... (还有 %d 个)\n", len(records)-3)
				break
			}
			fmt.Printf("  视频 %d:\n", record.Avid)
			fmt.Printf("    召回类型: %s\n", record.Reason)
			fmt.Printf("    分区ID: %d\n", record.ZoneId)
			fmt.Printf("    播放量: %d\n", record.Play)
			fmt.Printf("    点赞数: %d\n", record.Like)
			fmt.Printf("    标签数: %d\n", len(record.Tags))

			// 输出JSON以便查看详细信息
			jsonData, _ := json.MarshalIndent(record, "    ", "  ")
			fmt.Printf("    详细信息:\n%s\n", string(jsonData))
		}
		fmt.Println()
	}
}

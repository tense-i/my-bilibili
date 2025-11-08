package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"mybilibili/app/recall/cmd/rpc/recall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接到召回服务
	conn, err := grpc.Dial("127.0.0.1:9006", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	client := recall.NewRecallClient(conn)

	// 测试用户
	testUsers := []int64{1001, 1002, 1003}

	fmt.Println("========================================")
	fmt.Println("测试召回服务 (Recall RPC)")
	fmt.Println("========================================\n")

	for _, mid := range testUsers {
		fmt.Printf("=== 测试用户 %d ===\n", mid)

		// 构建召回请求
		req := &recall.RecallRequest{
			Mid:        mid,
			TotalLimit: 20,
			Infos: []*recall.RecallInfo{
				{
					Name:     "HotRecall",
					Tag:      "recall:hot:default",
					Limit:    10,
					Priority: 10,
					Filter:   "bloomfilter",
				},
			},
		}

		// 调用召回服务
		resp, err := client.Recall(context.Background(), req)
		if err != nil {
			fmt.Printf("  错误: %v\n\n", err)
			continue
		}

		// 输出结果
		fmt.Printf("  召回数量: %d\n", len(resp.List))
		if len(resp.List) > 0 {
			fmt.Println("  前5个结果:")
			for i, item := range resp.List {
				if i >= 5 {
					break
				}
				fmt.Printf("    %d. AVID=%d, 分数=%.2f, 类型=%s\n",
					i+1, item.Avid, item.Score, item.RecallType)
			}
		}

		// 输出 JSON（方便调试）
		if len(resp.List) > 0 {
			jsonData, _ := json.MarshalIndent(resp.List[0], "  ", "  ")
			fmt.Printf("  第一个结果详情:\n  %s\n", string(jsonData))
		}

		fmt.Println()
	}

	fmt.Println("========================================")
	fmt.Println("测试完成！")
	fmt.Println("========================================")
}

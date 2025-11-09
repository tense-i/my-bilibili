package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	modelPath := "../../app/recommend/cmd/rpc/models/0.0.13/model.json"

	data, err := ioutil.ReadFile(modelPath)
	if err != nil {
		fmt.Printf("❌ 读取文件失败: %v\n", err)
		os.Exit(1)
	}

	var xgbModel struct {
		Learner struct {
			GradientBooster struct {
				Model struct {
					Trees []struct {
						TreeParam struct {
							NumNodes string `json:"num_nodes"` // 注意：XGBoost JSON 中 num_nodes 是字符串
						} `json:"tree_param"`
						SplitIndices    []int     `json:"split_indices"`
						SplitConditions []float64 `json:"split_conditions"`
						DefaultLeft     []int     `json:"default_left"`
						LeftChildren    []int     `json:"left_children"`
						RightChildren   []int     `json:"right_children"`
						LeafValues      []float64 `json:"base_weights"`
					} `json:"trees"`
				} `json:"model"`
			} `json:"gradient_booster"`
		} `json:"learner"`
	}

	if err := json.Unmarshal(data, &xgbModel); err != nil {
		fmt.Printf("❌ 解析 JSON 失败: %v\n", err)
		os.Exit(1)
	}

	trees := xgbModel.Learner.GradientBooster.Model.Trees
	fmt.Printf("✓ 成功解析，共 %d 棵树\n\n", len(trees))

	if len(trees) > 0 {
		tree := trees[0]
		fmt.Printf("第一棵树信息:\n")
		fmt.Printf("  num_nodes: %s\n", tree.TreeParam.NumNodes)
		fmt.Printf("  left_children length: %d\n", len(tree.LeftChildren))
		fmt.Printf("  right_children length: %d\n", len(tree.RightChildren))
		fmt.Printf("  split_indices length: %d\n", len(tree.SplitIndices))
		fmt.Printf("  base_weights length: %d\n", len(tree.LeafValues))

		if len(tree.LeftChildren) > 0 {
			fmt.Printf("\n前5个节点:\n")
			for i := 0; i < 5 && i < len(tree.LeftChildren); i++ {
				isLeaf := tree.LeftChildren[i] == -1 && tree.RightChildren[i] == -1
				fmt.Printf("  节点 %d: left=%d, right=%d, is_leaf=%v",
					i, tree.LeftChildren[i], tree.RightChildren[i], isLeaf)
				if isLeaf {
					fmt.Printf(", leaf_value=%.6f", tree.LeafValues[i])
				} else {
					fmt.Printf(", feature=%d, threshold=%.6f",
						tree.SplitIndices[i], tree.SplitConditions[i])
				}
				fmt.Println()
			}
		}
	}
}

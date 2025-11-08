package rank

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/core/logx"
)

// XGBoostModel XGBoost 模型封装
type XGBoostModel struct {
	ModelPath    string
	ModelVersion string
	NumFeatures  int
	NumTrees     int
	Trees        []*DecisionTree
}

// DecisionTree 决策树结构
type DecisionTree struct {
	Nodes []TreeNode `json:"nodes"`
}

// TreeNode 树节点
type TreeNode struct {
	NodeID       int     `json:"node_id"`
	FeatureIndex int     `json:"feature_index"` // -1 表示叶子节点
	Threshold    float64 `json:"threshold"`
	LeftChild    int     `json:"left_child"`
	RightChild   int     `json:"right_child"`
	LeafValue    float64 `json:"leaf_value"`
	MissingGoTo  int     `json:"missing_go_to"` // 缺失值处理
}

// LoadXGBoostModel 加载 XGBoost 模型
func LoadXGBoostModel(modelDir string) (*XGBoostModel, error) {
	// 1. 读取模型配置
	configPath := filepath.Join(modelDir, "config.json")
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取模型配置失败: %v", err)
	}

	var config struct {
		Version     string `json:"version"`
		NumFeatures int    `json:"num_features"`
		NumTrees    int    `json:"num_trees"`
		Format      string `json:"format"`
	}

	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("解析模型配置失败: %v", err)
	}

	logx.Infof("加载模型配置: version=%s, num_features=%d, num_trees=%d",
		config.Version, config.NumFeatures, config.NumTrees)

	// 2. 加载模型文件
	model := &XGBoostModel{
		ModelPath:    modelDir,
		ModelVersion: config.Version,
		NumFeatures:  config.NumFeatures,
		NumTrees:     config.NumTrees,
		Trees:        make([]*DecisionTree, 0),
	}

	// 尝试加载 JSON 格式模型
	modelJSONPath := filepath.Join(modelDir, "model.json")
	if _, err := os.Stat(modelJSONPath); err == nil {
		if err := model.loadFromJSON(modelJSONPath); err != nil {
			logx.Errorf("从 JSON 加载模型失败: %v", err)
		} else {
			logx.Info("成功从 JSON 加载模型")
			return model, nil
		}
	}

	// 如果 JSON 加载失败，尝试加载 TreeLite Protobuf 格式
	protoPath := filepath.Join(modelDir, "model.proto")
	if _, err := os.Stat(protoPath); err == nil {
		if err := model.loadFromTreeLiteProto(protoPath); err != nil {
			logx.Errorf("从 TreeLite Protobuf 加载模型失败: %v", err)
		} else {
			logx.Info("成功从 TreeLite Protobuf 加载模型")
			return model, nil
		}
	}

	// 如果都失败，返回一个虚拟模型（使用规则排序）
	logx.Errorf("无法加载模型文件，将使用规则排序: %v", err)
	return model, nil
}

// loadFromJSON 从 XGBoost JSON 格式加载
func (m *XGBoostModel) loadFromJSON(jsonPath string) error {
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	// XGBoost JSON 格式解析
	var xgbModel struct {
		Learner struct {
			GradientBooster struct {
				Model struct {
					Trees []struct {
						TreeParam struct {
							NumNodes int `json:"num_nodes"`
						} `json:"tree_param"`
						SplitIndices    []int     `json:"split_indices"`
						SplitConditions []float64 `json:"split_conditions"`
						DefaultLeft     []int     `json:"default_left"`
						LeftChildren    []int     `json:"left_children"`
						RightChildren   []int     `json:"right_children"`
						LeafValues      []float64 `json:"base_weights"` // 叶子节点值
					} `json:"trees"`
				} `json:"model"`
			} `json:"gradient_booster"`
		} `json:"learner"`
	}

	if err := json.Unmarshal(data, &xgbModel); err != nil {
		return fmt.Errorf("解析 XGBoost JSON 失败: %v", err)
	}

	// 转换为内部格式
	trees := xgbModel.Learner.GradientBooster.Model.Trees
	m.Trees = make([]*DecisionTree, len(trees))

	for i, tree := range trees {
		numNodes := tree.TreeParam.NumNodes
		nodes := make([]TreeNode, numNodes)

		for j := 0; j < numNodes; j++ {
			node := TreeNode{
				NodeID: j,
			}

			// 判断是否为叶子节点
			if tree.LeftChildren[j] == -1 && tree.RightChildren[j] == -1 {
				// 叶子节点
				node.FeatureIndex = -1
				node.LeafValue = tree.LeafValues[j]
			} else {
				// 内部节点
				node.FeatureIndex = tree.SplitIndices[j]
				node.Threshold = tree.SplitConditions[j]
				node.LeftChild = tree.LeftChildren[j]
				node.RightChild = tree.RightChildren[j]
				node.MissingGoTo = tree.LeftChildren[j] // 缺失值默认走左子树
				if tree.DefaultLeft[j] == 0 {
					node.MissingGoTo = tree.RightChildren[j]
				}
			}

			nodes[j] = node
		}

		m.Trees[i] = &DecisionTree{Nodes: nodes}
	}

	logx.Infof("从 JSON 加载了 %d 棵树", len(m.Trees))
	return nil
}

// loadFromTreeLiteProto 从 TreeLite Protobuf 格式加载
func (m *XGBoostModel) loadFromTreeLiteProto(protoPath string) error {
	// TODO: 实现 TreeLite Protobuf 格式解析
	// 这里需要使用 protobuf 库来解析
	return fmt.Errorf("TreeLite Protobuf 格式暂未实现")
}

// Predict 预测单个样本
func (m *XGBoostModel) Predict(features []float64) float64 {
	if len(features) != m.NumFeatures {
		logx.Errorf("特征数量不匹配: expected=%d, got=%d", m.NumFeatures, len(features))
		return 0.5 // 返回默认分数
	}

	// 如果没有加载树，使用默认分数
	if len(m.Trees) == 0 {
		return 0.5
	}

	// 累加所有树的预测结果
	sum := 0.0
	for _, tree := range m.Trees {
		sum += m.predictTree(tree, features)
	}

	// XGBoost 的预测值需要通过 sigmoid 转换为概率
	// 对于二分类：p = 1 / (1 + exp(-sum))
	score := 1.0 / (1.0 + exp(-sum))

	return score
}

// predictTree 预测单棵树
func (m *XGBoostModel) predictTree(tree *DecisionTree, features []float64) float64 {
	nodeID := 0 // 从根节点开始

	for {
		node := tree.Nodes[nodeID]

		// 如果是叶子节点，返回值
		if node.FeatureIndex == -1 {
			return node.LeafValue
		}

		// 获取特征值
		featureValue := 0.0
		if node.FeatureIndex < len(features) {
			featureValue = features[node.FeatureIndex]
		}

		// 根据阈值决定走左子树还是右子树
		if featureValue < node.Threshold {
			nodeID = node.LeftChild
		} else {
			nodeID = node.RightChild
		}

		// 安全检查
		if nodeID < 0 || nodeID >= len(tree.Nodes) {
			return 0.0
		}
	}
}

// exp 指数函数
func exp(x float64) float64 {
	// 使用 Taylor 展开近似
	// e^x ≈ 1 + x + x²/2! + x³/3! + x⁴/4! + ...
	if x > 10 {
		return 22026.4657948 // e^10
	}
	if x < -10 {
		return 0.0000453999
	}

	result := 1.0
	term := 1.0

	for i := 1; i <= 20; i++ {
		term *= x / float64(i)
		result += term
		if term < 1e-10 && term > -1e-10 {
			break
		}
	}

	return result
}

// BatchPredict 批量预测
func (m *XGBoostModel) BatchPredict(featuresBatch [][]float64) []float64 {
	scores := make([]float64, len(featuresBatch))
	for i, features := range featuresBatch {
		scores[i] = m.Predict(features)
	}
	return scores
}

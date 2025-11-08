package rank

import (
	"math"
	"sort"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"mybilibili/app/recommend/cmd/rpc/model"
	"mybilibili/app/recommend/cmd/rpc/recommend"
)

// RankModelManager 排序模型管理器
type RankModelManager struct {
	modelDir string
	model    *XGBoostModel
}

// NewRankModelManager 创建排序模型管理器
func NewRankModelManager(modelDir string) *RankModelManager {
	// 尝试加载模型
	xgbModel, err := LoadXGBoostModel(modelDir)
	if err != nil {
		logx.Errorf("加载 XGBoost 模型失败: %v, 将使用规则排序", err)
	}

	return &RankModelManager{
		modelDir: modelDir,
		model:    xgbModel,
	}
}

// DoRank 执行排序（模型排序）
func (m *RankModelManager) DoRank(req *recommend.RecommendRequest, records []*model.RecommendRecord, profile *model.UserProfile) error {
	// 如果模型未加载或没有树，降级为规则排序
	if m.model == nil || len(m.model.Trees) == 0 {
		logx.Infof("模型未加载，使用规则排序")
		RuleBasedRank(records, profile)
		return nil
	}

	// 使用 XGBoost 模型进行排序
	logx.Infof("使用 XGBoost 模型排序: version=%s, num_trees=%d", m.model.ModelVersion, m.model.NumTrees)

	// 1. 为每个视频构建特征并预测分数
	for _, record := range records {
		features := BuildFeatureV13(record, profile)
		score := m.model.Predict(features)
		record.Score = score

		if req.Debug {
			logx.Infof("视频 %d 预测分数: %.4f", record.AVID, score)
		}
	}

	// 2. 按分数降序排序
	sort.Slice(records, func(i, j int) bool {
		return records[i].Score > records[j].Score
	})

	logx.Infof("模型排序完成: 排序数量=%d", len(records))
	return nil
}

// RuleBasedRank 规则排序（降级方案）
func RuleBasedRank(records []*model.RecommendRecord, profile *model.UserProfile) {
	// 计算每个视频的排序分数
	for _, record := range records {
		score := calculateScore(record, profile)
		record.Score = score
	}

	// 按分数降序排序
	sort.Slice(records, func(i, j int) bool {
		return records[i].Score > records[j].Score
	})
}

// calculateScore 计算排序分数（规则打分）
func calculateScore(record *model.RecommendRecord, profile *model.UserProfile) float64 {
	score := 0.0

	// 1. 召回类型权重
	if strings.Contains(record.RecallTypes, "SelectionRecall") {
		score += 10.0
	}
	if strings.Contains(record.RecallTypes, "LikeI2IRecall") {
		score += 8.0
	}
	if strings.Contains(record.RecallTypes, "FollowRecall") {
		score += 7.0
	}
	if strings.Contains(record.RecallTypes, "UserProfileRecall") {
		score += 5.0
	}
	if strings.Contains(record.RecallTypes, "HotRecall") {
		score += 3.0
	}

	// 2. 视频质量分数
	if record.State == model.StateSelection {
		score += 5.0
	} else if record.State == model.StateHighQuality {
		score += 3.0
	}

	// 3. 播放量分数（对数归一化）
	if record.PlayMonth > 0 {
		playScore := math.Log10(float64(record.PlayMonth)+1) / math.Log10(1000000)
		score += playScore * 10.0
	}

	// 4. 点赞率分数
	if record.PlayMonth > 0 && record.LikesMonth > 0 {
		likeRatio := float64(record.LikesMonth) / float64(record.PlayMonth)
		score += likeRatio * 50.0
	}

	// 5. 标签匹配分数
	matchCount := 0
	for _, tag := range record.Tags {
		if _, ok := profile.Tags[tag]; ok {
			matchCount++
		}
	}
	if matchCount > 0 {
		score += float64(matchCount) * 2.0
	}

	// 6. 时效性分数（新视频加权）
	// daysSincePublish := (time.Now().Unix() - record.PubTime) / 86400
	// if daysSincePublish < 3 {
	// 	score += 3.0
	// }

	return score
}

// BuildFeatureV13 构建 v13 模型的特征（65个特征）
func BuildFeatureV13(record *model.RecommendRecord, profile *model.UserProfile) []float64 {
	features := make([]float64, 65)
	idx := 0

	// 1. 分区特征 (35个 One-Hot)
	zoneFeatures := makeZoneBuckets(record.ZoneID, 35)
	copy(features[idx:], zoneFeatures)
	idx += 35

	// 2. 召回特征 (10个 One-Hot)
	recallFeatures := makeRecallFeatures(record.RecallTypes)
	copy(features[idx:], recallFeatures)
	idx += 10

	// 3. 状态特征 (5个 One-Hot)
	stateFeatures := makeStateBuckets(int32(record.State), 5)
	copy(features[idx:], stateFeatures)
	idx += 5

	// 4. 全站统计特征 (6个)
	features[idx] = math.Log10(float64(record.PlayHive)+1) / math.Log10(1000000)
	features[idx+1] = math.Log10(float64(record.LikesHive)+1) / math.Log10(100000)
	features[idx+2] = math.Log10(float64(record.FavHive)+1) / math.Log10(50000)
	features[idx+3] = math.Log10(float64(record.ReplyHive)+1) / math.Log10(10000)
	features[idx+4] = math.Log10(float64(record.ShareHive)+1) / math.Log10(5000)
	features[idx+5] = math.Log10(float64(record.CoinHive)+1) / math.Log10(20000)
	idx += 6

	// 5. 月度统计特征 (5个)
	features[idx] = math.Log10(float64(record.PlayMonth)+1) / math.Log10(100000)
	features[idx+1] = math.Log10(float64(record.LikesMonth)+1) / math.Log10(5000)
	features[idx+2] = math.Log10(float64(record.ReplyMonth)+1) / math.Log10(2000)
	features[idx+3] = math.Log10(float64(record.ShareMonth)+1) / math.Log10(1000)
	features[idx+4] = math.Log10(float64(record.PlayMonthFinish)+1) / math.Log10(50000)
	idx += 5

	// 6. 交叉特征 (4个)
	matchTagCount := countMatchTags(record.Tags, profile.Tags)
	features[idx] = float64(btoi(matchTagCount > 0))
	features[idx+1] = float64(matchTagCount)
	idx += 2

	return features
}

// makeZoneBuckets 创建分区特征桶（One-Hot）
func makeZoneBuckets(zoneID int32, size int) []float64 {
	buckets := make([]float64, size)
	if int(zoneID) < size {
		buckets[zoneID] = 1.0
	}
	return buckets
}

// makeRecallFeatures 创建召回特征（One-Hot）
func makeRecallFeatures(recallTypes string) []float64 {
	features := make([]float64, 10)
	types := []string{
		"HotRecall", "SelectionRecall", "LikeI2IRecall", "LikeTagRecall",
		"LikeUPRecall", "PosI2IRecall", "PosTagRecall", "UserProfileRecall",
		"FollowRecall", "RandomRecall",
	}

	for i, typ := range types {
		if strings.Contains(recallTypes, typ) {
			features[i] = 1.0
		}
	}

	return features
}

// makeStateBuckets 创建状态特征桶（One-Hot）
func makeStateBuckets(state int32, size int) []float64 {
	buckets := make([]float64, size)
	if int(state) < size && state > 0 {
		buckets[state-1] = 1.0
	}
	return buckets
}

// countMatchTags 统计匹配的标签数量
func countMatchTags(tags []string, userTags map[string]float64) int {
	count := 0
	for _, tag := range tags {
		if _, ok := userTags[tag]; ok {
			count++
		}
	}
	return count
}

// btoi bool to int
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

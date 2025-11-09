package rank

import (
	"fmt"
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
			// 输出前10个特征值用于调试
			featureStr := ""
			for i := 0; i < 10 && i < len(features); i++ {
				featureStr += fmt.Sprintf("f[%d]=%.4f ", i, features[i])
			}
			logx.Infof("视频 %d 预测分数: %.4f, 特征: %s", record.AVID, score, featureStr)
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
	logx.Infof("使用规则排序: 视频数量=%d", len(records))

	// 计算每个视频的排序分数
	for _, record := range records {
		score := calculateScore(record, profile)
		record.Score = score
		logx.Infof("视频 %d 规则分数: %.4f, 召回类型=%s, 状态=%d, 播放量=%d, 点赞数=%d",
			record.AVID, score, record.RecallTypes, record.State, record.PlayMonth, record.LikesMonth)
	}

	// 按分数降序排序
	sort.Slice(records, func(i, j int) bool {
		return records[i].Score > records[j].Score
	})

	logx.Infof("规则排序完成: 排序数量=%d", len(records))
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
// 特征顺序必须与训练脚本 train_v13_final.py 中的 FEATURES_V13 完全匹配
func BuildFeatureV13(record *model.RecommendRecord, profile *model.UserProfile) []float64 {
	features := make([]float64, 65)

	// 特征顺序与训练脚本完全匹配
	// 训练脚本中的特征顺序：
	// ["zone-bucket-168", "zone-bucket-75", "play_hive", "zone-bucket-95", "likes_month",
	//  "state-bucket-3", "share_month", "recall-PosTagRecall", "zone-bucket-124", "recall-PosI2IRecall",
	//  "state-bucket-4", "recall-SelectionRecall", "zone-bucket-156", "contains_tag_count", "zone-bucket-158",
	//  "zone-bucket-183", "zone-bucket-184", "zone-bucket-21", "zone-bucket-154", "zone-bucket-159",
	//  "zone-bucket-85", "recall-LikeUPRecall", "reply_month", "state-bucket-1", "zone-bucket-96",
	//  "has_tag_count", "zone-bucket-86", "zone-bucket-138", "zone-bucket-182", "play_month_finish",
	//  "recall-HotRecall", "zone-bucket-157", "zone-bucket-20", "zone-bucket-39", "zone-bucket-161",
	//  "reply_hive", "recall-LikeTagRecall", "zone-bucket-76", "zone-bucket-98", "state-bucket-5",
	//  "zone-bucket-22", "zone-bucket-27", "zone-bucket-122", "zone-bucket-176", "recall-UserProfileBBQ",
	//  "recall-UserProfileBili", "zone-bucket-163", "zone-bucket-30", "zone-bucket-31", "zone-bucket-59",
	//  "recall-LikeI2IRecall", "zone-bucket-25", "zone-bucket-28", "zone-bucket-24", "zone-bucket-29",
	//  "zone-bucket-164", "coin_hive", "play_month", "share_hive", "recall-RandomRecall",
	//  "fav_hive", "zone-bucket-162", "likes_hive", "recall-FollowRecall", "zone-bucket-47"]

	// 构建特征映射
	zoneBucketMap := makeZoneBucketMap(record.ZoneID)
	recallMap := makeRecallMap(record.RecallTypes)
	stateBucketMap := makeStateBucketMap(int32(record.State))
	matchTagCount := countMatchTags(record.Tags, profile.Tags)

	// 按照训练脚本的顺序填充特征
	idx := 0
	// 0: zone-bucket-168
	features[idx] = zoneBucketMap[168]
	idx++
	// 1: zone-bucket-75
	features[idx] = zoneBucketMap[75]
	idx++
	// 2: play_hive (使用 log1p，与训练脚本一致)
	features[idx] = math.Log1p(float64(record.PlayHive))
	idx++
	// 3: zone-bucket-95
	features[idx] = zoneBucketMap[95]
	idx++
	// 4: likes_month (使用 log1p，与训练脚本一致)
	features[idx] = math.Log1p(float64(record.LikesMonth))
	idx++
	// 5: state-bucket-3
	features[idx] = stateBucketMap[3]
	idx++
	// 6: share_month (使用 log1p，与训练脚本一致)
	features[idx] = math.Log1p(float64(record.ShareMonth))
	idx++
	// 7: recall-PosTagRecall
	features[idx] = recallMap["PosTagRecall"]
	idx++
	// 8: zone-bucket-124
	features[idx] = zoneBucketMap[124]
	idx++
	// 9: recall-PosI2IRecall
	features[idx] = recallMap["PosI2IRecall"]
	idx++
	// 10: state-bucket-4
	features[idx] = stateBucketMap[4]
	idx++
	// 11: recall-SelectionRecall
	features[idx] = recallMap["SelectionRecall"]
	idx++
	// 12: zone-bucket-156
	features[idx] = zoneBucketMap[156]
	idx++
	// 13: contains_tag_count
	features[idx] = float64(matchTagCount)
	idx++
	// 14: zone-bucket-158
	features[idx] = zoneBucketMap[158]
	idx++
	// 15: zone-bucket-183
	features[idx] = zoneBucketMap[183]
	idx++
	// 16: zone-bucket-184
	features[idx] = zoneBucketMap[184]
	idx++
	// 17: zone-bucket-21
	features[idx] = zoneBucketMap[21]
	idx++
	// 18: zone-bucket-154
	features[idx] = zoneBucketMap[154]
	idx++
	// 19: zone-bucket-159
	features[idx] = zoneBucketMap[159]
	idx++
	// 20: zone-bucket-85
	features[idx] = zoneBucketMap[85]
	idx++
	// 21: recall-LikeUPRecall
	features[idx] = recallMap["LikeUPRecall"]
	idx++
	// 22: reply_month (使用 log1p，与训练脚本一致)
	features[idx] = math.Log1p(float64(record.ReplyMonth))
	idx++
	// 23: state-bucket-1
	features[idx] = stateBucketMap[1]
	idx++
	// 24: zone-bucket-96
	features[idx] = zoneBucketMap[96]
	idx++
	// 25: has_tag_count
	features[idx] = float64(btoi(matchTagCount > 0))
	idx++
	// 26: zone-bucket-86
	features[idx] = zoneBucketMap[86]
	idx++
	// 27: zone-bucket-138
	features[idx] = zoneBucketMap[138]
	idx++
	// 28: zone-bucket-182
	features[idx] = zoneBucketMap[182]
	idx++
	// 29: play_month_finish (使用 log1p，与训练脚本一致)
	features[idx] = math.Log1p(float64(record.PlayMonthFinish))
	idx++
	// 30: recall-HotRecall
	features[idx] = recallMap["HotRecall"]
	idx++
	// 31: zone-bucket-157
	features[idx] = zoneBucketMap[157]
	idx++
	// 32: zone-bucket-20
	features[idx] = zoneBucketMap[20]
	idx++
	// 33: zone-bucket-39
	features[idx] = zoneBucketMap[39]
	idx++
	// 34: zone-bucket-161
	features[idx] = zoneBucketMap[161]
	idx++
	// 35: reply_hive (使用 log1p，与训练脚本一致)
	features[idx] = math.Log1p(float64(record.ReplyHive))
	idx++
	// 36: recall-LikeTagRecall
	features[idx] = recallMap["LikeTagRecall"]
	idx++
	// 37: zone-bucket-76
	features[idx] = zoneBucketMap[76]
	idx++
	// 38: zone-bucket-98
	features[idx] = zoneBucketMap[98]
	idx++
	// 39: state-bucket-5
	features[idx] = stateBucketMap[5]
	idx++
	// 40: zone-bucket-22
	features[idx] = zoneBucketMap[22]
	idx++
	// 41: zone-bucket-27
	features[idx] = zoneBucketMap[27]
	idx++
	// 42: zone-bucket-122
	features[idx] = zoneBucketMap[122]
	idx++
	// 43: zone-bucket-176
	features[idx] = zoneBucketMap[176]
	idx++
	// 44: recall-UserProfileBBQ
	features[idx] = recallMap["UserProfileBBQ"]
	idx++
	// 45: recall-UserProfileBili
	features[idx] = recallMap["UserProfileBili"]
	idx++
	// 46: zone-bucket-163
	features[idx] = zoneBucketMap[163]
	idx++
	// 47: zone-bucket-30
	features[idx] = zoneBucketMap[30]
	idx++
	// 48: zone-bucket-31
	features[idx] = zoneBucketMap[31]
	idx++
	// 49: zone-bucket-59
	features[idx] = zoneBucketMap[59]
	idx++
	// 50: recall-LikeI2IRecall
	features[idx] = recallMap["LikeI2IRecall"]
	idx++
	// 51: zone-bucket-25
	features[idx] = zoneBucketMap[25]
	idx++
	// 52: zone-bucket-28
	features[idx] = zoneBucketMap[28]
	idx++
	// 53: zone-bucket-24
	features[idx] = zoneBucketMap[24]
	idx++
	// 54: zone-bucket-29
	features[idx] = zoneBucketMap[29]
	idx++
	// 55: zone-bucket-164
	features[idx] = zoneBucketMap[164]
	idx++
	// 56: coin_hive (使用 log1p，与训练脚本一致)
	features[idx] = math.Log1p(float64(record.CoinHive))
	idx++
	// 57: play_month (使用 log1p，与训练脚本一致)
	features[idx] = math.Log1p(float64(record.PlayMonth))
	idx++
	// 58: share_hive (使用 log1p，与训练脚本一致)
	features[idx] = math.Log1p(float64(record.ShareHive))
	idx++
	// 59: recall-RandomRecall
	features[idx] = recallMap["RandomRecall"]
	idx++
	// 60: fav_hive (使用 log1p，与训练脚本一致)
	features[idx] = math.Log1p(float64(record.FavHive))
	idx++
	// 61: zone-bucket-162
	features[idx] = zoneBucketMap[162]
	idx++
	// 62: likes_hive (使用 log1p，与训练脚本一致)
	features[idx] = math.Log1p(float64(record.LikesHive))
	idx++
	// 63: recall-FollowRecall
	features[idx] = recallMap["FollowRecall"]
	idx++
	// 64: zone-bucket-47
	features[idx] = zoneBucketMap[47]
	idx++

	return features
}

// makeZoneBucketMap 创建分区特征映射
func makeZoneBucketMap(zoneID int32) map[int32]float64 {
	buckets := make(map[int32]float64)
	buckets[zoneID] = 1.0
	return buckets
}

// makeRecallMap 创建召回特征映射
func makeRecallMap(recallTypes string) map[string]float64 {
	recallMap := make(map[string]float64)
	types := []string{
		"HotRecall", "SelectionRecall", "LikeI2IRecall", "LikeTagRecall",
		"LikeUPRecall", "PosI2IRecall", "PosTagRecall", "UserProfileRecall",
		"UserProfileBBQ", "UserProfileBili", "FollowRecall", "RandomRecall",
	}

	for _, typ := range types {
		if strings.Contains(recallTypes, typ) {
			recallMap[typ] = 1.0
		}
	}

	return recallMap
}

// makeStateBucketMap 创建状态特征映射
func makeStateBucketMap(state int32) map[int32]float64 {
	buckets := make(map[int32]float64)
	if state > 0 {
		buckets[state] = 1.0
	}
	return buckets
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

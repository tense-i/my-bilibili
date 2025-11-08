package filter

import (
	"mybilibili/app/recommend/cmd/rpc/model"
	"mybilibili/app/recommend/cmd/rpc/recommend"
)

// FilterManager 过滤器管理器
type FilterManager struct {
	filters []Filter
}

// Filter 过滤器接口
type Filter interface {
	DoFilter(req *recommend.RecommendRequest, records []*model.RecommendRecord, profile *model.UserProfile) []*model.RecommendRecord
}

// NewFilterManager 创建过滤器管理器
func NewFilterManager() *FilterManager {
	return &FilterManager{
		filters: []Filter{
			&DefaultFilter{},
			&BlackFilter{},
			&DurationFilter{},
		},
	}
}

// PreFilter 初步过滤
func (m *FilterManager) PreFilter(req *recommend.RecommendRequest, records []*model.RecommendRecord, profile *model.UserProfile) []*model.RecommendRecord {
	for _, filter := range m.filters {
		records = filter.DoFilter(req, records, profile)
	}
	return records
}

// DefaultFilter 精确去重过滤器
type DefaultFilter struct{}

func (f *DefaultFilter) DoFilter(req *recommend.RecommendRequest, records []*model.RecommendRecord, profile *model.UserProfile) []*model.RecommendRecord {
	seen := make(map[int64]bool)
	filtered := make([]*model.RecommendRecord, 0, len(records))

	for _, record := range records {
		if !seen[record.AVID] {
			seen[record.AVID] = true
			filtered = append(filtered, record)
		}
	}

	return filtered
}

// BlackFilter 黑名单过滤器
type BlackFilter struct{}

func (f *BlackFilter) DoFilter(req *recommend.RecommendRequest, records []*model.RecommendRecord, profile *model.UserProfile) []*model.RecommendRecord {
	if len(profile.BlackUps) == 0 {
		return records
	}

	filtered := make([]*model.RecommendRecord, 0, len(records))

	for _, record := range records {
		if !profile.BlackUps[record.UPMID] {
			filtered = append(filtered, record)
		}
	}

	return filtered
}

// DurationFilter 时长过滤器
type DurationFilter struct{}

func (f *DurationFilter) DoFilter(req *recommend.RecommendRequest, records []*model.RecommendRecord, profile *model.UserProfile) []*model.RecommendRecord {
	filtered := make([]*model.RecommendRecord, 0, len(records))

	for _, record := range records {
		// B站主站视频时长控制：1分钟-60分钟
		if record.Duration >= 60 && record.Duration <= 3600 {
			filtered = append(filtered, record)
		}
	}

	return filtered
}

// BloomFilter 布隆过滤器（在后处理阶段使用）
type BloomFilter struct {
	checkFunc func(mid int64, avid int64) bool
}

func NewBloomFilter(checkFunc func(mid int64, avid int64) bool) *BloomFilter {
	return &BloomFilter{
		checkFunc: checkFunc,
	}
}

func (f *BloomFilter) DoFilter(req *recommend.RecommendRequest, records []*model.RecommendRecord, profile *model.UserProfile) []*model.RecommendRecord {
	if f.checkFunc == nil || req.Mid == 0 {
		return records
	}

	filtered := make([]*model.RecommendRecord, 0, len(records))

	for _, record := range records {
		if !f.checkFunc(req.Mid, record.AVID) {
			filtered = append(filtered, record)
		}
	}

	return filtered
}

// FollowsFilter 关注过滤器 - 过滤掉用户已关注的UP主的视频
// 应用场景：UP主推荐列表，避免推荐用户已关注的UP主
type FollowsFilter struct{}

func (f *FollowsFilter) DoFilter(req *recommend.RecommendRequest, records []*model.RecommendRecord, profile *model.UserProfile) []*model.RecommendRecord {
	if len(profile.FollowUps) == 0 {
		return records
	}

	filtered := make([]*model.RecommendRecord, 0, len(records))

	for _, record := range records {
		// 检查UP主是否在关注列表中
		if _, followed := profile.FollowUps[record.UPMID]; !followed {
			filtered = append(filtered, record)
		}
	}

	return filtered
}

// RelatedFilter 相关过滤器 - 过滤掉当前视频和同UP主的视频
// 应用场景：相关推荐，避免推荐当前视频本身和同一UP主的其他视频
type RelatedFilter struct {
	sourceAVID  int64 // 当前视频的AVID
	sourceUPMID int64 // 当前视频的UP主MID
}

func NewRelatedFilter(sourceAVID int64, sourceUPMID int64) *RelatedFilter {
	return &RelatedFilter{
		sourceAVID:  sourceAVID,
		sourceUPMID: sourceUPMID,
	}
}

func (f *RelatedFilter) DoFilter(req *recommend.RecommendRequest, records []*model.RecommendRecord, profile *model.UserProfile) []*model.RecommendRecord {
	filtered := make([]*model.RecommendRecord, 0, len(records))

	for _, record := range records {
		// 过滤掉当前视频本身
		if record.AVID == f.sourceAVID {
			continue
		}
		// 过滤掉同一UP主的其他视频
		if record.UPMID == f.sourceUPMID {
			continue
		}
		filtered = append(filtered, record)
	}

	return filtered
}

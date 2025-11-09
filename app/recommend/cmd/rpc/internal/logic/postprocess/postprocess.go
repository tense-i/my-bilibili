package postprocess

import (
	"context"

	"mybilibili/app/recommend/cmd/rpc/internal/svc"
	"mybilibili/app/recommend/cmd/rpc/model"
	"mybilibili/app/recommend/cmd/rpc/recommend"
)

// PostProcessor 后处理器
type PostProcessor struct {
	svcCtx *svc.ServiceContext
}

// NewPostProcessor 创建后处理器
func NewPostProcessor(svcCtx *svc.ServiceContext) *PostProcessor {
	return &PostProcessor{
		svcCtx: svcCtx,
	}
}

// Process 执行后处理
func (p *PostProcessor) Process(ctx context.Context, req *recommend.RecommendRequest, records []*model.RecommendRecord, profile *model.UserProfile) []*model.RecommendRecord {
	if len(records) == 0 {
		return records
	}

	// 1. 时长过滤（B站主站视频时长：60-3600秒）
	// 注意：如果配置为0，不进行时长过滤
	durationMin := p.svcCtx.BusinessConfig.DurationMin
	durationMax := p.svcCtx.BusinessConfig.DurationMax
	if durationMin > 0 || durationMax > 0 {
		records = p.durationFilterProcess(records)
	}

	// 2. 布隆过滤器去重（检查用户是否看过）
	if p.svcCtx.BusinessConfig.BloomFilterEnable && req.Mid > 0 {
		records = p.bloomFilterProcess(ctx, req.Mid, records)
	}

	// 3. 打散策略（避免连续同UP主/同标签）
	// 只在记录数大于3时才进行打散
	if len(records) > 3 {
		records = p.scatterProcess(records)
	}

	return records
}

// durationFilterProcess 时长过滤处理
func (p *PostProcessor) durationFilterProcess(records []*model.RecommendRecord) []*model.RecommendRecord {
	minDuration := p.svcCtx.BusinessConfig.DurationMin
	maxDuration := p.svcCtx.BusinessConfig.DurationMax

	// 如果配置为0，使用默认值
	if minDuration == 0 {
		minDuration = 60 // 1分钟
	}
	if maxDuration == 0 {
		maxDuration = 3600 // 60分钟
	}

	filtered := make([]*model.RecommendRecord, 0, len(records))
	for _, record := range records {
		if record.Duration >= int32(minDuration) && record.Duration <= int32(maxDuration) {
			filtered = append(filtered, record)
		}
	}

	return filtered
}

// bloomFilterProcess 布隆过滤器处理
func (p *PostProcessor) bloomFilterProcess(ctx context.Context, mid int64, records []*model.RecommendRecord) []*model.RecommendRecord {
	filtered := make([]*model.RecommendRecord, 0, len(records))

	for _, record := range records {
		// 检查用户是否看过该视频
		seen, err := p.svcCtx.Dao.CheckBloomFilter(ctx, mid, record.AVID)
		if err != nil {
			// 如果检查出错，保留记录（保守策略）
			filtered = append(filtered, record)
			continue
		}
		// 如果用户没看过（seen=false），保留记录
		if !seen {
			filtered = append(filtered, record)
		}
		// 如果用户看过（seen=true），过滤掉该记录
	}

	return filtered
}

// scatterProcess 打散处理
func (p *PostProcessor) scatterProcess(records []*model.RecommendRecord) []*model.RecommendRecord {
	if len(records) <= 3 {
		return records
	}

	scattered := make([]*model.RecommendRecord, 0, len(records))
	upMIDLastIdx := make(map[int64]int)
	tagLastIdx := make(map[string]int)

	minUPDistance := p.svcCtx.BusinessConfig.ScatterUPMinDistance
	minTagDistance := p.svcCtx.BusinessConfig.ScatterTagMinDistance

	for _, record := range records {
		shouldScatter := false

		// 检查UP主打散
		if lastIdx, exists := upMIDLastIdx[record.UPMID]; exists {
			if len(scattered)-lastIdx < minUPDistance {
				shouldScatter = true
			}
		}

		// 检查标签打散
		if !shouldScatter {
			for _, tag := range record.Tags {
				if lastIdx, exists := tagLastIdx[tag]; exists {
					if len(scattered)-lastIdx < minTagDistance {
						shouldScatter = true
						break
					}
				}
			}
		}

		if !shouldScatter {
			scattered = append(scattered, record)
			idx := len(scattered) - 1
			upMIDLastIdx[record.UPMID] = idx
			for _, tag := range record.Tags {
				tagLastIdx[tag] = idx
			}
		}
	}

	// 如果打散后数量太少，追加部分未打散的
	if len(scattered) < len(records)/2 {
		for _, record := range records {
			found := false
			for _, s := range scattered {
				if s.AVID == record.AVID {
					found = true
					break
				}
			}
			if !found {
				scattered = append(scattered, record)
			}
		}
	}

	return scattered
}

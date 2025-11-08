package logic

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"mybilibili/app/recall/cmd/rpc/internal/svc"
	"mybilibili/app/recall/cmd/rpc/model"
	"mybilibili/app/recall/cmd/rpc/recall"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecallLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecallLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecallLogic {
	return &RecallLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Recall 召回接口 - 多路召回并合并
func (l *RecallLogic) Recall(in *recall.RecallRequest) (*recall.RecallResponse, error) {
	logx.Infof("召回请求: mid=%d, total_limit=%d, infos_count=%d", in.Mid, in.TotalLimit, len(in.Infos))

	// 存储所有召回结果
	allItems := make([]*model.RecallItem, 0)
	itemMap := make(map[int64]*model.RecallItem) // 用于去重

	// 并发执行所有召回策略
	for _, info := range in.Infos {
		items, err := l.executeRecallStrategy(in, info)
		if err != nil {
			logx.Errorf("召回策略 %s 执行失败: %v", info.Name, err)
			continue
		}

		// 合并结果并去重
		for _, item := range items {
			item.RecallType = info.Name
			item.RecallTag = info.Tag
			item.Priority = info.Priority

			// 如果已存在，保留优先级高的
			if existing, ok := itemMap[item.AVID]; ok {
				if item.Priority > existing.Priority {
					itemMap[item.AVID] = item
				}
			} else {
				itemMap[item.AVID] = item
			}
		}
	}

	// 转换为列表
	for _, item := range itemMap {
		allItems = append(allItems, item)
	}

	// 按优先级和分数排序
	sort.Slice(allItems, func(i, j int) bool {
		if allItems[i].Priority != allItems[j].Priority {
			return allItems[i].Priority > allItems[j].Priority
		}
		return allItems[i].Score > allItems[j].Score
	})

	// 限制返回数量
	if len(allItems) > int(in.TotalLimit) {
		allItems = allItems[:in.TotalLimit]
	}

	// 转换为 proto 格式
	list := make([]*recall.RecallItem, 0, len(allItems))
	for _, item := range allItems {
		list = append(list, &recall.RecallItem{
			Avid:       item.AVID,
			Score:      item.Score,
			RecallType: item.RecallType,
			RecallTag:  item.RecallTag,
			Priority:   item.Priority,
			Extra:      item.Extra,
		})
	}

	logx.Infof("召回完成: 召回数量=%d", len(list))

	return &recall.RecallResponse{
		List:  list,
		Total: int32(len(list)),
		DebugInfo: map[string]string{
			"recall_count": fmt.Sprintf("%d", len(list)),
		},
	}, nil
}

// executeRecallStrategy 执行单个召回策略
func (l *RecallLogic) executeRecallStrategy(req *recall.RecallRequest, info *recall.RecallInfo) ([]*model.RecallItem, error) {
	var items []*model.RecallItem
	var err error

	switch info.Name {
	case model.HotRecall:
		items, err = l.hotRecall(req, info)
	case model.SelectionRecall:
		items, err = l.selectionRecall(req, info)
	case model.LikeI2IRecall:
		items, err = l.likeI2IRecall(req, info)
	case model.PosI2IRecall:
		items, err = l.posI2IRecall(req, info)
	case model.LikeTagRecall:
		items, err = l.tagRecall(req, info)
	case model.PosTagRecall:
		items, err = l.posTagRecall(req, info)
	case model.UserProfileRecall, model.UserProfileBili, model.UserProfileBBQ:
		items, err = l.tagRecall(req, info)
	case model.FollowRecall, model.LikeUPRecall:
		items, err = l.upRecall(req, info)
	case model.RandomRecall:
		items, err = l.randomRecall(req, info)
	default:
		items, err = l.defaultRecall(req, info)
	}

	if err != nil {
		return nil, err
	}

	// 应用布隆过滤器
	if info.Filter == "bloomfilter" && req.Mid > 0 {
		items = l.applyBloomFilter(req.Mid, items)
	}

	// 限制数量
	if len(items) > int(info.Limit) {
		items = items[:info.Limit]
	}

	logx.Infof("召回策略 %s 完成: 召回数量=%d", info.Name, len(items))

	return items, nil
}

// hotRecall 热门召回
func (l *RecallLogic) hotRecall(req *recall.RecallRequest, info *recall.RecallInfo) ([]*model.RecallItem, error) {
	avids, err := l.svcCtx.Dao.GetHotVideos(l.ctx, int64(info.Limit), 0)
	if err != nil {
		return nil, err
	}

	items := make([]*model.RecallItem, 0, len(avids))
	for _, avid := range avids {
		items = append(items, &model.RecallItem{
			AVID:  avid,
			Score: 1.0,
		})
	}
	return items, nil
}

// selectionRecall 精选召回
func (l *RecallLogic) selectionRecall(req *recall.RecallRequest, info *recall.RecallInfo) ([]*model.RecallItem, error) {
	avids, err := l.svcCtx.Dao.GetSelectionVideos(l.ctx, int64(info.Limit))
	if err != nil {
		return nil, err
	}

	items := make([]*model.RecallItem, 0, len(avids))
	for _, avid := range avids {
		items = append(items, &model.RecallItem{
			AVID:  avid,
			Score: 1.5, // 精选内容分数更高
		})
	}
	return items, nil
}

// likeI2IRecall 点赞I2I召回
func (l *RecallLogic) likeI2IRecall(req *recall.RecallRequest, info *recall.RecallInfo) ([]*model.RecallItem, error) {
	// 从 info.Tag 中提取视频ID（格式：RECALL:I2I:{avid}）
	var sourceAVID int64
	fmt.Sscanf(info.Tag, "RECALL:I2I:%d", &sourceAVID)

	if sourceAVID == 0 {
		return nil, fmt.Errorf("invalid I2I tag: %s", info.Tag)
	}

	// 获取相似视频
	avids, err := l.svcCtx.Dao.GetI2IVideos(l.ctx, sourceAVID, int64(info.Limit))
	if err != nil {
		return nil, err
	}

	items := make([]*model.RecallItem, 0, len(avids))
	for _, avid := range avids {
		items = append(items, &model.RecallItem{
			AVID:  avid,
			Score: 1.2,
			Extra: map[string]string{
				"source_avid": fmt.Sprintf("%d", sourceAVID),
			},
		})
	}

	return items, nil
}

// posI2IRecall 正反馈I2I召回
func (l *RecallLogic) posI2IRecall(req *recall.RecallRequest, info *recall.RecallInfo) ([]*model.RecallItem, error) {
	// 从 info.Tag 中提取视频ID（格式：RECALL:I2I:{avid}）
	var sourceAVID int64
	fmt.Sscanf(info.Tag, "RECALL:I2I:%d", &sourceAVID)

	if sourceAVID == 0 {
		return nil, fmt.Errorf("invalid I2I tag: %s", info.Tag)
	}

	// 获取相似视频
	avids, err := l.svcCtx.Dao.GetI2IVideos(l.ctx, sourceAVID, int64(info.Limit))
	if err != nil {
		return nil, err
	}

	items := make([]*model.RecallItem, 0, len(avids))
	for _, avid := range avids {
		items = append(items, &model.RecallItem{
			AVID:  avid,
			Score: 1.1, // 正反馈的分数略低于点赞
			Extra: map[string]string{
				"source_avid": fmt.Sprintf("%d", sourceAVID),
			},
		})
	}

	return items, nil
}

// posTagRecall 正反馈标签召回
func (l *RecallLogic) posTagRecall(req *recall.RecallRequest, info *recall.RecallInfo) ([]*model.RecallItem, error) {
	// info.Tag 可能是标签名或标签ID
	// 格式：RECALL:HOT_T:{tag_id} 或 RECALL:HOT_T:{tag_name}
	avids, err := l.svcCtx.Dao.GetTagVideos(l.ctx, info.Tag, int64(info.Limit))
	if err != nil {
		return nil, err
	}

	items := make([]*model.RecallItem, 0, len(avids))
	for _, avid := range avids {
		items = append(items, &model.RecallItem{
			AVID:  avid,
			Score: 1.0,
		})
	}
	return items, nil
}

// randomRecall 随机召回（新发布视频）
func (l *RecallLogic) randomRecall(req *recall.RecallRequest, info *recall.RecallInfo) ([]*model.RecallItem, error) {
	// 从 info.Tag 中提取标签名（格式：RECALL:T:{tag_name}）
	var tag string
	if strings.HasPrefix(info.Tag, "RECALL:T:") {
		tag = strings.TrimPrefix(info.Tag, "RECALL:T:")
	} else {
		tag = info.Tag
	}

	avids, err := l.svcCtx.Dao.GetRandomVideos(l.ctx, tag, int64(info.Limit))
	if err != nil {
		return nil, err
	}

	items := make([]*model.RecallItem, 0, len(avids))
	for _, avid := range avids {
		items = append(items, &model.RecallItem{
			AVID:  avid,
			Score: 0.8, // 随机召回分数较低
			Extra: map[string]string{
				"tag": tag,
			},
		})
	}
	return items, nil
}

// tagRecall 标签召回
func (l *RecallLogic) tagRecall(req *recall.RecallRequest, info *recall.RecallInfo) ([]*model.RecallItem, error) {
	// info.Tag 应该是标签名称
	avids, err := l.svcCtx.Dao.GetTagVideos(l.ctx, info.Tag, int64(info.Limit))
	if err != nil {
		return nil, err
	}

	items := make([]*model.RecallItem, 0, len(avids))
	for _, avid := range avids {
		items = append(items, &model.RecallItem{
			AVID:  avid,
			Score: 1.0,
		})
	}
	return items, nil
}

// upRecall UP主召回
func (l *RecallLogic) upRecall(req *recall.RecallRequest, info *recall.RecallInfo) ([]*model.RecallItem, error) {
	// 获取用户关注的 UP 主
	upMids, err := l.svcCtx.Dao.GetUserFollowUPs(l.ctx, req.Mid)
	if err != nil {
		return nil, err
	}

	allAvids := make([]int64, 0)
	for _, upMid := range upMids {
		avids, err := l.svcCtx.Dao.GetUPVideos(l.ctx, upMid, int64(info.Limit))
		if err != nil {
			continue
		}
		allAvids = append(allAvids, avids...)
	}

	// 去重并限制数量
	avidMap := make(map[int64]bool)
	items := make([]*model.RecallItem, 0)
	for _, avid := range allAvids {
		if !avidMap[avid] && len(items) < int(info.Limit) {
			avidMap[avid] = true
			items = append(items, &model.RecallItem{
				AVID:  avid,
				Score: 1.3,
			})
		}
	}

	return items, nil
}

// defaultRecall 默认召回（从 ZSet）
func (l *RecallLogic) defaultRecall(req *recall.RecallRequest, info *recall.RecallInfo) ([]*model.RecallItem, error) {
	// 默认使用热门召回
	return l.hotRecall(req, info)
}

// applyBloomFilter 应用布隆过滤器
func (l *RecallLogic) applyBloomFilter(mid int64, items []*model.RecallItem) []*model.RecallItem {
	filtered := make([]*model.RecallItem, 0)

	for _, item := range items {
		exists, err := l.svcCtx.Dao.CheckBloomFilter(l.ctx, mid, item.AVID)
		if err != nil || !exists {
			// 不存在于 Bloom Filter 中，可以推荐
			filtered = append(filtered, item)
		}
	}

	return filtered
}

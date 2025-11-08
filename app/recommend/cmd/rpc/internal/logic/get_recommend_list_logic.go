package logic

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"mybilibili/app/recall/cmd/rpc/recall"
	"mybilibili/app/recommend/cmd/rpc/internal/logic/filter"
	"mybilibili/app/recommend/cmd/rpc/internal/logic/postprocess"
	"mybilibili/app/recommend/cmd/rpc/internal/logic/rank"
	"mybilibili/app/recommend/cmd/rpc/internal/svc"
	"mybilibili/app/recommend/cmd/rpc/model"
	"mybilibili/app/recommend/cmd/rpc/recommend"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRecommendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRecommendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecommendListLogic {
	return &GetRecommendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetRecommendList 获取推荐列表 - 完整推荐流程
func (l *GetRecommendListLogic) GetRecommendList(in *recommend.RecommendRequest) (*recommend.RecommendResponse, error) {
	startTime := time.Now()
	logx.Infof("推荐请求开始: mid=%d, limit=%d, page=%d", in.Mid, in.Limit, in.Page)

	// 1. 加载用户画像
	userProfile, err := l.svcCtx.Dao.GetUserProfile(l.ctx, in.Mid)
	if err != nil {
		logx.Errorf("加载用户画像失败: %v", err)
		return nil, err
	}

	// 2. 加载用户黑名单（用户画像中已包含关注信息）
	blacklist, err := l.svcCtx.Dao.GetUserBlacklist(l.ctx, in.Mid)
	if err == nil {
		for targetID := range blacklist {
			if userProfile.BlackUps == nil {
				userProfile.BlackUps = make(map[int64]bool)
			}
			userProfile.BlackUps[targetID] = true
		}
	}

	logx.Infof("用户画像加载完成: follow_count=%d, black_count=%d", len(userProfile.FollowUps), len(userProfile.BlackUps))

	// 3. 构建召回请求并调用召回服务
	recallReq := l.buildRecallRequest(in, userProfile)
	recallResp, err := l.svcCtx.RecallRpc.Recall(l.ctx, recallReq)
	if err != nil || recallResp == nil || len(recallResp.List) == 0 {
		logx.Errorf("召回失败或结果为空: %v", err)
		// 降级召回
		recallResp, err = l.downGradeRecall(in, userProfile)
		if err != nil || recallResp == nil {
			return &recommend.RecommendResponse{
				List:    make([]*recommend.RecommendItem, 0),
				Total:   0,
				HasMore: false,
				Message: "召回失败",
			}, nil
		}
	}

	logx.Infof("召回完成: 召回数量=%d", len(recallResp.List))

	// 4. 转换召回结果并获取视频详细信息
	records := l.transformRecallToRecords(recallResp.List)
	records = l.enrichVideoInfo(records)

	logx.Infof("视频信息丰富完成: 数量=%d", len(records))

	// 5. 初步过滤（精确去重 + 黑名单）
	filterManager := filter.NewFilterManager()
	records = filterManager.PreFilter(in, records, userProfile)

	logx.Infof("初步过滤完成: 剩余数量=%d", len(records))

	// 6. 排序（XGBoost 模型或规则排序）
	if l.svcCtx.BusinessConfig.RankEnable {
		rankManager := rank.NewRankModelManager(l.svcCtx.Config.RankModel.ModelDir)
		err = rankManager.DoRank(in, records, userProfile)
		if err != nil && l.svcCtx.BusinessConfig.RankFallback {
			logx.Errorf("模型排序失败，降级为规则排序: %v", err)
			rank.RuleBasedRank(records, userProfile)
		}
	} else {
		rank.RuleBasedRank(records, userProfile)
	}

	logx.Infof("排序完成")

	// 7. 后处理（布隆过滤 + 时长过滤 + 打散）
	processor := postprocess.NewPostProcessor(l.svcCtx)
	records = processor.Process(l.ctx, in, records, userProfile)

	logx.Infof("后处理完成: 最终数量=%d", len(records))

	// 8. 分页
	total := len(records)
	if in.Limit == 0 {
		in.Limit = 20
	}
	if in.Page == 0 {
		in.Page = 1
	}

	start := int((in.Page - 1) * in.Limit)
	end := int(in.Page * in.Limit)
	if start >= total {
		start = 0
		end = 0
	}
	if end > total {
		end = total
	}

	pagedRecords := records
	if start < end {
		pagedRecords = records[start:end]
	} else {
		pagedRecords = make([]*model.RecommendRecord, 0)
	}

	// 9. 存储推荐结果（用于去重）
	avids := make([]int64, 0, len(pagedRecords))
	for _, record := range pagedRecords {
		avids = append(avids, record.AVID)
	}
	l.svcCtx.Dao.SaveRecommendRecord(l.ctx, in.Mid, avids)

	// 9.5 存储日志（特征日志和响应日志）
	l.storeLog(in, pagedRecords, userProfile)

	// 10. 转换为响应格式
	list := make([]*recommend.RecommendItem, 0, len(pagedRecords))
	for _, record := range pagedRecords {
		list = append(list, &recommend.RecommendItem{
			Avid:     record.AVID,
			Cid:      record.CID,
			Title:    record.Title,
			Cover:    record.Cover,
			Duration: record.Duration,
			PubTime:  record.PubTime,
			ZoneId:   record.ZoneID,
			ZoneName: record.ZoneName,
			UpMid:    record.UPMID,
			UpName:   record.UPName,
			Play:     record.Play,
			Like:     record.Like,
			Coin:     record.Coin,
			Fav:      record.Fav,
			Share:    record.Share,
			Score:    record.Score,
			Reason:   record.Reason,
			Tags:     record.Tags,
			Extra:    record.Extra,
		})
	}

	elapsed := time.Since(startTime).Milliseconds()
	logx.Infof("推荐请求完成: 耗时=%dms, 返回数量=%d", elapsed, len(list))

	debugInfo := make(map[string]string)
	if in.Debug {
		debugInfo["elapsed_ms"] = fmt.Sprintf("%d", elapsed)
		debugInfo["recall_count"] = fmt.Sprintf("%d", len(recallResp.List))
		debugInfo["filter_count"] = fmt.Sprintf("%d", len(records))
		debugInfo["final_count"] = fmt.Sprintf("%d", len(list))
	}

	return &recommend.RecommendResponse{
		List:      list,
		Total:     int32(total),
		HasMore:   end < total,
		DebugInfo: debugInfo,
		Message:   "success",
	}, nil
}

// buildRecallRequest 构建召回请求
func (l *GetRecommendListLogic) buildRecallRequest(req *recommend.RecommendRequest, profile *model.UserProfile) *recall.RecallRequest {
	infos := make([]*recall.RecallInfo, 0)

	// 构建会话特征（从点赞和正反馈视频中提取标签和UP主）
	likeVideoMap, likeUPMap, likeTagIDMap, posVideoMap, posTagIDMap := l.buildSessionFeature(profile)

	// 1. 精选召回（优先级最高）
	infos = append(infos, &recall.RecallInfo{
		Name:     "SelectionRecall",
		Tag:      "recall:selection",
		Limit:    100,
		Filter:   "bloomfilter",
		Priority: 10000,
		Scorer:   "default",
	})

	// 2. 热门召回
	infos = append(infos, &recall.RecallInfo{
		Name:     "HotRecall",
		Tag:      "recall:hot:default",
		Limit:    200,
		Filter:   "bloomfilter",
		Priority: 10,
		Scorer:   "default",
	})

	// 3. 点赞 I2I 召回（Top 10 点赞视频）
	count := 0
	for avid := range likeVideoMap {
		if count >= 10 {
			break
		}
		infos = append(infos, &recall.RecallInfo{
			Name:     "LikeI2IRecall",
			Tag:      fmt.Sprintf("RECALL:I2I:%d", avid),
			Limit:    40,
			Filter:   "bloomfilter",
			Priority: 10000,
			Scorer:   "default",
		})
		count++
	}

	// 4. 点赞 UP主召回（Top 10 UP主）
	count = 0
	for upMID := range likeUPMap {
		if count >= 10 {
			break
		}
		infos = append(infos, &recall.RecallInfo{
			Name:     "LikeUPRecall",
			Tag:      fmt.Sprintf("RECALL:HOT_UP:%d", upMID),
			Limit:    40,
			Filter:   "bloomfilter",
			Priority: 10000,
			Scorer:   "default",
		})
		count++
	}

	// 5. 关注UP主召回（Top 10）
	count = 0
	for upMID := range profile.FollowUps {
		if count >= 10 {
			break
		}
		infos = append(infos, &recall.RecallInfo{
			Name:     "FollowRecall",
			Tag:      fmt.Sprintf("RECALL:HOT_UP:%d", upMID),
			Limit:    40,
			Filter:   "bloomfilter",
			Priority: 1000,
			Scorer:   "default",
		})
		count++
	}

	// 6. 点赞标签召回（Top 10 标签）
	count = 0
	for tagID := range likeTagIDMap {
		if count >= 10 {
			break
		}
		infos = append(infos, &recall.RecallInfo{
			Name:     "LikeTagRecall",
			Tag:      fmt.Sprintf("RECALL:HOT_T:%d", tagID),
			Limit:    20,
			Filter:   "bloomfilter",
			Priority: 1000,
			Scorer:   "default",
		})
		count++
	}

	// 7. 正反馈 I2I 召回（Top 10 正反馈视频）
	count = 0
	for avid := range posVideoMap {
		if count >= 10 {
			break
		}
		infos = append(infos, &recall.RecallInfo{
			Name:     "PosI2IRecall",
			Tag:      fmt.Sprintf("RECALL:I2I:%d", avid),
			Limit:    20,
			Filter:   "bloomfilter",
			Priority: 100,
			Scorer:   "default",
		})
		count++
	}

	// 8. 正反馈标签召回（Top 10 标签）
	count = 0
	for tagID := range posTagIDMap {
		if count >= 10 {
			break
		}
		infos = append(infos, &recall.RecallInfo{
			Name:     "PosTagRecall",
			Tag:      fmt.Sprintf("RECALL:HOT_T:%d", tagID),
			Limit:    30,
			Filter:   "bloomfilter",
			Priority: 100,
			Scorer:   "default",
		})
		count++
	}

	// 9. 用户画像标签召回（历史标签）
	count = 0
	for tag := range profile.Tags {
		if count >= 10 {
			break
		}
		infos = append(infos, &recall.RecallInfo{
			Name:     "UserProfileBili",
			Tag:      fmt.Sprintf("RECALL:HOT_T:%s", tag),
			Limit:    20,
			Filter:   "bloomfilter",
			Priority: 100,
			Scorer:   "default",
		})
		count++
	}

	// 10. 用户画像分区召回
	count = 0
	for zoneID := range profile.Zones {
		if count >= 10 {
			break
		}
		infos = append(infos, &recall.RecallInfo{
			Name:     "UserProfileBili",
			Tag:      fmt.Sprintf("RECALL:HOT_T:%d", zoneID),
			Limit:    10,
			Filter:   "bloomfilter",
			Priority: 100,
			Scorer:   "default",
		})
		count++
	}

	// 11. 老用户随机召回（增加探索推荐）
	if len(profile.LastRecords) >= 20 && len(infos) < 50 {
		tagCountMap := make(map[string]int)
		for tag := range profile.Tags {
			tagCountMap[tag]++
		}

		randomTagCount := 0
		for tag := range tagCountMap {
			if randomTagCount >= 10 {
				break
			}
			randomTagCount++
			infos = append(infos, &recall.RecallInfo{
				Name:     "RandomRecall",
				Tag:      fmt.Sprintf("RECALL:T:%s", tag),
				Limit:    10,
				Filter:   "bloomfilter",
				Priority: 10,
				Scorer:   "default",
			})
		}
	}

	// 召回标签较多时，减少热门召回数量
	if len(infos) >= 5 {
		for _, info := range infos {
			if info.Name == "HotRecall" {
				info.Limit = 50
			}
		}
	}

	// 合并相同 Tag 的召回请求
	infos = l.mergeRecallKey(infos)

	return &recall.RecallRequest{
		Mid:        req.Mid,
		Buvid:      req.Buvid,
		TotalLimit: l.svcCtx.BusinessConfig.RecallLimit,
		Infos:      infos,
		TraceId:    req.TraceId,
	}
}

// mergeRecallKey 合并相同 Tag 的召回请求
// 目的：减少召回服务的重复查询，提高召回效率
func (l *GetRecommendListLogic) mergeRecallKey(recallInfos []*recall.RecallInfo) []*recall.RecallInfo {
	recallTagNameMap := make(map[string][]string)           // Tag -> []Name
	recallTagInfoMap := make(map[string]*recall.RecallInfo) // Tag -> Info
	recallTagPriorityMap := make(map[string]int32)          // Tag -> 最高优先级

	// 按 Tag 分组
	for _, recallInfo := range recallInfos {
		names := recallTagNameMap[recallInfo.Tag]
		names = append(names, recallInfo.Name)
		recallTagNameMap[recallInfo.Tag] = names
		recallTagInfoMap[recallInfo.Tag] = recallInfo

		// 保留最高优先级
		if priority, ok := recallTagPriorityMap[recallInfo.Tag]; ok {
			if recallInfo.Priority > priority {
				recallTagPriorityMap[recallInfo.Tag] = recallInfo.Priority
			}
		} else {
			recallTagPriorityMap[recallInfo.Tag] = recallInfo.Priority
		}
	}

	// 合并同 Tag 的召回请求
	newRecallInfos := make([]*recall.RecallInfo, 0, len(recallTagNameMap))
	for tag, names := range recallTagNameMap {
		recallInfo := recallTagInfoMap[tag]
		// 将多个召回策略名称用 | 连接
		recallInfo.Name = strings.Join(names, "|")
		// 使用最高优先级
		recallInfo.Priority = recallTagPriorityMap[tag]
		newRecallInfos = append(newRecallInfos, recallInfo)
	}

	logx.Infof("合并召回请求: 原数量=%d, 合并后=%d", len(recallInfos), len(newRecallInfos))
	return newRecallInfos
}

// buildSessionFeature 构建会话特征（从用户行为视频中提取标签和UP主）
// 返回值：likeVideoMap, likeUPMap, likeTagIDMap, posVideoMap, posTagIDMap
func (l *GetRecommendListLogic) buildSessionFeature(profile *model.UserProfile) (
	map[int64]int64, map[int64]int64, map[int64]int64, map[int64]int64, map[int64]int64) {

	// 收集所有需要查询的视频ID
	avids := make([]int64, 0)
	for avid := range profile.LikeVideos {
		avids = append(avids, avid)
	}
	for avid := range profile.PosVideos {
		avids = append(avids, avid)
	}
	for avid := range profile.NegVideos {
		avids = append(avids, avid)
	}

	if len(avids) == 0 {
		return make(map[int64]int64), make(map[int64]int64), make(map[int64]int64),
			make(map[int64]int64), make(map[int64]int64)
	}

	// 批量查询视频标签信息
	videoIndexReq := &recall.VideoIndexRequest{
		Avids: avids,
	}
	videoIndexResp, err := l.svcCtx.RecallRpc.VideoIndex(l.ctx, videoIndexReq)
	if err != nil || videoIndexResp == nil {
		logx.Errorf("查询视频索引失败: %v", err)
		return make(map[int64]int64), make(map[int64]int64), make(map[int64]int64),
			make(map[int64]int64), make(map[int64]int64)
	}

	// 初始化 TagIDs map（如果不存在）
	if profile.LikeTagIDs == nil {
		profile.LikeTagIDs = make(map[int64]int64)
	}
	if profile.PosTagIDs == nil {
		profile.PosTagIDs = make(map[int64]int64)
	}
	if profile.NegTagIDs == nil {
		profile.NegTagIDs = make(map[int64]int64)
	}
	if profile.LikeUPs == nil {
		profile.LikeUPs = make(map[int64]int64)
	}

	// 从视频索引中提取标签和UP主信息
	for _, videoIndex := range videoIndexResp.List {
		avid := videoIndex.BasicInfo.Avid

		// 处理点赞视频
		if timestamp, ok := profile.LikeVideos[avid]; ok {
			// 提取标签
			for _, tag := range videoIndex.BasicInfo.Tags {
				tagID := tag.TagId
				profile.LikeTagIDs[tagID]++
			}
			// 提取UP主
			upMID := videoIndex.BasicInfo.Mid
			profile.LikeUPs[upMID] = timestamp
		}

		// 处理正反馈视频
		if _, ok := profile.PosVideos[avid]; ok {
			for _, tag := range videoIndex.BasicInfo.Tags {
				tagID := tag.TagId
				profile.PosTagIDs[tagID]++
			}
		}

		// 处理负反馈视频
		if _, ok := profile.NegVideos[avid]; ok {
			for _, tag := range videoIndex.BasicInfo.Tags {
				tagID := tag.TagId
				profile.NegTagIDs[tagID]++
			}
		}
	}

	// 提取 Top N 的数据
	likeVideoMap := l.topNMap(profile.LikeVideos, 10)
	likeUPMap := l.topNMap(profile.LikeUPs, 10)
	likeTagIDMap := l.topNMap(profile.LikeTagIDs, 10)
	posVideoMap := l.topNMap(profile.PosVideos, 10)
	posTagIDMap := l.topNMap(profile.PosTagIDs, 10)

	return likeVideoMap, likeUPMap, likeTagIDMap, posVideoMap, posTagIDMap
}

// topNMap 提取 map 中前 N 个元素（按值排序）
func (l *GetRecommendListLogic) topNMap(m map[int64]int64, n int) map[int64]int64 {
	if len(m) <= n {
		return m
	}

	// 转换为切片并排序
	type pair struct {
		key   int64
		value int64
	}
	pairs := make([]pair, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, pair{k, v})
	}

	// 按值降序排序
	for i := 0; i < len(pairs)-1; i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[j].value > pairs[i].value {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}

	// 取 Top N
	result := make(map[int64]int64, n)
	for i := 0; i < n && i < len(pairs); i++ {
		result[pairs[i].key] = pairs[i].value
	}

	return result
}

// downGradeRecall 降级召回（放宽过滤条件）
func (l *GetRecommendListLogic) downGradeRecall(req *recommend.RecommendRequest, profile *model.UserProfile) (*recall.RecallResponse, error) {
	logx.Info("开始降级召回")

	// 只使用热门召回和精选召回，不使用布隆过滤器
	infos := []*recall.RecallInfo{
		{
			Name:     "SelectionRecall",
			Tag:      "recall:selection",
			Limit:    100,
			Filter:   "",
			Priority: 10000,
			Scorer:   "default",
		},
		{
			Name:     "HotRecall",
			Tag:      "recall:hot:default",
			Limit:    300,
			Filter:   "",
			Priority: 10,
			Scorer:   "default",
		},
	}

	recallReq := &recall.RecallRequest{
		Mid:        req.Mid,
		Buvid:      req.Buvid,
		TotalLimit: l.svcCtx.BusinessConfig.RecallLimit,
		Infos:      infos,
		TraceId:    req.TraceId,
	}

	return l.svcCtx.RecallRpc.Recall(l.ctx, recallReq)
}

// transformRecallToRecords 转换召回结果为推荐记录
func (l *GetRecommendListLogic) transformRecallToRecords(items []*recall.RecallItem) []*model.RecommendRecord {
	records := make([]*model.RecommendRecord, 0, len(items))

	for _, item := range items {
		record := &model.RecommendRecord{
			AVID:        item.Avid,
			Score:       item.Score,
			RecallTypes: item.RecallType,
			RecallTags:  item.RecallTag,
			Extra:       item.Extra,
		}
		records = append(records, record)
	}

	return records
}

// enrichVideoInfo 丰富视频信息
func (l *GetRecommendListLogic) enrichVideoInfo(records []*model.RecommendRecord) []*model.RecommendRecord {
	// 批量获取视频信息
	avids := make([]int64, 0, len(records))
	for _, record := range records {
		avids = append(avids, record.AVID)
	}

	videosMap, err := l.svcCtx.Dao.GetVideosInfo(l.ctx, avids)
	if err != nil {
		logx.Errorf("批量获取视频信息失败: %v", err)
		return records
	}

	enriched := make([]*model.RecommendRecord, 0, len(records))
	for _, record := range records {
		videoInfo, ok := videosMap[record.AVID]
		if !ok || videoInfo == nil {
			continue
		}

		// 合并信息
		record.Title = videoInfo.Title
		record.Cover = videoInfo.Cover
		record.Duration = videoInfo.Duration
		record.PubTime = videoInfo.PubTime
		record.ZoneID = videoInfo.ZoneID
		record.UPMID = videoInfo.UPMID
		record.State = videoInfo.State
		record.PlayHive = videoInfo.PlayHive
		record.LikesHive = videoInfo.LikesHive
		record.FavHive = videoInfo.FavHive
		record.ShareHive = videoInfo.ShareHive
		record.ReplyHive = videoInfo.ReplyHive
		record.CoinHive = videoInfo.CoinHive
		record.PlayMonth = videoInfo.PlayMonth
		record.LikesMonth = videoInfo.LikesMonth
		record.ShareMonth = videoInfo.ShareMonth
		record.ReplyMonth = videoInfo.ReplyMonth
		record.PlayMonthFinish = videoInfo.PlayMonthFinish
		record.Play = videoInfo.PlayMonth
		record.Like = videoInfo.LikesMonth
		// Tags 已经在 GetVideosInfo 中填充

		// 生成推荐理由
		record.Reason = l.generateReason(record)

		enriched = append(enriched, record)
	}

	return enriched
}

// generateReason 生成推荐理由
func (l *GetRecommendListLogic) generateReason(record *model.RecommendRecord) string {
	if strings.Contains(record.RecallTypes, "SelectionRecall") {
		return "精选推荐"
	}
	if strings.Contains(record.RecallTypes, "LikeI2IRecall") {
		return "根据你喜欢的视频推荐"
	}
	if strings.Contains(record.RecallTypes, "FollowRecall") {
		return "你关注的UP主发布的视频"
	}
	if strings.Contains(record.RecallTypes, "UserProfileRecall") {
		return "根据你的兴趣推荐"
	}
	if strings.Contains(record.RecallTypes, "HotRecall") {
		playStr := strconv.FormatInt(record.Play, 10)
		return fmt.Sprintf("热门视频 · %s播放", playStr)
	}
	return "为你推荐"
}

// storeLog 存储日志（特征日志和响应日志）
// 目的：1. 用于离线模型训练  2. 用于调试和分析推荐效果
func (l *GetRecommendListLogic) storeLog(req *recommend.RecommendRequest, records []*model.RecommendRecord, profile *model.UserProfile) {
	// 构建简化的记录（只保留必要字段，减少日志体积）
	simplifiedRecords := make([]map[string]interface{}, 0, len(records))
	for _, record := range records {
		simplifiedRecord := map[string]interface{}{
			"avid":         record.AVID,
			"cid":          record.CID,
			"score":        record.Score,
			"recall_types": record.RecallTypes,
			"recall_tags":  record.RecallTags,
			"state":        record.State,
			"up_mid":       record.UPMID,
			"zone_id":      record.ZoneID,
			"duration":     record.Duration,
			"play_month":   record.PlayMonth,
			"likes_month":  record.LikesMonth,
		}
		simplifiedRecords = append(simplifiedRecords, simplifiedRecord)
	}

	// 构建请求信息
	requestInfo := map[string]interface{}{
		"mid":      req.Mid,
		"buvid":    req.Buvid,
		"limit":    req.Limit,
		"page":     req.Page,
		"trace_id": req.TraceId,
	}

	// 构建响应信息
	responseInfo := map[string]interface{}{
		"total":     len(records),
		"records":   simplifiedRecords,
		"timestamp": time.Now().Unix(),
	}

	// 构建用户画像信息（用于特征分析）
	userInfo := map[string]interface{}{
		"mid":               profile.MID,
		"tags_count":        len(profile.Tags),
		"zones_count":       len(profile.Zones),
		"follow_ups_count":  len(profile.FollowUps),
		"black_ups_count":   len(profile.BlackUps),
		"like_videos_count": len(profile.LikeVideos),
		"pos_videos_count":  len(profile.PosVideos),
	}

	// 日志输出（结构化日志，便于日志收集系统采集）
	logx.WithContext(l.ctx).Infow("recommend_log",
		logx.Field("type", "recommend"),
		logx.Field("request", requestInfo),
		logx.Field("response", responseInfo),
		logx.Field("user_profile", userInfo),
	)

	// 如果是调试模式，输出更详细的信息
	if req.Debug {
		logx.WithContext(l.ctx).Infow("recommend_debug_log",
			logx.Field("type", "recommend_debug"),
			logx.Field("user_tags", profile.Tags),
			logx.Field("user_zones", profile.Zones),
			logx.Field("like_tag_ids", profile.LikeTagIDs),
			logx.Field("pos_tag_ids", profile.PosTagIDs),
		)
	}
}

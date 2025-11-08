// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package recommend

import (
	"context"

	"mybilibili/app/api/creative/internal/svc"
	"mybilibili/app/api/creative/internal/types"
	"mybilibili/app/recommend/cmd/rpc/recommend"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRecommendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取推荐列表
func NewGetRecommendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecommendListLogic {
	return &GetRecommendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRecommendListLogic) GetRecommendList(req *types.GetRecommendListReq) (resp *types.GetRecommendListResp, err error) {
	// 调用 recommend-rpc 服务
	rpcResp, err := l.svcCtx.RecommendRpc.GetRecommendList(l.ctx, &recommend.RecommendRequest{
		Mid:   req.Mid,
		Limit: req.Limit,
		Page:  req.Page,
		Debug: req.Debug,
	})
	if err != nil {
		l.Errorf("调用 recommend-rpc 失败: %v", err)
		return nil, err
	}

	// 转换响应
	items := make([]types.RecommendItem, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		items = append(items, types.RecommendItem{
			AVID:     item.Avid,
			Title:    item.Title,
			Cover:    item.Cover,
			Duration: item.Duration,
			PubTime:  item.PubTime,
			ZoneID:   item.ZoneId,
			ZoneName: getZoneName(item.ZoneId),
			UPMID:    item.UpMid,
			UPName:   "", // 需要从 user 服务获取
			Play:     item.Play,
			Like:     item.Like,
			Coin:     item.Coin,
			Fav:      item.Fav,
			Share:    item.Share,
			Score:    item.Score,
			Reason:   item.Reason,
			Tags:     item.Tags,
		})
	}

	resp = &types.GetRecommendListResp{
		List:      items,
		Total:     rpcResp.Total,
		HasMore:   rpcResp.HasMore,
		DebugInfo: rpcResp.DebugInfo,
	}

	l.Infof("推荐列表获取成功: mid=%d, count=%d", req.Mid, len(items))
	return resp, nil
}

// getZoneName 根据 zone_id 获取分区名称（简化版）
func getZoneName(zoneID int32) string {
	zoneMap := map[int32]string{
		1:   "动画",
		3:   "音乐",
		4:   "游戏",
		5:   "娱乐",
		11:  "电视剧",
		13:  "番剧",
		23:  "电影",
		36:  "科技",
		119: "鬼畜",
		129: "舞蹈",
		155: "时尚",
		160: "生活",
		167: "国创",
		168: "汽车",
		177: "纪录片",
		181: "影视",
		188: "数码",
		202: "资讯",
		203: "美食",
	}
	if name, ok := zoneMap[zoneID]; ok {
		return name
	}
	return "其他"
}

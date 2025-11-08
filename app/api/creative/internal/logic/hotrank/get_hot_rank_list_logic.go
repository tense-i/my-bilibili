// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package hotrank

import (
	"context"

	"mybilibili/app/api/creative/internal/svc"
	"mybilibili/app/api/creative/internal/types"
	"mybilibili/app/hotrank/cmd/rpc/hotrank_client"
	"mybilibili/app/video/cmd/rpc/video_client"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHotRankListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取全站热门排行榜
func NewGetHotRankListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHotRankListLogic {
	return &GetHotRankListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHotRankListLogic) GetHotRankList(req *types.GetHotRankListReq) (resp *types.GetHotRankListResp, err error) {
	// 1. 调用 hotrank-rpc 获取排行榜
	hotrankResp, err := l.svcCtx.HotrankRpc.GetHotRankList(l.ctx, &hotrank_client.GetHotRankListReq{
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		l.Errorf("GetHotRankList HotrankRpc.GetHotRankList error: %v", err)
		return nil, err
	}

	if len(hotrankResp.List) == 0 {
		return &types.GetHotRankListResp{
			List:  []types.HotRankItem{},
			Total: hotrankResp.Total,
		}, nil
	}

	// 2. 提取视频ID列表
	vids := make([]int64, len(hotrankResp.List))
	for i, item := range hotrankResp.List {
		vids[i] = item.Oid
	}

	// 3. 批量获取视频信息
	videoInfoResp, err := l.svcCtx.VideoRpc.BatchGetVideoInfo(l.ctx, &video_client.BatchGetVideoInfoReq{
		Vids: vids,
	})
	if err != nil {
		l.Errorf("GetHotRankList VideoRpc.BatchGetVideoInfo error: %v", err)
		return nil, err
	}

	// 4. 批量获取视频统计
	videoStatResp, err := l.svcCtx.VideoRpc.BatchGetVideoStat(l.ctx, &video_client.BatchGetVideoStatReq{
		Vids: vids,
	})
	if err != nil {
		l.Errorf("GetHotRankList VideoRpc.BatchGetVideoStat error: %v", err)
		return nil, err
	}

	// 5. 组合数据
	result := make([]types.HotRankItem, 0, len(hotrankResp.List))
	for _, item := range hotrankResp.List {
		rankItem := types.HotRankItem{
			Oid:  item.Oid,
			Hot:  item.Hot,
			Rank: item.Rank,
		}

		// 填充视频信息
		if info, ok := videoInfoResp.Infos[item.Oid]; ok && info != nil {
			rankItem.Title = info.Title
			rankItem.Cover = info.Cover
			rankItem.Author = info.AuthorName
		}

		// 填充统计数据
		if stat, ok := videoStatResp.Stats[item.Oid]; ok && stat != nil {
			rankItem.View = stat.View
			rankItem.Like = stat.Like
		}

		result = append(result, rankItem)
	}

	return &types.GetHotRankListResp{
		List:  result,
		Total: hotrankResp.Total,
	}, nil
}

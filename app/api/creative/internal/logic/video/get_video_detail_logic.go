// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package video

import (
	"context"

	"mybilibili/app/api/creative/internal/svc"
	"mybilibili/app/api/creative/internal/types"
	"mybilibili/app/hotrank/cmd/rpc/hotrank_client"
	"mybilibili/app/video/cmd/rpc/video_client"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取视频详情
func NewGetVideoDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoDetailLogic {
	return &GetVideoDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetVideoDetailLogic) GetVideoDetail(req *types.GetVideoDetailReq) (resp *types.GetVideoDetailResp, err error) {
	// 1. 获取视频信息
	videoInfoResp, err := l.svcCtx.VideoRpc.GetVideoInfo(l.ctx, &video_client.GetVideoInfoReq{
		Vid: req.Vid,
	})
	if err != nil {
		l.Errorf("GetVideoDetail VideoRpc.GetVideoInfo error: %v", err)
		return nil, err
	}

	// 2. 获取视频统计
	videoStatResp, err := l.svcCtx.VideoRpc.GetVideoStat(l.ctx, &video_client.GetVideoStatReq{
		Vid: req.Vid,
	})
	if err != nil {
		l.Errorf("GetVideoDetail VideoRpc.GetVideoStat error: %v", err)
		return nil, err
	}

	// 3. 获取热度值
	hotResp, err := l.svcCtx.HotrankRpc.GetHotByOID(l.ctx, &hotrank_client.GetHotByOIDReq{
		Oid:      req.Vid,
		Business: 1, // 视频
	})
	if err != nil {
		l.Errorf("GetVideoDetail HotrankRpc.GetHotByOID error: %v", err)
		// 热度值获取失败不影响主流程，设为0
		hotResp = &hotrank_client.GetHotByOIDResp{Hot: 0}
	}

	// 4. 组合数据
	return &types.GetVideoDetailResp{
		Video: types.VideoDetail{
			Vid:        videoInfoResp.Info.Vid,
			Title:      videoInfoResp.Info.Title,
			Cover:      videoInfoResp.Info.Cover,
			AuthorId:   videoInfoResp.Info.AuthorId,
			AuthorName: videoInfoResp.Info.AuthorName,
			RegionId:   videoInfoResp.Info.RegionId,
			PubTime:    videoInfoResp.Info.PubTime,
			Duration:   videoInfoResp.Info.Duration,
			Desc:       videoInfoResp.Info.Desc,
			View:       videoStatResp.Stat.View,
			Like:       videoStatResp.Stat.Like,
			Coin:       videoStatResp.Stat.Coin,
			Fav:        videoStatResp.Stat.Fav,
			Share:      videoStatResp.Stat.Share,
			Reply:      videoStatResp.Stat.Reply,
			Danmaku:    videoStatResp.Stat.Danmaku,
			Hot:        hotResp.Hot,
		},
	}, nil
}

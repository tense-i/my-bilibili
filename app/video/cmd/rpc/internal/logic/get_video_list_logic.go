package logic

import (
	"context"

	"mybilibili/app/video/cmd/rpc/internal/svc"
	"mybilibili/app/video/cmd/rpc/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoListLogic {
	return &GetVideoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取视频列表（游标分页）
// 获取视频列表（游标分页）
func (l *GetVideoListLogic) GetVideoList(in *video.GetVideoListReq) (*video.GetVideoListResp, error) {
	// 查询视频列表
	infos, err := l.svcCtx.VideoInfoModel.FindListByLastVid(l.ctx, in.LastVid, int(in.Limit))
	if err != nil {
		l.Errorf("GetVideoList FindListByLastVid error: %v", err)
		return nil, err
	}

	if len(infos) == 0 {
		return &video.GetVideoListResp{List: []*video.VideoData{}}, nil
	}

	// 提取 vid 列表
	vids := make([]int64, len(infos))
	for i, info := range infos {
		vids[i] = info.Vid
	}

	// 批量查询统计数据
	stats, err := l.svcCtx.VideoStatModel.FindByVids(l.ctx, vids)
	if err != nil {
		l.Errorf("GetVideoList FindByVids stats error: %v", err)
		return nil, err
	}

	// 构建 vid -> stat 映射
	statMap := make(map[int64]*video.VideoStat)
	for _, stat := range stats {
		statMap[stat.Vid] = &video.VideoStat{
			Vid:     stat.Vid,
			View:    stat.View,
			Like:    stat.LikeCount,
			Coin:    stat.Coin,
			Fav:     stat.Fav,
			Share:   stat.Share,
			Reply:   stat.Reply,
			Danmaku: stat.Danmaku,
		}
	}

	// 组合结果
	result := make([]*video.VideoData, 0, len(infos))
	for _, info := range infos {
		videoData := &video.VideoData{
			Info: &video.VideoInfo{
				Vid:        info.Vid,
				Title:      info.Title,
				Cover:      info.Cover,
				AuthorId:   info.AuthorId,
				AuthorName: info.AuthorName,
				RegionId:   int64(info.RegionId),
				PubTime:    info.PubTime,
				Duration:   int32(info.Duration),
				Desc:       info.Desc.String,
				State:      int32(info.State),
			},
			Stat: statMap[info.Vid],
		}
		result = append(result, videoData)
	}

	return &video.GetVideoListResp{List: result}, nil
}

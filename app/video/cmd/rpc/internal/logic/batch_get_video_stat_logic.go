package logic

import (
	"context"

	"mybilibili/app/video/cmd/rpc/internal/svc"
	"mybilibili/app/video/cmd/rpc/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetVideoStatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetVideoStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetVideoStatLogic {
	return &BatchGetVideoStatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量获取视频统计数据（⭐hotrank-job 会调用此接口进行热度计算）
func (l *BatchGetVideoStatLogic) BatchGetVideoStat(in *video.BatchGetVideoStatReq) (*video.BatchGetVideoStatResp, error) {
	// 参数校验
	if len(in.Vids) == 0 {
		return &video.BatchGetVideoStatResp{Stats: make(map[int64]*video.VideoStat)}, nil
	}

	// 批量查询统计数据
	stats, err := l.svcCtx.VideoStatModel.FindByVids(l.ctx, in.Vids)
	if err != nil {
		l.Errorf("BatchGetVideoStat FindByVids error: %v", err)
		return nil, err
	}

	// 转换为 proto 结构
	result := make(map[int64]*video.VideoStat)
	for _, stat := range stats {
		result[stat.Vid] = &video.VideoStat{
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

	return &video.BatchGetVideoStatResp{Stats: result}, nil
}

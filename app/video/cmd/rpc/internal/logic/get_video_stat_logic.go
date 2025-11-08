package logic

import (
	"context"

	"mybilibili/app/video/cmd/rpc/internal/svc"
	"mybilibili/app/video/cmd/rpc/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoStatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoStatLogic {
	return &GetVideoStatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取视频统计数据
func (l *GetVideoStatLogic) GetVideoStat(in *video.GetVideoStatReq) (*video.GetVideoStatResp, error) {
	// 查询统计数据
	stat, err := l.svcCtx.VideoStatModel.FindOneByVid(l.ctx, in.Vid)
	if err != nil {
		l.Errorf("GetVideoStat FindOneByVid error: %v", err)
		return nil, err
	}

	// 转换为 proto 结构
	return &video.GetVideoStatResp{
		Stat: &video.VideoStat{
			Vid:     stat.Vid,
			View:    stat.View,
			Like:    stat.LikeCount,
			Coin:    stat.Coin,
			Fav:     stat.Fav,
			Share:   stat.Share,
			Reply:   stat.Reply,
			Danmaku: stat.Danmaku,
		},
	}, nil
}

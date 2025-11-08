package logic

import (
	"context"

	"mybilibili/app/video/cmd/rpc/internal/svc"
	"mybilibili/app/video/cmd/rpc/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoInfoLogic {
	return &GetVideoInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取视频信息
func (l *GetVideoInfoLogic) GetVideoInfo(in *video.GetVideoInfoReq) (*video.GetVideoInfoResp, error) {
	// 查询视频信息
	info, err := l.svcCtx.VideoInfoModel.FindOneByVid(l.ctx, in.Vid)
	if err != nil {
		l.Errorf("GetVideoInfo FindOneByVid error: %v", err)
		return nil, err
	}

	// 转换为 proto 结构
	return &video.GetVideoInfoResp{
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
	}, nil
}

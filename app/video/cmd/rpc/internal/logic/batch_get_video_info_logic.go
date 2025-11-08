package logic

import (
	"context"

	"mybilibili/app/video/cmd/rpc/internal/svc"
	"mybilibili/app/video/cmd/rpc/video"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetVideoInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetVideoInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetVideoInfoLogic {
	return &BatchGetVideoInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量获取视频信息（⭐hotrank-job 会调用此接口）
func (l *BatchGetVideoInfoLogic) BatchGetVideoInfo(in *video.BatchGetVideoInfoReq) (*video.BatchGetVideoInfoResp, error) {
	// 参数校验
	if len(in.Vids) == 0 {
		return &video.BatchGetVideoInfoResp{Infos: make(map[int64]*video.VideoInfo)}, nil
	}

	// 批量查询视频信息
	infos, err := l.svcCtx.VideoInfoModel.FindByVids(l.ctx, in.Vids)
	if err != nil {
		l.Errorf("BatchGetVideoInfo FindByVids error: %v", err)
		return nil, err
	}

	// 转换为 proto 结构
	result := make(map[int64]*video.VideoInfo)
	for _, info := range infos {
		result[info.Vid] = &video.VideoInfo{
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
		}
	}

	return &video.BatchGetVideoInfoResp{Infos: result}, nil
}

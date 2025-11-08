package logic

import (
	"context"

	"mybilibili/app/recall/cmd/rpc/internal/svc"
	"mybilibili/app/recall/cmd/rpc/recall"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVideoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoIndexLogic {
	return &VideoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// VideoIndex 视频索引查询
func (l *VideoIndexLogic) VideoIndex(in *recall.VideoIndexRequest) (*recall.VideoIndexResponse, error) {
	logx.Infof("视频索引查询: avids_count=%d", len(in.Avids))

	// 从数据库查询视频基础信息
	videos, err := l.svcCtx.Dao.GetVideosBasicInfo(l.ctx, in.Avids)
	if err != nil {
		logx.Errorf("查询视频索引失败: %v", err)
		return nil, err
	}

	// 转换为 proto 格式
	list := make([]*recall.VideoIndex, 0, len(videos))
	for _, video := range videos {
		// 转换标签
		tags := make([]*recall.Tag, 0, len(video.Tags))
		for _, tag := range video.Tags {
			tags = append(tags, &recall.Tag{
				TagId:   tag.TagID,
				TagName: tag.TagName,
			})
		}

		list = append(list, &recall.VideoIndex{
			Avid: video.AVID,
			BasicInfo: &recall.BasicInfo{
				Mid:             video.MID,
				Title:           video.Title,
				ZoneId:          video.ZoneID,
				Duration:        video.Duration,
				PubTime:         video.PubTime,
				State:           int32(video.State),
				Tags:            tags,
				PlayHive:        video.PlayHive,
				LikesHive:       video.LikesHive,
				FavHive:         video.FavHive,
				ReplyHive:       video.ReplyHive,
				ShareHive:       video.ShareHive,
				CoinHive:        video.CoinHive,
				PlayMonth:       video.PlayMonth,
				LikesMonth:      video.LikesMonth,
				ReplyMonth:      video.ReplyMonth,
				ShareMonth:      video.ShareMonth,
				PlayMonthFinish: video.PlayMonthFinish,
			},
		})
	}

	logx.Infof("视频索引查询完成: 返回数量=%d", len(list))

	return &recall.VideoIndexResponse{
		List: list,
	}, nil
}

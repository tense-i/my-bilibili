package dao

import (
	"context"

	pb "mybilibili/app/video/cmd/rpc/video"
	"mybilibili/app/video/cmd/rpc/video_client"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoDao struct {
	videoRpc video_client.Video
}

func NewVideoDao(videoRpc video_client.Video) *VideoDao {
	return &VideoDao{
		videoRpc: videoRpc,
	}
}

// Archives 批量获取视频元数据（参考主项目）
func (d *VideoDao) Archives(ctx context.Context, vids []int64) (map[int64]*pb.VideoInfo, error) {
	resp, err := d.videoRpc.BatchGetVideoInfo(ctx, &pb.BatchGetVideoInfoReq{
		Vids: vids,
	})
	if err != nil {
		logx.Errorf("VideoDao.Archives BatchGetVideoInfo error(%v)", err)
		return nil, err
	}

	return resp.Infos, nil
}

// Stats 批量获取视频统计数据（参考主项目）
func (d *VideoDao) Stats(ctx context.Context, vids []int64) (map[int64]*pb.VideoStat, error) {
	resp, err := d.videoRpc.BatchGetVideoStat(ctx, &pb.BatchGetVideoStatReq{
		Vids: vids,
	})
	if err != nil {
		logx.Errorf("VideoDao.Stats BatchGetVideoStat error(%v)", err)
		return nil, err
	}

	return resp.Stats, nil
}

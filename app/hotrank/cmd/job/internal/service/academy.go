package service

import (
	"context"
	"math"
	"time"

	pb "mybilibili/app/video/cmd/rpc/video"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	BusinessForVideo   = 1
	BusinessForArticle = 2
)

// FlushHot 计算热度值（参考主项目的 FlushHot）
func (s *Service) FlushHot(bs int) {
	defer s.wg.Done()

	var (
		ctx   = context.Background()
		id    int64
		limit = 30 // 每批30条
	)

	logx.Infof("FlushHot started for business: %d", bs)

	for {
		// 1. 游标分页查询（参考主项目）
		arcs, err := s.aca.Archives(ctx, id, bs, limit)
		if err != nil {
			logx.Errorf("s.aca.Archives id(%d) error(%v)", id, err)
			time.Sleep(time.Second * 10)
			continue
		}

		count := len(arcs)
		if count == 0 {
			id = 0 // 重新从头开始（参考主项目）
			logx.Info("FlushHot: no more records, sleep 1 hour and restart from beginning")
			time.Sleep(time.Hour * 1)
			continue
		}

		// 2. 提取OID列表
		oids := make([]int64, 0)
		for _, a := range arcs {
			oids = append(oids, a.OID)
			id = a.ID // 更新游标（参考主项目）
		}

		// 3. 批量计算热度（参考主项目）
		hots, err := s.computeHotByOIDs(ctx, oids, bs)
		if err != nil {
			logx.Errorf("s.computeHotByOIDs error(%v)", err)
			time.Sleep(time.Second * 10)
			continue
		}

		// 4. 批量更新数据库（参考主项目）
		if err := s.aca.UPHotByAIDs(ctx, hots); err != nil {
			logx.Errorf("s.aca.UPHotByAIDs hots(%+v) error(%v)", hots, err)
			time.Sleep(time.Second * 10)
			continue
		}

		logx.Infof("FlushHot success: processed %d videos, last_id=%d", count, id)

		// 休息一下，避免频繁查询
		time.Sleep(time.Second * 5)
	}
}

// computeHotByOIDs 批量计算热度（参考主项目）
func (s *Service) computeHotByOIDs(ctx context.Context, oids []int64, bs int) (map[int64]int64, error) {
	res := make(map[int64]int64)

	if bs == BusinessForVideo {
		// 批量获取视频元数据（参考主项目）
		arcs, err := s.arc.Archives(ctx, oids)
		if err != nil {
			return nil, err
		}

		// 批量获取统计数据（参考主项目）
		stats, err := s.arc.Stats(ctx, oids)
		if err != nil {
			logx.Errorf("s.arc.Stats oids(%+v) error(%v)", oids, err)
			return nil, err
		}

		// 计算每个视频的热度（参考主项目）
		for _, oid := range oids {
			if arc, ok := arcs[oid]; ok && arc != nil {
				if stat, ok := stats[oid]; ok && stat != nil {
					res[oid] = countArcHot(stat, arc.PubTime)
				}
			}
		}
	}

	return res, nil
}

// countArcHot 热度计算公式（完全参考主项目）
// 公式：硬币×0.4 + 收藏×0.3 + 弹幕×0.4 + 评论×0.4 + 播放×0.25 + 点赞×0.4 + 分享×0.6
// 新视频提权：24小时内发布的视频热度×1.5
func countArcHot(stat *pb.VideoStat, ptime int64) int64 {
	if stat == nil {
		return 0
	}

	// 多维度加权计算（完全参考主项目）
	hot := float64(stat.Coin)*0.4 +
		float64(stat.Fav)*0.3 +
		float64(stat.Danmaku)*0.4 +
		float64(stat.Reply)*0.4 +
		float64(stat.View)*0.25 +
		float64(stat.Like)*0.4 +
		float64(stat.Share)*0.6

	// 新视频提权（完全参考主项目）
	if ptime >= time.Now().AddDate(0, 0, -1).Unix() && ptime <= time.Now().Unix() {
		hot *= 1.5
	}

	return int64(math.Floor(hot))
}

package service

import (
	"sync"

	"mybilibili/app/hotrank/cmd/job/internal/config"
	"mybilibili/app/hotrank/cmd/job/internal/dao"
	"mybilibili/app/video/cmd/rpc/video_client"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// Service 主服务（参考主项目）
type Service struct {
	c   config.Config
	arc *dao.VideoDao   // RPC客户端（调用video服务）
	aca *dao.AcademyDao // 数据库DAO（热度表）
	wg  sync.WaitGroup
}

// New 初始化服务（参考主项目）
func New(c config.Config, videoRpc video_client.Video) *Service {
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	s := &Service{
		c:   c,
		arc: dao.NewVideoDao(videoRpc),
		aca: dao.NewAcademyDao(conn),
	}

	// 启动热度计算任务（参考主项目）
	if c.HotSwitch {
		s.wg.Add(1)
		go s.FlushHot(BusinessForVideo) // 计算视频热度
		// go s.FlushHot(BusinessForArticle)  // 计算专栏热度（暂不实现）
	}

	return s
}

// Close 关闭服务
func (s *Service) Close() {
	s.wg.Wait()
}

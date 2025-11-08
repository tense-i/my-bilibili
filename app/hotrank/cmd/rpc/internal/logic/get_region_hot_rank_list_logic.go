package logic

import (
	"context"

	"mybilibili/app/hotrank/cmd/rpc/hotrank"
	"mybilibili/app/hotrank/cmd/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRegionHotRankListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRegionHotRankListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRegionHotRankListLogic {
	return &GetRegionHotRankListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取分区热门排行榜
// 获取分区热门排行榜
func (l *GetRegionHotRankListLogic) GetRegionHotRankList(in *hotrank.GetRegionHotRankListReq) (*hotrank.GetRegionHotRankListResp, error) {
	// 默认值处理
	if in.Limit <= 0 {
		in.Limit = 50
	}
	if in.Limit > 100 {
		in.Limit = 100
	}

	// 查询分区排行榜数据（业务类型1表示视频）
	list, err := l.svcCtx.AcademyArchiveModel.FindRegionHotRankList(l.ctx, in.RegionId, 1, in.Offset, in.Limit)
	if err != nil {
		l.Errorf("GetRegionHotRankList FindRegionHotRankList error: %v", err)
		return nil, err
	}

	// 查询总数
	total, err := l.svcCtx.AcademyArchiveModel.CountRegionHotRank(l.ctx, in.RegionId, 1)
	if err != nil {
		l.Errorf("GetRegionHotRankList CountRegionHotRank error: %v", err)
		return nil, err
	}

	// 转换为 proto 结构
	result := make([]*hotrank.HotRankItem, 0, len(list))
	for i, item := range list {
		result = append(result, &hotrank.HotRankItem{
			Oid:      item.Oid,
			Business: int32(item.Business),
			Hot:      item.Hot,
			Rank:     in.Offset + int64(i) + 1, // 计算排名
		})
	}

	return &hotrank.GetRegionHotRankListResp{
		List:  result,
		Total: total,
	}, nil
}

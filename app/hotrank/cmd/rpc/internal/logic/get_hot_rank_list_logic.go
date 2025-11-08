package logic

import (
	"context"

	"mybilibili/app/hotrank/cmd/rpc/hotrank"
	"mybilibili/app/hotrank/cmd/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHotRankListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetHotRankListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHotRankListLogic {
	return &GetHotRankListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取全站热门排行榜
func (l *GetHotRankListLogic) GetHotRankList(in *hotrank.GetHotRankListReq) (*hotrank.GetHotRankListResp, error) {
	// 默认值处理
	if in.Limit <= 0 {
		in.Limit = 50
	}
	if in.Limit > 100 {
		in.Limit = 100
	}

	// 查询排行榜数据（业务类型1表示视频）
	list, err := l.svcCtx.AcademyArchiveModel.FindHotRankList(l.ctx, 1, in.Offset, in.Limit)
	if err != nil {
		l.Errorf("GetHotRankList FindHotRankList error: %v", err)
		return nil, err
	}

	// 查询总数
	total, err := l.svcCtx.AcademyArchiveModel.CountHotRank(l.ctx, 1)
	if err != nil {
		l.Errorf("GetHotRankList CountHotRank error: %v", err)
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

	return &hotrank.GetHotRankListResp{
		List:  result,
		Total: total,
	}, nil
}

package search

import (
	"context"

	"mybilibili/app/search/cmd/api/internal/svc"
	"mybilibili/app/search/cmd/api/internal/types"
	"mybilibili/app/search/cmd/rpc/search"

	"github.com/zeromicro/go-zero/core/logx"
)

type PgcUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPgcUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PgcUpdateLogic {
	return &PgcUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PgcUpdateLogic) PgcUpdate(req *types.PgcUpdateReq) (*types.UpdateResp, error) {
	// 转换请求
	items := make([]*search.PgcUpdateItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = &search.PgcUpdateItem{
			MediaId: item.MediaId,
			Field:   item.Field,
		}
	}

	// 调用 RPC
	resp, err := l.svcCtx.SearchRpc.PgcUpdate(l.ctx, &search.PgcUpdateReq{
		Items: items,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateResp{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}

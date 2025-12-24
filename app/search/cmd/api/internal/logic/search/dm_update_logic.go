package search

import (
	"context"

	"mybilibili/app/search/cmd/api/internal/svc"
	"mybilibili/app/search/cmd/api/internal/types"
	"mybilibili/app/search/cmd/rpc/search"

	"github.com/zeromicro/go-zero/core/logx"
)

type DmUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDmUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DmUpdateLogic {
	return &DmUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DmUpdateLogic) DmUpdate(req *types.DmUpdateReq) (*types.UpdateResp, error) {
	// 转换请求
	items := make([]*search.DmUpdateItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = &search.DmUpdateItem{
			Id:    item.Id,
			Oid:   item.Oid,
			Field: item.Field,
		}
	}

	// 调用 RPC
	resp, err := l.svcCtx.SearchRpc.DmUpdate(l.ctx, &search.DmUpdateReq{
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

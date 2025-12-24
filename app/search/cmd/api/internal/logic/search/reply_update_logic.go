package search

import (
	"context"

	"mybilibili/app/search/cmd/api/internal/svc"
	"mybilibili/app/search/cmd/api/internal/types"
	"mybilibili/app/search/cmd/rpc/search"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReplyUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyUpdateLogic {
	return &ReplyUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReplyUpdateLogic) ReplyUpdate(req *types.ReplyUpdateReq) (*types.UpdateResp, error) {
	// 转换请求
	items := make([]*search.ReplyUpdateItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = &search.ReplyUpdateItem{
			Id:    item.Id,
			Oid:   item.Oid,
			Mid:   item.Mid,
			State: item.State,
		}
	}

	// 调用 RPC
	resp, err := l.svcCtx.SearchRpc.ReplyUpdate(l.ctx, &search.ReplyUpdateReq{
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

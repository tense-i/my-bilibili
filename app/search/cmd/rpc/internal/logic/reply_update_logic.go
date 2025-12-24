package logic

import (
	"context"
	"fmt"

	"mybilibili/app/search/cmd/rpc/internal/svc"
	"mybilibili/app/search/cmd/rpc/search"
	"mybilibili/app/search/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyUpdateLogic {
	return &ReplyUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReplyUpdateLogic) ReplyUpdate(in *search.ReplyUpdateReq) (*search.UpdateResp, error) {
	if len(in.Items) == 0 {
		return &search.UpdateResp{
			Success: true,
			Message: "no items to update",
		}, nil
	}

	// 构建批量更新项
	bulkItems := make([]model.BulkUpdateItem, 0, len(in.Items))
	for _, item := range in.Items {
		indexName := fmt.Sprintf("replyrecord_%d", item.Mid%100)
		indexID := fmt.Sprintf("%d_%d", item.Id, item.Oid)

		fields := map[string]interface{}{
			"state": item.State,
		}

		bulkItems = append(bulkItems, model.BulkUpdateItem{
			IndexName: indexName,
			IndexID:   indexID,
			Fields:    fields,
		})
	}

	// 执行批量更新
	err := l.svcCtx.ESClient.BulkUpdate(l.ctx, "replyExternal", bulkItems)
	if err != nil {
		l.Errorf("ReplyUpdate failed: %v", err)
		return &search.UpdateResp{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &search.UpdateResp{
		Success: true,
		Message: fmt.Sprintf("updated %d items", len(in.Items)),
	}, nil
}

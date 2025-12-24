package logic

import (
	"context"
	"fmt"

	"mybilibili/app/search/cmd/rpc/internal/svc"
	"mybilibili/app/search/cmd/rpc/search"
	"mybilibili/app/search/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type DmUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDmUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DmUpdateLogic {
	return &DmUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DmUpdateLogic) DmUpdate(in *search.DmUpdateReq) (*search.UpdateResp, error) {
	if len(in.Items) == 0 {
		return &search.UpdateResp{
			Success: true,
			Message: "no items to update",
		}, nil
	}

	// 构建批量更新项
	bulkItems := make([]model.BulkUpdateItem, 0, len(in.Items))
	for _, item := range in.Items {
		indexName := fmt.Sprintf("dm_search_%d", item.Oid%1000)
		indexID := fmt.Sprintf("%d", item.Id)

		fields := make(map[string]interface{})
		for k, v := range item.Field {
			fields[k] = v
		}

		bulkItems = append(bulkItems, model.BulkUpdateItem{
			IndexName: indexName,
			IndexID:   indexID,
			Fields:    fields,
		})
	}

	// 执行批量更新
	err := l.svcCtx.ESClient.BulkUpdate(l.ctx, "dmExternal", bulkItems)
	if err != nil {
		l.Errorf("DmUpdate failed: %v", err)
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

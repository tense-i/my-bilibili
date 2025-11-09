// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package wallet

import (
	"context"

	"mybilibili/app/wallet/cmd/api/internal/svc"
	"mybilibili/app/wallet/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStreamListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询流水
func NewGetStreamListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStreamListLogic {
	return &GetStreamListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetStreamListLogic) GetStreamList(req *types.GetStreamListReq) (resp *types.GetStreamListResp, err error) {
	// todo: add your logic here and delete this line

	return
}

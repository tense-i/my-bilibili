// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package wallet

import (
	"context"

	"mybilibili/app/wallet/cmd/api/internal/svc"
	"mybilibili/app/wallet/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询余额
func NewGetDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDetailLogic {
	return &GetDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDetailLogic) GetDetail(req *types.GetDetailReq) (resp *types.GetDetailResp, err error) {
	// todo: add your logic here and delete this line

	return
}

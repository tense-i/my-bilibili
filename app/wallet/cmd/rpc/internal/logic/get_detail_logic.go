package logic

import (
	"context"

	"mybilibili/app/wallet/cmd/rpc/internal/svc"
	"mybilibili/app/wallet/cmd/rpc/wallet"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDetailLogic {
	return &GetDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询余额详情
func (l *GetDetailLogic) GetDetail(in *wallet.GetDetailReq) (*wallet.GetDetailResp, error) {
	// todo: add your logic here and delete this line

	return &wallet.GetDetailResp{}, nil
}

package logic

import (
	"context"

	"mybilibili/app/wallet/cmd/rpc/internal/svc"
	"mybilibili/app/wallet/cmd/rpc/wallet"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStreamListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStreamListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStreamListLogic {
	return &GetStreamListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询流水列表
func (l *GetStreamListLogic) GetStreamList(in *wallet.GetStreamListReq) (*wallet.GetStreamListResp, error) {
	// todo: add your logic here and delete this line

	return &wallet.GetStreamListResp{}, nil
}

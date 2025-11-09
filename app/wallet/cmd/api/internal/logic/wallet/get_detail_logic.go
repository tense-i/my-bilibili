// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package wallet

import (
	"context"

	"mybilibili/app/wallet/cmd/api/internal/svc"
	"mybilibili/app/wallet/cmd/api/internal/types"
	walletpb "mybilibili/app/wallet/cmd/rpc/wallet"

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
	// 调用RPC服务
	rpcResp, err := l.svcCtx.WalletRpc.GetDetail(l.ctx, &walletpb.GetDetailReq{
		Uid: req.Uid,
	})
	if err != nil {
		l.Errorf("get detail rpc failed: uid=%d, err=%v", req.Uid, err)
		return nil, err
	}

	if rpcResp.Detail == nil {
		return &types.GetDetailResp{}, nil
	}

	return &types.GetDetailResp{
		Detail: types.WalletDetail{
			Uid:             rpcResp.Detail.Uid,
			Gold:            rpcResp.Detail.Gold,
			IapGold:         rpcResp.Detail.IapGold,
			Silver:          rpcResp.Detail.Silver,
			GoldRechargeCnt: rpcResp.Detail.GoldRechargeCnt,
			GoldPayCnt:      rpcResp.Detail.GoldPayCnt,
			SilverPayCnt:    rpcResp.Detail.SilverPayCnt,
			CostBase:        rpcResp.Detail.CostBase,
		},
	}, nil
}

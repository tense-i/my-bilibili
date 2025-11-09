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

type ExchangeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 兑换
func NewExchangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeLogic {
	return &ExchangeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExchangeLogic) Exchange(req *types.ExchangeReq) (resp *types.WalletResp, err error) {
	// 调用RPC服务
	rpcResp, err := l.svcCtx.WalletRpc.Exchange(l.ctx, &walletpb.ExchangeReq{
		Uid:           req.Uid,
		SrcCoinType:   req.SrcCoinType,
		SrcCoinNum:    req.SrcCoinNum,
		DestCoinType:  req.DestCoinType,
		DestCoinNum:   req.DestCoinNum,
		TransactionId: req.TransactionId,
		Platform:      req.Platform,
	})
	if err != nil {
		l.Errorf("exchange rpc failed: uid=%d, err=%v", req.Uid, err)
		return nil, err
	}

	return &types.WalletResp{
		Gold:    rpcResp.Gold,
		IapGold: rpcResp.IapGold,
		Silver:  rpcResp.Silver,
	}, nil
}

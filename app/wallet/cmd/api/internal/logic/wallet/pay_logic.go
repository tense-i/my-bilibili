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

type PayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 消费
func NewPayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayLogic {
	return &PayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayLogic) Pay(req *types.PayReq) (resp *types.WalletResp, err error) {
	// 调用RPC服务
	rpcResp, err := l.svcCtx.WalletRpc.Pay(l.ctx, &walletpb.PayReq{
		Uid:           req.Uid,
		CoinType:      req.CoinType,
		CoinNum:       req.CoinNum,
		TransactionId: req.TransactionId,
		ExtendTid:     req.ExtendTid,
		Platform:      req.Platform,
		Timestamp:     req.Timestamp,
		BizCode:       req.BizCode,
	})
	if err != nil {
		l.Errorf("pay rpc failed: uid=%d, err=%v", req.Uid, err)
		return nil, err
	}

	return &types.WalletResp{
		Gold:    rpcResp.Gold,
		IapGold: rpcResp.IapGold,
		Silver:  rpcResp.Silver,
	}, nil
}

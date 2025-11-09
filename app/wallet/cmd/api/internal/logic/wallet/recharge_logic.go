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

type RechargeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 充值
func NewRechargeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RechargeLogic {
	return &RechargeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RechargeLogic) Recharge(req *types.RechargeReq) (resp *types.WalletResp, err error) {
	// 调用RPC服务
	rpcResp, err := l.svcCtx.WalletRpc.Recharge(l.ctx, &walletpb.RechargeReq{
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
		l.Errorf("recharge rpc failed: uid=%d, err=%v", req.Uid, err)
		return nil, err
	}

	return &types.WalletResp{
		Gold:    rpcResp.Gold,
		IapGold: rpcResp.IapGold,
		Silver:  rpcResp.Silver,
	}, nil
}

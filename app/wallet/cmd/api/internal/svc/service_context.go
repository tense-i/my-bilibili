// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"mybilibili/app/wallet/cmd/api/internal/config"
	"mybilibili/app/wallet/cmd/rpc/wallet_client"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	// Wallet RPC客户端
	WalletRpc wallet_client.Wallet
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		WalletRpc: wallet_client.NewWallet(zrpc.MustNewClient(c.WalletRpc)),
	}
}

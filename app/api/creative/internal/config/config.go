// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	// Video RPC 客户端配置
	VideoRpc zrpc.RpcClientConf

	// Hotrank RPC 客户端配置
	HotrankRpc zrpc.RpcClientConf

	// Recommend RPC 客户端配置
	RecommendRpc zrpc.RpcClientConf
}

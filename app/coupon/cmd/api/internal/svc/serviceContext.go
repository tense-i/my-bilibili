// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"mybilibili/app/coupon/cmd/api/internal/config"
	"mybilibili/app/coupon/cmd/rpc/coupon_client"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	CouponRpc coupon_client.Coupon
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		CouponRpc: coupon_client.NewCoupon(zrpc.MustNewClient(c.CouponRpc)),
	}
}

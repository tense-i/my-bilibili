package main

import (
	"flag"
	"fmt"

	"mybilibili/app/recommend/cmd/rpc/internal/config"
	"mybilibili/app/recommend/cmd/rpc/internal/server"
	"mybilibili/app/recommend/cmd/rpc/internal/svc"
	"mybilibili/app/recommend/cmd/rpc/recommend"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/recommend.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		recommend.RegisterRecommendServer(grpcServer, server.NewRecommendServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting recommend rpc server at %s...\n", c.ListenOn)
	s.Start()
}

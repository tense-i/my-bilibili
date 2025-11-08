package main

import (
	"flag"
	"fmt"

	"mybilibili/app/recall/cmd/rpc/internal/config"
	"mybilibili/app/recall/cmd/rpc/internal/server"
	"mybilibili/app/recall/cmd/rpc/internal/svc"
	"mybilibili/app/recall/cmd/rpc/recall"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/recall.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		recall.RegisterRecallServer(grpcServer, server.NewRecallServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting recall rpc server at %s...\n", c.ListenOn)
	s.Start()
}

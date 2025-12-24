package main

import (
	"flag"
	"fmt"

	"mybilibili/app/search/cmd/rpc/internal/config"
	"mybilibili/app/search/cmd/rpc/internal/server"
	"mybilibili/app/search/cmd/rpc/internal/svc"
	"mybilibili/app/search/cmd/rpc/search"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/search.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		search.RegisterSearchServer(grpcServer, server.NewSearchServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting search rpc server at %s...\n", c.ListenOn)
	s.Start()
}

// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"
	"net/http"

	"mybilibili/app/coupon/cmd/api/internal/config"
	"mybilibili/app/coupon/cmd/api/internal/handler"
	"mybilibili/app/coupon/cmd/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/coupon.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 添加 swagger.json 静态文件路由，供 Apifox 通过 URL 导入
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/doc/coupon/swagger.json",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "doc/coupon.json")
		}),
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

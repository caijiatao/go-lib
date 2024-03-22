package main

import (
	"chat-app-svr/rpc/notification/internal/logic"
	"flag"
	"fmt"

	"chat-app-svr/rpc/notification/internal/config"
	"chat-app-svr/rpc/notification/internal/server"
	"chat-app-svr/rpc/notification/internal/svc"
	"chat-app-svr/rpc/notification/notification"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/notification.yaml", "the config file")

func main() {
	flag.Parse()

	conf.MustLoad(*configFile, config.Conf())
	c := config.Conf()
	ctx := svc.NewServiceContext(*config.Conf())

	logic.InitPushServers()

	s := zrpc.MustNewServer(config.Conf().RpcServerConf, func(grpcServer *grpc.Server) {
		notification.RegisterNotificationServer(grpcServer, server.NewNotificationServer(ctx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

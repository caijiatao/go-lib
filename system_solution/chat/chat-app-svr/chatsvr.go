package main

import (
	"chat-app-svr/internal/logic/chat"
	"flag"
	"fmt"
	"net/http"

	"chat-app-svr/internal/config"
	"chat-app-svr/internal/handler"
	"chat-app-svr/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/chatsvr.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	initWebsocketServer(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

func initWebsocketServer(server *rest.Server, ctx *svc.ServiceContext) {
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/ws/chat",
		Handler: func(writer http.ResponseWriter, request *http.Request) {
			chat.NewWebsocketHandler(request.Context(), ctx).ServeWs(writer, request)
		},
	})
}

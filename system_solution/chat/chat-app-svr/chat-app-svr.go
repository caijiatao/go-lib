package main

import (
	"chat-app-svr/internal/config"
	"chat-app-svr/internal/handler"
	service "chat-app-svr/internal/logic/chat"
	"chat-app-svr/internal/svc"
	"flag"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/chat-app-svr-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	initWebsocketServer(server)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

func initWebsocketServer(server *rest.Server) {
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/ws/chat",
		Handler: func(writer http.ResponseWriter, request *http.Request) {
			service.ServeWs(writer, request)
		},
	})
}

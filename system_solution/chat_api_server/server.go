package chat_api_server

import (
	"context"
	"fmt"
	"golib/libs/logger"
	"golib/libs/naming"
	chat "golib/system_solution/chat/api"
	"golib/system_solution/chat_api_server/api"
	"golib/system_solution/cmd/chat_api_server/config"
	"google.golang.org/grpc"
	"net"
	"sync"
)

var (
	server         *Server
	chatServerOnce sync.Once
)

type Server struct {
	register  *naming.ServiceRegister
	discovery *naming.ServiceDiscovery

	chatClients map[string]chat.ChatClient

	api.UnimplementedChatApiServer
}

func (s *Server) PushMessage(ctx context.Context, request *chat.PushMessageRequest) (*chat.PushMessageReply, error) {
	logger.CtxInfof(ctx, "receive push message request: %v", request)
	return &chat.PushMessageReply{}, nil
}

func NewServer(conf *config.Config) *Server {
	chatServerOnce.Do(func() {
		var err error
		server = &Server{}
		server.runRPCServer(conf)

		registerKey := fmt.Sprintf("/api_server/%s", conf.Env.GetTarget())
		server.register, err = naming.NewServiceRegister(conf.ETCDEndpoints, registerKey, conf.Env.GetTarget(), 5)
		if err != nil {
			panic(err)
		}

		server.discovery = naming.NewServiceDiscovery(conf.ETCDEndpoints, "/chat/")
		if err != nil {
			panic(err)
		}

		go server.watchChatClients()
	})
	return server
}

func (s *Server) Close() {
	_ = s.register.Close()
	_ = s.discovery.Close()
}

func (s *Server) runRPCServer(conf *config.Config) {
	lis, err := net.Listen(conf.RPCServer.Network, conf.Env.GetTarget())
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	api.RegisterChatApiServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Errorf("rpc server serve error: %v", err)
		panic(err)
	}
}

func (s *Server) watchChatClients() {
	for {
		select {
		case _, ok := <-s.discovery.Watch():
			if !ok {
				logger.Infof("discovery exit")
				return
			}
			s.renewChatClients()
		}
	}
}

func (s *Server) renewChatClients() {
	serviceList := s.discovery.GetServiceList()
	chatClients := map[string]chat.ChatClient{}
	for k, addr := range serviceList {
		chatClient, err := newChatClient(addr)
		if err != nil {
			logger.Errorf("new api server client error: %v", err)
			continue
		}
		chatClients[k] = chatClient
	}
	for key, old := range s.chatClients {
		if _, ok := chatClients[key]; !ok {
			logger.CtxInfof(context.Background(), "api server %s(%v) offline", key, old)
		}
	}
	s.chatClients = chatClients
}

func newChatClient(addr string) (chat.ChatClient, error) {
	conn, err := grpc.Dial(addr)
	if err != nil {
		return nil, err
	}
	return chat.NewChatClient(conn), nil
}

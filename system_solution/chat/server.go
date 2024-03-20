package chat

import (
	"context"
	"fmt"
	"golib/libs/logger"
	"golib/libs/naming"
	chat "golib/system_solution/chat/api"
	"golib/system_solution/chat_api_server/api"
	"golib/system_solution/cmd/chat/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	apiServers map[string]api.ChatApiClient

	chat.UnimplementedChatServer

	rpcSvr *grpc.Server
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

		registerKey := fmt.Sprintf("/chat/%s", conf.Env.GetTarget())
		server.register, err = naming.NewServiceRegister(conf.ETCDEndpoints, registerKey, conf.Env.GetTarget(), 5)
		if err != nil {
			panic(err)
		}

		server.discovery = naming.NewServiceDiscovery(conf.ETCDEndpoints, "/api_server/")
		if err != nil {
			panic(err)
		}

		go server.watchAPIServers()
	})
	return server
}

func (s *Server) Close() {
	_ = s.register.Close()
	_ = s.discovery.Close()
	s.rpcSvr.GracefulStop()
}

func (s *Server) runRPCServer(conf *config.Config) {
	lis, err := net.Listen(conf.RPCServer.Network, conf.Env.GetTarget())
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	chat.RegisterChatServer(grpcServer, s)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Errorf("rpc server serve error: %v", err)
			panic(err)
		}
	}()
	s.rpcSvr = grpcServer
}

func (s *Server) watchAPIServers() {
	s.renewAPIServers()
	for {
		select {
		case _, ok := <-s.discovery.Watch():
			if !ok {
				logger.Infof("discovery exit")
				return
			}
			s.renewAPIServers()
		}
	}
}

func (s *Server) renewAPIServers() {
	serviceList := s.discovery.GetServiceList()
	apiServers := map[string]api.ChatApiClient{}
	for k, addr := range serviceList {
		apiServer, err := newAPIServerClient(addr)
		if err != nil {
			logger.Errorf("new api server client error: %v", err)
			continue
		}
		apiServers[k] = apiServer
	}
	for key, old := range s.apiServers {
		if _, ok := apiServers[key]; !ok {
			logger.CtxInfof(context.Background(), "api server %s(%v) offline", key, old)
		}
	}
	s.apiServers = apiServers
}

func newAPIServerClient(addr string) (api.ChatApiClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return api.NewChatApiClient(conn), nil
}

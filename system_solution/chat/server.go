package chat

import (
	"context"
	"fmt"
	"golib/libs/logger"
	"golib/libs/naming"
	"golib/system_solution/api_server/api"
	chat "golib/system_solution/chat/api"
	"golib/system_solution/cmd/chat/config"
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

	apiServers map[string]api.ApiServerClient

	chat.UnimplementedChatServer
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
}

func (s *Server) runRPCServer(conf *config.Config) {
	lis, err := net.Listen(conf.RPCServer.Network, conf.Env.GetTarget())
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	chat.RegisterChatServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Errorf("rpc server serve error: %v", err)
		panic(err)
	}
}

func (s *Server) watchAPIServers() {
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
	apiServers := map[string]api.ApiServerClient{}
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

func newAPIServerClient(addr string) (api.ApiServerClient, error) {
	conn, err := grpc.Dial(addr)
	if err != nil {
		return nil, err
	}
	return api.NewApiServerClient(conn), nil
}

package chat

import (
	"context"
	"fmt"
	"golib/libs/logger"
	"golib/libs/naming"
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
	register *naming.ServiceRegister

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
	})
	return server
}

func (s *Server) Close() {
	_ = s.register.Close()
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

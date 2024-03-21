package third_part

import (
	"context"
	"golib/libs/logger"
	"golib/libs/naming"
	"golib/system_solution/chat_api_server/api"
	"golib/system_solution/cmd/chat/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

var (
	chatAPIServer     *ChatAPIServer
	chatAPIServerOnce sync.Once
)

type ChatAPIServer struct {
	discovery *naming.ServiceDiscovery

	apiServers map[string]api.ChatApiClient
}

func NewChatApiServer() *ChatAPIServer {
	chatAPIServerOnce.Do(func() {
		chatAPIServer = &ChatAPIServer{}
		chatAPIServer.discovery = naming.NewServiceDiscovery(config.Conf().ETCDEndpoints, "/api_server/")
		go chatAPIServer.watchAPIServers()
	})
	return chatAPIServer

}

func (cs *ChatAPIServer) Close() {
	_ = cs.discovery.Close()
}

func (cs *ChatAPIServer) watchAPIServers() {
	cs.renewClients()
	for {
		select {
		case _, ok := <-cs.discovery.Watch():
			if !ok {
				logger.Infof("discovery exit")
				return
			}
			cs.renewClients()
		}
	}
}

func (cs *ChatAPIServer) renewClients() {
	serviceList := cs.discovery.GetServiceList()
	apiServers := map[string]api.ChatApiClient{}
	for k, addr := range serviceList {
		apiServer, err := newChatAPIClient(addr)
		if err != nil {
			logger.Errorf("new api server client error: %v", err)
			continue
		}
		apiServers[k] = apiServer
	}
	for key, old := range cs.apiServers {
		if _, ok := apiServers[key]; !ok {
			logger.CtxInfof(context.Background(), "api server %cs(%v) offline", key, old)
		}
	}
	cs.apiServers = apiServers
}

func newChatAPIClient(addr string) (api.ChatApiClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return api.NewChatApiClient(conn), nil
}

func (cs *ChatAPIServer) UserDetail(ctx context.Context, userId int64) (*api.UserDetailReply, error) {
	for _, apiServer := range cs.apiServers {
		reply, err := apiServer.UserDetail(ctx, &api.UserDetailRequest{UserId: userId})
		if err != nil {
			logger.CtxErrorf(ctx, "user details error: %v", err)
			continue
		}
		return reply, nil
	}
	return &api.UserDetailReply{}, nil
}

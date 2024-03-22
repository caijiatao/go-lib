package svc

import (
	"chat-app-svr/internal/config"
	"chat-app-svr/rpc/chat/chatclient"
	"chat-app-svr/rpc/user/userclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	User   userclient.User
	Chat   chatclient.Chat
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		User:   userclient.NewUser(zrpc.MustNewClient(c.User)),
		Chat:   chatclient.NewChat(zrpc.MustNewClient(c.Chat)),
	}
}

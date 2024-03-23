package svc

import (
	"chat-app-svr/internal/config"
	"chat-app-svr/rpc/receive_svr/receiver"
	"chat-app-svr/rpc/send_svr/sender"
	"chat-app-svr/rpc/user/userclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	User   userclient.User

	Send    sender.Sender
	Receive receiver.Receiver
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		User:    userclient.NewUser(zrpc.MustNewClient(c.User)),
		Send:    sender.NewSender(zrpc.MustNewClient(c.Send)),
		Receive: receiver.NewReceiver(zrpc.MustNewClient(c.Receive)),
	}
}

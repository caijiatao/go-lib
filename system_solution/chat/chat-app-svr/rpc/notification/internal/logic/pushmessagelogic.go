package logic

import (
	"chat-app-svr/rpc/notification/internal/config"
	"chat-app-svr/rpc/notification/internal/svc"
	"chat-app-svr/rpc/notification/notification"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"golib/libs/net_helper"
)

type PushMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func InitPushServers() {
	subscriber, err := discov.NewSubscriber(config.Conf().Chat.Etcd.Hosts, config.Conf().Chat.Etcd.Key)
	if err != nil {
		panic(err)
	}
	values := subscriber.Values()

	fmt.Println(config.Conf())

	for _, value := range values {
		fmt.Println("push server: ", value)
	}

	net_helper.GetFigureOutListenOn(config.Conf().RpcServerConf.ListenOn)

	subscriber.AddListener(func() {
		values := subscriber.Values()
		for _, value := range values {
			fmt.Println("push server: ", value)
		}
	})
}

func NewPushMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushMessageLogic {
	return &PushMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PushMessageLogic) PushMessage(in *notification.PushMessageRequest) (*notification.PushMessageReply, error) {
	// todo: add your logic here and delete this line

	return &notification.PushMessageReply{}, nil
}

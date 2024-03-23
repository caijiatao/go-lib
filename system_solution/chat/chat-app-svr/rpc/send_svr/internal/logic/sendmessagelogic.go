package logic

import (
	"context"

	"chat-app-svr/rpc/send_svr/internal/svc"
	"chat-app-svr/rpc/send_svr/send"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMessageLogic) SendMessage(in *send.SendMessageRequest) (*send.SendMessageReply, error) {
	// todo: add your logic here and delete this line

	return &send.SendMessageReply{}, nil
}

package logic

import (
	"context"

	"chat-app-svr/rpc/send_svr/internal/svc"
	"chat-app-svr/rpc/send_svr/send"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendGroupMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendGroupMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendGroupMessageLogic {
	return &SendGroupMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendGroupMessageLogic) SendGroupMessage(in *send.SendGroupMessageRequest) (*send.SendGroupMessageReply, error) {
	// todo: add your logic here and delete this line

	return &send.SendGroupMessageReply{}, nil
}

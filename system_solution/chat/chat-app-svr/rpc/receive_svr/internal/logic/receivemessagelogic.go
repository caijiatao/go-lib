package logic

import (
	"context"

	"chat-app-svr/rpc/receive_svr/internal/svc"
	"chat-app-svr/rpc/receive_svr/receive"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReceiveMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReceiveMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReceiveMessageLogic {
	return &ReceiveMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReceiveMessageLogic) ReceiveMessage(in *receive.ReceiveMessageRequest) (*receive.ReceiveMessageReply, error) {
	// todo: add your logic here and delete this line

	return &receive.ReceiveMessageReply{}, nil
}

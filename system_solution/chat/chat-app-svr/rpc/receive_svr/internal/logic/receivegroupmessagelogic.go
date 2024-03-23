package logic

import (
	"context"

	"chat-app-svr/rpc/receive_svr/internal/svc"
	"chat-app-svr/rpc/receive_svr/receive"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReceiveGroupMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReceiveGroupMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReceiveGroupMessageLogic {
	return &ReceiveGroupMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReceiveGroupMessageLogic) ReceiveGroupMessage(in *receive.ReceiveGroupMessageRequest) (*receive.ReceiveGroupMessageReply, error) {
	// todo: add your logic here and delete this line

	return &receive.ReceiveGroupMessageReply{}, nil
}

package logic

import (
	"context"

	"chat-app-svr/rpc/chat/chat"
	"chat-app-svr/rpc/chat/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushMessageLogic {
	return &PushMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PushMessageLogic) PushMessage(in *chat.PushMessageRequest) (*chat.PushMessageReply, error) {
	// todo: add your logic here and delete this line

	return &chat.PushMessageReply{}, nil
}

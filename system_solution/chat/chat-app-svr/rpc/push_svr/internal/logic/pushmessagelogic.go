package logic

import (
	"context"

	"chat-app-svr/rpc/push_svr/internal/svc"
	"chat-app-svr/rpc/push_svr/push"

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

func (l *PushMessageLogic) PushMessage(in *push.PushMessageRequest) (*push.PushMessageReply, error) {
	// todo: add your logic here and delete this line

	return &push.PushMessageReply{}, nil
}

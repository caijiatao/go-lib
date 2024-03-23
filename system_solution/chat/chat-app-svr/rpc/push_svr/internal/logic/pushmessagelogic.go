package logic

import (
	"context"
	"errors"

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

func (l *PushMessageLogic) checkMessage(message *push.PushMessageRequest) error {
	if message == nil {
		return errors.New("message is empty")
	}
	if message.ToUserId == 0 {
		return errors.New("to_user is empty")
	}
	if message.FromUserId == 0 {
		return errors.New("from_user is empty")
	}
	if len(message.Content) == 0 {
		return errors.New("content is empty")
	}
	return nil
}

// PushMessage 如果用户不在线则，则看是否有屏蔽消息，没有屏蔽消息则不推入消息中心，等待用户拉取
func (l *PushMessageLogic) PushMessage(in *push.PushMessageRequest) (*push.PushMessageReply, error) {
	err := l.checkMessage(in)
	if err != nil {
		return &push.PushMessageReply{}, err
	}
	return &push.PushMessageReply{}, nil
}

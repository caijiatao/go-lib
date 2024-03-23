package logic

import (
	"context"

	"chat-app-svr/rpc/user/internal/svc"
	"chat-app-svr/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserOnlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserOnlineLogic {
	return &UserOnlineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserOnlineLogic) UserOnline(in *user.UserOnlineRequest) (*user.UserOnlineReply, error) {
	// todo: add your logic here and delete this line

	return &user.UserOnlineReply{}, nil
}

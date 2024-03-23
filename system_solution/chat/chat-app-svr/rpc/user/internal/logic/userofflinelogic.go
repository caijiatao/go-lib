package logic

import (
	"context"

	"chat-app-svr/rpc/user/internal/svc"
	"chat-app-svr/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserOfflineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserOfflineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserOfflineLogic {
	return &UserOfflineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserOfflineLogic) UserOffline(in *user.UserOfflineRequest) (*user.UserOfflineReply, error) {
	// todo: add your logic here and delete this line

	return &user.UserOfflineReply{}, nil
}

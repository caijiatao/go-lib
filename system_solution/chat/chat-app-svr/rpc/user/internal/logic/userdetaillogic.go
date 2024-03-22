package logic

import (
	"context"

	"chat-app-svr/rpc/user/internal/svc"
	"chat-app-svr/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDetailLogic {
	return &UserDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserDetailLogic) UserDetail(in *user.UserDetailRequest) (*user.UserDetailReply, error) {
	u, err := l.svcCtx.UserDao.QueryUserById(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	return &user.UserDetailReply{
		UserId: u.Id,
		Phone:  u.PhoneNumber,
	}, nil
}

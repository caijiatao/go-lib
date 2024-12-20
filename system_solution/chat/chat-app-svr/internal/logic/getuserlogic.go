package logic

import (
	"chat-app-svr/internal/svc"
	"chat-app-svr/internal/types"
	"chat-app-svr/rpc/user/userclient"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(req *types.GetUserReq) (resp *types.GetUserResp, err error) {
	detail, err := l.svcCtx.User.UserDetail(l.ctx, &userclient.UserDetailRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.GetUserResp{
		UserId:      detail.UserInfo.UserId,
		PhoneNumber: detail.UserInfo.Phone,
	}

	return resp, nil
}

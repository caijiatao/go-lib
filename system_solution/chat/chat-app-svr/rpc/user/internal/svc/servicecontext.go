package svc

import (
	"chat-app-svr/rpc/model"
	"chat-app-svr/rpc/user/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(sqlx.NewSqlConn("postgres", c.DataSource)),
	}
}

package svc

import (
	"chat-app-svr/rpc/user/internal/config"
	"chat-app-svr/rpc/user/internal/dao"
)

type ServiceContext struct {
	Config config.Config

	UserDao *dao.UserDao
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserDao: dao.NewUserDao(),
	}
}

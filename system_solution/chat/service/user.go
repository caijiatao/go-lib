package service

import (
	"context"
	"sync"
)

var (
	userSvc     *userService
	userSvcOnce sync.Once
)

type userService struct{}

func UserService() *userService {
	userSvcOnce.Do(func() {
		userSvc = &userService{}
	})
	return userSvc
}

// Online
//
//	@Description: TODO 用户上线后将信息维护到ETCD中，通过ETCD可以获取用户在哪个服务器上
func (us *userService) Online(ctx context.Context, userId int64) error {
	return nil
}

// Offline
//
//	@Description: TODO 用户下线后将信息从ETCD中删除
func (us *userService) Offline(ctx context.Context, userId int64) error {
	return nil
}

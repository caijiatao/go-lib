package dao

import (
	"chat-app-svr/rpc/user/internal/model"
	"context"
	"sync"
)

var (
	userDao     *UserDao
	userDaoOnce sync.Once
)

type UserDao struct {
}

func NewUserDao() *UserDao {
	userDaoOnce.Do(func() {
		userDao = &UserDao{}
	})
	return userDao
}

func (d *UserDao) QueryUserById(ctx context.Context, userId int64) (*model.User, error) {
	return &model.User{
		Id: userId,
	}, nil
}

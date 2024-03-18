package dao

import (
	"context"
	"golib/libs/orm"
	"golib/system_solution/user/model"
	"sync"
)

var (
	userD                 *userDao
	authManageDaoInitOnce sync.Once
)

type userDao struct{}

func UserDao() *userDao {
	authManageDaoInitOnce.Do(func() {
		userD = &userDao{}
	})
	return userD
}

func (self *userDao) GetAllUsers(ctx context.Context) ([]model.User, error) {
	users := make([]model.User, 0)
	err := orm.Context(ctx).Model(model.User{}).Find(&users).Error()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (self *userDao) GetUsersByPage(ctx context.Context, pageLimit orm.QueryPageCondition) ([]model.User, int64, error) {
	total := int64(0)
	err := orm.Context(ctx).Model(model.User{}).Where("role=0").Count(&total).Error()
	if err != nil {
		return nil, 0, err
	}
	users := make([]model.User, 0)
	err = orm.Context(ctx).Model(model.User{}).Where("role=0").Limit(*pageLimit.Limit).Offset(*pageLimit.Offset).Find(&users).Error()
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (self *userDao) GetUserByUser(ctx context.Context, user model.User) (*model.User, error) {
	userFind := &model.User{}
	err := orm.Context(ctx).Model(model.User{}).Where(user).First(userFind).Error()
	if err != nil {
		return nil, err
	}
	return userFind, nil
}

func (self *userDao) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	err := orm.Context(ctx).Model(user).Create(user).Error()
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (self *userDao) UpdateUser(ctx context.Context, user *model.User) error {
	err := orm.Context(ctx).Model(user).Where("id = ?", user.ID).Updates(user).Error()
	if err != nil {
		return err
	}
	return nil
}

func (self *userDao) GetUserDevice(ctx context.Context, device *model.UserDevice) ([]model.UserDevice, error) {
	devices := make([]model.UserDevice, 0)
	err := orm.Context(ctx).Model(device).Where("user_id=?", device.UserID).Find(&devices).Error()
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (self *userDao) CreateUserDevice(ctx context.Context, device *model.UserDevice) error {
	err := orm.Context(ctx).Model(device).Create(device).Error()
	if err != nil {
		return err
	}
	return nil
}

package service

import (
	"github.com/gin-gonic/gin"
	"golib/libs/orm"
	"golib/system_solution/user/dao"
	"golib/system_solution/user/define"
	"golib/system_solution/user/model"
	"sync"
)

var (
	userS         *userService
	userSInitOnce sync.Once
)

type userService struct{}

func UserService() *userService {
	userSInitOnce.Do(func() {
		userS = &userService{}
	})
	return userS
}

func (self *userService) UserList(ctx *gin.Context, req define.UserListReq) (resp *define.UserListResp, err error) {
	users, total, err := dao.UserDao().GetUsersByPage(ctx, orm.GetPageCondition(req.PageNum, req.PageSize))
	if err != nil {
		return nil, err
	}
	resp = &define.UserListResp{
		Total: total,
	}

	for _, user := range users {
		resp.UserListResp = append(resp.UserListResp, define.UserInfoResp{
			UserID:      user.ID,
			PhoneNumber: user.PhoneNumber,
			NickName:    user.NickName,
		})
	}
	return resp, nil
}

func (self *userService) UserDetail(ctx *gin.Context, req define.AdminUserDetailReq) (resp *define.UserInfoResp, err error) {
	user, err := dao.UserDao().GetUserByUser(ctx, model.User{
		ID: req.UserID,
	})
	if err != nil {
		return nil, err
	}
	resp = &define.UserInfoResp{
		UserID:      user.ID,
		PhoneNumber: user.PhoneNumber,
		NickName:    user.NickName,
	}

	return resp, nil
}

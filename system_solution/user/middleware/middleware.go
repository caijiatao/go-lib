package middleware

import (
	"github.com/gin-gonic/gin"
	"golib/libs/gin_helper"
	"golib/system_solution/user/dao"
	"golib/system_solution/user/model"
)

func Authentication() gin.HandlerFunc {
	return AuthenticationVarious(func(ctx *gin.Context) error {
		userCtxPhone, exit := ctx.Get(CtxUserPhone)
		if !exit {
			return nil
		}
		userPhone := userCtxPhone.(string)
		if userPhone == "" {
			return nil
		}
		// todo 引入缓存
		userInfo, err := dao.UserDao().GetUserByUser(ctx, model.User{
			PhoneNumber: userPhone,
		})
		if err != nil {
			return gin_helper.INNER_ERROR
		}
		if userInfo == nil {
			return gin_helper.PARAM_ERROR
		}
		return nil
	})
}

type authMiddleware struct {
	ctx *gin.Context
}

func AuthMiddleware(ctx *gin.Context) *authMiddleware {
	return &authMiddleware{
		ctx: ctx,
	}
}

func (am authMiddleware) GetUser() *User {
	return GetUser(am.ctx)
}

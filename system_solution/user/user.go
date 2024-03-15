package user

import "sync"

var (
	userController         *UserController
	userControllerInitOnce sync.Once
)

type UserController struct{}

func GetUserController() *UserController {
	userControllerInitOnce.Do(func() {
		userController = &UserController{}
	})
	return userController
}

func (u *UserController) RegisterRoutes() {
}

package good_t

type Login interface {
	Login(user User) bool
}

type wxLogin struct{}

func (w *wxLogin) Login(user User) bool {
	return true
}

type phoneLogin struct{}

func (p *phoneLogin) Login(user User) bool {
	return true
}

type Register interface {
	Register(user User) bool
}

type wxRegister struct{}

func (w *wxRegister) Register(user User) bool {
	return true
}

type phoneRegister struct{}

func (p *phoneRegister) Register(user User) bool {
	return true
}

type UserService struct {
	Login
	Register
}

func NewUserService() *UserService {
	return &UserService{}
}

package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	goRedis "github.com/redis/go-redis/v9"
	"golib/libs/gin_helper"
	"golib/libs/logger"
	"golib/libs/redis"
	"golib/libs/sms"
	"golib/libs/util"
	"golib/system_solution/user/dao"
	"golib/system_solution/user/define"
	"golib/system_solution/user/middleware"
	"golib/system_solution/user/model"
	"gorm.io/gorm"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var (
	authS         *authService
	authSInitOnce sync.Once
)

type authService struct{}

func AuthService() *authService {
	authSInitOnce.Do(func() {
		authS = &authService{}
	})
	return authS
}

func (self *authService) GenerateCaptcha(ctx *gin.Context, width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}

	return sb.String()
}

func (self *authService) SendSms(ctx *gin.Context, req define.SmsReq, ip string) error {
	err := self.limitSms(ctx, ip, 100)
	if err != nil {
		return err
	}
	// 手机号限制；单个手机号一天只可以发20次
	err = self.limitSms(ctx, req.PhoneNumber, 20)
	if err != nil {
		return err
	}

	smsCode := self.GenerateCaptcha(ctx, 6)

	err = redis.ClusterClient().SetEx(ctx, fmt.Sprintf(redis.CAPTCHA_KEY, req.PhoneNumber), smsCode, 5*time.Minute).Err()
	if err != nil {
		return err
	}

	smsConfig := sms.Config()
	captchaCode := define.SmsTemplateCode{
		smsCode,
	}
	code, _ := json.Marshal(captchaCode)
	err = sms.Client().SendSms(req.PhoneNumber, smsConfig.SignName, smsConfig.TemplateCode, string(code))
	if err != nil {
		logger.Error("send Sms reqID:%s ERR:%s", logger.CtxTraceID(ctx), err.Error())
		return err
	}
	return nil
}

func (self *authService) VerifySms(ctx *gin.Context, phoneNumber string, smsCode string) error {
	code, err := redis.ClusterClient().Get(ctx, fmt.Sprintf(redis.CAPTCHA_KEY, phoneNumber)).Result()
	if err != nil {
		return gin_helper.SMSCODE_ERROR
	}
	if code != smsCode {
		return gin_helper.SMSCODE_ERROR
	}
	return nil
}

func (self *authService) IsNewUser(ctx *gin.Context, phoneNumber string) (bool, error) {
	_, err := dao.UserDao().GetUserByUser(ctx, model.User{
		PhoneNumber: phoneNumber,
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (self *authService) MockLogin() (*define.LoginRsp, error) {
	return &define.LoginRsp{
		ID:          rand.Int63(),
		PhoneNumber: "12345678901",
		NickName:    "mock",
		IsAdmin:     true,
		Token:       "123456",
	}, nil
}

func (self *authService) Login(ctx *gin.Context, req define.LoginReq) (*define.LoginRsp, error) {
	return self.MockLogin()
	// 校验密码
	if req.Password != "" && req.SmsCode == "" {
		err := self.VerfiyAccountPwd(ctx, req)
		if err != nil {
			return nil, err
		}
	}

	modelUser := &model.User{
		PhoneNumber: req.PhoneNumber,
	}
	user, err := dao.UserDao().GetUserByUser(ctx, *modelUser)
	if err != nil {
		if err == gorm.ErrRecordNotFound && req.SmsCode != "" && req.Password != "" {
			// 未找到账号 && 验证码情况下新建账号
			user, err = self.Register(ctx, req)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// 执行登录操作
	token, err := self.DoLogin(ctx, user)
	if err != nil {
		logger.Errorf("DoLogin fail ,reqID:%s ,err:%s", logger.CtxTraceID(ctx), err.Error())
		return nil, gin_helper.INNER_ERROR
	}

	return &define.LoginRsp{
		ID:          user.ID,
		PhoneNumber: user.PhoneNumber,
		NickName:    user.NickName,
		Token:       token,
	}, nil
}

func (self *authService) VerfiyAccountPwd(ctx *gin.Context, req define.LoginReq) error {
	modelUser := &model.User{
		PhoneNumber: req.PhoneNumber,
	}
	user, err := dao.UserDao().GetUserByUser(ctx, *modelUser)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = gin_helper.LOGIN_UNKNOWN
		}
		return err
	}
	// 账号锁定时间未到
	now := time.Now().Unix()
	lockTime := user.LockTime.Add(-8 * time.Hour).Unix()
	if lockTime > now {
		restTime := time.Unix(lockTime-now, 0)
		lockErr := gin_helper.NewError(gin_helper.LOGIN_LOCK.Code, fmt.Sprintf(gin_helper.LOGIN_LOCK.Msg, restTime.Minute(), restTime.Second()))
		return lockErr
	}

	if user.Password != util.EncryptPassword(req.Password, define.PasswordSalt) {
		return self.LockUser(ctx, user)
	}
	return nil
}

func (self *authService) LockUser(ctx *gin.Context, user *model.User) error {
	var lockErr *gin_helper.MyError
	userUpdate := &model.User{
		ID:         user.ID,
		UpdateTime: time.Now(),
	}
	// 更新user记录
	err := dao.UserDao().UpdateUser(ctx, userUpdate)
	if err != nil {
		logger.Errorf("database update user error: %s", err)
		return err
	}
	return lockErr
}

func (self *authService) DoLogin(ctx *gin.Context, user *model.User) (token string, err error) {
	// 生成token
	token, err = middleware.GenerateToken(user.ID, user.PhoneNumber)
	if err != nil {
		logger.Errorf("token generate error: %s", err)
		return "", err
	}

	err = middleware.DoLogin(ctx, token, user.PhoneNumber)
	if err != nil {
		return "", err
	}
	// 更新user记录
	userUpdate := &model.User{
		ID:         user.ID,
		UpdateTime: time.Now(),
	}
	err = dao.UserDao().UpdateUser(ctx, userUpdate)
	if err != nil {
		logger.Errorf("database update user error: %s", err)
		return "", err
	}
	return token, nil
}

func (self *authService) Logout(ctx *gin.Context, user middleware.User, token string) error {
	err := middleware.DoLogOut(ctx, token, user.PhoneNumber)
	if err != nil {
		return err
	}
	return nil

}

func (self *authService) Register(ctx *gin.Context, req define.LoginReq) (*model.User, error) {
	user := &model.User{
		PhoneNumber: req.PhoneNumber,
		Password:    util.EncryptPassword(req.Password, define.PasswordSalt),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
		ExpireTime:  time.Now().Add(7 * 24 * time.Hour),
	}
	_, err := dao.UserDao().CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (self *authService) Password(ctx *gin.Context, userID int64, req define.PasswordReq) error {
	user, err := dao.UserDao().GetUserByUser(ctx, model.User{
		ID: userID,
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = gin_helper.LOGIN_UNKNOWN
		}
		return err
	}
	if user.Password != "" {
		if user.Password != util.EncryptPassword(req.OldPassword, define.PasswordSalt) {
			return gin_helper.NewError(gin_helper.ERR_PASSWORD.Code, "原始密码输入错误，请重新输入")
		}
		if user.Password == util.EncryptPassword(req.Password, define.PasswordSalt) {
			return gin_helper.SAME_PASSWORD
		}
	}
	user.Password = util.EncryptPassword(req.Password, define.PasswordSalt)
	user.UpdateTime = time.Now()

	return dao.UserDao().UpdateUser(ctx, user)
}

func (self *authService) ForgetPassword(ctx *gin.Context, req define.ForgetPasswordReq) error {
	user, err := dao.UserDao().GetUserByUser(ctx, model.User{
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = gin_helper.LOGIN_UNKNOWN
		}
		return err
	}
	err = self.VerifySms(ctx, req.PhoneNumber, req.SmsCode)
	if err != nil {
		return err
	}

	if user.Password == util.EncryptPassword(req.Password, define.PasswordSalt) {
		return gin_helper.SAME_PASSWORD
	}
	user.Password = util.EncryptPassword(req.Password, define.PasswordSalt)
	user.UpdateTime = time.Now()
	return dao.UserDao().UpdateUser(ctx, user)
}

func (self *authService) limitSms(ctx *gin.Context, key string, limitCount int64) error {
	key = fmt.Sprintf("%s:%s", "sms:limit:", key)
	if limitCount <= 0 {
		return nil
	}
	c, err := redis.ClusterClient().Get(ctx, key).Int64()
	if err != nil && err != goRedis.Nil {
		return err
	}
	if c >= limitCount {
		return gin_helper.SMSCODE_LIMIT
	}
	t := redis.ClusterClient().TTL(ctx, key).Val()
	if t < 0 {
		t = 24 * time.Hour
	}
	err = redis.ClusterClient().SetEx(ctx, key, c+1, t).Err()
	if err != nil {
		return err
	}
	return nil
}

func (self *authService) UserInfo(ctx *gin.Context, userID int64) (*define.UserInfoResp, error) {
	userModel := &model.User{
		ID: userID,
	}
	userInfo, err := dao.UserDao().GetUserByUser(ctx, *userModel)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = gin_helper.LOGIN_UNKNOWN
		}
		return nil, err
	}

	var newUser bool
	if userInfo.Password == "" {
		newUser = true
	}
	resp := &define.UserInfoResp{
		UserID:      userInfo.ID,
		PhoneNumber: userInfo.PhoneNumber,
		NickName:    userInfo.NickName,
		NewUser:     newUser,
	}
	return resp, nil
}

func (self *authService) UserDevices(ctx *gin.Context, phoneNumber string) ([]model.UserDevice, error) {
	userModel := &model.User{
		PhoneNumber: phoneNumber,
	}
	userInfo, err := dao.UserDao().GetUserByUser(ctx, *userModel)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return dao.UserDao().GetUserDevice(ctx, &model.UserDevice{
		UserID: userInfo.ID,
	})
}

func (self *authService) GetUserIp(ctx *gin.Context) string {
	ip := ctx.Request.Header.Get("X-Customize-Forwarded-For")
	if strings.TrimSpace(ip) == "" {
		remoteAddr := ctx.ClientIP()
		addrPort := strings.Split(remoteAddr, ":")
		ip = addrPort[0]
	}
	return ip
}

func (self *authService) AddUserDevice(ctx *gin.Context, userID int64) error {
	devices, err := dao.UserDao().GetUserDevice(ctx, &model.UserDevice{
		UserID: userID,
	})
	if err != nil {
		return nil
	}
	ip := self.GetUserIp(ctx)
	userAgent := ctx.GetHeader("User-Agent")
	var exist bool
	for _, val := range devices {
		if val.IP == ip && val.Browser == userAgent {
			exist = true
		}
	}
	if !exist {
		return dao.UserDao().CreateUserDevice(ctx, &model.UserDevice{
			UserID:     userID,
			IP:         ip,
			Browser:    userAgent,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		})
	}
	return nil
}

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golib/libs/captcha"
	"golib/libs/gin_helper"
	"golib/libs/limiter"
	"golib/libs/util"
	"golib/system_solution/user/define"
	"golib/system_solution/user/middleware"
	"golib/system_solution/user/service"
	"regexp"
	"strings"
)

type UserController struct{}

func (uc *UserController) RegisterRoutes(r gin.IRouter) {
	r = r.Group("/api/user")
	r.Use(gin_helper.CorsMiddleware())
	r.Use(limiter.IpLimiter())
	r.POST("/sms", uc.SendSms)
	r.POST("/sms/verify", uc.VerifySms)
	r.POST("/login", uc.Login)
	r.POST("/password/forget", uc.ForgetPassword)

	r.Use(middleware.Authentication())
	r.GET("/user", uc.UserInfo)
	r.POST("/logout", uc.Logout)
	r.POST("/password", uc.Password) //第一次设置密码或者更改密码

}

func (uc *UserController) Captcha(ctx *gin.Context) {
	captchaID := captcha.NewCaptchaID()
	gin_helper.SendSuccessResp(ctx, define.CaptchaResp{
		CaptchaID: captchaID,
	})
}

func (uc *UserController) SendSms(ctx *gin.Context) {
	req := define.SmsReq{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}
	ok := uc.VerifyPhoneNumber(req.PhoneNumber)
	if !ok {
		ctx.Error(errors.New("请输入正确的手机号码"))
		return
	}
	remoteAddr := ctx.ClientIP()
	addrPort := strings.Split(remoteAddr, ":")
	ip := addrPort[0]
	err = service.AuthService().SendSms(ctx, req, ip)
	if err != nil {
		ctx.Error(err)
		return
	}

	gin_helper.SendSuccessResp(ctx, nil)
}

func (uc *UserController) VerifySms(ctx *gin.Context) {
	req := define.VerifySmsReq{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}
	ok := uc.VerifyPhoneNumber(req.PhoneNumber)
	if !ok {
		ctx.Error(errors.New("请输入正确的手机号码"))
		return
	}
	err = service.AuthService().VerifySms(ctx, req.PhoneNumber, req.SmsCode)
	if err != nil {
		ctx.Error(err)
		return
	}

	newUser, err := service.AuthService().IsNewUser(ctx, req.PhoneNumber)
	if err != nil {
		ctx.Error(err)
		return
	}

	gin_helper.SendSuccessResp(ctx, &define.NewUserRsp{
		NewUser: newUser,
	})
}

func (uc *UserController) Login(ctx *gin.Context) {
	req := define.LoginReq{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	//newDevice, err := uc.IsNewDevice(ctx, req.PhoneNumber)
	//if err != nil {
	//	ctx.Error(err)
	//	return
	//}
	//if newDevice && req.SmsCode == "" {
	//	ctx.Error(errors.New("该设备首次登录，请进行验证码验证"))
	//	return
	//}
	//ok := uc.VerifyPhoneNumber(req.PhoneNumber)
	//if !ok {
	//	ctx.Error(errors.New("请输入正确的手机号码"))
	//	return
	//}
	//if req.Password == "" && req.SmsCode == "" {
	//	ctx.Error(errors.New("请输入密码或者验证码"))
	//	return
	//}
	//if req.Password != "" {
	//	req.Password, err = util.Decrypt(req.Password)
	//	if err != nil {
	//		ctx.Error(err)
	//		return
	//	}
	//	if !uc.VerifyPwd(req.Password) {
	//		ctx.Error(errors.New("密码不合法"))
	//		return
	//	}
	//}
	rsp, err := service.AuthService().Login(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}
	gin_helper.SendSuccessResp(ctx, rsp)
}

func (uc *UserController) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("token")
	user := middleware.GetUser(ctx)
	err := service.AuthService().Logout(ctx, *user, token)
	if err != nil {
		ctx.Error(err)
		return
	}
	gin_helper.SendSuccessResp(ctx, nil)
}

func (uc *UserController) ForgetPassword(ctx *gin.Context) {
	req := define.ForgetPasswordReq{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(gin_helper.PARAM_ERROR)
		return
	}
	if req.Password != "" {
		req.Password, err = util.Decrypt(req.Password)
		if err != nil {
			ctx.Error(gin_helper.ERR_PASSWORD)
			return
		}
	}
	if !uc.VerifyPwd(req.Password) {
		ctx.Error(gin_helper.ILLEGAL_PASSWORD)
		return
	}
	err = service.AuthService().ForgetPassword(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	gin_helper.SendSuccessResp(ctx, nil)
}

func (uc *UserController) Password(ctx *gin.Context) {
	req := define.PasswordReq{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(gin_helper.PARAM_ERROR)
		return
	}
	req.Password, err = util.Decrypt(req.Password)
	if err != nil {
		ctx.Error(gin_helper.ERR_PASSWORD)
		return
	}
	if req.OldPassword != "" {
		req.OldPassword, err = util.Decrypt(req.OldPassword)
		if err != nil {
			ctx.Error(gin_helper.ERR_PASSWORD)
			return
		}
	}
	if !uc.VerifyPwd(req.Password) {
		ctx.Error(gin_helper.ILLEGAL_PASSWORD)
		return
	}

	user := middleware.GetUser(ctx)
	err = service.AuthService().Password(ctx, user.UserID, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	gin_helper.SendSuccessResp(ctx, nil)
}

func (uc *UserController) UserInfo(ctx *gin.Context) {
	user := middleware.GetUser(ctx)
	resp, err := service.AuthService().UserInfo(ctx, user.UserID)
	if err != nil {
		ctx.Error(err)
		return
	}
	resp.NewDevice, err = uc.IsNewDevice(ctx, resp.PhoneNumber)
	if err != nil {
		ctx.Error(err)
		return
	}

	gin_helper.SendSuccessResp(ctx, resp)
}

func (uc *UserController) VerifyPwd(pwd string) bool {
	if len(pwd) < 6 || len(pwd) > 18 {
		return false
	}
	// 过滤掉这四类字符以外的密码串,直接判断不合法
	re, err := regexp.Compile(`^[a-zA-Z0-9.@$!%*#_~?&^]{6,18}$`)
	if err != nil {
		return false
	}
	match := re.MatchString(pwd)
	if !match {
		return false
	}

	var level = 0
	patternList := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`, `[.@$!%*#_~?&^]+`}
	for _, pattern := range patternList {
		match, _ := regexp.MatchString(pattern, pwd)
		if match {
			level++
		}
	}
	if level < 2 {
		return false
	}

	return true
}

func (uc *UserController) VerifyPhoneNumber(phone string) bool {
	regRuler := "^1[345789]{1}\\d{9}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(phone)
}

func (uc *UserController) IsNewDevice(ctx *gin.Context, phoneNumber string) (bool, error) {
	ip := service.AuthService().GetUserIp(ctx)
	devices, err := service.AuthService().UserDevices(ctx, phoneNumber)
	if err != nil {
		return false, err
	}
	for _, val := range devices {
		if val.IP == ip {
			return false, nil
		}
	}
	return true, nil
}

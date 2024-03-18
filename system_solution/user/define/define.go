package define

import "golib/libs/orm"

const PasswordSalt = "salt"

const MaxFailedCount = 5

type LoginReq struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	SmsCode     string `json:"smsCode"`
	Password    string `json:"password"`
}

type VerifySmsReq struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	SmsCode     string `json:"smsCode"  binding:"required"`
}

type SmsReq struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
}

type SmsTemplateCode struct {
	Code string `json:"code"`
}

type CaptchaReq struct {
	CaptchaID   string `json:"captchaID" binding:"required"`
	CaptchaCode string `json:"captchaCode" binding:"required"`
}

type CaptchaResp struct {
	CaptchaID string `json:"captchaID"`
}

type LoginRsp struct {
	ID          int64  `json:"id"`
	PhoneNumber string `json:"phoneNumber"`
	NickName    string `json:"nickName"`
	IsAdmin     bool   `json:"isAdmin"`
	Token       string `json:"token"`
}

type NewUserRsp struct {
	NewUser bool `json:"newUser"`
}

type ForgetPasswordReq struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	SmsCode     string `json:"smsCode" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type PasswordReq struct {
	Password    string `json:"password" binding:"required"`
	OldPassword string `json:"oldPassword"`
}

type ApplicationReq struct {
	UserID    int64
	UserName  string `json:"userName" binding:"required"`
	Telephone string `json:"telephone" binding:"required"`
	Industry  string `json:"industry" binding:"required"`
	Email     string `json:"email"`
	Company   string `json:"company"`
	Position  string `json:"position"`
	Note      string `json:"note"`
}

type UserInfoResp struct {
	UserID      int64  `json:"userID"`
	PhoneNumber string `json:"phoneNumber"`
	NickName    string `json:"nickName"`
	NewUser     bool   `json:"newUser"`
	NewDevice   bool   `json:"newDevice"`
}

type UserListReq struct {
	orm.PageParams
}

type UserListWhereCondition struct {
	*orm.QueryPageCondition
}

type UserListResp struct {
	UserListResp []UserInfoResp `json:"userList"`
	Total        int64          `json:"total"`
}

type AdminUserDetailReq struct {
	UserID int64 `json:"userID" binding:"required" form:"userID"`
}

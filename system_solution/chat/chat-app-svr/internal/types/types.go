// Code generated by goctl. DO NOT EDIT.
package types

type GetUserReq struct {
	UserId string `form:"userId"`
}

type GetUserResp struct {
	UserId      string `json:"userId"`
	PhoneNumber string `json:"phoneNumber"`
}

type LoginReq struct {
	PhoneNumber string `form:"phoneNumber"`
	Password    string `form:"password"`
}

type LoginResp struct {
	UserId string `json:"userId"`
	Token  string `json:"token"`
}
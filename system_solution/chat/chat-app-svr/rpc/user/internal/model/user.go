package model

import "time"

type User struct {
	Id          int64     `json:"id"`
	NickName    string    `json:"nickName"`
	PhoneNumber string    `json:"phoneNumber"`
	Profile     string    `json:"profile"`
	Password    string    `json:"password"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}

type UserFriend struct {
	Id         int64     `json:"id"`
	UserId     int64     `json:"userId"`
	FriendId   int64     `json:"friendId"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

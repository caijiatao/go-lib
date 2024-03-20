package model

import (
	"time"
)

// User 用户
type User struct {
	ID          int64     `gorm:"column:id"`
	NickName    string    `gorm:"column:nick_name"`
	PhoneNumber string    `gorm:"column:phone_number"`
	Password    string    `gorm:"column:password"`
	Profile     string    `gorm:"column:profile"`
	LockTime    time.Time `gorm:"column:lock_time"`
	CreateTime  time.Time `gorm:"column:create_time"`
	UpdateTime  time.Time `gorm:"column:update_time"`
}

func (user *User) TableName() string {
	return "user"
}

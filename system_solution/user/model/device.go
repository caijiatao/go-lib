package model

import (
	"time"
)

// UserDevice 用户设备
type UserDevice struct {
	ID         int64     `gorm:"column:id"`
	UserID     int64     `gorm:"column:user_id"`
	IP         string    `gorm:"column:ip"`
	Browser    string    `gorm:"column:browser"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

func (self *UserDevice) TableName() string {
	return "user_device"
}

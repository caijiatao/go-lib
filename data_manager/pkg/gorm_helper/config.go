package gorm_helper

import "gorm.io/gorm"

type Config struct {
	DBClientName string
	Dial         gorm.Dialector
	*gorm.Config
}

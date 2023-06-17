package orm

import (
	"gorm.io/gorm"
	"time"
)

type GormClientProxy struct {
	DriverOperator
	Client
	*gorm.DB
}

func NewGormClientProxy(config Config) (Client, error) {
	db, err := gorm.Open(config.Dial, config.Config)
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if len(config.DBClientName) == 0 {
		config.DBClientName = defaultDBClientName
	}

	return &GormClientProxy{
		DB: db,
	}, nil
}

func (p *GormClientProxy) Check() (ok bool) {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return false
	}
	err = sqlDB.Ping()
	if err != nil {
		return false
	}
	return true
}

func (p *GormClientProxy) GetSchema() Schema {
	// TODO xorm 有实现
	return nil
}

package gorm_helper

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type gormContextKeyType string

const dbContextKey gormContextKeyType = "db-context-key"

type DB struct {
	*gorm.DB
}

func InitDb(config Config) (err error) {
	globalDBManagerInitOnce.Do(func() {
		globalDBManager = &dbManager{}
	})

	db, err := gorm.Open(config.Dial, config.Config)
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if len(config.DBClientName) == 0 {
		config.DBClientName = defaultDBClientName
	}

	err = globalDBManager.add(config.DBClientName, &DB{
		DB: db,
	})
	if err != nil {
		return err
	}
	return nil
}

// BindContext
//
//	@Description: 将DB访问资源绑定到context中
func BindContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, dbContextKey, globalDBManager)
}

func Context(ctx context.Context) *DB {
	value := ctx.Value(dbContextKey)
	if value == nil {
		return nil
	}
	return value.(*dbManager).get(defaultDBClientName)
}

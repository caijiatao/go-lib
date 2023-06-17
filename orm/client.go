package orm

import (
	"context"
)

type gormContextKeyType string

const defaultKey gormContextKeyType = "default"

type Client interface {
	DriverOperator
	// TODO ORM 要暴露出来的方法
}

type DriverOperator interface {
	Check() bool
	GetSchema() Schema
}

func NewClient(config Config) (client Client, err error) {
	globalDBManagerInitOnce.Do(func() {
		globalDBManager = &ormManger{}
	})

	client, err = NewGormClientProxy(config)
	if err != nil {
		return nil, err
	}

	err = globalDBManager.add(config.DBClientName, client)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// BindContext
//
//	@Description: 将DB访问资源绑定到context中
func BindContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, defaultKey, globalDBManager)
}

func Context(ctx context.Context) Client {
	value := ctx.Value(defaultKey)
	if value == nil {
		return nil
	}
	return value.(*ormManger).get(defaultDBClientName)
}

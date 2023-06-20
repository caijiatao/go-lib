package orm

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type contextKeyType string

const dbContextKey contextKeyType = "default"

// Client
// @Description: ORM client 接口
type Client interface {
	DriverMetaOperator
	Create(value interface{}) Client
	CreateInBatches(value interface{}, batchSize int) Client
	Save(value interface{}) Client
	First(dest interface{}, conds ...interface{}) Client
	Take(dest interface{}, conds ...interface{}) Client
	Last(dest interface{}, conds ...interface{}) Client
	Find(dest interface{}, conds ...interface{}) Client
	FirstOrInit(dest interface{}, conds ...interface{}) Client
	FirstOrCreate(dest interface{}, conds ...interface{}) Client
	Update(column string, value interface{}) Client
	Updates(values interface{}) Client
	UpdateColumn(column string, value interface{}) Client
	UpdateColumns(values interface{}) Client
	Delete(value interface{}, conds ...interface{}) Client
	Count(count *int64) Client
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	Scan(dest interface{}) Client
	Pluck(column string, dest interface{}) Client
	ScanRows(rows *sql.Rows, dest interface{}) error
	Begin(opts ...*sql.TxOptions) Client
	Commit() Client
	Rollback() Client
	SavePoint(name string) Client
	RollbackTo(name string) Client
	Exec(sql string, values ...interface{}) Client
	Model(value interface{}) Client
	Where(query interface{}, args ...interface{}) Client
	Table(name string, args ...interface{}) Client
}

func NewOrmClient(config Config) (err error) {
	globalClientManagerInitOnce.Do(func() {
		globalClientManager = &clientManager{}
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

	err = globalClientManager.add(config.DBClientName, &GormClientProxy{
		db: db,
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
	return context.WithValue(ctx, dbContextKey, globalClientManager)
}

// Context
//
//	@Description: 从上下文中获取Client
func Context(ctx context.Context) Client {
	return getClientByName(ctx, defaultDBClientName)
}

func ContextByClientName(ctx context.Context, clientName string) Client {
	return getClientByName(ctx, clientName)
}

func GetClientByClientName(clientName string) Client {
	return globalClientManager.get(clientName)
}

func getClientByName(ctx context.Context, clientName string) Client {
	value := ctx.Value(dbContextKey)
	if value == nil {
		return nil
	}
	return value.(*clientManager).get(clientName)
}

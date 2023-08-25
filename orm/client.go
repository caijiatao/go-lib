package orm

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
)

type contextKeyType string

const dbContextKey contextKeyType = "default"
const SourceDBClientNamePrefix = "source: "

type Query interface {
	First(dest interface{}, conds ...interface{}) Client
	Take(dest interface{}, conds ...interface{}) Client
	Last(dest interface{}, conds ...interface{}) Client
	Find(dest interface{}, conds ...interface{}) Client
	FirstOrInit(dest interface{}, conds ...interface{}) Client
	FirstOrCreate(dest interface{}, conds ...interface{}) Client
	Count(count *int64) Client
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	Scan(dest interface{}) Client
	Pluck(column string, dest interface{}) Client
	ScanRows(rows *sql.Rows, dest interface{}) error
	Joins(sql string, args ...interface{}) Client
	Select(query interface{}, args ...interface{}) Client
}

type Create interface {
	Create(value interface{}) Client
	CreateInBatches(value interface{}, batchSize int) Client
	Save(value interface{}) Client
	DB() *gorm.DB
}

type Update interface {
	Update(column string, value interface{}) Client
	Updates(values interface{}) Client
	UpdateColumn(column string, value interface{}) Client
	UpdateColumns(values interface{}) Client
	Begin(opts ...*sql.TxOptions) Client
	Commit() Client
	Rollback() Client
	SavePoint(name string) Client
	RollbackTo(name string) Client
	Exec(sql string, values ...interface{}) Client
	Transaction(fc func(tx Client) error, opts ...*sql.TxOptions) error
	Delete(value interface{}, conds ...interface{}) Client
}

type Basic interface {
	Model(value interface{}) Client
	Where(query interface{}, args ...interface{}) Client
	Limit(limit int) Client
	Offset(offset int) Client
	Order(value interface{}) Client
	Table(name string, args ...interface{}) Client
	Raw(sql string, values ...interface{}) Client
	RowsAffected() int64
	Preload(sql string, args ...interface{}) Client
	Error() error
}

// Client
// @Description: ORM client 接口 ，底层可以替换实现
type Client interface {
	DriverMetaOperator
	AdvanceQuery
	Create
	Query
	Update
	Basic
}

func Init() error {
	configs := readConfigs()
	for _, config := range configs {
		err := NewOrmClient(config)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewOrmClient(config *Config) (err error) {
	if len(config.DBClientName) == 0 {
		config.DBClientName = defaultDBClientName
	}

	client, err := NewGormClientProxy(*config)
	if err != nil {
		return err
	}

	err = globalClientManager.add(config.DBClientName, client)
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
	if clientName == "" {
		return globalClientManager.get(defaultDBClientName)
	}
	return globalClientManager.get(clientName)
}

func getClientByName(ctx context.Context, clientName string) Client {
	value := ctx.Value(dbContextKey)
	if value == nil {
		return nil
	}
	return value.(*clientManager).get(clientName)
}

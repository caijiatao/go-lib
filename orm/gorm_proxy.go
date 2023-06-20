package orm

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"time"
)

type GormClientProxy struct {
	driverOperator DriverMetaOperator
	client         Client
	db             *gorm.DB
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

	client := &GormClientProxy{
		db: db,
	}

	metaDriverOperatorConstructor, ok := driverMetaConstructorMap[db.Dialector.Name()]
	if ok {
		client.driverOperator = metaDriverOperatorConstructor(db)
	}

	return client, nil
}

func (clientProxy *GormClientProxy) Create(value interface{}) Client {
	c := clientProxy.clone()
	c.db = clientProxy.db.Create(value)
	return c
}

func (clientProxy *GormClientProxy) CreateInBatches(value interface{}, batchSize int) Client {
	c := clientProxy.clone()
	c.db = c.db.CreateInBatches(value, batchSize)
	return c
}

func (clientProxy *GormClientProxy) Save(value interface{}) Client {
	c := clientProxy.clone()
	c.db = c.db.Save(value)
	return c
}

func (clientProxy *GormClientProxy) First(dest interface{}, conds ...interface{}) Client {
	c := clientProxy.clone()
	c.db = c.db.First(dest, conds...)
	return c
}

func (clientProxy *GormClientProxy) Take(dest interface{}, conds ...interface{}) Client {
	c := clientProxy.clone()
	c.db = c.db.Take(dest, conds...)
	return c
}

func (clientProxy *GormClientProxy) Last(dest interface{}, conds ...interface{}) Client {
	c := clientProxy.clone()
	c.db = c.db.Last(dest, conds...)
	return c
}

func (clientProxy *GormClientProxy) Find(dest interface{}, conds ...interface{}) Client {
	c := clientProxy.clone()
	c.db = c.db.Find(dest, conds...)
	return c
}

func (clientProxy *GormClientProxy) FirstOrInit(dest interface{}, conds ...interface{}) Client {
	c := clientProxy.clone()
	c.db = c.db.FirstOrInit(dest, conds...)
	return c
}

func (clientProxy *GormClientProxy) FirstOrCreate(dest interface{}, conds ...interface{}) Client {
	c := clientProxy.clone()
	c.db = c.db.FirstOrCreate(dest, conds...)
	return c
}

func (clientProxy *GormClientProxy) Update(column string, value interface{}) Client {
	c := clientProxy.clone()
	c.db = c.db.Update(column, value)
	return c
}

func (clientProxy *GormClientProxy) Updates(values interface{}) Client {
	c := clientProxy.clone()
	c.db = c.db.Updates(values)
	return c
}

func (clientProxy *GormClientProxy) UpdateColumn(column string, value interface{}) Client {
	c := clientProxy.clone()
	c.db = c.db.UpdateColumn(column, value)
	return c
}

func (clientProxy *GormClientProxy) UpdateColumns(values interface{}) Client {
	c := clientProxy.clone()
	c.db = c.db.UpdateColumns(values)
	return c
}

func (clientProxy *GormClientProxy) Delete(value interface{}, conds ...interface{}) Client {
	c := clientProxy.clone()
	c.db = c.db.Delete(value, conds...)
	return c
}

func (clientProxy *GormClientProxy) Count(count *int64) Client {
	c := clientProxy.clone()
	c.db = c.db.Count(count)
	return c
}

func (clientProxy *GormClientProxy) Row() *sql.Row {
	return clientProxy.db.Row()
}

func (clientProxy *GormClientProxy) Rows() (*sql.Rows, error) {
	return clientProxy.db.Rows()
}

func (clientProxy *GormClientProxy) Scan(dest interface{}) Client {
	c := clientProxy.clone()
	c.db = c.db.Scan(dest)
	return c
}

func (clientProxy *GormClientProxy) Pluck(column string, dest interface{}) Client {
	c := clientProxy.clone()
	c.db = c.db.Pluck(column, dest)
	return c
}

func (clientProxy *GormClientProxy) ScanRows(rows *sql.Rows, dest interface{}) error {
	return clientProxy.db.ScanRows(rows, dest)
}

func (clientProxy *GormClientProxy) Begin(opts ...*sql.TxOptions) Client {
	c := clientProxy.clone()
	c.db = c.db.Begin(opts...)
	return c
}

func (clientProxy *GormClientProxy) Commit() Client {
	c := clientProxy.clone()
	c.db = c.db.Commit()
	return c
}

func (clientProxy *GormClientProxy) Rollback() Client {
	c := clientProxy.clone()
	c.db = c.db.Rollback()
	return c
}

func (clientProxy *GormClientProxy) SavePoint(name string) Client {
	c := clientProxy.clone()
	c.db = c.db.SavePoint(name)
	return c
}

func (clientProxy *GormClientProxy) RollbackTo(name string) Client {
	c := clientProxy.clone()
	c.db = c.db.RollbackTo(name)
	return c
}

func (clientProxy *GormClientProxy) Exec(sql string, values ...interface{}) Client {
	c := clientProxy.clone()
	c.db = c.db.Exec(sql, values...)
	return c
}

func (clientProxy *GormClientProxy) GetDBMetas() ([]ITableMeta, error) {
	if clientProxy.driverOperator == nil {
		return nil, errors.New("driver not support get db metas")
	}
	return clientProxy.driverOperator.GetDBMetas()
}

func (clientProxy *GormClientProxy) clone() *GormClientProxy {
	return &GormClientProxy{
		driverOperator: clientProxy.driverOperator,
		client:         clientProxy.client,
		db:             clientProxy.db,
	}
}

func (clientProxy *GormClientProxy) Check() (ok bool) {
	sqlDB, err := clientProxy.db.DB()
	if err != nil {
		return false
	}
	err = sqlDB.Ping()
	if err != nil {
		return false
	}
	return true
}

func (clientProxy *GormClientProxy) Model(value interface{}) Client {
	c := clientProxy.clone()
	c.db = clientProxy.db.Model(value)
	return c
}

func (clientProxy *GormClientProxy) Where(query interface{}, args ...interface{}) Client {
	c := clientProxy.clone()
	c.db = c.db.Where(query, args...)
	return c
}

func (clientProxy *GormClientProxy) Table(name string, args ...interface{}) Client {
	c := clientProxy.clone()
	c.db = c.db.Table(name, args...)
	return c
}

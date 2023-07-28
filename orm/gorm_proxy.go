package orm

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"time"
)

type GormClientProxy struct {
	driverOperator DriverMetaOperator

	advanceQuery AdvanceQuery
	client       Client
	db           *gorm.DB
}

func (clientProxy *GormClientProxy) GetSchemaNames() ([]string, error) {
	if clientProxy.driverOperator == nil {
		return nil, errors.New("driver not support get db metas")
	}
	return clientProxy.driverOperator.GetSchemaNames()
}

func (clientProxy *GormClientProxy) GetSchemas() ([]SchemaMeta, error) {
	if clientProxy.driverOperator == nil {
		return nil, errors.New("driver not support get db metas")
	}
	return clientProxy.driverOperator.GetSchemas()
}

func (clientProxy *GormClientProxy) GetTableNames(schemaName string) ([]string, error) {
	if clientProxy.driverOperator == nil {
		return nil, errors.New("driver not support get db metas")
	}
	return clientProxy.driverOperator.GetTableNames(schemaName)
}

func (clientProxy *GormClientProxy) GetTables(schemaName string) ([]TableMeta, error) {
	if clientProxy.driverOperator == nil {
		return nil, errors.New("driver not support get db metas")
	}
	return clientProxy.driverOperator.GetTables(schemaName)
}

func (clientProxy *GormClientProxy) GetColumnNames(tableName, schemaName string) ([]string, error) {
	if clientProxy.driverOperator == nil {
		return nil, errors.New("driver not support get db metas")
	}
	return clientProxy.driverOperator.GetColumnNames(tableName, schemaName)
}

func (clientProxy *GormClientProxy) GetColumns(tableName, schemaName string) ([]ColumnMeta, error) {
	if clientProxy.driverOperator == nil {
		return nil, errors.New("driver not support get db metas")
	}
	return clientProxy.driverOperator.GetColumns(tableName, schemaName)
}

func (clientProxy *GormClientProxy) GetConstraintMeta(tableName, schemaName string) ([]ConstraintMeta, error) {
	if clientProxy.driverOperator == nil {
		return nil, errors.New("driver not support get db metas")
	}
	return clientProxy.driverOperator.GetConstraintMeta(tableName, schemaName)
}

func (clientProxy *GormClientProxy) CreateTableByMeta(tableMeta TableMeta) (err error) {
	return clientProxy.driverOperator.CreateTableByMeta(tableMeta)
}

func (clientProxy *GormClientProxy) JoinQueryTablesByCursor(joinQueryParams []JoinQueryParam, batchSize int, fc func(data []map[string]interface{})) error {
	c := clientProxy.clone()
	return c.advanceQuery.JoinQueryTablesByCursor(joinQueryParams, batchSize, fc)
}

func (clientProxy *GormClientProxy) QueryByCursor(tableName string, batchSize int, selectFields []string, orderBy []string, fc func(data []map[string]interface{})) error {
	c := clientProxy.clone()
	return c.advanceQuery.QueryByCursor(tableName, batchSize, selectFields, orderBy, fc)
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
	sqlDB.SetMaxIdleConns(config.MaxIdle)
	sqlDB.SetMaxOpenConns(config.MaxOpen)
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

	advanceQueryConstructor, ok := advanceQueryConstructorMap[db.Dialector.Name()]
	if ok {
		client.advanceQuery = advanceQueryConstructor(db)
	}

	return client, nil
}

func (clientProxy *GormClientProxy) Create(value interface{}) Client {
	c := clientProxy.clone()
	c.setDB(clientProxy.db.Create(value))
	return c
}

func (clientProxy *GormClientProxy) CreateInBatches(value interface{}, batchSize int) Client {
	c := clientProxy.clone()
	c.setDB(c.db.CreateInBatches(value, batchSize))
	return c
}

func (clientProxy *GormClientProxy) Save(value interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Save(value))
	return c
}

func (clientProxy *GormClientProxy) First(dest interface{}, conds ...interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.First(dest, conds...))
	return c
}

func (clientProxy *GormClientProxy) Take(dest interface{}, conds ...interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Take(dest, conds...))
	return c
}

func (clientProxy *GormClientProxy) Last(dest interface{}, conds ...interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Last(dest, conds...))
	return c
}

func (clientProxy *GormClientProxy) Find(dest interface{}, conds ...interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Find(dest, conds...))
	return c
}

func (clientProxy *GormClientProxy) FirstOrInit(dest interface{}, conds ...interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.FirstOrInit(dest, conds...))
	return c
}

func (clientProxy *GormClientProxy) FirstOrCreate(dest interface{}, conds ...interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.FirstOrCreate(dest, conds...))
	return c
}

func (clientProxy *GormClientProxy) Update(column string, value interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Update(column, value))
	return c
}

func (clientProxy *GormClientProxy) Updates(values interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Updates(values))
	return c
}

func (clientProxy *GormClientProxy) UpdateColumn(column string, value interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.UpdateColumn(column, value))
	return c
}

func (clientProxy *GormClientProxy) UpdateColumns(values interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.UpdateColumns(values))
	return c
}

func (clientProxy *GormClientProxy) Delete(value interface{}, conds ...interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Delete(value, conds...))
	return c
}

func (clientProxy *GormClientProxy) Count(count *int64) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Count(count))
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
	c.setDB(c.db.Scan(dest))
	return c
}

func (clientProxy *GormClientProxy) Pluck(column string, dest interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Pluck(column, dest))
	return c
}

func (clientProxy *GormClientProxy) ScanRows(rows *sql.Rows, dest interface{}) error {
	return clientProxy.db.ScanRows(rows, dest)
}

func (clientProxy *GormClientProxy) Begin(opts ...*sql.TxOptions) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Begin(opts...))
	return c
}

func (clientProxy *GormClientProxy) Commit() Client {
	c := clientProxy.clone()
	c.setDB(c.db.Commit())
	return c
}

func (clientProxy *GormClientProxy) Rollback() Client {
	c := clientProxy.clone()
	c.setDB(c.db.Rollback())
	return c
}

func (clientProxy *GormClientProxy) SavePoint(name string) Client {
	c := clientProxy.clone()
	c.setDB(c.db.SavePoint(name))
	return c
}

func (clientProxy *GormClientProxy) RollbackTo(name string) Client {
	c := clientProxy.clone()
	c.setDB(c.db.RollbackTo(name))
	return c
}

func (clientProxy *GormClientProxy) Exec(sql string, values ...interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Exec(sql, values...))
	return c
}

func (clientProxy *GormClientProxy) clone() *GormClientProxy {
	gp := &GormClientProxy{
		driverOperator: clientProxy.driverOperator,
		advanceQuery:   clientProxy.advanceQuery,
		client:         clientProxy.client,
		db:             clientProxy.db,
	}
	gp.setDB(clientProxy.db)
	return gp
}

func (clientProxy *GormClientProxy) setDB(db *gorm.DB) {
	clientProxy.db = db
	clientProxy.driverOperator.setDB(db)
	clientProxy.advanceQuery.setDB(db)
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
	c.setDB(clientProxy.db.Model(value))
	return c
}

func (clientProxy *GormClientProxy) Where(query interface{}, args ...interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Where(query, args...))
	return c
}

func (clientProxy *GormClientProxy) Limit(limit int) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Limit(limit))
	return c
}

func (clientProxy *GormClientProxy) Offset(offset int) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Offset(offset))
	return c
}

func (clientProxy *GormClientProxy) Order(value interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Order(value))
	return c
}

func (clientProxy *GormClientProxy) Table(name string, args ...interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Table(name, args...))
	return c
}

func (clientProxy *GormClientProxy) GetPrimaryKey(tableName string) (string, error) {
	return clientProxy.driverOperator.GetPrimaryKey(tableName)
}

func (clientProxy *GormClientProxy) Raw(sql string, values ...interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Raw(sql, values...))
	return c
}

func (clientProxy *GormClientProxy) Joins(sql string, args ...interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Joins(sql, args...))
	return c
}

func (clientProxy *GormClientProxy) Preload(sql string, args ...interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Preload(sql, args...))
	return c
}

func (clientProxy *GormClientProxy) Error() error {
	return clientProxy.db.Error
}

func (clientProxy *GormClientProxy) RowsAffected() int64 {
	return clientProxy.db.RowsAffected
}

func (clientProxy *GormClientProxy) Select(query interface{}, args ...interface{}) Client {
	c := clientProxy.clone()
	c.setDB(c.db.Select(query, args))
	return c
}

func (clientProxy *GormClientProxy) DB() *gorm.DB {
	return clientProxy.db
}

func (clientProxy *GormClientProxy) Transaction(fc func(tx Client) error, opts ...*sql.TxOptions) error {
	var c Client = clientProxy.clone()
	return c.DB().Transaction(func(tx *gorm.DB) error {
		c.setDB(tx)
		err := fc(c)
		return err
	}, opts...)
}

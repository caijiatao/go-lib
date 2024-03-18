package orm

import (
	"database/sql"
	"github.com/beltran/gohive"
	"github.com/pkg/errors"
	"golib/libs/logger"
	"gorm.io/gorm"
	"strconv"
)

const HIVE_DB_NAME = "hive"

type HiveClientProxy struct {
	driverOperator *hiveDriverMetaOperator

	config Config

	advanceQuery *hiveAdvanceQuery
	client       Client
	connection   *gohive.Connection
}

func NewHiveClientProxy(config Config) (*HiveClientProxy, error) {
	client := &HiveClientProxy{
		config: config,
	}
	err := client.initConnection()
	if err != nil {
		return nil, err
	}

	metaDriverOperatorConstructor, ok := driverMetaConstructorMap[HIVE_DB_NAME]
	if ok {
		driverOperator := metaDriverOperatorConstructor(client.connection)
		client.driverOperator = driverOperator.(*hiveDriverMetaOperator)
	}

	advanceQueryConstructor, ok := advanceQueryConstructorMap[HIVE_DB_NAME]
	if ok {
		advanceQuery := advanceQueryConstructor(client.connection)
		client.advanceQuery = advanceQuery.(*hiveAdvanceQuery)
	}

	return client, nil
}

func (hiveProxy *HiveClientProxy) initConnection() error {
	configuration := gohive.NewConnectConfiguration()
	configuration.Username = hiveProxy.config.SourceConfig.User
	configuration.Password = hiveProxy.config.SourceConfig.Password
	configuration.Database = hiveProxy.config.SourceConfig.DbName
	port, err := strconv.Atoi(hiveProxy.config.SourceConfig.Port)
	if err != nil {
		return errors.New("connection err")
	}
	hiveProxy.connection, err = gohive.Connect(hiveProxy.config.SourceConfig.Host, port, "NONE", configuration)
	if err != nil {
		return errors.New("connection err")
	}
	return nil
}

func (hiveProxy *HiveClientProxy) clone() *HiveClientProxy {
	clientProxy, err := NewHiveClientProxy(hiveProxy.config)
	if err != nil {
		return nil
	}
	return clientProxy
}

func (hiveProxy *HiveClientProxy) GetSchemaNames() ([]string, error) {
	if hiveProxy.driverOperator == nil {
		return nil, errors.New("driver not support get db metas")
	}
	return hiveProxy.driverOperator.GetSchemaNames()
}

func (hiveProxy *HiveClientProxy) GetSchemas() ([]SchemaMeta, error) {
	if hiveProxy.driverOperator == nil {
		return nil, errors.New("driver not support get db metas")
	}
	return hiveProxy.driverOperator.GetSchemas()
}

func (hiveProxy *HiveClientProxy) GetTableNames(schemaName string) ([]string, error) {
	if hiveProxy.driverOperator == nil {
		return nil, errors.New("driver not support get db metas")
	}
	return hiveProxy.driverOperator.GetTableNames(schemaName)
}

func (hiveProxy *HiveClientProxy) GetTables(schemaName string) ([]TableMeta, error) {
	if hiveProxy.driverOperator == nil {
		return nil, errors.New("driver not support get db metas")
	}
	return hiveProxy.driverOperator.GetTables(schemaName)
}

func (hiveProxy *HiveClientProxy) GetColumnNames(tableName string) ([]string, error) {
	if hiveProxy.driverOperator == nil {
		return nil, errors.New("driver not support get db metas")
	}
	return hiveProxy.driverOperator.GetColumnNames(tableName)
}

func (hiveProxy *HiveClientProxy) GetColumns(tableName, schemaName string) ([]ColumnMeta, error) {
	if hiveProxy.driverOperator == nil {
		return nil, errors.New("driver not support get db metas")
	}
	return hiveProxy.driverOperator.GetColumns(tableName, schemaName)
}

func (hiveProxy *HiveClientProxy) GetConstraintMeta(tableName, schemaName string) ([]ConstraintMeta, error) {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) GetPrimaryKey(tableName string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) CreateTableByMeta(tableMeta TableMeta) (err error) {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) QueryByCursor(tableName string, batchSize int, selectFields []string, where string, orderBy []string, fc func(data []map[string]interface{})) (map[string]interface{}, error) {
	hiveClientProxy := hiveProxy.clone()
	defer func() {
		err := hiveClientProxy.Close()
		if err != nil {
			logger.Error("close hive client error", err)
		}
	}()
	lastData, err := hiveClientProxy.advanceQuery.QueryByCursor(tableName, batchSize, selectFields, "", orderBy, fc)
	if err != nil {
		return nil, err
	}
	return lastData, nil
}

func (hiveProxy *HiveClientProxy) JoinQueryTablesByCursor(joinQueryParams []JoinQueryParam, batchSize int, fc func(data []map[string]interface{})) error {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) setDB(db *gorm.DB) {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Create(value interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) CreateInBatches(value interface{}, batchSize int) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Save(value interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) First(dest interface{}, conds ...interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Take(dest interface{}, conds ...interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Last(dest interface{}, conds ...interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Find(dest interface{}, conds ...interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) FirstOrInit(dest interface{}, conds ...interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) FirstOrCreate(dest interface{}, conds ...interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Update(column string, value interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Updates(values interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) UpdateColumn(column string, value interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) UpdateColumns(values interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Delete(value interface{}, conds ...interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Count(count *int64) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Row() *sql.Row {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Rows() (*sql.Rows, error) {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Scan(dest interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Pluck(column string, dest interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) ScanRows(rows *sql.Rows, dest interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Begin(opts ...*sql.TxOptions) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Commit() Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Rollback() Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) SavePoint(name string) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) RollbackTo(name string) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Exec(sql string, values ...interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Model(value interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Where(query interface{}, args ...interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Limit(limit int) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Offset(offset int) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Order(value interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Table(name string, args ...interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Raw(sql string, values ...interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Joins(sql string, args ...interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Preload(sql string, args ...interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Error() error {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) RowsAffected() int64 {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Select(query interface{}, args ...interface{}) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) DB() *gorm.DB {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Transaction(fc func(tx Client) error, opts ...*sql.TxOptions) error {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Group(name string) Client {
	//TODO implement me
	panic("implement me")
}

func (hiveProxy *HiveClientProxy) Check() (ok bool) {
	// hive 每次都是新的连接，失败会在初始化报错
	return true
}

func (hiveProxy *HiveClientProxy) Close() error {
	if hiveProxy.connection == nil {
		return nil
	}
	err := hiveProxy.connection.Close()
	if err != nil {
		return err
	}
	return nil
}

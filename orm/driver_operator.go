package orm

import (
	"gorm.io/gorm"
	"strings"
)

type driverMetaConstructor func(db *gorm.DB) DriverMetaOperator

type TableMeta struct {
	TableName   string
	Columns     []ColumnMeta
	Constraints []ConstraintMeta
}

type ColumnMeta struct {
	ColumnName string
	ColumnType string
	Comment    string
}

type ConstraintMeta struct {
}

const (
	postgreDialName = "postgres"
	mysqlDialName   = "mysql"
)

var (
	driverMetaConstructorMap = map[string]driverMetaConstructor{
		postgreDialName: newPostgreDriverMetaOperator,
		mysqlDialName:   newMysqlDriverMetaOperator,
	}
)

// DriverMetaOperator
// @Description: 驱动元信息
type DriverMetaOperator interface {
	GetDBMetas() ([]ITableMeta, error)
	GetPrimaryKey(tableName string) (string, error)
	CreateTableByMeta(tableMeta TableMeta) (err error)

	setDB(db *gorm.DB)
}

type mysqlDriverMetaOperator struct {
	db *gorm.DB
}

func (m *mysqlDriverMetaOperator) setDB(db *gorm.DB) {
	//TODO implement me
	panic("implement me")
}

func (m *mysqlDriverMetaOperator) GetPrimaryKey(tableName string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func newMysqlDriverMetaOperator(db *gorm.DB) DriverMetaOperator {
	return &mysqlDriverMetaOperator{
		db: db,
	}
}

func (m *mysqlDriverMetaOperator) GetDBMetas() ([]ITableMeta, error) {
	var tableNames []string
	result := m.db.Raw("SHOW TABLES").Scan(&tableNames)
	if result.Error != nil {
		return nil, result.Error
	}

	var tables []ITableMeta
	for _, tableName := range tableNames {
		var columns []gormColumnMeta
		result := m.db.Raw("DESCRIBE " + tableName).Scan(&columns)
		if result.Error != nil {
			return nil, result.Error
		}
		var icolumns []IColumnMeta
		for _, column := range columns {
			icolumns = append(icolumns, &gormColumnMeta{Field: column.Field, Type: column.Type})
		}
		tables = append(tables, &gormTableMeta{name: tableName, columns: icolumns})
	}

	return tables, nil
}

func (p *mysqlDriverMetaOperator) CreateTableByMeta(tableMeta TableMeta) (err error) {
	return nil
}

type postgreDriverMetaOperator struct {
	db *gorm.DB
}

func (p *postgreDriverMetaOperator) setDB(db *gorm.DB) {
	p.db = db
}

func newPostgreDriverMetaOperator(db *gorm.DB) DriverMetaOperator {
	return &postgreDriverMetaOperator{db: db}
}

func (p *postgreDriverMetaOperator) GetDBMetas() ([]ITableMeta, error) {
	var tableNames []string
	p.db.Raw("SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = 'public'").Scan(&tableNames)
	// Iterate over the table names and retrieve metadata for each table
	var tables []ITableMeta
	for _, tableName := range tableNames {
		var columns []IColumnMeta

		// Retrieve the column names and types for the current table
		rows, err := p.db.Raw("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = ?", tableName).Rows()
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var columnName string
			var dataType string
			if err := rows.Scan(&columnName, &dataType); err != nil {
				return nil, err
			}
			columns = append(columns, &gormColumnMeta{Field: columnName, Type: dataType})
		}

		tables = append(tables, &gormTableMeta{name: tableName, columns: columns})
	}

	return tables, nil
}

func (p *postgreDriverMetaOperator) GetPrimaryKey(tableName string) (string, error) {
	var primaryKey string
	tableNameSplit := strings.Split(tableName, ".")
	err := p.db.Raw(`SELECT column_name FROM information_schema.key_column_usage WHERE table_name = ? AND constraint_name LIKE '%_pkey'`, tableNameSplit[len(tableNameSplit)-1]).Scan(&primaryKey).Error
	if err != nil {
		return "", err
	}
	return primaryKey, nil
}

// CreateTableByMeta
//
//	@Description: TODO 通过源数据创建表结构
func (p *postgreDriverMetaOperator) CreateTableByMeta(tableMeta TableMeta) (err error) {
	return nil
}

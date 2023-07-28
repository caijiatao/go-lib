package orm

import (
	"gorm.io/gorm"
	"strings"
)

type driverMetaConstructor func(db *gorm.DB) DriverMetaOperator

type SchemaMeta struct {
	SchemaName string
	Tables     []TableMeta
}

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
	GetSchemaNames() ([]string, error)
	GetSchemas() ([]SchemaMeta, error)
	GetTableNames(schemaName string) ([]string, error)
	GetTables(schemaName string) ([]TableMeta, error)
	GetColumnNames(tableName, schemaName string) ([]string, error)
	GetColumns(tableName, schemaName string) ([]ColumnMeta, error)
	GetConstraintMeta(tableName, schemaName string) ([]ConstraintMeta, error)

	GetPrimaryKey(tableName string) (string, error)
	CreateTableByMeta(tableMeta TableMeta) (err error)

	setDB(db *gorm.DB)
}

type mysqlDriverMetaOperator struct {
	db *gorm.DB
}

func newMysqlDriverMetaOperator(db *gorm.DB) DriverMetaOperator {
	operator := new(mysqlDriverMetaOperator)
	operator.setDB(db)
	return operator

}

func (m *mysqlDriverMetaOperator) GetSchemaNames() ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mysqlDriverMetaOperator) GetSchemas() ([]SchemaMeta, error) {
	return make([]SchemaMeta, 0), nil
}

func (m *mysqlDriverMetaOperator) GetTableNames(schemaName string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mysqlDriverMetaOperator) GetTables(schemaName string) ([]TableMeta, error) {
	tables := make([]TableMeta, 0)
	tableNames := make([]string, 0)
	m.db.Raw("SHOW TABLES").Scan(tableNames)
	for _, name := range tableNames {
		cols, err := m.GetColumns(name, schemaName)
		if err != nil {
			return tables, err
		}
		cons, err := m.GetConstraintMeta(name, schemaName)
		if err != nil {
			return tables, err
		}
		tables = append(tables, TableMeta{
			TableName:   name,
			Columns:     cols,
			Constraints: cons,
		})
	}
	return tables, nil
}

func (m *mysqlDriverMetaOperator) GetColumnNames(tableName, schemaName string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mysqlDriverMetaOperator) GetColumns(tableName, schemaName string) ([]ColumnMeta, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mysqlDriverMetaOperator) GetConstraintMeta(tableName, schemaName string) ([]ConstraintMeta, error) {
	//TODO implement me
	return nil, nil
}

func (m *mysqlDriverMetaOperator) CreateTableByMeta(tableMeta TableMeta) (err error) {
	//TODO implement me
	panic("implement me")
}

func (m *mysqlDriverMetaOperator) setDB(db *gorm.DB) {
	m.db = db
}

func (m *mysqlDriverMetaOperator) GetPrimaryKey(tableName string) (string, error) {
	//TODO implement me
	panic("implement me")
}

type postgreDriverMetaOperator struct {
	db *gorm.DB
}

func newPostgreDriverMetaOperator(db *gorm.DB) DriverMetaOperator {
	operator := new(postgreDriverMetaOperator)
	operator.setDB(db)
	return operator
}

func (p *postgreDriverMetaOperator) GetSchemaNames() ([]string, error) {
	sql := `	
SELECT schema_name
FROM information_schema.schemata
WHERE schema_name NOT LIKE 'pg_%'
  AND schema_name <> 'information_schema';`
	schemaNames := make([]string, 0)
	rows, err := p.db.Raw(sql).Rows()
	if err != nil {
		return schemaNames, err
	}
	defer rows.Close()
	for rows.Next() {
		var schemaName string
		err = rows.Scan(&schemaName)
		if err != nil {
			return schemaNames, err
		}
		schemaNames = append(schemaNames, schemaName)
	}
	return schemaNames, err
}

func (p *postgreDriverMetaOperator) GetSchemas() ([]SchemaMeta, error) {
	schemas := make([]SchemaMeta, 0)
	SchemaNames, err := p.GetSchemaNames()
	if err != nil {
		return schemas, err
	}
	for _, name := range SchemaNames {
		tables, err := p.GetTables(name)
		if err != nil {
			return schemas, err
		}
		schemas = append(schemas, SchemaMeta{
			SchemaName: name,
			Tables:     tables,
		})
	}
	return schemas, err
}

func (p *postgreDriverMetaOperator) GetTableNames(schemaName string) ([]string, error) {
	tableNames := make([]string, 0)
	rows, err := p.db.Raw("select table_name from information_schema.tables where table_schema = ?", schemaName).Rows()
	if err != nil {
		return tableNames, err
	}
	defer rows.Close()
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			return tableNames, err
		}
		tableNames = append(tableNames, tableName)
	}
	return tableNames, err
}

func (p *postgreDriverMetaOperator) GetTables(schemaName string) ([]TableMeta, error) {
	tables := make([]TableMeta, 0)
	tableNames, err := p.GetTableNames(schemaName)
	if err != nil {
		return tables, err
	}
	for _, name := range tableNames {
		cols, err := p.GetColumns(name, schemaName)
		if err != nil {
			return tables, err
		}
		constraints, err := p.GetConstraintMeta(name, schemaName)
		if err != nil {
			return tables, err
		}
		tables = append(tables, TableMeta{
			TableName:   name,
			Columns:     cols,
			Constraints: constraints,
		})
	}
	return tables, err
}

func (p *postgreDriverMetaOperator) GetColumnNames(tableName, schemaName string) ([]string, error) {
	columnNames := make([]string, 0)
	err := p.db.Raw("SELECT column_name FROM information_schema.Columns WHERE table_schema = ? AND table_name = ?", schemaName, tableName).Scan(&columnNames).Error
	return columnNames, err

}

func (p *postgreDriverMetaOperator) GetColumns(tableName, schemaName string) ([]ColumnMeta, error) {
	columns := make([]ColumnMeta, 0)
	rows, err := p.db.Raw("SELECT column_name, data_type  FROM information_schema.Columns WHERE table_schema = ? AND table_name = ?", schemaName, tableName).Rows()
	if err != nil {
		return columns, err
	}
	for rows.Next() {
		column := ColumnMeta{}
		err := rows.Scan(&column.ColumnName, &column.ColumnType)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}
	return columns, nil
}

func (p *postgreDriverMetaOperator) GetConstraintMeta(tableName, schemaName string) ([]ConstraintMeta, error) {
	// fixme
	return nil, nil
}

func (p *postgreDriverMetaOperator) setDB(db *gorm.DB) {
	p.db = db
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

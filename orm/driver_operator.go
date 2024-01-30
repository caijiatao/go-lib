package orm

import (
	"context"
	"fmt"
	"github.com/beltran/gohive"
	"golib/orm/postgres_util"
	"gorm.io/gorm"
	"strings"
)

type driverMetaConstructor func(db interface{}) DriverMetaOperator

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
	hiveDialName    = "hive"
)

var (
	driverMetaConstructorMap = map[string]driverMetaConstructor{
		postgreDialName: newPostgreDriverMetaOperator,
		mysqlDialName:   newMysqlDriverMetaOperator,
		hiveDialName:    newHiveDriverMetaOperator,
	}
)

// DriverMetaOperator
// @Description: 驱动元信息
type DriverMetaOperator interface {
	GetSchemaNames() ([]string, error)
	GetSchemas() ([]SchemaMeta, error)
	GetTableNames(schemaName string) ([]string, error)
	GetTables(schemaName string) ([]TableMeta, error)
	GetColumnNames(tableName string) ([]string, error)
	GetColumns(tableName, schemaName string) ([]ColumnMeta, error)
	GetConstraintMeta(tableName, schemaName string) ([]ConstraintMeta, error)

	GetPrimaryKey(tableName string) (string, error)
	CreateTableByMeta(tableMeta TableMeta) (err error)

	setDB(db *gorm.DB)
}

type hiveDriverMetaOperator struct {
	cn *gohive.Connection
}

func newHiveDriverMetaOperator(db interface{}) DriverMetaOperator {
	operator := new(hiveDriverMetaOperator)
	operator.cn = db.(*gohive.Connection)
	return operator
}

func (i *hiveDriverMetaOperator) GetSchemaNames() ([]string, error) {
	return []string{}, nil
}

func (i *hiveDriverMetaOperator) GetSchemas() ([]SchemaMeta, error) {
	schemas := make([]SchemaMeta, 0)
	tables, err := i.GetTables("")
	if err != nil {
		return schemas, err
	}
	schemas = append(schemas, SchemaMeta{
		SchemaName: "",
		Tables:     tables,
	})
	return schemas, err
}

func (i *hiveDriverMetaOperator) GetTableNames(schemaName string) ([]string, error) {
	ctx := context.Background()
	cursor := i.cn.Cursor()
	cursor.Exec(ctx, "show tables")
	if cursor.Err != nil {
		return nil, cursor.Err
	}
	dataList := make([]string, 0)
	for cursor.HasMore(ctx) {
		resultMap := cursor.RowMap(ctx)
		if cursor.Err != nil {
			return nil, cursor.Err
		}
		for _, value := range resultMap {
			dataList = append(dataList, value.(string))
		}
	}
	return dataList, nil
}

func (i *hiveDriverMetaOperator) GetTables(schemaName string) ([]TableMeta, error) {
	tables := make([]TableMeta, 0)
	tableNames, err := i.GetTableNames(schemaName)
	if err != nil {
		return tables, err
	}
	for _, name := range tableNames {
		cols, err := i.GetColumns(name, schemaName)
		if err != nil {
			return tables, err
		}
		//constraints, err := i.GetConstraintMeta(name, schemaName)
		//if err != nil {
		//	return tables, err
		//}
		tables = append(tables, TableMeta{
			TableName: name,
			Columns:   cols,
			//Constraints: constraints,
		})
	}
	return tables, err
}

func (i *hiveDriverMetaOperator) GetColumnNames(tableName string) ([]string, error) {
	res := make([]string, 0)
	columnMetas, err := i.GetColumns(tableName, "")
	if err != nil {
		return nil, err
	}
	for _, meta := range columnMetas {
		res = append(res, meta.ColumnName)
	}
	return res, nil
}

func (i *hiveDriverMetaOperator) GetColumns(tableName, schemaName string) ([]ColumnMeta, error) {
	ctx := context.Background()
	cursor := i.cn.Cursor()
	cursor.Exec(ctx, fmt.Sprintf("describe %s", tableName))
	if cursor.Err != nil {
		return nil, cursor.Err
	}
	dataList := make([]ColumnMeta, 0)
	dulMap := make(map[string]bool)
	for cursor.HasMore(ctx) {
		resultMap := cursor.RowMap(ctx)
		if cursor.Err != nil {
			return nil, cursor.Err
		}
		columnMeta := ColumnMeta{}
		for key, value := range resultMap {
			switch key {
			case "col_name":
				columnMeta.ColumnName = value.(string)
			case "data_type":
				columnType := fmt.Sprintf("%v", value)
				if value == nil || fmt.Sprintf("%v", value) == "data_type" {
					columnType = ""
				}
				columnMeta.ColumnType = columnType
			}
		}
		if columnMeta.ColumnType != "" && !dulMap[columnMeta.ColumnName] {
			dataList = append(dataList, columnMeta)
			dulMap[columnMeta.ColumnName] = true
		}
	}
	return dataList, nil
}

func (i *hiveDriverMetaOperator) GetConstraintMeta(tableName, schemaName string) ([]ConstraintMeta, error) {
	//TODO implement me
	panic("implement me")
}

func (i *hiveDriverMetaOperator) GetPrimaryKey(tableName string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (i *hiveDriverMetaOperator) CreateTableByMeta(tableMeta TableMeta) (err error) {
	//TODO implement me
	panic("implement me")
}

func (i *hiveDriverMetaOperator) setDB(db *gorm.DB) {
	//TODO implement me
	panic("implement me")
}

type mysqlDriverMetaOperator struct {
	db *gorm.DB
}

func newMysqlDriverMetaOperator(db interface{}) DriverMetaOperator {
	operator := new(mysqlDriverMetaOperator)
	operator.setDB(db.(*gorm.DB))
	return operator
}

func (m *mysqlDriverMetaOperator) GetSchemaNames() ([]string, error) {
	return []string{}, nil
}

func (m *mysqlDriverMetaOperator) GetSchemas() ([]SchemaMeta, error) {
	schemas := make([]SchemaMeta, 0)
	tables, err := m.GetTables("")
	if err != nil {
		return schemas, err
	}
	schemas = append(schemas, SchemaMeta{
		SchemaName: "",
		Tables:     tables,
	})
	return schemas, err
}

func (m *mysqlDriverMetaOperator) GetTableNames(schemaName string) ([]string, error) {
	tableNames := make([]string, 0)
	rows, err := m.db.Raw("show tables").Rows()
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

func (m *mysqlDriverMetaOperator) GetTables(schemaName string) ([]TableMeta, error) {
	tables := make([]TableMeta, 0)
	tableNames := make([]string, 0)
	m.db.Raw("SHOW TABLES").Scan(tableNames)
	for _, name := range tableNames {
		cols, err := m.GetColumns(name, schemaName)
		if err != nil {
			return tables, err
		}
		//cons, err := m.GetConstraintMeta(name, schemaName)
		//if err != nil {
		//	return tables, err
		//}
		tables = append(tables, TableMeta{
			TableName: name,
			Columns:   cols,
			//Constraints: cons,
		})
	}
	return tables, nil
}

func (m *mysqlDriverMetaOperator) GetColumnNames(tableName string) ([]string, error) {
	columnNames := make([]map[string]interface{}, 0)
	res := make([]string, 0)
	err := m.db.Raw(fmt.Sprintf("show columns from %s", tableName)).Scan(&columnNames).Error
	for _, columnMap := range columnNames {
		for key, value := range columnMap {
			if key == "Field" {
				res = append(res, value.(string))
			}
		}
	}
	return res, err
}

func (m *mysqlDriverMetaOperator) GetColumns(tableName, schemaName string) ([]ColumnMeta, error) {
	columnNames := make([]map[string]interface{}, 0)
	res := make([]ColumnMeta, 0)
	err := m.db.Raw(fmt.Sprintf("show columns from %s", tableName)).Scan(&columnNames).Error
	for _, columnMap := range columnNames {
		meta := ColumnMeta{}
		for key, value := range columnMap {
			switch key {
			case "Field":
				meta.ColumnName = value.(string)
			case "Type":
				meta.ColumnType = value.(string)
			}
		}
		res = append(res, meta)
	}
	return res, err
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

func newPostgreDriverMetaOperator(db interface{}) DriverMetaOperator {
	operator := new(postgreDriverMetaOperator)
	operator.setDB(db.(*gorm.DB))
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

func (p *postgreDriverMetaOperator) GetColumnNames(tableName string) ([]string, error) {
	schemaName, tableName := postgres_util.GetSchemaAndTableName(tableName)
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

type yashanDriverMetaOperator struct {
	db *gorm.DB
}

func newYashanDriverMetaOperator(db interface{}) DriverMetaOperator {
	return &yashanDriverMetaOperator{
		db: db.(*gorm.DB),
	}
}

func (y *yashanDriverMetaOperator) GetSchemaNames() ([]string, error) {
	return []string{}, nil
}

func (y *yashanDriverMetaOperator) GetSchemas() ([]SchemaMeta, error) {
	schemas := make([]SchemaMeta, 0)
	tables, err := y.GetTables("")
	if err != nil {
		return schemas, err
	}
	schemas = append(schemas, SchemaMeta{
		SchemaName: "",
		Tables:     tables,
	})
	return schemas, err
}

func (y *yashanDriverMetaOperator) GetTableNames(schemaName string) ([]string, error) {
	tableNames := make([]string, 0)
	err := y.db.Raw("select table_name from user_tables").Scan(&tableNames).Error
	if err != nil {
		return tableNames, err
	}
	return tableNames, nil
}

func (y *yashanDriverMetaOperator) GetTables(schemaName string) ([]TableMeta, error) {
	tables := make([]TableMeta, 0)
	tableNames := make([]string, 0)
	y.db.Raw("select table_name from user_tables").Scan(&tableNames)
	for _, name := range tableNames {
		cols, err := y.GetColumns(name, schemaName)
		if err != nil {
			return tables, err
		}
		//cons, err := m.GetConstraintMeta(name, schemaName)
		//if err != nil {
		//	return tables, err
		//}
		tables = append(tables, TableMeta{
			TableName: name,
			Columns:   cols,
			//Constraints: cons,
		})
	}
	return tables, nil
}

func (y *yashanDriverMetaOperator) GetColumnNames(tableName string) ([]string, error) {
	columnNames := make([]string, 0)
	err := y.db.Raw(fmt.Sprintf("select column_name from user_tab_columns where table_name = '%s'", tableName)).Scan(&columnNames).Error
	if err != nil {
		return columnNames, err
	}
	return columnNames, nil
}

func (y *yashanDriverMetaOperator) GetColumns(tableName, schemaName string) ([]ColumnMeta, error) {
	columns := make([]ColumnMeta, 0)
	rows, err := y.db.Raw("SELECT column_name, data_type  FROM user_tab_columns WHERE table_name = ?", tableName).Rows()
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

func (y *yashanDriverMetaOperator) GetConstraintMeta(tableName, schemaName string) ([]ConstraintMeta, error) {
	//TODO implement me
	panic("implement me")
}

func (y *yashanDriverMetaOperator) GetPrimaryKey(tableName string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (y *yashanDriverMetaOperator) CreateTableByMeta(tableMeta TableMeta) (err error) {
	//TODO implement me
	panic("implement me")
}

func (y *yashanDriverMetaOperator) setDB(db *gorm.DB) {
	y.db = db
}

package orm

import (
	"context"
	"fmt"
	"github.com/beltran/gohive"
	"golib/libs/orm/postgres_util"
	"gorm.io/gorm"
	"strings"
	"time"
)

type AdvanceQuery interface {
	QueryByCursor(tableName string, batchSize int, selectFields []string, where string, orderBy []string, fc func(data []map[string]interface{})) (map[string]interface{}, error)
	JoinQueryTablesByCursor(joinQueryParams []JoinQueryParam, batchSize int, fc func(data []map[string]interface{})) error
	setDB(db *gorm.DB)
}

type advanceQueryConstructorFunc func(connection interface{}) AdvanceQuery

var (
	advanceQueryConstructorMap = map[string]advanceQueryConstructorFunc{
		postgreDialName: newPostgreQuery,
		mysqlDialName:   newMysqlAdvanceQuery,
		hiveDialName:    newHiveAdvanceQuery,
	}
)

type postgreQuery struct {
	db *gorm.DB
}

func (p *postgreQuery) setDB(db *gorm.DB) {
	p.db = db
}

func newPostgreQuery(db interface{}) AdvanceQuery {
	return &postgreQuery{
		db: db.(*gorm.DB),
	}
}

func (p *postgreQuery) getTableName(rawTableName string) string {
	splitTable := strings.Split(rawTableName, ".")
	return splitTable[len(splitTable)-1]
}

func (p *postgreQuery) wrapQuotesTableName(tableName string) string {
	splitTable := strings.Split(tableName, ".")
	tableNames := make([]string, 0)
	for _, s := range splitTable {
		firstChar := s[0]
		lastChar := s[len(s)-1]

		if firstChar == '"' && lastChar == '"' {
			tableNames = append(tableNames, s)
		} else {
			tableNames = append(tableNames, fmt.Sprintf(`"%s"`, s))
		}
	}
	return strings.Join(tableNames, ".")
}

func getSelectFieldsString(selectFields []string) string {
	// 适配增量数据，查询所有的key
	selectFieldStr := "*"
	if len(selectFields) > 0 {
		selectFieldStr = strings.Join(selectFields, ",")
	}
	return selectFieldStr
}

func getOrderByString(orderBy []string) string {
	orderByStr := ""
	if len(orderBy) > 0 {
		orderByStr = strings.Join(orderBy, ",")
		orderByStr = " ORDER BY " + orderByStr
	}
	return orderByStr
}

func (p *postgreQuery) getOrderByString(orderBy []string) string {
	orderByStr := ""
	if len(orderBy) > 0 {
		newOrderBy := make([]string, 0, len(orderBy))
		for i := 0; i < len(orderBy); i++ {
			parts := strings.Fields(orderBy[i])
			if len(parts) != 2 {
				continue
			}
			quotedColumn := fmt.Sprintf(`"%s" %s nulls first `, parts[0], parts[1])
			newOrderBy = append(newOrderBy, quotedColumn)
		}
		orderByStr = strings.Join(newOrderBy, ",")
		orderByStr = " ORDER BY " + orderByStr
	}
	return orderByStr
}

func (p *postgreQuery) getSelectFieldsString(selectFields []string) string {
	// 适配增量数据。查询所有的key。
	selectFieldStr := "*"
	if len(selectFields) > 0 {
		newSelectFields := make([]string, 0, len(selectFieldStr))
		for i := 0; i < len(selectFields); i++ {
			newSelectFields = append(newSelectFields, fmt.Sprintf(`"%s"`, selectFields[i]))
		}
		selectFieldStr = strings.Join(newSelectFields, ",")
	}
	return selectFieldStr
}

func (p *postgreQuery) QueryByCursor(tableName string, batchSize int, selectFields []string, where string, orderBy []string, fc func(data []map[string]interface{})) (map[string]interface{}, error) {
	splitTableName := p.getTableName(tableName)
	// 避免同个表名并发创建出相同的游标名
	cursorName := fmt.Sprintf(`"%s_cursor_%d"`, splitTableName, time.Now().UnixNano())
	tableName = p.wrapQuotesTableName(tableName)
	tx := p.db.Begin()
	defer func() {
		closeCursor := fmt.Sprintf("CLOSE %s", cursorName)
		tx = tx.Exec(closeCursor)
		tx = tx.Commit()
	}()
	selectFieldStr := p.getSelectFieldsString(selectFields)
	orderByStr := p.getOrderByString(orderBy)

	createCursor := fmt.Sprintf("DECLARE %s CURSOR FOR SELECT %s FROM %s %s %s", cursorName, selectFieldStr, tableName, where, orderByStr)

	tx = tx.Exec(createCursor)
	err := tx.Error
	if err != nil {
		return nil, err
	}
	queryCursor := fmt.Sprintf("FETCH FORWARD %d FROM %s", batchSize, cursorName)

	lastData := make(map[string]interface{})
	for {
		data := make([]map[string]interface{}, 0)
		tx = tx.Raw(queryCursor).Scan(&data)
		err := tx.Error
		if err != nil {
			return nil, err
		}
		if len(data) == 0 {
			return lastData, nil
		}
		lastData = data[len(data)-1]
		fc(data)
	}
}

type JoinQueryParam struct {
	TableName string
	JoinCol   string
	Columns   []string
}

func (j *JoinQueryParam) getFullJoinSql(joinParam JoinQueryParam) string {
	joinTemplate := "FULL JOIN %s ON %s.%s = " + joinParam.TableName + "." + joinParam.JoinCol
	return fmt.Sprintf(joinTemplate, j.TableName, j.TableName, j.JoinCol)
}

func (j *JoinQueryParam) getQueryColSql() string {
	if len(j.Columns) == 0 {
		return fmt.Sprintf("%s.*", j.TableName)
	}
	colsSql := make([]string, 0)
	for _, column := range j.Columns {
		colsSql = append(colsSql, fmt.Sprintf("%s.%s AS %s", j.TableName, column, postgres_util.GetColAlias(j.TableName, column)))
	}

	return strings.Join(colsSql, ",")
}

func (p *postgreQuery) getJoinQueryCols(joinQueryParam []JoinQueryParam) string {
	colsSql := make([]string, 0)
	for _, param := range joinQueryParam {
		colSql := param.getQueryColSql()
		colsSql = append(colsSql, colSql)
	}
	return strings.Join(colsSql, ",")
}

func (p *postgreQuery) JoinQueryTablesByCursor(joinQueryParams []JoinQueryParam, batchSize int, fc func(data []map[string]interface{})) error {
	tableNames := make([]string, len(joinQueryParams))
	for i, param := range joinQueryParams {
		tableNames[i] = p.getTableName(param.TableName)
	}
	cursorName := fmt.Sprintf("%s_cursor", strings.Join(tableNames, "_"))
	p.db = p.db.Begin()
	defer func() {
		closeCursor := fmt.Sprintf("CLOSE %s", cursorName)
		p.db = p.db.Exec(closeCursor)
		p.db = p.db.Commit()
	}()
	querySql := fmt.Sprintf("SELECT %s FROM %s", p.getJoinQueryCols(joinQueryParams), joinQueryParams[0].TableName)

	joinSqls := []string{querySql}
	for i := 1; i < len(joinQueryParams); i++ {
		joinSqls = append(joinSqls, joinQueryParams[i].getFullJoinSql(joinQueryParams[0]))
	}
	querySql = strings.Join(joinSqls, " ")

	createCursor := fmt.Sprintf("DECLARE %s CURSOR FOR %s", cursorName, querySql)

	p.db = p.db.Exec(createCursor)

	queryCursor := fmt.Sprintf("FETCH FORWARD %d FROM %s", batchSize, cursorName)

	for {
		data := make([]map[string]interface{}, 0)
		p.db = p.db.Raw(queryCursor).Scan(&data)
		if len(data) == 0 {
			return nil
		}
		fc(data)
	}
}

type mysqlAdvanceQuery struct {
	db *gorm.DB
}

func newMysqlAdvanceQuery(db interface{}) AdvanceQuery {
	m := &mysqlAdvanceQuery{
		db: db.(*gorm.DB),
	}
	return m
}

func (m *mysqlAdvanceQuery) wrapQuotesTableName(tableName string) string {
	if strings.HasPrefix(tableName, "`") && strings.HasSuffix(tableName, "`") {
		return tableName
	}
	return fmt.Sprintf("`%s`", tableName)
}

func (m *mysqlAdvanceQuery) QueryByCursor(tableName string, batchSize int, selectFields []string, where string, orderBy []string, fc func(data []map[string]interface{})) (map[string]interface{}, error) {
	return m.QueryByRows(tableName, batchSize, selectFields, where, orderBy, fc)
	//var (
	//	queryTableSqlFormat = "SELECT %s FROM %s %s %s LIMIT %d OFFSET %d"
	//	offset              = 0
	//	selectFieldsString  = getSelectFieldsString(selectFields)
	//	orderByString       = getOrderByString(orderBy)
	//)
	//lastData := make(map[string]interface{})
	//for {
	//	data := make([]map[string]interface{}, 0)
	//	querySQL := fmt.Sprintf(queryTableSqlFormat, selectFieldsString, m.wrapQuotesTableName(tableName), where, orderByString, batchSize, offset)
	//	m.db = m.db.Raw(querySQL).Scan(&data)
	//	err := m.db.Error
	//	if err != nil {
	//		return nil, err
	//	}
	//	if len(data) == 0 {
	//		return lastData, nil
	//	}
	//	lastData = data[len(data)-1]
	//	fc(data)
	//	offset += len(data)
	//}
}

func (m *mysqlAdvanceQuery) QueryByRows(tableName string, batchSize int, selectFields []string, where string, orderBy []string, fc func(data []map[string]interface{})) (map[string]interface{}, error) {
	var (
		queryTableSqlFormat = "SELECT %s FROM %s %s %s"
		offset              = 0
		selectFieldsString  = getSelectFieldsString(selectFields)
		orderByString       = getOrderByString(orderBy)
	)
	querySQL := fmt.Sprintf(queryTableSqlFormat, selectFieldsString, m.wrapQuotesTableName(tableName), where, orderByString)
	rows, err := m.db.Raw(querySQL).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	lastData := make(map[string]interface{})
	for {
		data := make([]map[string]interface{}, 0, batchSize)
		for rows.Next() {
			item := make(map[string]interface{})
			scanRowsErr := m.db.ScanRows(rows, &item)
			if scanRowsErr != nil {
				return nil, scanRowsErr
			}
			data = append(data, item)
			if len(data) == batchSize {
				break
			}
		}
		if len(data) == 0 {
			return lastData, nil
		}
		lastData = data[len(data)-1]
		fc(data)
		offset += len(data)
	}
}

func (m *mysqlAdvanceQuery) JoinQueryTablesByCursor(joinQueryParams []JoinQueryParam, batchSize int, fc func(data []map[string]interface{})) error {
	//TODO implement me
	panic("implement me")
}

func (m *mysqlAdvanceQuery) setDB(db *gorm.DB) {
	m.db = db
}

type hiveAdvanceQuery struct {
	connection *gohive.Connection
}

func newHiveAdvanceQuery(db interface{}) AdvanceQuery {
	return &hiveAdvanceQuery{
		connection: db.(*gohive.Connection),
	}
}

func (m *hiveAdvanceQuery) wrapQuotesTableName(tableName string) string {
	tableNames := strings.Split(tableName, ".")
	wrapTableNames := make([]string, 0, len(tableNames))
	for _, name := range tableNames {
		if strings.HasPrefix(name, "`") && strings.HasSuffix(name, "`") {
			wrapTableNames = append(wrapTableNames, name)
		} else {
			name = fmt.Sprintf("`%s`", name)
			wrapTableNames = append(wrapTableNames, name)
		}
	}
	return strings.Join(wrapTableNames, ".")
}

// QueryByCursor
//
//	@Description: hive 游标查询所有数据，不使用order by
func (h *hiveAdvanceQuery) QueryByCursor(tableName string, batchSize int, selectFields []string, where string, orderBy []string, fc func(data []map[string]interface{})) (map[string]interface{}, error) {
	querySQLTemplate := "SELECT %s FROM %s LIMIT %d OFFSET %d"
	offset := 0
	selectFieldsString := getSelectFieldsString(selectFields)
	for {
		ctx := context.Background()
		cursor := h.connection.Cursor()
		querySQL := fmt.Sprintf(querySQLTemplate, selectFieldsString, h.wrapQuotesTableName(tableName), batchSize, offset)

		cursor.Exec(ctx, querySQL)
		data := make([]map[string]interface{}, 0, batchSize)
		for cursor.HasMore(ctx) {
			if cursor.Err != nil {
				return nil, cursor.Err
			}
			data = append(data, cursor.RowMap(ctx))
		}
		if len(data) == 0 {
			break
		}
		fc(data)
		offset += len(data)
	}

	return nil, nil
}

func (h *hiveAdvanceQuery) JoinQueryTablesByCursor(joinQueryParams []JoinQueryParam, batchSize int, fc func(data []map[string]interface{})) error {
	//TODO implement me
	panic("implement me")
}

func (h *hiveAdvanceQuery) setDB(db *gorm.DB) {
	return
}

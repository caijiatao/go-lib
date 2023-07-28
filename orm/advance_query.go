package orm

import (
	"fmt"
	"golib/orm/postgres_util"
	"gorm.io/gorm"
	"strings"
)

type AdvanceQuery interface {
	QueryByCursor(tableName string, batchSize int, selectFields []string, orderBy []string, fc func(data []map[string]interface{})) error
	JoinQueryTablesByCursor(joinQueryParams []JoinQueryParam, batchSize int, fc func(data []map[string]interface{})) error
	setDB(db *gorm.DB)
}

type advanceQueryConstructorFunc func(db *gorm.DB) AdvanceQuery

var (
	advanceQueryConstructorMap = map[string]advanceQueryConstructorFunc{
		postgreDialName: newPostgreQuery,
	}
)

type postgreQuery struct {
	db *gorm.DB
}

func (p *postgreQuery) setDB(db *gorm.DB) {
	p.db = db
}

func newPostgreQuery(db *gorm.DB) AdvanceQuery {
	return &postgreQuery{db: db}
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

func (p *postgreQuery) QueryByCursor(tableName string, batchSize int, selectFields []string, orderBy []string, fc func(data []map[string]interface{})) error {
	splitTableName := p.getTableName(tableName)
	cursorName := fmt.Sprintf("%s_cursor", splitTableName)
	tableName = p.wrapQuotesTableName(tableName)
	p.db = p.db.Begin()
	defer func() {
		closeCursor := fmt.Sprintf("CLOSE %s", cursorName)
		p.db = p.db.Exec(closeCursor)
		p.db = p.db.Commit()
	}()
	selectFieldStr := "*"
	if len(selectFields) > 0 {
		selectFieldStr = strings.Join(selectFields, ",")
	}
	orderByStr := ""
	if len(orderBy) > 0 {
		orderByStr = strings.Join(orderBy, ",")
		orderByStr = " ORDER BY " + orderByStr
	}

	createCursor := fmt.Sprintf("DECLARE %s CURSOR FOR SELECT %s FROM %s %s", cursorName, selectFieldStr, tableName, orderByStr)

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

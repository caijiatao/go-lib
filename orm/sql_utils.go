package orm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

const (
	// 扩展条件
	operateNe  = "ne"
	operateLt  = "lt"
	operateLte = "lte"
	operateGt  = "gt"
	operateGte = "gte"
)

var (
	operateValueMap = map[string]string{
		operateNe:  "!=",
		operateLt:  "<",
		operateLte: "<=",
		operateGt:  ">",
		operateGte: ">=",
	}
)

var defaultPlanDetailsOrderBy = []string{"(id + 0) desc"}

type fieldInjectWhereCondition func(engine Client, col string, value conditionColumnInfo) Client

func injectInt(engine Client, col string, columnInfo conditionColumnInfo) Client {
	if columnInfo.v.Int() == 0 {
		return engine
	}
	engine = engine.Where(fmt.Sprintf("%s %s ?", col, columnInfo.getOperate()), columnInfo.v.Int())
	return engine
}

func injectUint(engine Client, col string, columnInfo conditionColumnInfo) Client {
	if columnInfo.v.Uint() == 0 {
		return engine
	}
	engine = engine.Where(fmt.Sprintf("%s %s ?", col, columnInfo.getOperate()), columnInfo.v.Uint())
	return engine
}

func injectString(engine Client, col string, columnInfo conditionColumnInfo) Client {
	if columnInfo.v.String() == "" {
		return engine
	}
	engine = engine.Where(fmt.Sprintf("%s %s ?", col, columnInfo.getOperate()), columnInfo.v.String())
	return engine
}

func injectComplex(engine Client, col string, columnInfo conditionColumnInfo) Client {
	if columnInfo.v.Complex() == 0 {
		return engine
	}
	engine = engine.Where(fmt.Sprintf("%s %s ?", col, columnInfo.getOperate()), columnInfo.v.Complex())
	return engine
}

func injectSlice(engine Client, col string, columnInfo conditionColumnInfo) Client {
	if columnInfo.v.IsNil() {
		return engine
	}
	if columnInfo.v.Len() == 0 {
		return engine
	}
	kind := columnInfo.v.Index(0).Type().Kind()
	values := make([]interface{}, 0)
	for i := 0; i < columnInfo.v.Len(); i++ {
		if convFunc, ok := supportArrayFieldType[kind]; ok {
			convValue, isEmpty := convFunc(columnInfo.v, i)
			if !isEmpty {
				values = append(values, convValue)
			}
		}
	}
	if len(values) <= 0 {
		return engine
	}
	engine = engine.Where(fmt.Sprintf("%s %s (?)", col, columnInfo.getOperate()), values)
	return engine
}

type convArrayItem func(value reflect.Value, idx int) (interface{}, bool)

func convIntArrayItem(value reflect.Value, idx int) (interface{}, bool) {
	return value.Index(idx).Int(), value.Index(idx).Int() == 0
}
func convUIntArrayItem(value reflect.Value, idx int) (interface{}, bool) {
	return value.Index(idx).Uint(), value.Index(idx).Uint() == 0
}
func convStringArrayItem(value reflect.Value, idx int) (interface{}, bool) {
	return value.Index(idx).String(), value.Index(idx).String() == ""
}

var supportArrayFieldType = map[reflect.Kind]convArrayItem{
	reflect.Int:    convIntArrayItem,
	reflect.Int8:   convIntArrayItem,
	reflect.Int16:  convIntArrayItem,
	reflect.Int32:  convIntArrayItem,
	reflect.Int64:  convIntArrayItem,
	reflect.Uint:   convUIntArrayItem,
	reflect.Uint8:  convUIntArrayItem,
	reflect.Uint16: convUIntArrayItem,
	reflect.Uint32: convUIntArrayItem,
	reflect.Uint64: convUIntArrayItem,
	reflect.String: convStringArrayItem,
}

var supportFieldType = map[reflect.Kind]fieldInjectWhereCondition{
	reflect.Int:        injectInt,
	reflect.Int8:       injectInt,
	reflect.Int16:      injectInt,
	reflect.Int32:      injectInt,
	reflect.Int64:      injectInt,
	reflect.Uint:       injectUint,
	reflect.Uint8:      injectUint,
	reflect.Uint16:     injectUint,
	reflect.Uint32:     injectUint,
	reflect.Uint64:     injectUint,
	reflect.Complex64:  injectComplex,
	reflect.Complex128: injectComplex,
	reflect.Array:      injectSlice,
	reflect.Slice:      injectSlice,
	reflect.String:     injectString,
}

type PageQuery interface {
	GetPageCondition() (QueryPageCondition, []string)
}

func (q QueryPageCondition) GetPageCondition() (QueryPageCondition, []string) {
	return q, q.OrderBy
}

func (q QueryPageCondition) SetOrderBy(param []string) {
	q.OrderBy = param
}

func (q QueryPageCondition) SetPageCondition(pageNo, pageCount int) QueryPageCondition {
	condition := GetPageCondition(pageNo, pageCount)
	q.Limit = condition.Limit
	q.Offset = condition.Offset
	return q
}

type QueryPageCondition struct {
	OrderBy []string
	Limit   *int
	Offset  *int
}

type SelectValuesInterface interface {
	GetSelectValue() string
}

type SelectValues struct {
	values []string
}

func NewSelectValues(values []string) *SelectValues {
	return &SelectValues{values: values}
}

func (s *SelectValues) GetSelectValue() string {
	if len(s.values) == 0 {
		return ""
	}
	return strings.Join(s.values, ",")
}

func GetPageCondition(pageNo, pageCount int) QueryPageCondition {
	condition := QueryPageCondition{
		Limit: &pageCount,
	}
	if pageNo > 0 {
		i := (pageNo - 1) * pageCount
		condition.Offset = &i
	} else {
		i := 0
		condition.Offset = &i
	}
	return condition

}

func GetPageConditionWithDefaultValue(pageNo, pageCount int) QueryPageCondition {
	if pageNo == 0 {
		pageNo = 1
	}
	if pageCount == 0 {
		pageCount = 10
	}
	condition := QueryPageCondition{
		Limit: &pageCount,
	}
	if pageNo > 0 {
		i := (pageNo - 1) * pageCount
		condition.Offset = &i
	} else {
		i := 0
		condition.Offset = &i
	}
	return condition

}

type QueryConditionTest struct {
	order string `queryField:"at_id"`
	page  QueryPageCondition
}

func (q *QueryConditionTest) GetPageCondition() (QueryPageCondition, []string) {
	return q.page, []string{"id desc", "ctime asc"}
}
func ClientWhereIgnoreOrderBy(client Client, condition interface{}) Client {
	queryCondition, _ := fetchQueryAndPageInfo(condition)

	//set where
	client = extendDetailWhereCondition(client, queryCondition)

	//set order by
	//client = extendSQLCommonOrderBy(client, orderBy, defaultPlanDetailsOrderBy)

	return client
}

func ClientWhere(client Client, condition interface{}) Client {
	queryCondition, orderBy := fetchQueryAndPageInfo(condition)

	//set where
	client = extendDetailWhereCondition(client, queryCondition)

	//set order by
	client = extendSQLCommonOrderBy(client, orderBy, defaultPlanDetailsOrderBy)

	return client
}

func clientByCount(client Client, count *int64) error {
	return client.Offset(0).Count(count).Error()
}

func ClientByPage(client Client, condition interface{}, tableName string, result interface{}) (int64, error) {
	if tableName == "" {
		return 0, errors.New("table name can't be empty")
	}
	client = client.Table(tableName)
	client = ClientWhere(client, condition)
	count := int64(0)
	err := clientByCount(client, &count)
	if err != nil {
		return count, err
	}
	if count == 0 {
		return count, nil
	}
	// set limit offset
	client = extendSQLCommonPageCondition(client, condition)
	client = extendSQLCommonSelectValues(client, condition)
	err = client.Find(result).Error()

	return count, err
}

func ClientByCondition(client Client, condition interface{}) Client {
	// set where and order by
	client = ClientWhere(client, condition)

	return client
}

func ClientWhereAndValues(client Client, condition interface{}) Client {
	// set where and order by
	client = ClientWhere(client, condition)
	// select values
	client = extendSQLCommonSelectValues(client, condition)

	return client
}

type conditionColumnInfo struct {
	v                    reflect.Value
	operate              string
	queryFieldColumnName string
}

func newWhereColumnInfo(
	rawQueryField string,
	v reflect.Value,
	field reflect.StructField,
) conditionColumnInfo {
	queryFields := strings.Split(rawQueryField, ",")
	columnInfo := conditionColumnInfo{
		v:                    v.FieldByName(field.Name),
		queryFieldColumnName: queryFields[0],
	}
	if len(queryFields) == 2 {
		columnInfo.operate = queryFields[1]
	}
	return columnInfo
}

func (c conditionColumnInfo) getOperate() string {
	if c.v.Kind() == reflect.Array || c.v.Kind() == reflect.Slice {
		if c.operate == operateNe {
			return "NOT IN"
		}
		return "IN"
	}
	if v, ok := operateValueMap[c.operate]; ok {
		return v
	}
	// 默认返回相等查询
	return "="
}

func fetchQueryAndPageInfo(condition interface{}) ([]conditionColumnInfo, []string) {
	qc2v := make([]conditionColumnInfo, 0)
	ct := reflect.TypeOf(condition)
	value := reflect.ValueOf(condition)

	for i := 0; i < ct.NumField(); i++ {
		field := ct.Field(i)
		queryField := field.Tag.Get("column")
		//判断是否有匿名字段，如果有则需要取匿名字段中的字段去转换
		if field.Type.Kind() == reflect.Struct && field.Type.NumField() > 0 {
			temporarySlice, _ := fetchQueryAndPageInfo(value.FieldByName(field.Name).Interface())
			qc2v = append(qc2v, temporarySlice...)
		}
		if queryField != "" {
			columnInfo := newWhereColumnInfo(queryField, value, field)
			qc2v = append(qc2v, columnInfo)
		}
	}

	info, ok := condition.(PageQuery)
	if ok {
		_, orderBy := info.GetPageCondition()
		return qc2v, orderBy

	}

	return qc2v, nil
}

func extendDetailWhereCondition(common Client, condition []conditionColumnInfo) Client {
	for _, columnInfo := range condition {
		kind := columnInfo.v.Kind()
		if f, ok := supportFieldType[kind]; ok {
			common = f(common, columnInfo.queryFieldColumnName, columnInfo)
		}
	}
	return common
}

func extendSQLCommonOrderBy(client Client, orderBy []string, defaultOrderBy []string) Client {
	orderByStr := ""
	if len(orderBy) == 0 {
		orderByStr = strings.Join(defaultOrderBy, ",")
	} else {
		orderByStr = strings.Join(orderBy, ",")
	}
	return client.Order(orderByStr)
}

func extendSQLCommonPageCondition(client Client, condition interface{}) Client {
	info, ok := condition.(PageQuery)
	if ok {
		pageCondition, _ := info.GetPageCondition()

		if pageCondition.Offset != nil {
			client = client.Offset(*pageCondition.Offset)
		}
		if pageCondition.Limit != nil {
			client = client.Limit(*pageCondition.Limit)
		}
		return client
	}
	return client
}

func extendSQLCommonSelectValues(client Client, condition interface{}) Client {
	info, ok := condition.(SelectValuesInterface)
	if ok {
		selectValue := info.GetSelectValue()
		if len(selectValue) > 0 {
			client = client.Select(selectValue)
		}
	}
	return client
}

func ClientWithOtherWhereCondition(client Client, condition interface{}, otherCondition map[string]interface{}) Client {
	queryCondition, orderBy := fetchQueryAndPageInfo(condition)

	//set where
	client = extendDetailWhereCondition(client, queryCondition)

	// other condition
	for k, v := range otherCondition {
		client = client.Where(k, v)
	}

	//set order by
	client = extendSQLCommonOrderBy(client, orderBy, defaultPlanDetailsOrderBy)

	return client
}

func ClientByPageAndExtendWhereCondition(client Client, condition interface{}, otherCondition map[string]interface{}, tableName string, result interface{}) (int64, error) {
	if tableName == "" {
		return 0, errors.New("table name can't be empty")
	}
	client = client.Table(tableName)
	client = ClientWithOtherWhereCondition(client, condition, otherCondition)
	count := int64(0)
	err := clientByCount(client, &count)
	if err != nil {
		return count, err
	}
	if count == 0 {
		return count, nil
	}
	// set limit offset
	client = extendSQLCommonPageCondition(client, condition)
	client = extendSQLCommonSelectValues(client, condition)
	err = client.Find(result).Error()

	return count, err
}

type PageIn struct {
	Pageno     int    // 页码
	Count      int    // 数量
	OrderBy    string // 为空字符串 代表不需要排序
	IsGetTotal bool   // 为False代表不需要获取总数
}

func Paginator(qs Client, pageIn *PageIn, out interface{}) (int64, error) {
	// 如果pageIn 为nil，则不需要分页，当pageIn 不为nil时，pageno和count都大于0时才会用到offset和limit
	if pageIn == nil {
		err := qs.Find(out).Error()
		if err != nil {
			return 0, errors.WithStack(err)
		}
		return 0, nil
	}
	total := int64(0)
	if pageIn.IsGetTotal {
		err := qs.Count(&total).Error()
		if err != nil {
			return 0, errors.WithStack(err)
		}
	}
	if pageIn.OrderBy != "" {
		qs = qs.Order(pageIn.OrderBy)
	}
	if pageIn.Pageno > 0 && pageIn.Count > 0 {
		qs = qs.Offset(int((pageIn.Pageno - 1) * pageIn.Count)).Limit(pageIn.Count)
	}
	err := qs.Find(out).Error()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return total, nil
}

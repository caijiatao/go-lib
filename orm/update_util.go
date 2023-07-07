package orm

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

var (
	updateColumnOperate = map[string]string{
		"add": "+",
		"sub": "-",
	}
)

func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

type updateColumnInfo struct {
	v          reflect.Value
	operate    string
	columnName string
}

func newUpdateColumnInfo(
	rawQueryField string,
	v reflect.Value,
	field reflect.StructField,
) updateColumnInfo {
	updateFields := strings.Split(rawQueryField, ",")
	updateColumnInfo := updateColumnInfo{
		v:          v.FieldByName(field.Name),
		columnName: updateFields[0],
	}
	if len(updateFields) == 2 {
		updateColumnInfo.operate = updateFields[1]
	}
	return updateColumnInfo
}

func (u *updateColumnInfo) GetUpdateField() interface{} {
	o, ok := updateColumnOperate[u.operate]
	if !ok {
		return u.v.Interface()
	}
	return gorm.Expr(fmt.Sprintf("%s %s ?", u.columnName, o), u.v.Interface())
}

func GetUpdateFieldsMap(updateFields interface{}) map[string]interface{} {
	ut := reflect.TypeOf(updateFields)
	value := reflect.ValueOf(updateFields)
	updateColumnInfos := make([]updateColumnInfo, 0)

	for i := 0; i < ut.NumField(); i++ {
		field := ut.Field(i)
		updateField := field.Tag.Get("column")
		if updateField != "" {
			columnInfo := newUpdateColumnInfo(updateField, value, field)
			updateColumnInfos = append(updateColumnInfos, columnInfo)
		}
	}
	updateFieldsMap := make(map[string]interface{})
	for _, columnInfo := range updateColumnInfos {
		if isBlank(columnInfo.v) {
			continue
		}
		updateFieldsMap[columnInfo.columnName] = columnInfo.GetUpdateField()
	}
	return updateFieldsMap
}

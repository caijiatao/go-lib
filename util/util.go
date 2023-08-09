package util

import (
	"encoding/json"
	"time"
)

func CheckValueType(value interface{}, columnType string) (ok bool) {
	switch columnType {
	case "bool", "boolean":
		_, ok = value.(bool)
	case "int", "integer", "smallserial", "enum":
		_, ok = value.(int)
	case "int16":
		_, ok = value.(int16)
	case "smallint":
		_, ok = value.(int8)
	case "bigint", "real", "bigserial":
		_, ok = value.(int64)
	case "decimal", "numeric", "double precision", "money":
		_, ok = value.(float64)
	case "string", "character varying", "character(n) char(n)", "text", "UUID":
		_, ok = value.(string)
	case "json":
		var v string
		if v, ok = value.(string); !ok {
			return
		}
		ok = json.Valid([]byte(v))
	case "timestamp":
		_, ok = value.(time.Time)
	case "date", "time":
		_, ok = value.(string)
	}
	return ok
}

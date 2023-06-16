package service

import (
	"encoding/json"
	lua "github.com/yuin/gopher-lua"
)

func CheckValueType(value interface{}, columnType string) (ok bool) {
	switch columnType {
	case "bool":
		_, ok = value.(bool)
	case "int":
		_, ok = value.(int)
	case "int16":
		_, ok = value.(int16)
	case "string":
		_, ok = value.(string)
	case "json":
		var v string
		if v, ok = value.(string); !ok {
			return
		}
		ok = json.Valid([]byte(v))
	}
	return ok
}

func GetLuaType(value interface{}, columnType string) lua.LValue {
	switch columnType {
	case "bool":
		return lua.LBool(value.(bool))
	case "int":
		return lua.LNumber(value.(int))
	case "string":
		return lua.LString(value.(string))
	case "json":
		return lua.LString(value.(string))
	}
	return nil
}

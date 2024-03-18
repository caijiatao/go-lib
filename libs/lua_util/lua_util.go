package lua_util

import (
	lua "github.com/yuin/gopher-lua"
)

func GetLuaType(value interface{}, columnType string) lua.LValue {
	switch columnType {
	case "bool", "boolean":
		return lua.LBool(value.(bool))
	case "int", "integer", "smallserial", "enum", "int16", "smallint", "int8", "bigint", "real", "bigserial":
		return lua.LNumber(value.(int64))
	case "string", "character varying", "character(n) char(n)", "text", "UUID":
		return lua.LString(value.(string))
	case "json":
		return lua.LString(value.(string))
	case "decimal", "numeric", "double precision", "money":
		return lua.LNumber(value.(float64))
	case "timestamp", "date", "time":
		return lua.LString(value.(string))
	}
	return nil
}

func GetLuaTypeValue(value interface{}) lua.LValue {
	if v, ok := value.(bool); ok {
		return lua.LBool(v)
	}
	if v, ok := value.(int64); ok {
		return lua.LNumber(v)
	}
	if v, ok := value.(int); ok {
		return lua.LNumber(v)
	}
	if v, ok := value.(float32); ok {
		return lua.LNumber(v)
	}
	if v, ok := value.(int32); ok {
		return lua.LNumber(v)
	}
	if v, ok := value.(int16); ok {
		return lua.LNumber(v)
	}
	if v, ok := value.(int8); ok {
		return lua.LNumber(v)
	}
	if v, ok := value.(uint64); ok {
		return lua.LNumber(v)
	}
	if v, ok := value.(uint32); ok {
		return lua.LNumber(v)
	}
	if v, ok := value.(uint16); ok {
		return lua.LNumber(v)
	}
	if v, ok := value.(uint8); ok {
		return lua.LNumber(v)
	}
	if v, ok := value.(uint); ok {
		return lua.LNumber(v)
	}
	if v, ok := value.(string); ok {
		return lua.LString(v)
	}
	if v, ok := value.(string); ok {
		return lua.LString(v)
	}
	if v, ok := value.(float64); ok {
		return lua.LNumber(v)
	}
	if v, ok := value.(string); ok {
		return lua.LString(v)
	}
	return nil
}

package service

import (
	"errors"
	lua "github.com/yuin/gopher-lua"
)

type StandardTableMapping struct {
	MappingTableName            string
	ColumnName2ColumnMappingMap map[string]ColumnStandardMapping
}

type ColumnStandardMapping struct {
	StandardColumnName string
	StandardColumnType string
	ConvertScript      string
}

func NewStandardMapping(standardColumnType string, convertScript string, standardColumnName string) ColumnStandardMapping {
	return ColumnStandardMapping{StandardColumnType: standardColumnType, ConvertScript: convertScript, StandardColumnName: standardColumnName}
}

func (s *ColumnStandardMapping) MapToStandValue(value interface{}) (result interface{}, err error) {
	if ok := CheckValueType(value, s.StandardColumnType); !ok {
		return nil, errors.New("type error")
	}

	if result, err = s.convertValue(value); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *ColumnStandardMapping) convertValue(value interface{}) (interface{}, error) {
	if len(s.ConvertScript) == 0 {
		return value, nil
	}
	L := lua.NewState()
	defer L.Close()
	err := L.DoString(s.ConvertScript)
	if err != nil {
		return nil, err
	}
	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("mapValue"),
		NRet:    1,
		Protect: true,
	}, GetLuaType(value, s.StandardColumnType)); err != nil {
		panic(err)
	}
	lValue := L.Get(-1)
	return lValue.String(), nil
}

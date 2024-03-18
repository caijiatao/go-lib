package lua_util

import (
	"fmt"
	"golib/libs/logger"
)

type LuaError struct {
	// 这里0表示不知道哪出错，-1表示最后一行(列)非空行(列)
	StartLine   int    `json:"start_line"`
	EndLine     int    `json:"end_line"`
	StartColumn int    `json:"start_column"`
	EndColumn   int    `json:"end_column"`
	Info        string `json:"info"`
}

func NewLuaError(info string) *LuaError {
	return unknownError(info)
}

func (e *LuaError) Error() string {
	if e.StartLine == 0 && e.EndLine == 0 {
		return e.Info
	}
	return fmt.Sprintf("start:%d:%d, end:%d:%d, info:%s", e.StartLine, e.StartColumn, e.EndLine, e.EndColumn, e.Info)
}

func (e *LuaError) parse(err error) error {
	for _, parser := range errorParsers {
		if parserErr := parser.parse(err); parserErr != nil {
			return parserErr
		}
	}
	return err
}

func unknownError(info string) *LuaError {
	return &LuaError{
		StartLine:   0,
		EndLine:     0,
		StartColumn: 0,
		EndColumn:   0,
		Info:        info,
	}
}

func parseErr(err error) error {
	logger.Error("origin error is:", err)
	luaError := NewLuaError(err.Error())
	err = luaError.parse(err)
	return err
}

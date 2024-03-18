package lua_util

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"regexp"
	"strconv"
)

var (
	// 对可能的错误类型做解析
	errorParsers = []*errorParser{
		// 1.类似这样的<string> line:2(column:3) near '+':   parse error
		newErrorParser(`line:(\d+)\(column:(\d+)\)\s*(.+)`, func(pattern *regexp.Regexp, err error) error {
			if try := pattern.FindStringSubmatch(err.Error()); len(try) == 4 {
				linePos, _ := strconv.Atoi(try[1])
				colPos, _ := strconv.Atoi(try[2])
				return &LuaError{
					StartLine:   linePos,
					EndLine:     linePos,
					StartColumn: colPos,
					EndColumn:   colPos,
					Info:        try[3],
				}
			}
			return nil
		}),

		/* 2.类似这样的
		origin error is: <string>:3: cannot perform add operation between number and string
		stack traceback:
			<string>:3: in function 'test'
			<string>:5: in main chunk
			[G]: ?
		*/
		newErrorParser(`:(\d+):\s*(.+)`, func(pattern *regexp.Regexp, err error) error {
			if try := pattern.FindStringSubmatch(err.Error()); len(try) == 3 {
				linePos, _ := strconv.Atoi(try[1])
				return &LuaError{
					StartLine:   linePos,
					EndLine:     linePos,
					StartColumn: 0,
					EndColumn:   0,
					Info:        try[2],
				}
			}
			return nil
		}),

		// 3.类似这样的origin error is: <string> at EOF:   parse error
		newErrorParser(`<\S+>(.+)`, func(pattern *regexp.Regexp, err error) error {
			if try := pattern.FindStringSubmatch(err.Error()); len(try) == 2 {
				info := try[1]
				luaCheckErrPatternEOF := regexp.MustCompile(`at EOF:`)
				if len(luaCheckErrPatternEOF.FindStringSubmatch(info)) > 0 {
					// 在最后
					return &LuaError{
						StartLine:   -1,
						EndLine:     -1,
						StartColumn: -1,
						EndColumn:   -1,
						Info:        info,
					}
				}
			}
			return nil
		}),
	}
)

type parserFunc func(pattern *regexp.Regexp, err error) error

type errorParser struct {
	pattern *regexp.Regexp
	f       parserFunc
}

func (e *errorParser) parse(err error) error {
	return e.f(e.pattern, err)
}

func newErrorParser(patternStr string, f parserFunc) *errorParser {
	return &errorParser{
		pattern: regexp.MustCompile(patternStr),
		f:       f,
	}
}

type CheckConfig struct {
	FuncNames []string
}

func NewCheckConfig() *CheckConfig {
	return &CheckConfig{FuncNames: make([]string, 0)}
}

type CheckOption func(config *CheckConfig)

func WithCheckConfigFuncName(funcName string) CheckOption {
	return func(config *CheckConfig) {
		config.FuncNames = append(config.FuncNames, funcName)
	}
}

func WithCheckConfigFuncNames(funcName []string) CheckOption {
	return func(config *CheckConfig) {
		config.FuncNames = append(config.FuncNames, funcName...)
	}
}

func Check(code string, checkOpts ...CheckOption) error {
	l := lua.NewState()
	err := l.DoString(code)
	if err != nil {
		// todo: 解析一下出错位置
		return parseErr(err)
	}

	checkConfig := NewCheckConfig()
	for _, opt := range checkOpts {
		opt(checkConfig)
	}
	err = checkFuncExists(l, checkConfig)
	if err != nil {
		return err
	}

	return nil
}

// checkFuncExists
//
//	@Description:检查一下是否有指定的函数
//	@param state
//	@return error
func checkFuncExists(l *lua.LState, checkConfig *CheckConfig) error {
	for _, funcName := range checkConfig.FuncNames {
		attrGenFunc := l.GetGlobal(funcName) // 还是这样指定名字最简单
		if attrGenFunc == lua.LNil || attrGenFunc.Type() != lua.LTFunction {
			err := NewLuaError(fmt.Sprintf("expect %v as the generation function", funcName))
			return parseErr(err)
		}
	}
	return nil
}

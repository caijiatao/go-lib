package util

import "runtime"

func GetCallerFunctionName() string {
	pc, _, _, _ := runtime.Caller(2) // 获取上层调用栈的信息
	funcName := runtime.FuncForPC(pc).Name()
	return funcName
}

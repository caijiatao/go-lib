package logger

import (
	"context"
)

var defaultLog Logger

func init() {
	// TODO 读取不同的配置来替换不同的log底层实现
	defaultLog = NewLogrusProxy()
}

func getDefaultLogger() Logger {
	return defaultLog
}

type Logger interface {
	// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
	Info(args ...interface{})
	// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
	Infof(format string, args ...interface{})
	// CtxInfof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
	CtxInfof(ctx context.Context, format string, args ...interface{})
	// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
	Error(args ...interface{})
	// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
	Errorf(format string, args ...interface{})
	// CtxErrorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
	CtxErrorf(ctx context.Context, format string, args ...interface{})
}

func Error(v ...any) {
	getDefaultLogger().Error(v...)
}

func Info(v ...any) {
	getDefaultLogger().Info(v...)
}

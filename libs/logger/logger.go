package logger

import (
	"context"
	"golib/libs/util"
	"time"
)

var defaultLogger = &logger{}

type logger struct {
	L Logger
}

func init() {
	InitLogger()
}

func InitLogger() {
	Init()
	defaultLogger.L = newZapLogger()
}

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	CtxDebugf(ctx context.Context, format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})
	CtxInfof(ctx context.Context, format string, v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	CtxWarnf(ctx context.Context, format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	CtxErrorf(ctx context.Context, format string, v ...interface{})
}

// Debug log
func Debug(v ...any) {
	defaultLogger.L.Debug(v...)
}

// Debugf log
func Debugf(format string, v ...any) {
	defaultLogger.L.Debugf(format, v...)
}

// Info log
func Info(v ...any) {
	defaultLogger.L.Info(v...)
}

// Infof log
func Infof(format string, v ...any) {
	defaultLogger.L.Infof(format, v...)
}

// Warn log
func Warn(v ...any) {
	defaultLogger.L.Warn(v...)
}

// Warnf log
func Warnf(format string, v ...any) {
	defaultLogger.L.Warnf(format, v...)
}

// Error log
func Error(v ...any) {
	defaultLogger.L.Error(v...)
}

// Errorf log
func Errorf(format string, v ...any) {
	defaultLogger.L.Errorf(format, v...)
}

func CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.L.CtxDebugf(ctx, format, v...)
}

func CtxInfof(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.L.CtxInfof(ctx, format, v...)
}

func CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.L.CtxWarnf(ctx, format, v...)
}

func CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	defaultLogger.L.CtxErrorf(ctx, format, v...)
}

// LogDuration
//
//	@Description:
//	  使用的时候直接defer LogDuration(ctx,format,xxx)()
func LogDuration(ctx context.Context, format string, v ...interface{}) func() {
	start := time.Now()
	functionName := util.GetCallerFunctionName()
	return func() {
		CtxInfof(ctx, "Time taken by %s function is %v", functionName, time.Since(start))
		CtxInfof(ctx, format, v...)
	}
}

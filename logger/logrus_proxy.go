package logger

import (
	"context"
	"github.com/sirupsen/logrus"
)

type LogrusProxy struct {
	l *logrus.Logger
}

func (l *LogrusProxy) Debug(v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *LogrusProxy) Debugf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *LogrusProxy) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *LogrusProxy) Warn(v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *LogrusProxy) Warnf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *LogrusProxy) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *LogrusProxy) Info(args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *LogrusProxy) Infof(format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *LogrusProxy) CtxInfof(ctx context.Context, format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *LogrusProxy) Error(args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *LogrusProxy) Errorf(format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *LogrusProxy) CtxErrorf(ctx context.Context, format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func NewLogrusProxy() Logger {
	return &LogrusProxy{l: logrus.New()}
}

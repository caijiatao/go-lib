package logger

import (
	"context"
	"github.com/sirupsen/logrus"
)

type LogrusProxy struct {
	l *logrus.Logger
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

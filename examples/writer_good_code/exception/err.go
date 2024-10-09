package main

import (
	"github.com/pkg/errors"
)

func foo() error {
	return errors.New("something went wrong")
}

func bar() error {
	return errors.WithStack(foo()) // 将堆栈信息附加到错误上
}

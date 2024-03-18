package concurrency

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestConcurrencyGroup(t *testing.T) {
	cg := NewConcurrencyGroup()
	ctx := context.Background()
	cg.Go(quickOperate, ctx)
	cg.Go(errOperate, ctx)
	time.Sleep(time.Second)
	cg.Go(slowOperate, ctx)

	errs := cg.Wait()
	for _, err := range errs {
		fmt.Println(err)
	}
}

func TestConcurrencyGroupWithTimeout(t *testing.T) {
	cg := NewConcurrencyGroup(WithConcurrencyGroupParamTimeoutOpt(3 * time.Second))
	ctx := context.Background()
	cg.Go(quickOperate, ctx)
	cg.Go(slowOperate, ctx)

	errs := cg.Wait()
	for _, err := range errs {
		fmt.Println(err)
	}
}

func TestConcurrencyGroupWithNotInterrupt(t *testing.T) {
	cg := NewConcurrencyGroup(WithConcurrencyGroupParamInterruptTypeOpt(UninterruptedType))
	ctx := context.Background()
	// 错误在最前面，后面即便检查到也会继续操作
	cg.Go(errOperate, ctx)
	time.Sleep(time.Second)
	cg.Go(quickOperate, ctx)
	cg.Go(slowOperate, ctx)
	errs := cg.Wait()
	for _, err := range errs {
		fmt.Println(err)
	}
}

func TestConcurrencyGroupWithFastFail(t *testing.T) {
	cg := NewConcurrencyGroup(WithConcurrencyGroupParamInterruptTypeOpt(FastFail))
	ctx := context.Background()
	cg.Go(quickOperate, ctx)
	cg.Go(slowOperate, ctx)
	cg.Go(slowOperate, ctx)
	// 错误再最后，但是前面只要检查到立即退出
	cg.Go(errOperate, ctx)
	errs := cg.Wait()
	for _, err := range errs {
		fmt.Println(err)
	}
}

func quickOperate(ctx context.Context) error {
	// 模拟业务耗时操作
	time.Sleep(time.Second)
	return nil
}

func slowOperate(ctx context.Context) error {
	time.Sleep(5 * time.Second)
	return nil
}

func errOperate(ctx context.Context) error {
	return errors.New("test err")
}

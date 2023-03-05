package concurrency

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"golib/goasync"
)

// InterruptType 失败类型
type InterruptType int64

const (
	FastFail          InterruptType = iota + 1 // 只要有一个错误立马终止执行
	PreExecFail                                // 按加入方法顺序来判断错误，e.g. 第二个方法执行前，如果第一个方法已经出现错误，则不执行直接返回，如果第三个执行失败，则仍然会执行第二个方法
	UninterruptedType                          // 不中断执行
)

var (
	ConcurrencyGroupInterruptErr   = errors.New("interrupted!")
	ConcurrencyGroupContextDoneErr = errors.New("Context Done err!")
)

func IsInterruptErr(err error) bool {
	return errors.Is(err, ConcurrencyGroupInterruptErr)
}

func IsContextDoneErr(err error) bool {
	return errors.Is(err, ConcurrencyGroupContextDoneErr)
}

type concurrencyGroupParam struct {
	timeout       time.Duration // 执行超时
	interruptType InterruptType // 执行的中断类型
	concurrency   int           // 控制并发数
}

type GroupParamOptFunc func(param *concurrencyGroupParam)

func WithConcurrencyGroupParamTimeoutOpt(timeout time.Duration) GroupParamOptFunc {
	return func(param *concurrencyGroupParam) {
		param.timeout = timeout
	}
}

func WithConcurrencyGroupParamInterruptTypeOpt(interruptType InterruptType) GroupParamOptFunc {
	return func(param *concurrencyGroupParam) {
		param.interruptType = interruptType
	}
}

func WithConcurrencyGroupParamConcurrency(concurrency int) GroupParamOptFunc {
	return func(param *concurrencyGroupParam) {
		param.concurrency = concurrency
	}
}

type concurrencyGroupErrs struct {
	sync.Mutex // 扩容时保证并发安全
	errs       []error
}

func newConcurrencyGroupErrs() *concurrencyGroupErrs {
	return &concurrencyGroupErrs{
		errs: make([]error, 8), // 默认保存8个错误
	}
}

func (cge *concurrencyGroupErrs) growsErrSlice(index int) {
	// 不需要扩容直接返回
	if len(cge.errs) > index {
		return
	}
	cge.Lock()
	defer cge.Unlock()
	// 二次检查，避免重复扩容
	if len(cge.errs) > index {
		return
	}

	// 真正开始执行扩容
	growLen := 0
	if len(cge.errs) < 1024 {
		// 小于1024按原长度两倍进行扩容，避免需要频繁进行扩容
		growLen = len(cge.errs)
	} else {
		// 大于1024 则按原长度的0.25进行扩容
		growLen = int(math.Ceil(float64(len(cge.errs)) * 0.25))
	}
	groupErrs := make([]error, growLen)
	cge.errs = append(cge.errs, groupErrs...)
}

func (cge *concurrencyGroupErrs) setErr(err error, index int) {
	cge.growsErrSlice(index)
	cge.errs[index] = err
}

// 判断是否出现过错误
func (cge *concurrencyGroupErrs) hasError() bool {
	errs := cge.errs
	for _, err := range errs {
		if err != nil {
			return true
		}
	}
	return false
}

// 前面的函数是否出现错误
func (cge *concurrencyGroupErrs) beforeHasError(index int) bool {
	errs := cge.errs
	// 执行到后面时可能前面有些错误还没有加入，所以取错误的长度进行遍历
	if index > len(errs) {
		index = len(errs)
	}
	for i := 0; i < index-1; i++ {
		if errs[i] != nil {
			return true
		}
	}
	return false
}

// GetAllRealErrs
// @Description: 获取执行结果真正的error
func (cge *concurrencyGroupErrs) GetAllRealErrs() []error {
	errs := make([]error, 0, len(cge.errs))
	for _, err := range cge.errs {
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// GetFirstErr
// @Description: 获取第一个错误
func (cge *concurrencyGroupErrs) GetFirstErr() error {
	for _, err := range cge.errs {
		if err != nil {
			return err
		}
	}
	return nil
}

type groupGoFunc func(ctx context.Context) error

type groupFuncInfo struct {
	index  int         // 执行函数的顺序信息
	f      groupGoFunc // 真正需要执行的函数
	cancel func()
}

// ConcurrencyGroup
// @Description: 并发执行的group结构体
type ConcurrencyGroup struct {
	cancel func()

	wg            sync.WaitGroup
	concurrencyCh chan struct{}
	groupFuncs    []*groupFuncInfo

	*concurrencyGroupParam

	*concurrencyGroupErrs
}

func NewConcurrencyGroup(opts ...GroupParamOptFunc) *ConcurrencyGroup {
	cg := &ConcurrencyGroup{
		cancel:               nil,
		wg:                   sync.WaitGroup{},
		groupFuncs:           make([]*groupFuncInfo, 0),
		concurrencyGroupErrs: newConcurrencyGroupErrs(),
		concurrencyGroupParam: &concurrencyGroupParam{
			concurrency: 8, //默认并发数为8 ，这里是指真正执行操作的并发数
		},
	}
	for _, opt := range opts {
		opt(cg.concurrencyGroupParam)
	}
	cg.concurrencyCh = make(chan struct{}, cg.concurrency)
	return cg
}

// addGoFunc
// @Description: 增加执行的函数
func (cg *ConcurrencyGroup) addGoFunc(f groupGoFunc, ctx context.Context) *groupFuncInfo {
	cg.wg.Add(1)
	goFunc := &groupFuncInfo{
		index: len(cg.groupFuncs), // 下标进行顺序索引
		f:     f,
	}
	cg.groupFuncs = append(cg.groupFuncs, goFunc)
	return goFunc
}

func (cg *ConcurrencyGroup) Wait() []error {
	cg.wg.Wait()
	// 执行退出操作
	_ = cg.Close()
	// 裁切掉过多扩容的部分
	return cg.errs[0:len(cg.groupFuncs)]
}

// isInterruptExec
// @Description: 是否中断执行
func (cg *ConcurrencyGroup) isInterruptExec(funcIndex int) bool {
	switch cg.interruptType {
	case PreExecFail:
		return cg.beforeHasError(funcIndex)
	case FastFail:
		// 有错误立即中断
		return cg.hasError()
	case UninterruptedType:
		return false
	default:
		// 默认策略是前置有执行失败则中断
		return cg.beforeHasError(funcIndex)
	}
}

func (cg *ConcurrencyGroup) doFunc(f *groupFuncInfo, ctx context.Context) chan error {
	errChan := make(chan error, 1)
	// 不执行直接返回
	if cg.isInterruptExec(f.index) {
		errChan <- ConcurrencyGroupInterruptErr
		return errChan
	}
	go func() {
		defer func() {
			r := recover()
			err := goasync.PanicErrHandler(r)
			if err != nil {
				errChan <- errors.New(fmt.Sprintf("%q", err))
			}
		}()
		// 二次检查，被调度时产生错误了
		if cg.isInterruptExec(f.index) {
			errChan <- ConcurrencyGroupInterruptErr
			return
		}
		// 读取错误
		errChan <- f.f(ctx)
	}()
	return errChan
}

func (cg *ConcurrencyGroup) Close() error {
	for _, goFunc := range cg.groupFuncs {
		if goFunc.cancel != nil {
			goFunc.cancel()
		}
	}
	if cg.cancel != nil {
		cg.cancel()
	}
	return nil
}

func (cg *ConcurrencyGroup) withContextTimeout(ctx context.Context) (context.Context, func()) {
	if cg.timeout <= 0 {
		return ctx, nil
	}
	return context.WithTimeout(ctx, cg.timeout)
}

func (cg *ConcurrencyGroup) Go(f groupGoFunc, ctx context.Context) {
	goFunc := cg.addGoFunc(f, ctx)
	go func() {
		// 最后执行done,避免recover中还有收尾工作没做
		defer func() {
			// 执行完成读走缓冲区数据
			<-cg.concurrencyCh
			cg.wg.Done()
		}()

		defer func() {
			r := recover()
			_ = goasync.PanicErrHandler(r)
		}()
		// 如果已经达到最大并发数，则进入等待
		cg.concurrencyCh <- struct{}{}
		// 开始执行则设置超时
		ctx, goFunc.cancel = cg.withContextTimeout(ctx)
		errChan := cg.doFunc(goFunc, ctx)
		select {
		case err := <-errChan:
			cg.setErr(err, goFunc.index)
		case <-ctx.Done(): // context被done则直接返回
			cg.setErr(ConcurrencyGroupContextDoneErr, goFunc.index)
		}
	}()
}

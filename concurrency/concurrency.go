package concurrency

import (
	"context"
	"sync"

	"golib/goasync"
)

const (
	gConcurrencyCount = 32
)

type Result struct {
	Res interface{}
	Err error
}

type ConcurrencyTasks struct {
	Ctx         context.Context
	IDs         []interface{}
	Results     []Result
	Func        func(ID interface{}) (interface{}, error)
	Concurrency int
	ch          chan struct{}
	wg          sync.WaitGroup
}

func (t *ConcurrencyTasks) Run() {
	rs := make([]Result, len(t.IDs))
	t.Results = rs
	if len(t.IDs) == 0 {
		return
	}
	t.wg.Add(len(t.IDs))
	count := gConcurrencyCount
	if t.Concurrency > 0 {
		count = t.Concurrency
	}
	t.ch = make(chan struct{}, count)
	for i := 0; i < len(t.IDs); i++ {
		t.ch <- struct{}{}
		go t.doFun(i)
	}
	close(t.ch)
	t.wg.Wait()
}

func (t *ConcurrencyTasks) doFun(i int) {
	var ID interface{}
	defer func() {
		// 从panic中恢复
		var r any = recover()
		goasync.PanicErrHandler(r)
	}()
	defer t.wg.Done()
	//执行
	ID = t.IDs[i]
	r, err := t.Func(ID)
	t.Results[i] = Result{
		Res: r,
		Err: err,
	}
	<-t.ch
}

func (t *ConcurrencyTasks) Successes() int {
	successes := 0
	for _, r := range t.Results {
		if r.Err != nil {
			continue
		}
		successes += 1
	}
	return successes
}

func (t *ConcurrencyTasks) GetErr() error {
	for _, r := range t.Results {
		if r.Err != nil {
			return r.Err
		}
	}
	return nil
}

package gopool

import (
	"golib/libs/goasync"
	"sync"
)

var workerPool sync.Pool

func initWorker() {
	workerPool.New = newWorker
}

func newWorker() interface{} {
	return &worker{}
}

type worker struct {
	pool *pool
	sync.Mutex
}

func (w *worker) run() {
	go func() {
		for {
			t := w.pool.tasks.getTask()
			if t == nil {
				w.Recycle()
				return
			}
			func() {
				defer func() {
					r := recover()
					_ = goasync.PanicErrHandler(r)
				}()
				t.fc()
			}()
			t.Recycle()
		}
	}()
}

func (w *worker) close() {
	w.Lock()
	defer w.Unlock()
	w.pool.decWorkerCount()
}

func (w *worker) Recycle() {
	w.close()
	w.pool = nil
	workerPool.Put(w)
}

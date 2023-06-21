package gopool

import (
	"golang.org/x/net/context"
	"sync"
	"sync/atomic"
)

var taskPool sync.Pool

func initTask() {
	taskPool.New = newTask
}

type taskManager struct {
	taskHead *task
	taskTail *task

	taskCount int64
	sync.Mutex
}

func (tm *taskManager) getTask() *task {
	tm.Lock()
	defer tm.Unlock()
	if tm.taskHead == nil {
		return nil
	}
	t := tm.taskHead
	tm.taskHead = tm.taskHead.next
	atomic.AddInt64(&tm.taskCount, -1)
	return t
}

func (tm *taskManager) addTask(t *task) {
	tm.Lock()
	defer tm.Unlock()
	if tm.taskHead == nil {
		tm.taskHead = t
		tm.taskTail = t
	} else {
		tm.taskTail.next = t
		tm.taskTail = t
	}
	atomic.AddInt64(&tm.taskCount, 1)
}

func (tm *taskManager) getTaskCount() int64 {
	return atomic.LoadInt64(&tm.taskCount)
}

type task struct {
	ctx context.Context
	fc  func()

	next *task
}

func newTask() interface{} {
	return &task{}
}

func (t *task) zero() {
	t.ctx = nil
	t.fc = nil
	t.next = nil
}

func (t *task) Recycle() {
	t.zero()
	taskPool.Put(t)
}

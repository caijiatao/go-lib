package gopool

import (
	"context"
	"sync/atomic"
)

type Pool interface {
	Name() string
	SetCap(cap int64)
	Go(fc func())
	GoWithContext(ctx context.Context, fc func())
	WorkerCount() int64
}
type pool struct {
	name string

	config *Config

	tasks *taskManager

	workerCount int64
}

func NewPool(name string, config *Config) Pool {
	return &pool{name: name, config: config, tasks: &taskManager{}}
}

func (p *pool) Name() string {
	return p.name
}

func (p *pool) SetCap(cap int64) {
	p.config.setCap(cap)
}

func (p *pool) Go(fc func()) {
	p.GoWithContext(context.Background(), fc)
}

func (p *pool) GoWithContext(ctx context.Context, fc func()) {
	t := taskPool.Get().(*task)
	t.ctx = ctx
	t.fc = fc
	p.tasks.addTask(t)
	if p.isNeedScaleWorker() {
		p.scaleWorker()
	}
}

func (p *pool) isNeedScaleWorker() bool {
	return p.tasks.getTaskCount() >= p.config.ScaleThreshold && p.WorkerCount() <= p.config.getCap()
}

func (p *pool) scaleWorker() {
	p.incrWorkerCount()
	w := workerPool.Get().(*worker)
	w.pool = p
	w.run()
}

func (p *pool) incrWorkerCount() {
	atomic.AddInt64(&p.workerCount, 1)
}

func (p *pool) decWorkerCount() {
	atomic.AddInt64(&p.workerCount, -1)
}

func (p *pool) WorkerCount() int64 {
	return atomic.LoadInt64(&p.workerCount)
}

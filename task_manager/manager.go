package task_manager

import (
	"context"
	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golib/concurrency"
	"golib/concurrency/gopool"
	"golib/etcd_helper"
	"golib/logger"
	"sync"
	"time"
)

var (
	HandlerNameRepeatErr = errors.New("handler name repeat")
)

var (
	taskPreemptFail = errors.New("task preempt fail")
)

var globalTaskManager *TaskManager
var defaultTaskManagerInitOnce sync.Once

func Init() {
	defaultTaskManagerInitOnce.Do(func() {
		globalTaskManager = NewTaskManager()
	})
}

type TaskManager struct {
	handlerMap sync.Map

	wg sync.WaitGroup
}

func NewTaskManager() *TaskManager {
	manager := &TaskManager{}
	ctx, cancelCtx := context.WithCancel(context.Background())
	ctx = etcd_helper.BindContext(ctx)

	manager.wg.Add(1)
	gopool.Go(func() {
		err := manager.watchTasks(ctx)
		if err != nil {
			logger.CtxErrorf(ctx, "watchTasks err: %+v", err)
			return
		}
	})
	gopool.Go(func() {
		err := manager.scanTasks(ctx)
		if err != nil {
			logger.CtxErrorf(ctx, "scanTasks err: %+v", err)
			return
		}
	})
	manager.wg.Wait()

	// 服务退出时清理任务
	cancel := func() {
		cancelCtx()
		err := manager.Close()
		if err != nil {
			logger.CtxErrorf(ctx, "manager.Close err: %+v", err)
		}
	}
	concurrency.GracefulQuit(cancel)

	return manager
}

func (m *TaskManager) RegisterHandler(handler *TaskHandler) error {
	_, ok := m.handlerMap.Load(handler.Name)
	if ok {
		return HandlerNameRepeatErr
	}
	m.handlerMap.Store(handler.Name, handler)
	return nil
}

func (m *TaskManager) getHandler(handlerName TaskHandlerName) *TaskHandler {
	value, ok := m.handlerMap.Load(handlerName)
	if !ok {
		return nil
	}
	handler, ok := value.(*TaskHandler)
	if !ok {
		return nil
	}
	return handler
}

func (m *TaskManager) handleTask(ctx context.Context, k, v []byte, kvVersion int64) {
	key := taskKey(k)
	handler := m.getHandler(key.getHandlerName())
	if handler == nil {
		logger.Errorf("handler not found, handlerName: %s", key.getHandlerName())
		return
	}

	task := decodeTask(string(v), kvVersion, handler.ParamsType)
	if task == nil {
		logger.CtxInfof(ctx, "decodeTask fail, value:%s", v)
		return
	}
	handler.handleTask(ctx, task)
}

func (m *TaskManager) watchTasks(ctx context.Context) (err error) {
	wch := etcd_helper.Watch(ctx, taskKeyPrefix, clientv3.WithPrefix())
	m.wg.Done()
	for wc := range wch {
		for _, event := range wc.Events {
			if event.Type == clientv3.EventTypePut {
				m.handleTask(ctx, event.Kv.Key, event.Kv.Value, event.Kv.Version)
			}
		}
	}
	return nil
}

func (m *TaskManager) scanTasks(ctx context.Context) error {
	for {
		afterTime := time.After(time.Minute)
		select {
		case <-afterTime:
			logger.CtxInfof(ctx, "scanTasks start")
		case <-ctx.Done():
			logger.CtxInfof(ctx, "scanTasks ctx.Done")
			return nil
		}
		getResponse, err := etcd_helper.Get(ctx, taskKeyPrefix, clientv3.WithPrefix())
		if err != nil {
			continue
		}
		logger.CtxInfof(ctx, "scanTasks len(getResponse.Kvs): %d", len(getResponse.Kvs))
		for _, kv := range getResponse.Kvs {
			m.handleTask(ctx, kv.Key, kv.Value, kv.Version)
		}
	}
}

// Close
//
//	@Description: 停掉所有正在跑的任务
func (m *TaskManager) Close() error {
	m.handlerMap.Range(func(key, value interface{}) bool {
		taskHandler, ok := value.(*TaskHandler)
		if !ok {
			return true
		}
		err := taskHandler.Close()
		if err != nil {
			logger.Errorf("close interruptRunningTask fail: %+v", err)
		}
		return true
	})
	return nil
}

package task_manager

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golib/libs/concurrency"
	"golib/libs/concurrency/gopool"
	"golib/libs/etcd_helper"
	"golib/libs/logger"
	"golib/libs/orm"
	"sync"
	"time"
)

var (
	HandlerNameRepeatErr = errors.New("handler name repeat")
)

var (
	taskPreemptFail = errors.New("task preempt fail")
)

var (
	scanTaskDelayTime = time.Minute
)

var globalTaskManager *TaskManager
var defaultTaskManagerInitOnce sync.Once

func Init(prefix string) {
	defaultTaskManagerInitOnce.Do(func() {
		taskKeyPrefix = prefix
		logger.Infof("task manager taskKeyPrefix: %s", taskKeyPrefix)
		globalTaskManager = NewTaskManager()
	})
}

type TaskManager struct {
	handlerMap sync.Map

	wg sync.WaitGroup

	cancelCtx func()
}

func NewTaskManager() *TaskManager {
	manager := &TaskManager{}
	ctx, cancelCtx := context.WithCancel(context.Background())
	ctx = orm.BindContext(ctx)
	ctx = etcd_helper.BindContext(ctx)
	manager.cancelCtx = cancelCtx

	manager.wg.Add(1)
	gopool.Go(func() {
		manager.watchTasks(ctx)
	})
	gopool.Go(func() {
		manager.scanTasks(ctx)
	})
	manager.wg.Wait()

	// 服务退出时清理任务
	cancel := func() {
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

func (m *TaskManager) handleTask(ctx context.Context, k, v []byte, kvVersion int64) error {
	key := taskKey(k)
	handler := m.getHandler(key.getHandlerName())
	if handler == nil {
		return errors.New(fmt.Sprintf("handler not found, handlerName: %s", key.getHandlerName()))
	}

	task := &Task{}
	err := task.decode(string(v))
	if err != nil {
		return err
	}
	task.TaskVersion = kvVersion
	handler.handleTask(ctx, task)
	return nil
}

func (m *TaskManager) watchTasks(ctx context.Context) {
	wch := etcd_helper.Watch(ctx, taskKeyPrefix, clientv3.WithPrefix())
	logger.Infof("watchTasks prefix: %s", taskKeyPrefix)
	m.wg.Done()
	for wc := range wch {
		for _, event := range wc.Events {
			if event.Type == clientv3.EventTypePut {
				err := m.handleTask(ctx, event.Kv.Key, event.Kv.Value, event.Kv.Version)
				if err != nil {
					logger.CtxErrorf(ctx, "watchTasks handleTask err: %+v", err)
				}
			}
		}
	}
	return
}

func (m *TaskManager) scanTasks(ctx context.Context) {
	for {
		afterTime := time.After(scanTaskDelayTime)
		select {
		case <-afterTime:
			logger.CtxInfof(ctx, "scanTasks start")
		case <-ctx.Done():
			logger.CtxInfof(ctx, "scanTasks ctx.Done")
			return
		}
		getResponse, err := etcd_helper.Get(ctx, taskKeyPrefix, clientv3.WithPrefix())
		if err != nil {
			continue
		}
		logger.CtxInfof(ctx, "scanTasks len(getResponse.Kvs): %d, task key prefix:%s", len(getResponse.Kvs), taskKeyPrefix)
		for _, kv := range getResponse.Kvs {
			err = m.handleTask(ctx, kv.Key, kv.Value, kv.Version)
			if err != nil {
				logger.CtxErrorf(ctx, "scanTasks handleTask err: %+v", err)
			}
		}
	}
}

// Close
//
//	@Description: 停掉所有正在跑的任务
func (m *TaskManager) Close() error {
	m.cancelCtx()
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

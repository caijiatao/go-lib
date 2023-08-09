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
)

var (
	HandlerNameRepeatErr = errors.New("handler name repeat")
)

var (
	taskPreemptFail = errors.New("task preempt fail")
)

var defaultTaskManager *TaskManager
var defaultTaskManagerInitOnce sync.Once

func Init() {
	defaultTaskManagerInitOnce.Do(func() {
		defaultTaskManager = NewTaskManager()
	})
}

type TaskManager struct {
	handlerMap   sync.Map
	runningTasks sync.Map

	wg sync.WaitGroup
}

func NewTaskManager() *TaskManager {
	manager := &TaskManager{}
	ctx, cancelCtx := context.WithCancel(context.Background())
	ctx = etcd_helper.BindContext(ctx)

	manager.wg.Add(1)
	gopool.Go(func() {
		err := manager.watchExecutableTasks(ctx)
		if err != nil {
			logger.CtxErrorf(ctx, "watchExecutableTasks err: %+v", err)
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

func (m *TaskManager) preemptTask(ctx context.Context, task *Task) (err error) {
	task.markRunning()

	txn := etcd_helper.Context(ctx).
		Txn(ctx).
		If(clientv3.Compare(clientv3.Version(task.getTaskKey().String()), "=", task.TaskVersion)).
		Then(clientv3.OpPut(task.getTaskKey().String(), task.encode()))
	txnResponse, err := txn.Commit()
	if err != nil {
		return err
	}
	if !txnResponse.Succeeded {
		return taskPreemptFail
	}
	m.runningTasks.Store(task.getTaskKey().String(), task)
	return nil
}

func (m *TaskManager) completeTask(ctx context.Context, task *Task) (err error) {
	putResponse, err := etcd_helper.Put(ctx, task.getTaskKey().String(), task.encode())
	if err != nil {
		logger.CtxErrorf(ctx, "completeTask Put err: %+v", err)
		return err
	}
	task.TaskVersion = putResponse.Header.GetRevision()
	m.runningTasks.Delete(task.getTaskKey().String())
	return nil
}

func (m *TaskManager) handlePutEvent(ctx context.Context, event *clientv3.Event) {
	key := taskKey(event.Kv.Key)
	handler := m.getHandler(key.getHandlerName())
	if handler == nil {
		logger.Errorf("handler not found, handlerName: %s", key.getHandlerName())
		return
	}

	task := DecodeTask(string(event.Kv.Value), event.Kv.Version, handler.ParamsType)
	if task.Status != pending {
		return
	}

	ctx = logger.WithTraceId(ctx)

	err := m.preemptTask(ctx, task)
	if err != nil {
		logger.CtxInfof(ctx, "preemptTask fail: %+v", err)
		return
	}
	// 抢占到了，执行完后需要释放
	defer func() {
		err := m.completeTask(ctx, task)
		if err != nil {
			logger.CtxErrorf(ctx, "completeTask fail: %+v", err)
		}
	}()
	err = handler.Exec(ctx, task)
	if err != nil {
		logger.CtxErrorf(ctx, "handler.Exec err: %+v", err)
		return
	}
}

func (m *TaskManager) watchExecutableTasks(ctx context.Context) (err error) {
	wch := etcd_helper.Watch(ctx, taskKeyPrefix, clientv3.WithPrefix())
	m.wg.Done()
	for wc := range wch {
		for _, event := range wc.Events {
			if event.Type == clientv3.EventTypePut {
				m.handlePutEvent(ctx, event)
			}
		}
	}
	return nil
}

func (m *TaskManager) interruptTask(ctx context.Context, task *Task) error {
	task.markPending()
	_, err := etcd_helper.Put(ctx, task.getTaskKey().String(), task.encode())
	if err != nil {
		return err
	}
	return nil
}

// Close
//
//	@Description: 停掉所有正在跑的任务
func (m *TaskManager) Close() error {
	ctx := etcd_helper.BindContext(logger.WithTraceId(context.Background()))
	m.runningTasks.Range(func(key, value interface{}) bool {
		task, ok := value.(*Task)
		if !ok {
			return true
		}
		err := m.interruptTask(ctx, task)
		if err != nil {
			logger.CtxErrorf(ctx, "close interruptTask fail: %+v", err)
		}
		return true
	})
	return nil
}

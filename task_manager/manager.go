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
	handlerMap   sync.Map
	runningTasks sync.Map

	taskExecGoPool gopool.Pool // 控制任务执行的并发度

	wg sync.WaitGroup
}

func NewTaskManager() *TaskManager {
	manager := &TaskManager{
		taskExecGoPool: gopool.NewPool(taskKeyPrefix, gopool.NewConfig(gopool.WithConfigCap(1000))),
	}
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

func (m *TaskManager) completeTask(ctx context.Context, task *Task) error {
	if task.Status == success {
		deleteResponse, err := etcd_helper.Delete(ctx, task.getTaskKey().String())
		if err != nil {
			logger.CtxErrorf(ctx, "completeTask Delete err: %+v", err)
			return err
		}
		logger.CtxInfof(ctx, "completeTask DeleteResponse.deleted: %d", deleteResponse.Deleted)
	} else {
		putResponse, err := etcd_helper.Put(ctx, task.getTaskKey().String(), task.encode())
		if err != nil {
			logger.CtxErrorf(ctx, "completeTask Put err: %+v", err)
			return err
		}
		logger.CtxInfof(ctx, "completeTask PutResponse.header.revision: %d", putResponse.Header.Revision)
	}
	m.runningTasks.Delete(task.getTaskKey().String())
	return nil
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
	ctx = logger.CtxWithTraceId(ctx, task.TraceId)
	m.taskExecGoPool.Go(func() {
		m.handlerPendingTask(ctx, task, handler)
	})
}

func (m *TaskManager) handlerPendingTask(ctx context.Context, task *Task, handler *TaskHandler) {
	if task.Status != pending {
		return
	}
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
	// 延迟执行
	after := task.NextExecuteTime.Sub(time.Now())
	select {
	case <-time.After(after):
	}
	err = handler.exec(ctx, task)
	if err != nil {
		logger.CtxErrorf(ctx, "handler.exec err: %s", err.Error())
		return
	}
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

func (m *TaskManager) interruptTask(ctx context.Context, task *Task) error {
	task.markPending(0)
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

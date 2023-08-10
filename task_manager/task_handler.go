package task_manager

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golib/concurrency/gopool"
	"golib/etcd_helper"
	"golib/logger"
	"reflect"
	"sync"
	"time"
)

type TaskFuncType func(ctx context.Context, params interface{}) (err error)

type TaskHandler struct {
	Name       TaskHandlerName
	Config     *TaskConfig
	TaskFunc   TaskFuncType
	ParamsType reflect.Type

	runningTasks sync.Map

	taskHandlerGoPool gopool.Pool // 控制handler的并发度
}

func NewTaskHandler(name TaskHandlerName, config *TaskConfig, taskFunc TaskFuncType, paramsType interface{}) *TaskHandler {
	return &TaskHandler{
		Name:              name,
		Config:            config,
		TaskFunc:          taskFunc,
		ParamsType:        reflect.TypeOf(paramsType),
		taskHandlerGoPool: gopool.NewPool("taskHandlerGoPool", gopool.NewConfig(gopool.WithConfigCap(100))),
	}
}

func (handler *TaskHandler) exec(ctx context.Context, task *Task) (err error) {
	if handler.Config.Timeout > 0 {
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithTimeout(ctx, handler.Config.Timeout)
		defer cancelFunc()
	}

	errChan := make(chan error)
	handler.taskHandlerGoPool.Go(func() {
		select {
		case <-ctx.Done():
			return
		default:
		}
		errChan <- handler.TaskFunc(ctx, task.Params)
	})
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-errChan:
	}
	if err != nil {
		if task.ExecCount >= handler.Config.Retry {
			task.markFail()
		} else {
			task.markPending(handler.Config.DelayTime)
		}
		return err
	}
	task.markSuccess()
	return nil
}

func (handler *TaskHandler) handleTask(ctx context.Context, task *Task) {
	ctx = logger.CtxWithTraceId(ctx, task.TraceId)
	switch task.Status {
	case pending:
		handler.taskHandlerGoPool.Go(func() {
			handler.handlePendingTask(ctx, task)
		})
	}
}

func (handler *TaskHandler) handlePendingTask(ctx context.Context, task *Task) {
	err := handler.preemptTask(ctx, task)
	if err != nil {
		logger.CtxInfof(ctx, "preemptTask fail: %+v", err)
		return
	}
	// 抢占到了，执行完后需要释放
	defer func() {
		err := handler.completeTask(ctx, task)
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

func (handler *TaskHandler) preemptTask(ctx context.Context, task *Task) (err error) {
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
	handler.runningTasks.Store(task.getTaskKey().String(), task)
	return nil
}

func (handler *TaskHandler) completeTask(ctx context.Context, task *Task) error {
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
	handler.runningTasks.Delete(task.getTaskKey().String())
	return nil
}

func (handler *TaskHandler) Close() (err error) {
	err = handler.interruptRunningTask()
	if err != nil {
		return err
	}
	return nil
}

func (handler *TaskHandler) interruptRunningTask() error {
	handler.runningTasks.Range(func(key, value interface{}) bool {
		task, ok := value.(*Task)
		if !ok {
			return true
		}
		ctx := etcd_helper.BindContext(logger.CtxWithTraceId(context.Background(), task.TraceId))
		task.markPending(0)
		_, err := etcd_helper.Put(ctx, task.getTaskKey().String(), task.encode())
		if err != nil {
			logger.CtxErrorf(ctx, "close interruptRunningTask fail: %+v", err)
		}
		// 有错误也不暂停，继续中断其他任务
		return true
	})
	return nil
}

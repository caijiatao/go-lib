package task_manager

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golib/concurrency/gopool"
	"golib/etcd_helper"
	"golib/logger"
	"reflect"
	"runtime"
	"sync"
	"time"
)

type TaskHandler struct {
	Name     TaskHandlerName
	Config   *TaskConfig
	TaskFunc interface{}

	runningTasks sync.Map

	taskHandlerGoPool gopool.Pool // 控制handler的并发度
}

func NewTaskHandler(name TaskHandlerName, config *TaskConfig, taskFunc interface{}) *TaskHandler {
	t := reflect.TypeOf(taskFunc)
	if t.Kind() != reflect.Func {
		panic(fmt.Sprintf("%s 非法方法", runtime.FuncForPC(reflect.ValueOf(taskFunc).Pointer()).Name()))
	}
	inNum := t.NumIn()
	if inNum != 2 {
		panic(fmt.Sprintf("%s 非法方法入参", runtime.FuncForPC(reflect.ValueOf(taskFunc).Pointer()).Name()))
	}
	if t.In(0).Name() != "Context" {
		panic(fmt.Sprintf("%s 非法方法入参", runtime.FuncForPC(reflect.ValueOf(taskFunc).Pointer()).Name()))
	}
	outNum := t.NumOut()
	if outNum != 1 {
		panic(fmt.Sprintf("%s 非法方法出参", runtime.FuncForPC(reflect.ValueOf(taskFunc).Pointer()).Name()))
	}
	if t.Out(0).Name() != "error" {
		panic(fmt.Sprintf("%s 方法出参应为error", runtime.FuncForPC(reflect.ValueOf(taskFunc).Pointer()).Name()))
	}

	return &TaskHandler{
		Name:              name,
		Config:            config,
		TaskFunc:          taskFunc,
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
		reValues := make([]reflect.Value, 0)
		reValueCtx := reflect.ValueOf(ctx)
		reValues = append(reValues, reValueCtx)

		fT := reflect.TypeOf(handler.TaskFunc)
		v := reflect.New(fT.In(1))
		paramBytes, err := json.Marshal(task.Params)
		if err != nil {
			errChan <- err
		}
		err = json.Unmarshal(paramBytes, v.Interface())
		if err != nil {
			errChan <- err
		}
		reValues = append(reValues, v.Elem())

		fv := reflect.ValueOf(handler.TaskFunc)
		resVal := fv.Call(reValues)
		//NewTaskHandler做了前置判断
		if len(resVal) == 1 {
			if !resVal[0].IsNil() {
				errChan <- errors.New(fmt.Sprintf("%v", resVal[0]))
			} else {
				errChan <- nil
			}
		}

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
			task.markPending()
		}
		return err
	}
	task.markSuccess()
	return nil
}

func (handler *TaskHandler) handleTask(ctx context.Context, task *Task) {
	ctx = logger.CtxWithTraceId(ctx, task.TraceId)
	task.ctx, task.cancelCtx = context.WithCancel(ctx)
	switch task.Status {
	case pending:
		handler.taskHandlerGoPool.Go(func() {
			handler.handlePendingTask(ctx, task)
		})
	case del:
		value, ok := handler.runningTasks.Load(task.getTaskKey().String())
		if !ok {
			return
		}
		runningTask, ok := value.(*Task)
		if !ok {
			return
		}
		// 删除正在内存中跑的任务，再去删除ETCD中的任务
		runningTask.cancelCtx()
		handler.runningTasks.Delete(task.getTaskKey().String())
		_, err := etcd_helper.Context(ctx).
			Txn(ctx).
			If(clientv3.Compare(clientv3.Version(task.getTaskKey().String()), "=", task.TaskVersion)).
			Then(clientv3.OpDelete(task.getTaskKey().String())).
			Commit()
		if err != nil {
			logger.CtxErrorf(ctx, "delete task fail: %+v", err)
		}
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
	case <-task.ctx.Done(): // 任务的context被取消则直接结束
		return
	}
	err = handler.exec(task.ctx, task)
	if err != nil {
		logger.CtxErrorf(ctx, "handler.exec err: %s, task:%+v", err.Error(), task)
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
	_, ok := handler.runningTasks.Load(task.getTaskKey().String())
	if !ok {
		// 任务已经被删除，被强制中断，不需要做任何操作
		return nil
	}

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
		task.markPending()
		// 被打断，不算执行次数
		task.ExecCount--
		_, err := etcd_helper.Put(ctx, task.getTaskKey().String(), task.encode())
		if err != nil {
			logger.CtxErrorf(ctx, "close interruptRunningTask fail: %+v", err)
		}
		// 有错误也不暂停，继续中断其他任务
		return true
	})
	return nil
}

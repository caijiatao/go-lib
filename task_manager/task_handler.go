package task_manager

import (
	"context"
	"golib/concurrency/gopool"
	"reflect"
)

type TaskFuncType func(ctx context.Context, params interface{}) (err error)

type TaskHandler struct {
	Name       TaskHandlerName
	Config     *TaskConfig
	TaskFunc   TaskFuncType
	ParamsType reflect.Type

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

func (h *TaskHandler) exec(ctx context.Context, task *Task) (err error) {
	if h.Config.Timeout > 0 {
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithTimeout(ctx, h.Config.Timeout)
		defer cancelFunc()
	}

	errChan := make(chan error)
	h.taskHandlerGoPool.Go(func() {
		select {
		case <-ctx.Done():
			return
		default:
		}
		errChan <- h.TaskFunc(ctx, task.Params)
	})
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-errChan:
	}
	if err != nil {
		if task.ExecCount >= h.Config.Retry {
			task.markFail()
		} else {
			task.markPending(h.Config.DelayTime)
		}
		return err
	}
	task.markSuccess()
	return nil
}

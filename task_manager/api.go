package task_manager

import (
	"context"
	"github.com/pkg/errors"
	"golib/etcd_helper"
)

func RegisterTaskHandler(handler *TaskHandler) (err error) {
	err = defaultTaskManager.RegisterHandler(handler)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// SubmitTask
//
//	@Description: 提交任务等待执行
func SubmitTask(ctx context.Context, task *Task) (err error) {
	if task == nil {
		return nil
	}
	_, err = etcd_helper.Put(ctx, task.getTaskKey().String(), task.encode())
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func AllTaskDone(ctx context.Context) bool {
	isAllDone := true
	defaultTaskManager.runningTasks.Range(func(key, value interface{}) bool {
		isAllDone = false
		return false // 返回 false 以停止遍历
	})
	return isAllDone
}

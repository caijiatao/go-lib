package task_manager

import (
	"context"
	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golib/etcd_helper"
	"golib/logger"
)

func RegisterTaskHandler(handler *TaskHandler) (err error) {
	err = globalTaskManager.RegisterHandler(handler)
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
	return SubmitTasks(ctx, []*Task{task})
}

func SubmitTasks(ctx context.Context, tasks []*Task) (err error) {
	if len(tasks) == 0 {
		return nil
	}
	traceId := logger.CtxTraceID(ctx)
	for i, _ := range tasks {
		tasks[i].TraceId = traceId
	}

	putOps := make([]clientv3.Op, 0, len(tasks))
	for _, task := range tasks {
		putOp := clientv3.OpPut(task.getTaskKey().String(), task.encode())
		putOps = append(putOps, putOp)
	}
	txnResponse, err := etcd_helper.Context(ctx).Txn(ctx).If().Then(putOps...).Commit()
	if err != nil {
		return errors.WithStack(err)
	}
	if !txnResponse.Succeeded {
		return errors.New("SubmitTasks put err")
	}
	return nil
}

// DeleteTask
//
//	@Description: 删除任务则会中断现在在内存中还没执行的任务，有延迟执行占用内存的才需要强制删除
func DeleteTask(ctx context.Context, task *Task) (err error) {
	task.Status = del
	putOp := clientv3.OpPut(task.getTaskKey().String(), task.encode())
	txnResponse, err := etcd_helper.Context(ctx).
		Txn(ctx).
		If(clientv3.Compare(clientv3.Value(task.getTaskKey().String()), "!=", "")).
		Then(putOp).
		Commit()
	if err != nil {
		return errors.WithStack(err)
	}

	if !txnResponse.Succeeded {
		logger.CtxErrorf(ctx, "Task maneger DeleteTask :%+v, task:%+v", txnResponse, task)
		return errors.New("DeleteTask put err")
	}
	return nil
}

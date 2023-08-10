package task_manager

import (
	"context"
	errors2 "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"golib/etcd_helper"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	etcd_helper.InitTestSuite()
	Init()
	m.Run()
}

func TestTaskManager(t *testing.T) {
	var (
		value       int64
		err         error
		handlerName TaskHandlerName = "test_task"
	)

	type testParamType struct {
		AddValue int64
	}

	taskHandler := NewTaskHandler(handlerName, NewTaskConfig(), func(ctx context.Context, params interface{}) (err error) {
		p := params.(*testParamType)
		value += p.AddValue
		return nil
	}, testParamType{})
	err = RegisterTaskHandler(taskHandler)
	assert.Nil(t, err)
	err = RegisterTaskHandler(taskHandler)
	assert.Equal(t, HandlerNameRepeatErr, errors2.Cause(err))

	var (
		taskID   = "1"
		testTask = NewTask(taskID, handlerName, &testParamType{AddValue: int64(1)})
	)
	ctx := etcd_helper.BindContext(context.Background())
	err = SubmitTask(ctx, testTask)
	assert.Nil(t, err)

	time.Sleep(1 * time.Second)

	for {
		done := AllTaskDone(ctx)
		if done {
			break
		}
		time.Sleep(2 * time.Second)
	}
	assert.Equal(t, int64(1), value)
}

func TestTaskRetrySuccess(t *testing.T) {
	var (
		value       int64
		err         error
		handlerName TaskHandlerName = "test_task_retry"
		retryCount                  = 3
	)

	// 1.构造handler
	type testParamType struct {
		AddValue int64
	}
	taskHandler := NewTaskHandler(handlerName, NewTaskConfig(WithTaskConfigRetry(retryCount)), func(ctx context.Context, params interface{}) (err error) {
		p := params.(*testParamType)
		// 最后一次重试成功
		if value == int64(retryCount-1) {
			return nil
		}
		value += p.AddValue
		return errors2.New("test error")
	}, testParamType{})
	err = RegisterTaskHandler(taskHandler)
	assert.Nil(t, err)

	// 2.构造task
	var (
		taskID   = "1"
		testTask = NewTask(taskID, handlerName, &testParamType{AddValue: int64(1)})
	)
	ctx := etcd_helper.BindContext(context.Background())
	err = SubmitTask(ctx, testTask)
	assert.Nil(t, err)

	// 3.等待任务完成
	for {
		getResponse, err := etcd_helper.Get(ctx, testTask.getTaskKey().String())
		assert.Nil(t, err)
		if len(getResponse.Kvs) == 0 {
			break
		}
	}
	assert.Equal(t, int64(retryCount-1), value)
}

func TestTaskRetryFail(t *testing.T) {
	var (
		value       int64
		err         error
		handlerName TaskHandlerName = "test_task_retry_fail"
		retryCount                  = 3
	)

	// 1.构造handler
	type testParamType struct {
		AddValue int64
	}
	taskHandler := NewTaskHandler(handlerName, NewTaskConfig(WithTaskConfigRetry(retryCount)), func(ctx context.Context, params interface{}) (err error) {
		p := params.(*testParamType)
		value += p.AddValue
		return errors2.New("test error")
	}, testParamType{})
	err = RegisterTaskHandler(taskHandler)
	assert.Nil(t, err)

	// 2.提交任务
	var (
		taskID   = "1"
		testTask = NewTask(taskID, handlerName, &testParamType{AddValue: int64(1)})
	)
	ctx := etcd_helper.BindContext(context.Background())
	err = SubmitTask(ctx, testTask)
	assert.Nil(t, err)

	// 3.判断结果
	isFail := false
	for !isFail {
		getResponse, err := etcd_helper.Get(ctx, testTask.getTaskKey().String())
		assert.Nil(t, err)
		for _, kv := range getResponse.Kvs {
			task := decodeTask(string(kv.Value), kv.Version, taskHandler.ParamsType)
			isFail = task.Status == fail
		}
	}
	assert.Equal(t, int64(retryCount), value)
}

func TestTaskExecTimeout(t *testing.T) {
	var (
		value       int64
		err         error
		handlerName TaskHandlerName = "test_task_timeout"
	)

	// 1.构造handler
	type testParamType struct {
		AddValue int64
	}
	taskHandler := NewTaskHandler(handlerName, NewTaskConfig(WithTaskConfigTimeout(time.Second), WithTaskConfigRetry(3)), func(ctx context.Context, params interface{}) (err error) {
		p := params.(*testParamType)
		time.Sleep(2 * time.Second)
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:

		}
		value += p.AddValue
		return nil
	}, testParamType{})
	err = RegisterTaskHandler(taskHandler)
	assert.Nil(t, err)

	// 2.提交任务
	var (
		taskID   = "1"
		testTask = NewTask(taskID, handlerName, &testParamType{AddValue: int64(1)})
	)
	ctx := etcd_helper.BindContext(context.Background())
	err = SubmitTask(ctx, testTask)
	assert.Nil(t, err)

	// 3.判断结果
	isTimeout := false
	for !isTimeout {
		getResponse, err := etcd_helper.Get(ctx, testTask.getTaskKey().String())
		assert.Nil(t, err)
		for _, kv := range getResponse.Kvs {
			task := decodeTask(string(kv.Value), kv.Version, taskHandler.ParamsType)
			isTimeout = task.Status == fail
		}
	}
	assert.Equal(t, int64(0), value)
}

func TestExecTaskDelayTime(t *testing.T) {
	var (
		value       int64
		err         error
		handlerName TaskHandlerName = "test_task_delay"
	)

	// 1.构造handler
	type testParamType struct {
		AddValue int64
	}
	taskHandler := NewTaskHandler(handlerName, NewTaskConfig(WithTaskConfigDelayTime(5*time.Second)), func(ctx context.Context, params interface{}) (err error) {
		p := params.(*testParamType)
		value += p.AddValue
		return nil
	}, testParamType{})
	err = RegisterTaskHandler(taskHandler)
	assert.Nil(t, err)

	// 2.提交任务
	var (
		taskID   = "1"
		testTask = NewTask(taskID, handlerName, &testParamType{AddValue: int64(1)})
	)
	ctx := etcd_helper.BindContext(context.Background())
	err = SubmitTask(ctx, testTask)
	assert.Nil(t, err)

	// 3.判断结果
	for {
		getResponse, err := etcd_helper.Get(ctx, testTask.getTaskKey().String())
		assert.Nil(t, err)
		if len(getResponse.Kvs) == 0 {
			break
		}
	}
	assert.Equal(t, int64(1), value)
}

func TestPreemptTask(t *testing.T) {

}

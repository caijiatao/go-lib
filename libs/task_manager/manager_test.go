package task_manager

import (
	"context"
	"fmt"
	errors2 "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	etcd_helper2 "golib/libs/etcd_helper"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	etcd_helper2.InitTestSuite()
	m.Run()
}

func testReset(t *testing.T) {
	err := globalTaskManager.Close()
	assert.Nil(t, err)
	globalTaskManager = nil
	defaultTaskManagerInitOnce = sync.Once{}
}

func TestTaskManager(t *testing.T) {
	Init("")
	defer testReset(t)
	var (
		value       int64
		err         error
		handlerName TaskHandlerName = "test_task"
	)

	type testParamType struct {
		AddValue int64
	}

	taskHandler := NewTaskHandler(handlerName, NewTaskConfig(), func(ctx context.Context, params *testParamType) (err error) {
		p := params
		value += p.AddValue
		return nil
	})
	err = RegisterTaskHandler(taskHandler)
	assert.Nil(t, err)
	err = RegisterTaskHandler(taskHandler)
	assert.Equal(t, HandlerNameRepeatErr, errors2.Cause(err))

	var (
		taskID   = fmt.Sprintf("%d", time.Now().Unix())
		testTask = NewTask(taskID, handlerName, &testParamType{AddValue: int64(1)})
	)
	ctx := etcd_helper2.BindContext(context.Background())
	err = SubmitTask(ctx, testTask)
	assert.Nil(t, err)

	for {
		getResponse, err := etcd_helper2.Get(ctx, testTask.getTaskKey().String())
		assert.Nil(t, err)
		if len(getResponse.Kvs) == 0 {
			break
		}
	}
	assert.Equal(t, int64(1), value)
}

func TestTaskRetrySuccess(t *testing.T) {
	Init("")
	defer testReset(t)
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
	taskHandler := NewTaskHandler(handlerName, NewTaskConfig(WithTaskConfigRetry(retryCount)), func(ctx context.Context, params *testParamType) (err error) {
		p := params
		// 最后一次重试成功
		if value == int64(retryCount-1) {
			return nil
		}
		value += p.AddValue
		return errors2.New("test error")
	})
	err = RegisterTaskHandler(taskHandler)
	assert.Nil(t, err)

	// 2.构造task
	var (
		taskID   = fmt.Sprintf("%d", time.Now().Unix())
		testTask = NewTask(taskID, handlerName, &testParamType{AddValue: int64(1)})
	)
	ctx := etcd_helper2.BindContext(context.Background())
	err = SubmitTask(ctx, testTask)
	assert.Nil(t, err)

	// 3.等待任务完成
	for {
		getResponse, err := etcd_helper2.Get(ctx, testTask.getTaskKey().String())
		assert.Nil(t, err)
		if len(getResponse.Kvs) == 0 {
			break
		}
	}
	assert.Equal(t, int64(retryCount-1), value)
}

func TestTaskRetryFail(t *testing.T) {
	Init("")
	defer testReset(t)
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
	taskHandler := NewTaskHandler(handlerName, NewTaskConfig(WithTaskConfigRetry(retryCount)), func(ctx context.Context, params *testParamType) (err error) {
		p := params
		value += p.AddValue
		return errors2.New("test error")
	})
	err = RegisterTaskHandler(taskHandler)
	assert.Nil(t, err)

	// 2.提交任务
	var (
		taskID   = fmt.Sprintf("%d", time.Now().Unix())
		testTask = NewTask(taskID, handlerName, &testParamType{AddValue: int64(1)})
	)
	ctx := etcd_helper2.BindContext(context.Background())
	err = SubmitTask(ctx, testTask)
	assert.Nil(t, err)

	// 3.判断结果
	isFail := false
	for !isFail {
		getResponse, err := etcd_helper2.Get(ctx, testTask.getTaskKey().String())
		assert.Nil(t, err)
		for _, kv := range getResponse.Kvs {
			task := &Task{}
			err = task.decode(string(kv.Value))
			if err != nil {
				assert.Nil(t, err)
			}
			task.TaskVersion = kv.Version
			isFail = task.Status == fail
		}
	}
	assert.Equal(t, int64(retryCount), value)
}

func TestTaskExecTimeout(t *testing.T) {
	Init("")
	defer testReset(t)
	var (
		value       int64
		err         error
		handlerName TaskHandlerName = "test_task_timeout"
	)

	// 1.构造handler
	type testParamType struct {
		AddValue int64
	}
	taskHandler := NewTaskHandler(handlerName, NewTaskConfig(WithTaskConfigTimeout(time.Second), WithTaskConfigRetry(3)), func(ctx context.Context, params *testParamType) (err error) {
		p := params
		time.Sleep(2 * time.Second)
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:

		}
		value += p.AddValue
		return nil
	})
	err = RegisterTaskHandler(taskHandler)
	assert.Nil(t, err)

	// 2.提交任务
	var (
		taskID   = fmt.Sprintf("%d", time.Now().Unix())
		testTask = NewTask(taskID, handlerName, &testParamType{AddValue: int64(1)})
	)
	ctx := etcd_helper2.BindContext(context.Background())
	err = SubmitTask(ctx, testTask)
	assert.Nil(t, err)

	// 3.判断结果
	isTimeout := false
	for !isTimeout {
		getResponse, err := etcd_helper2.Get(ctx, testTask.getTaskKey().String())
		assert.Nil(t, err)
		for _, kv := range getResponse.Kvs {
			task := &Task{}
			err = task.decode(string(kv.Value))
			if err != nil {
				assert.Nil(t, err)
			}
			task.TaskVersion = kv.Version
			isTimeout = task.Status == fail
		}
	}
	assert.Equal(t, int64(0), value)
}

func TestExecTaskDelayTime(t *testing.T) {
	Init("")
	defer testReset(t)
	var (
		value       int64
		err         error
		handlerName TaskHandlerName = "test_task_delay"
	)

	// 1.构造handler
	type testParamType struct {
		AddValue int64
	}
	taskHandler := NewTaskHandler(handlerName, NewTaskConfig(WithTaskConfigDelayTime(1*time.Second)), func(ctx context.Context, params *testParamType) (err error) {
		p := params
		value += p.AddValue
		return nil
	})
	err = RegisterTaskHandler(taskHandler)
	assert.Nil(t, err)

	// 2.提交任务
	var (
		taskID   = fmt.Sprintf("%d", time.Now().Unix())
		testTask = NewTask(taskID, handlerName, &testParamType{AddValue: int64(1)})
	)
	ctx := etcd_helper2.BindContext(context.Background())
	err = SubmitTask(ctx, testTask)
	assert.Nil(t, err)

	// 3.判断结果
	for {
		getResponse, err := etcd_helper2.Get(ctx, testTask.getTaskKey().String())
		assert.Nil(t, err)
		if len(getResponse.Kvs) == 0 {
			break
		}
	}
	assert.Equal(t, int64(1), value)
}

func TestInterruptTask(t *testing.T) {
	Init("")
	defer testReset(t)
	var (
		err error
	)
	type testParamType struct {
		AddValue int64
	}
	var (
		handlerName     = TaskHandlerName("test_task_interrupt")
		testHandlerFunc = func() interface{} {
			return func(ctx context.Context, params interface{}) (err error) {
				_ = globalTaskManager.Close()
				return nil
			}
		}
		taskCloseHandler = NewTaskHandler(handlerName, NewTaskConfig(), testHandlerFunc())
	)

	err = RegisterTaskHandler(taskCloseHandler)
	assert.Nil(t, err)
	ctx := etcd_helper2.BindContext(context.Background())
	task := NewTask("1", handlerName, &testParamType{AddValue: int64(1)})
	err = SubmitTask(ctx, task)
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)

	getResponse, err := etcd_helper2.Get(ctx, task.getTaskKey().String())
	assert.Nil(t, err)
	for _, kv := range getResponse.Kvs {
		task := &Task{}
		err = task.decode(string(kv.Value))
		if err != nil {
			assert.Nil(t, err)
		}
		task.TaskVersion = kv.Version
		assert.Equal(t, task.Status, pending)
	}

	_, err = etcd_helper2.Delete(ctx, task.getTaskKey().String())
	assert.Nil(t, err)
}

func TestDeleteDelayTask(t *testing.T) {
	Init("")
	defer testReset(t)
	var (
		value       int64
		err         error
		handlerName TaskHandlerName = "test_delay_task_delete"
	)

	// 1.构造handler
	type testParamType struct {
		AddValue int64
	}
	taskHandler := NewTaskHandler(handlerName, NewTaskConfig(WithTaskConfigDelayTime(5*time.Second)), func(ctx context.Context, params *testParamType) (err error) {
		p := params
		value += p.AddValue
		return nil
	})
	err = RegisterTaskHandler(taskHandler)
	assert.Nil(t, err)

	// 2.提交任务
	var (
		taskID   = fmt.Sprintf("%d", time.Now().Unix())
		testTask = NewTask(taskID, handlerName, &testParamType{AddValue: int64(1)})
	)
	ctx := etcd_helper2.BindContext(context.Background())
	err = SubmitTask(ctx, testTask)
	assert.Nil(t, err)

	// 3.强制删除任务，立即结束
	for {
		// 等待任务已经被调度
		_, ok := globalTaskManager.getHandler(handlerName).runningTasks.Load(testTask.getTaskKey().String())
		if ok {
			break
		}
	}
	err = DeleteTask(ctx, testTask)
	assert.Nil(t, err)

	// 4.判断结果
	time.Sleep(time.Second)
	getResponse, err := etcd_helper2.Get(ctx, testTask.getTaskKey().String())
	assert.Nil(t, err)
	assert.Len(t, getResponse.Kvs, 0)
	assert.Equal(t, int64(0), value) // 没有执行任务
	runningTask, ok := globalTaskManager.getHandler(handlerName).runningTasks.Load(testTask.getTaskKey().String())
	assert.False(t, ok)
	assert.Nil(t, runningTask)
}

func TestQuitAndRestart(t *testing.T) {
	// 1.第一次初始化
	Init("")
	ctx := etcd_helper2.BindContext(context.Background())
	_, err := etcd_helper2.Delete(ctx, taskKeyPrefix, clientv3.WithPrefix())
	assert.Nil(t, err)
	var (
		value       int64
		handlerName TaskHandlerName = "test_delay_task_delete"
	)

	// 2.构造任务
	type testParamType struct {
		AddValue int64
	}
	taskHandler := NewTaskHandler(handlerName, NewTaskConfig(WithTaskConfigDelayTime(time.Second)), func(ctx context.Context, params *testParamType) (err error) {
		p := params
		value += p.AddValue
		return nil
	})
	err = RegisterTaskHandler(taskHandler)
	assert.Nil(t, err)
	var (
		taskID   = fmt.Sprintf("%d", time.Now().Unix())
		testTask = NewTask(taskID, handlerName, &testParamType{AddValue: int64(1)})
	)

	err = SubmitTask(ctx, testTask)
	assert.Nil(t, err)
	for {
		// 等待任务被抢占到
		_, ok := globalTaskManager.getHandler(handlerName).runningTasks.Load(testTask.getTaskKey().String())
		if ok {
			break
		}
	}

	// 3.模拟程序退出 ，判断执行结果
	testReset(t)
	assert.Equal(t, int64(0), value) // 第一次退出后没有执行
	getResponse, err := etcd_helper2.Get(ctx, testTask.getTaskKey().String())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(getResponse.Kvs))
	// 将不需要判断的字段处理成相同的值
	notExecTask := &Task{}
	err = notExecTask.decode(string(getResponse.Kvs[0].Value))
	assert.Nil(t, err)
	testTask.TaskVersion = notExecTask.TaskVersion
	testTask.TraceId = notExecTask.TraceId
	testTask.LastUpdateTime = notExecTask.LastUpdateTime
	assert.Equal(t, string(getResponse.Kvs[0].Value), testTask.encode())

	// 4.第二次初始化
	//patches := gomonkey.ApplyGlobalVar(&scanTaskDelayTime, time.Second) // 缩短任务扫描时间，快速执行任务
	//defer patches.Reset()
	Init("")
	err = RegisterTaskHandler(taskHandler)
	assert.Nil(t, err)
	// 等待执行完成判断结果
	for {
		getResponse, err = etcd_helper2.Get(ctx, testTask.getTaskKey().String())
		assert.Nil(t, err)
		if len(getResponse.Kvs) == 0 {
			break
		}
	}
	assert.Equal(t, int64(1), value) // 第二次执行成功
}

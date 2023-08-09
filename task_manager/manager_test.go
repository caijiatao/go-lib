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

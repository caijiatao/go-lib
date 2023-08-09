package task_manager

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	taskKeyPrefix = "/task_manager"
)

type taskKey string

const (
	taskKeyPrefixIndex = iota + 1
	taskKeyHandleNameIndex
)

func (k taskKey) String() string {
	return string(k)
}

func (k taskKey) getHandlerName() TaskHandlerName {
	values := strings.Split(string(k), "/")
	if len(values) > taskKeyHandleNameIndex {
		return TaskHandlerName(values[taskKeyHandleNameIndex])
	}
	return ""
}

type TaskHandlerName string
type TaskStatus int

const (
	pending TaskStatus = iota + 1
	running
	success
	fail
)

type Task struct {
	TaskId      string
	HandlerName TaskHandlerName
	Params      interface{} // 任务参数

	Status      TaskStatus
	ExecCount   int
	TaskVersion int64 // 任务版本号，用于任务抢占

	CreateTime     time.Time
	LastUpdateTime time.Time
}

func NewTask(taskId string, handlerName TaskHandlerName, params interface{}) *Task {
	return &Task{
		TaskId:         taskId,
		HandlerName:    handlerName,
		Params:         params,
		Status:         pending,
		ExecCount:      0,
		CreateTime:     time.Now(),
		LastUpdateTime: time.Now(),
	}
}

func DecodeTask(taskString string, taskVersion int64, paramType reflect.Type) *Task {
	t := &Task{
		Params: reflect.New(paramType).Interface(),
	}
	err := json.Unmarshal([]byte(taskString), t)
	if err != nil {
		return nil
	}
	t.TaskVersion = taskVersion
	return t
}

func (t *Task) encode() string {
	taskString, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(taskString)
}

func (t *Task) getTaskKey() taskKey {
	k := fmt.Sprintf("%s/%s/%s", taskKeyPrefix, t.HandlerName, t.TaskId)
	return taskKey(k)
}

func (t *Task) markRunning() {
	t.ExecCount++
	t.LastUpdateTime = time.Now()
	t.Status = running
}

func (t *Task) markPending() {
	t.LastUpdateTime = time.Now()
	t.Status = pending
}

func (t *Task) markFail() {
	t.LastUpdateTime = time.Now()
	t.Status = fail
}

func (t *Task) markSuccess() {
	t.LastUpdateTime = time.Now()
	t.Status = success
}

type TaskConfig struct {
	Timeout   time.Duration // 任务超时时间
	DelayTime time.Duration // 延迟多长时间再执行
	Retry     int
}

func NewTaskConfig() *TaskConfig {
	return &TaskConfig{
		Timeout:   0,
		DelayTime: 0,
		Retry:     0,
	}
}

type TaskFuncType func(ctx context.Context, params interface{}) (err error)

type TaskHandler struct {
	Name       TaskHandlerName
	Config     *TaskConfig
	TaskFunc   TaskFuncType
	ParamsType reflect.Type
}

func NewTaskHandler(name TaskHandlerName, config *TaskConfig, taskFunc TaskFuncType, paramsType interface{}) *TaskHandler {
	return &TaskHandler{
		Name:       name,
		Config:     config,
		TaskFunc:   taskFunc,
		ParamsType: reflect.TypeOf(paramsType),
	}
}

func (h *TaskHandler) Exec(ctx context.Context, task *Task) (err error) {
	err = h.TaskFunc(ctx, task.Params)
	if err != nil {
		if task.ExecCount >= h.Config.Retry {
			task.markFail()
		} else {
			task.markPending()
		}
		return err
	}
	task.markSuccess()
	return nil
}

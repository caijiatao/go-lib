package task_manager

import (
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
	TaskVersion int64  // 任务版本号，用于任务抢占
	TraceId     string // 日志id

	CreateTime      time.Time
	NextExecuteTime time.Time
	LastUpdateTime  time.Time
}

func NewTask(taskId string, handlerName TaskHandlerName, params interface{}) *Task {
	handler := globalTaskManager.getHandler(handlerName)
	if handler == nil {
		return nil
	}
	return &Task{
		TaskId:          taskId,
		HandlerName:     handlerName,
		Params:          params,
		Status:          pending,
		ExecCount:       0,
		CreateTime:      time.Now(),
		NextExecuteTime: time.Now().Add(handler.Config.DelayTime),
		LastUpdateTime:  time.Now(),
	}
}

func decodeTask(taskString string, taskVersion int64, paramType reflect.Type) *Task {
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

func (t *Task) markPending(delayTime time.Duration) {
	t.LastUpdateTime = time.Now()
	t.NextExecuteTime = time.Now().Add(delayTime)
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

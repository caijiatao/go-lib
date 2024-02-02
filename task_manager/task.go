package task_manager

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var (
	taskKeyPrefix = "/task_manager"
)

type taskKey string

func (k taskKey) String() string {
	return string(k)
}

func (k taskKey) getHandlerName() TaskHandlerName {
	newKey := strings.Replace(string(k), taskKeyPrefix+"/", "", 1)
	values := strings.Split(newKey, "/")
	if len(values) > 0 {
		return TaskHandlerName(values[0])
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
	del
)

type Task struct {
	ctx       context.Context
	cancelCtx func()

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

func (t *Task) decode(taskString string) error {
	err := json.Unmarshal([]byte(taskString), &t)
	if err != nil {
		return err
	}
	return nil
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

func (t *Task) markDel() {
	t.LastUpdateTime = time.Now()
	t.Status = del
}

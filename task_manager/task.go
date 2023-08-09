package task_manager

import (
	"context"
	"time"
)

type Task struct {
}

type TaskParams interface {
}

type TaskHandler struct {
	Timeout  time.Duration
	TaskFunc func(ctx context.Context, params interface{}) error
}

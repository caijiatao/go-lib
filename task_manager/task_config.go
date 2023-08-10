package task_manager

import "time"

type TaskConfigOpt func(config *TaskConfig)

type TaskConfig struct {
	Timeout   time.Duration // 任务超时时间
	DelayTime time.Duration // 延迟多长时间再执行
	Retry     int
}

func WithTaskConfigTimeout(timeout time.Duration) TaskConfigOpt {
	return func(config *TaskConfig) {
		config.Timeout = timeout
	}
}

func WithTaskConfigDelayTime(delayTime time.Duration) TaskConfigOpt {
	return func(config *TaskConfig) {
		config.DelayTime = delayTime
	}
}

func WithTaskConfigRetry(retry int) TaskConfigOpt {
	return func(config *TaskConfig) {
		config.Retry = retry
	}
}

func NewTaskConfig(opts ...TaskConfigOpt) *TaskConfig {
	config := &TaskConfig{}
	for _, opt := range opts {
		opt(config)
	}
	return config
}

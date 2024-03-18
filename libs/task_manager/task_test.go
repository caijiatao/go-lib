package task_manager

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_taskKey_getHandlerName(t *testing.T) {
	tests := []struct {
		name          string
		k             taskKey
		want          TaskHandlerName
		taskKeyPrefix string
	}{
		{
			name:          "test1",
			k:             taskKey("/task_manager/test1/123"),
			want:          TaskHandlerName("test1"),
			taskKeyPrefix: "/task_manager",
		},
		{
			name:          "test2",
			k:             taskKey("/task_manager/13726285213/test2/123"),
			want:          TaskHandlerName("test2"),
			taskKeyPrefix: "/task_manager/13726285213",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskKeyPrefix = tt.taskKeyPrefix
			assert.Equalf(t, tt.want, tt.k.getHandlerName(), "getHandlerName()")
		})
	}
}

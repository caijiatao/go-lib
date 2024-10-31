package good_util

import "testing"

func Test_loChannelDispatcher(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test loChannelDispatcher",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loChannelDispatcher()
		})
	}
}

func Test_loBuffer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test loBuffer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loBuffer()
		})
	}
}

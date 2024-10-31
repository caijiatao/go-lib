package good_util

import "testing"

func Test_loDuration(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test_loDuration",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loDuration()
		})
	}
}

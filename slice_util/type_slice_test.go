package slice_util

import (
	"fmt"
	"testing"
)

func TestDeltasFIFO_Add(t *testing.T) {
	type args struct {
		key   string
		delta Delta
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				key: "test",
				delta: Delta{
					Type:  "test",
					Value: "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewDeltasFIFO()
			f.Add(tt.args.key, tt.args.delta)
			fmt.Println(f.Get(tt.args.key))
		})
	}
}

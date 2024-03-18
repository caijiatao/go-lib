package string_util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitAndTransferStrSliceToInterfaceSlice(t *testing.T) {
	type args struct {
		stringArray []string
		maxSize     int
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "边界切分",
			args: args{
				stringArray: []string{"1", "2", "3", "4", "5"},
				maxSize:     2,
			},
			want: []interface{}{
				[]string{"1", "2"},
				[]string{"3", "4"},
				[]string{"5"},
			},
		},
		{
			name: "小于最大size",
			args: args{
				stringArray: []string{"1", "2", "3", "4", "5"},
				maxSize:     10,
			},
			want: []interface{}{
				[]string{"1", "2", "3", "4", "5"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := SplitAndTransferStrSliceToInterfaceSlice(tt.args.stringArray, tt.args.maxSize)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestFilterStrings(t *testing.T) {
	type args struct {
		originStrings []string
		filterStrings []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "正常过滤",
			args: args{
				originStrings: []string{"1", "2", "3"},
				filterStrings: []string{"1"},
			},
			want: []string{"2", "3"},
		},
		{
			name: "originStrings空值测试",
			args: args{
				originStrings: []string{},
				filterStrings: []string{"1", "2", "3"},
			},
			want: []string{},
		},
		{
			name: "filterStrings空值测试",
			args: args{
				originStrings: []string{"1", "2", "3"},
				filterStrings: []string{},
			},
			want: []string{"1", "2", "3"},
		},
		{
			name: "全部为空值测试",
			args: args{
				originStrings: []string{},
				filterStrings: []string{},
			},
			want: []string{},
		},
		{
			name: "过滤为空值测试",
			args: args{
				originStrings: []string{"1", "2", "3"},
				filterStrings: []string{"1", "2", "3"},
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := FilterStrings(tt.args.originStrings, tt.args.filterStrings)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestSplitStringSlice(t *testing.T) {
	type args struct {
		stringArray []string
		chunkSize   int
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "t1",
			args: args{
				stringArray: []string{"1", "2", "3"},
				chunkSize:   1,
			},
			want: [][]string{
				{"1"},
				{"2"},
				{"3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, SplitStringSlice(tt.args.stringArray, tt.args.chunkSize), "SplitStringSlice(%v, %v)", tt.args.stringArray, tt.args.chunkSize)
		})
	}
}

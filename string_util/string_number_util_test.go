package string_util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindFirstNumberStr(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "字符开头",
			args: args{
				str: "AT202201",
			},
			want: "202201",
		},
		{
			name: "数字开头",
			args: args{
				str: "202202AT",
			},
			want: "202202",
		},
		{
			name: "字符和数字掺杂",
			args: args{
				str: "AT202203AT1234",
			},
			want: "202203",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := FindFirstNumberStr(tt.args.str)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestIsNumber(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "t1",
			args: args{"123"},
			want: true,
		},
		{
			name: "t2",
			args: args{"1sss"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, IsNumber(tt.args.s), "IsNumber(%v)", tt.args.s)
		})
	}
}

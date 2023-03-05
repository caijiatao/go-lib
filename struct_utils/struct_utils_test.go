package struct_utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type aStruct struct {
	Name string
	Male string
}

func (s *aStruct) IsEmpty() bool {
	return s.Male == "" && s.Name == ""
}

type bStruct struct {
	Name string
	Male string
	Age  int
}

func TestCopyIntersectionStruct(t *testing.T) {
	a := &aStruct{
		Name: "derrick",
		Male: "male",
	}
	b := &bStruct{
		Name: "xu",
		Male: "male",
		Age:  100,
	}
	CopyIntersectionStruct(b, a)
	println(a)
	println(b)
}

type complexSt struct {
	A        aStruct
	S        []string
	IntValue int
}

func (c *complexSt) IsEmpty() bool {
	return c.A.IsEmpty() && len(c.S) == 0 && c.IntValue == 0
}

func TestIsStructEmpty(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "正常为空",
			args: args{
				v: complexSt{
					A:        aStruct{},
					S:        make([]string, 0),
					IntValue: 0,
				},
			},
			want: true,
		},
		{
			name: "数组不为空",
			args: args{
				v: complexSt{
					A:        aStruct{},
					S:        []string{"1"},
					IntValue: 0,
				},
			},
			want: false,
		},
		{
			name: "结构体不为空",
			args: args{
				v: complexSt{
					A: aStruct{
						Name: "111",
						Male: "111222",
					},
					S:        []string{},
					IntValue: 0,
				},
			},
			want: false,
		},
		{
			name: "字面量不为空",
			args: args{
				v: complexSt{
					A:        aStruct{},
					S:        make([]string, 0),
					IntValue: 1,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, IsStructEmpty(tt.args.v), "IsStructEmpty(%v)", tt.args.v)
		})
	}
}

func BenchmarkReflectIsStructEmpty(b *testing.B) {
	s := complexSt{
		A:        aStruct{},
		S:        make([]string, 0),
		IntValue: 0,
	}
	for i := 0; i < b.N; i++ {
		IsStructEmpty(s)
	}
}

func BenchmarkNormalIsStructEmpty(b *testing.B) {
	s := complexSt{
		A:        aStruct{},
		S:        make([]string, 0),
		IntValue: 0,
	}
	for i := 0; i < b.N; i++ {
		s.IsEmpty()
	}
}

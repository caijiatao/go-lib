package goasync

import (
	"errors"
	"fmt"
	"testing"
)

func TestPanicErrHandler(t *testing.T) {
	tests := []struct {
		name string
		err  any
	}{
		{
			name: "t1",
			err:  errors.New("111"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				e := PanicErrHandler(r)
				fmt.Println(e)
			}()
			if tt.err != nil {
				panic(tt.err)
			}
		})
	}
}

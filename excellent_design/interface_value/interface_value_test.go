package interface_value

import (
	"fmt"
	"testing"
)

func TestInterfaceValue(t *testing.T) {
	r := &RouteImpl{}
	r.Register(func() {
		fmt.Println("hello")
	})
	r.Register(func() {
		fmt.Println("world")
	})
	for _, f := range r.fs {
		f()
	}
	setFunc(r)
	for _, f := range r.fs {
		f()
	}
}

func setFunc(r Route) {
	r = r.Register(func() {
		fmt.Println("hello1")
	})
}

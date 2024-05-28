package main

import (
	"fmt"
	"reflect"
)

type test interface {
	test()
}

type A struct{}

func (receiver *A) test() {
	println("test")
}

func main() {
	a := &A{}
	aType := reflect.TypeOf(a)

	fmt.Println(aType.Kind())
	fmt.Println(aType.Implements(reflect.TypeOf((*test)(nil)).Elem()))
}

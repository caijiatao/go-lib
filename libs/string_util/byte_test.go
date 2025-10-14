package string_util

import (
	"fmt"
	"testing"
)

func TestString2Bytes(t *testing.T) {
	bytes := []byte{104, 101, 108, 108, 111}
	str := Bytes2String(bytes)
	fmt.Println(str)

	bytes[1] = 100
	fmt.Println(str)
	String2Bytes(str)[1] = 101
	fmt.Println(str)
}

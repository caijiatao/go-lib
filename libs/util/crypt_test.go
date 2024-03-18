package util

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	text, err := Encrypt("hello go")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(text)
}

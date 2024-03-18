package runtime

import (
	"fmt"
	"testing"
)

func TestGo(t *testing.T) {
	Go(func() {
		fmt.Println("Test Panic")
		panic("Test Panic")
	})
}

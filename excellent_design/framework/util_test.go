package framework

import (
	"fmt"
	"testing"
)

func Test_chunkSizeFor(t *testing.T) {
	fmt.Println(chunkSizeFor(100, 16))
}

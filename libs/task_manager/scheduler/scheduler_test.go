package scheduler

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	for i := 1; i < 100000; i *= 10 {
		numAllNodes := i
		fmt.Println(numAllNodes, numAllNodes*(50-numAllNodes/125)/100)
	}
}

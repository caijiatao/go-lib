package exception

import (
	"fmt"
	"testing"
)

func TestBar(t *testing.T) {
	err := bar()
	fmt.Printf("err: %+v\n", err)
}

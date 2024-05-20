package concurrency

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func makeShutdownFunc() ShutdownFunc {
	return func(ctx context.Context) error {
		time.Sleep(1 * time.Second)
		fmt.Println("quit success")
		return nil
	}
}

func TestQuit(t *testing.T) {
	RegisterShutdown(makeShutdownFunc())
	RegisterShutdown(makeShutdownFunc())
	for {
		time.Sleep(1 * time.Second)
		fmt.Println("running")
	}
}

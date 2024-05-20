package main

import (
	"context"
	"fmt"
	"golib/libs/concurrency"
	"runtime"
	"time"
)

func add(a, b int) int {
	return a + b
}

func deadLoop() {
	for {
		add(3, 5)
	}
}

func scheduler() {
	runtime.GOMAXPROCS(1)
	go deadLoop()
	for {
		time.Sleep(1 * time.Second)
		fmt.Println("main")
	}
}

func makeShutdownFunc() concurrency.ShutdownFunc {
	return func(ctx context.Context) error {
		time.Sleep(1 * time.Second)
		fmt.Println("quit process success")
		return nil
	}
}

func main() {
	concurrency.RegisterShutdown(makeShutdownFunc())
	concurrency.RegisterShutdown(makeShutdownFunc())

	for {
		time.Sleep(time.Second)
		fmt.Println("running")
	}
}

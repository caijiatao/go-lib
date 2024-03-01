package concurrency

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBoundedFrequencyRunner_Run(t *testing.T) {
	count := 0
	runner := NewBoundedFrequencyRunner(func() {
		fmt.Println("running process")
		time.Sleep(time.Second * 1)
		count++
		fmt.Println("running process end")
	})

	go runner.Loop()
	// 第一个正常启动
	go runner.Run()
	// 第二个不会启动，但是run信号正常写入
	go runner.Run()
	// 第三个直接丢弃信号不会启动
	go runner.Run()

	time.Sleep(time.Second * 2)

	assert.Equal(t, 1, count)
}

func TestTimer(t *testing.T) {
	timer := time.NewTimer(time.Second)
	go func() {
		<-timer.C
		fmt.Println("1")
		time.Sleep(time.Second * 1)
		timer.Reset(time.Second * 2)
		fmt.Println("reset")
		<-timer.C
		fmt.Println("2")
	}()

	time.Sleep(time.Second * 5)
}

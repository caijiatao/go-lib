package concurrency

import (
	"testing"
	"time"
)

func TestClick(t *testing.T) {
	rTime := newRealTimer(time.NewTimer(0), time.Now())
	go func() {
		for {
			select {
			case <-rTime.C():
				t.Log("go timer click")
				//default:
				//	t.Log("go timer not click")
			}
		}
	}()
	rTime.Reset(time.Second)
	time.Sleep(time.Second * 2)
	rTime.timer.Stop()
	time.Sleep(time.Second)
}

func TestStop(t *testing.T) {
	rTime := newRealTimer(time.NewTimer(0), time.Now())
	go func() {
		for {
			select {
			case <-rTime.C():
				t.Log("go timer click")
			}
		}
	}()
	time.Sleep(time.Second)
	rTime.timer.Stop()
	rTime.Reset(time.Second)
	time.Sleep(time.Second * 2)
}

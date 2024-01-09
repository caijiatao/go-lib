package concurrency

import "time"

type realTimer struct {
	timer *time.Timer
	next  time.Time
}

func (rt *realTimer) C() <-chan time.Time {
	return rt.timer.C
}

func (rt *realTimer) Reset(d time.Duration) bool {
	rt.next = time.Now().Add(d)
	return rt.timer.Reset(d)
}

package concurrency

import "time"

type realTimer struct {
	timer *time.Timer
	next  time.Time
}

func newRealTimer(timer *time.Timer, next time.Time) *realTimer {
	return &realTimer{timer: timer, next: next}
}

func (rt *realTimer) C() <-chan time.Time {
	return rt.timer.C
}

func (rt *realTimer) Stop() bool {
	return rt.timer.Stop()
}

func (rt *realTimer) Reset(d time.Duration) bool {
	rt.next = time.Now().Add(d)
	return rt.timer.Reset(d)
}

package concurrency

import "fmt"

type BoundedFrequencyRunner struct {
	run chan struct{}
	fn  func()
}

func NewBoundedFrequencyRunner(fn func()) *BoundedFrequencyRunner {
	return &BoundedFrequencyRunner{
		run: make(chan struct{}, 1),
		fn:  fn,
	}
}

func (b *BoundedFrequencyRunner) Run() {
	select {
	case b.run <- struct{}{}:
		fmt.Println("b is empty")
	default:
		fmt.Println("b is not empty")
	}
}

func (b *BoundedFrequencyRunner) Loop() {
	for {
		select {
		case <-b.run:
			b.fn()
		}
	}
}

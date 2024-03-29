package concurrency

import "sync"

type Parallelizer struct {
	Concurrency int
	ch          chan struct{}
}

func NewParallelizer(concurrency int) *Parallelizer {
	return &Parallelizer{
		Concurrency: concurrency,
		ch:          make(chan struct{}, concurrency),
	}
}

type DoWorkerPieceFunc func(piece int)

func (p *Parallelizer) Until(pices int, f DoWorkerPieceFunc) {
	wg := sync.WaitGroup{}
	for i := 0; i < pices; i++ {
		p.ch <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer func() {
				<-p.ch
				wg.Done()
			}()
			f(i)
		}(i)
	}
	wg.Wait()
}

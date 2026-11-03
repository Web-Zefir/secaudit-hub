package worker

import "sync"

// Pool limits how many scans run at once.
type Pool struct {
	sem chan struct{}
	wg  sync.WaitGroup
}

func NewPool(size int) *Pool {
	if size < 1 {
		size = 1
	}
	return &Pool{sem: make(chan struct{}, size)}
}

func (p *Pool) Go(fn func()) {
	p.wg.Add(1)
	p.sem <- struct{}{}
	go func() {
		defer p.wg.Done()
		defer func() { <-p.sem }()
		defer func() { _ = recover() }()
		fn()
	}()
}

package throttle

import (
	"sync"
	"sync/atomic"
)

type Tick struct {
	n atomic.Int64
	c atomic.Int64
	m sync.Mutex
}

// Go calls fn every interval of d.
// Caller is waiting for completion of the fn.
// On the first call, all the callers are waits until the fn is completed
func (it *Tick) Do(d int, fn func()) {
	var last = it.n.Load()
	var next = it.c.Add(1)
	if next < last+int64(d) {
		return
	}
	switch {
	case last == 0:
		it.m.Lock()
	case !it.m.TryLock():
		return
	}
	if it.n.Load() != last {
		it.m.Unlock()
		return
	}
	fn()
	it.n.Store(next)
	it.m.Unlock()
}

// Go calls the fn function no more than once every n calls.
// Caller isn't waiting for completion of the fn.
// On the first call, all the callers are waits until the fn is completed
func (it *Tick) Go(d int, fn func()) {
	var last = it.n.Load()
	var next = it.c.Add(1)
	if next < last+int64(d) {
		return
	}
	switch {
	case last == 0:
		it.m.Lock()
	case !it.m.TryLock():
		return
	}
	if it.n.Load() != last {
		it.m.Unlock()
		return
	}
	if last == 0 {
		fn()
		it.n.Store(next)
		it.m.Unlock()
		return
	}
	go func() {
		fn()
		it.n.Store(next)
		it.m.Unlock()
	}()
}

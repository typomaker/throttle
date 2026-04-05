package throttle

import (
	"sync"
	"sync/atomic"
	"time"
)

type Time struct {
	n atomic.Int64
	m sync.Mutex
}

// Go calls fn every interval of d.
// Caller is waiting for completion of the fn.
// On the first call, all the callers are waits until the fn is completed
func (it *Time) Do(d time.Duration, fn func()) {
	var last = it.n.Load()
	var next = time.Now().UnixNano()
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

// Go calls fn every interval of d.
// Caller isn't waiting for completion of the fn.
// On the first call, all the callers are waits until the fn is completed
func (it *Time) Go(d time.Duration, fn func()) {
	var last = it.n.Load()
	var next = time.Now().UnixNano()
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
	go func() {
		fn()
		it.n.Store(next)
		it.m.Unlock()
	}()
}

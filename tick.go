package throttle

import (
	"sync"
	"sync/atomic"
)

// Tick throttles handler execution. Its zero value is ready to use.
// A Tick must not be copied after first use.
// Callers and their handlers are responsible for synchronizing user data.
type Tick struct {
	n atomic.Int64
	c atomic.Int64
	m sync.Mutex
}

// Do invokes fn on the first call, then at most once per d calls
// counted from the last successful run, including calls skipped during execution.
// Concurrent calls skip an active run, including during initialization.
// The caller waits for its own run to complete.
func (it *Tick) Do(d int, fn func()) {
	last := it.n.Load()
	next := it.c.Add(1)
	if last != 0 && next-last < int64(d) {
		return
	}
	if !it.m.TryLock() {
		return
	}
	if it.n.Load() != last {
		it.m.Unlock()
		return
	}
	defer it.m.Unlock()
	fn()
	it.n.Store(next)
}

// Go invokes fn on the first call, then at most once per d calls
// counted from the last successful run, including calls skipped during execution.
// Concurrent calls skip an active run, including during initialization.
// Every run executes in a goroutine.
func (it *Tick) Go(d int, fn func()) {
	last := it.n.Load()
	next := it.c.Add(1)
	if last != 0 && next-last < int64(d) {
		return
	}
	if !it.m.TryLock() {
		return
	}
	if it.n.Load() != last {
		it.m.Unlock()
		return
	}
	go func() {
		defer it.m.Unlock()
		fn()
		it.n.Store(next)
	}()
}

// DoWait invokes fn on the first call, then at most once per d calls
// counted from the last successful run, including calls skipped during execution.
// Until the first successful run, callers wait and fn runs synchronously.
// A panic leaves initialization pending for the next caller.
// After initialization, concurrent calls skip an active run.
// Subsequent runs execute synchronously.
// During initialization, fn must not recursively call DoWait or GoWait
// on the same instance: doing so deadlocks.
func (it *Tick) DoWait(d int, fn func()) {
	if it.n.Load() != 0 {
		it.Do(d, fn)
		return
	}
	next := it.c.Add(1)
	it.m.Lock()
	defer it.m.Unlock()
	if it.n.Load() != 0 {
		return
	}
	fn()
	it.n.Store(next)
}

// GoWait invokes fn on the first call, then at most once per d calls
// counted from the last successful run, including calls skipped during execution.
// Until the first successful run, callers wait and fn runs synchronously.
// A panic leaves initialization pending for the next caller.
// After initialization, concurrent calls skip an active run.
// Subsequent runs execute in a goroutine.
// During initialization, fn must not recursively call DoWait or GoWait
// on the same instance: doing so deadlocks.
func (it *Tick) GoWait(d int, fn func()) {
	if it.n.Load() != 0 {
		it.Go(d, fn)
		return
	}
	next := it.c.Add(1)
	it.m.Lock()
	defer it.m.Unlock()
	if it.n.Load() != 0 {
		return
	}
	fn()
	it.n.Store(next)
}

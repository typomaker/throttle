package throttle

import (
	"sync"
	"sync/atomic"
	"time"
)

// timeOrigin retains the monotonic clock reading for elapsed-time comparisons.
var timeOrigin = time.Now()

// Time throttles handler execution. Its zero value is ready to use.
// A Time must not be copied after first use.
// Callers and their handlers are responsible for synchronizing user data.
type Time struct {
	n atomic.Int64
	m sync.Mutex
}

// Do invokes fn at most once per interval d.
// Concurrent calls skip an active run, including during initialization.
// The caller waits for its own run to complete.
func (it *Time) Do(d time.Duration, fn func()) {
	last := it.n.Load()
	now := time.Since(timeOrigin).Nanoseconds() + 1
	if last != 0 && now-last < int64(d) {
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
	// Reserve zero for instances with no successful runs.
	next := time.Since(timeOrigin).Nanoseconds() + 1
	fn()
	it.n.Store(next)
}

// Go invokes fn at most once per interval d.
// Concurrent calls skip an active run, including during initialization.
// Every run executes in a goroutine.
func (it *Time) Go(d time.Duration, fn func()) {
	last := it.n.Load()
	now := time.Since(timeOrigin).Nanoseconds() + 1
	if last != 0 && now-last < int64(d) {
		return
	}
	if !it.m.TryLock() {
		return
	}
	if it.n.Load() != last {
		it.m.Unlock()
		return
	}
	// Reserve zero for instances with no successful runs.
	next := time.Since(timeOrigin).Nanoseconds() + 1
	go func() {
		defer it.m.Unlock()
		fn()
		it.n.Store(next)
	}()
}

// DoWait invokes fn at most once per interval d.
// Until the first successful run, callers wait and fn runs synchronously.
// A panic leaves initialization pending for the next caller.
// After initialization, concurrent calls skip an active run.
// Subsequent runs execute synchronously.
// During initialization, fn must not recursively call DoWait or GoWait
// on the same instance: doing so deadlocks.
func (it *Time) DoWait(d time.Duration, fn func()) {
	if it.n.Load() != 0 {
		it.Do(d, fn)
		return
	}
	it.m.Lock()
	defer it.m.Unlock()
	if it.n.Load() != 0 {
		return
	}
	// Reserve zero for instances with no successful runs.
	next := time.Since(timeOrigin).Nanoseconds() + 1
	fn()
	it.n.Store(next)
}

// GoWait invokes fn at most once per interval d.
// Until the first successful run, callers wait and fn runs synchronously.
// A panic leaves initialization pending for the next caller.
// After initialization, concurrent calls skip an active run.
// Subsequent runs execute in a goroutine.
// During initialization, fn must not recursively call DoWait or GoWait
// on the same instance: doing so deadlocks.
func (it *Time) GoWait(d time.Duration, fn func()) {
	if it.n.Load() != 0 {
		it.Go(d, fn)
		return
	}
	it.m.Lock()
	defer it.m.Unlock()
	if it.n.Load() != 0 {
		return
	}
	// Reserve zero for instances with no successful runs.
	next := time.Since(timeOrigin).Nanoseconds() + 1
	fn()
	it.n.Store(next)
}

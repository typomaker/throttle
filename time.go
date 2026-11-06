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
// It returns true if fn completes normally, or false if the call is skipped.
func (it *Time) Do(d time.Duration, fn func()) bool {
	last := it.n.Load()
	now := time.Since(timeOrigin).Nanoseconds() + 1
	if last != 0 && now-last < int64(d) {
		return false
	}
	if !it.m.TryLock() {
		return false
	}
	if it.n.Load() != last {
		it.m.Unlock()
		return false
	}
	defer it.m.Unlock()
	// Reserve zero for instances with no successful runs.
	next := time.Since(timeOrigin).Nanoseconds() + 1
	fn()
	it.n.Store(next)
	return true
}

// Go invokes fn at most once per interval d.
// Concurrent calls skip an active run, including during initialization.
// Every run executes in a goroutine.
// It returns true if fn is scheduled, or false if the call is skipped.
// A true result does not confirm completion.
func (it *Time) Go(d time.Duration, fn func()) bool {
	last := it.n.Load()
	now := time.Since(timeOrigin).Nanoseconds() + 1
	if last != 0 && now-last < int64(d) {
		return false
	}
	if !it.m.TryLock() {
		return false
	}
	if it.n.Load() != last {
		it.m.Unlock()
		return false
	}
	// Reserve zero for instances with no successful runs.
	next := time.Since(timeOrigin).Nanoseconds() + 1
	go func() {
		defer it.m.Unlock()
		fn()
		it.n.Store(next)
	}()
	return true
}

// DoWait invokes fn at most once per interval d.
// Until the first successful run, callers wait and fn runs synchronously.
// A panic leaves initialization pending for the next caller.
// After initialization, concurrent calls skip an active run.
// Subsequent runs execute synchronously.
// During initialization, fn must not recursively call DoWait or GoWait
// on the same instance: doing so deadlocks.
// It returns true after running or waiting for initialization; afterwards it
// returns the result of Do.
func (it *Time) DoWait(d time.Duration, fn func()) bool {
	if it.n.Load() != 0 {
		return it.Do(d, fn)
	}
	it.m.Lock()
	defer it.m.Unlock()
	if it.n.Load() != 0 {
		return true
	}
	// Reserve zero for instances with no successful runs.
	next := time.Since(timeOrigin).Nanoseconds() + 1
	fn()
	it.n.Store(next)
	return true
}

// GoWait invokes fn at most once per interval d.
// Until the first successful run, callers wait and fn runs synchronously.
// A panic leaves initialization pending for the next caller.
// After initialization, concurrent calls skip an active run.
// Subsequent runs execute in a goroutine.
// During initialization, fn must not recursively call DoWait or GoWait
// on the same instance: doing so deadlocks.
// It returns true after running or waiting for initialization; afterwards it
// returns the result of Go.
func (it *Time) GoWait(d time.Duration, fn func()) bool {
	if it.n.Load() != 0 {
		return it.Go(d, fn)
	}
	it.m.Lock()
	defer it.m.Unlock()
	if it.n.Load() != 0 {
		return true
	}
	// Reserve zero for instances with no successful runs.
	next := time.Since(timeOrigin).Nanoseconds() + 1
	fn()
	it.n.Store(next)
	return true
}

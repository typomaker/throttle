package throttle

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestTimeDoZeroInterval(t *testing.T) {
	var v Time
	calls := 0
	for range 3 {
		v.Do(0, func() { calls++ })
		v.m.Lock()
		v.m.Unlock()
	}
	if calls != 3 {
		t.Fatalf("got %d calls, want 3", calls)
	}
}
func TestTimeDoNegativeInterval(t *testing.T) {
	var v Time
	calls := 0
	for range 3 {
		v.Do(-1, func() { calls++ })
		v.m.Lock()
		v.m.Unlock()
	}
	if calls != 3 {
		t.Fatalf("got %d calls, want 3", calls)
	}
}
func TestTimeGoZeroInterval(t *testing.T) {
	var v Time
	calls := 0
	for range 3 {
		v.Go(0, func() { calls++ })
		v.m.Lock()
		v.m.Unlock()
	}
	if calls != 3 {
		t.Fatalf("got %d calls, want 3", calls)
	}
}
func TestTimeGoNegativeInterval(t *testing.T) {
	var v Time
	calls := 0
	for range 3 {
		v.Go(-1, func() { calls++ })
		v.m.Lock()
		v.m.Unlock()
	}
	if calls != 3 {
		t.Fatalf("got %d calls, want 3", calls)
	}
}
func TestTimeDoWaitZeroInterval(t *testing.T) {
	var v Time
	calls := 0
	for range 3 {
		v.DoWait(0, func() { calls++ })
		v.m.Lock()
		v.m.Unlock()
	}
	if calls != 3 {
		t.Fatalf("got %d calls, want 3", calls)
	}
}
func TestTimeDoWaitNegativeInterval(t *testing.T) {
	var v Time
	calls := 0
	for range 3 {
		v.DoWait(-1, func() { calls++ })
		v.m.Lock()
		v.m.Unlock()
	}
	if calls != 3 {
		t.Fatalf("got %d calls, want 3", calls)
	}
}
func TestTimeGoWaitZeroInterval(t *testing.T) {
	var v Time
	calls := 0
	for range 3 {
		v.GoWait(0, func() { calls++ })
		v.m.Lock()
		v.m.Unlock()
	}
	if calls != 3 {
		t.Fatalf("got %d calls, want 3", calls)
	}
}
func TestTimeGoWaitNegativeInterval(t *testing.T) {
	var v Time
	calls := 0
	for range 3 {
		v.GoWait(-1, func() { calls++ })
		v.m.Lock()
		v.m.Unlock()
	}
	if calls != 3 {
		t.Fatalf("got %d calls, want 3", calls)
	}
}
func TestTimeDoConcurrentStateChanges(t *testing.T) {
	var v Time
	var active, calls atomic.Int64
	var callers sync.WaitGroup
	for range 32 {
		callers.Add(1)
		go func() {
			defer callers.Done()
			for range 4000 {
				v.Do(0, func() {
					if active.Add(1) != 1 {
						t.Error("handlers overlap")
					}
					calls.Add(1)
					active.Add(-1)
				})
			}
		}()
	}
	callers.Wait()
	v.m.Lock()
	defer v.m.Unlock()
	if calls.Load() == 0 {
		t.Fatal("no handler executed")
	}
	if active.Load() != 0 {
		t.Fatal("handler still active")
	}
}
func TestTimeGoConcurrentStateChanges(t *testing.T) {
	var v Time
	var active, calls atomic.Int64
	var callers sync.WaitGroup
	for range 32 {
		callers.Add(1)
		go func() {
			defer callers.Done()
			for range 4000 {
				v.Go(0, func() {
					if active.Add(1) != 1 {
						t.Error("handlers overlap")
					}
					calls.Add(1)
					active.Add(-1)
				})
			}
		}()
	}
	callers.Wait()
	v.m.Lock()
	defer v.m.Unlock()
	if calls.Load() == 0 {
		t.Fatal("no handler executed")
	}
	if active.Load() != 0 {
		t.Fatal("handler still active")
	}
}
func TestTimeGoSkipDoesNotScheduleHandler(t *testing.T) {
	var v Time
	v.DoWait(0, func() {})
	calls := 0
	v.Go(time.Hour, func() { calls++ })
	v.m.Lock()
	defer v.m.Unlock()
	if calls != 0 {
		t.Fatalf("unexpected calls: %d", calls)
	}
}
func TestTimeGoWaitSkipDoesNotScheduleHandler(t *testing.T) {
	var v Time
	v.DoWait(0, func() {})
	calls := 0
	v.GoWait(time.Hour, func() { calls++ })
	v.m.Lock()
	defer v.m.Unlock()
	if calls != 0 {
		t.Fatalf("unexpected calls: %d", calls)
	}
}
func TestTickDoZeroInterval(t *testing.T) {
	var v Tick
	calls := 0
	for range 3 {
		v.Do(0, func() { calls++ })
		v.m.Lock()
		v.m.Unlock()
	}
	if calls != 3 {
		t.Fatalf("got %d calls, want 3", calls)
	}
}
func TestTickDoNegativeInterval(t *testing.T) {
	var v Tick
	calls := 0
	for range 3 {
		v.Do(-1, func() { calls++ })
		v.m.Lock()
		v.m.Unlock()
	}
	if calls != 3 {
		t.Fatalf("got %d calls, want 3", calls)
	}
}
func TestTickGoZeroInterval(t *testing.T) {
	var v Tick
	calls := 0
	for range 3 {
		v.Go(0, func() { calls++ })
		v.m.Lock()
		v.m.Unlock()
	}
	if calls != 3 {
		t.Fatalf("got %d calls, want 3", calls)
	}
}
func TestTickGoNegativeInterval(t *testing.T) {
	var v Tick
	calls := 0
	for range 3 {
		v.Go(-1, func() { calls++ })
		v.m.Lock()
		v.m.Unlock()
	}
	if calls != 3 {
		t.Fatalf("got %d calls, want 3", calls)
	}
}
func TestTickDoWaitZeroInterval(t *testing.T) {
	var v Tick
	calls := 0
	for range 3 {
		v.DoWait(0, func() { calls++ })
		v.m.Lock()
		v.m.Unlock()
	}
	if calls != 3 {
		t.Fatalf("got %d calls, want 3", calls)
	}
}
func TestTickDoWaitNegativeInterval(t *testing.T) {
	var v Tick
	calls := 0
	for range 3 {
		v.DoWait(-1, func() { calls++ })
		v.m.Lock()
		v.m.Unlock()
	}
	if calls != 3 {
		t.Fatalf("got %d calls, want 3", calls)
	}
}
func TestTickGoWaitZeroInterval(t *testing.T) {
	var v Tick
	calls := 0
	for range 3 {
		v.GoWait(0, func() { calls++ })
		v.m.Lock()
		v.m.Unlock()
	}
	if calls != 3 {
		t.Fatalf("got %d calls, want 3", calls)
	}
}
func TestTickGoWaitNegativeInterval(t *testing.T) {
	var v Tick
	calls := 0
	for range 3 {
		v.GoWait(-1, func() { calls++ })
		v.m.Lock()
		v.m.Unlock()
	}
	if calls != 3 {
		t.Fatalf("got %d calls, want 3", calls)
	}
}
func TestTickDoConcurrentStateChanges(t *testing.T) {
	var v Tick
	var active, calls atomic.Int64
	var callers sync.WaitGroup
	for range 32 {
		callers.Add(1)
		go func() {
			defer callers.Done()
			for range 4000 {
				v.Do(0, func() {
					if active.Add(1) != 1 {
						t.Error("handlers overlap")
					}
					calls.Add(1)
					active.Add(-1)
				})
			}
		}()
	}
	callers.Wait()
	v.m.Lock()
	defer v.m.Unlock()
	if calls.Load() == 0 {
		t.Fatal("no handler executed")
	}
	if active.Load() != 0 {
		t.Fatal("handler still active")
	}
}
func TestTickGoConcurrentStateChanges(t *testing.T) {
	var v Tick
	var active, calls atomic.Int64
	var callers sync.WaitGroup
	for range 32 {
		callers.Add(1)
		go func() {
			defer callers.Done()
			for range 4000 {
				v.Go(0, func() {
					if active.Add(1) != 1 {
						t.Error("handlers overlap")
					}
					calls.Add(1)
					active.Add(-1)
				})
			}
		}()
	}
	callers.Wait()
	v.m.Lock()
	defer v.m.Unlock()
	if calls.Load() == 0 {
		t.Fatal("no handler executed")
	}
	if active.Load() != 0 {
		t.Fatal("handler still active")
	}
}
func TestTickGoSkipDoesNotScheduleHandler(t *testing.T) {
	var v Tick
	v.DoWait(0, func() {})
	calls := 0
	v.Go(100, func() { calls++ })
	v.m.Lock()
	defer v.m.Unlock()
	if calls != 0 {
		t.Fatalf("unexpected calls: %d", calls)
	}
}
func TestTickGoWaitSkipDoesNotScheduleHandler(t *testing.T) {
	var v Tick
	v.DoWait(0, func() {})
	calls := 0
	v.GoWait(100, func() { calls++ })
	v.m.Lock()
	defer v.m.Unlock()
	if calls != 0 {
		t.Fatalf("unexpected calls: %d", calls)
	}
}

func TestTimeDoWaitRecordsTimeAfterLock(t *testing.T) {
	var v Time
	v.m.Lock()
	started, done := make(chan struct{}), make(chan struct{})
	go func() {
		close(started)
		v.DoWait(time.Second, func() {})
		close(done)
	}()
	<-started
	// Let the caller reach the locked initialization path.
	time.Sleep(20 * time.Millisecond)
	unlockedAt := time.Since(timeOrigin).Nanoseconds() + 1
	v.m.Unlock()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("initialization did not complete")
	}
	if got := v.n.Load(); got < unlockedAt {
		t.Fatalf("recorded time %d precedes lock release %d", got, unlockedAt)
	}
}

func TestTimeGoWaitRecordsTimeAfterLock(t *testing.T) {
	var v Time
	v.m.Lock()
	started, done := make(chan struct{}), make(chan struct{})
	go func() {
		close(started)
		v.GoWait(time.Second, func() {})
		close(done)
	}()
	<-started
	// Let the caller reach the locked initialization path.
	time.Sleep(20 * time.Millisecond)
	unlockedAt := time.Since(timeOrigin).Nanoseconds() + 1
	v.m.Unlock()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("initialization did not complete")
	}
	if got := v.n.Load(); got < unlockedAt {
		t.Fatalf("recorded time %d precedes lock release %d", got, unlockedAt)
	}
}

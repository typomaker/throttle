package throttle_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/typomaker/throttle"
)

func TestTimeDo(t *testing.T) {
	t.Parallel()
	o := atomic.Int64{}
	e := throttle.Time{}
	c := make(chan bool)

	const n = 100
	for range n {
		go func() {
			e.Do(2, func() {
				o.Add(1)
			})
			c <- true
		}()
	}
	for range n {
		<-c
	}
	require.GreaterOrEqual(t, o.Load(), int64(1))
	require.LessOrEqual(t, o.Load(), int64(n))
}

func BenchmarkTimeDo(b *testing.B) {
	o := atomic.Int64{}
	e := throttle.Time{}
	f := func() { o.Add(1) }
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			e.Do(10, f)
		}
	})
}

func TestTimeGo(t *testing.T) {
	t.Parallel()
	o := atomic.Int64{}
	e := throttle.Time{}
	c := make(chan bool)

	const n = 10
	done := make(chan struct{}, n)
	for range n {
		go func() {
			e.Go(2, func() {
				o.Add(1)
				done <- struct{}{}
			})
			c <- true
		}()
	}
	for range n {
		<-c
	}
	await(t, done)
	require.GreaterOrEqual(t, o.Load(), int64(1))
	require.LessOrEqual(t, o.Load(), int64(n))
}

func BenchmarkTimeGo(b *testing.B) {
	o := atomic.Int64{}
	f := func() { o.Add(1) }
	e := throttle.Time{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			e.Go(10, f)
		}
	})
}
func TestTimeGoWaitFirstSync(t *testing.T) {
	t.Parallel()
	o := atomic.Int64{}
	e := throttle.Time{}
	e.GoWait(0, func() {
		o.Add(1)
	})
	require.EqualValues(t, 1, o.Load())
}

func TestTimeDoFirstRunIgnoresInterval(t *testing.T) {
	var v throttle.Time
	called := false
	v.Do(24*time.Hour, func() { called = true })
	require.True(t, called)
}

func TestTimeDoRespectsIntervalAfterFirstRun(t *testing.T) {
	var v throttle.Time
	calls := 0
	v.Do(24*time.Hour, func() { calls++ })
	v.Do(24*time.Hour, func() { calls++ })
	require.Equal(t, 1, calls)
}

func TestTimeGoWaitFirstRunIgnoresInterval(t *testing.T) {
	var v throttle.Time
	called := false
	v.GoWait(24*time.Hour, func() { called = true })
	require.True(t, called)
}

func TestTimeGoWaitRespectsIntervalAfterFirstRun(t *testing.T) {
	var v throttle.Time
	calls := 0
	v.GoWait(24*time.Hour, func() { calls++ })
	v.GoWait(24*time.Hour, func() { calls++ })
	require.Equal(t, 1, calls)
}

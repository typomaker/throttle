package throttle_test

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typomaker/throttle"
)

func TestTickDo(t *testing.T) {
	t.Parallel()
	o := atomic.Int64{}
	e := throttle.Tick{}
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

func BenchmarkTickDo(b *testing.B) {
	o := atomic.Int64{}
	e := throttle.Tick{}
	f := func() { o.Add(1) }
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			e.Do(10, f)
		}
	})
}

func TestTickGo(t *testing.T) {
	t.Parallel()
	o := atomic.Int64{}
	e := throttle.Tick{}
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

func BenchmarkTickGo(b *testing.B) {
	o := atomic.Int64{}
	f := func() { o.Add(1) }
	e := throttle.Tick{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			e.Go(10, f)
		}
	})
}
func TestTickGoWaitFirstSync(t *testing.T) {
	t.Parallel()
	o := atomic.Int64{}
	e := throttle.Tick{}
	e.GoWait(0, func() {
		o.Add(1)
	})
	require.EqualValues(t, 1, o.Load())
}

func TestTickDoFirstRunIgnoresInterval(t *testing.T) {
	var v throttle.Tick
	called := false
	v.Do(10, func() { called = true })
	require.True(t, called)
}

func TestTickGoFirstRunIgnoresInterval(t *testing.T) {
	var v throttle.Tick
	done := make(chan struct{})
	v.Go(10, func() { close(done) })
	await(t, done)
}

func TestTickDoCountsTenCallsFromLastRun(t *testing.T) {
	var v throttle.Tick
	calls := 0
	handler := func() { calls++ }
	v.Do(10, handler)
	require.Equal(t, 1, calls)
	for range 9 {
		v.Do(10, handler)
	}
	require.Equal(t, 1, calls)
	v.Do(10, handler)
	require.Equal(t, 2, calls)
}

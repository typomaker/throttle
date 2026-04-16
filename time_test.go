package throttle_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/typomaker/throttle"
)

type one int

func (o *one) Increment() {
	*o++
}

func TestDo(t *testing.T) {
	t.Parallel()
	o := one(0)
	e := throttle.Time{}
	c := make(chan bool)

	const n = 100
	for range n {
		go func() {
			e.Do(time.Minute, func() {
				o.Increment()
			})
			c <- true
		}()
	}
	for range n {
		<-c
	}
	require.EqualValues(t, 1, o)
}
func BenchmarkDo(b *testing.B) {
	e := throttle.Time{}
	f := func() {}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			e.Do(time.Hour, f)
		}
	})
}
func TestGo(t *testing.T) {
	t.Parallel()
	o := one(0)
	e := throttle.Time{}
	const n = 100
	for range n {
		e.Go(time.Minute, func() {
			o.Increment()
		})
	}
	time.Sleep(time.Millisecond)
	require.EqualValues(t, 1, o)
}
func BenchmarkGo(b *testing.B) {
	e := throttle.Time{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			e.Go(time.Hour, func() {})
		}
	})
}
func TestDoZeroD(t *testing.T) {
	t.Parallel()
	o := one(0)
	e := throttle.Time{}
	c := make(chan bool)

	const n = 100
	for range n {
		go func() {
			e.Do(0, func() {
				o.Increment()
			})
			c <- true
		}()
	}
	for range n {
		<-c
	}
	require.LessOrEqual(t, 1, o)
	require.GreaterOrEqual(t, n, o)
}
func TestGoZeroD(t *testing.T) {
	t.Parallel()
	o := one(0)
	e := throttle.Time{}
	c := make(chan bool)

	const n = 100
	for range n {
		go func() {
			e.Go(0, func() {
				o.Increment()
			})
			c <- true
		}()
	}
	for range n {
		<-c
	}
	require.LessOrEqual(t, 1, o)
	require.GreaterOrEqual(t, n, o)
}
func TestGoFirstSync(t *testing.T) {
	t.Parallel()
	o := one(0)
	e := throttle.Time{}
	e.Go(0, func() {
		o.Increment()
	})
	require.EqualValues(t, 1, o)
}

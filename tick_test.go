package throttle_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typomaker/throttle"
)

func TestTickDo(t *testing.T) {
	t.Parallel()
	o := one(0)
	e := throttle.Tick{}
	c := make(chan bool)

	const n = 100
	for range n {
		go func() {
			e.Do(2, func() {
				o.Increment()
			})
			c <- true
		}()
	}
	for range n {
		<-c
	}
	require.LessOrEqual(t, 1, o)
}

func BenchmarkTickDo(b *testing.B) {
	o := one(0)
	e := throttle.Tick{}
	f := func() { o.Increment() }
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			e.Do(10, f)
		}
	})
}

func TestTickGo(t *testing.T) {
	t.Parallel()
	o := one(0)
	e := throttle.Tick{}
	c := make(chan bool)

	const n = 10
	for range n {
		go func() {
			e.Go(2, func() {
				o.Increment()
			})
			c <- true
		}()
	}
	for range n {
		<-c
	}
	require.LessOrEqual(t, 1, o)
}

func BenchmarkTickGo(b *testing.B) {
	o := one(0)
	f := func() { o.Increment() }
	e := throttle.Tick{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			e.Go(10, f)
		}
	})
}
func TestTickGoFirstSync(t *testing.T) {
	t.Parallel()
	o := one(0)
	e := throttle.Tick{}
	e.Go(0, func() {
		o.Increment()
	})
	require.EqualValues(t, 1, o)
}

package throttle_test

import (
	"github.com/stretchr/testify/require"
	"github.com/typomaker/throttle"
	"testing"
	"time"
)

func TestTimeDoReturnsTrueOnFirstRun(t *testing.T) {
	var v throttle.Time
	done := make(chan struct{})
	require.True(t, v.Do(time.Hour, func() { close(done) }))
	await(t, done)
}
func TestTimeDoReturnsFalseBeforeInterval(t *testing.T) {
	var v throttle.Time
	require.True(t, v.DoWait(0, func() {}))
	require.False(t, v.Do(time.Hour, func() { t.Error("unexpected handler") }))
}
func TestTimeGoReturnsTrueOnFirstRun(t *testing.T) {
	var v throttle.Time
	done := make(chan struct{})
	require.True(t, v.Go(time.Hour, func() { close(done) }))
	await(t, done)
}
func TestTimeGoReturnsFalseBeforeInterval(t *testing.T) {
	var v throttle.Time
	require.True(t, v.DoWait(0, func() {}))
	require.False(t, v.Go(time.Hour, func() { t.Error("unexpected handler") }))
}
func TestTimeDoWaitReturnsTrueOnFirstRun(t *testing.T) {
	var v throttle.Time
	done := make(chan struct{})
	require.True(t, v.DoWait(time.Hour, func() { close(done) }))
	await(t, done)
}
func TestTimeDoWaitReturnsFalseBeforeInterval(t *testing.T) {
	var v throttle.Time
	require.True(t, v.DoWait(0, func() {}))
	require.False(t, v.DoWait(time.Hour, func() { t.Error("unexpected handler") }))
}
func TestTimeGoWaitReturnsTrueOnFirstRun(t *testing.T) {
	var v throttle.Time
	done := make(chan struct{})
	require.True(t, v.GoWait(time.Hour, func() { close(done) }))
	await(t, done)
}
func TestTimeGoWaitReturnsFalseBeforeInterval(t *testing.T) {
	var v throttle.Time
	require.True(t, v.DoWait(0, func() {}))
	require.False(t, v.GoWait(time.Hour, func() { t.Error("unexpected handler") }))
}
func TestTickDoReturnsTrueOnFirstRun(t *testing.T) {
	var v throttle.Tick
	done := make(chan struct{})
	require.True(t, v.Do(10, func() { close(done) }))
	await(t, done)
}
func TestTickDoReturnsFalseBeforeInterval(t *testing.T) {
	var v throttle.Tick
	require.True(t, v.DoWait(0, func() {}))
	require.False(t, v.Do(10, func() { t.Error("unexpected handler") }))
}
func TestTickGoReturnsTrueOnFirstRun(t *testing.T) {
	var v throttle.Tick
	done := make(chan struct{})
	require.True(t, v.Go(10, func() { close(done) }))
	await(t, done)
}
func TestTickGoReturnsFalseBeforeInterval(t *testing.T) {
	var v throttle.Tick
	require.True(t, v.DoWait(0, func() {}))
	require.False(t, v.Go(10, func() { t.Error("unexpected handler") }))
}
func TestTickDoWaitReturnsTrueOnFirstRun(t *testing.T) {
	var v throttle.Tick
	done := make(chan struct{})
	require.True(t, v.DoWait(10, func() { close(done) }))
	await(t, done)
}
func TestTickDoWaitReturnsFalseBeforeInterval(t *testing.T) {
	var v throttle.Tick
	require.True(t, v.DoWait(0, func() {}))
	require.False(t, v.DoWait(10, func() { t.Error("unexpected handler") }))
}
func TestTickGoWaitReturnsTrueOnFirstRun(t *testing.T) {
	var v throttle.Tick
	done := make(chan struct{})
	require.True(t, v.GoWait(10, func() { close(done) }))
	await(t, done)
}
func TestTickGoWaitReturnsFalseBeforeInterval(t *testing.T) {
	var v throttle.Tick
	require.True(t, v.DoWait(0, func() {}))
	require.False(t, v.GoWait(10, func() { t.Error("unexpected handler") }))
}

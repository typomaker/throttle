package throttle_test

import (
	"github.com/stretchr/testify/require"
	"github.com/typomaker/throttle"
	"testing"
	"time"
)

func TestTimeDoAfterInitialization(t *testing.T) {
	var v throttle.Time
	v.DoWait(0, func() {})
	entered, release, returned, finished := make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{})
	go func() {
		v.Do(0, func() { close(entered); <-release; close(finished) })
		close(returned)
	}()
	await(t, entered)
	select {
	case <-returned:
		t.Error("synchronous call returned before handler")
	default:
	}
	require.False(t, v.Do(0, func() { t.Error("active handler must be skipped") }))
	close(release)
	await(t, returned)
	await(t, finished)
}
func TestTimeDoRunsAfterInterval(t *testing.T) {
	var v throttle.Time
	v.DoWait(0, func() {})
	time.Sleep(time.Millisecond)
	done := make(chan struct{})
	v.Do(time.Millisecond, func() { close(done) })
	await(t, done)
}
func TestTimeDoSkipsBeforeInterval(t *testing.T) {
	var v throttle.Time
	v.DoWait(0, func() {})
	v.Do(time.Hour, func() { t.Error("interval has not elapsed") })
}
func TestTimeDoPanicAfterInitialization(t *testing.T) {
	var v throttle.Time
	v.DoWait(0, func() {})
	require.Panics(t, func() { v.Do(0, func() { panic("failed") }) })
	done := make(chan struct{})
	go v.Do(0, func() { close(done) })
	await(t, done)
}
func TestTimeGoAfterInitialization(t *testing.T) {
	var v throttle.Time
	v.DoWait(0, func() {})
	entered, release, returned, finished := make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{})
	go func() {
		v.Go(0, func() { close(entered); <-release; close(finished) })
		close(returned)
	}()
	await(t, entered)
	await(t, returned)
	require.False(t, v.Go(0, func() { t.Error("active handler must be skipped") }))
	close(release)
	await(t, returned)
	await(t, finished)
}
func TestTimeGoRunsAfterInterval(t *testing.T) {
	var v throttle.Time
	v.DoWait(0, func() {})
	time.Sleep(time.Millisecond)
	done := make(chan struct{})
	v.Go(time.Millisecond, func() { close(done) })
	await(t, done)
}
func TestTimeGoSkipsBeforeInterval(t *testing.T) {
	var v throttle.Time
	v.DoWait(0, func() {})
	v.Go(time.Hour, func() { t.Error("interval has not elapsed") })
}
func TestTimeDoWaitAfterInitialization(t *testing.T) {
	var v throttle.Time
	v.DoWait(0, func() {})
	entered, release, returned, finished := make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{})
	go func() {
		v.DoWait(0, func() { close(entered); <-release; close(finished) })
		close(returned)
	}()
	await(t, entered)
	select {
	case <-returned:
		t.Error("synchronous call returned before handler")
	default:
	}
	require.False(t, v.DoWait(0, func() { t.Error("active handler must be skipped") }))
	close(release)
	await(t, returned)
	await(t, finished)
}
func TestTimeDoWaitRunsAfterInterval(t *testing.T) {
	var v throttle.Time
	v.DoWait(0, func() {})
	time.Sleep(time.Millisecond)
	done := make(chan struct{})
	v.DoWait(time.Millisecond, func() { close(done) })
	await(t, done)
}
func TestTimeDoWaitSkipsBeforeInterval(t *testing.T) {
	var v throttle.Time
	v.DoWait(0, func() {})
	v.DoWait(time.Hour, func() { t.Error("interval has not elapsed") })
}
func TestTimeDoWaitPanicAfterInitialization(t *testing.T) {
	var v throttle.Time
	v.DoWait(0, func() {})
	require.Panics(t, func() { v.DoWait(0, func() { panic("failed") }) })
	done := make(chan struct{})
	go v.DoWait(0, func() { close(done) })
	await(t, done)
}
func TestTimeGoWaitAfterInitialization(t *testing.T) {
	var v throttle.Time
	v.DoWait(0, func() {})
	entered, release, returned, finished := make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{})
	go func() {
		v.GoWait(0, func() { close(entered); <-release; close(finished) })
		close(returned)
	}()
	await(t, entered)
	await(t, returned)
	require.False(t, v.GoWait(0, func() { t.Error("active handler must be skipped") }))
	close(release)
	await(t, returned)
	await(t, finished)
}
func TestTimeGoWaitRunsAfterInterval(t *testing.T) {
	var v throttle.Time
	v.DoWait(0, func() {})
	time.Sleep(time.Millisecond)
	done := make(chan struct{})
	v.GoWait(time.Millisecond, func() { close(done) })
	await(t, done)
}
func TestTimeGoWaitSkipsBeforeInterval(t *testing.T) {
	var v throttle.Time
	v.DoWait(0, func() {})
	v.GoWait(time.Hour, func() { t.Error("interval has not elapsed") })
}
func TestTickDoAfterInitialization(t *testing.T) {
	var v throttle.Tick
	v.DoWait(0, func() {})
	entered, release, returned, finished := make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{})
	go func() {
		v.Do(0, func() { close(entered); <-release; close(finished) })
		close(returned)
	}()
	await(t, entered)
	select {
	case <-returned:
		t.Error("synchronous call returned before handler")
	default:
	}
	require.False(t, v.Do(0, func() { t.Error("active handler must be skipped") }))
	close(release)
	await(t, returned)
	await(t, finished)
}
func TestTickDoExactCallInterval(t *testing.T) {
	var v throttle.Tick
	v.DoWait(0, func() {})
	v.Do(3, func() { t.Error("second call must be skipped") })
	v.Do(3, func() { t.Error("third call must be skipped") })
	done := make(chan struct{})
	v.Do(3, func() { close(done) })
	await(t, done)
}
func TestTickDoPanicAfterInitialization(t *testing.T) {
	var v throttle.Tick
	v.DoWait(0, func() {})
	require.Panics(t, func() { v.Do(0, func() { panic("failed") }) })
	done := make(chan struct{})
	go v.Do(0, func() { close(done) })
	await(t, done)
}
func TestTickGoAfterInitialization(t *testing.T) {
	var v throttle.Tick
	v.DoWait(0, func() {})
	entered, release, returned, finished := make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{})
	go func() {
		v.Go(0, func() { close(entered); <-release; close(finished) })
		close(returned)
	}()
	await(t, entered)
	await(t, returned)
	require.False(t, v.Go(0, func() { t.Error("active handler must be skipped") }))
	close(release)
	await(t, returned)
	await(t, finished)
}
func TestTickGoExactCallInterval(t *testing.T) {
	var v throttle.Tick
	v.DoWait(0, func() {})
	v.Go(3, func() { t.Error("second call must be skipped") })
	v.Go(3, func() { t.Error("third call must be skipped") })
	done := make(chan struct{})
	v.Go(3, func() { close(done) })
	await(t, done)
}
func TestTickDoWaitAfterInitialization(t *testing.T) {
	var v throttle.Tick
	v.DoWait(0, func() {})
	entered, release, returned, finished := make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{})
	go func() {
		v.DoWait(0, func() { close(entered); <-release; close(finished) })
		close(returned)
	}()
	await(t, entered)
	select {
	case <-returned:
		t.Error("synchronous call returned before handler")
	default:
	}
	require.False(t, v.DoWait(0, func() { t.Error("active handler must be skipped") }))
	close(release)
	await(t, returned)
	await(t, finished)
}
func TestTickDoWaitExactCallInterval(t *testing.T) {
	var v throttle.Tick
	v.DoWait(3, func() {})
	v.DoWait(3, func() { t.Error("second call must be skipped") })
	v.DoWait(3, func() { t.Error("third call must be skipped") })
	done := make(chan struct{})
	v.DoWait(3, func() { close(done) })
	await(t, done)
}
func TestTickDoWaitPanicAfterInitialization(t *testing.T) {
	var v throttle.Tick
	v.DoWait(0, func() {})
	require.Panics(t, func() { v.DoWait(0, func() { panic("failed") }) })
	done := make(chan struct{})
	go v.DoWait(0, func() { close(done) })
	await(t, done)
}
func TestTickGoWaitAfterInitialization(t *testing.T) {
	var v throttle.Tick
	v.DoWait(0, func() {})
	entered, release, returned, finished := make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{})
	go func() {
		v.GoWait(0, func() { close(entered); <-release; close(finished) })
		close(returned)
	}()
	await(t, entered)
	await(t, returned)
	require.False(t, v.GoWait(0, func() { t.Error("active handler must be skipped") }))
	close(release)
	await(t, returned)
	await(t, finished)
}
func TestTickGoWaitExactCallInterval(t *testing.T) {
	var v throttle.Tick
	v.GoWait(3, func() {})
	v.GoWait(3, func() { t.Error("second call must be skipped") })
	v.GoWait(3, func() { t.Error("third call must be skipped") })
	done := make(chan struct{})
	v.GoWait(3, func() { close(done) })
	await(t, done)
}

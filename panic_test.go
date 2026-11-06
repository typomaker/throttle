package throttle_test

import (
	"github.com/stretchr/testify/require"
	"github.com/typomaker/throttle"
	"testing"
	"time"
)

func TestTimeDoReusableAfterPanic(t *testing.T) {
	var v throttle.Time
	require.PanicsWithValue(t, "handler panic", func() {
		v.Do(0, func() { panic("handler panic") })
	})
	done := make(chan struct{})
	go v.Do(0, func() { close(done) })
	await(t, done)
}

func TestTimeDoWaitReusableAfterPanic(t *testing.T) {
	var v throttle.Time
	require.PanicsWithValue(t, "handler panic", func() {
		v.DoWait(0, func() { panic("handler panic") })
	})
	done := make(chan struct{})
	go v.DoWait(0, func() { close(done) })
	await(t, done)
}

func TestTimeGoWaitReusableAfterPanic(t *testing.T) {
	var v throttle.Time
	require.PanicsWithValue(t, "handler panic", func() {
		v.GoWait(0, func() { panic("handler panic") })
	})
	done := make(chan struct{})
	go v.GoWait(0, func() { close(done) })
	await(t, done)
}

func TestTickDoReusableAfterPanic(t *testing.T) {
	var v throttle.Tick
	require.PanicsWithValue(t, "handler panic", func() {
		v.Do(0, func() { panic("handler panic") })
	})
	done := make(chan struct{})
	go v.Do(0, func() { close(done) })
	await(t, done)
}

func TestTickDoWaitReusableAfterPanic(t *testing.T) {
	var v throttle.Tick
	require.PanicsWithValue(t, "handler panic", func() {
		v.DoWait(0, func() { panic("handler panic") })
	})
	done := make(chan struct{})
	go v.DoWait(0, func() { close(done) })
	await(t, done)
}

func TestTickGoWaitReusableAfterPanic(t *testing.T) {
	var v throttle.Tick
	require.PanicsWithValue(t, "handler panic", func() {
		v.GoWait(0, func() { panic("handler panic") })
	})
	done := make(chan struct{})
	go v.GoWait(0, func() { close(done) })
	await(t, done)
}
func TestTimeDoWaitWaitingCallerRetriesAfterPanic(t *testing.T) {
	var v throttle.Time
	entered, release, failed, retried, returned := make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{})
	go func() {
		defer close(failed)
		defer func() {
			if value := recover(); value != "handler panic" {
				t.Errorf("unexpected panic: %v", value)
			}
		}()
		v.DoWait(0, func() { close(entered); <-release; panic("handler panic") })
	}()
	await(t, entered)
	go func() { v.DoWait(0, func() { close(retried) }); close(returned) }()
	select {
	case <-returned:
		t.Error("caller did not wait for initialization")
	case <-time.After(20 * time.Millisecond):
	}
	close(release)
	await(t, failed)
	await(t, retried)
	await(t, returned)
}
func TestTimeGoWaitWaitingCallerRetriesAfterPanic(t *testing.T) {
	var v throttle.Time
	entered, release, failed, retried, returned := make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{})
	go func() {
		defer close(failed)
		defer func() {
			if value := recover(); value != "handler panic" {
				t.Errorf("unexpected panic: %v", value)
			}
		}()
		v.GoWait(0, func() { close(entered); <-release; panic("handler panic") })
	}()
	await(t, entered)
	go func() { v.GoWait(0, func() { close(retried) }); close(returned) }()
	select {
	case <-returned:
		t.Error("caller did not wait for initialization")
	case <-time.After(20 * time.Millisecond):
	}
	close(release)
	await(t, failed)
	await(t, retried)
	await(t, returned)
}
func TestTickDoWaitWaitingCallerRetriesAfterPanic(t *testing.T) {
	var v throttle.Tick
	entered, release, failed, retried, returned := make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{})
	go func() {
		defer close(failed)
		defer func() {
			if value := recover(); value != "handler panic" {
				t.Errorf("unexpected panic: %v", value)
			}
		}()
		v.DoWait(0, func() { close(entered); <-release; panic("handler panic") })
	}()
	await(t, entered)
	go func() { v.DoWait(0, func() { close(retried) }); close(returned) }()
	select {
	case <-returned:
		t.Error("caller did not wait for initialization")
	case <-time.After(20 * time.Millisecond):
	}
	close(release)
	await(t, failed)
	await(t, retried)
	await(t, returned)
}
func TestTickGoWaitWaitingCallerRetriesAfterPanic(t *testing.T) {
	var v throttle.Tick
	entered, release, failed, retried, returned := make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{})
	go func() {
		defer close(failed)
		defer func() {
			if value := recover(); value != "handler panic" {
				t.Errorf("unexpected panic: %v", value)
			}
		}()
		v.GoWait(0, func() { close(entered); <-release; panic("handler panic") })
	}()
	await(t, entered)
	go func() { v.GoWait(0, func() { close(retried) }); close(returned) }()
	select {
	case <-returned:
		t.Error("caller did not wait for initialization")
	case <-time.After(20 * time.Millisecond):
	}
	close(release)
	await(t, failed)
	await(t, retried)
	await(t, returned)
}

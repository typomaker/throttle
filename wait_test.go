package throttle_test

import (
	"testing"
	"time"

	"github.com/typomaker/throttle"
)

func await(t *testing.T, ch <-chan struct{}) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("call did not complete")
	}
}

func TestTimeDoSkipsActiveFirstRun(t *testing.T) {
	var v throttle.Time
	entered, release := make(chan struct{}), make(chan struct{})
	firstReturned := make(chan struct{})
	go func() {
		v.Do(0, func() { close(entered); <-release })
		close(firstReturned)
	}()
	await(t, entered)
	returned := make(chan struct{})
	go func() {
		if got := v.Do(0, func() { t.Error("concurrent handler must be skipped") }); got != false {
			t.Errorf("concurrent call returned %v, want false", got)
		}
		close(returned)
	}()
	await(t, returned)
	select {
	case <-firstReturned:
		t.Error("first caller returned before synchronous handler completed")
	default:
	}
	close(release)
	await(t, firstReturned)
	await(t, returned)
}

func TestTimeGoSkipsActiveFirstRun(t *testing.T) {
	var v throttle.Time
	entered, release := make(chan struct{}), make(chan struct{})
	firstReturned := make(chan struct{})
	go func() {
		v.Go(0, func() { close(entered); <-release })
		close(firstReturned)
	}()
	await(t, entered)
	returned := make(chan struct{})
	go func() {
		if got := v.Go(0, func() { t.Error("concurrent handler must be skipped") }); got != false {
			t.Errorf("concurrent call returned %v, want false", got)
		}
		close(returned)
	}()
	await(t, returned)
	await(t, firstReturned)
	close(release)
	await(t, firstReturned)
	await(t, returned)
}

func TestTimeDoWaitWaitsForInitialization(t *testing.T) {
	var v throttle.Time
	entered, release := make(chan struct{}), make(chan struct{})
	firstReturned := make(chan struct{})
	go func() {
		v.DoWait(0, func() { close(entered); <-release })
		close(firstReturned)
	}()
	await(t, entered)
	returned := make(chan struct{})
	go func() {
		if got := v.DoWait(0, func() { t.Error("concurrent handler must be skipped") }); got != true {
			t.Errorf("concurrent call returned %v, want true", got)
		}
		close(returned)
	}()
	select {
	case <-returned:
		t.Error("caller returned before initialization completed")
	case <-time.After(20 * time.Millisecond):
	}
	select {
	case <-firstReturned:
		t.Error("first caller returned before synchronous handler completed")
	default:
	}
	close(release)
	await(t, firstReturned)
	await(t, returned)
}

func TestTimeGoWaitWaitsForInitialization(t *testing.T) {
	var v throttle.Time
	entered, release := make(chan struct{}), make(chan struct{})
	firstReturned := make(chan struct{})
	go func() {
		v.GoWait(0, func() { close(entered); <-release })
		close(firstReturned)
	}()
	await(t, entered)
	returned := make(chan struct{})
	go func() {
		if got := v.GoWait(0, func() { t.Error("concurrent handler must be skipped") }); got != true {
			t.Errorf("concurrent call returned %v, want true", got)
		}
		close(returned)
	}()
	select {
	case <-returned:
		t.Error("caller returned before initialization completed")
	case <-time.After(20 * time.Millisecond):
	}
	select {
	case <-firstReturned:
		t.Error("first caller returned before synchronous handler completed")
	default:
	}
	close(release)
	await(t, firstReturned)
	await(t, returned)
}

func TestTickDoSkipsActiveFirstRun(t *testing.T) {
	var v throttle.Tick
	entered, release := make(chan struct{}), make(chan struct{})
	firstReturned := make(chan struct{})
	go func() {
		v.Do(0, func() { close(entered); <-release })
		close(firstReturned)
	}()
	await(t, entered)
	returned := make(chan struct{})
	go func() {
		if got := v.Do(0, func() { t.Error("concurrent handler must be skipped") }); got != false {
			t.Errorf("concurrent call returned %v, want false", got)
		}
		close(returned)
	}()
	await(t, returned)
	select {
	case <-firstReturned:
		t.Error("first caller returned before synchronous handler completed")
	default:
	}
	close(release)
	await(t, firstReturned)
	await(t, returned)
}

func TestTickGoSkipsActiveFirstRun(t *testing.T) {
	var v throttle.Tick
	entered, release := make(chan struct{}), make(chan struct{})
	firstReturned := make(chan struct{})
	go func() {
		v.Go(0, func() { close(entered); <-release })
		close(firstReturned)
	}()
	await(t, entered)
	returned := make(chan struct{})
	go func() {
		if got := v.Go(0, func() { t.Error("concurrent handler must be skipped") }); got != false {
			t.Errorf("concurrent call returned %v, want false", got)
		}
		close(returned)
	}()
	await(t, returned)
	await(t, firstReturned)
	close(release)
	await(t, firstReturned)
	await(t, returned)
}

func TestTickDoWaitWaitsForInitialization(t *testing.T) {
	var v throttle.Tick
	entered, release := make(chan struct{}), make(chan struct{})
	firstReturned := make(chan struct{})
	go func() {
		v.DoWait(0, func() { close(entered); <-release })
		close(firstReturned)
	}()
	await(t, entered)
	returned := make(chan struct{})
	go func() {
		if got := v.DoWait(0, func() { t.Error("concurrent handler must be skipped") }); got != true {
			t.Errorf("concurrent call returned %v, want true", got)
		}
		close(returned)
	}()
	select {
	case <-returned:
		t.Error("caller returned before initialization completed")
	case <-time.After(20 * time.Millisecond):
	}
	select {
	case <-firstReturned:
		t.Error("first caller returned before synchronous handler completed")
	default:
	}
	close(release)
	await(t, firstReturned)
	await(t, returned)
}

func TestTickGoWaitWaitsForInitialization(t *testing.T) {
	var v throttle.Tick
	entered, release := make(chan struct{}), make(chan struct{})
	firstReturned := make(chan struct{})
	go func() {
		v.GoWait(0, func() { close(entered); <-release })
		close(firstReturned)
	}()
	await(t, entered)
	returned := make(chan struct{})
	go func() {
		if got := v.GoWait(0, func() { t.Error("concurrent handler must be skipped") }); got != true {
			t.Errorf("concurrent call returned %v, want true", got)
		}
		close(returned)
	}()
	select {
	case <-returned:
		t.Error("caller returned before initialization completed")
	case <-time.After(20 * time.Millisecond):
	}
	select {
	case <-firstReturned:
		t.Error("first caller returned before synchronous handler completed")
	default:
	}
	close(release)
	await(t, firstReturned)
	await(t, returned)
}

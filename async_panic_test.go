package throttle_test

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/typomaker/throttle"
)

// An unhandled goroutine panic terminates the process, so isolate these checks.
func assertAsyncPanic(t *testing.T, scenario string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestAsyncPanicProcess$")
	cmd.Env = append(os.Environ(), "THROTTLE_PANIC_CASE="+scenario)
	output, err := cmd.CombinedOutput()
	require.NoError(t, ctx.Err(), "child process timed out")
	require.Error(t, err)
	require.Contains(t, string(output), "panic: background handler panic")
}

func TestTimeGoPanicAfterInitialization(t *testing.T) {
	assertAsyncPanic(t, "TimeGo")
}

func TestTimeGoWaitPanicAfterInitialization(t *testing.T) {
	assertAsyncPanic(t, "TimeGoWait")
}

func TestTickGoPanicAfterInitialization(t *testing.T) {
	assertAsyncPanic(t, "TickGo")
}

func TestTickGoWaitPanicAfterInitialization(t *testing.T) {
	assertAsyncPanic(t, "TickGoWait")
}

func TestAsyncPanicProcess(t *testing.T) {
	scenario := os.Getenv("THROTTLE_PANIC_CASE")
	if scenario == "" {
		t.Skip("subprocess helper")
	}
	fail := func() { panic("background handler panic") }
	var timed throttle.Time
	var counted throttle.Tick
	timed.DoWait(0, func() {})
	counted.DoWait(0, func() {})
	switch scenario {
	case "TimeGo":
		timed.Go(0, fail)
	case "TimeGoWait":
		timed.GoWait(0, fail)
	case "TickGo":
		counted.Go(0, fail)
	case "TickGoWait":
		counted.GoWait(0, fail)
	default:
		t.Fatal("unknown scenario")
	}
	select {}
}

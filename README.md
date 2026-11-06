# Throttle

Throttle function calls by elapsed time (`Time`) or by number of calls (`Tick`).
Both types are safe for concurrent use and ready to use with their zero value.

Handlers run only when you call a method. There is no background scheduler,
queue, or automatic retry of skipped calls.

## Install

Requires Go 1.25 or later.

```sh
go get github.com/typomaker/throttle
```

## Choose a method

The same methods exist on `Time` and `Tick`:

| Method | First handler | Concurrent callers during initialization | Later handlers |
| --- | --- | --- | --- |
| `Do` | Synchronous | Return immediately | Synchronous |
| `Go` | Goroutine | Return immediately | Goroutine |
| `DoWait` | Synchronous | Wait for successful initialization | Synchronous |
| `GoWait` | Synchronous | Wait for successful initialization | Goroutine |

All methods admit the first call immediately, regardless of the interval.
Only one handler executes per instance at a time. After initialization,
concurrent calls skip an active handler, including with waiting methods.
A synchronous call waits only when it executes its own handler, except for
initialization waiting in `DoWait` and `GoWait`.

Methods on the same instance share the same execution history and lock.
Use a separate instance for each independently throttled operation.

## Return values

All methods return a `bool`:

- `Do`: `true` after the handler returns normally; `false` if skipped.
- `Go`: `true` if a goroutine was scheduled; `false` if skipped. It does not
  report handler completion or success.
- `DoWait` and `GoWait`: `true` if this call ran the initial handler or waited
  for another caller to initialize successfully. After initialization they
  return the result of `Do` or `Go`, including `false` on skipped calls.

A synchronous panic propagates without returning a result. Existing calls may
ignore the result; method values assigned to function types must include the
new `bool` return value.

## Time intervals

Pass a `time.Duration` to limit how often a handler can start:

```go
package main

import (
    "fmt"
    "time"

    "github.com/typomaker/throttle"
)

func main() {
    var limiter throttle.Time
    handler := func() { fmt.Println("refresh") }

    limiter.Do(10*time.Second, handler) // Runs immediately.
    limiter.Do(10*time.Second, handler) // Skipped: interval has not elapsed.
    limiter.Do(0, handler)              // Runs immediately if no handler is active.
}
```

`Time` uses a monotonic elapsed clock. The timestamp is captured after acquiring
the lock, before invoking the handler (`Do`, `DoWait`, and initial `GoWait`) or
creating its goroutine (`Go` and subsequent `GoWait`). It is committed only after
the handler returns normally.

Handler duration counts toward the interval: with a 10-second interval and a
3-second handler, the next run is eligible about 7 seconds after completion.
If the handler takes longer than the interval, a new run can be admitted on the
next call after completion. Goroutine scheduling can delay actual execution.

Intervals are measured from the last successful run's recorded timestamp.
Missed intervals do not accumulate runs, and there is no fixed wall-clock schedule.

## Call intervals

Pass an `int` to `Tick`. After the first run, it counts calls from the call that
started the last successful run:

```go
package main

import (
    "fmt"

    "github.com/typomaker/throttle"
)

func main() {
    var limiter throttle.Tick
    for call := 1; call <= 21; call++ {
        limiter.Do(10, func() { fmt.Println(call) })
    }
    // Prints 1, 11, 21.
}
```

Calls skipped while a handler is active also count toward the next interval.
Concurrent calls can still be skipped because another run has won admission;
this is a throttle, not a guarantee to execute every tenth call.

## Wait for initialization, then refresh asynchronously

Use `GoWait` when every caller needs initialized data before returning, while
later refreshes should run in the background. Here the handler publishes an
immutable snapshot through an atomic pointer:

```go
package main

import (
    "fmt"
    "sync/atomic"
    "time"

    "github.com/typomaker/throttle"
)

type Snapshot struct {
    Value string
}

type Cache struct {
    refresh throttle.Time
    current atomic.Pointer[Snapshot]
}

func (c *Cache) Get() *Snapshot {
    c.refresh.GoWait(time.Minute, func() {
        // Load a new snapshot from your data source here.
        c.current.Store(&Snapshot{Value: "loaded"})
    })
    return c.current.Load()
}

func main() {
    var cache Cache
    fmt.Println(cache.Get().Value)
}
```

Until the first handler returns normally, concurrent `Get` calls wait. Later
calls can return the previous snapshot while a refresh is running. Callers must
not modify the published snapshot. Use `DoWait` if later refreshes should execute
synchronously for the caller admitted to run them.

Initialization means the handler returned normally. These methods do not inspect
application errors: a handler that returns after a failed load still counts as
successful. Use waiting methods consistently for callers that require initialization.

## Bypass the interval

An interval of zero or less bypasses the time/call threshold. It does not bypass
the lock or guarantee execution under contention:

```go
var limiter throttle.Tick
if refreshNormally {
    limiter.Do(10, handler)
} else {
    limiter.Do(0, handler)
}
```

If admitted, a normally returning handler updates the shared execution history.
The next call with a positive interval counts from that run. `Go(0, handler)`
works the same way but executes asynchronously. This also applies to `Time`.

## Panics and safe use

- The library does not recover panics. Synchronous panics propagate to the caller;
  an unhandled panic in a background goroutine terminates the process. Recovery
  for background handlers must happen inside the handler itself.
- The lock is released even if a handler panics. A propagated panic does not
  update execution history. During initialization, a waiting caller can retry.
  If the handler recovers its own panic and returns normally, the run counts as
  successful.
- During initialization, a handler must not recursively call `DoWait` or
  `GoWait` on the same instance: this deadlocks.
- Callers and their handlers are responsible for synchronizing user data.
  The throttle serializes handlers, but does not protect concurrent readers.
- Do not copy `Time` or `Tick` after first use. Pass pointers instead.

<h3 align="center">Throttle</h3>

---

<p align="center"> Call the function once in a certain time interval.
    <br>
</p>

## 📝 Table of Contents

- [📝 Table of Contents](#-table-of-contents)
- [🧐 About ](#-about-)
- [🏁 Getting Started ](#-getting-started-)
	- [Prerequisites](#prerequisites)
	- [Installing](#installing)
- [🎈 Usage ](#-usage-)

## 🧐 About <a name = "about"></a>

Calls the function once in the specified time interval. For example, it can be used to load data from the database into memory every 10 minutes without raising background tasks.

## 🏁 Getting Started <a name = "getting_started"></a>

### Prerequisites

```
go version 1.25
```

### Installing
```
go get github.com/typomaker/throttle
```
## 🎈 Usage <a name="usage"></a>
```go
package main

import (
	"fmt"
	"time"

	"github.com/typomaker/throttle"
)

func main() {
	var oncePer2Second throttle.Time
	var oncePer4Second throttle.Time
	for {
		oncePer2Second.Do(2*time.Second, func() {
			fmt.Println("Do once per 2 second, blocking", time.Now())
		})
		oncePer4Second.Go(4*time.Second, func() {
			fmt.Println("Go once per 4 second, non blocking", time.Now())
		})
	}
}
```

## First-run behavior

`Do` executes the handler synchronously. `Go` executes it in a goroutine,
including the first run. Concurrent calls skip an active handler.

Use `DoWait` or `GoWait` when callers must wait for initialization:

```go
var cache throttle.Time
cache.GoWait(time.Minute, loadCache)
```

Both methods execute the first handler synchronously and block concurrent
waiting calls until it completes successfully. Subsequent `DoWait` runs are
synchronous; subsequent `GoWait` runs execute in a goroutine. After initialization,
concurrent calls skip active runs as usual. The same methods exist on `Tick`. All its methods admit the first call immediately;
subsequent intervals count calls from the last successful run, including calls
skipped while a handler is active.
Use waiting methods consistently for callers that require initialization.

Panics propagate to the caller (or the executing goroutine). The lock is released,
and a panicking handler is not recorded as a successful run. Waiting callers can
retry initialization.

`Time` and `Tick` must not be copied after first use. Callers and their handlers
are responsible for synchronizing user data. During initialization, a handler
must not recursively call `DoWait` or `GoWait` on the same instance: this deadlocks.

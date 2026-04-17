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

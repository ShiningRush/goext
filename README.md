# goext
[![Go Reference](https://pkg.go.dev/badge/github.com/shiningrush/goext.svg)](https://pkg.go.dev/github.com/shiningrush/goext)
[![Go Report Card](https://goreportcard.com/badge/github.com/shiningrush/goext)](https://goreportcard.com/report/github.com/shiningrush/goext)

`goext` is a small collection of practical helpers for Go projects, including:

- set and slice utilities
- generic zero-value helpers
- lightweight parallel processing helpers
- in-memory event bus utilities
- simple job scheduling helpers
- time helpers that are easy to mock in tests
- small error aggregation helpers

## Requirements

- Go `1.18+`

## Installation

```bash
go get github.com/shiningrush/goext
```

## Packages

| Package | What it provides |
| --- | --- |
| `datax` | A simple `Set` type and generic slice set-operations |
| `errx` | Batch error aggregation |
| `gtx` | Generic zero-value helpers |
| `intx` | Integer slice helpers |
| `int32x` | `int32` slice helpers |
| `parallel` | Parallel workers and streaming worker sessions |
| `runx/eventx` | In-memory pub/sub event bus |
| `runx/jobx` | Once, interval, and cron-style jobs |
| `stringx` | String slice helpers |
| `timex` | Mockable clock helpers and duration parsing |

## datax

`datax` contains a simple set implementation plus generic helpers for common slice operations.

### Set

Main capabilities:

- `NewSet()` and `NewSetFrom(...)`
- `Add`, `Remove`, `TryAdd`, `Has`, `Len`, `All`
- `Equal`, `IsSubsetOf`, `IsProperSubsetOf`, `IsSupersetOf`
- `Intersect`, `Union`, `Diff`

```go
package main

import (
	"fmt"

	"github.com/shiningrush/goext/datax"
)

func main() {
	s := datax.NewSet().Add("a", "b", "c", 4)

	fmt.Println(s.Len())      // 4
	fmt.Println(s.Has("b"))   // true
	fmt.Println(s.TryAdd(4))  // false

	s.Remove("c")
	fmt.Println(s.All())

	another := datax.NewSet().Add("b", 4)
	fmt.Println(s.Intersect(another).All())
	fmt.Println(another.IsSubsetOf(s))
}
```

### Generic slice helpers

`datax` also provides generic helpers for comparable or set-like slice operations:

- `HasItem`
- `IsSuperset`
- `IsSubset`
- `IsProperSubset`
- `Intersect`
- `Diff`

```go
package main

import (
	"fmt"

	"github.com/shiningrush/goext/datax"
)

func main() {
	a := []string{"a", "b", "c"}
	b := []string{"b", "c", "d"}

	fmt.Println(datax.HasItem(a, "a"))   // true
	fmt.Println(datax.IsSubset([]string{"b"}, a))
	fmt.Println(datax.Intersect(a, b))   // [b c]
	fmt.Println(datax.Diff(a, b))        // [a]
}
```

## errx

`errx` provides a small error collector for accumulating multiple errors and returning them as one value.

```go
package main

import (
	"errors"
	"fmt"

	"github.com/shiningrush/goext/errx"
)

func main() {
	batch := &errx.BatchErrors{}
	batch.Append(errors.New("create user failed"))
	batch.Append(errors.New("send email failed"))

	if batch.HasError() {
		fmt.Println(batch.Len())   // 2
		fmt.Println(batch.Error()) // formatted multi-line message
	}
}
```

## gtx

`gtx` provides generic helpers around zero values.

```go
package main

import (
	"fmt"

	"github.com/shiningrush/goext/gtx"
)

func main() {
	fmt.Println(gtx.Zero[int]())      // 0
	fmt.Println(gtx.Zero[string]())   // ""
	fmt.Println(gtx.IsZero(0))        // true
	fmt.Println(gtx.IsZero("hello"))  // false
}
```

## stringx, intx, int32x

These packages provide typed wrappers around the generic slice helpers in `datax`.

Available helpers:

- `HasItem`
- `IsSuperset`
- `IsSubset`
- `IsProperSubset`
- `Intersect`
- `Diff`

### stringx

```go
package main

import (
	"fmt"

	"github.com/shiningrush/goext/stringx"
)

func main() {
	a := []string{"prod", "gray", "test"}
	b := []string{"gray", "test"}

	fmt.Println(stringx.HasItem(a, "prod")) // true
	fmt.Println(stringx.IsSuperset(a, b))   // true
	fmt.Println(stringx.Diff(a, b))         // [prod]
}
```

### intx

```go
package main

import (
	"fmt"

	"github.com/shiningrush/goext/intx"
)

func main() {
	a := []int{1, 2, 3}
	b := []int{2, 3, 4}

	fmt.Println(intx.Intersect(a, b)) // [2 3]
	fmt.Println(intx.Diff(a, b))      // [1]
}
```

### int32x

```go
package main

import (
	"fmt"

	"github.com/shiningrush/goext/int32x"
)

func main() {
	a := []int32{1, 2, 3}
	b := []int32{1, 3}

	fmt.Println(int32x.IsSuperset(a, b)) // true
}
```

## parallel

`parallel` provides two ways to process work concurrently:

- `Do`: start `N` workers to run the same function
- `StreamDo`: create a worker session and stream items into it

### Do

`Do` runs a `Work` function once per worker and returns all non-nil errors.

- The default worker count is `runtime.GOMAXPROCS(0)`
- Use `parallel.WorkerNumber(n)` to override it

```go
package main

import (
	"errors"
	"fmt"

	"github.com/shiningrush/goext/parallel"
)

func main() {
	errs := parallel.Do(func(workerIdx int) error {
		if workerIdx == 1 {
			return errors.New("worker 1 failed")
		}
		return nil
	}, parallel.WorkerNumber(4))

	fmt.Println(len(errs))
}
```

### StreamDo

`StreamDo` is useful when you have a stream of input items and want workers to process them continuously.

Key types:

- `StreamWork[P, R]`
- `StreamSession[P, R]`
- `StreamPayload[R]`
- `StreamPayloads[R]`

Common session methods:

- `Send`
- `CompleteSend`
- `Wait`
- `ReceivedPayloads`
- `ReceiveChan`

```go
package main

import (
	"fmt"
	"strings"

	"github.com/shiningrush/goext/parallel"
)

func main() {
	session := parallel.StreamDo(func(workerIdx int, item string) (string, error) {
		return fmt.Sprintf("worker=%d result=%s", workerIdx, strings.ToUpper(item)), nil
	}, parallel.WorkerNumber(2))

	session.Send("alpha")
	session.Send("beta")
	session.Send("gamma")

	payloads := session.ReceivedPayloads()
	if err := payloads.AggregateErr(); err != nil {
		panic(err)
	}

	for _, payload := range payloads {
		fmt.Println(payload.Result)
	}
}
```

### Stream options

- `parallel.WorkerNumber(n)`: set worker count
- `parallel.IgnoreResult()`: do not collect worker results
- `parallel.ReceiveDataFromChan()`: consume payloads manually from `ReceiveChan()`

When the input type implements `parallel.KeyOwner`, items with the same key are dispatched to the same worker.

```go
type UserJob struct {
	UserID string
}

func (j UserJob) GetKey() string {
	return j.UserID
}
```

This is useful when the same key must be processed in order by the same worker.

## runx/eventx

`eventx` provides a simple in-memory pub/sub event bus.

Common flow:

1. Define an event type by implementing `Topic() []string`
2. Define a handler by implementing `Topic() []string` and `Handle(...)`
3. Register handlers with `eventx.Subscribe(...)`
4. Publish events with `eventx.Publish(...)` or `eventx.PublishSync(...)`

Notes:

- `Publish` dispatches asynchronously and returns immediately
- `PublishSync` waits until all matched handlers finish
- `Close` waits for async handlers to complete before shutdown
- Publishing after `Close` will panic
- `Subscribe` is safe, but should not be called at the same time as `Publish`

```go
package main

import (
	"context"
	"fmt"

	"github.com/shiningrush/goext/runx/eventx"
)

type UserCreatedEvent struct {
	ID   int64
	Name string
}

func (e *UserCreatedEvent) Topic() []string {
	return []string{"user.created"}
}

type AuditHandler struct{}

func (h *AuditHandler) Topic() []string {
	return []string{"user.created"}
}

func (h *AuditHandler) Handle(ctx context.Context, event eventx.Event) {
	if userEvent, ok := event.(*UserCreatedEvent); ok {
		fmt.Printf("audit log: user created, id=%d, name=%s\n", userEvent.ID, userEvent.Name)
	}
}

func main() {
	if err := eventx.Subscribe(&AuditHandler{}); err != nil {
		panic(err)
	}

	e := &UserCreatedEvent{
		ID:   1001,
		Name: "alice",
	}

	eventx.Publish(e)
	eventx.PublishSync(context.TODO(), e)
	eventx.Close()
}
```

If you need an isolated bus instead of the package-level default bus, create one with `eventx.NewInMemoryEventBus()`.

## runx/jobx

`jobx` provides a lightweight job runner for three kinds of jobs:

- once jobs
- interval jobs
- cron jobs

You can use the package-level singleton via:

- `RegisterJob`
- `RegisterJobDesc`
- `Start`
- `Stop`
- `RemoveJob`
- `ClearJobs`

You can also create an isolated scheduler with `jobx.NewJobDemon()`.

### Once job

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/shiningrush/goext/runx/jobx"
)

func main() {
	jobx.RegisterJob("warm-cache", jobx.JobType{
		Once: &jobx.OnceJobDesc{
			Delay: 3 * time.Second,
		},
	}, func(ctx context.Context) {
		fmt.Println("cache warmed")
	})

	jobx.Start()
	defer jobx.Stop()
}
```

### Interval job

```go
jobx.RegisterJob("sync-metrics", jobx.JobType{
	Interval: &jobx.IntervalJobDesc{
		Interval: 30 * time.Second,
	},
}, func(ctx context.Context) {
	// do work
})
```

### Cron job

`jobx` uses [`robfig/cron/v3`](https://github.com/robfig/cron) internally.

```go
jobx.RegisterJob("daily-report", jobx.JobType{
	Cron: &jobx.CronJobDesc{
		Spec: "0 0 9 * * *",
	},
}, func(ctx context.Context) {
	// run at 09:00:00 every day
})
```

Notes:

- `Start()` is idempotent while the scheduler is already running
- `Stop()` stops interval and cron jobs by closing the shared channel
- A once job runs only once by default; set `AlwaysStart: true` to run it again after restarting

## timex

`timex` provides a mockable clock abstraction and a small duration parser.

Main APIs:

- `Now()`
- `After(d)`
- `ParseDuration(text)`
- `MockClockImpl(clock)`

### Mockable clock

```go
package main

import (
	"fmt"
	"time"

	"github.com/shiningrush/goext/timex"
)

func main() {
	now := timex.Now()
	fmt.Println(now)

	<-timex.After(100 * time.Millisecond)
	fmt.Println("done")
}
```

### Duration parsing

`ParseDuration` supports these units:

- `d`
- `h`
- `m`
- `s`
- `ms`

```go
package main

import (
	"fmt"

	"github.com/shiningrush/goext/timex"
)

func main() {
	d, err := timex.ParseDuration("4h")
	if err != nil {
		panic(err)
	}

	fmt.Println(d)
}
```

## Development

Common commands:

```bash
make install
make tidy
make test
```

`make install` installs tools used by this repository, including `mockgen` and `goimports`.

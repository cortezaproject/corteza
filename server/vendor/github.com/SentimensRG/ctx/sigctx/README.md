# sigctx

Package `sigctx` provides contexts for graceful shutdown.

The `sigctx` package provides a context that terminates when it receives a
SIGINT or SIGTERM.  This provides a convenient mechanism for triggering
graceful application shutdown.

`sigctx.New` returns a `ctx.C`, which implements the ubiquitous `ctx.Doner`
interface.  It fires when either SIGINT or SIGTERM is caught.

## Examples

```go
import (
    "log"
    "github.com/SentimensRG/ctx/sigctx"
)

func main() {
    ctx := sigctx.New()  // returns a regular context.Context

    <-ctx.Done()  // will unblock on SIGINT and SIGTERM
    log.Println("exiting.")
}
```

`sigctx.Tick` can be used to react to streams of signals.  For example, you can
implement a graceful shutdown attempt, followed by a forced shutdown.

```go
import (
    "log"
    "github.com/SentimensRG/ctx/sigctx"
    "github.com/SentimensRG/ctx"
)

func main() {
    t := sigctx.Tick()
    d, cancel := ctx.WithCancel(ctx.Background())

    go func() {
        defer cancel()

        go func() {
            // business logic goes here
        }()

        <-t
        log.Println("attempting graceful shutdown - press Ctrl + c again to force quit")
        go func() {
            defer cancel()
            // cleanup logic goes here
        }()

        <-t
        log.Println("forcing close")
    }()

    <-d.Done()
}

```
package sigctx

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var ctx context.Context

func init() {
	var cancel func()
	ctx, cancel = context.WithCancel(context.Background())

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ch:
			cancel()
		case <-ctx.Done():
		}
	}()
}

// New signal-bound context
func New() context.Context {
	return ctx
}

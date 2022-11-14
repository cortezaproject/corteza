// Package sigctx provides a context that expires when a SIGINT or SIGTERM is
// received.
package sigctx

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/SentimensRG/ctx"
)

var (
	c              ctx.C
	sigCh          chan os.Signal
	initC, initSig sync.Once
)

func initSigCh() {
	sigCh = make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
}

// New signal-bound ctx.C that terminates when either SIGINT or SIGTERM
// is caught.
func New() ctx.C {
	initC.Do(func() {
		initSig.Do(initSigCh)

		dc := make(chan struct{})
		c = dc

		go func() {
			select {
			case <-sigCh:
				close(dc)
			case <-c.Done():
			}
		}()
	})

	return c
}

// Tick returns a channel that recvs each time a either SIGINT or SIGTERM are
// caught.
func Tick() <-chan struct{} {
	initSig.Do(initSigCh)

	dc := make(chan struct{})
	go func() {
		for {
			<-sigCh
			dc <- struct{}{}
		}
	}()

	return dc
}

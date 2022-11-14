package ctx

import (
	"context"
	"sync"
	"time"
)

// Binder is the interface that wraps the basic Bind method.
// Bind executes logic until the Doner completes.  Implementations of Bind must
// not return until the Doner has completed.
type Binder interface {
	Bind(Doner)
}

// BindFunc is an adapter to allow the use of ordinary functions as Binders.
type BindFunc func(Doner)

// Bind executes logic until the Doner completes.  It satisfies the Binder
// interface.
func (f BindFunc) Bind(d Doner) { f(d) }

// Doner can block until something is done
type Doner interface {
	Done() <-chan struct{}
}

// C is a basic implementation of Doner
type C <-chan struct{}

// Background is the ctx analog to context.Background().  It never fires.
func Background() C {
	return nil
}

// Done returns a channel that receives when an action is complete
func (dc C) Done() <-chan struct{} { return dc }

type ctx struct {
	Doner
}

// Deadline returns the time when work done on behalf of this context
// should be canceled. Deadline returns ok==false when no deadline is
// set. Successive calls to Deadline return the same results.
func (ctx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c ctx) Err() error {
	select {
	case <-c.Done():
		return context.Canceled
	default:
		return nil
	}
}

func (c ctx) Value(interface{}) (v interface{}) {
	return
}

// AsContext creates a context that fires when the Doner fires
func AsContext(d Doner) context.Context {
	return ctx{d}
}

// After time time has elapsed, the Doner fires
func After(d time.Duration) C {
	ch := make(chan struct{})
	go func() {
		<-time.After(d)
		close(ch)
	}()
	return ch
}

// WithCancel returns a new Doner that can be cancelled via the associated
// function
func WithCancel(d Doner) (C, func()) {
	var closer sync.Once
	cq := make(chan struct{})
	cancel := func() { closer.Do(func() { close(cq) }) }

	go func() {
		select {
		case <-cq:
		case <-d.Done():
			cancel()
		}
	}()

	return cq, cancel
}

// Tick returns a <-chan whose range ends when the underlying context cancels
func Tick(d Doner) <-chan struct{} {
	c := make(chan struct{})
	cq := d.Done()
	go func() {
		for {
			select {
			case <-cq:
				close(c)
				return
			default:
				select {
				case c <- struct{}{}:
				case <-cq:
				}
			}
		}
	}()
	return c
}

// Defer guarantees that a function will be called after a context has cancelled
func Defer(d Doner, cb func()) {
	go func() {
		<-d.Done()
		cb()
	}()
}

// Link returns a channel that fires if ANY of the constituent Doners have fired
func Link(doners ...Doner) C {
	c := make(chan struct{})
	cancel := func() { close(c) }

	var once sync.Once
	for _, d := range doners {
		Defer(d, func() { once.Do(cancel) })
	}

	return c
}

// Join returns a channel that receives when all constituent Doners have fired
func Join(doners ...Doner) C {
	var wg sync.WaitGroup
	wg.Add(len(doners))
	for _, d := range doners {
		Defer(d, wg.Done)
	}

	cq := make(chan struct{})
	go func() {
		wg.Wait()
		close(cq)
	}()
	return cq
}

// FTick calls a function in a loop until the Doner has fired
func FTick(d Doner, f func()) {
	for range Tick(d) {
		f()
	}
}

// FTickInterval calls a function repeatedly at a given internval, until the Doner
// has fired.  Note that FTickInterval ignores the time spent executing a function,
// and instead guarantees an interval of `t` between of return of the previous
// function call and the invocation of the next function call.
func FTickInterval(d Doner, t time.Duration, f func()) {
	for {
		select {
		case <-d.Done():
			return
		case <-time.After(t):
			f()
		}
	}
}

// FDone returns a doner that fires when the function returns or panics
func FDone(f func()) C {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		f()
	}()
	return ch
}

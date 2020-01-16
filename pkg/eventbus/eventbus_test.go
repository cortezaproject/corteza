package eventbus

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func TestEventbusRegUnreg(t *testing.T) {
	var (
		a = assert.New(t)

		bus = New()
	)

	a.Empty(bus.handlers)
	h1 := bus.Register(nil)
	a.NotZero(h1)
	h2 := bus.Register(nil)
	a.NotZero(h2)
	h3 := bus.Register(nil)
	a.NotZero(h3)
	a.Len(bus.handlers, 3)
	bus.Unregister(h1)
	a.Len(bus.handlers, 2)
	bus.Unregister(h2)
	a.Len(bus.handlers, 1)
	bus.Unregister(h3)
	a.Empty(bus.handlers)
}

func BenchmarkEventbusHandlerLookup(b *testing.B) {
	var (
		bus  = New()
		ptrs []uintptr
	)

	// Register handlers
	for n := 0; n < b.N; n++ {
		ptrs = append(ptrs, bus.Register(nil))
	}

	// Register find & fire handlers
	for n := 0; n < b.N; n++ {
		bus.find(nil)
	}
}

func BenchmarkEventbusRegUnreg(b *testing.B) {
	var (
		bus  = New()
		ptrs []uintptr
	)

	// Register handlers
	for n := 0; n < b.N; n++ {
		ptrs = append(ptrs, bus.Register(nil))
	}

	// Shuffle stored pinters
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(ptrs), func(i, j int) { ptrs[i], ptrs[j] = ptrs[j], ptrs[i] })

	// Unregister all pointers
	for _, ptr := range ptrs {
		bus.Unregister(ptr)
	}
}

func TestEventFiring(t *testing.T) {
	var (
		a = assert.New(t)

		// Let's make sure our simple tests don't
		// cause any extra stress with dat race conditions
		i = &atomic.Int32{}

		ev = func(ev string) Event {
			return &mockEvent{rType: "resource", eType: ev}
		}

		ctx = context.Background()
		bus = New()
	)

	bus.Register(func(ctx context.Context, ev Event) error { i.Inc(); return nil }, On("inc"), For("resource"))
	bus.Register(func(ctx context.Context, ev Event) error { i.Dec(); return nil }, On("dec"), For("resource"))
	bus.Register(func(ctx context.Context, ev Event) error { return fmt.Errorf("handl-err") }, On("err"), For("resource"))

	a.Equal(int32(0), i.Load())
	a.NoError(bus.WaitFor(ctx, ev("inc")))
	a.Equal(int32(1), i.Load())
	a.NoError(bus.WaitFor(ctx, ev("inc")))
	a.Equal(int32(2), i.Load())
	a.NoError(bus.WaitFor(ctx, ev("dec")))
	a.Equal(int32(1), i.Load())
	a.EqualError(bus.WaitFor(ctx, ev("err")), "handl-err")

	bus.Dispatch(ctx, ev("dec"))
	bus.Dispatch(ctx, ev("inc"))
	bus.Dispatch(ctx, ev("inc"))
	bus.wait()
	a.Equal(int32(2), i.Load())
}

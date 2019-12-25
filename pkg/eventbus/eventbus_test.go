package eventbus

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventbusRegUnreg(t *testing.T) {
	var (
		a = assert.New(t)

		bus = New()
	)

	a.Empty(bus.triggers)
	h1 := bus.Register(nil)
	a.NotZero(h1)
	h2 := bus.Register(nil)
	a.NotZero(h2)
	h3 := bus.Register(nil)
	a.NotZero(h3)
	a.Len(bus.triggers, 3)
	bus.Unregister(h1)
	a.Len(bus.triggers, 2)
	bus.Unregister(h2)
	a.Len(bus.triggers, 1)
	bus.Unregister(h3)
	a.Empty(bus.triggers)
}

func BenchmarkEventbusTriggerLookup(b *testing.B) {
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

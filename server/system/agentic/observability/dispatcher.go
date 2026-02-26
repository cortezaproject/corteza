package observability

import "sync"

type (
	Dispatcher interface {
		ID() string
		OnSpan(span AgentSpan) error
		OnEvent(event AgentEvent) error
		Flush() error
		Shutdown() error
	}

	Bus struct {
		l           sync.RWMutex
		dispatchers []Dispatcher
	}
)

func NewBus(dd ...Dispatcher) *Bus {
	return &Bus{dispatchers: dd}
}

func (b *Bus) Register(d Dispatcher) {
	b.l.Lock()
	defer b.l.Unlock()
	b.dispatchers = append(b.dispatchers, d)
}

func (b *Bus) EmitSpan(span AgentSpan) {
	b.l.RLock()
	defer b.l.RUnlock()
	for _, d := range b.dispatchers {
		_ = d.OnSpan(span)
	}
}

func (b *Bus) EmitEvent(event AgentEvent) {
	b.l.RLock()
	defer b.l.RUnlock()
	for _, d := range b.dispatchers {
		_ = d.OnEvent(event)
	}
}

func (b *Bus) Flush() {
	b.l.RLock()
	defer b.l.RUnlock()
	for _, d := range b.dispatchers {
		_ = d.Flush()
	}
}

func (b *Bus) Shutdown() {
	b.l.RLock()
	defer b.l.RUnlock()
	for _, d := range b.dispatchers {
		_ = d.Shutdown()
	}
}

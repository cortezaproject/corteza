package gatekeep

import "sync"

type (
	listenerWrap struct {
		id int
		fn eventListener
	}

	inMemBuss struct {
		mux       sync.RWMutex
		listeners []listenerWrap
	}
)

var (
	crtID = int(0)
)

const (
	inMemBussListenerCap = 5000
)

// Subscribe adds a new listener to the in-memory event bus
// The returned reference should be used to manage the subscribed listener.
func (q *inMemBuss) Subscribe(listener eventListener) int {
	q.mux.Lock()
	defer q.mux.Unlock()

	// @todo should we have some more graceful fail for this?
	// The listener is meant to be more light weight so this will probably be fine.
	// We could consider (rather, the architecture permits) more hard core EBs
	if len(q.listeners)+1 == inMemBussListenerCap {
		panic("inMemBuss: too many listeners")
	}

	q.listeners = append(q.listeners, listenerWrap{
		id: crtID,
		fn: listener,
	})

	crtID++
	return q.listeners[len(q.listeners)-1].id
}

// Unsubscribe removes a listener from the in-memory event bus
func (q *inMemBuss) Unsubscribe(id int) {
	q.mux.Lock()
	defer q.mux.Unlock()

	aux := make([]listenerWrap, 0, len(q.listeners))

	for _, l := range q.listeners {
		if l.id == id {
			continue
		}
		aux = append(aux, l)
	}

	q.listeners = aux
}

// Publish sends an event to all listeners
// The publish doesn't filter anything; it's up to the listeners to decide
// if they are interested in the event or not.
func (q *inMemBuss) Publish(event Event) {
	q.mux.RLock()
	defer q.mux.RUnlock()

	// @todo parallelize
	for _, wrap := range q.listeners {
		wrap.fn(event)
	}
}

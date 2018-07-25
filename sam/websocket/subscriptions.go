package websocket

import (
	"sync"
)

type (
	// A subscription holds a "channel" that the user is joined to
	Subscription struct {
		// for tracking new messages
		lastCommentID string
	}

	// A list of all user-joined channels
	Subscriptions struct {
		sync.RWMutex

		Subscriptions map[string]*Subscription
	}
)

func (Subscriptions) New() *Subscriptions {
	return &Subscriptions{sync.RWMutex{}, make(map[string]*Subscription)}
}

// @todo: load/save all subscriptions from database

func (s *Subscriptions) Add(name string, sub *Subscription) string {
	s.Lock()
	defer s.Unlock()
	s.Subscriptions[name] = sub
	return name
}

func (s *Subscriptions) Get(name string) *Subscription {
	s.RLock()
	defer s.RUnlock()
	return s.Subscriptions[name]
}

func (s *Subscriptions) Delete(name string) {
	s.Lock()
	defer s.Unlock()
	delete(s.Subscriptions, name)
}

func (s *Subscriptions) DeleteAll() {
	s.Lock()
	defer s.Unlock()
	for index, _ := range s.Subscriptions {
		delete(s.Subscriptions, index)
	}
}

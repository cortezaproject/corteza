package websocket

import (
	"sync"
)

type (
	// A subscription holds a "channel" that the user is joined to
	Subscription struct {
		// for tracking new messages
		// lastCommentID string
	}

	// A list of all user-joined channels
	Subscriptions struct {
		sync.RWMutex

		Subscriptions map[string]*Subscription
	}
)

func NewSubscriptions() *Subscriptions {
	return &Subscriptions{sync.RWMutex{}, make(map[string]*Subscription)}
}

// @todo: load/save all subscriptions from database

func (s *Subscriptions) Add(channelID string) *Subscription {
	s.Lock()
	defer s.Unlock()
	s.Subscriptions[channelID] = &Subscription{}
	return s.Subscriptions[channelID]
}

func (s *Subscriptions) Get(channelID string) *Subscription {
	s.RLock()
	defer s.RUnlock()
	return s.Subscriptions[channelID]
}

func (s *Subscriptions) Delete(channelID string) {
	s.Lock()
	defer s.Unlock()
	delete(s.Subscriptions, channelID)
}

func (s *Subscriptions) DeleteAll() {
	s.Lock()
	defer s.Unlock()
	for index := range s.Subscriptions {
		delete(s.Subscriptions, index)
	}
}

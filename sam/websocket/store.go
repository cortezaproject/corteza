package websocket

import (
	"github.com/crusttech/crust/sam/websocket/outgoing"
	"github.com/titpetric/factory"
	"sync"
)

type (
	Store struct {
		sync.RWMutex

		Sessions map[uint64]*Session
	}
)

func (Store) New() *Store {
	return &Store{sync.RWMutex{}, make(map[uint64]*Session)}
}

var store *Store

func init() {
	store = Store{}.New()
}

func (s *Store) Save(session *Session) *Session {
	session.id = factory.Sonyflake.NextID()
	s.Lock()
	defer s.Unlock()
	s.Sessions[session.id] = session
	return session
}

func (s *Store) Get(id uint64) *Session {
	s.RLock()
	defer s.RUnlock()
	return s.Sessions[id]
}

func (s *Store) Delete(id uint64) {
	s.Lock()
	defer s.Unlock()
	delete(s.Sessions, id)
}

func (s *Store) MessageFanout(messages ...*outgoing.WsMessage) {
	// @todo this should probably implement some logic behind...
	for _, message := range messages {
		for _, sess := range s.Sessions {
			sess.send <- message
		}
	}
}

package websocket

import (
	"github.com/cortezaproject/corteza-server/pkg/id"
	"sync"
)

var store *Store

type (
	Store struct {
		sync.RWMutex

		Sessions map[uint64]*Session
	}
)

func init() {
	store = NewStore()
}

func NewStore() *Store {
	return &Store{sync.RWMutex{}, make(map[uint64]*Session)}
}

func (s *Store) Save(session *Session) *Session {
	session.id = id.Next()
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

func (s *Store) Walk(callback func(*Session)) {
	s.RLock()
	defer s.RUnlock()
	for _, sess := range s.Sessions {
		callback(sess)
	}
}

func (s *Store) CountConnections(userID uint64) (count uint) {
	s.Walk(func(session *Session) {
		if session.user.Identity() == userID {
			count++
		}
	})

	return
}

// GetConnectedUsers gets all connected user to ws session
func GetConnectedUsers() []uint64 {
	var chk = map[uint64]bool{}

	store.Walk(func(session *Session) {
		chk[session.user.Identity()] = true
	})

	var out = make([]uint64, 0)
	for ID := range chk {
		out = append(out, ID)
	}

	return out
}

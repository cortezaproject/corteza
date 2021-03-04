package request

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"net/http"
)

type (
	SessionManager struct {
		cstore store.AuthSessions
		sstore sessions.Store
		opt    options.AuthOpt
		log    *zap.Logger
	}
)

func NewSessionManager(store store.AuthSessions, opt options.AuthOpt, log *zap.Logger) *SessionManager {
	m := &SessionManager{opt: opt, log: log}
	m.cstore = store
	m.sstore = CortezaSessionStore(store, opt)
	return m
}

func (m SessionManager) Store() sessions.Store { return m.sstore }

func (m *SessionManager) Get(r *http.Request) *sessions.Session {
	ses, _ := m.sstore.Get(r, m.opt.SessionCookieName)
	return ses
}

func (m *SessionManager) Save(w http.ResponseWriter, r *http.Request) {
	if err := m.Get(r).Save(r, w); err != nil {
		m.log.Warn("failed to save sessions", zap.Error(err))
	}
}

// Returns all users sessions
func (m *SessionManager) Search(ctx context.Context, userID uint64) (set []*types.AuthSession, err error) {
	set, _, err = m.cstore.SearchAuthSessions(ctx, types.AuthSessionFilter{UserID: userID})
	for i := range set {
		set[i].Data = nil
	}
	return
}

// Returns all users sessions
func (m *SessionManager) DeleteByUserID(ctx context.Context, userID uint64) (err error) {
	return m.cstore.DeleteAuthSessionsByUserID(ctx, userID)
}

// Returns all users sessions
func (m *SessionManager) DeleteByID(ctx context.Context, sessionID string) (err error) {
	return m.cstore.DeleteAuthSessionByID(ctx, sessionID)
}

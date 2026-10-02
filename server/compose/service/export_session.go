package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/auth"
)

const (
	exportSessionTTL = time.Hour
)

type (
	RecordExportSession struct {
		SessionID uint64 `json:"sessionID,string"`
		UserID    uint64 `json:"userID,string"`

		NamespaceID uint64 `json:"namespaceID,string"`
		ModuleID    uint64 `json:"moduleID,string"`

		Filename string `json:"filename"`
		Ext      string `json:"ext"`

		Filter              string   `json:"filter"`
		Fields              []string `json:"fields"`
		Timezone            string   `json:"timezone"`
		MultiValueDelimiter string   `json:"multiValueDelimiter"`
		WrapMultiValue      string   `json:"wrapMultiValue"`
		ResolveRefs         bool     `json:"resolveRefs"`
		IncludeRefID        bool     `json:"includeRefID"`

		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	exportSessionSet []*RecordExportSession

	exportSession struct {
		l        sync.Mutex
		sessions exportSessionSet
	}

	ExportSessionService interface {
		Create(ctx context.Context, s *RecordExportSession) (*RecordExportSession, error)
		FindByID(ctx context.Context, sessionID uint64) (*RecordExportSession, error)
		DeleteByID(ctx context.Context, sessionID uint64) error
	}
)

func ExportSession() *exportSession {
	svc := &exportSession{
		sessions: exportSessionSet{},
	}

	go svc.cleanupLoop(context.Background())

	return svc
}

func (svc *exportSession) indexOf(userID, sessionID uint64) int {
	for i, r := range svc.sessions {
		if r.SessionID == sessionID && r.UserID == userID {
			return i
		}
	}
	return -1
}

func (svc *exportSession) Create(ctx context.Context, s *RecordExportSession) (*RecordExportSession, error) {
	svc.l.Lock()
	defer svc.l.Unlock()

	s.SessionID = nextID()
	s.UserID = auth.GetIdentityFromContext(ctx).Identity()
	s.CreatedAt = time.Now()
	s.UpdatedAt = s.CreatedAt

	svc.sessions = append(svc.sessions, s)
	return s, nil
}

func (svc *exportSession) FindByID(ctx context.Context, sessionID uint64) (*RecordExportSession, error) {
	svc.l.Lock()
	defer svc.l.Unlock()

	userID := auth.GetIdentityFromContext(ctx).Identity()
	i := svc.indexOf(userID, sessionID)
	if i < 0 {
		return nil, fmt.Errorf("compose.service.RecordExportSessionNotFound")
	}

	s := svc.sessions[i]
	if time.Since(s.UpdatedAt) > exportSessionTTL {
		svc.sessions = removeExportSession(svc.sessions, i)
		return nil, fmt.Errorf("compose.service.RecordExportSessionNotFound")
	}

	return s, nil
}

func (svc *exportSession) DeleteByID(ctx context.Context, sessionID uint64) error {
	svc.l.Lock()
	defer svc.l.Unlock()

	userID := auth.GetIdentityFromContext(ctx).Identity()
	i := svc.indexOf(userID, sessionID)
	if i >= 0 {
		svc.sessions = removeExportSession(svc.sessions, i)
	}
	return nil
}

func removeExportSession(s exportSessionSet, i int) exportSessionSet {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func (svc *exportSession) cleanupLoop(ctx context.Context) {
	t := time.NewTicker(10 * time.Minute)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			svc.clean()
		}
	}
}

func (svc *exportSession) clean() {
	svc.l.Lock()
	defer svc.l.Unlock()

	for i := len(svc.sessions) - 1; i >= 0; i-- {
		if time.Since(svc.sessions[i].UpdatedAt) > exportSessionTTL {
			svc.sessions = removeExportSession(svc.sessions, i)
		}
	}
}

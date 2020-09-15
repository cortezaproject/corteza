package service

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"sync"
	"time"
)

type (
	recordSet []*RecordImportSession

	importSession struct {
		l       sync.Mutex
		records recordSet
	}

	ImportSessionService interface {
		FindByID(ctx context.Context, sessionID uint64) (*RecordImportSession, error)
		SetByID(ctx context.Context, sessionID, namespaceID, moduleID uint64, fields map[string]string, progress *RecordImportProgress, decoder Decoder) (*RecordImportSession, error)
		DeleteByID(ctx context.Context, sessionID uint64) error
	}
)

func ImportSession() *importSession {
	return &importSession{
		records: recordSet{},
	}
}

func (svc importSession) indexOf(userID, sessionID uint64) int {
	for i, r := range svc.records {
		if r.SessionID == sessionID && r.UserID == userID {
			return i
		}
	}

	return -1
}

func (svc *importSession) FindByID(ctx context.Context, sessionID uint64) (*RecordImportSession, error) {
	svc.l.Lock()
	defer svc.l.Unlock()

	userID := auth.GetIdentityFromContext(ctx).Identity()
	i := svc.indexOf(userID, sessionID)
	if i >= 0 {
		return svc.records[i], nil
	}
	return nil, fmt.Errorf("compose.service.RecordImportSessionNotFound")
}

func (svc *importSession) SetByID(ctx context.Context, sessionID, namespaceID, moduleID uint64, fields map[string]string, progress *RecordImportProgress, decoder Decoder) (*RecordImportSession, error) {
	svc.l.Lock()
	defer svc.l.Unlock()

	userID := auth.GetIdentityFromContext(ctx).Identity()
	i := svc.indexOf(userID, sessionID)
	var ris *RecordImportSession

	if i >= 0 {
		ris = svc.records[i]
	} else {
		ris = &RecordImportSession{
			SessionID: nextID(),
			CreatedAt: time.Now(),
		}
		svc.records = append(svc.records, ris)
		ris.UserID = userID
	}
	ris.UpdatedAt = time.Now()

	if namespaceID > 0 {
		ris.NamespaceID = namespaceID
	}
	if moduleID > 0 {
		ris.ModuleID = moduleID
	}
	if fields != nil {
		ris.Fields = fields
	}
	if progress != nil {
		ris.Progress = *progress
	}

	if ris.Progress.FinishedAt != nil {
		ris.Decoder = nil
	} else if decoder != nil {
		ris.Decoder = decoder
	}

	return ris, nil
}

// https://stackoverflow.com/a/37335777
func remove(s recordSet, i int) recordSet {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func (svc *importSession) DeleteByID(ctx context.Context, sessionID uint64) error {
	svc.l.Lock()
	defer svc.l.Unlock()

	userID := auth.GetIdentityFromContext(ctx).Identity()
	i := svc.indexOf(userID, sessionID)

	if i >= 0 {
		svc.records = remove(svc.records, i)
	}
	return nil
}

// @todo run this in some interval
func (svc *importSession) clean(ctx context.Context) {
	svc.l.Lock()
	defer svc.l.Unlock()

	for i := len(svc.records) - 1; i >= 0; i-- {
		r := svc.records[i]
		if time.Now().After(r.UpdatedAt.Add(time.Hour * 24 * 3)) {
			svc.records = remove(svc.records, i)
		}
	}
}

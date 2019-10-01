package service

import (
	"context"
	"sync"
	"time"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/pkg/auth"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type (
	recordSet []*RecordImportSession

	importSession struct {
		l      sync.Mutex
		logger *zap.Logger

		records recordSet
	}

	ImportSessionService interface {
		FindRecordByID(ctx context.Context, sessionID uint64) (*RecordImportSession, error)
		SetRecordByID(ctx context.Context, sessionID, namespaceID, moduleID uint64, fields map[string]string, progress *RecordImportProgress, decoder Decoder) (*RecordImportSession, error)
		DeleteRecordByID(ctx context.Context, sessionID uint64) error
	}
)

func ImportSession() *importSession {
	return &importSession{
		logger:  DefaultLogger.Named("importSession"),
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

func (svc *importSession) FindRecordByID(ctx context.Context, sessionID uint64) (*RecordImportSession, error) {
	svc.l.Lock()
	defer svc.l.Unlock()

	userID := auth.GetIdentityFromContext(ctx).Identity()
	i := svc.indexOf(userID, sessionID)
	if i >= 0 {
		return svc.records[i], nil
	}
	return nil, errors.New("Can't access session: session not found")
}

func (svc *importSession) SetRecordByID(ctx context.Context, sessionID, namespaceID, moduleID uint64, fields map[string]string, progress *RecordImportProgress, decoder Decoder) (*RecordImportSession, error) {
	svc.l.Lock()
	defer svc.l.Unlock()

	userID := auth.GetIdentityFromContext(ctx).Identity()
	i := svc.indexOf(userID, sessionID)
	var ris *RecordImportSession

	if i >= 0 {
		ris = svc.records[i]
	} else {
		ris = &RecordImportSession{
			SessionID: factory.Sonyflake.NextID(),
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

func (svc *importSession) DeleteRecordByID(ctx context.Context, sessionID uint64) error {
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

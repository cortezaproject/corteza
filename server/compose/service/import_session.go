package service

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"

	envoyCsv "github.com/cortezaproject/corteza/server/pkg/envoyx/csv"
	envoyJson "github.com/cortezaproject/corteza/server/pkg/envoyx/json"
)

type (
	recordSet []*RecordImportSession

	importSession struct {
		l       sync.Mutex
		records recordSet
	}

	countableProvider interface {
		envoyx.Provider
		Count() uint64
		Fields() []string
	}

	ImportSessionService interface {
		Create(ctx context.Context, f io.ReadSeeker, name, contentType string, namespaceID, moduleID uint64) (*RecordImportSession, error)
		FindByID(ctx context.Context, sessionID uint64) (*RecordImportSession, error)
		DeleteByID(ctx context.Context, sessionID uint64) error
	}
)

func ImportSession() *importSession {
	return &importSession{
		records: recordSet{},
	}
}

func (svc *importSession) indexOf(userID, sessionID uint64) int {
	for i, r := range svc.records {
		if r.SessionID == sessionID && r.UserID == userID {
			return i
		}
	}

	return -1
}

func (svc *importSession) Create(ctx context.Context, f io.ReadSeeker, name, contentType string, namespaceID, moduleID uint64) (_ *RecordImportSession, err error) {
	svc.l.Lock()
	defer svc.l.Unlock()

	// Prepare the session
	sh := &RecordImportSession{
		Name:        name,
		SessionID:   nextID(),
		UserID:      auth.GetIdentityFromContext(ctx).Identity(),
		NamespaceID: namespaceID,
		ModuleID:    moduleID,

		OnError:  IMPORT_ON_ERROR_FAIL,
		Fields:   make(map[string]string),
		Progress: &RecordImportProgress{},

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rd, err := func() (pp countableProvider, err error) {
		if envoyCsv.CanDecodeFile(f) || envoyCsv.CanDecodeMime(contentType) {
			f.Seek(0, 0)
			return envoyCsv.Decoder(f, name)
		}

		if envoyJson.CanDecodeFile(f) || envoyJson.CanDecodeMime(contentType) {
			f.Seek(0, 0)
			return envoyJson.Decoder(f, name)
		}

		return nil, fmt.Errorf("compose.service.RecordImportFormatNotSupported")
	}()
	if err != nil {
		return
	}

	sh.Providers = append(sh.Providers, rd)

	prepKey := func(k string) string {
		return strings.TrimSpace(strings.ToLower(k))
	}

	sh.Progress.EntryCount = rd.Count()
	for _, f := range rd.Fields() {
		sh.Fields[f] = ""

		// @todo improve this bit
		k := prepKey(f)
		if k == "id" || k == "recordid" {
			sh.Key = f
		}
	}

	// Create it
	svc.records = append(svc.records, sh)
	return sh, nil
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

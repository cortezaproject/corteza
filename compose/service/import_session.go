package service

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/csv"
	"github.com/cortezaproject/corteza-server/pkg/envoy/json"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	recordSet []*recordImportSession

	importSession struct {
		l       sync.Mutex
		records recordSet
	}

	ImportSessionService interface {
		Create(ctx context.Context, f io.ReadSeeker, name, contentType string, namespaceID, moduleID uint64) (*recordImportSession, error)
		FindByID(ctx context.Context, sessionID uint64) (*recordImportSession, error)
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

func (svc *importSession) Create(ctx context.Context, f io.ReadSeeker, name, contentType string, namespaceID, moduleID uint64) (*recordImportSession, error) {
	svc.l.Lock()
	defer svc.l.Unlock()

	// Prepare the session
	sh := &recordImportSession{
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

	// Decoders; We only need to do csv & yaml here
	cd := csv.Decoder()
	jd := json.Decoder()

	// This will really be at most 1
	var err error
	do := &envoy.DecoderOpts{
		Name: name,
		Path: "",
	}

	sh.Resources, err = func() ([]resource.Interface, error) {
		if cd.CanDecodeFile(f) || cd.CanDecodeMime(contentType) {
			f.Seek(0, 0)
			return cd.Decode(ctx, f, do)
		}

		f.Seek(0, 0)
		if jd.CanDecodeFile(f) || jd.CanDecodeMime(contentType) {
			f.Seek(0, 0)
			return jd.Decode(ctx, f, do)
		}

		return nil, fmt.Errorf("compose.service.RecordImportFormatNotSupported")
	}()

	if err != nil {
		return nil, err
	}

	// Get some metadata
	n, ok := (sh.Resources[0]).(*resource.ResourceDataset)
	if !ok {
		// @todo move this logic to service and use action/error pattern
		return nil, fmt.Errorf("compose.service.RecordImportFormatNotSupported")
	}

	prepKey := func(k string) string {
		return strings.TrimSpace(strings.ToLower(k))
	}

	sh.Progress.EntryCount = n.P.Count()
	for _, f := range n.P.Fields() {
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

func (svc *importSession) FindByID(ctx context.Context, sessionID uint64) (*recordImportSession, error) {
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

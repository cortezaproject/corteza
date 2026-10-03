package compose

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/compose/dalutils"
	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/types"
)

// Saving a module (which replaces its DAL model) while records of that
// module are being read used to crash the server; the model registry must
// be safe to read and replace concurrently. Run with -race.
func TestRecordAccess_whileModelReplaced(t *testing.T) {
	h := newHelper(t)
	h.clearRecords()

	module := h.repoMakeRecordModuleWithFields("model replace module",
		&types.ModuleField{Name: "name", Kind: "String"},
	)
	for i := 0; i < 5; i++ {
		h.makeRecord(module, &types.RecordValue{Name: "name", Value: "r"})
	}
	ns := h.lookupNamespaceByID(module.NamespaceID)

	var (
		ctx      = context.Background()
		stop     = make(chan struct{})
		wg       sync.WaitGroup
		mux      sync.Mutex
		readErrs []error
	)

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}

				_, _, err := dalutils.ComposeRecordsList(ctx, defDal, module, types.RecordFilter{
					ModuleID:    module.ID,
					NamespaceID: module.NamespaceID,
				})
				if err != nil {
					mux.Lock()
					readErrs = append(readErrs, err)
					mux.Unlock()
				}
			}
		}()
	}

	deadline := time.Now().Add(400 * time.Millisecond)
	for time.Now().Before(deadline) {
		h.noError(service.DalModelReplace(ctx, service.DefaultStore, nil, defDal, ns, module))
	}

	close(stop)
	wg.Wait()

	// reads may briefly fail while the model is swapped, but nothing may panic
	// and reads must work again once the replacement is done
	_, _, err := dalutils.ComposeRecordsList(ctx, defDal, module, types.RecordFilter{ModuleID: module.ID, NamespaceID: module.NamespaceID})
	h.noError(err)
	t.Logf("reads that failed during replacement: %d", len(readErrs))
}

package workflows

import (
	"context"
	"fmt"
	"testing"
	"time"

	autTypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/compose/dalutils"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/stretchr/testify/require"
)

// Test_iterator_records_total iterates with incTotal over more records than
// the buffer holds; total is fetched with the first page and reported on every
// iteration
func Test_iterator_records_total(t *testing.T) {
	wfexec.MaxIteratorBufferSize = 2
	defer func() {
		wfexec.MaxIteratorBufferSize = wfexec.DefaultMaxIteratorBufferSize
	}()

	var (
		ctx = bypassRBAC(context.Background())
		req = require.New(t)
	)

	loadNewScenario(ctx, t)

	ns, err := store.LookupComposeNamespaceBySlug(ctx, defStore, "ns1")
	req.NoError(err)
	mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, defStore, ns.ID, "mod1")
	req.NoError(err)
	mod.Fields, _, err = store.SearchComposeModuleFields(ctx, defStore, types.ModuleFieldFilter{ModuleID: []uint64{mod.ID}})
	req.NoError(err)

	for j := 0; j < 5; j++ {
		rec := &types.Record{ID: id.Next(), NamespaceID: ns.ID, ModuleID: mod.ID, CreatedAt: time.Now()}
		rec.SetModule(mod)
		rec.Values = types.RecordValueSet{{Name: "f1", Value: fmt.Sprintf("%d", j+1), RecordID: rec.ID}}
		req.NoError(dalutils.ComposeRecordCreate(ctx, defDal, mod, rec))
	}

	var (
		_, trace = mustExecWorkflow(ctx, t, "testing", autTypes.WorkflowExecParams{})
	)

	// 6x iterator, 5x continue, 1x terminator, 1x completed
	req.Len(trace, 13)

	for j := 0; j <= 4; j++ {
		total, err := expr.Integer{}.Cast(trace[j*2].Results.GetValue()["t"])
		req.NoError(err)
		req.Equal(int64(5), total.Get().(int64), "iteration %d reports the total", j)
	}
}

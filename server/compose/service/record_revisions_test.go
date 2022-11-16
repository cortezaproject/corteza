package service

import (
	"context"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/revisions"
	"github.com/stretchr/testify/require"
	"testing"
)

type (
	mockRecordRevisionsDAL struct {
		search func(ctx context.Context, mf dal.ModelRef, f filter.Filter) (_ dal.Iterator, err error)
		create func(ctx context.Context, mf dal.ModelRef, revision *revisions.Revision) error
	}
)

func (svc *mockRecordRevisionsDAL) Search(ctx context.Context, mf dal.ModelRef, f filter.Filter) (_ dal.Iterator, err error) {
	return svc.search(ctx, mf, f)
}

func (svc *mockRecordRevisionsDAL) Create(ctx context.Context, mf dal.ModelRef, rev *revisions.Revision) error {
	return svc.create(ctx, mf, rev)
}

func TestRecordRevisions(t *testing.T) {
	var (
		req  = require.New(t)
		ctx  = context.Background()
		skip = types.ModuleFieldConfig{RecordRevisions: types.ModuleFieldConfigRecordRevisions{Skip: true}}
		mod  = &types.Module{
			ID:   1,
			Name: "test",
			Fields: []*types.ModuleField{
				{Name: "rev1", Kind: "string"},
				{Name: "rev2", Kind: "string"},
				{Name: "rev3", Kind: "string", Config: skip},
			},
		}

		rec = &types.Record{ID: 2, ModuleID: 1}

		dalMock = &mockRecordRevisionsDAL{}

		svc = &recordRevisions{r: dalMock}
	)

	rec.SetModule(mod)
	req.NoError(rec.SetValue("rev1", 0, "val1"))
	req.NoError(rec.SetValue("rev2", 0, "val2"))
	req.NoError(rec.SetValue("rev3", 0, "val3"))

	var changes int
	dalMock.create = func(ctx context.Context, mf dal.ModelRef, r *revisions.Revision) error {
		changes = len(r.Changes)
		return nil
	}

	req.NoError(svc.created(ctx, rec))
	req.Equal(3, changes, "expecting ownedBy, rev1 and rev2")
}

package gig

import (
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/gig"
	"github.com/cortezaproject/corteza-server/store"
)

func Test_envoy_worker_import(t *testing.T) {
	var (
		ctx, svc, h, s, g = setupWithImportGig(t)
		err               error
	)
	_ = s
	_ = svc
	_ = g

	t.Run("parsing", func(_ *testing.T) {
		g, err = svc.AddSources(ctx, g, gig.SourceWrapSet{{
			Uri:   testSource(t),
			IsDir: true,
		}})
		h.a.NoError(err)
	})

	t.Run("prepare", func(_ *testing.T) {
		g, err = svc.Prepare(ctx, g)
		h.a.NoError(err)

		state, err := g.Worker.State(ctx)
		h.a.NoError(err)

		ews := state.(gig.WorkerStateEnvoy)

		h.a.NotEqual(0, len(ews.Resources))
	})

	t.Run("exec", func(_ *testing.T) {
		cleanup(ctx, h, s)
		g, err = svc.Exec(ctx, g)
		h.a.NoError(err)

		nn, _, err := store.SearchComposeNamespaces(ctx, s, types.NamespaceFilter{})
		h.a.NoError(err)
		h.a.Len(nn, 1)
		h.a.Equal("ns1 name", nn[0].Name)

		mm, _, err := store.SearchComposeModules(ctx, s, types.ModuleFilter{NamespaceID: nn[0].ID})
		h.a.NoError(err)
		h.a.Len(mm, 1)
		h.a.Equal("mod1 name", mm[0].Name)

		ff, _, err := store.SearchComposeModuleFields(ctx, s, types.ModuleFieldFilter{ModuleID: []uint64{mm[0].ID}})
		h.a.NoError(err)
		h.a.Len(ff, 2)
		h.a.Equal("f1 label", ff[0].Label)
		h.a.Equal("f2 label", ff[1].Label)
		cleanup(ctx, h, s)
	})
}

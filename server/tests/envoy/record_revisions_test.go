package envoy

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/stretchr/testify/require"
)

func TestModuleWithRecordRevisionsRegistersBothModels(t *testing.T) {
	var (
		ctx = context.Background()
		req = require.New(t)
	)

	cleanup(t)

	nodes, providers, err := defaultEnvoy.Decode(ctx, envoyx.DecodeParams{
		Type:   envoyx.DecodeTypeURI,
		Params: map[string]any{"uri": "file://testdata/record_revisions"},
	})
	req.NoError(err)

	params := envoyx.EncodeParams{
		Type:   envoyx.EncodeTypeStore,
		Params: map[string]any{"storer": defaultStore, "dal": defaultDal},
	}
	gg, err := defaultEnvoy.Bake(ctx, params, providers, nodes...)
	req.NoError(err)
	req.NoError(defaultEnvoy.Encode(ctx, params, gg))

	ns, err := store.LookupComposeNamespaceBySlug(ctx, defaultStore, "test_ns_rev")
	req.NoError(err)
	mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, defaultStore, ns.ID, "test_ns_rev_mod")
	req.NoError(err)
	req.True(mod.Config.RecordRevisions.Enabled)

	models, err := defaultDal.SearchModels(ctx)
	req.NoError(err)

	var records, revisions int
	for _, m := range models {
		if m.Resource != mod.RbacResource() {
			continue
		}
		if m.ResourceID == mod.ID {
			records++
		} else {
			revisions++
		}
	}
	req.Equal(1, records, "the module's records model")
	req.Equal(1, revisions, "the module's revisions model")
}

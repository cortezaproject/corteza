package envoy

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestSettingImportExport(t *testing.T) {
	var (
		ctx   = context.Background()
		req   = require.New(t)
		nodes envoyx.NodeSet
		gg    *envoyx.DepGraph
		err   error
	)

	cleanup(t)

	// The test
	//
	// * imports some YAML files
	// * checks the DB state
	// * exports the DB into a YAML
	// * clears the DB
	// * imports the YAML
	// * checks the DB state the same way as before
	//
	// The above outlined flow allows us to trivially check if the data is both
	// imported and exported correctly.
	//
	// The initial step could also manually populate the DB but the YAML import
	// is more convenient.

	t.Run("initial import", func(t *testing.T) {
		t.Run("parse configs", func(t *testing.T) {
			nodes, _, err = defaultEnvoy.Decode(ctx, envoyx.DecodeParams{
				Type: envoyx.DecodeTypeURI,
				Params: map[string]any{
					"uri": "file://testdata/settings",
				},
			})
			req.NoError(err)
		})

		t.Run("bake", func(t *testing.T) {
			gg, err = defaultEnvoy.Bake(ctx, envoyx.EncodeParams{
				Type: envoyx.EncodeTypeStore,
				Params: map[string]any{
					"storer": defaultStore,
					"dal":    defaultDal,
				},
			}, nil, nodes...)
			req.NoError(err)
		})

		t.Run("import into DB", func(t *testing.T) {
			err = defaultEnvoy.Encode(ctx, envoyx.EncodeParams{
				Type: envoyx.EncodeTypeStore,
				Params: map[string]any{
					"storer": defaultStore,
					"dal":    defaultDal,
				},
			}, gg)
			req.NoError(err)
		})

		assertSettingState(ctx, t, defaultStore, req)
	})

	// Prepare a temp file where we'll dump the YAML into
	auxFile, err := os.CreateTemp(os.TempDir(), "*.yaml")
	req.NoError(err)
	spew.Dump(auxFile.Name())
	// defer os.Remove(auxFile.Name())
	defer auxFile.Close()

	var rules envoyx.NodeSet
	t.Run("export", func(t *testing.T) {
		t.Run("export from DB", func(t *testing.T) {
			nodes, _, err = defaultEnvoy.Decode(ctx, envoyx.DecodeParams{
				Type: envoyx.DecodeTypeStore,
				Params: map[string]any{
					"storer": defaultStore,
					"dal":    defaultDal,
				},
				Filter: map[string]envoyx.ResourceFilter{
					types.SettingValueResourceType: {},
				},
			})
			req.NoError(err)
		})

		var rr rbac.RuleSet
		t.Run("get all rules", func(t *testing.T) {
			rr, _, err = store.SearchRbacRules(ctx, defaultStore, rbac.RuleFilter{})
			req.NoError(err)
		})

		t.Run("connect rules to resources", func(t *testing.T) {
			rules, err = envoyx.RBACRulesForNodes(rr, nodes...)
			req.NoError(err)

			nodes = append(nodes, rules...)
		})

		t.Run("bake", func(t *testing.T) {
			gg, err = defaultEnvoy.Bake(ctx, envoyx.EncodeParams{
				Type: envoyx.EncodeTypeStore,
				Params: map[string]any{
					"storer": defaultStore,
					"dal":    defaultDal,
				},
			}, nil, nodes...)
			req.NoError(err)
		})

		t.Run("write file", func(t *testing.T) {
			err = defaultEnvoy.Encode(ctx, envoyx.EncodeParams{
				Type: envoyx.EncodeTypeIo,
				Params: map[string]any{
					"writer": auxFile,
				},
			}, gg)
			req.NoError(err)
		})
	})

	cleanup(t)

	t.Run("second import", func(t *testing.T) {
		t.Run("yaml parse", func(t *testing.T) {
			nodes, _, err = defaultEnvoy.Decode(ctx, envoyx.DecodeParams{
				Type: envoyx.DecodeTypeURI,
				Params: map[string]any{
					"uri": fmt.Sprintf("file://%s", auxFile.Name()),
				},
			})
			req.NoError(err)
		})

		t.Run("bake", func(t *testing.T) {
			gg, err = defaultEnvoy.Bake(ctx, envoyx.EncodeParams{
				Type: envoyx.EncodeTypeStore,
				Params: map[string]any{
					"storer": defaultStore,
					"dal":    defaultDal,
				},
			}, nil, nodes...)
			req.NoError(err)
		})

		t.Run("run import", func(t *testing.T) {
			err = defaultEnvoy.Encode(ctx, envoyx.EncodeParams{
				Type: envoyx.EncodeTypeStore,
				Params: map[string]any{
					"storer": defaultStore,
					"dal":    defaultDal,
				},
			}, gg)
			req.NoError(err)
		})

		assertSettingState(ctx, t, defaultStore, req)
	})
}

func assertSettingState(ctx context.Context, t *testing.T, s store.Storer, req *require.Assertions) {
	t.Run("check state", func(t *testing.T) {
		ss, _, err := store.SearchSettingValues(ctx, s, types.SettingsFilter{})
		req.NoError(err)

		var st *types.SettingValue
		st = ss.FindByName("auth.external.enabled")
		req.NotNil(st)

		st = ss.FindByName("auth.mail.from-address")
		req.NotNil(st)

		st = ss.FindByName("auth.internal.enabled")
		req.NotNil(st)

	})
}

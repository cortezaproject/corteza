package envoy

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/dalutils"
	"github.com/cortezaproject/corteza/server/compose/envoy"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestRecordsImportExport(t *testing.T) {
	var (
		ctx       = context.Background()
		req       = require.New(t)
		nodes     envoyx.NodeSet
		providers []envoyx.Provider
		gg        *envoyx.DepGraph
		err       error
	)
	_ = gg

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
			nodes, providers, err = defaultEnvoy.Decode(ctx, envoyx.DecodeParams{
				Type: envoyx.DecodeTypeURI,
				Params: map[string]any{
					"uri": "file://testdata/datasource_records",
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
			}, providers, nodes...)
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

		assertRecordState(ctx, t, defaultStore, defaultDal, req)
	})

	// Prepare a temp file where we'll dump the YAML into
	auxFile, err := os.CreateTemp(os.TempDir(), "*.csv")
	req.NoError(err)
	spew.Dump(auxFile.Name())
	// defer os.Remove(auxFile.Name())
	defer auxFile.Close()

	t.Run("export", func(t *testing.T) {
		t.Run("export from DB", func(t *testing.T) {
			nodes, _, err = defaultEnvoy.Decode(ctx, envoyx.DecodeParams{
				Type: envoyx.DecodeTypeStore,
				Params: map[string]any{
					"storer": defaultStore,
					"dal":    defaultDal,
				},
				Filter: map[string]envoyx.ResourceFilter{
					envoy.ComposeRecordDatasourceAuxType: {
						Refs: map[string]envoyx.Ref{
							"NamespaceID": {
								ResourceType: types.NamespaceResourceType,
								Identifiers:  envoyx.MakeIdentifiers("test_ns_1"),
								Scope: envoyx.Scope{
									ResourceType: types.NamespaceResourceType,
									Identifiers:  envoyx.MakeIdentifiers("test_ns_1")},
							},
							"ModuleID": {
								ResourceType: types.ModuleResourceType,
								Identifiers:  envoyx.MakeIdentifiers("test_ns_1_mod_1"),
								Scope: envoyx.Scope{
									ResourceType: types.NamespaceResourceType,
									Identifiers:  envoyx.MakeIdentifiers("test_ns_1"),
								},
							},
						},
						Scope: envoyx.Scope{
							ResourceType: types.NamespaceResourceType,
							Identifiers:  envoyx.MakeIdentifiers("test_ns_1"),
						},
					},
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
					"uri": "file://testdata/datasource_records",
				},
			})
			req.NoError(err)

			_, providers, err = defaultEnvoy.Decode(ctx, envoyx.DecodeParams{
				Type: envoyx.DecodeTypeURI,
				Params: map[string]any{
					"uri": fmt.Sprintf("file://%s", auxFile.Name()),
				},
			})
			req.NoError(err)
			for _, p := range providers {
				p.SetIdent("records.csv")
			}
		})

		t.Run("bake", func(t *testing.T) {
			gg, err = defaultEnvoy.Bake(ctx, envoyx.EncodeParams{
				Type: envoyx.EncodeTypeStore,
				Params: map[string]any{
					"storer": defaultStore,
					"dal":    defaultDal,
				},
			}, providers, nodes...)
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

		assertRecordState(ctx, t, defaultStore, defaultDal, req)
	})
}

func assertRecordState(ctx context.Context, t *testing.T, s store.Storer, dl dal.FullService, req *require.Assertions) {
	t.Run("check state", func(t *testing.T) {
		ns, err := store.LookupComposeNamespaceBySlug(ctx, defaultStore, "test_ns_1")
		req.NoError(err)

		mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, defaultStore, ns.ID, "test_ns_1_mod_1")
		req.NoError(err)

		rr, _, err := dalutils.ComposeRecordsList(ctx, dl, mod, types.RecordFilter{})
		req.NoError(err)

		req.Len(rr, 8)

		compareValues(req, rr[0].Values, types.RecordValueSet{{Name: "test_ns_1_mod_1_f1", Value: "row_1_src_c1"}, {Name: "test_ns_1_mod_1_f2", Value: "row_1_src_c2"}})
		compareValues(req, rr[1].Values, types.RecordValueSet{{Name: "test_ns_1_mod_1_f1", Value: "row_2_src_c1"}, {Name: "test_ns_1_mod_1_f2", Value: "row_2_src_c2"}})
		compareValues(req, rr[2].Values, types.RecordValueSet{{Name: "test_ns_1_mod_1_f1", Value: "row_3_src_c1"}, {Name: "test_ns_1_mod_1_f2", Value: "row_3_src_c2"}})
		compareValues(req, rr[3].Values, types.RecordValueSet{{Name: "test_ns_1_mod_1_f1", Value: "row_4_src_c1"}, {Name: "test_ns_1_mod_1_f2", Value: "row_4_src_c2"}})
		compareValues(req, rr[4].Values, types.RecordValueSet{{Name: "test_ns_1_mod_1_f1", Value: "row_5_src_c1"}, {Name: "test_ns_1_mod_1_f2", Value: "row_5_src_c2"}})
		compareValues(req, rr[5].Values, types.RecordValueSet{{Name: "test_ns_1_mod_1_f1", Value: "row_6_src_c1"}, {Name: "test_ns_1_mod_1_f2", Value: "row_6_src_c2"}})
		compareValues(req, rr[6].Values, types.RecordValueSet{{Name: "test_ns_1_mod_1_f1", Value: "row_7_src_c1"}, {Name: "test_ns_1_mod_1_f2", Value: "row_7_src_c2"}})
		compareValues(req, rr[7].Values, types.RecordValueSet{{Name: "test_ns_1_mod_1_f1", Value: "row_8_src_c1"}, {Name: "test_ns_1_mod_1_f2", Value: "row_8_src_c2"}})

	})
}

func compareValues(req *require.Assertions, a, b types.RecordValueSet) {
	for _, va := range a {
		vb := b.Get(va.Name, va.Place)
		req.NotNil(vb)

		req.Equal(va.Value, vb.Value)
	}

}

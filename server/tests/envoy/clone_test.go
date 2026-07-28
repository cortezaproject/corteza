package envoy

import (
	"context"
	"sort"
	"testing"

	composeTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/envoyx"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	"github.com/stretchr/testify/require"
)

// TestCloneNamespace_WithModuleField reproduces the namespace-clone step of
// the project-revision flow (system/service/project_revision.go
// CreateRevision -> compose/service/namespace.go
// Namespace.CloneFromStore -> Namespace.Clone -> Namespace.envoyRun) for a
// namespace that has a module with at least one field -- i.e. any project
// with a real data model.
//
// Root cause: CloneFromStore's decode Filter listed ModuleFieldResourceType
// as its own top-level scope, in addition to ModuleResourceType (which
// already yields each module's fields as nested children via
// extendedModuleDecoder in compose/envoy/store_decode.go, correctly scoped
// by ModuleID). makeModuleFieldFilter (compose/envoy/store_decode.go)
// ignores the scope it's given -- unlike its siblings extendModuleFilter /
// extendPageFilter / extendNamespaceFilter -- so that extra top-level entry
// decoded every module field in the store, unscoped. Every field that
// belonged to a module actually being cloned therefore got a duplicate
// node. Because the destination namespace doesn't exist yet when Prepare
// runs (all Prepare calls run before any Encode call, see
// pkg/envoyx/store.go encodeStore), matchupModuleFields never matches
// either copy against an existing row, so both got a fresh ID and both were
// upserted under the same (new) module -- tripping the (name, module_id)
// unique index on module_field ("not unique").
//
// This test drives the same StoreDecoder/StoreEncoder machinery
// CloneFromStore uses (decode with CloneFromStore's filter shape, then the
// same namespace-rename/reference-remap envoyRun performs, then
// Bake+Encode) against a namespace seeded with one module carrying one
// field. It lives here -- rather than in compose/service, where
// CloneFromStore itself is defined -- because compose/service's test
// suite currently panics on unrelated, pre-existing breakage (a nil
// eventbus dereference in chart_test.go) that aborts the whole test
// binary before any test in the package can report a result. This package
// already wires up a real StoreDecoder/StoreEncoder pair against an
// in-memory sqlite store with a working DAL (see main_test.go), which is
// exactly what full clone encoding needs (module encode registers a DAL
// model).
func TestCloneNamespace_WithModuleField(t *testing.T) {
	ctx := context.Background()
	req := require.New(t)
	cleanup(t)

	// Seed a source namespace with one module carrying one field -- the
	// minimal "real data model" that triggers the bug.
	srcNs := &composeTypes.Namespace{
		ID:      id.Next(),
		Name:    "clone-src",
		Slug:    "clone-src",
		Enabled: true,
	}
	req.NoError(store.CreateComposeNamespace(ctx, defaultStore, srcNs))

	mod := &composeTypes.Module{
		ID:          id.Next(),
		NamespaceID: srcNs.ID,
		Handle:      "lead",
		Name:        "Lead",
	}
	req.NoError(store.CreateComposeModule(ctx, defaultStore, mod))

	fld := &composeTypes.ModuleField{
		ID:          id.Next(),
		NamespaceID: srcNs.ID,
		ModuleID:    mod.ID,
		Kind:        "String",
		Name:        "email",
	}
	req.NoError(store.CreateComposeModuleField(ctx, defaultStore, fld))

	// --- Mirror CloneFromStore (compose/service/namespace.go) ---
	nsScope := envoyx.ResourceFilter{
		Scope: envoyx.Scope{
			ResourceType: composeTypes.NamespaceResourceType,
			Identifiers:  envoyx.MakeIdentifiers(srcNs.Slug, srcNs.ID),
		},
	}
	nn, _, err := defaultEnvoy.Decode(ctx, envoyx.DecodeParams{
		Type: envoyx.DecodeTypeStore,
		Params: map[string]any{
			"storer": defaultStore,
			"dal":    defaultDal,
		},
		Filter: map[string]envoyx.ResourceFilter{
			composeTypes.NamespaceResourceType:  {Identifiers: envoyx.MakeIdentifiers(srcNs.Slug, srcNs.ID)},
			composeTypes.ModuleResourceType:     nsScope,
			composeTypes.PageResourceType:       nsScope,
			composeTypes.PageLayoutResourceType: nsScope,
			composeTypes.ChartResourceType:      nsScope,
		},
	})
	req.NoError(err)

	// Sanity: exactly one ModuleField node must have been decoded for the
	// seeded field. If this ever fails again with 2, the decode is
	// producing duplicates again.
	fieldNodes := 0
	for _, n := range nn {
		if n.ResourceType == composeTypes.ModuleFieldResourceType {
			fieldNodes++
		}
	}
	req.Equal(1, fieldNodes, "decode must yield exactly one node per module field, not a duplicate")

	// --- Mirror envoyRun's namespace rename + reference remap (compose/service/namespace.go) ---
	oldRef := envoyx.Ref{
		ResourceType: composeTypes.NamespaceResourceType,
		Identifiers:  envoyx.MakeIdentifiers(srcNs.Slug, srcNs.ID),
		Scope: envoyx.Scope{
			ResourceType: composeTypes.NamespaceResourceType,
			Identifiers:  envoyx.MakeIdentifiers(srcNs.Slug, srcNs.ID),
		},
	}
	nsNode := envoyx.NodeForRef(oldRef, nn...)
	req.NotNil(nsNode)
	newNs := nsNode.Resource.(*composeTypes.Namespace)
	newNs.ID = 0
	newNs.Name = "clone-src (revision 1)"
	newNs.Slug = "clone-src-rev1"

	nsNode.Identifiers = envoyx.MakeIdentifiers(newNs.Slug)
	nsNode.Scope.Identifiers = nsNode.Identifiers

	for _, n := range nn {
		nr := make(map[string]envoyx.Ref)
		for k, r := range n.References {
			if r.ResourceType == composeTypes.NamespaceResourceType {
				r.Identifiers = nsNode.Identifiers
			}
			if r.Scope.ResourceType == composeTypes.NamespaceResourceType {
				r.Scope = nsNode.Scope
			}
			nr[k] = r
		}
		n.References = nr
		if n.Scope.ResourceType == nsNode.Scope.ResourceType {
			n.Scope = nsNode.Scope
		}
	}

	gg, err := defaultEnvoy.Bake(ctx, envoyx.EncodeParams{
		Type: envoyx.EncodeTypeStore,
		Params: map[string]any{
			"storer": defaultStore,
			"dal":    defaultDal,
		},
	}, nil, nn...)
	req.NoError(err)

	err = defaultEnvoy.Encode(ctx, envoyx.EncodeParams{
		Type: envoyx.EncodeTypeStore,
		Params: map[string]any{
			"storer": defaultStore,
			"dal":    defaultDal,
		},
	}, gg)
	req.NoError(err, "clone must not fail with a ModuleField uniqueness error")

	// --- Assertions on the resulting DB state ---
	clonedNs, err := store.LookupComposeNamespaceBySlug(ctx, defaultStore, newNs.Slug)
	req.NoError(err)
	req.NotNil(clonedNs)
	req.NotEqual(srcNs.ID, clonedNs.ID)

	clonedMods, _, err := store.SearchComposeModules(ctx, defaultStore, composeTypes.ModuleFilter{NamespaceID: clonedNs.ID})
	req.NoError(err)
	req.Len(clonedMods, 1)
	req.NotEqual(mod.ID, clonedMods[0].ID, "cloned module must get a new ID, not reuse the source module's row")

	clonedFields, _, err := store.SearchComposeModuleFields(ctx, defaultStore, composeTypes.ModuleFieldFilter{ModuleID: []uint64{clonedMods[0].ID}})
	req.NoError(err)
	req.Len(clonedFields, 1, "cloned module must have exactly one field, not a duplicate")
	req.Equal(fld.Name, clonedFields[0].Name)
	req.NotEqual(fld.ID, clonedFields[0].ID, "cloned field must get a new ID, not reuse the source field's row")

	// The source module/field must be untouched by the clone.
	srcFields, _, err := store.SearchComposeModuleFields(ctx, defaultStore, composeTypes.ModuleFieldFilter{ModuleID: []uint64{mod.ID}})
	req.NoError(err)
	req.Len(srcFields, 1)
	req.Equal(fld.ID, srcFields[0].ID)
}

// cloneNamespaceViaStore drives the same StoreDecoder/StoreEncoder machinery
// CloneFromStore (compose/service/namespace.go) uses -- decode with
// CloneFromStore's filter shape, the namespace-rename/reference-remap
// envoyRun performs, then Bake+Encode -- and returns the resulting cloned
// namespace. Shared by the field-order regression tests below; see
// TestCloneNamespace_WithModuleField above for why this lives in tests/envoy
// rather than compose/service.
func cloneNamespaceViaStore(t *testing.T, req *require.Assertions, srcNs *composeTypes.Namespace, newSlug, newName string) *composeTypes.Namespace {
	ctx := context.Background()

	nsScope := envoyx.ResourceFilter{
		Scope: envoyx.Scope{
			ResourceType: composeTypes.NamespaceResourceType,
			Identifiers:  envoyx.MakeIdentifiers(srcNs.Slug, srcNs.ID),
		},
	}
	nn, _, err := defaultEnvoy.Decode(ctx, envoyx.DecodeParams{
		Type: envoyx.DecodeTypeStore,
		Params: map[string]any{
			"storer": defaultStore,
			"dal":    defaultDal,
		},
		Filter: map[string]envoyx.ResourceFilter{
			composeTypes.NamespaceResourceType:  {Identifiers: envoyx.MakeIdentifiers(srcNs.Slug, srcNs.ID)},
			composeTypes.ModuleResourceType:     nsScope,
			composeTypes.PageResourceType:       nsScope,
			composeTypes.PageLayoutResourceType: nsScope,
			composeTypes.ChartResourceType:      nsScope,
		},
	})
	req.NoError(err)

	oldRef := envoyx.Ref{
		ResourceType: composeTypes.NamespaceResourceType,
		Identifiers:  envoyx.MakeIdentifiers(srcNs.Slug, srcNs.ID),
		Scope: envoyx.Scope{
			ResourceType: composeTypes.NamespaceResourceType,
			Identifiers:  envoyx.MakeIdentifiers(srcNs.Slug, srcNs.ID),
		},
	}
	nsNode := envoyx.NodeForRef(oldRef, nn...)
	req.NotNil(nsNode)
	newNs := nsNode.Resource.(*composeTypes.Namespace)
	newNs.ID = 0
	newNs.Name = newName
	newNs.Slug = newSlug

	nsNode.Identifiers = envoyx.MakeIdentifiers(newNs.Slug)
	nsNode.Scope.Identifiers = nsNode.Identifiers

	for _, n := range nn {
		nr := make(map[string]envoyx.Ref)
		for k, r := range n.References {
			if r.ResourceType == composeTypes.NamespaceResourceType {
				r.Identifiers = nsNode.Identifiers
			}
			if r.Scope.ResourceType == composeTypes.NamespaceResourceType {
				r.Scope = nsNode.Scope
			}
			nr[k] = r
		}
		n.References = nr
		if n.Scope.ResourceType == nsNode.Scope.ResourceType {
			n.Scope = nsNode.Scope
		}
	}

	gg, err := defaultEnvoy.Bake(ctx, envoyx.EncodeParams{
		Type: envoyx.EncodeTypeStore,
		Params: map[string]any{
			"storer": defaultStore,
			"dal":    defaultDal,
		},
	}, nil, nn...)
	req.NoError(err)

	err = defaultEnvoy.Encode(ctx, envoyx.EncodeParams{
		Type: envoyx.EncodeTypeStore,
		Params: map[string]any{
			"storer": defaultStore,
			"dal":    defaultDal,
		},
	}, gg)
	req.NoError(err)

	clonedNs, err := store.LookupComposeNamespaceBySlug(ctx, defaultStore, newNs.Slug)
	req.NoError(err)
	req.NotNil(clonedNs)

	return clonedNs
}

// TestCloneNamespace_FieldOrder_ExplicitPlace reproduces the field-order
// defect for a module whose fields carry real, distinct Place values (e.g.
// set by updateModuleFields on every edit -- compose/service/module.go).
//
// Root cause: decodeModuleField (compose/envoy/store_decode.gen.go, backed by
// store.SearchComposeModuleFields) has no ORDER BY, so extendedModuleDecoder
// (compose/envoy/store_decode.go) received -- and handed on to the encoder --
// fields in whatever order the store happened to return them, which need not
// match Place order. Each field's own Place value round-tripped into the
// clone unchanged (clone does carry Place across correctly), but nothing
// re-derives a fresh, meaningful Place sequence for the destination, so the
// destination's raw row order -- which a later unsorted read can surface --
// need not match the source, even when the source's Place values are
// perfectly good.
//
// This seeds the fields in the DB in a scrambled order relative to their
// intended (Place) sequence -- value, stage, closeDate, title, with Place
// 1, 2, 3, 0 respectively, i.e. exactly the reordering reported live -- and
// asserts the clone reproduces the source's Place-sorted order (title,
// value, stage, closeDate) with a clean, distinct Place sequence.
func TestCloneNamespace_FieldOrder_ExplicitPlace(t *testing.T) {
	ctx := context.Background()
	req := require.New(t)
	cleanup(t)

	srcNs := &composeTypes.Namespace{
		ID:      id.Next(),
		Name:    "clone-src-order-a",
		Slug:    "clone-src-order-a",
		Enabled: true,
	}
	req.NoError(store.CreateComposeNamespace(ctx, defaultStore, srcNs))

	mod := &composeTypes.Module{
		ID:          id.Next(),
		NamespaceID: srcNs.ID,
		Handle:      "deal",
		Name:        "Deal",
	}
	req.NoError(store.CreateComposeModule(ctx, defaultStore, mod))

	// Inserted (and, on sqlite, physically stored) in this order, but Place
	// says the intended order is title(0), value(1), stage(2), closeDate(3).
	type seed struct {
		name  string
		place int
	}
	for _, s := range []seed{
		{"value", 1},
		{"stage", 2},
		{"closeDate", 3},
		{"title", 0},
	} {
		f := &composeTypes.ModuleField{
			ID:          id.Next(),
			NamespaceID: srcNs.ID,
			ModuleID:    mod.ID,
			Kind:        "String",
			Name:        s.name,
			Place:       s.place,
		}
		req.NoError(store.CreateComposeModuleField(ctx, defaultStore, f))
	}

	clonedNs := cloneNamespaceViaStore(t, req, srcNs, "clone-src-order-a-rev1", "clone-src-order-a (revision 1)")

	clonedMods, _, err := store.SearchComposeModules(ctx, defaultStore, composeTypes.ModuleFilter{NamespaceID: clonedNs.ID})
	req.NoError(err)
	req.Len(clonedMods, 1)

	clonedFields, _, err := store.SearchComposeModuleFields(ctx, defaultStore, composeTypes.ModuleFieldFilter{ModuleID: []uint64{clonedMods[0].ID}})
	req.NoError(err)
	req.Len(clonedFields, 4)

	sort.Sort(clonedFields)
	gotNames := make([]string, len(clonedFields))
	gotPlaces := make([]int, len(clonedFields))
	for i, f := range clonedFields {
		gotNames[i] = f.Name
		gotPlaces[i] = f.Place
	}
	req.Equal([]string{"title", "value", "stage", "closeDate"}, gotNames, "cloned fields must sort into the source's intended order")
	req.Equal([]int{0, 1, 2, 3}, gotPlaces, "cloned fields must get a clean, distinct Place sequence")
}

// TestCloneNamespace_FieldOrder_DegeneratePlace covers the more common real
// case: a module created with all of its fields in a single createModule
// call (compose/service/module.go) never has Place assigned at all -- unlike
// updateModuleFields, createModule has no `f.Place = idx` step, so every
// field lands at Place=0. A raw read of such a module only "looks" ordered
// because the source table has never been rewritten since that initial
// insert; cloning it is a fresh batch of writes with no such accidental
// guarantee. This seeds fields exactly that way (Place=0 across the board,
// inserted in the intended display order) and asserts the clone still
// reproduces that order with a newly meaningful, distinct Place sequence --
// i.e. the destination is no longer relying on accidental row order at all.
func TestCloneNamespace_FieldOrder_DegeneratePlace(t *testing.T) {
	ctx := context.Background()
	req := require.New(t)
	cleanup(t)

	srcNs := &composeTypes.Namespace{
		ID:      id.Next(),
		Name:    "clone-src-order-b",
		Slug:    "clone-src-order-b",
		Enabled: true,
	}
	req.NoError(store.CreateComposeNamespace(ctx, defaultStore, srcNs))

	mod := &composeTypes.Module{
		ID:          id.Next(),
		NamespaceID: srcNs.ID,
		Handle:      "deal",
		Name:        "Deal",
	}
	req.NoError(store.CreateComposeModule(ctx, defaultStore, mod))

	for _, name := range []string{"title", "value", "stage", "closeDate"} {
		f := &composeTypes.ModuleField{
			ID:          id.Next(),
			NamespaceID: srcNs.ID,
			ModuleID:    mod.ID,
			Kind:        "String",
			Name:        name,
			Place:       0,
		}
		req.NoError(store.CreateComposeModuleField(ctx, defaultStore, f))
	}

	clonedNs := cloneNamespaceViaStore(t, req, srcNs, "clone-src-order-b-rev1", "clone-src-order-b (revision 1)")

	clonedMods, _, err := store.SearchComposeModules(ctx, defaultStore, composeTypes.ModuleFilter{NamespaceID: clonedNs.ID})
	req.NoError(err)
	req.Len(clonedMods, 1)

	clonedFields, _, err := store.SearchComposeModuleFields(ctx, defaultStore, composeTypes.ModuleFieldFilter{ModuleID: []uint64{clonedMods[0].ID}})
	req.NoError(err)
	req.Len(clonedFields, 4)

	placesSeen := make(map[int]bool, len(clonedFields))
	for _, f := range clonedFields {
		req.False(placesSeen[f.Place], "cloned fields must not share a Place value")
		placesSeen[f.Place] = true
	}

	sort.Sort(clonedFields)
	gotNames := make([]string, len(clonedFields))
	for i, f := range clonedFields {
		gotNames[i] = f.Name
	}
	req.Equal([]string{"title", "value", "stage", "closeDate"}, gotNames, "cloned fields must sort into the source's insertion order once Place is no longer degenerate")
}

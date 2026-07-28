package envoy

import (
	"context"
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

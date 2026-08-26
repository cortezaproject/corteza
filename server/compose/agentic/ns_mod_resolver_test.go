package agentic

import (
	"context"
	"fmt"
	"testing"

	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/stretchr/testify/require"
)

// Only the two methods resolution uses are stubbed; the embedded interface
// makes reaching for any other one a panic rather than a silent pass.
type stubNamespaceSvc struct {
	cmpService.NamespaceService
	exact func(any) (*cmpTypes.Namespace, error)
	found cmpTypes.NamespaceSet
}

func (s stubNamespaceSvc) FindByAny(_ context.Context, id any) (*cmpTypes.Namespace, error) {
	if s.exact == nil {
		return nil, fmt.Errorf("not found")
	}
	return s.exact(id)
}

func (s stubNamespaceSvc) Search(_ context.Context, f cmpTypes.NamespaceFilter) (cmpTypes.NamespaceSet, cmpTypes.NamespaceFilter, error) {
	if f.Query == "" {
		return s.found, f, nil
	}

	out := cmpTypes.NamespaceSet{}
	for _, ns := range s.found {
		if contains(ns.Slug, f.Query) || contains(ns.Name, f.Query) {
			out = append(out, ns)
		}
	}
	return out, f, nil
}

type stubModuleSvc struct {
	cmpService.ModuleService
	found cmpTypes.ModuleSet
}

func (s stubModuleSvc) FindByAny(_ context.Context, _ uint64, _ any) (*cmpTypes.Module, error) {
	return nil, fmt.Errorf("not found")
}

func (s stubModuleSvc) Search(_ context.Context, f cmpTypes.ModuleFilter) (cmpTypes.ModuleSet, cmpTypes.ModuleFilter, error) {
	if f.Query == "" {
		return s.found, f, nil
	}

	out := cmpTypes.ModuleSet{}
	for _, m := range s.found {
		if contains(m.Handle, f.Query) || contains(m.Name, f.Query) {
			out = append(out, m)
		}
	}
	return out, f, nil
}

func contains(hay, needle string) bool {
	return len(needle) > 0 && len(hay) >= len(needle) && indexFold(hay, needle) >= 0
}

func indexFold(hay, needle string) int {
	lower := func(s string) string {
		b := []byte(s)
		for i := range b {
			if b[i] >= 'A' && b[i] <= 'Z' {
				b[i] += 'a' - 'A'
			}
		}
		return string(b)
	}
	h, n := lower(hay), lower(needle)
	for i := 0; i+len(n) <= len(h); i++ {
		if h[i:i+len(n)] == n {
			return i
		}
	}
	return -1
}

func withNamespaces(t *testing.T, svc cmpService.NamespaceService) {
	t.Helper()
	prev := cmpService.DefaultNamespace
	cmpService.DefaultNamespace = svc
	t.Cleanup(func() { cmpService.DefaultNamespace = prev })
}

func withModules(t *testing.T, svc cmpService.ModuleService) {
	t.Helper()
	prev := cmpService.DefaultModule
	cmpService.DefaultModule = svc
	t.Cleanup(func() { cmpService.DefaultModule = prev })
}

var namespaces = cmpTypes.NamespaceSet{
	{ID: 1, Slug: "mtg_collection", Name: "MTG Collection"},
	{ID: 2, Slug: "demo_crm", Name: "Demo CRM"},
	{ID: 3, Slug: "", Name: "Validation"},
}

// TL;DR: a namespace named the way a person says it out loud resolves.
// Example: a model writes what the user said — "my mtg namespace" — and an
// exact match on handle, slug or name does not find "MTG Collection". This runs
// before the tool does, so the call died here and compose_namespace_lookup
// never reached the listing it falls back to. The model was handed an error
// with no route forward and invented five namespaces that did not exist.
func TestResolveNamespace(t *testing.T) {
	t.Run("an exact name still wins without searching", func(t *testing.T) {
		withNamespaces(t, stubNamespaceSvc{
			exact: func(any) (*cmpTypes.Namespace, error) { return namespaces[1], nil },
		})

		ns, err := resolveNamespaceRef(t.Context(), "demo_crm")
		require.NoError(t, err)
		require.Equal(t, uint64(2), ns.ID)
	})

	t.Run("one partial match is the namespace that was meant", func(t *testing.T) {
		withNamespaces(t, stubNamespaceSvc{found: namespaces})

		for _, ref := range []string{"mtg", "MTG", "collection", "MTG Collection"} {
			ns, err := resolveNamespaceRef(t.Context(), ref)
			require.NoError(t, err, ref)
			require.Equal(t, uint64(1), ns.ID, ref)
		}
	})

	// Picking either would be a guess about which data the agent reads.
	t.Run("two partial matches is a question, not a guess", func(t *testing.T) {
		withNamespaces(t, stubNamespaceSvc{found: cmpTypes.NamespaceSet{
			{ID: 1, Slug: "mtg_cards", Name: "MTG Cards"},
			{ID: 2, Slug: "mtg_decks", Name: "MTG Decks"},
		}})

		_, err := resolveNamespaceRef(t.Context(), "mtg")
		require.Error(t, err)
		require.Contains(t, err.Error(), "matches 2")
		require.Contains(t, err.Error(), "mtg_cards")
		require.Contains(t, err.Error(), "mtg_decks")
	})

	// An error with no route forward is what the model filled in by inventing.
	t.Run("no match names what does exist", func(t *testing.T) {
		withNamespaces(t, stubNamespaceSvc{found: namespaces})

		_, err := resolveNamespaceRef(t.Context(), "magic")
		require.Error(t, err)
		require.Contains(t, err.Error(), `"magic" not found`)
		require.Contains(t, err.Error(), "mtg_collection")
		require.Contains(t, err.Error(), "demo_crm")
		require.Contains(t, err.Error(), "Validation", "a namespace with no slug is still named")
	})
}

func TestResolveModule(t *testing.T) {
	ns := &cmpTypes.Namespace{ID: 1, Slug: "mtg_collection"}
	modules := cmpTypes.ModuleSet{
		{ID: 10, Handle: "card", Name: "Card"},
		{ID: 11, Handle: "collection_item", Name: "Collection Item"},
		{ID: 12, Handle: "wishlist_item", Name: "Wishlist Item"},
	}

	t.Run("one partial match resolves", func(t *testing.T) {
		withModules(t, stubModuleSvc{found: modules})

		mod, err := resolveModuleRef(t.Context(), ns, "collection")
		require.NoError(t, err)
		require.Equal(t, uint64(11), mod.ID)
	})

	t.Run("an ambiguous one names both", func(t *testing.T) {
		withModules(t, stubModuleSvc{found: modules})

		_, err := resolveModuleRef(t.Context(), ns, "item")
		require.Error(t, err)
		require.Contains(t, err.Error(), "matches 2 modules in mtg_collection")
		require.Contains(t, err.Error(), "collection_item")
		require.Contains(t, err.Error(), "wishlist_item")
	})

	t.Run("no match names the modules that exist", func(t *testing.T) {
		withModules(t, stubModuleSvc{found: modules})

		_, err := resolveModuleRef(t.Context(), ns, "nonsense")
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found in mtg_collection")
		require.Contains(t, err.Error(), "card")
	})
}

// The handle is what every other argument wants back; the name is what the
// person recognises. A namespace with neither should not read as ", ,".
func TestDescribe(t *testing.T) {
	require.Equal(t, "mtg_collection (MTG Collection)", describe("mtg_collection", "MTG Collection"))
	require.Equal(t, "Validation", describe("", "Validation"))
	require.Equal(t, "card", describe("card", "card"))
	require.Equal(t, "card", describe("card", "Card"), "same word, different case, is not worth repeating")
	require.Equal(t, "card", describe("card", ""))
}

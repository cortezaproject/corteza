package agentic

import (
	"context"
	"fmt"
	"strings"

	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	agenticRuntime "github.com/crusttech/human/server/system/agentic/runtime"
)

// looseMatchLimit caps the search a partial name runs. A term matching more
// than a handful is not a name, and the answer is the same either way: say
// which ones and ask.
const looseMatchLimit = 25

type nsModResolver struct{}

func NsModResolver() *nsModResolver {
	return &nsModResolver{}
}

func (r *nsModResolver) Resolve(ctx context.Context, namespace, module string) (nsID, modID uint64, err error) {
	ns, err := resolveNamespaceRef(ctx, namespace)
	if err != nil {
		return 0, 0, err
	}

	if module == "" {
		return ns.ID, 0, nil
	}

	mod, err := resolveModuleRef(ctx, ns, module)
	if err != nil {
		return 0, 0, err
	}

	return ns.ID, mod.ID, nil
}

// resolveNamespaceRef answers to the name a person would use out loud.
//
// A model writes what the person said — "my mtg namespace" — and an exact match
// on handle, slug or name does not find "MTG Collection". Failing there used to
// end the call, because this runs before the tool does: compose_namespace_lookup
// never reached the listing it falls back to, and the model was left with an
// error and no way forward.
//
// One unambiguous partial match is the namespace that was meant. Two is a
// question, and the error asks it rather than guessing. Search is bounded by
// what the invoking user may read, so this widens what can be named, never what
// can be reached.
func resolveNamespaceRef(ctx context.Context, ref string) (*cmpTypes.Namespace, error) {
	if ns, err := cmpService.DefaultNamespace.FindByAny(ctx, ref); err == nil {
		return ns, nil
	}

	set, err := findNamespaces(ctx, ref)
	if err != nil {
		return nil, fmt.Errorf("namespace lookup failed: %w", err)
	}

	switch len(set) {
	case 1:
		return set[0], nil

	case 0:
		all, err := findNamespaces(ctx, "")
		if err != nil || len(all) == 0 {
			return nil, fmt.Errorf("namespace %q not found", ref)
		}
		return nil, fmt.Errorf("namespace %q not found. These exist: %s", ref, namespaceNames(all))

	default:
		return nil, fmt.Errorf("namespace %q matches %d of them: %s. Name one exactly", ref, len(set), namespaceNames(set))
	}
}

// resolveModuleRef is resolveNamespaceRef within one namespace: "card" naming both
// 'card' and 'deck_card' is the ordinary case, and picking either would be a
// guess about which records the agent reads.
func resolveModuleRef(ctx context.Context, ns *cmpTypes.Namespace, ref string) (*cmpTypes.Module, error) {
	if mod, err := cmpService.DefaultModule.FindByAny(ctx, ns.ID, ref); err == nil {
		return mod, nil
	}

	set, err := findModules(ctx, ns.ID, ref)
	if err != nil {
		return nil, fmt.Errorf("module lookup failed: %w", err)
	}

	switch len(set) {
	case 1:
		return set[0], nil

	case 0:
		all, err := findModules(ctx, ns.ID, "")
		if err != nil || len(all) == 0 {
			return nil, fmt.Errorf("module %q not found in %s", ref, ns.Slug)
		}
		return nil, fmt.Errorf("module %q not found in %s. These exist: %s", ref, ns.Slug, moduleNames(all))

	default:
		return nil, fmt.Errorf("module %q matches %d modules in %s: %s. Name one exactly", ref, len(set), ns.Slug, moduleNames(set))
	}
}

// resolveModuleArg is resolveModuleRef for a raw tool argument, which arrives as
// any and may be an ID rather than a name.
func resolveModuleArg(ctx context.Context, ns *cmpTypes.Namespace, v any) (*cmpTypes.Module, error) {
	ref, ok := v.(string)
	if !ok {
		mod, err := cmpService.DefaultModule.FindByAny(ctx, ns.ID, v)
		if err != nil {
			return nil, toolkit.Errf("module lookup", err)
		}
		return mod, nil
	}

	return resolveModuleRef(ctx, ns, ref)
}

// Search, not Find: it is the entry point the handlers use, and it is the one
// that carries the access checks that keep this to what the caller may read.
func findNamespaces(ctx context.Context, query string) (cmpTypes.NamespaceSet, error) {
	f := cmpTypes.NamespaceFilter{Query: query}

	var err error
	if f.Paging, err = filter.NewPaging(looseMatchLimit, ""); err != nil {
		return nil, err
	}

	set, _, err := cmpService.DefaultNamespace.Search(ctx, f)
	return set, err
}

func findModules(ctx context.Context, nsID uint64, query string) (cmpTypes.ModuleSet, error) {
	f := cmpTypes.ModuleFilter{NamespaceID: nsID, Query: query}

	var err error
	if f.Paging, err = filter.NewPaging(looseMatchLimit, ""); err != nil {
		return nil, err
	}

	set, _, err := cmpService.DefaultModule.Search(ctx, f)
	return set, err
}

// The handle is what every other argument wants back, so it leads; the name is
// what the person recognises.
func namespaceNames(set cmpTypes.NamespaceSet) string {
	out := make([]string, 0, len(set))
	for _, ns := range set {
		out = append(out, describe(ns.Slug, ns.Name))
	}
	return strings.Join(out, ", ")
}

func moduleNames(set cmpTypes.ModuleSet) string {
	out := make([]string, 0, len(set))
	for _, mod := range set {
		out = append(out, describe(mod.Handle, mod.Name))
	}
	return strings.Join(out, ", ")
}

func describe(handle, name string) string {
	if handle == "" {
		return name
	}
	if name == "" || strings.EqualFold(handle, name) {
		return handle
	}
	return fmt.Sprintf("%s (%s)", handle, name)
}

func (r *nsModResolver) LookupNamespace(ctx context.Context, id uint64) (agenticRuntime.NsHandle, error) {
	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, id)
	if err != nil {
		return agenticRuntime.NsHandle{}, err
	}
	handle := ns.Slug
	if handle == "" {
		handle = ns.Name
	}
	return agenticRuntime.NsHandle{ID: ns.ID, Handle: handle, Name: ns.Name, CreatedByAgent: ns.CreatedByAgent}, nil
}

func (r *nsModResolver) LookupModule(ctx context.Context, nsID, modID uint64) (agenticRuntime.ModHandle, error) {
	mod, err := cmpService.DefaultModule.FindByAny(ctx, nsID, modID)
	if err != nil {
		return agenticRuntime.ModHandle{}, err
	}
	modHandle := mod.Handle
	if modHandle == "" {
		modHandle = mod.Name
	}
	return agenticRuntime.ModHandle{ID: mod.ID, NamespaceID: mod.NamespaceID, Handle: modHandle, Name: mod.Name, CreatedByAgent: mod.CreatedByAgent}, nil
}

// lookupNamespaceArg resolves the "namespace" argument, naming it when it is
// simply absent.
//
// FindByAny is given the raw argument, so a missing one arrives as nil and
// comes back as "invalid ID" — which reads as though a namespace was supplied
// and rejected. A model told that supplies an ID it invented, rather than the
// namespace it forgot.
func lookupNamespaceArg(ctx context.Context, args map[string]any) (*cmpTypes.Namespace, error) {
	if v, ok := args["namespace"]; !ok || v == nil || v == "" {
		return nil, fmt.Errorf("namespace is required — its handle, slug, name, or ID as a string")
	}

	ref, ok := args["namespace"].(string)
	if !ok {
		ns, err := cmpService.DefaultNamespace.FindByAny(ctx, args["namespace"])
		if err != nil {
			return nil, toolkit.Errf("namespace lookup", err)
		}
		return ns, nil
	}

	// One rule for naming a namespace, whether the runtime resolves it before
	// the tool or the tool resolves it itself.
	return resolveNamespaceRef(ctx, ref)
}
